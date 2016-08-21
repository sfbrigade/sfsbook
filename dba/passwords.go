package dba

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
)

type password struct {
	Username     string `json:"username"`
	Passwordhash string `passwordhash:"uid"`
}

func (p password) Type() string {
	return "password"
}

func buildPasswordMapping() {
	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword_analyzer.Name

	// TODO(rjk): I might need to adjust this.
	// TODO(rjk): Move password data to a separate file.
	passwordFieldMapping := bleve.NewTextFieldMapping()
	passwordFieldMapping.Index = false
	passwordFieldMapping.IncludeTermVectors = false
	passwordFieldMapping.IncludeInAll = false

	// passwordEntryMapping is a document for each user.
	passwordEntryMapping := bleve.NewDocumentMapping()
	passwordEntryMapping.AddFieldMappingsAt("uid", keywordFieldMapping)
	passwordEntryMapping.AddFieldMappingsAt("passwordhash", passwordFieldMapping)
}
