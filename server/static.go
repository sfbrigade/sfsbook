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
type staticServer GlobalState


func MakeStaticServer(global *GlobalState) *staticServer {
	return (*staticServer)(global)
}

func (gs *staticServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	sn := req.URL.Path

	// Filename-specific actions.
	switch path.Ext(sn) {
	case ".js":
		w.Header().Add("Content-Type", "application/javascript")
	}

	str, err := gs.GetAsString(sn)
	if err != nil {
		// TODO(rjk): Rationalize error handling here. There needs to be a 404 page.
		respondWithError(w, fmt.Sprintln("Server error", err))
	}

	// TODO(rjk): Refactor this.
	gs.ServeForString(str, w, req)
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
