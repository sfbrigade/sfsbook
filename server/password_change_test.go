package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/pborman/uuid"
	"github.com/rjkroege/mocking"
	"golang.org/x/crypto/bcrypt"
)

func makeUnderTestHandler(tape *mocking.Tape) *passwordChange {
	undertesthandler := &passwordChange{
		// Always use the embedded resource.
		embr:         makeEmbeddableResource(""),
		passwordfile: (*mockPasswordIndex)(tape),
	}
	return undertesthandler
}

func TestUsermgtNotsignedIn(t *testing.T) {
	defer resourceHelper()()

	undertesthandler := makeUnderTestHandler(nil)

	testreq := httptest.NewRequest("GET", "https://sfsbook.org/usermgt/changepasswd.html", nil)
	recorder := httptest.NewRecorder()
	testreq = testreq.WithContext(context.WithValue(testreq.Context(), UserCookieStateName, new(UserCookie)))

	undertesthandler.ServeHTTP(recorder, testreq)

	result := recorder.Result()
	if got, want := result.StatusCode, 200; got != want {
		t.Errorf("bad response code: got %v, want %v", got, want)
	}
	resultAsString, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal("couldn't read recorded response", err)
	}

	if got, want := string(resultAsString), "\n\tIsAuthed: false\n\tDisplayName: \n\tChangeAttemptedAndFailed: false\n\tChangeAttemptedAndSucceeded: false\n\tReasonForFailure: \n"; got != want {
		t.Errorf("bad response body: got %v\n(%#v)\nwant %v\n(%#v)", got, got, want, want)
	}
}

func TestUsermgtSignedIn(t *testing.T) {
	uuid := uuid.NewRandom()
	defer resourceHelper()()
	undertesthandler := makeUnderTestHandler(nil)

	testreq := httptest.NewRequest("POST", "https://sfsbook.org/usermgt/changepasswd.html", nil)
	recorder := httptest.NewRecorder()

	usercookie := &UserCookie{
		Uuid:        uuid,
		Capability:  CapabilityViewPublicResourceEntry,
		Displayname: "Homer Simpson",
		// Time not needed.
	}
	testreq = testreq.WithContext(context.WithValue(testreq.Context(), UserCookieStateName, usercookie))

	// Run handler.
	undertesthandler.ServeHTTP(recorder, testreq)

	result := recorder.Result()
	if got, want := result.StatusCode, 200; got != want {
		t.Errorf("bad response code: got %v, want %v", got, want)
	}
	resultAsString, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal("couldn't read recorded response", err)
	}

	if got, want := string(resultAsString), "\n\tIsAuthed: true\n\tDisplayName: Homer Simpson\n\tChangeAttemptedAndFailed: true\n\tChangeAttemptedAndSucceeded: false\n\tReasonForFailure: Need to enter the previous password\n"; got != want {
		t.Errorf("bad response body: got %v\n(%#v)\nwant %v\n(%#v)", got, got, want, want)
	}

	// With posted args: empty new password fields.
	postargs := strings.NewReader("oldpassword=o&newpassword=&newpasswordagain=")
	testreq = httptest.NewRequest("POST", "https://sfsbook.org/usermgt/changepasswd.html", postargs)
	testreq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	recorder = httptest.NewRecorder()
	testreq = testreq.WithContext(context.WithValue(testreq.Context(), UserCookieStateName, usercookie))

	// Run handler.
	undertesthandler.ServeHTTP(recorder, testreq)

	result = recorder.Result()
	if got, want := result.StatusCode, 200; got != want {
		t.Errorf("bad response code: got %v, want %v", got, want)
	}
	resultAsString, err = ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal("couldn't read recorded response", err)
	}

	if got, want := string(resultAsString), "\n\tIsAuthed: true\n\tDisplayName: Homer Simpson\n\tChangeAttemptedAndFailed: true\n\tChangeAttemptedAndSucceeded: false\n\tReasonForFailure: Need to enter a new password\n"; got != want {
		t.Errorf("bad response body: got %v\n(%#v)\nwant %v\n(%#v)", got, got, want, want)
	}

	// With posted args: differering new password fields.
	postargs = strings.NewReader("oldpassword=o&newpassword=p&newpasswordagain=na")
	testreq = httptest.NewRequest("POST", "https://sfsbook.org/usermgt/changepasswd.html", postargs)
	testreq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	recorder = httptest.NewRecorder()
	testreq = testreq.WithContext(context.WithValue(testreq.Context(), UserCookieStateName, usercookie))

	// Run handler.
	undertesthandler.ServeHTTP(recorder, testreq)

	result = recorder.Result()
	if got, want := result.StatusCode, 200; got != want {
		t.Errorf("bad response code: got %v, want %v", got, want)
	}
	resultAsString, err = ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal("couldn't read recorded response", err)
	}

	if got, want := string(resultAsString), "\n\tIsAuthed: true\n\tDisplayName: Homer Simpson\n\tChangeAttemptedAndFailed: true\n\tChangeAttemptedAndSucceeded: false\n\tReasonForFailure: New password fields need to match.\n"; got != want {
		t.Errorf("bad response body: got %v\n(%#v)\nwant %v\n(%#v)", got, got, want, want)
	}

	// With posted args: newpassword is too short.
	postargs = strings.NewReader("oldpassword=o&newpassword=pw0&newpasswordagain=pw0")
	testreq = httptest.NewRequest("POST", "https://sfsbook.org/usermgt/changepasswd.html", postargs)
	testreq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	recorder = httptest.NewRecorder()
	testreq = testreq.WithContext(context.WithValue(testreq.Context(), UserCookieStateName, usercookie))

	// Run handler.
	undertesthandler.ServeHTTP(recorder, testreq)

	result = recorder.Result()
	if got, want := result.StatusCode, 200; got != want {
		t.Errorf("bad response code: got %v, want %v", got, want)
	}
	resultAsString, err = ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal("couldn't read recorded response", err)
	}

	if got, want := string(resultAsString), "\n\tIsAuthed: true\n\tDisplayName: Homer Simpson\n\tChangeAttemptedAndFailed: true\n\tChangeAttemptedAndSucceeded: false\n\tReasonForFailure: New password is too easily guessed.\n"; got != want {
		t.Errorf("bad response body: got %v\n(%#v)\nwant %v\n(%#v)", got, got, want, want)
	}
}

