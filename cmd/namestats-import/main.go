// This is the namestats-import command.  See README.md for details.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
	"github.com/shutej/namestats/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	from = flag.Int("from", 1880, "first year to include")
	to   = flag.Int("to", time.Now().Year()-1, "last year to include")
	uri  = flag.String("uri", "", "URI to connect to MongoDB")
)

type PopularNames struct{}

func (_ PopularNames) valuesForYear(year int) url.Values {
	values := url.Values{}
	values.Add("year", strconv.Itoa(year))
	values.Add("top", "1000")
	values.Add("number", "n")
	return values
}

const ssaNamesUrl = "http://www.ssa.gov/cgi-bin/popularnames.cgi"

var (
	trMatcher = cascadia.MustCompile("table table tr")
	thMatcher = cascadia.MustCompile("td:not([colspan])")
)

type Record struct {
	Year  int
	Name  string
	Rank  int
	Count int
}

type MF struct {
	M, F Record
}

type SSA []MF

func parseSelection(sel *goquery.Selection) (int, error) {
	return strconv.Atoi(strings.Replace(sel.Text(), ",", "", -1))
}

func (self PopularNames) Year(year int) (SSA, error) {
	response, err := http.PostForm(ssaNamesUrl, self.valuesForYear(year))
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}

	retval := SSA{}

	// We asked for 1000 rows.  We skip the first row, and the last row is non-inclusive.
	doc.FindMatcher(trMatcher).Slice(1, 1001).EachWithBreak(
		func(row int, rowSel *goquery.Selection) bool {
			mf := MF{
				M: Record{
					Year: year,
				},
				F: Record{
					Year: year,
				},
			}

			// See the example input above to understand column position.
			rowSel.FindMatcher(thMatcher).EachWithBreak(
				func(col int, colSel *goquery.Selection) bool {
					switch col {
					case 0:
						if rank, err := parseSelection(colSel); err != nil {
							return false
						} else {
							mf.M.Rank = rank
							mf.F.Rank = rank
						}
					case 1:
						mf.M.Name = colSel.Text()
					case 2:
						mf.M.Count, err = parseSelection(colSel)
						if err != nil {
							return false
						}
					case 3:
						mf.F.Name = colSel.Text()
					case 4:
						mf.F.Count, err = parseSelection(colSel)
						if err != nil {
							return false
						}
					}
					return true
				})

			if err != nil {
				return false
			}

			retval = append(retval, mf)
			return true
		})

	if err != nil {
		return nil, err
	}

	return retval, nil
}

func add(session *mgo.Session, gender models.Gender, record Record) {
	c := session.DB("").C(models.Collection)
	c.EnsureIndex(mgo.Index{
		Key:    []string{"name", "gender"},
		Unique: true,
	})

	name := &models.Name{}
	c.Find(bson.M{"name": record.Name, "gender": gender}).One(&name)
	if name.Id == "" {
		name.Id = bson.NewObjectId()
		name.Name = record.Name
		name.Gender = gender
	}
	name.Rank.Set(record.Year, record.Rank)
	name.Count.Set(record.Year, record.Count)
	c.UpsertId(name.Id, bson.M{"$set": name})
}

func main() {
	flag.Parse()

	session, err := mgo.Dial(*uri)
	if err != nil {
		log.Fatal(err)
	}

	ssaNames := PopularNames{}
	for year := *from; year <= *to; year++ {
		fmt.Fprintf(os.Stderr, "%d\n", year)
		ssa, err := ssaNames.Year(year)
		if err != nil {
			log.Fatal(err)
		}

		for _, mf := range ssa {
			add(session, models.Male, mf.M)
			add(session, models.Female, mf.F)
		}
	}
}
