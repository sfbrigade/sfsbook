package server

import (
	"log"
	"net/http"
	"time"

	"gopkg.in/tylerb/graceful.v1"
)

// MakeServer creates a graceful.Server serving from the specificed address.
// The contents of pathroot are served.
func MakeServer(address, pathroot string) *graceful.Server {
	log.Println("MakeServer", address, pathroot)
	m := http.NewServeMux()

	// TODO(rjk): Expand appropriately.
	m.Handle("/", MakeStaticServer(pathroot))

	srv := &graceful.Server{
		Timeout: 5 * time.Second,
		Server: &http.Server{
			Addr:    address,
			Handler: m,
		},
	}
	return srv
}

// helper function. Re-write me.
func respondWithError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(message))
}
