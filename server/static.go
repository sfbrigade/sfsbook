package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
)


// TODO(rjk): This will probably require additional fields.
type staticServer struct {
	ff *fileFinder
}



func MakeStaticServer(ff *fileFinder) *staticServer {
	return &staticServer{ff: ff}
}


func (gs *staticServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	sn := req.URL.Path

	// Filename-specific actions.
	switch path.Ext(sn) {
	case ".js":
		w.Header().Add("Content-Type", "application/javascript")
	}

	if err := gs.ff.StreamOrString(sn, gs, w, req); err != nil {
		respondWithError(w, fmt.Sprintln("Server error", err))
	}
}

func (gs *staticServer) ServeForString(s string, w http.ResponseWriter, req interface{}) {
	reader := strings.NewReader(s)
	gs.ServeForStream(reader, w, req)
}

func (gs *staticServer) ServeForStream(reader io.Reader, w http.ResponseWriter, req interface{}) {
	if _, err := io.Copy(w, reader); err != nil {
		log.Println("could not copy to the request body ", err)
		respondWithError(w, fmt.Sprintln("Can't copy: ", err))
		return
	}
}
