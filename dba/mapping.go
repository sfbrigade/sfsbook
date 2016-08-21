package dba

import (
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
	"github.com/blevesearch/bleve/analysis/language/en"
)

// This code is largely inspired by the bleve beer-search demo application.
// Keeps all the state in the bleve database.


type comment struct {
	// The uuid of the associated resource card.
	ResourceUuid string    `json:"resourceuuid"`
	CreationTime time.Time `json:"creationtime"`
	UpdateTime   time.Time `json:"updatetime"`
	Owner        string    `json:"owner"`
	Viewability  string    `json:"viewability"`
	Body         string    `json:"body"`
}

func (p comment) Type() string {
	return "comment"
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

	// a date/time mapping
	// I believe that this is good like this. I will have to experiment.
	dateTimeMapping := bleve.NewDateTimeFieldMapping()

	// a generic reusable mapping for booleans
	boolFieldMapping := bleve.NewBooleanFieldMapping()

	// TODO(rjk): There is a an open-tail of effort to do here.
	// resourceEntryMapping is the mappings for each of the resource
	// entries.
	resourceEntryMapping := bleve.NewDocumentMapping()
	// TODO(rjk): Make sure that I have full language support enabled.
	resourceEntryMapping.DefaultAnalyzer = en.AnalyzerName

	// With a default analyzer specified, we don't need to list the english field mappings.
	resourceEntryMapping.AddFieldMappingsAt("uuid", keywordFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("email", keywordFieldMapping)

	// TODO(rjk): Support the indexing of the hand_sort later. At the moment, this is not
	// well structured.
	resourceEntryMapping.AddFieldMappingsAt("hand_sort", ignoredFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("website", keywordFieldMapping)

	// I note in passing that this can be populated from the hand_sort data.
	// I might consider adding additional structure to this.
	// TODO(rjk): Perhaps it could be auto-populated.
	resourceEntryMapping.AddFieldMappingsAt("wheelchair", ignoredFieldMapping)

	// To track if we have been reviewed.
	resourceEntryMapping.AddFieldMappingsAt("reviewed", boolFieldMapping)

	// Time when this resource was first added to the database.
	resourceEntryMapping.AddFieldMappingsAt("date_indexed", dateTimeMapping)
	resourceEntryMapping.AddFieldMappingsAt("date_modified", dateTimeMapping)

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

	// I don't support multiple types. But I could use this to address CSV, web and JSON
	// updated documents? Particularly given that they can have different fields.
	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("resource", resourceEntryMapping)

	// TODO(rjk): Implement some fieldstuff.
	indexMapping.TypeField = "_type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}
