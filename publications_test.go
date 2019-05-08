package goil

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var group = flag.Uint("group", 0, "Group ID to publish as")
var picpath = flag.String("pic", "", "Path to picture to publish")

func Test_PostPublication_Basic(t *testing.T) {
	skipIfNoSession(t)

	message := fmt.Sprintf("Testing at %v", time.Now())
	post := CreatePublication(message, Divers)

	err := session.PostPublication(post)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_PostPublication_AsGroup(t *testing.T) {
	skipIfNoSession(t)

	if *group == 0 {
		t.Skip("No group indicated, skipping...")
	}

	publi := CreatePublication("Lorem Ipsum", Divers)
	publi.PublishAs(Group(*group), true)

	err := session.PostPublication(publi)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_PostPublication_WithPicture(t *testing.T) {
	skipIfNoSession(t)

	if *picpath == "" {
		t.Skip("No photo filepath indicated, skipping...")
	}

	publi := CreatePublication("Picture Test", Divers)

	// Open the attachment file
	file, err := os.Open(*picpath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Add the picture
	publi.Attachments.AttachPicture(filepath.Base(*picpath), file)

	// Publish the publication
	err = session.PostPublication(publi)
	if err != nil {
		t.Fatal(err)
	}
}
