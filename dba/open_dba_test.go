package dba

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/sfbrigade/sfsbook/dba/fieldmap"

	"log"
)

const testdata = `[
    {
        "address": "100 Json Blvd, JavaScript, BrowserLand",
        "categories": "Test data",
        "description": "A filler record",
        "email": "me@you.com",
        "hand_sort": [
            "Crisis line: Business line: Fax:",
            "TDD:",
            "Insurance: Fees: Hours: Ages:",
            "800-522-0925",
            "414-274-0925",
            "414-272-2870",
            "FREE",
            "",
            ""
        ],
        "languges": "english spanish",
        "name": "test record 1",
        "pops_served": "This entry needs a populations served list",
        "services": "This entry needs a services list",
        "website": "www.9to5.org",
        "wheelchair": "This entry needs Wheelchair accessibliy info"
    }
]
`

func TestIndexResourcet(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "sfsbook")
	log.Println(tmpdir)
	if err != nil {
		t.Fatal("can't make a temporary directory", err)
	}
	// defer os.RemoveAll(tmpdir)

	// Create a database. Should fail because one or more of the paths doesn't exist.
	if _, err := OpenBleve(tmpdir, fieldmap.RefGuide); err == nil {
		t.Fatal("OpenBleve succeeded adding a non-existent starter database when it should have failed", err)
	}

	// Stick some data in the file.
	file, err := os.Create(filepath.Join(tmpdir, "refguide.json"))
	if err != nil {
		t.Fatal("can't openfile in tmp directory", tmpdir, "because", err)
	}
	if n, err := io.WriteString(file, testdata); err != nil || n != len(testdata) {
		t.Fatal("can't write testdata to tmpdir file because", err)
	}
	file.Close()

	// Create a database. Should succeed.
	db, err := OpenBleve(tmpdir, fieldmap.RefGuide)
	if err != nil {
		t.Fatal("OpenBleve failed to open and index some testdata", err)
	}
	defer db.Close()

	// There is a datum in the database.
	if n, err := db.DocCount(); n != 1 || err != nil {
		t.Error("expected to find data in the database, count is", n, "or error getting count", err)
	}

	fields, err := db.Fields()
	if err != nil {
		t.Fatal("couldn't retrieve the fields from the Bleve database because", err)
	}

	want_fields := []string{"_all", "_type", "address", "categories", "date_indexed", "description", "email", "hand_sort", "languges", "name", "pops_served", "reviewed", "services", "website", "wheelchair"}

	sort.Strings(want_fields)
	sort.Strings(fields)

	if want, got := want_fields, fields; !reflect.DeepEqual(want, got) {
		t.Errorf("wanted %#v but got %#v\n", want_fields, fields)
	}

}
