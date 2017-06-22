package dba

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
	"github.com/pborman/uuid"
)

// buildResourceDocumentMapping builds the mappings needed for resource guide
// entries.
func buildResourceDocumentMapping() *mapping.DocumentMapping {
	resourceEntryMapping := bleve.NewDocumentMapping()

	// TODO(rjk): Make sure that I have full language support enabled.
	resourceEntryMapping.DefaultAnalyzer = en.AnalyzerName

	// With a default analyzer specified, we don't need to list the english field mappings.
	// resourceEntryMapping.AddFieldMappingsAt("uuid", keywordFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("email", KeywordFieldMapping)

	// TODO(rjk): Support the indexing of the hand_sort later. At the moment, this is not
	// well structured. Later code will use the better-structured version of the data found
	// in the csv.
	resourceEntryMapping.AddFieldMappingsAt("hand_sort", IgnoredFieldMapping)
	resourceEntryMapping.AddFieldMappingsAt("website", KeywordFieldMapping)

	// I note in passing that this can be populated from the hand_sort data.
	// I might consider adding additional code to automatically freshen the data.
	resourceEntryMapping.AddFieldMappingsAt("wheelchair", IgnoredFieldMapping)

	// To track if we have been reviewed.
	resourceEntryMapping.AddFieldMappingsAt("reviewed", BoolFieldMapping)

	// Time when this resource was first added to the database and last modified.
	// TODO(rjk): Note need to track the edits.
	resourceEntryMapping.AddFieldMappingsAt("date_indexed", DateTimeMapping)
	resourceEntryMapping.AddFieldMappingsAt("date_modified", DateTimeMapping)

	return resourceEntryMapping
}

const sourcefile = "refguide.json"

type RefGuideType string

func (g RefGuideType) Name() string {
	return string(g)
}

func (_ RefGuideType) Mapping() *mapping.IndexMappingImpl {
	return AllDocumentMapping(IndexDocumentMap{
		"resource": buildResourceDocumentMapping(),
	})
}

var RefGuide = RefGuideType("sfsbook.bleve")

func (_ RefGuideType) LoadStartData(i bleve.Index, pathroot string) error {
	log.Println("RefGuideType LoadStartData")

	jsonBytes, err := ioutil.ReadFile(filepath.Join(pathroot, sourcefile))
	if err != nil {
		return err
	}

	log.Println("LoadStartData: read the database")

	// parse bytes as json
	var parsedResources []map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsedResources)
	if err != nil {
		return err
	}

	log.Println("LoadStartData: parsed the database records")

	// So: how do I maintain flexibility in the handling of the fields?
	// Can unmarshal into a map of interface{}
	// I can set reasonable defaults.
	// Documents can have sub-documents...

	uuidsList := make([]string, len(parsedResources))

	batch := i.NewBatch()
	for i, r := range parsedResources {
		rid := uuid.NewRandom().String()
		r["reviewed"] = false
		// This can be adapted to specify different types.
		r["_type"] = "resource"
		r["date_indexed"] = time.Now()
		batch.Index(rid, r)

		uuidsList[i] = rid
	}

	if err := i.Batch(batch); err != nil {
		return err
	}
	log.Println("LoadStartData: indexed resources")

	// index uuid slice since bleve doesn't have a way to iterate through all keys
	if err := i.Index(UUIDsIndexName, uuidsList); err != nil {
		return err
	}
	log.Println("LoadStartData: indexed uuid of resources")

	return nil
}
