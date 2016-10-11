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

	"github.com/gorilla/securecookie"
	"github.com/pborman/uuid"
)

// Capability is a set of capabilities. Sized because we may need
// to transmit it over an IPC boundary.
type Capability int64

const (
	// Finding this Capability in a UserCookie implies that the uuid is empty.
	CapabilityAnonymous Capability = 0
)

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

// UserCookie is data stored in the auth cookie. It's encrypted via
// the securecookie facilities and set as a cookie on the interaction with
// the remote UA. It holds the identity and capability of the UA.
// I presume that the cookie encrypting is sufficiently secure to
// defeat an attempt to decrypt, alter the capability and reencode.
type UserCookie struct {
	// The user identifier.
	uuid uuid.UUID

	// A mask of capabilities.
	capability Capability

	// The time that the cookie was created.
	timestamp time.Time

	// The user's display_name
	display_name string
}

// TODO(rjk): Add the ability to check that a given uuid needs to be
// revalidated.

// cookieHandler is the state for an implementation of http.Handler that
// can invoke its delegatehandler with a decoded auth cookie context.
type cookieHandler struct {
	securecookie.SecureCookie

	// TODO(rjk): Implement revocation.
	revokelist []uuid.UUID
	delegate   http.Handler
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

// cookieTooling contains the cookie-related tooling for the HandlerFactory.
type cookieTooling struct {
	hashkey, blockkey []byte
}

// makeCookieTooling constructs cookie tooling for the HandlerFactory.
func makeCookieTooling(statepath string) (*cookieTooling, error) {
	hashkey, err := makeCookie(statepath, "hashkey.dat")
	if err != nil {
		return nil, err
	}
	blockkey, err := makeCookie(statepath, "blockkey.dat")
	if err != nil {
		return nil, err
	}
	return &cookieTooling{
		hashkey:  hashkey,
		blockkey: blockkey,
	}, nil
}

// MakeUserStateHandler builds a http.Handler that can
// decrypt auth cookies. See ServeHTTP below.
func (hf *HandlerFactory) makeCookieHandler(delegate http.Handler) http.Handler {
	return &cookieHandler{
		SecureCookie: *securecookie.New(hf.cookietool.hashkey, hf.cookietool.blockkey),
		revokelist:   make([]uuid.UUID, 0, 10),
		delegate:     delegate,
	}
}

const SessionCookieName = "session"
const UserCookieStateName = "usercookiestate"

// GetCookie retrieves the UserState from the context.
func GetCookie(req *http.Request) *UserCookie {
	// It's a fatal error to put anything other than a *UserCookie in slot
	// UserCookieStateName.
	return req.Context().Value(UserCookieStateName).(*UserCookie)
}

// ServeHTTP updates the given http.Request with a decoded
// instance of the auth cookie or updates the response to redirect
// appropriately.
func (cf *cookieHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(SessionCookieName)
	usercookie := new(UserCookie)
	if err == nil {
		// This request has a cookie.
		if err = cf.Decode(SessionCookieName, cookie.Value, usercookie); err != nil {
			log.Println("request had a cookie but it was not decodeable:", err)
			// TODO(rjk):
			// redirect to the login page with an appropriate error message.
		}
		log.Println("request had a cookie and I could decode it", *usercookie)
		// TODO(rjk): Test here for revocation, cookie rotation, etc.
	} else {
		log.Println("anonymous access")
		usercookie.capability = CapabilityAnonymous
	}

	cf.delegate.ServeHTTP(w, req.WithContext(
		context.WithValue(
			req.Context(), UserCookieStateName, usercookie)))
}

// TODO(rjk): Need a mechanism for revoking credentials.

func (u *UserCookie) IsAuthed() bool {
	return u.capability != CapabilityAnonymous
}

func (u *UserCookie) DisplayName() string {
	return u.display_name
}
