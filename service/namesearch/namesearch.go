package namesearch

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shutej/namestats/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Option func(*Service)

type Service struct {
	session *mgo.Session
}

func Session(session *mgo.Session) Option {
	return func(self *Service) {
		self.session = session
	}
}

func New(options ...Option) *Service {
	self := &Service{}
	for _, option := range options {
		option(self)
	}
	return self
}

func (self *Service) db() *mgo.Database {
	return self.session.DB("")
}

type response struct {
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
	Total      int           `json:"total"`
	NameSearch []models.Name `json:"namesearch"`
}

func (self *Service) NameSearch(ctx *gin.Context) {
	c := self.db().C(models.Collection)

	query := ctx.Param("query")

	response := response{
		NameSearch: []models.Name{},
	}
	response.Offset, _ = strconv.Atoi(ctx.Query("offset"))
	response.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	if response.Limit == 0 {
		response.Limit = 10
	}

	// For scans...
	c.EnsureIndex(mgo.Index{
		Key: []string{"totalCount"},
	})

	// For prefix scans...
	c.EnsureIndex(mgo.Index{
		Key: []string{"lowerName", "totalCount"},
	})

	var q *mgo.Query
	if query == "" {
		q = c.Find(nil)
	} else {
		q = c.Find(bson.M{
			"lowerName": bson.RegEx{
				Pattern: "^" + regexp.QuoteMeta(strings.ToLower(query)),
			},
		})
	}

	var err error

	if response.Total, err = q.Count(); err != nil {
		ctx.String(500, "%v", err)
		return
	}

	var name models.Name
	iter := q.Skip(response.Offset).Limit(response.Limit).Sort("-totalCount").Iter()
	for iter.Next(&name) {
		if name.RelatedNames == nil {
			name.RelatedNames = []string{}
		}
		response.NameSearch = append(response.NameSearch, name)
	}

	if err := iter.Close(); err != nil {
		ctx.String(500, "%v", err)
	}

	ctx.JSON(200, response)
}
