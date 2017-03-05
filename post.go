// post.go manages the posting part of goil

package goil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// The post URI
const postURI string = "http://iseplive.fr/post/add"

type Category uint8

const (
	Divers     Category = 8
	News                = 6
	Photos              = 1
	Videos              = 2
	Journal             = 3
	Gazettes            = 10
	Podcasts            = 4
	Evenements          = 5
	Sondages            = 9
	Annales             = 7
)

func (c Category) format() string {
	return strconv.FormatUint(uint64(c), 10)
}

type Group uint8

const (
	NoGroup    Group = 0
	HustleISEP       = 46 // Verifier
)

func (g Group) format() string {
	return strconv.FormatUint(uint64(g), 10)
}

// A post on iseplive.fr
type Post struct {
	Message  string   // message
	Category Category // category
	Group    Group    // association/group
	Official bool
	Private  bool
	Dislike  bool // Activate or not the dislike button

	// Attachments
	Photos    []string // name="attachment_photo[]"
	Videos    []string // name="attachment_video[]"
	Audio     []string // name="attachment_audio[]"
	Documents []string // name="attachment_file[]"

	// Event
	EventTitle string
	EventStart time.Time // Format in "DD/MM/YYYY à HH:MM"
	EventEnd   time.Time // Format in "DD/MM/YYYY à HH:MM"

	// Survey
	SurveyQuestion string
	SurveyEnd      time.Time // Format in "DD/MM/YYYY à HH:MM"
	SurveyAnswers  []string
}

// TODO Check if when post with no group "official" field exists
func NewPost(message string, category Category, private bool, dislike bool) *Post {
	return &Post{Message: message, Category: category, Group: NoGroup, Private: private, Dislike: dislike}
}

// Post as a group
func (p *Post) PostAs(group Group, official bool) {
	p.Group = group
	p.Official = official
}

func (p *Post) AddEvent(title string, start time.Time, end time.Time) {
	p.EventTitle = title
	p.EventStart = start
	p.EventEnd = end
}

// Example:
// post.AddSurvey("Who's a good boy?", time.Now.Add(3*time.Hour), "I am", "How should I know ? I'm a dog", "This question isn't well defined, can't answer")
func (p *Post) AddSurvey(question string, end time.Time, answers ...string) {
	p.SurveyQuestion = question
	p.SurveyEnd = end
	p.SurveyAnswers = answers
}

type FileType uint8

const (
	Photo    FileType = 0
	Video    FileType = 1
	Audio    FileType = 2
	Document FileType = 3
)

func (p *Post) AttachPhoto(filepath string) {
	p.Photos = append(p.Photos, filepath)
}
func (p *Post) AttachVideo(filepath string) {
	p.Photos = append(p.Photos, filepath)
}
func (p *Post) AttachAudio(filepath string) {
	p.Photos = append(p.Photos, filepath)
}
func (p *Post) AttachDocument(filepath string) {
	p.Photos = append(p.Photos, filepath)
}

// Generic file add
/*
func (p *Post) AddFile(filetype fileType, filepath string) {
	switch filetype {
		case Photo:
			p.AttachPhoto(filepath)
		case Video:
			p.AttachVideo(filepath)
		case Audio:
			p.AttachAudio(filepath)
		case Document:
			p.AttachDocument(filepath)
	}
} */

func bts(from bool) string {
	var output string
	if from {
		output = "1"
	} else {
		output = "0"
	}
	return output
}

func (s *Session) PublishPost(post *Post) error {

	// Prepare the body of the request
	body := &bytes.Buffer{}
	// Prepare it as a multipart/form-data
	writer := multipart.NewWriter(body)

	// Prepare the non-file parameters
	params := map[string]string{
		"message":  post.Message,
		"category": post.Category.format(),
		"group":    post.Group.format(),
		"private":  bts(post.Dislike),
		"dislike":  bts(post.Dislike),
	}
	if post.Group != 0 {
		params["official"] = bts(post.Official)
	}

	// Is there an event or a survey ? Then add them too
	//TODO

	// Add them
	for key, val := range params {
		err := writer.WriteField(key, val)
		if err != nil {
			return err
		}
	}

	// Helper func that adds each attachment to its corresponding fields
	addToForm := func(name string, paths []string) error {
		// Now add the files for the pictures
		for _, path := range paths {
			// Open one
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Create the form field
			part, err := writer.CreateFormFile(name, filepath.Base(path))
			if err != nil {
				return err
			}

			// Copy the file to the form
			_, err = io.Copy(part, file)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// Add the files
	err := addToForm("attachment_photo[]", post.Photos)
	if err != nil {
		return err
	}
	err = addToForm("attachment_video[]", post.Videos)
	if err != nil {
		return err
	}
	err = addToForm("attachment_audio[]", post.Audio)
	if err != nil {
		return err
	}
	err = addToForm("attachment_file[]", post.Documents)
	if err != nil {
		return err
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return err
	}

	// Create the request
	req, err := http.NewRequest("POST", postURI, body)
	if err != nil {
		return err
	}

	// Add the right header
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute it
	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}

	// Check for good feedback
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status Code of response isn't right: %d instead of 200", resp.StatusCode))
	}

	return nil

}

const (
	delBasePath string = "http://iseplive.fr/ajax/post/"
	delEndPath  string = "/delete"
)

// GET http://iseplive.fr/ajax/post/POSTID/delete
// TODO: More info can be retrieved before redirect, notably the success/failure of the delete. retrieve it.
func (s *Session) DeletePost(postID uint) error {
	path := delBasePath + strconv.FormatUint(uint64(postID), 10) + delEndPath
	resp, err := s.Client.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status Code of response isn't right: %d instead of 200", resp.StatusCode))
	}

	fmt.Println(resp.Body)
	return nil
}
