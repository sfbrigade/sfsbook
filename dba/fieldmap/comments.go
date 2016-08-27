package fieldmap

import (
	"time"

	"github.com/blevesearch/bleve"
)

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
