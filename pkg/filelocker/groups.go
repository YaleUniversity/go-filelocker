package filelocker

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Group is a user's group respresentation from filelocker
type Group struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// GroupsResponse is the response from filelocker for a list of the user's groups
type GroupsResponse struct {
	Groups        []Group  `xml:"data>group"`
	ErrorMessages []string `xml:"messages>error"`
	InfoMessages  []string `xml:"messages>info"`
}

// Groups lists a users groups in filelocker.  It requires an authenticated client.
func (c *Client) Groups() (*GroupsResponse, error) {
	form := url.Values{}
	form.Add("format", "cli")

	url := fmt.Sprintf("%s/account/get_groups", c.BaseURL)
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", defaultContentTypeHeader)
	req.Header.Add("Accept", defaultAcceptHeader)

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

	var v GroupsResponse
	err = xml.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	if len(v.ErrorMessages) > 0 {
		return &v, errors.New("error listing groups")
	}

	return &v, nil
}
