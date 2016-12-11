package server

import (
	"net/http"
	"time"

	"github.com/sfbrigade/sfsbook/dba"
	"github.com/sfbrigade/sfsbook/setup"
)

// MakeServer creates a Server serving from the specified address.
// The contents of pathroot are served.
// Conceivably, it's possible that passing the bi through here is a layering violation?
// TODO(rjk): I'm convinced, it's a layering violation. Make it go away.
// TODO(rjk): redirect to from http to https.
func MakeServer(address string, hf *HandlerFactory, cf *setup.CertFactory) *http.Server {
	m := http.NewServeMux()

	m.Handle("/js/", hf.makeCookieHandler(hf.makeStaticHandler()))
	m.Handle("/css/", hf.makeCookieHandler(hf.makeStaticHandler()))

	m.Handle("/resources/",
		hf.makeCookieHandler(
			hf.makeResourceHandler(dba.MakeResourceResultsGenerator(hf.resourceguide))))

	m.Handle("/search.html",
		hf.makeCookieHandler(
			hf.makeTemplatedHandler(dba.MakeQueryResultsGenerator(hf.resourceguide))))

	// TODO(rjk): need to wire up the login data
	// Having the cookie data lets me handle the situation of someone navigating here
	// in error.
	m.Handle("/login.html",
		hf.makeCookieHandler(hf.makeLoginHandler()))

	m.Handle("/usermgt/changepasswd.html",
		hf.makeCookieHandler(hf.makePasswdChangeHandler()))
	m.Handle("/usermgt/listusers.html",
		hf.makeCookieHandler(hf.makeListUsersHandler()))
	m.Handle("/usermgt/adduser.html",
		hf.makeCookieHandler(hf.makeAddUserHandler()))

	m.Handle("/",
		hf.makeCookieHandler(
			hf.makeTemplatedHandler(dba.MakeStubGenerator(hf.resourceguide))))

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Addr:         address,
		Handler:      m,
		TLSConfig:    cf.GetTLSConfig(),
	}
	return srv
}

// helper function. Re-write me.
func respondWithError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(message))
}
