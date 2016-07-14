package dba

import (
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
	"github.com/blevesearch/bleve/analysis/language/en"
)

// This code is largely inspired by the bleve beer-search demo application.
// Keeps all the state in the bleve database.

// password contains the user identity records.
type password struct {
	Username string  `json:"username"`
	Passwordhash string  `passwordhash:"uid"`
}

func (p password) Type() string {
	return "password"
}

type comment struct {
	// The uuid of the associated resource card.
	ResourceUuid string  `json:"resourceuuid"`
	CreationTime time.Time `json:"creationtime"`
	UpdateTime time.Time `json:"updatetime"`
	Owner string `json:"owner"`
	Viewability string `json:"viewability"`
	Body string `json:"body"`
}

func (p comment) Type() string {
	return "comment"
}

type resource struct {
	Uuid string `json:"uuid"`
	// TODO(rjk): can rationalize this?
	Address string `json:"address"`

	// TODO(rjk): try to make it a structured list.
	Categories string `json:"categories"`
	
	Description string `json:"description"`
	Email string `json:"email"`
	// TODO(rjk): try to make it a structured list.
	Languages string `json:"languages"`
	Name string `json:"name"`
	PopsServed string `json:"pops_served"`
	Services string `json:"services"`

	// TODO(rjk): need a validator
	// Aside: there should be a validator service provided by the app.
	Website string `json:"website"`

	// TODO(rjk): this should be a boolean.
	Wheelchair string `json:"wheelchair"`

	// TODO(rjk): we need some way to encode phone numbers.
	// TODO(rjk): we need some way to encode when a facility is open.
}

func (p resource) Type() string {
	return "resource"
}

// buildIndexMapping 
func buildIndexMapping() (*bleve.IndexMapping, error) {
	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword_analyzer.Name

	// a generic reusable mapping for ignored content.
	ignoredFieldMapping := bleve.NewTextFieldMapping()
	ignoredFieldMapping.Store = false
	ignoredFieldMapping.Index = false
	ignoredFieldMapping.IncludeTermVectors = false
	ignoredFieldMapping.IncludeInAll = false

	// TODO(rjk): I might need to adjust this.
	passwordFieldMapping := bleve.NewTextFieldMapping()
	ignoredFieldMapping.Index = false
	ignoredFieldMapping.IncludeTermVectors = false
	ignoredFieldMapping.IncludeInAll = false

	// TODO(rjk): There is a an open-tail of effort to do here.
	// resourceEntryMapping is the mappings for each of the resource
	// entries.
	resourceEntryMapping := bleve.NewDocumentMapping()
	resourceEntryMapping.AddFieldMappingsAt("uuid", keywordFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("address", englishTextFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("categories", englishTextFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("description", englishTextFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("email", keywordFieldMapping)

	// TODO(rjk): Support the indexing of the hand_sort later. At the moment, this is not
	// well structured. 
	// resourceEntryMapping.AddFieldMappingsAt("hand_sort", ignoredFieldMapping)
	
	resourceEntryMapping.AddFieldMappingsAt("languages", englishTextFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("pops_served", englishTextFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("services", englishTextFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("website", keywordFieldMapping)

	// I note in passing that this can be populated from the hand_sort data.
	// I might consider adding additional structure to this.
	// TODO(rjk): Perhaps it could be auto-populated.
	resourceEntryMapping.AddFieldMappingsAt("wheelchair", ignoredFieldMapping)


	// TODO(rjk): Create structures for comments. Create structures for resourceEntries.

	// commentEntryMapping is a document for each comment.
	commentEntryMapping := bleve.NewDocumentMapping()
	commentEntryMapping.AddFieldMappingsAt("uuid", keywordFieldMapping)
	// comment creation date.
	// comment update date.
	// the uid of the user.
	commentEntryMapping.AddFieldMappingsAt("owner", keywordFieldMapping)
	// comments can be vieweable by signed in user ("me", signed-in "volunteers", "world")
	commentEntryMapping.AddFieldMappingsAt("viewability", keywordFieldMapping)
	commentEntryMapping.AddFieldMappingsAt("body", englishTextFieldMapping)

	// passwordEntryMapping is a document for each user.
	passwordEntryMapping := bleve.NewDocumentMapping()
	passwordEntryMapping.AddFieldMappingsAt("uid", keywordFieldMapping)
	passwordEntryMapping.AddFieldMappingsAt("passwordhash", passwordFieldMapping)


	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("resource", resourceEntryMapping)
	indexMapping.AddDocumentMapping("comment", commentEntryMapping)
	indexMapping.AddDocumentMapping("password", passwordEntryMapping)

	// TODO(rjk): Implement some fieldstuff.
//	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}
