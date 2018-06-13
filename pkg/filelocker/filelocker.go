package filelocker

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// Client is a filelocker client
type Client struct {
	BaseURL  *url.URL     // base URL for "API" calls to filelocker.  ie. https://files.yale.edu/cli/
	Client   *http.Client // HTTP client used to communicate with the "API".
	Errors   []string
	Messages []string
	Origin   string
}

const (
	defaultAcceptHeader      = "text/xml"
	defaultContentTypeHeader = "application/x-www-form-urlencoded"
)

// NewClient returns a new Filelocker "API" client. If a nil httpClient is
// provided, http.DefaultClient will be used with a 30s timeout.  A userID, apiKey,
// and baseURL must also be passed and a session will be established.
func NewClient(userID, apiKey, baseURL string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	if httpClient.Jar == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		httpClient.Jar = jar
	}

	bURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Add("CLIkey", apiKey)
	form.Add("userId", userID)

	url := fmt.Sprintf("%s/cli/CLI_login", baseURL)
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", defaultContentTypeHeader)
	req.Header.Add("Accept", defaultAcceptHeader)

	resp, err := httpClient.Do(req)
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

	type Result struct {
		ErrorMessages []string `xml:"messages>error"`
		InfoMessages  []string `xml:"messages>info"`
	}
	var v Result
	err = xml.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	client := Client{
		Client:   httpClient,
		BaseURL:  bURL,
		Errors:   v.ErrorMessages,
		Messages: v.InfoMessages,
	}

	if len(v.ErrorMessages) > 0 || len(v.InfoMessages) == 0 {
		return &client, errors.New("error logging into filelocker")
	}
	client.Origin = v.InfoMessages[0]

	return &client, nil
}
