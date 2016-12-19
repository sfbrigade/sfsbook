package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pborman/uuid"
	"github.com/sfbrigade/sfsbook/dba"
	"golang.org/x/crypto/bcrypt"
)

type addUser struct {
	embr         *embeddableResources
	passwordfile dba.PasswordIndex
}

// makeAddUserHandler returns a handler for editing user data.
func (hf *HandlerFactory) makeAddUserHandler() http.Handler {
	return &addUser{
		embr:         makeEmbeddableResource(hf.sitedir),
		passwordfile: hf.passwordfile,
	}
}

type addUserResult struct {
	Useradded        bool
	Usernotadded     bool
	ReasonForFailure string
	Username         string
}

func (gs *addUser) ender(w http.ResponseWriter, req *http.Request, listusersresult interface{}) {
	sn := req.URL.Path
	templates := []string{sn, "/head.html", "/header.html","/searchbar.html", "/footer.html"}
	parseAndExecuteTemplate(gs.embr, w, req, templates, listusersresult)
}

// TODO(rjk): Note refactoring opportunity with resource search.
func (gs *addUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("ServeHTTP, addUser")

	adduserresult := new(addUserResult)

	// Skip immediately to end if not authed.
	if !GetCookie(req).HasCapabilityIniviteUsers() {
		gs.ender(w, req, adduserresult)
		return
	}

	if req.Method == "POST" {
		log.Println("handling Post of added user")
		if err := req.ParseForm(); err != nil {
			respondWithError(w, fmt.Sprintln("bad uploaded form data: ", err))
			return
		}

		// We only use the posted info.
		log.Println("parsing form:")
		for k, v := range req.PostForm {
			log.Println("	", k, "		", v)
		}

		// TODO(rjk): Could be generalized nicely with a loop. With password
		// validation added separately.
		username, err := getValidatedString("username", req.PostForm)
		if err != nil {
			log.Println("no username in Post", err)
			adduserresult.Usernotadded = true
			adduserresult.ReasonForFailure = "Need to enter a username"
			goto end
		}
		displayname, err := getValidatedString("displayname", req.PostForm)
		if err != nil {
			log.Println("no displayname in Post", err)
			adduserresult.Usernotadded = true
			adduserresult.ReasonForFailure = "Need to enter a display name for user"
			goto end
		}
		role, err := getValidatedString("role", req.PostForm)
		if err != nil {
			log.Println("no role in Post", err)
			adduserresult.Usernotadded = true
			adduserresult.ReasonForFailure = "Need to enter a display name for user"
			goto end
		}

		// TODO(rjk): Refactor with usemgt
		newpassword, err := getValidatedString("newpassword", req.PostForm)
		if err != nil {
			log.Println("no newpassword in form")
			adduserresult.Usernotadded = true
			adduserresult.ReasonForFailure = "Need to enter a new password"
			goto end
		}
		newpasswordagain, err := getValidatedString("newpasswordagain", req.PostForm)
		if err != nil {
			log.Println("no newpassword in form")
			adduserresult.Usernotadded = true
			adduserresult.ReasonForFailure = "Need to enter the new password overagain."
			goto end
		}

		if newpassword != newpasswordagain {
			log.Println("newpasswords don't match")
			adduserresult.Usernotadded = true
			adduserresult.ReasonForFailure = "New password fields need to match."
			goto end
		}

		// TODO(rjk): Invoke a more sophisticated policy for validating that
		// the password is useful.
		if len(newpassword) < 5 {
			log.Println("newpassword is weak")
			adduserresult.Usernotadded = true
			adduserresult.ReasonForFailure = "New password is too easily guessed."
			goto end
		}

		log.Println("data", newpassword, displayname, username, role)
		sr := map[string]interface{}{
			"name":         username,
			"cost":         string(bcrypt.DefaultCost),
			"role":         role,
			"display_name": displayname,
			"_type":        "password",
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.DefaultCost)
		if err != nil {
			log.Println("bcrypt.GenerateFromPassword failed:", err)
			respondWithError(w, fmt.Sprintln("bcrypt.GenerateFromPassword failed:", err))
			return
		}
		sr["passwordhash"] = string(hash)

		suuid := uuid.NewRandom()
		if err := gs.passwordfile.Index(string(suuid), sr); err != nil {
			log.Println("passwordfile.Index failed", err)
			adduserresult.Usernotadded = true
			adduserresult.ReasonForFailure = "Failed to index new user into password store"
		}
		adduserresult.Useradded = true
		adduserresult.Username = displayname
	}

end:
	gs.ender(w, req, adduserresult)
}
