// Package goil provides an interface to iseplive.fr, it's currently a WIP
// Supported features:
// 	- Posting
//	- Deleting a post by ID
//	- Retrieving a student by URL
//	- Retrieving a list of all students along with their data
// Incoming
//	- Retrieving only one student's data
//	- Retrieving students pictures
//	- Retrieving publications
//	- Commenting
//	- Liking / Disliking
package goil

import (
	"net/url"
)

const BaseURLString string = "http://iseplive.fr/"

var BaseURL *url.URL = &url.URL{
	Scheme: "http",
	Host:   "iseplive.fr",
}
