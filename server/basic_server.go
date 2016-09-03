package server

import (
	"net/http"
	"time"

	"github.com/sfbrigade/sfsbook/dba"
)

// MakeServer creates a Server serving from the specified address.
// The contents of pathroot are served.
// Conceivably, it's possible that passing the bi through here is a layering violation?
// TODO(rjk): I'm convinced, it's a layering violation. Make it go away.
// TODO(rjk): redirect to from http to https.
func MakeServer(address string, global *GlobalState) *http.Server {
	m := http.NewServeMux()

	m.Handle("/js/", MakeStaticServer(global))
	m.Handle("/resources/", MakeResourceServer(global, dba.MakeResourceResultsGenerator(global.ResourceGuide)))
	m.Handle("/search.html", MakeTemplatedServer(global, dba.MakeQueryResultsGenerator(global.ResourceGuide)))
	m.Handle("/", MakeTemplatedServer(global, dba.MakeStubGenerator(global.ResourceGuide)))

	// TODO(rjk): why no https config here?
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Addr:         address,
		Handler:      m,

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
