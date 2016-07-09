package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//go:generate go run ../generator/tool/generateresources.go -output embedded_resources.go -prefix ../site/ ../site

// TODO(rjk): This will probably require additional fields.
type staticServer struct {
	s string

	// insert runtime processors here
}

// RuntimeProcessor provides entry points that the basic server uses to 
// process and return content. Different extensions can have different
//  implementations of the RuntimeProcessor. The simplest interface
// simply copies the file to the output or fails and is provided internally
// by staticServer itself.
// can (obviously) generate content dynamically.
type RuntimeProcessor interface {
	ServeStream(reader io.Reader, w http.ResponseWriter, req *http.Request)
	ServeString(s string, w http.ResponseWriter, req *http.Request) 
}

func MakeStaticServer(pathroot string) *staticServer {
	return &staticServer{filepath.Join(pathroot, "site")}
}

func (gs *staticServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	sn := req.URL.Path
	if sn == "" {
		sn = "index.html"
	}

	// I might generalize this too
	fpath := filepath.Join(gs.s, sn)

	// Filename-specific actions.
	var processor RuntimeProcessor
	switch path.Ext(sn) {
	case ".js":
		processor = gs
		w.Header().Add("Content-Type", "application/javascript")
	default:
		processor = gs
	}

	log.Println(sn, fpath)
	if _, err := os.Stat(fpath); err != nil {
		res, ok := Resources[sn]
		if !ok {
			log.Println("file", fpath, "missing from disk", err, "and also missing", sn, "from compiled in resource")
			respondWithError(w, fmt.Sprintln("Can't open file: ", err))
			return
		}
		processor.ServeString(res, w, req)
	} else {
		f, err := os.Open(fpath)
		// Set up processing chain for file here. There are two kinds of processing. Static
		// (i.e. processing that happens at generate time) and dynamic (templating) that
		// happen at serve time. Generation-time processing (i.e. JavaScript minification)
		// would happen here by running code provided by the generator's library.

		// TODO(rjk): I'm not sure that we necessarily want to expose our innards
		// in the 404 responses. Perhaps only if you're signed in.
		if err != nil {
			log.Println("problem opening file", fpath, err)
			respondWithError(w, fmt.Sprintln("Can't open file: ", err))
			return
		}
		defer f.Close()
		processor.ServeStream(f, w, req)
	}
	log.Println("finished. supposedly copied the content")
}

func (gs *staticServer) ServeString(s string, w http.ResponseWriter, req *http.Request) {
	reader := strings.NewReader(s)
	gs.ServeStream(reader, w, req)
}

func (gs *staticServer) ServeStream(reader io.Reader, w http.ResponseWriter, req *http.Request) {
	if _, err := io.Copy(w, reader); err != nil {
		log.Println("could not copy to the request body ", err)
		respondWithError(w, fmt.Sprintln("Can't copy: ", err))
		return
	}
}
