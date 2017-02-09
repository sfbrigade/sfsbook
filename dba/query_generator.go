package dba

import (
	"log"
	"net/http"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

// QueryResultsGenerator does search queries against the resource book.
type QueryResultsGenerator struct {
	index bleve.Index
}

func MakeQueryResultsGenerator(bi bleve.Index) *QueryResultsGenerator {
	return &QueryResultsGenerator{index: bi}
}

// I would obviously want to expand on this.
type ResourceResult struct {
	Uuid        string
	Categories  string
	Description string
	Name        string
	Services    string
	Address     string
}

type queryResults struct {
	generatedResultCore

	Query     string
	Resources []ResourceResult
}

func (qr *QueryResultsGenerator) ForRequest(param interface{}) GeneratedResult {
	req := param.(*http.Request)
	// TODO(rjk): manage cookies etc.

	// TODO(rjk): web doesn't belong in DBA. And dba doesn't belong in server so fix up your layering issues.
	// The correct middle-tier interface would be map of values?
	// And some of this code would then be movable?
	// Maybe there's some kind of existing library to map form params to struct entries?
	// Or it could be generated automatically? Use: https://github.com/gorilla/schema
	log.Println("req.Form", req.Form)

	// This is dorky flow... I might want to do something different.
	querystring := ""

	// TODO(rjk): validation of querystring parameters should happen in the web layer.
	if q, ok := req.Form["query"]; ok {
		if len(q) > 0 {
			querystring = q[0]
		}
	}

	results := &queryResults{
		generatedResultCore: generatedResultCore{
			success:     false,
			failureText: "query had a sad",
			debug:       false,
		},
		Query: querystring,
	}

	// Query goals (long term)?
	// I think we need some way for the user to specify terms for search.

	// Actually do a query against the database.
	// TODO(rjk): Refine the search handling.
	middleq := make([]query.Query, 0, 5)
	phrases := strings.Split(querystring, ", ")
	for _, phrase := range phrases {
		//  for _, k := range []string{"description", "services", "categories", "name", "website", "email", "address" } {
		q := query.NewMatchPhraseQuery(phrase)
		//	  q.SetField(k)
		middleq = append(middleq, q)
		//  }
	}

	bq := query.NewBooleanQuery(
		middleq,
		[]query.Query{},
		[]query.Query{})

	// Makes a search request.
	log.Println("querystring:", querystring)
	searchRequest := bleve.NewSearchRequest(bq)

	// Modify the search request to only retrieve some fields.
	searchRequest.Fields = []string{"name", "categories", "description", "services", "website", "email", "address"}
	searchRequest.Size = 10
	// Advance this to move forward through the result set...
	searchRequest.From = 0

	// Get a highlighted output.
	//	searchRequest.Highlight = bleve.NewHighlight()

	// Search.
	searchresults, err := qr.index.Search(searchRequest)

	// I need to rationalize the error handling..
	if err != nil {
		// TODO(rjk): update the FailureText
		log.Println("query failed", err)
		return results
	}

	results.generatedResultCore.success = true
	results.generatedResultCore.failureText = ""

	// Verbose... but handy
	log.Println(searchresults)

	results.Resources = make([]ResourceResult, len(searchresults.Hits))

	c := 0
	for _, sr := range searchresults.Hits {
		//		log.Println(sr.Fields)
		// TODO(rjk): There is a lot of boilerplate here. Maybe I can be clever.
		results.Resources[c].Uuid = sr.ID
		results.Resources[c].Name = sr.Fields["name"].(string)
		results.Resources[c].Services = sr.Fields["services"].(string)
		results.Resources[c].Categories = sr.Fields["categories"].(string)
		results.Resources[c].Description = sr.Fields["description"].(string)
		results.Resources[c].Address = sr.Fields["address"].(string)

		c++
		if c > 10 {
			break
		}
	}

	return results
}

