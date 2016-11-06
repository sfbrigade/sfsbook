package dba

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/lang/en"
)

// IndexDocumentMap associates the document type (key) with
// the document type's index mapping (value) for a single database.
type IndexDocumentMap map[string]*mapping.DocumentMapping

// Standard field mappings. Use them everywhere.
var EnglishTextFieldMapping *mapping.FieldMapping
var KeywordFieldMapping *mapping.FieldMapping
var IgnoredFieldMapping *mapping.FieldMapping
var DateTimeMapping *mapping.FieldMapping
var BoolFieldMapping *mapping.FieldMapping

// init makes all of the mappings that we use. A single instance of a
// mapping such as englishTextMapping can be used for any number of
// fields.
func init() {
	// a generic reusable mapping for english text
	EnglishTextFieldMapping = bleve.NewTextFieldMapping()
	EnglishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	KeywordFieldMapping = bleve.NewTextFieldMapping()
	KeywordFieldMapping.Analyzer = keyword.Name

	// a generic reusable mapping for ignored content.
	IgnoredFieldMapping = bleve.NewTextFieldMapping()
	IgnoredFieldMapping.Store = false
	IgnoredFieldMapping.Index = false
	IgnoredFieldMapping.IncludeTermVectors = false
	IgnoredFieldMapping.IncludeInAll = false

	// a date/time mapping
	// I believe that this is good like this. I will have to experiment.
	DateTimeMapping = bleve.NewDateTimeFieldMapping()

	// a generic reusable mapping for booleans
	BoolFieldMapping = bleve.NewBooleanFieldMapping()
}

// AllDocumentMapping creates a new top-level mapping for an entire database
// from the provided map of per-document mappings. (Per Bleve terminology, the
// database is conceptually an array of documents with an index of the document
// contents that permits finding sets of documents. In sfsbook, the fundamental
// document type is a single resource in the resource guide.)
//
// Each document type has a string name that should be the key in docMappings.
// Each document must have a _type key that specifies this key value in order
// to select the type of the document. Features like comments and edit auditing
// will introduce additional document types.
func AllDocumentMapping(docMappings IndexDocumentMap) *mapping.IndexMappingImpl {
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
