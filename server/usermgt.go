package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sfbrigade/sfsbook/dba"
	"golang.org/x/crypto/bcrypt"
)

type passwordChange struct {
	embr         *embeddableResources
	passwordfile dba.PasswordIndex
}

// makePasswdChangeHandler returns a handler for changing the user password.
func (hf *HandlerFactory) makePasswdChangeHandler() http.Handler {
	return &passwordChange{
		embr:         makeEmbeddableResource(hf.sitedir),
		passwordfile: hf.passwordfile,
	}
}

type passwordChangeResult struct {
	ChangeAttemptedAndSucceeded bool
	ChangeAttemptedAndFailed    bool
	ReasonForFailure            string
}

// TODO(rjk): Note refactoring opportunity with loginServer.
func (gs *passwordChange) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	changeresult := new(passwordChangeResult)

	// Skip immediately to end if not authed.
	if !GetCookie(req).IsAuthed() {
		goto end
	}

	if req.Method == "POST" {
		log.Println("handling Post of resource in usermgt")
		// TODO(rjk): rational logging of failed attempts.

		if err := req.ParseForm(); err != nil {
			respondWithError(w, fmt.Sprintln("bad uploaded form data to login: ", err))
			return
		}

		// We only use the posted info.
		// This dumps passwords in the clear in the log.
		// TODO(rjk): delete before landing this code.
		log.Println("parsing form:")
		for k, v := range req.PostForm {
			log.Println("	", k, "		", v)
		}

		// TODO: must make sure that we count how many times the UA hammers
		// on us. How do I track this if they clear the cookies? I presume that we
		// do it by IP. Too many requests (of any kind) from any one IP need to
		// get the connection dropped.

		oldpassword, err := getValidatedString("oldpassword", req.PostForm)
		if err != nil {
			log.Println("no oldpassword in Post", err)
			changeresult.ChangeAttemptedAndFailed = true
			changeresult.ReasonForFailure = "Need to enter the previous password"
			goto end
		}
		newpassword, err := getValidatedString("newpassword", req.PostForm)
		if err != nil {
			log.Println("no newpassword in form")
			changeresult.ChangeAttemptedAndFailed = true
			changeresult.ReasonForFailure = "Need to enter a new password"
			goto end
		}
		newpasswordagain, err := getValidatedString("newpasswordagain", req.PostForm)
		if err != nil {
			log.Println("no newpassword in form")
			changeresult.ChangeAttemptedAndFailed = true
			changeresult.ReasonForFailure = "Need to enter the new password overagain."
			goto end
		}

		if newpassword != newpasswordagain {
			log.Println("newpasswords don't match")
			changeresult.ChangeAttemptedAndFailed = true
			changeresult.ReasonForFailure = "New password fields need to match."
			goto end
		}

		// TODO(rjk): Invoke a more sophisticated policy for validating that
		// the password is useful.
		if len(newpassword) < 5 {
			log.Println("newpassword is weak")
			changeresult.ChangeAttemptedAndFailed = true
			changeresult.ReasonForFailure = "New password is too easily guessed."
			goto end
		}

		uuid := GetCookie(req).Uuid
		suuid := string(uuid)
		// TODO(rjk): Would be nice if the bleve could use []byte as ids.
		sr, err := gs.passwordfile.MapForDocument(suuid)
		if err != nil {
			// This signifies something perturbing: the uuid exists but is not in
			// the DB. The most likely cause is that the user has been deleted.
			// TODO(rjk): Add uuid to revocation list. This would force the user
			// to attempt to sign-in. Revocation must be server side so that it
			// cannot be tampered with.
			log.Println("oops! uuid in sesssion cookie not in database")
			changeresult.ChangeAttemptedAndFailed = true
			changeresult.ReasonForFailure = "Account error. Please sign-out and sign-in again."
			goto end
		}

		// TODO(rjk): I should test this?
		pw := sr["passwordhash"].(string)
		if err := bcrypt.CompareHashAndPassword([]byte(pw), []byte(oldpassword)); err != nil {
			// This means that the user has entered an invalid password.
			// It's ok to tell the user this because the user has already signed in
			// once.
			// TODO(rjk): Count how many times we've attempted it. Too many
			// attempts will a) revoke cookie and b) temporaily ban the user.
			changeresult.ChangeAttemptedAndFailed = true
			changeresult.ReasonForFailure = "Old password is incorrect"
			goto end
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.DefaultCost)
		if err != nil {
			log.Println("bcrypt.GenerateFromPassword failed:", err)
			respondWithError(w, fmt.Sprintln("bcrypt.GenerateFromPassword failed:", err))
			return
		}
		sr["passwordhash"] = string(hash)

		if err := gs.passwordfile.Index(string(suuid), sr); err != nil {
			log.Println("passwordfile.Index failed", err)
			respondWithError(w, fmt.Sprintln("passwordfile.Index failed", err))
			return
		}
		changeresult.ChangeAttemptedAndSucceeded = true
	}

end:

	sn := req.URL.Path
	templates := []string{sn}

	parseAndExecuteTemplate(gs.embr, w, req, templates, changeresult)
}