func testStateSetup(uuid uuid.UUID, poststring string) (*httptest.ResponseRecorder, *http.Request) {
	postargs := strings.NewReader(poststring)
	testreq := httptest.NewRequest("POST", "https://sfsbook.org/usermgt/changepasswd.html", postargs)
	testreq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()

	usercookie := &UserCookie{
		Uuid:        uuid,
		Capability:  CapabilityViewPublicResourceEntry,
		Displayname: "Homer Simpson",
		// Time not needed.
	}
	testreq = testreq.WithContext(context.WithValue(testreq.Context(), UserCookieStateName, usercookie))
	return recorder, testreq
}

func TestUsermgtAffectingDatabase(t *testing.T) {
	uuid := uuid.NewRandom()
	defer resourceHelper()()
	tape := mocking.NewTape()
	undertesthandler := makeUnderTestHandler(tape)

	// uuid is missing test.
	tape.SetResponses(
		fmt.Errorf("internal database error"),
	)
	recorder, testreq := testStateSetup(uuid, "oldpassword=op&newpassword=pinky0&newpasswordagain=pinky0")

	// Run handler.
	undertesthandler.ServeHTTP(recorder, testreq)

	result := recorder.Result()
	if got, want := result.StatusCode, 200; got != want {
		t.Errorf("bad response code: got %v, want %v", got, want)
	}
	resultAsString, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal("couldn't read recorded response", err)
	}

	if got, want := string(resultAsString), "\n\tIsAuthed: true\n\tDisplayName: Homer Simpson\n\tChangeAttemptedAndFailed: true\n\tChangeAttemptedAndSucceeded: false\n\tReasonForFailure: Account error. Please sign-out and sign-in again.\n"; got != want {
		t.Errorf("bad response body: got %v\n(%#v)\nwant %v\n(%#v)", got, got, want, want)
	}

	if got, expected := tape.Play(), []interface{}{docStim{"PasswordIndex.Document", string(uuid)}}; !reflect.DeepEqual(expected, got) {
		t.Errorf("invalid call sequence. Got %#v, want %#v", got, expected)
	}

	// password mismatch test
	tape.Rewind()
	tape.SetResponses(
		map[string]interface{}{
			"passwordhash": "badpassword",
		},
	)
	recorder, testreq = testStateSetup(uuid, "oldpassword=op&newpassword=pinky0&newpasswordagain=pinky0")

	undertesthandler.ServeHTTP(recorder, testreq)

	result = recorder.Result()
	if got, want := result.StatusCode, 200; got != want {
		t.Errorf("bad response code: got %v, want %v", got, want)
	}
	resultAsString, err = ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal("couldn't read recorded response", err)
	}

	if got, want := string(resultAsString), "\n\tIsAuthed: true\n\tDisplayName: Homer Simpson\n\tChangeAttemptedAndFailed: true\n\tChangeAttemptedAndSucceeded: false\n\tReasonForFailure: Old password is incorrect\n"; got != want {
		t.Errorf("bad response body: got %v\n(%#v)\nwant %v\n", got, got, want)
	}

	if got, expected := tape.Play(), []interface{}{docStim{"PasswordIndex.Document", string(uuid)}}; !reflect.DeepEqual(expected, got) {
		t.Errorf("invalid call sequence. Got %#v, want %#v", got, expected)
	}

	// password update test with error
	encodedold, err := bcrypt.GenerateFromPassword([]byte("op"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("can't encode password", err)
	}
	tape.Rewind()
	tape.SetResponses(
		map[string]interface{}{
			"passwordhash": string(encodedold),
		},
		fmt.Errorf("database failed!"),
	)
	recorder, testreq = testStateSetup(uuid, "oldpassword=op&newpassword=pinky0&newpasswordagain=pinky0")

	undertesthandler.ServeHTTP(recorder, testreq)

	result = recorder.Result()
	if got, want := result.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("bad response code: got %v, want %v", got, want)
	}

	successfulTape := []interface{}{
		docStim{"PasswordIndex.Document", string(uuid)},
		indexStim{
			"PasswordIndex.Index",
			string(uuid),
			map[string]interface{}{
				"passwordhash": "foo",
			},
		},
	}

	playedtape := tape.Play()
	if got, expected := playedtape[0], successfulTape[0]; !reflect.DeepEqual(expected, got) {
		t.Errorf("invalid call sequence. Got %#v, want %#v", got, expected)
	}
	is, ok := playedtape[1].(indexStim)
	if !ok {
		t.Errorf("Expected an indexStim in the tape but got an %#v instead", playedtape[1])
	}
	ormmap, ok := is.data.(map[string]interface{})
	if !ok {
		t.Errorf("indexedStim data  is of wrong type. Is %#v instead", is.data)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(ormmap["passwordhash"].(string)), []byte("pinky0")); err != nil {
		t.Error("did not write a valid encrypted password into the tape.")
	}

	// successful update
	tape.Rewind()
	tape.SetResponses(
		map[string]interface{}{
			"passwordhash": string(encodedold),
			"extrafield":   "we copy extrafields",
		},
		nil,
	)
	recorder, testreq = testStateSetup(uuid, "oldpassword=op&newpassword=pinky0&newpasswordagain=pinky0")

	// Run handler.
	undertesthandler.ServeHTTP(recorder, testreq)

	result = recorder.Result()
	if got, want := result.StatusCode, 200; got != want {
		t.Errorf("bad response code: got %v, want %v", got, want)
	}
	resultAsString, err = ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal("couldn't read recorded response", err)
	}

	if got, want := string(resultAsString), "\n\tIsAuthed: true\n\tDisplayName: Homer Simpson\n\tChangeAttemptedAndFailed: false\n\tChangeAttemptedAndSucceeded: true\n\tReasonForFailure: \n"; got != want {
		t.Errorf("bad response body: got %v\n(%#v)\nwant %v\n(%#v)", got, got, want, want)
	}

	successfulTape = []interface{}{
		docStim{"PasswordIndex.Document", string(uuid)},
		indexStim{
			"PasswordIndex.Index",
			string(uuid),
			map[string]interface{}{
				"passwordhash": "foo",
				"extrafield":   "we copy extrafields",
			},
		},
	}
	playedtape = tape.Play()
	if got, expected := playedtape[0], successfulTape[0]; !reflect.DeepEqual(expected, got) {
		t.Errorf("invalid call sequence. Got %#v, want %#v", got, expected)
	}
	is, ok = playedtape[1].(indexStim)
	if !ok {
		t.Errorf("Expected an indexStim in the tape but got an %#v instead", playedtape[1])
	}
	ormmap, ok = is.data.(map[string]interface{})
	if !ok {
		t.Errorf("indexedStim data  is of wrong type. Is %#v instead", is.data)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(ormmap["passwordhash"].(string)), []byte("pinky0")); err != nil {
		t.Error("did not write a valid encrypted password into the tape.")
	}
	if got, expected := ormmap["extrafield"].(string), "we copy extrafields"; !reflect.DeepEqual(expected, got) {
		t.Errorf("invalid call sequence. Got %#v, want %#v", got, expected)
	}

}
