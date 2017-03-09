package goil

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Session struct {
	Client *http.Client
}

// Endpoint used by Login
const loginEndpoint string = BaseURLString + "signin/"

// Login creates a session given a username/password combo and a *http.Client
// Warning: It is your responsibility to add a timeout
func Login(username string, password string, client *http.Client) (*Session, error) {
	// Init
	sess := &Session{}

	// Prepare login form
	loginForm := url.Values{}
	loginForm.Add("password", password)
	loginForm.Add("username", username)

	// Make cookiejar
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client.Jar = cookieJar

	// Prepare request
	req, err := http.NewRequest("POST", loginEndpoint, bytes.NewBufferString(loginForm.Encode()))
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

// Endpoint used by Logout
const logoutEndpoint string = BaseURLString + "logout/"

// Logout does what its name says
func (s *Session) Logout() error {
	resp, err := s.Client.Get(logoutEndpoint)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Logout failed: Status code returned was %d instead of %d", resp.StatusCode, 200)
	}
	return nil
}

// CreateSessionByCookie creates a Session with the given http.Cookie and http.Client
func CreateSessionByCookie(cookie *http.Cookie, client *http.Client) *Session {
	// Create a slice around it
	cookies := make([]*http.Cookie, 1)
	cookies[0] = cookie

	// Make cookiejar
	cookieJar, _ := cookiejar.New(nil)
	client.Jar = cookieJar

	// Add cookie to the cookie jar
	cookieJar.SetCookies(BaseURL, cookies)

	// Create a session
	sess := &Session{client}
	return sess
}

// CreateSessionByCookieValue creates a Session with the given cookie value and http.Client
func CreateSessionByCookieValue(cookieValue string, client *http.Client) *Session {
	// Create a cookie
	cookie := &http.Cookie{
		Name:  "login",
		Value: cookieValue,
	}

	// Use CreateSessionByCookie & return its output
	return CreateSessionByCookie(cookie, client)
}

// Retrieve the login cookie from a session
func (s *Session) Cookie() (*http.Cookie, error) {
	if s == nil {
		return nil, fmt.Errorf("Can't retrieve a cookie from a non-existent session")
	}

	// Iterate through the cookies in order to find the right one
	for _, k := range s.Client.Jar.Cookies(BaseURL) {
		if k.Name == "login" {
			return k, nil
		}
	}

	// If no logincookie was found, return an error
	return nil, fmt.Errorf("No valid login cookie found :(")
}
