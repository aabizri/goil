package goilscrap

import (
	"bytes"
	"errors"
	"fmt"
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
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

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

	// Check if response is OK
	if resp.StatusCode != 302 {
		return nil, errors.New(fmt.Sprintf("Status code isn't OK: %d", resp.StatusCode))
	}

	sess.Client = client
	return sess, nil
}

func ExportToCSV(studentList StudentList, writer io.Writer) error {
	err := gocsv.Marshal(&studentList, writer)
	return err
}
