package server

import (
    "context"
    "log"
    "net/http"
    "net/http/httptest"
    "net/url"
    "strings"
    "testing"

    "github.com/gorilla/securecookie"
    "github.com/rjkroege/mocking"
    "github.com/pborman/uuid"
)


// TODO: generalize / extract & import:
// copied & modified from usermgt_test.go:makeUnderTestHandler
// (and changed passwordChange to loginServer)
func makeMockServer(tape *mocking.Tape) *loginServer {

    // not secure (not random) but OK for testing
    var hashKey = []byte("very-secret")
    var blockKey = []byte(nil)

    undertesthandler := &loginServer{
        // Always use the embedded resource.
        embr:         makeEmbeddableResource(""),
        passwordfile: (*mockPasswordIndex)(tape),
        cookiecodec: securecookie.New(hashKey, blockKey),
    }
    return undertesthandler
}


// attempt login and ensure that response code is redirect to index.html
func TestServeHTTP(t *testing.T) {

    log.SetFlags(log.LstdFlags | log.Lshortfile) // debug: verbose logs including filename

    uuid := uuid.NewRandom()
    defer testHelper()()
    undertesthandler := makeMockServer(nil)

    usercookie := &UserCookie{
        Uuid:        uuid,
        Capability:  CapabilityVolunteer,
        Displayname: "Homer Simpson",
        // Time not needed.
    }

    recorder := httptest.NewRecorder()

    // initialize our form values and set up the test
    values := url.Values{
        "username": {"test_user"},
        "password": {"password"},
    }
    testreq := httptest.NewRequest("POST", "/login.html", strings.NewReader(values.Encode()))
    testreq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    testreq = testreq.WithContext(context.WithValue(testreq.Context(), UserCookieStateName, usercookie))

    // Run handler.
    undertesthandler.ServeHTTP(recorder, testreq)

    // The import part of the test -- did we get redirect (http.StatusFound) to /index.html?
    result := recorder.Result()
    if got, want := result.StatusCode, http.StatusFound; got != want {
        t.Errorf("bad response code: got %v (%v), want %v (%v)", got, http.StatusText(got), want, http.StatusText(want))
    }
    if loc, err := result.Location(); err != nil || loc.RequestURI() != "/index.html" {
        if err != nil {
            t.Errorf("Expected Location header set; but it isn't: ", err)
        } else {
            t.Errorf("expected redirect to /index.html, got location: ", loc.RequestURI())
        }
    }
}
