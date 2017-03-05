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
	"errors"
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

// Convenience wrapper around Login & *Session.Cookie
func LoginRetrieveCookie(username string, password string) (*http.Cookie, error) {
	// Login
	sess, err := Login(username, password, &http.Client{})
	if err != nil {
		return nil, err
	}

	// Retrieve cookie
	return sess.Cookie()
}

// Retrieve the login cookie
func (s *Session) Cookie() (*http.Cookie, error) {
	// Parse the url
	url, err := url.Parse("http://iseplive.fr")
	if err != nil {
		return nil, err
	}

	// Iterate through the cookies in order to find the right one
	for _, k := range s.Client.Jar.Cookies(url) {
		if k.Name == "login" {
			return k, nil
		}
	}

	// If no logincookie was found, return an error
	return nil, errors.New("No valid login cookie found :(")
}
