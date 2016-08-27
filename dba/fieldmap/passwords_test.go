package fieldmap

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/blevesearch/bleve"
	"github.com/sfbrigade/sfsbook/dba"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordFile(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "sfsbook")
	// log.Println(tmpdir)
	if err != nil {
		t.Fatal("can't make a temporary directory", err)
	}
	defer os.RemoveAll(tmpdir)

	// Create a database. Should fail because we refuse to create an
	// insecure password file without a command line flag.
	if _, err := dba.OpenBleve(tmpdir, PasswordFile); err == nil {
		t.Fatal("OpenBleve succeeded adding a password database when it should have refused", err)
	}

	original_init_passwords := *init_passwords
	*init_passwords = true
	defer func() { *init_passwords = original_init_passwords }()

	db, err := dba.OpenBleve(tmpdir, PasswordFile)
	if err != nil {
		t.Fatal("OpenBleve failed to add a password file", err)
	}

	// Records exist.
	if n, err := db.DocCount(); n != 2 || err != nil {
		t.Error("expected to find data in the database, count is", n, "or error getting count", err)
	}

	// Search the database.
	for _, uname_passwd := range [][]string{[]string{"admin", "sesame"}, []string{"volunteer", "open"}} {
		uname := uname_passwd[0]
		sreq := bleve.NewSearchRequest(bleve.NewMatchQuery(uname))
		// Note that the result only contains the fields specified here.
		sreq.Fields = []string{"name", "cost", "passwordhash"}

		searchResults, err := db.Search(sreq)
		if err != nil {
			t.Error("couldn't search the password file", err)
			continue
		}

		if len(searchResults.Hits) != 1 {
			t.Error("expected 1 match for", uname, "but got", len(searchResults.Hits))
			continue
		}

		// validate data
		sr := searchResults.Hits[0]

		if got, want := int(sr.Fields["cost"].(string)[0]), bcrypt.DefaultCost; got != want {
			t.Error("cost wrong got", got, "want", want)
		}

		pw := sr.Fields["passwordhash"].(string)
		passwd := uname_passwd[1]
		if pw == passwd {
			t.Error("passwords not actually encrypted?")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(pw), []byte(passwd)); err != nil {
			t.Errorf("password %s not encoded in a way (%#v) that we can use it to decode %v", passwd, pw, err)
		}
	}
}
