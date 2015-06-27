package models

import (
	"gopkg.in/mgo.v2/bson"
)

const Collection = "names"

type Gender string

const (
	Male   Gender = "m"
	Female        = "f"
)

type Name struct {
	Id     bson.ObjectId `bson:"_id"    json:"-"`
	Gender Gender        `bson:"gender" json:"gender"`
	Name   string        `bson:"name"   json:"name"`
	Rank   Series        `bson:"rank"   json:"rank"`
	Count  Series        `bson:"count"  json:"count"`

	// Filled in by postprocessor
	LowerName    string   `bson:"lowerName,omitempty"    json:"-"`
	RelatedNames []string `bson:"relatedNames,omitempty" json:"relatedNames"`
	TotalCount   int      `bson:"totalCount,omitempty"   json:"totalCount"`
}
