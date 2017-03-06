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

var BaseURL *url.URL = &url.URL{
	Scheme: "http",
	Host:   "iseplive.fr",
}

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

// Login using only a cookie
func CreateSessionByCookieValue(cookieValue string, client *http.Client) *Session {
	// Create a cookie
	cookie := http.Cookie{
		Name:  "login",
		Value: cookieValue,
	}

	// Create a slice around it
	cookies := make([]*http.Cookie, 1)
	cookies[0] = &cookie

	// Make cookiejar
	cookieJar, _ := cookiejar.New(nil)
	client.Jar = cookieJar

	// Add cookie to the cookie jar
	cookieJar.SetCookies(BaseURL, cookies)

	// Create a session
	sess := &Session{client}
	return sess
}

// Retrieve the login cookie
func (s *Session) Cookie() (*http.Cookie, error) {
	// Iterate through the cookies in order to find the right one
	for _, k := range s.Client.Jar.Cookies(BaseURL) {
		if k.Name == "login" {
			return k, nil
		}
	}

	// If no logincookie was found, return an error
	return nil, errors.New("No valid login cookie found :(")
}
