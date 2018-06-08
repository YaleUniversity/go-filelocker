package filelocker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SecureMessage is an encrypted message sent through filelocker
type SecureMessage struct {
	Body       string   `json:"body"`
	Created    string   `json:"creationDatetime"`
	Expiration string   `json:"expirationDatetime"`
	ID         int      `json:"id"`
	OwnerID    string   `json:"ownerId"`
	Recipients []string `json:"messageRecipients"`
	Subject    string   `json:"subject"`
	Viewed     string   `json:"viewedDatetime"`
}

// SecureMessagesResponse is the response for the list of secure messages
type SecureMessagesResponse struct {
	Messages      [][]SecureMessage `json:"data"`
	ErrorMessages []string          `json:"fMessages"`
	InfoMessages  []string          `json:"sMessages"`
}

// SecureMessages gets the list of messages for a user
func (c *Client) SecureMessages() (*SecureMessagesResponse, error) {
	url := fmt.Sprintf("%s/message/get_messages", c.BaseURL)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", defaultContentTypeHeader)
	req.Header.Add("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			// TODO: log event
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var v SecureMessagesResponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	if len(v.ErrorMessages) > 0 {
		return nil, errors.New("error listing secure messages")
	}

	return &v, nil
}

// SecureMessageReadResponse is the response for marking a message as read
type SecureMessageReadResponse struct {
	ErrorMessages []string `json:"fMessages"`
	InfoMessages  []string `json:"sMessages"`
}

// SecureMessageRead marks a message ID as read
func (c *Client) SecureMessageRead(id int) (*SecureMessageReadResponse, error) {
	form := url.Values{}
	form.Add("messageId", strconv.Itoa(id))

	url := fmt.Sprintf("%s/message/read_message", c.BaseURL)
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", defaultContentTypeHeader)
	req.Header.Add("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			// TODO: log event
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var v SecureMessageReadResponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	if len(v.ErrorMessages) > 0 {
		return nil, errors.New("error marking secure message as read")
	}

	return &v, nil
}

// SecureMessageCountResponse is a response for the count of new messages
type SecureMessageCountResponse struct {
	Count         int      `json:"data"`
	ErrorMessages []string `json:"fMessages"`
	InfoMessages  []string `json:"sMessages"`
}

// SecureMessagesCount gets the new messages count
func (c *Client) SecureMessagesCount() (*SecureMessageCountResponse, error) {
	url := fmt.Sprintf("%s/message/get_new_message_count", c.BaseURL)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", defaultContentTypeHeader)
	req.Header.Add("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			// TODO: log event
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var v SecureMessageCountResponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	if len(v.ErrorMessages) > 0 {
		return nil, errors.New("error listing secure messages")
	}

	return &v, nil
}

// NewMessageResponse is the response from filelocker sending a new secure message
type NewMessageResponse struct {
	ErrorMessages []string `json:"fMessages"`
	InfoMessages  []string `json:"sMessages"`
}

// NewSecureMessage sends a new secure message via filelocker
func (c *Client) NewSecureMessage(subject, msg string, recipients []string, expire time.Time) (*NewMessageResponse, error) {
	form := url.Values{}
	form.Add("requestOrigin", c.Origin)
	form.Add("subject", subject)
	form.Add("body", msg)
	form.Add("expiration", expire.Format("01/02/2006"))
	recipientIds := strings.Join(recipients, ",")
	form.Add("recipientIds", recipientIds)

	url := fmt.Sprintf("%s/message/create_message", c.BaseURL)
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", defaultContentTypeHeader)
	req.Header.Add("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			// TODO: log event
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var v NewMessageResponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	if len(v.ErrorMessages) > 0 {
		return nil, errors.New("error sending secure messages")
	}

	return &v, nil
}

// SecureMessagesDeleteResponse is the response for deleting a secure message
type SecureMessagesDeleteResponse struct {
	ErrorMessages []string `json:"fMessages"`
	InfoMessages  []string `json:"sMessages"`
}

// SecureMessagesDelete marks a message ID as read
func (c *Client) SecureMessagesDelete(ids []int) (*SecureMessagesDeleteResponse, error) {
	var idList []string
	for _, i := range ids {
		idList = append(idList, strconv.Itoa(i))
	}

	form := url.Values{}
	form.Add("messageIds", strings.Join(idList, ","))
	form.Add("requestOrigin", c.Origin)

	url := fmt.Sprintf("%s/message/delete_messages", c.BaseURL)
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", defaultContentTypeHeader)
	req.Header.Add("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			// TODO: log event
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var v SecureMessagesDeleteResponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	if len(v.ErrorMessages) > 0 {
		return nil, errors.New("error deleting secure messages")
	}

	return &v, nil
}
