package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sfbrigade/sfsbook/setup"
)

//go:generate go run ../generator/tool/generateresources.go -output embedded_resources.go -prefix ../site/ ../site

// This file contains the code needed to either serve a resource from the
// embedded resources or find in the site directory. The above generate
// directory actually constructs the embedded resources using tooling
// from the generator package.

// Serve provides entry points that the file finder uses to actually
// serve content.
type Serve interface {
	// The desired content is available as a file.
	ServeForStream(reader io.Reader, w http.ResponseWriter, req interface{})
	ServeForString(s string, w http.ResponseWriter, req interface{}) 
}

// TODO(rjk): This is wrong. It should really be that the Global state is constructed by
// including a FileFinder instance. And then the global state would provide the
// methods of the FileFinder. And all of this should get relocated to server.
type FileFinder setup.GlobalState

func MakeFileFinder(global *setup.GlobalState) *FileFinder {
	// Types have to align. This isn't very good code.
	return (*FileFinder)(global)
}

// Given an Url, 
func (ff *FileFinder) StreamOrString(upath string, serve Serve, w http.ResponseWriter, req interface{}) error {
	fpath := filepath.Join(ff.Sitedir, upath)
	log.Println(upath, fpath)
		
	if _, err := os.Stat(fpath); err != nil {
		res, ok := Resources[upath]
		if !ok {
			log.Println("file", fpath, "missing from disk", err, "and also missing", upath, "from compiled in resource")
			return fmt.Errorf("file %s missing from site directory: %v and also not compiled in", upath, err)
		}
		// TODO(rjk): Revisit/rationalize the handling of errors.
		serve.ServeForString(res, w, req)
		return nil
	}

	f, err := os.Open(fpath)
	if err != nil {
		log.Println("problem opening file", fpath, err)
		return fmt.Errorf("file %s missing from site: %v", upath,  err)
	}
	defer f.Close()
	serve.ServeForStream(f, w, req)
	return nil
}
