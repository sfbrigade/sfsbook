package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/rjkroege/mocking"
)

// TestEditUsersActions shows that a user with capability can edit
// values on a user structure.
func TestEditUsersActions(t *testing.T) {
	defer resourceHelper()()
	tape := mocking.NewTape()
	undertesthandler := makeUnderTestHandlerListUsers(tape)

	testPatterns := []testPattern{
		// Basic editing: change the role from volunteer to admin.
		{
			"?userquery=&selected-0=31E946C1-7F1A-491D-BAAE-6BAEA3641FC8&rolechange=admin",
			http.StatusOK,
			[]interface{}{
				map[string]interface{} {
					"display_name": "Homer Simpson",
					"role":  "volunteer",
					"name": "homer.simpson",
				},
				nil,
				[]map[string]interface{}{
					{
						"display_name": "Homer Simpson",
					},
					{
						"display_name": "Lisa Simpson",
					},
				},
			},
			[]interface {}{
				docStim{name:"PasswordIndex.Document", uuid:"1\xe9F\xc1\u007f\x1aI\x1d\xba\xaek\xae\xa3d\x1f\xc8"}, 
				indexStim{fun:"PasswordIndex.Index", id:"1\xe9F\xc1\u007f\x1aI\x1d\xba\xaek\xae\xa3d\x1f\xc8", data:map[string]interface {}{"name":"homer.simpson", "role":"admin", "display_name":"Homer Simpson"}},
				listUsersStim{name: "PasswordIndex.ListUsers", query:"", size:10, from:0},
			},
			"\n\tIsAuthed: true\n\tDisplayName: Homer Simpson\n\n\tUserquery: \n\tUsers: [map[display_name:Homer Simpson] map[display_name:Lisa Simpson]]\n\tQuerysuccess: true\n\tDiagnosticmessage: Showing all...\n",
		},
		// Test that the deletion works.
		{
			"?userquery=&selected-0=31E946C1-7F1A-491D-BAAE-6BAEA3641FC8&rolechange=nochange&deleteuser=Delete",
			http.StatusOK,
			[]interface{}{
				nil,
				[]map[string]interface{}{
					{
						"display_name": "Lisa Simpson",
					},
				},
			},
			[]interface {}{
				deleteStim{name:"PasswordIndex.Delete", uuid:"1\xe9F\xc1\u007f\x1aI\x1d\xba\xaek\xae\xa3d\x1f\xc8"}, 
				listUsersStim{name: "PasswordIndex.ListUsers", query:"", size:10, from:0},
			},
			"\n\tIsAuthed: true\n\tDisplayName: Homer Simpson\n\n\tUserquery: \n\tUsers: [map[display_name:Lisa Simpson]]\n\tQuerysuccess: true\n\tDiagnosticmessage: Showing all...\n",
		},
		// Test that an invalid role is ignored.
		{
			"?userquery=&selected-0=31E946C1-7F1A-491D-BAAE-6BAEA3641FC8&rolechange=fuzzypeaches",
			http.StatusBadRequest,
			[]interface{}{},
			[]interface {}{},
			"client is attempting something wrong",
		},
		// Test that the deletion failure is correctly handled.
		{
			"?userquery=&selected-0=31E946C1-7F1A-491D-BAAE-6BAEA3641FC8&rolechange=nochange&deleteuser=Delete",
			http.StatusOK,
			[]interface{}{
				fmt.Errorf("PasswordIndex.Delete failed"),
				[]map[string]interface{}{
					{
						"display_name": "Lisa Simpson",
					},
				},
			},
			[]interface {}{
				deleteStim{name:"PasswordIndex.Delete", uuid:"1\xe9F\xc1\u007f\x1aI\x1d\xba\xaek\xae\xa3d\x1f\xc8"}, 
				listUsersStim{name: "PasswordIndex.ListUsers", query:"", size:10, from:0},
			},
			"\n\tIsAuthed: true\n\tDisplayName: Homer Simpson\n\n\tUserquery: \n\tUsers: [map[display_name:Lisa Simpson]]\n\tQuerysuccess: true\n\tDiagnosticmessage: Couldn&#39;t successfully delete all of the selected users.\n",
		},
	}

	for _, tp  := range testPatterns {
		testreq := httptest.NewRequest("GET",
			"https://sfsbook.org/usermgt/listusers.html" + tp.urlargs, nil)
		recorder := httptest.NewRecorder()
		testreq = addCookie(testreq, CapabilityViewUsers | CapabilityEditUsers)

		// Note simplified user data to avoid the issue that the maps are not
		// emitted in a consistent order.
		tape.Rewind()
		tape.SetResponses(tp.tapeResponse...)

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

