package server

import (
	"context"
	"net/http"

	"github.com/pborman/uuid"
	"github.com/rjkroege/mocking"
)

const embeddedResourceForPasswdchg = `
	IsAuthed: {{.DecodedCookie.IsAuthed}}
	DisplayName: {{.DecodedCookie.DisplayName}}
	ChangeAttemptedAndFailed: {{.Results.ChangeAttemptedAndFailed}}
	ChangeAttemptedAndSucceeded: {{.Results.ChangeAttemptedAndSucceeded}}
	ReasonForFailure: {{.Results.ReasonForFailure}}
`
const embeddedResourceForListusers = `
	IsAuthed: {{.DecodedCookie.IsAuthed}}
	DisplayName: {{.DecodedCookie.DisplayName}}

	Userquery: {{.Results.Userquery}}
	Users: {{.Results.Users}}
	Querysuccess: {{.Results.Querysuccess}}
	Diagnosticmessage: {{.Results.Diagnosticmessage}}
`

type testPattern struct {
	urlargs string
	statuscode int
	tapeResponse []interface{}
	tapeRecord []interface{}
	outputString string	
}

// resourceHelper installs the above set of constant resources in place
// of the resources read from the site directory (or compiled in.)
// These templates make parseAndExecuteTemplate echo its
// template arguments to show that the middleware is
// generating the correct values.
func resourceHelper() func() {
	stashedResources := Resources

	Resources = map[string]string{
		"/usermgt/changepasswd.html": embeddedResourceForPasswdchg,
		"/usermgt/listusers.html":    embeddedResourceForListusers,
		"/head.html":                 "",
		"/header.html":               "",
		"/searchbar.html":            "",
		"/footer.html":               "",
	}

	return func() { Resources = stashedResources }
}

// addCookie augments req with the context data that specifies that the
// user is allowed to view users from the admin dialog.
func addCookie(req *http.Request, capability CapabilityType) *http.Request {
	uuid := uuid.NewRandom()
	// User does have the capability to view users.
	usercookie := &UserCookie{
		Uuid:        uuid,
		Capability:  capability,
		Displayname: "Homer Simpson",
	}
	return req.WithContext(context.WithValue(req.Context(),
		UserCookieStateName, usercookie))
}

// makeUnderTestHandlerListUsers creates a listUsers structure that
// uses a mock (tape based) implementation of PasswordIndex.
func makeUnderTestHandlerListUsers(tape *mocking.Tape) *listUsers {
	undertesthandler := &listUsers{
		// Always use the embedded resource.
		embr:         makeEmbeddableResource(""),
		passwordfile: (*mockPasswordIndex)(tape),
	}
	return undertesthandler
}

