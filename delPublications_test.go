package goil

import (
	"flag"
	"testing"
)

var publicationid = flag.Uint("postid", 0, "Publication ID to delete")

func Test_DeletePublication(t *testing.T) {
	skipIfNoSession(t)

	if *publicationid == 0 {
		t.Skip("No post ID to delete indicated, skipping ...")
	}
	err := session.DeletePublication(PublicationID(*publicationid))
	if err != nil {
		t.Fatal(err)
	}
}
