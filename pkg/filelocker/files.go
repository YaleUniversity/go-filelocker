package filelocker

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// File is a file respresentation from filelocker
type File struct {
	ID           string `xml:"id,attr"`
	Name         string `xml:"name,attr"`
	Size         int    `xml:"size,attr"`
	PassedAvScan bool   `xml:"passedAvScan,attr"`
}

// FilesResponse is the response from filelocker for a list of the user's files
type FilesResponse struct {
	Files         []File   `xml:"data>file"`
	ErrorMessages []string `xml:"messages>error"`
	InfoMessages  []string `xml:"messages>info"`
}

// Files lists the uploaded files in filelocker.  It requires an authenticated client.
func (c *Client) Files() (*FilesResponse, error) {
	form := url.Values{}
	form.Add("format", "cli")

	url := fmt.Sprintf("%s/file/get_user_file_list", c.BaseURL)
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

	var v FilesResponse
	err = xml.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	if len(v.ErrorMessages) > 0 {
		return &v, errors.New("error listing file")
	}

	return &v, nil
}

// UploadResponse is the response from filelocker upload
type UploadResponse struct {
	ErrorMessages []string `xml:"messages>error"`
	InfoMessages  []string `xml:"messages>info"`
	File          File     `xml:"data>file"`
}

// Upload takes an io.ReadCloser , reads and uploads from it and then returns the response
func (c *Client) Upload(name, notes string, scan bool, f io.Reader) (*UploadResponse, error) {
	// yucky to take a Reader and then just read all the bytes =(
	file, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/file/upload", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(file))
	if err != nil {
		return nil, err
	}

	// setup query parameters
	params := req.URL.Query()
	params.Add("format", "cli")
	params.Add("fileName", name)

	if scan {
		params.Add("scanFile", strconv.FormatBool(scan))
	}

	if notes == "" {
		params.Add("fileNotes", "Uploaded by #golang")
	} else {
		params.Add("fileNotes", notes)
	}

	req.URL.RawQuery = params.Encode()

	// setup request headers
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Accept", defaultAcceptHeader)
	req.Header.Add("Content-Length", strconv.Itoa(len(file)))
	req.Header.Add("X-File-Name", name)

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

	var v UploadResponse
	err = xml.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	if len(v.ErrorMessages) > 0 {
		return &v, errors.New("error uploading")
	}

	return &v, nil
}

// DeleteResponse is the response from filelocker for a list of the user's groups
type DeleteResponse struct {
	ErrorMessages []string `xml:"messages>error"`
	InfoMessages  []string `xml:"messages>info"`
}

// Delete removes a file from filelocker.  It requires an authenticated client.
func (c *Client) Delete(files []string) (*DeleteResponse, error) {
	form := url.Values{}
	form.Add("format", "cli")
	form.Add("requestOrigin", c.Origin)
	fileIDs := strings.Join(files, ",")
	form.Add("fileIds", fileIDs)

	url := fmt.Sprintf("%s/file/delete_files", c.BaseURL)
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

	var v DeleteResponse
	err = xml.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	if len(v.ErrorMessages) > 0 {
		return &v, errors.New("error deleting file")
	}

	return &v, nil
}
