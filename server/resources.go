package server
// This module of package server is responsible for processing resuts for
// a specific named resource.

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/sfbrigade/sfsbook/dba"
)

type resourceServer struct {
	templatedServer
}

func MakeResourceServer(ff *fileFinder, g dba.Generator) *resourceServer {
	return &resourceServer{ 
		templatedServer: *MakeTemplatedServer(ff, g),
	}
}


func (gs *resourceServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	sn := req.URL.Path

	log.Println("ResourceServer is handling request for", sn)

	// Path is expected to be of the form /resources/<uuid>.html
	if path.Ext(sn) != ".html" {
		respondWithError(w, "bad extension: " + path.Ext(sn))
	}

	// Re-use req's payload.
	uuid := strings.TrimSuffix(path.Base(sn), path.Ext(sn))
	sn = "/resources/resource.html"

	if err := gs.ff.StreamOrString(sn, gs, w, uuid); err != nil {
		respondWithError(w, fmt.Sprintln("Server error", err))
	}
}



