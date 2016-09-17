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

type resourceServer templatedServer

func MakeResourceServer(global *HandlerFactory, g dba.Generator) *resourceServer {
	return (*resourceServer)(makeTemplatedHandler(global, g))
}

func (gs *resourceServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	sn := req.URL.Path

	log.Println("ResourceServer is handling request for", sn)

	// Path is expected to be of the form /resources/<uuid>.html
	if path.Ext(sn) != ".html" {
		respondWithError(w, "bad extension: "+path.Ext(sn))
		return
	}

	// Re-use req's payload.
	uuid := strings.TrimSuffix(path.Base(sn), path.Ext(sn))
	sn = "/resources/resource.html"

	// TODO(rjk): Validate the uuid here and error-out if it's non-sensical.

	dbreq := &dba.ResourceRequest{
		Uuid: uuid,
	}

	if req.Method == "POST" {
		log.Println("handling Post of resource")
		if err := req.ParseForm(); err != nil {
			respondWithError(w, fmt.Sprintln("bad uploaded form data: ", err))
			return
		}

		// We only use the posted info.
		log.Println("parsing form:")
		for k, v := range req.PostForm {
			log.Println("	", k, "		", v)
		}

		dbreq.IsPost = true
		dbreq.PostArgs = req.PostForm
	}

	str, err := gs.embr.GetAsString(sn)
	if err != nil {
		// TODO(rjk): Rationalize error handling here. There needs to be a 404 page.
		respondWithError(w, fmt.Sprintln("Server error", err))
	}

	((*templatedServer)(gs)).ServeForStrings(str, w, dbreq)
}
