package server

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
