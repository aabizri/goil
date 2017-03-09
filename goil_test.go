package goil

import (
	"flag"
	"net/http"
	"testing"
)

var username = flag.String("u", "", "Username on ISEPLive")
var password = flag.String("p", "", "Password on ISEPLive")
var loginCookie = flag.String("c", "", "Login cookie (Alternative to username/password combo)")
var session *Session

func init() {
	flag.Parse()

	// Create a session if either a username/password combo or a loginCookie is provided
	var err error
	if *username != "" && *password != "" {
		session, err = Login(*username, *password, &http.Client{})
	} else if *loginCookie != "" {
		session = CreateSessionByCookieValue(*loginCookie, &http.Client{})
	}
	if err != nil {
		panic(err)
	}
}

func skipIfNoSession(t *testing.T) {
	if session == nil {
		t.Skip("No session could be created, probably missing username/password combo or cookie")
	}
}
