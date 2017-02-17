package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/gorilla/securecookie"
	"github.com/pborman/uuid"
	"github.com/sfbrigade/sfsbook/dba"
	"golang.org/x/crypto/bcrypt"
)

// This is written possibly incorrectly. Refactor later.
// I need more tests and more refactoring. my code arrangement leaves
// much to be desired.
type loginServer struct {
	embr         *embeddableResources
	passwordfile dba.PasswordIndex
	cookiecodec  *securecookie.SecureCookie
}

// makeLoginHandler returns a handler for login actions.
func (hf *HandlerFactory) makeLoginHandler() http.Handler {
	return &loginServer{
		embr:         makeEmbeddableResource(hf.sitedir),
		passwordfile: hf.passwordfile,
		cookiecodec:  hf.cookiecodec,
	}
}

// getValidatedString returns a string from the postform or an error
// if there's something wrong with it.
func getValidatedString(key string, postform url.Values) (string, error) {
	if len(postform[key]) != 1 {
		return "", fmt.Errorf("key %s in POST data is invalid", key)
	}
	value := postform[key][0]

	// Values need non-zero length.
	if len(value) == 0 {
		return "", fmt.Errorf("value for key %s in POST data is of 0 length", key)
	}
	return value, nil
}

// TODO(rjk): It is conceivable that this could be computed from
// the cookie state and this code could be simplified.
type loginResult struct {
	ValidSignOne    bool
	AttemptedSignOn bool
}

func (gs *loginServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	sn := req.URL.Path

	log.Println("loginServer is handling request for", sn)

	loginresult := new(loginResult)

	if req.Method == "POST" {
		loginresult.AttemptedSignOn = true
		log.Println("handling Post of resource")

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

		username, err := getValidatedString("username", req.PostForm)
		if err != nil {
			log.Println("no username in Post", err)
			goto end
		}
		password, err := getValidatedString("password", req.PostForm)
		if err != nil {
			log.Println("no password in form")
			goto end
		}

		// Do database access.
		sreq := bleve.NewSearchRequest(bleve.NewMatchQuery(username))
		sreq.Fields = []string{"name", "cost", "passwordhash", "role", "display_name"}

		// This is an error case (something is wrong internally)
		searchResults, err := gs.passwordfile.Search(sreq)
		if err != nil {
			respondWithError(w, fmt.Sprintln("database couldn't respond with useful results", err))
		}

		if len(searchResults.Hits) != 1 {
			// This means that the user has entered an invalid username. But we don't
			// tell the UA this.
			log.Println("username mismatch")
			goto end
		}

		sr := searchResults.Hits[0]
		pw := sr.Fields["passwordhash"].(string)
		if err := bcrypt.CompareHashAndPassword([]byte(pw), []byte(password)); err != nil {
			// This means that the user has entered an invalid password.
			// But we don't tell the UA this either.
			log.Println("password mismatch")
			goto end
		}

		// User has successfully signed on
		log.Println("username: ", username, "has signed in")
		loginresult.ValidSignOne = true

		// TODO(rjk): force updating of password if DefaultCost has changed

		// Build the cookie.
		role := sr.Fields["role"].(string)

		// We're downstream of the cookieHandler and so already have a
		// usercookie. We've signed in successfully. So augment it. That
		// way, all downstream code will have the correct context.
		usercookie := GetCookie(req)
		usercookie.Uuid = uuid.UUID(sr.ID)
		usercookie.Displayname = sr.Fields["display_name"].(string)
		usercookie.Timestamp = time.Now()

		// TODO(rjk): Consider storing the capability in the user data record.
		switch role {
		case "admin":
			usercookie.Capability = CapabilityAdministrator
		case "volunteer":
			usercookie.Capability = CapabilityVolunteer
		default:
			usercookie.Capability = CapabilityViewPublicResourceEntry
		}

		log.Println("usercookie", usercookie)

		if encoded, err := gs.cookiecodec.Encode(SessionCookieName, usercookie); err == nil {
			cookie := &http.Cookie{
				Name:  SessionCookieName,
				Value: encoded,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
		} else {
			// I'm not sure what to do here.
			// I think this means that the user can't have a cookie.
			// i.e. we make a sad page.
			respondWithError(w, fmt.Sprintln("Server cookie error", err))
		}
	}

end:

	templates := []string{sn, "/header.html", "/footer.html"}
	// do the redirect?
	parseAndExecuteTemplate(gs.embr, w, req, templates, loginresult)
}
