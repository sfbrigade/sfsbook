package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/pborman/uuid"
)

type delegateHandlerBasic testing.T

func (dh *delegateHandlerBasic) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t := (*testing.T)(dh)

	uc := GetCookie(req)

	if uc.IsAuthed() {
		t.Error("delegate should not be authed")
	}

	if got, want := uc.Capability, CapabilityAnonymous; got != want {
		t.Errorf("capability mismatch: got %v, want %v", got, want)
	}
}

func TestNoHazCookieHandler(t *testing.T) {
	statepath, err := ioutil.TempDir("", "sfsbook")
	if err != nil {
		t.Fatal("can't make a temporary directory", err)
	}
	defer os.RemoveAll(statepath)

	cookiecodec, err := makeCookieTooling(statepath)
	if err != nil {
		t.Error("can't make cookie keys")
		return
	}

	undertesthandler := &cookieHandler{
		cookiecodec: cookiecodec,
		delegate:    (*delegateHandlerBasic)(t),
	}

	testreq := httptest.NewRequest("GET", "https://sfsbook.org/index.html", nil)
	recorder := httptest.NewRecorder()

	undertesthandler.ServeHTTP(recorder, testreq)
	result := recorder.Result()
	if got, want := result.StatusCode, 200; got != want {
		t.Errorf("bad response code: got %v, want %v", got, want)
	}

}

type delegateHandlerAuthedCookie struct {
	t  *testing.T
	uc *UserCookie
}

func (dh *delegateHandlerAuthedCookie) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t := dh.t
	uc := GetCookie(req)
	wantuc := dh.uc

	if !uc.IsAuthed() {
		t.Error("delegate should be authed")
	}

	if got, want := uc.Capability, wantuc.Capability; got != want {
		t.Errorf("capability mismatch: got %v, want %v", got, want)
	}
	if got, want := uc.Uuid, wantuc.Uuid; !reflect.DeepEqual(got, want) {
		t.Errorf("Uuid mismatch: got %v, want %v", got, want)
	}
	if got, want := uc.Displayname, wantuc.Displayname; got != want {
		t.Errorf("Displayname mismatch: got %v, want %v", got, want)
	}
}

func TestHazCookieHandler(t *testing.T) {
	statepath, err := ioutil.TempDir("", "sfsbook")
	if err != nil {
		t.Fatal("can't make a temporary directory", err)
	}
	defer os.RemoveAll(statepath)

	cookiecodec, err := makeCookieTooling(statepath)
	if err != nil {
		t.Error("can't make cookie keys")
		return
	}

	usercookie := &UserCookie{
		Uuid:        uuid.NewRandom(),
		Capability:  CapabilityViewOtherVolunteerComment | CapabilityInviteNewAdmin,
		Displayname: "Spiffy Tester",
	}

	undertesthandler := &cookieHandler{
		cookiecodec: cookiecodec,
		delegate: &delegateHandlerAuthedCookie{
			t:  t,
			uc: usercookie,
		},
	}

	testreq := httptest.NewRequest("GET", "https://sfsbook.org/index.html", nil)

	encodedcookie, err := cookiecodec.Encode(SessionCookieName, usercookie)
	if err != nil {
		t.Log(cookiecodec)
		t.Errorf("failed to encode cookie %v", err)
		return
	}

	// Add encoded cookie to the testreq.
	testreq.AddCookie(&http.Cookie{
		Name:  SessionCookieName,
		Value: encodedcookie,
		Path:  "/",
	})

	recorder := httptest.NewRecorder()
	undertesthandler.ServeHTTP(recorder, testreq)
	result := recorder.Result()
	if got, want := result.StatusCode, 200; got != want {
		t.Errorf("bad response code: got %v, want %v", got, want)
	}
}

func TestMakeCookieCodec(t *testing.T) {
	statepath, err := ioutil.TempDir("", "sfsbook")
	t.Log(statepath)
	if err != nil {
		t.Fatal("can't make a temporary directory", err)
	}
	defer os.RemoveAll(statepath)

	cookiecodec, err := makeCookieTooling(statepath)
	if err != nil {
		t.Fatal("can't make cookie", err)
	}

	encodedcookie, err := cookiecodec.Encode(SessionCookieName, "should be encrypted")
	if err != nil {
		t.Log(cookiecodec)
		t.Errorf("failed to encode cookie %v", err)
	}

	if encodedcookie == "should be encrypted" {
		t.Error("not encrypted?")
	}

	cp := filepath.Join(statepath, "hashkey.dat")
	if _, err := os.Stat(cp); err != nil {
		t.Fatalf("cookie not written to %s: %v", cp, err)
	}
	cp = filepath.Join(statepath, "blockkey.dat")
	if _, err := os.Stat(cp); err != nil {
		t.Fatalf("cookie not written to %s: %v", cp, err)
	}

	newcookiecodec, err := makeCookieTooling(statepath)
	if err != nil {
		t.Fatal("can't make cookie", err)
	}
	newencodedcookie, err := newcookiecodec.Encode(SessionCookieName, "should also be encrypted")
	if err != nil {
		t.Log(newcookiecodec)
		t.Errorf("failed to encode cookie %v", err)
	}

	if newencodedcookie == "should also be encrypted" {
		t.Error("not encrypted?")
	}
}
