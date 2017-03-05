// General testing system

package goil

import (
	//"fmt"
	"net/http"
	//"os"
	"flag"
	"testing"
)

var username = flag.String("uname", "", "Username on ISEPLive")
var password = flag.String("pass", "", "Password on ISEPLive")
var loginCookie = flag.String("cookie", "", "Login cookie (Alternative to username/password combo)")
var session *Session

func init() {
	flag.Parse()
	/*if (*username == "" || *password == "") && loginCookie=="" {
		return
	}*/
	if *username != "" && *password != "" {
		var err error
		session, err = Login(*username, *password, &http.Client{})
		if err != nil {
			panic(err)
		}
	} /*else {
		// NOT YET IMPLEMENTED: Directly get the login cookie
		// use session.Jar.SetCookies(http://iseplive.fr,cookies)
	}*/
}

func skipIfNoSession(t *testing.T) {
	if session == nil {
		t.Skip("No session could be created, probably missing username/password combo or cookie")
	}
}

func Test_Login(t *testing.T) {
	if *username == "" || *password == "" {
		t.Skip("Username or Password missing")
	}

	sess, err := Login(*username, *password, &http.Client{})
	if err != nil {
		t.Error(err)
	}

	t.Logf("We got cookies: %v", sess.Client.Jar)
}
