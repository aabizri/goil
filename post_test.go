// Testing the publication system

package goil

import (
	"flag"
	"testing"
)

var group = flag.Uint("group", 0, "Group ID to post as")
var photopath = flag.String("photo", "", "Path to photo to publish")
var postid = flag.Uint("postid", 0, "Post ID to delete")

func Test_PublishPost(t *testing.T) {
	skipIfNoSession(t)

	post := NewPost("Lorem Ipsum", Divers, true, false)
	//post.PostAs(NoGroup,false)

	err := session.PublishPost(post)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_PublishPostAsGroup(t *testing.T) {
	skipIfNoSession(t)

	if *group == 0 {
		t.Skip("No group indicated, skipping...")
	}

	post := NewPost("Lorem Ipsum", Divers, true, false)
	post.PostAs(Group(*group), true)

	err := session.PublishPost(post)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_PublishPhoto(t *testing.T) {
	skipIfNoSession(t)

	if *photopath == "" {
		t.Skip("No photo filepath indicated, skipping...")
	}

	post := NewPost("Picture Test", Divers, true, false)
	post.AttachPhoto(*photopath)

	err := session.PublishPost(post)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Delete(t *testing.T) {
	skipIfNoSession(t)

	if *postid == 0 {
		t.Skip("No post ID to delete indicated, skipping ...")
	}
	err := session.DeletePost(*postid)
	if err != nil {
		t.Fatal(err)
	}
}
