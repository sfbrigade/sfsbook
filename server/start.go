package server

import (
	"path/filepath"

	"gopkg.in/tylerb/graceful.v1"
)

func Start(pathroot string, srv *graceful.Server) error {
	return srv.ListenAndServeTLS(
		filepath.Join(pathroot, "state", "cert.pem"),
		filepath.Join(pathroot, "state", "key.pem"))
}
