package server

import (
	"log"
	"net/http"
	"time"

	"github.com/sfbrigade/sfsbook/dba"
	"github.com/sfbrigade/sfsbook/setup"

)

// MakeServer creates a Server serving from the specificed address.
// The contents of pathroot are served.
// Conceivably, it's possible that passing the bi through here is a layering violation?
// TODO(rjk): I'm convinced, it's a layering violation. Make it go away.
// TODO(rjk): redirect to from http to https.
func MakeServer(address string, global *setup.GlobalState) *http.Server {
	log.Println("MakeServer", address, global.Sitedir)
	m := http.NewServeMux()

	// Have I chained this in the right direction?
	// i.e.: why is the file-finder at the bottom?
	// Because... it has to be?
	// File-finder is inside the server.

	ff := MakeFileFinder(global)
	m.Handle("/js/", MakeStaticServer(ff))
	m.Handle("/resources/", MakeResourceServer(ff, dba.MakeResourceResultsGenerator(global.ResourceGuide)))
	m.Handle("/search.html", MakeTemplatedServer(ff, dba.MakeQueryResultsGenerator(global.ResourceGuide)))
	m.Handle("/", MakeTemplatedServer(ff, dba.MakeStubGenerator(global.ResourceGuide)))

	// TODO(rjk): why no https config here?
	srv := &http.Server{
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Addr:    address,
		Handler: m,

		// TLS config?
		
	}
	return srv
}

// helper function. Re-write me.
func respondWithError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(message))
}
