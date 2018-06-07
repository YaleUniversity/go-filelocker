package filelocker

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testUser = "testuser"
	testKey  = "secretsecret"
)

func TestNewClient(t *testing.T) {
	fl := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("Set-Cookie", "filelocker=123sessiontoken321")
		t.Logf("request url: %s", r.URL)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("error reading body", err)
		}
		t.Logf("request body: %s", string(body))

		loginResp := `
		<?xml version="1.0"?>
		<cli_response>
			<messages><info>123requestorigin321</info></messages>
			<data></data>
		</cli_response>
		`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(loginResp))
	}))
	defer fl.Close()

	filelockerURL := fl.URL
	client, err := NewClient(testUser, testKey, filelockerURL, nil)
	if err != nil {
		t.Error("error creating new filelocker client:", err)
	}

	if client.Origin != "123requestorigin321" {
		t.Errorf("Expected origin: 123requestorigin321, got %s", client.Origin)
	}

	t.Log(client)
}
