package fieldmap

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/language/en"
	"github.com/pborman/uuid"
)

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

const sourcefile = "refguide.json"

type RefGuideType string

func (g RefGuideType) Name() string {
	return string(g)
}

func (_ RefGuideType) Mapping() *bleve.IndexMapping {
	return allDocumentMapping(IndexDocumentMap{
		"resource": buildResourceDocumentMapping(),
	})
}

var RefGuide = RefGuideType("sfsbook.bleve")

func (_ RefGuideType) LoadStartData(i bleve.Index, pathroot string) error {
	log.Println("Indexing... now")

	jsonBytes, err := ioutil.ReadFile(filepath.Join(pathroot, sourcefile))
	if err != nil {
		return err
	}

	log.Println("read the database")

	// parse bytes as json
	var parsedResources []map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsedResources)
	if err != nil {
		return err
	}

	log.Println("parsed the database record")

	// So: how do I maintain flexibility in the handling of the fields?
	// Can unmarshal into a map of interface{}
	// I can set reasonable defaults.
	// Documents can have sub-documents...

	batch := i.NewBatch()
	for _, r := range parsedResources {
		rid := uuid.NewRandom().String()
		r["reviewed"] = false
		// This can be adapted to specify different types.
		r["_type"] = "resource"
		r["date_indexed"] = time.Now()
		batch.Index(rid, r)
	}

	log.Println("built a batch")

	err = i.Batch(batch)
	if err != nil {
		return err
	}
	log.Println("done Indexing...")
	return nil
}
