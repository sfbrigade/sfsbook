package server

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestEmbeddableResource(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "sfsbook")
	if err != nil {
		t.Fatal("can't make a temporary directory", err)
	}
	defer os.RemoveAll(tmpdir)

	if _, ok := Resources["/file_finder_test.html"]; ok {
		t.Fatal("entry", "/file_finder_test.html", "already exists. Something is weird")
	}
	if _, ok := Resources["/file_finder_missing.html"]; ok {
		t.Fatal("entry", "/file_finder_missing.html", "already exists. Something is weird")
	}
	// Add something to Resources for this test to avoid fragility.
	Resources["/file_finder_test.html"] = "test content embedded"

	embr := makeEmbeddableResource(tmpdir)
	_, err = embr.GetAsString("/file_finder_missing.html")
	if err.(Error) != Error(ErrorNoSuchEmbeddedResource) {
		t.Fatal("access to missing resource", "/file_finder_missing.html", "ought to have failed but didn't", err)
	}

	got, err := embr.GetAsString("/file_finder_test.html")
	if want := "test content embedded"; err != nil || got != want {
		t.Fatalf("access to embedded resource failed. %v, got: %v, want: %v\n", err, got, want)
	}

	w, err := os.Create(filepath.Join(tmpdir, "/file_finder_test.html"))
	if err != nil {
		t.Fatalf("can't make %s in %s", "/file_finder_test.html", tmpdir)
	}
	if n, err := io.WriteString(w, "test content file"); err != nil || n != len("test content file") {
		t.Fatalf("can't write %s to %s", "test content file", filepath.Join(tmpdir, "/file_finder_test.html"))
	}

	got, err = embr.GetAsString("/file_finder_test.html")
	if want := "test content file"; err != nil || got != want {
		t.Fatalf("access to embedded resource failed. %v, got: %v, want: %v\n", err, got, want)
	}

}
