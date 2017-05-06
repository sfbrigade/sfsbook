package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/rjkroege/mocking"
)

// TestListUsersShowBasicList shows that a user with capability can list
// the currently configured users.
func TestEditUsersActions(t *testing.T) {
	defer resourceHelper()()
	tape := mocking.NewTape()
	undertesthandler := makeUnderTestHandlerListUsers(tape)

	testPatterns := []testPattern{
		// Sets the roll to admin
		{
			"?userquery=&selected-0=31E946C1-7F1A-491D-BAAE-6BAEA3641FC8&rolechange=admin",
			http.StatusOK,
			[]map[string]interface{}{},
			[]interface {}{},
			"\n\tIsAuthed: true\n\tDisplayName: Homer Simpson\n\n\tUserquery: \n\tUsers: []\n\tQuerysuccess: false\n\tDiagnosticmessage: Sign in as an admin to edit users.\n",
		},
		// TODO(rjk): delete a user
	}

	for _, tp  := range testPatterns {
		testreq := httptest.NewRequest("GET",
			"https://sfsbook.org/usermgt/listusers.html" + tp.urlargs, nil)
		recorder := httptest.NewRecorder()
		testreq = addCookie(testreq, CapabilityViewUsers | CapabilityEditUsers)

		// Note simplified user data to avoid the issue that the maps are not
		// emitted in a consistent order.
		tape.Rewind()
		tape.SetResponses(tp.tapeResponse)

		// Run handler.
		undertesthandler.ServeHTTP(recorder, testreq)

		result := recorder.Result()
		if got, want := result.StatusCode, tp.statuscode; got != want {
			t.Errorf("bad response code: got %v, want %v", got, want)
		}
		resultAsString, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal("couldn't read recorded response", err)
		}

		if got, want := string(resultAsString), tp.outputString; got != want {
			t.Errorf("bad response body: got %v\n(%#v)\nwant %v\n(%#v)", got, got, want, want)
		}
		playedtape := tape.Play()
		if got, want := playedtape, tp.tapeRecord; !reflect.DeepEqual(want, got) {
			t.Errorf("invalid call sequence. Got %#v, want %#v", got, want)
		}
	}
}

