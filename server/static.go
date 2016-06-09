package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:generate go run ../generator/generateresources.go ../sites

// TODO(rjk): This will probably require additional fields.
type staticServer struct {
	s string
}

func MakeStaticServer(pathroot string) *staticServer {
	return &staticServer{filepath.Join(pathroot, "site")}
}

func (gs *staticServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO(rjk): This doesn't handle sub-directories correctly.
	_, sn := filepath.Split(req.URL.Path)
	if sn == "" {
		sn = "index.html"
	}
	filepath := filepath.Join(gs.s, sn)

	var resreader  io.Reader

	log.Println(filepath)
	if _, err := os.Stat(filepath); err != nil {
		res, ok := Resources[sn]
		if !ok {
			log.Println("file", filepath, "missing from disk", err, "and also missing", sn, "from compiled in resource")
			respondWithError(w, fmt.Sprintln("Can't open file: ", err))
			return
		}
	
		resreader = strings.NewReader(res)
	} else {
		f, err := os.Open(filepath)
		// Set up processing chain for file here. There are two kinds of processing. Static
		// (i.e. processing that happens at generate time) and dynamic (templating) that
		// happen at serve time.

		// TODO(rjk): I'm not sure that we necessarily want to expose our innards
		// like this. Perhaps I need to revise my opinions about how the error handling
		// should work.
		if err != nil {
			log.Println("problem opening file", filepath, err)
			respondWithError(w, fmt.Sprintln("Can't open file: ", err))
			return
		}
		defer f.Close()
		resreader = f		
	}

	// What about templating? Or generative content? We need to insert a phase here
	// here (some refactoring required) that can re-process the file.

	// Processing of the content has to happen here.	
	if _, err := io.Copy(w, resreader); err != nil {
		log.Println("could not copy to the request body ", err)
		respondWithError(w, fmt.Sprintln("Can't copy: ", err))
		return
	}
	log.Println("finished. supposedly copied the content")
}
