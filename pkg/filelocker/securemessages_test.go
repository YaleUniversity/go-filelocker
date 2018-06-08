package filelocker_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/YaleUniversity/go-filelocker/pkg/filelocker"
)

func TestSecureMessages(t *testing.T) {
	messagesList := `
	{
		"sMessages": [],
		"fMessages": [],
		"data": [
			[{
				"id": 1,
				"ownerId": "bossman",
				"body": "some secret message",
				"subject": "shh",
				"messageRecipients": ["peon1", "peon2"],
				"expirationDatetime": "07/07/2018",
				"viewedDatetime": null,
				"creationDatetime": "06/07/2018"
			}], [{
				"id": 2,
				"ownerId": "bossman",
				"body": "you are fired",
				"subject": "ohai",
				"messageRecipients": ["peon1"],
				"expirationDatetime": "07/07/2018",
				"viewedDatetime": null,
				"creationDatetime": "06/07/2018"
			}]
		]
	}
	`

	expected := &filelocker.SecureMessagesResponse{
		ErrorMessages: []string{},
		InfoMessages:  []string{},
		Messages: [][]filelocker.SecureMessage{
			[]filelocker.SecureMessage{
				filelocker.SecureMessage{
					Body:       "some secret message",
					Created:    "06/07/2018",
					Expiration: "07/07/2018",
					ID:         1,
					OwnerID:    "bossman",
					Recipients: []string{
						"peon1",
						"peon2",
					},
					Subject: "shh",
					Viewed:  "",
				},
			},
			[]filelocker.SecureMessage{
				filelocker.SecureMessage{
					Body:       "you are fired",
					Created:    "06/07/2018",
					Expiration: "07/07/2018",
					ID:         2,
					OwnerID:    "bossman",
					Recipients: []string{
						"peon1",
					},
					Subject: "ohai",
					Viewed:  "",
				},
			},
		},
	}

	fl := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.String() != "/message/get_messages" {
			t.Errorf("got url %s, expected '/message/get_messages'", r.URL)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(messagesList))
	}))
	defer fl.Close()

	bURL, err := url.Parse(fl.URL)
	if err != nil {
		t.Error(err)
	}

	client := filelocker.Client{
		Client:  http.DefaultClient,
		Origin:  "123requestorigin321",
		BaseURL: bURL,
	}

	actual, err := client.SecureMessages()
	if err != nil {
		t.Error("error listing secure messages", err)
	}

	if !reflect.DeepEqual(*expected, *actual) {
		t.Errorf("expected: %+v\ngot: %+v", *expected, *actual)
	}
}

func TestSecureMessagesError(t *testing.T) {
	messages := `
	{
		"sMessages": [],
		"fMessages": ["stuff broke"]
	}
	`

	fl := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.String() != "/message/get_messages" {
			t.Errorf("got url %s, expected '/message/get_messages'", r.URL)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(messages))
	}))
	defer fl.Close()

	bURL, err := url.Parse(fl.URL)
	if err != nil {
		t.Error(err)
	}

	client := filelocker.Client{
		Client:  http.DefaultClient,
		Origin:  "123requestorigin321",
		BaseURL: bURL,
	}

	_, err = client.SecureMessages()
	if err == nil {
		t.Error("expected secure messages list error, got nil")
	}
}

func TestSecureMessagesCount(t *testing.T) {
	messagesCount := `
	{
		"sMessages": [],
		"fMessages": [],
		"data": 1337
	}
	`

	expected := &filelocker.SecureMessageCountResponse{
		ErrorMessages: []string{},
		InfoMessages:  []string{},
		Count:         1337,
	}

	fl := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.String() != "/message/get_new_message_count" {
			t.Errorf("got url %s, expected '/message/get_new_message_count'", r.URL)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(messagesCount))
	}))
	defer fl.Close()

	bURL, err := url.Parse(fl.URL)
	if err != nil {
		t.Error(err)
	}

	client := filelocker.Client{
		Client:  http.DefaultClient,
		Origin:  "123requestorigin321",
		BaseURL: bURL,
	}

	actual, err := client.SecureMessagesCount()
	if err != nil {
		t.Error("error counting secure messages", err)
	}

	if !reflect.DeepEqual(*expected, *actual) {
		t.Errorf("expected: %+v\ngot: %+v", *expected, *actual)
	}
}

func TestSecureMessagesCountError(t *testing.T) {
	messages := `
	{
		"sMessages": [],
		"fMessages": ["stuff broke"]
	}
	`

	fl := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.String() != "/message/get_new_message_count" {
			t.Errorf("got url %s, expected '/message/get_new_message_count'", r.URL)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(messages))
	}))
	defer fl.Close()

	bURL, err := url.Parse(fl.URL)
	if err != nil {
		t.Error(err)
	}

	client := filelocker.Client{
		Client:  http.DefaultClient,
		Origin:  "123requestorigin321",
		BaseURL: bURL,
	}

	_, err = client.SecureMessagesCount()
	if err == nil {
		t.Error("expected secure messages count error, got nil")
	}
}

