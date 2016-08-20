package dba

import (
	"log"

	"github.com/blevesearch/bleve"
)

// ResourceResultsGenerator finds a specific resource by its uuid.
type ResourceResultsGenerator struct {
	index bleve.Index
}

func MakeResourceResultsGenerator(bi bleve.Index) *ResourceResultsGenerator {
	return &ResourceResultsGenerator{index: bi}
}

type resourceResults struct {
	generatedResultCore

	// The requested resource.
	Uuid string

	// The actual fields in the document.
	Document map[string]interface{}
}

// TODO(rjk): Would it be so wrong to use the http Request? Database layer (aka middle) should
// be devoid of http concepts?
type ResourceRequest struct {
	Uuid     string
	IsPost   bool
	PostArgs map[string][]string
}

var immutableFields map[string]struct{}
var mustInitializeFields map[string]struct{}

func init() {
	o := struct{}{}
	immutableFields = map[string]struct{}{
		"_type":              o,
		" date_indexed":      o,
		"date_last_modified": o,
	}

	mustInitializeFields = map[string]struct{}{
		"reviewed": o,
	}
}

// ForRequest generates the data comprising a result page showing a single
// resource guide entry.
func (qr *ResourceResultsGenerator) ForRequest(req interface{}) GeneratedResult {
	request := req.(*ResourceRequest)
	request.validateIfNecessary()
	uuid := request.Uuid

	log.Println("uuid", uuid)

	// Code quality comment: Writing the templates requires knowing what I've
	// produced here. I feel that I have not layered this code very well.
	results := &resourceResults{
		generatedResultCore: generatedResultCore{
			// TODO(rjk): Should use error type. This is not idiomatic Go.
			success:     false,
			failureText: "query had a sad",
			debug:       false,
		},
		Uuid: uuid,
	}

	if request.IsPost {
		results.generatedResultCore.failureText = "update had a sad"
	}

	doc, err := qr.index.Document(uuid)
	if err != nil || doc == nil {
		log.Println("query failed", err)
		results.generatedResultCore.failureText = err.Error()
		return results
	}

	resultsMap, err := resultsMapFromDocument(doc)
	if err != nil {
		log.Println("couldn't convert doc to resultsMap")
		results.generatedResultCore.failureText = err.Error()
		return results
	}

	if err := qr.mergeAndUpdateIfNecessary(request, resultsMap); err != nil {
		log.Println("query failed", err)
		results.generatedResultCore.failureText = err.Error()
		return results
	}

	log.Println("succeeded", doc)
	results.generatedResultCore.success = true
	results.generatedResultCore.failureText = "No error."
	results.Document = resultsMap

	// Need to support showing the comments.
	return results
}
