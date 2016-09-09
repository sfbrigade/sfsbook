package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"	
	"time"

	"github.com/pborman/uuid"
	"github.com/gorilla/securecookie"

)

// Capability is a set of capabilities. Sized because the 
type Capability int64

const (
	// Finding this Capability in a UserCookie implies that the uuid is empty.
	CapabilityAnonymous Capability = 0
)

// Constant Capability makes a bitmask.
const (
	// This is also the default.
	CapabilityViewPublicResourceEntry Capability = 1 << iota
	CapabilityViewOwnVolunteerComment
	CapabilityViewOtherVolunteerComment

	// Edit includes adding or removing.
	CapabilityEditOwnVolunteerComment
	CapabilityEditOtherVolunteerComment

	CapabilityEditResource

	CapabilityViewUsers	
	CapabilityInivteNewUser
	CapabilityEditUsers

	// The user has been altered. Finding this key in
	// a cookie suggests an altered auth flow when
	// redirected to the authentication page.
	CapabilityReauthenticate
)


// UserCookie is encrypted via securecookie facilities
// and set as a cookie on the interaction with the remote UA.
type UserCookie struct {
	// The user identifier. 
	uuid uuid.UUID

	// A mask of capabilities.
	capability Capability
	
	// The time that the cookie was created.
	timestamp time.Time
}


// TODO(rjk): Add the ability to check that a given uuid needs to be
// revalidated.

type UserState struct {
	securecookie.SecureCookie
	revokelist []uuid.UUID
}

// makeCookie builds and saves a cookie.
// TODO(rjk): Add automatic cookie rotation with aging and batches.
func makeCookie(statepath, cookiename string) ([]byte, error) {
	path := filepath.Join(statepath, cookiename)
	key, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("making new cookie", cookiename)
		key := securecookie.GenerateRandomKey(32)
		if key == nil {
			return nil, fmt.Errorf("No cookie for %s and can't make one", cookiename)
		}

		// TODO(rjk): Make sure that the umask is set appropriately.
		cookiefile, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("Can't create a %s to hold new cookie: %v",
				path, err)
		}

		if n, err := cookiefile.Write(key); err != nil || n != len(key) {
			return nil, fmt.Errorf("Can't write new cookie %s.  len is %d instead of %d or error: %v",
				path, n, len(key), err)
		}
	}
	return key, nil
}

// MakeUserState builds an instance of the user mangement facility.
func MakeUserState(statepath string) (*UserState, error) {
	// Make cookie keys.
	hashKey, err := makeCookie(statepath, "hashkey.dat")
	if err != nil {
		return nil, err
	}
	blockKey, err := makeCookie(statepath, "blockkey.dat")
	if err != nil {
		return nil, err
	}

	return &UserState{
		SecureCookie:        *securecookie.New(hashKey, blockKey),
		revokelist: make([]uuid.UUID, 0, 10),
	}, nil
}


const SessionCookieName = "session"
const COOKIE = 1

// WithDecodedUserCookie updates the given http.Request with a decoded
// instance of the cookie or updates the response to redirect
// appropriately and returns true if the req was successfully udpated.
// TODO(rjk): Fix up this comment.
func (um *UserState) WithDecodedUserCookie(w http.ResponseWriter, req *http.Request)  (*http.Request, bool) {
	cookie, err := req.Cookie(SessionCookieName)
	usercookie := new(UserCookie)
	if err == nil {
		// This request has a cookie.
		if err = um.Decode(SessionCookieName, cookie.Value, usercookie); err != nil {
			log.Println("request had a cookie but it was not decodeable:", err)
			// redirect to the login page with an appropriate error message.
		}
		log.Println("request had a cookie and I could decode it", *usercookie)

		// TODO(rjk): Handle revocation here, re-sign-in etc.
		return req, false

        } else {
		log.Println("anonymous access")
		usercookie.capability = CapabilityAnonymous
	}

	// Do I actually need to put this on the request? Only if the user state is shipped to another
	// service? i.e. This code is unnecessary? I should return a usercookie object?
	// the real reason to stick this data in the context is to support a not-in-process database layer?
	return req.WithContext(context.WithValue(req.Context(), COOKIE, usercookie)), true
}

// TODO(rjk): write me. Select users who have a cookie vended with capabilities but 
// have had their rights changed. Set this list of uuids. Force the complicated case. And handle
// TODO(rjk): this needs lots of tests (that aren't flaky?)
// TODO(rjk): this is on GlobalState.
//func (GlobalState*) InitializeRevokeUuidList() {
//}


// Might want versions on Context?
func  RequestIsAnonymous(req *http.Request) bool {
	uc := (req.Context().Value(COOKIE)).(*UserCookie)
	return (uc.capability & CapabilityAnonymous) != 0
}

func RequestHasCapability(req *http.Request, cap Capability) bool {
	uc := (req.Context().Value(COOKIE)).(*UserCookie)
	return (uc.capability & cap)!= 0
}

//func RequestUuid(req *http.Request) uuid.UUID {
//}
