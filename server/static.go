package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
)

type staticServer GlobalState


func MakeStaticServer(global *GlobalState) *staticServer {
	return (*staticServer)(global)
}

func (gs *staticServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO(rjk): Do auth lookup here.

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

	// TODO(rjk): Auth validation here.
	// TODO(rjk): Figure out how I describe the auth requirements.

	if n, err := io.WriteString(w, str); err != nil || n != len(str) {
		log.Println("couldn't write string to ResponseWriter, wrote",
			n, "of", len(str), "or received error", err)
	}
}
