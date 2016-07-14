package dba

import (
	"log"
	"net/http"

	"github.com/blevesearch/bleve"
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
	Uuid string
	Categories string
	Description string
	Name string
	Services string
}

type queryResults struct {
	Query string
	Success bool

	// Text to display if something has gone wrong.
	// Pull from internal resource table for i18n
	FailureText string
	Resources []ResourceResult
}

func (qr *QueryResultsGenerator) ForRequest(req *http.Request) interface{} {
	// TODO(rjk): manage cookies etc.

	// TODO(rjk): web doesn't belong in DBA. And dba doesn't belong in server so fix up your layering issues.
	// The correct middle-tier interface would be map of values?
	// And some of this code would then be movable?
	// Maybe there's some kind of existing library to map form params to struct entries?
	// Or it could be generated automatically? Use: https://github.com/gorilla/schema
	log.Println("req.Form", req.Form)

	// This is dorky flow... I might want to do something different.
	query := ""

	// TODO(rjk): validation of query parameters should happen in the web layer.
	if q, ok := req.Form["query"]; ok {
		if len(q) > 0 {
			query = q[0]
		}
	}

	results := &queryResults{
		Query: query,
		FailureText: "query had a sad",
	}	

	// Query goals (long term)?
	// I think we need some way for the user to specify terms for search.	

	// Actually do a query against the database.
	// how do I limit this to only some document types and fields?
	// I don't know how to do that. SetField() on the query?
	// create a more complicated query


	// Add specific terms as "must"
	// This makes a match query for the description field.
	bq := bleve.NewBooleanQuery(
		[]bleve.Query{},
		[]bleve.Query{
			bleve.NewMatchPhraseQuery(query).SetField("description"),
			bleve.NewMatchPhraseQuery(query).SetField("services"),
			bleve.NewMatchPhraseQuery(query).SetField("categories"),
			bleve.NewMatchPhraseQuery(query).SetField("name"),
		},
		[]bleve.Query{})

	// Makes a search request.
	searchRequest := bleve.NewSearchRequest(bq)

	// Modify the search request to only retrieve some fields.
	searchRequest.Fields = []string{ "name", "categories", "description", "services" }
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
		
		c++
		if c > 10 {
			break
		}
	}

	return results
}
