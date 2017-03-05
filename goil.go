// Scrapping the face book

// Package goil provides an interface to iseplive.fr, it's currently a WIP
// Supported features:
// 	- Posting
//	- Deleting a post by ID
//	- Retrieving the list of all students along with their data, but not pictures
// Incoming
//	- Retrieving only one student's data
//	- Retrieving students pictures
//	- Retrieving publications
//	- Commenting
//	- Liking / Disliking
package goil

import (
	"bytes"
	"github.com/gocarina/gocsv"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Session struct {
	Client *http.Client
}

func Login(username string, password string, client *http.Client) (*Session, error) {
	// Init
	sess := &Session{}

	// Prepare login form
	loginForm := url.Values{}
	loginForm.Add("password", password)
	loginForm.Add("username", username)

	// Make cookiejar
	cookieJar, _ := cookiejar.New(nil)
	client.Jar = cookieJar

	// Prepare request
	req, err := http.NewRequest("POST", "http://iseplive.fr/signin/", bytes.NewBufferString(loginForm.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	sess.Client = client
	return sess, nil
}

func ExportToCSV(studentList StudentList, writer io.Writer) error {
	err := gocsv.Marshal(&studentList, writer)
	return err
}
