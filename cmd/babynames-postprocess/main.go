package main

import (
	"flag"
	"log"
	"strings"

	"github.com/shutej/babynames/helpers"
	"github.com/shutej/babynames/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	uri      = flag.String("uri", "", "URI to connect to MongoDB")
	distance = flag.Int("distance", 1, "maximum string edit distance of related names")
)

func main() {
	flag.Parse()

	session, err := mgo.Dial(*uri)
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB("").C(models.Collection)

	names := []models.Name{}
	c.Find(nil).All(&names)

	for i := 0; i < len(names); i++ {
		a := strings.ToLower(names[i].Name)
		names[i].LowerName = a
		names[i].TotalCount = 0
		names[i].Count.Each(func(key int, data interface{}) {
			names[i].TotalCount += data.(int)
		})

		for j := 0; j < i; j++ {
			b := names[j].LowerName
			if len(a)-len(b) > *distance ||
				len(b)-len(a) > *distance {
				continue
			}

			if helpers.LevenshteinDistance(&a, &b) <= *distance {
				names[i].RelatedNames = append(names[i].RelatedNames, names[j].Name)
				names[j].RelatedNames = append(names[j].RelatedNames, names[i].Name)
			}
		}
	}

	for _, name := range names {
		c.UpsertId(name.Id, bson.M{"$set": name})
	}
}
