package goil

import (
	"net/http"
	"testing"
)

func Test_Login(t *testing.T) {
	if *username == "" || *password == "" {
		t.Skip("Username or Password missing")
	}

	_, err := Login(*username, *password, &http.Client{})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Cookie(t *testing.T) {
	skipIfNoSession(t)

	cookie, err := session.Cookie()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Cookie: %s", cookie.String())
}
