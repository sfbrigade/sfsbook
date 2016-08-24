package dba

import (
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
	"github.com/blevesearch/bleve/analysis/language/en"
)

// IndexDocumentMap associates the document type (key) with
// the document type's index mapping (value) for a single database.
type IndexDocumentMap map[string]*bleve.DocumentMapping

// This code is largely inspired by the bleve beer-search demo application.
// Keeps all the state in the bleve database.

// TODO(rjk): Relocate the comment code.
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

// Standard field mappings. Use them everywhere.
var englishTextFieldMapping *bleve.FieldMapping
var keywordFieldMapping *bleve.FieldMapping
var ignoredFieldMapping *bleve.FieldMapping
var dateTimeMapping *bleve.FieldMapping
var boolFieldMapping *bleve.FieldMapping

// init makes all of the mappings that we use. They can be reused.
func init() {
	// a generic reusable mapping for english text
	englishTextFieldMapping = bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping = bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword_analyzer.Name

	// a generic reusable mapping for ignored content.
	ignoredFieldMapping = bleve.NewTextFieldMapping()
	ignoredFieldMapping.Store = false
	ignoredFieldMapping.Index = false
	ignoredFieldMapping.IncludeTermVectors = false
	ignoredFieldMapping.IncludeInAll = false

	// a date/time mapping
	// I believe that this is good like this. I will have to experiment.
	dateTimeMapping = bleve.NewDateTimeFieldMapping()

	// a generic reusable mapping for booleans
	boolFieldMapping = bleve.NewBooleanFieldMapping()
}

// buildCommentIndexMapping
// TODO(rjk): Comments are not supported yet. This code requires additional
// attention.
func buildCommentIndexMapping() *bleve.DocumentMapping {
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
	return commentEntryMapping
}

// buildResourceDocumentMapping builds the mappings needed for resource guide
// entries.
func buildResourceDocumentMapping() *bleve.DocumentMapping {
	resourceEntryMapping := bleve.NewDocumentMapping()

	// TODO(rjk): Make sure that I have full language support enabled.
	resourceEntryMapping.DefaultAnalyzer = en.AnalyzerName

	// With a default analyzer specified, we don't need to list the english field mappings.
	// resourceEntryMapping.AddFieldMappingsAt("uuid", keywordFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("email", keywordFieldMapping)

	// TODO(rjk): Support the indexing of the hand_sort later. At the moment, this is not
	// well structured. Later code will use the better-structured version of the data found
	// in the csv.
	resourceEntryMapping.AddFieldMappingsAt("hand_sort", ignoredFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("website", keywordFieldMapping)

	// I note in passing that this can be populated from the hand_sort data.
	// I might consider adding additional code to automatically freshen the data.
	resourceEntryMapping.AddFieldMappingsAt("wheelchair", ignoredFieldMapping)

	// To track if we have been reviewed.
	resourceEntryMapping.AddFieldMappingsAt("reviewed", boolFieldMapping)

	// Time when this resource was first added to the database and last modified.
	// TODO(rjk): Note need to track the edits.
	resourceEntryMapping.AddFieldMappingsAt("date_indexed", dateTimeMapping)
	resourceEntryMapping.AddFieldMappingsAt("date_modified", dateTimeMapping)

	return resourceEntryMapping
}

// allDocumentMapping creates a new top-level mapping for an entire database
// from the provided map of per-document mappings. (Per Bleve terminology, the
// database is conceptually an array of documents with an index of the document
// contents that permits finding sets of documents. In sfsbook, the fundamental
// document type is a single resource in the resource guide.)
//
// Each document type has a string name that should be the key in docMappings.
// Each document must have a _type key that specifies this key value in order
// to select the type of the document. Features like comments and edit auditing
// will introduce additional document types.
func allDocumentMapping(docMappings IndexDocumentMap) *bleve.IndexMapping {
	indexMapping := bleve.NewIndexMapping()

	for k, v := range docMappings {
		indexMapping.AddDocumentMapping(k, v)
	}

	// The document type is found (k in the map) by accessing the field named "_type"
	indexMapping.TypeField = "_type"

	// TODO(rjk): Currently the default language is english but I should do language
	// detection
	indexMapping.DefaultAnalyzer = "en"
	return indexMapping
}
