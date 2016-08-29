package server

import (
	"net/http"
	"path/filepath"
)

func Start(pathroot string, srv *http.Server) error {
	return srv.ListenAndServeTLS(
		filepath.Join(pathroot, "state", "cert.pem"),
		filepath.Join(pathroot, "state", "key.pem"))
}