func TestSecureMessageDelete(t *testing.T) {
	messageDelete := `
	{
		"sMessages": ["Message(s) deleted"],
		"fMessages": []
	}
	`

	expected := &filelocker.SecureMessagesDeleteResponse{
		ErrorMessages: []string{},
		InfoMessages:  []string{"Message(s) deleted"},
	}

	fl := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.String() != "/message/delete_messages" {
			t.Errorf("got url %s, expected '/message/delete_messages'", r.URL)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("error reading body", err)
		}

		values, err := url.ParseQuery(string(body))
		if err != nil {
			t.Error(err)
		}

		if v, ok := values["messageIds"]; !ok {
			t.Error("expected messageIds parameter!")
		} else if v[0] != "1,2,3" {
			t.Errorf("expected messageIds parameter to be '1,2,3', got %s", v[0])
		}

		if v, ok := values["requestOrigin"]; !ok {
			t.Error("expected requestOrigin parameter!")
		} else if v[0] != "123requestorigin321" {
			t.Errorf("expected requestOrigin parameter to be '123requestorigin321', got %s", v[0])
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(messageDelete))
	}))
	defer fl.Close()

	bURL, err := url.Parse(fl.URL)
	if err != nil {
		t.Error(err)
	}

	client := filelocker.Client{
		Client:  http.DefaultClient,
		Origin:  "123requestorigin321",
		BaseURL: bURL,
	}

	actual, err := client.SecureMessagesDelete([]int{1, 2, 3})
	if err != nil {
		t.Error("error deleting secure messages", err)
	}

	if !reflect.DeepEqual(*expected, *actual) {
		t.Errorf("expected: %+v\ngot: %+v", *expected, *actual)
	}
}

func TestSecureMessagesDeleteError(t *testing.T) {
	messages := `
	{
		"sMessages": [],
		"fMessages": ["stuff broke"],
	}
	`

	fl := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.String() != "/message/delete_messages" {
			t.Errorf("got url %s, expected '/message/delete_messages'", r.URL)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(messages))
	}))
	defer fl.Close()

	bURL, err := url.Parse(fl.URL)
	if err != nil {
		t.Error(err)
	}

	client := filelocker.Client{
		Client:  http.DefaultClient,
		Origin:  "123requestorigin321",
		BaseURL: bURL,
	}

	_, err = client.SecureMessagesDelete([]int{1, 2, 3})
	if err == nil {
		t.Error("expected secure messages delete error, got nil")
	}
}

func TestNewSecureMessage(t *testing.T) {
	messageSend := `
	{
		"sMessages": ["Message \"test test\" sent."],
		"fMessages": []
	}
	`

	expected := &filelocker.NewMessageResponse{
		ErrorMessages: []string{},
		InfoMessages:  []string{"Message \"test test\" sent."},
	}

	expireTime := time.Now().Add(108000 * time.Hour)
	fl := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.String() != "/message/create_message" {
			t.Errorf("got url %s, expected '/message/create_message'", r.URL)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("error reading body", err)
		}

		values, err := url.ParseQuery(string(body))
		if err != nil {
			t.Error(err)
		}

		if v, ok := values["requestOrigin"]; !ok {
			t.Error("expected requestOrigin parameter!")
		} else if v[0] != "123requestorigin321" {
			t.Errorf("expected requestOrigin parameter to be '123requestorigin321', got %s", v[0])
		}

		if v, ok := values["subject"]; !ok {
			t.Error("expected subject parameter!")
		} else if v[0] != "test test" {
			t.Errorf("expected subject parameter to be 'test test', got %s", v[0])
		}

		if v, ok := values["body"]; !ok {
			t.Error("expected body parameter!")
		} else if v[0] != "secret" {
			t.Errorf("expected body parameter to be 'secret', got %s", v[0])
		}

		if v, ok := values["expiration"]; !ok {
			t.Error("expected expiration parameter!")
		} else if v[0] != expireTime.Format("01/02/2006") {
			t.Errorf("expected expiration parameter to be '%s', got %s", expireTime.Format("01/02/2006"), v[0])
		}

		if v, ok := values["recipientIds"]; !ok {
			t.Error("expected recipientIds parameter!")
		} else if v[0] != "user1,user2" {
			t.Errorf("expected recipientIds parameter to be 'user1,user2', got %s", v[0])
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(messageSend))
	}))
	defer fl.Close()

	bURL, err := url.Parse(fl.URL)
	if err != nil {
		t.Error(err)
	}

	client := filelocker.Client{
		Client:  http.DefaultClient,
		Origin:  "123requestorigin321",
		BaseURL: bURL,
	}

	actual, err := client.NewSecureMessage("test test", "secret", []string{"user1", "user2"}, expireTime)
	if err != nil {
		t.Error("error sending secure message", err)
	}

	if !reflect.DeepEqual(*expected, *actual) {
		t.Errorf("expected: %+v\ngot: %+v", *expected, *actual)
	}
}

func TestNewSecureMessagesError(t *testing.T) {
	messages := `
	{
		"sMessages": [],
		"fMessages": ["stuff broke"],
	}
	`

	fl := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.String() != "/message/create_message" {
			t.Errorf("got url %s, expected '/message/create_message'", r.URL)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(messages))
	}))
	defer fl.Close()

	bURL, err := url.Parse(fl.URL)
	if err != nil {
		t.Error(err)
	}

	client := filelocker.Client{
		Client:  http.DefaultClient,
		Origin:  "123requestorigin321",
		BaseURL: bURL,
	}

	_, err = client.NewSecureMessage("test test err", "secret", []string{"user1"}, time.Now())
	if err == nil {
		t.Error("expected secure messages delete error, got nil")
	}
}
