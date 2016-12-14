package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pborman/uuid"
	"github.com/sfbrigade/sfsbook/dba"
)

const (
	_ = iota
	_SEARCHACTION
	_RESETPASSWORD
	_DELETEUSERS
	_ROLECHANGE_TO_ADMIN
	_ROLECHANGE_TO_VOLUNTEER
	_ROLECHANGE_TO_NOROLE
	_BADACTORREQUEST
)

type listUsers struct {
	embr         *embeddableResources
	passwordfile dba.PasswordIndex
}

// makeListUsersHandler returns a handler for editing user data.
func (hf *HandlerFactory) makeListUsersHandler() http.Handler {
	return &listUsers{
		embr:         makeEmbeddableResource(hf.sitedir),
		passwordfile: hf.passwordfile,
	}
}

type listUsersResult struct {
	Userquery string
	// TODO(rjk): Consider making this typed in some fashion.
	Users             []map[string]interface{}
	Querysuccess      bool
	Diagnosticmessage string
}

func (gs *listUsers) ender(w http.ResponseWriter, req *http.Request, listusersresult interface{}) {
	sn := req.URL.Path
	templates := []string{sn, "/head.html", "/header.html", "/footer.html"}
	parseAndExecuteTemplate(gs.embr, w, req, templates, listusersresult)
}

// deleteUsers deletes the specified users. Attempts all and returns
// the last failure if any deletion failed.
func (gs *listUsers) deleteUsers(uuids []uuid.UUID) error {
	var err error
	for _, u := range uuids {
		// TODO(rjk): Revoke the cookies associated with this user.
		err = gs.passwordfile.Delete(string(u))
	}
	return err
}

// rolechangeUsers changes the specifed user role. Attempts all
// and returns the last failure if any role change failed.
// TODO(rjk): This might not be the best policy. Adjust this
func (gs *listUsers) rolechangeUsers(uuids []uuid.UUID, newrole string) error {
	var lasterror error
	for _, u := range uuids {
		suuid := string(u)
		passwordmap, err := gs.passwordfile.MapForDocument(suuid)
		if err != nil {
			// oops. Probably the article was deleted out from under us.
			// Just move on.
			log.Println("error reading user entry for", suuid, err)
			lasterror = err
			continue
		}
		// Have a map...
		log.Println("got a map", passwordmap)

		passwordmap["role"] = newrole

		log.Println(">> transformed map", passwordmap)
		// TODO(rjk): Revoke the cookies associated with this user.
		if err := gs.passwordfile.Index(suuid, passwordmap); err != nil {
			lasterror = err
			log.Println("error writing user entry with updated roll", suuid, err)
		}
	}
	return lasterror
}

