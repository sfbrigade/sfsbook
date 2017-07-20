package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func SearchHandler(c *gin.Context) {
	query := c.Query("query")

	c.HTML(http.StatusOK, "search.tmpl", gin.H{
		"Query":     query,
		"Resources": []interface{}{},
	})
}

//----------------------------------------------------------
// /resources routes
//----------------------------------------------------------

func ResourcesIndexHandler(c *gin.Context) {
	resources, err := GetResources()
	if err != nil {
		log.Fatal(err)
	}

	c.HTML(http.StatusOK, "resource.tmpl", gin.H{
		"Resources": resources,
	})
}

func ResourceShowHandler(c *gin.Context) {
	slug := c.Param("slug")

	resource, err := GetResourceBySlug(slug)
	if err != nil {
		log.Fatal(err)
	}

	c.HTML(http.StatusOK, "resource.tmpl", gin.H{
		"Query":     slug,
		"Resources": []*Resource{resource},
	})
}

func ResourcesJsIndexHandler(c *gin.Context) {
	resources, err := GetResources()
	if err != nil {
		log.Fatal(err)
	}

	c.HTML(http.StatusOK, "resources.js.tmpl", gin.H{
		"Resources": resources,
	})
}
