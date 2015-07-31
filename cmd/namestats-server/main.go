// This is the namestats-server command.  See README.md for details.
package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/shutej/namestats/service/namesearch"
	mgo "gopkg.in/mgo.v2"
)

var (
	listen = flag.String("listen", ":80", "port to listen on")
	uri    = flag.String("uri", "", "URI of MongoDB")
)

func main() {
	flag.Parse()

	session, err := mgo.Dial(*uri)
	if err != nil {
		log.Fatal(err)
	}

	service := namesearch.New(namesearch.Session(session))

	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(static.Serve("/", static.LocalFile("static", true)))

	router.GET("/v1/namesearch/", service.NameSearch)
	router.GET("/v1/namesearch/:query", service.NameSearch)

	router.Run(*listen)
}
