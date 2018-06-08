package filelocker_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/YaleUniversity/go-filelocker/pkg/filelocker"
)

var (
	testUser  = "testuser"
	testKey   = "secretsecret"
	loginResp = `
	<?xml version="1.0"?>
	<cli_response>
		<messages><info>123requestorigin321</info></messages>
		<data></data>
	</cli_response>
	`
	loginCookie = "filelocker=123sessiontoken321"
)

func TestNewClient(t *testing.T) {
	fl := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("Set-Cookie", loginCookie)

		if r.URL.String() != "/cli/CLI_login" {
			t.Errorf("got url %s, expected '/cli/CLI_login'", r.URL)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("error reading body", err)
		}

		values, err := url.ParseQuery(string(body))
		if err != nil {
			t.Error(err)
		}

		if v, ok := values["CLIkey"]; !ok {
			t.Error("expected CLIkey parameter!")
		} else if v[0] != "secretsecret" {
			t.Errorf("expected CLIkey parameter to be 'secretsecret', got %s", v[0])
		}

		if v, ok := values["userId"]; !ok {
			t.Error("expected userId parameter!")
		} else if v[0] != "testuser" {
			t.Errorf("expected userId parameter to be 'testuser', got %s", v[0])
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(loginResp))
	}))
	defer fl.Close()

	filelockerURL := fl.URL
	client, err := filelocker.NewClient(testUser, testKey, filelockerURL, nil)
	if err != nil {
		t.Error("error creating new filelocker client:", err)
	}

	if client.Origin != "123requestorigin321" {
		t.Errorf("Expected origin: 123requestorigin321, got %s", client.Origin)
	}

	t.Log(client)
}
