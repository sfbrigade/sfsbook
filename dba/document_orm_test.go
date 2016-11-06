package dba

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/pborman/uuid"
	"github.com/sfbrigade/sfsbook/dba/fieldmap"

	"github.com/blevesearch/bleve/registry"
)

// TestDatabaseType implements IndexFactory for test data.
type TestDatabaseType struct {
	name string
	data []map[string]interface{}
}

func (g *TestDatabaseType) Name() string {
	return g.name
}

func buildTestDocumentMapping() *mapping.DocumentMapping {
	numberMapping := bleve.NewNumericFieldMapping()

	// testDocumentMapping is a document mapping for tests.
	testDocumentMapping := bleve.NewDocumentMapping()
	testDocumentMapping.DefaultAnalyzer = keyword.Name

	testDocumentMapping.AddFieldMappingsAt("textfield", fieldmap.KeywordFieldMapping)
	testDocumentMapping.AddFieldMappingsAt("numberfield", numberMapping)
	testDocumentMapping.AddFieldMappingsAt("datefield", fieldmap.DateTimeMapping)
	testDocumentMapping.AddFieldMappingsAt("istruefield", fieldmap.BoolFieldMapping)

	return testDocumentMapping
}

// Need a mapping for each type.
func (_ TestDatabaseType) Mapping() *mapping.IndexMappingImpl {
	return fieldmap.AllDocumentMapping(fieldmap.IndexDocumentMap{
		"testdoc": buildTestDocumentMapping(),
	})
}

func (g *TestDatabaseType) LoadStartData(i bleve.Index, pathroot string) error {
	s := g.data

	batch := i.NewBatch()
	for _, r := range s {
		rid := uuid.NewRandom()
		r["_type"] = "testdoc"
		batch.Index(string(rid), r)
	}

	if err := i.Batch(batch); err != nil {
		return err
	}
	return nil
}

func TestMakeMapFromDocument(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "dba")
	// log.Println(tmpdir)
	if err != nil {
		t.Fatal("can't make a temporary directory", err)
	}
	defer os.RemoveAll(tmpdir)

	// NB: the '_' precluded the use of an MatchQuery. This
	// suggests that I should change password search to Term.
	testfilesetup := &TestDatabaseType{
		name: "test.bleve",
		data: []map[string]interface{}{
			map[string]interface{}{
				"textfield":   "string_value",
				"numberfield": 100.1,
				"datefield":   time.Unix(1000, 0),
				"istruefield": true,
			},
			map[string]interface{}{
				"textfield":   "second_string_value",
				"numberfield": 200.1,
				"datefield":   time.Unix(2000, 0),
				"istruefield": false,
			},
		},
	}

	a, b := registry.AnalyzerTypesAndInstances()
	t.Log("AnalyzerTypesAndInstances:", a, b)

	// Open and populate a database.
	db, err := OpenBleve(tmpdir, testfilesetup)
	if err != nil {
		t.Fatal("OpenBleve failed to add a test database:", err)
	}

	// Records exist.
	if n, err := db.DocCount(); n != 2 || err != nil {
		t.Error("expected to find data in the database, count is", n, "or error getting count", err)
	}

	// Search the database.
	for _, r := range testfilesetup.data {
		searchval := r["textfield"].(string)
		t.Log("searchval", searchval)
		sreq := bleve.NewSearchRequest(bleve.NewTermQuery(searchval))
		// Note that the result Fields contain only the ones listed here.
		sreq.Fields = []string{"textfield", "numberfield"}

		searchResults, err := db.Search(sreq)
		if err != nil {
			// Is this what lies underneath my error cases.
			t.Fatal("couldn't search the test database", err)
			continue
		}

		f, _ := db.Fields()
		t.Log("fields:", f)

		if len(searchResults.Hits) != 1 {
			t.Fatal("expected 1 match for", searchval, "but got", len(searchResults.Hits))
			continue
		}

		// check that everything worked.
		uuid := searchResults.Hits[0].ID
		doc, err := db.Document(uuid)
		if err != nil {
			t.Fatalf("failed to find %v in database: %v", uuid, err)
		}

		ormdoc, err := MakeMapFromDocument(doc)
		if err != nil {
			t.Fatalf("failed to ORM convert %v because %v", ormdoc, err)
		}

		if got, want := ormdoc["_type"].(string), r["_type"].(string); got != want {
			t.Errorf("invalid ORM mapping got %v but want %v", got, want)
		}
		if got, want := ormdoc["istruefield"].(bool), r["istruefield"].(bool); got != want {
			t.Errorf("invalid ORM mapping got %v but want %v", got, want)
		}
		if got, want := ormdoc["textfield"].(string), r["textfield"].(string); got != want {
			t.Errorf("invalid ORM mapping got %v but want %v", got, want)
		}
		if got, want := ormdoc["numberfield"].(float64), r["numberfield"].(float64); got != want {
			t.Errorf("invalid ORM mapping got %v but want %v", got, want)
		}
		if got, want := ormdoc["datefield"].(time.Time), r["datefield"].(time.Time); !got.Equal(want) {
			t.Errorf("invalid ORM mapping got %v but want %v", got, want)
		}
	}
}
