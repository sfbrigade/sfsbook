package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// TODO(rjk): This will probably require additional fields.
type staticServer struct {
	s string
}

func MakeStaticServer(pathroot string) *staticServer {
	return &staticServer{filepath.Join(pathroot, "site")}
}

func (gs *staticServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	_, sn := filepath.Split(req.URL.Path)

	filepath := filepath.Join(gs.s, sn)

	log.Println(filepath)

	if _, err := os.Stat(filepath); err != nil {
		log.Println("file", filepath, "missing", err, "ought to be trying the wired-in content")
		// TODO(rjk): If the file doesn't exist, we default to using the content
		// compiled into the application.
		respondWithError(w, fmt.Sprintln("Can't open file: ", err))
		return
	}

	f, err := os.Open(filepath)
	defer f.Close()
	// TODO(rjk): I'm not sure that we necessarily want to expose our innards
	// like this. Perhaps I need to revise my opinions about how the error handling
	// should work.
	if err != nil {
		log.Println("problem opening file", filepath, err)
		respondWithError(w, fmt.Sprintln("Can't open file: ", err))
		return
	}

	if _, err := io.Copy(w, f); err != nil {
		log.Println("could not copy to the request body ", err)
		respondWithError(w, fmt.Sprintln("Can't copy: ", err))
		return
	}
	log.Println("finished. supposedly copied the content")
}