// TODO(rjk): Note refactoring opportunity with resource search.
func (gs *listUsers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("listUsers.ServeHTTP", req)

	listusersresult := new(listUsersResult)

	// Skip immediately to end if not authed.
	if !GetCookie(req).HasCapability(CapabilityViewUsers) {
		listusersresult.Diagnosticmessage = "Sign in as an admin to list users."
		gs.ender(w, req, listusersresult)
		return
	}

	// Setup a query. The query is different if we have specified it.
	userquery := ""
	log.Println("not a post but proceeding anyway, req:", req)

	if err := req.ParseForm(); err != nil {
		respondWithError(w, fmt.Sprintln("bad uploaded form data to login: ", err))
		return
	}

	// Remember that nothing provided as part of the form can be
	// trusted to be valid. All must be validated. I need a general
	// mechanism to address clients that are broken.
	selecteduuids := make([]uuid.UUID, 0, 10)
	action := _SEARCHACTION
httpargloop:
	for k, v := range req.Form {
		log.Println("listusers ServeHTTP processing form item:", k, v, "current action", action)
		switch {
		case strings.HasPrefix(k, "selected-"):
			uuidstring, err := getValidatedString(k, req.Form)
			if err != nil {
				log.Println("posted form contained invalid k,v pair:", k, v, err)
				action = _BADACTORREQUEST
				break httpargloop
			}
			uuid := uuid.Parse(uuidstring)
			if uuid == nil {
				log.Println("posted form contained invalid uuid", v)
				action = _BADACTORREQUEST
				break httpargloop
			}
			selecteduuids = append(selecteduuids, uuid)
		case k == "rolechange":
			// parse the desired rolechange
			rc := v[0]
			switch rc {
			case "nochange":
				// Don't have to do anything
			case "admin":
				action = _ROLECHANGE_TO_ADMIN
			case "norole":
				action = _ROLECHANGE_TO_NOROLE
			case "volunteer":
				action = _ROLECHANGE_TO_VOLUNTEER
			default:
				log.Println("posted form contained invalid rolechange", k, v)
				action = _BADACTORREQUEST
				break httpargloop
			}
		case k == "deleteuser":
			action = _DELETEUSERS
		case k == "userquery":
			uq, err := getValidatedString("userquery", req.Form)
			if err != nil {
				log.Println("no userquery", err)
				listusersresult.Diagnosticmessage = "Showing all..."
			} else {
				// We had an argument to search with.
				listusersresult.Userquery = uq
				userquery = uq
				// TODO(rjk): Improve database indexing.
			}
		case k == "resetpassword":
			action = _RESETPASSWORD
			// I don't currently have enough information in the password
			// database to do this. And adding the additional data may
			// require interaction with how we implement oauth integration.

		default:
			// TODO(rjk): This is a convenience to simplify implementation.
			log.Println("ignoring extra fields in form", k, v)
		}
	}

	if !GetCookie(req).HasCapability(CapabilityEditUsers) && action != _SEARCHACTION && action != _BADACTORREQUEST {
		listusersresult.Diagnosticmessage = "Sign in as an admin to edit users."
		gs.ender(w, req, listusersresult)
		return
	}

	// Do the updates
	switch action {
	case _SEARCHACTION:
		// Or do nothing if only searching.
	case _RESETPASSWORD:
		// I have no idea how this should work. And I'm not going to build it
		// until we have discussed with SFWAR.
		log.Println("notimplemented: resetpassword applied to", selecteduuids)
	case _DELETEUSERS:
		if err := gs.deleteUsers(selecteduuids); err != nil {
			log.Println("failed to delete selected users", selecteduuids, err)
			listusersresult.Diagnosticmessage = "Couldn't successfully delete all of the selected users."
		}
	case _ROLECHANGE_TO_ADMIN:
		if err := gs.rolechangeUsers(selecteduuids, "admin"); err != nil {
			log.Println("failed to rolechange selected users to", "admin", selecteduuids, err)
			listusersresult.Diagnosticmessage = "Couldn't successfully rolechange all of the selected users to admin."
		}
	case _ROLECHANGE_TO_VOLUNTEER:
		if err := gs.rolechangeUsers(selecteduuids, "volunteer"); err != nil {
			log.Println("failed to rolechange selected users to", "volunteer", selecteduuids, err)
			listusersresult.Diagnosticmessage = "Couldn't successfully rolechange all of the selected users to volunteer."
		}
	case _ROLECHANGE_TO_NOROLE:
		// Conceivably, we might want to have a separation between a user with
		// no role and no such user. But I don't see why at the moment.
		if err := gs.deleteUsers(selecteduuids); err != nil {
			log.Println("failed to rolechange selected users to nothing", selecteduuids, err)
			listusersresult.Diagnosticmessage = "Couldn't successfully rolechange all of the selected users to norole."
		}
	default: // includes _BADACTORREQUEST
		respondWithError(w, "client is attempting something wrong")
		return
	}

	users, err := gs.passwordfile.ListUsers(userquery, 10, 0)
	if err != nil {
		respondWithError(w, fmt.Sprintln("database couldn't respond with useful results", err))
		return
	}
	if len(users) < 1 {
		// This probably means that the user has entered an invalid query.
		// Or that we've paged past the end.
		listusersresult.Diagnosticmessage = "Userquery matches no users."
		gs.ender(w, req, listusersresult)
		return
	}

	listusersresult.Querysuccess = true
	listusersresult.Users = users

	gs.ender(w, req, listusersresult)
}

