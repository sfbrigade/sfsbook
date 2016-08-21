package dba

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
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
	if err != nil {
		t.Fatal("can't make a temporary directory", err)
	}
	defer os.RemoveAll(tmpdir)

	// Create a database. Should fail because one or more of the paths doesn't exist.
	if _, err := OpenBleve(tmpdir);  err == nil {
		t.Fatal("OpenBleve succeeded adding a non-existent starter database when it should have failed", err)
	}

	// Stick some data in the file.
	file, err := os.Create(filepath.Join(tmpdir, sourcefile))
	if err != nil {
		t.Fatal("can't openfile in tmp directory",  tmpdir, "because", err)
	}
	if n, err := io.WriteString(file, testdata); err != nil || n != len(testdata)  {
		t.Fatal("can't write testdata to tmpdir file because", err)
	}
	file.Close()
	
	// Create a database. Should succeed.
	db, err := OpenBleve(tmpdir)
	if err != nil {
		t.Fatal("OpenBleve failed to open an index a testdata", err)
	}
	
	// There is datum in the database.
	if n, err := db.DocCount(); n != 1 || err != nil {
		t.Error("expected to find data in the database, count is", n , "or error getting count", err)
	}

}