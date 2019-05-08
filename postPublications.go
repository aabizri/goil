package goil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// TODO Check if when post with no group "official" field exists
func CreatePublication(message string, category Category) *Publication {
	return &Publication{Message: message, Category: category, Group: NoGroup}
}

// Post as a group
func (p *Publication) PublishAs(group Group, official bool) {
	p.Group = group
	p.Official = official
}

// Add an event to a publication
func (p *Publication) AddEvent(event Event) {
	p.Event = event
}

// CreateSurvey creates a survey structure
func CreateSurvey(question string, end time.Time, multiple bool, answers ...string) Survey {
	return Survey{
		Question: question,
		End:      end, Answers: answers,
		Multiple: multiple,
	}
}

// Adds a survey to a publication
func (p *Publication) AddSurvey(survey Survey) {
	p.Survey = survey
}

// Converts bool to string
func bts(from bool) string {
	var output string
	if from {
		output = "1"
	} else {
		output = "0"
	}
	return output
}

// Write the publication to the body of the request
func (publication *Publication) write(w io.Writer) (string, error) {

	// Prepare it as a multipart/form-data
	writer := multipart.NewWriter(w)

	// Prepare the base parameters
	params := map[string]string{
		"message":  publication.Message,
		"category": publication.Category.format(),
		"group":    publication.Group.format(),
		"private":  bts(publication.Dislike),
		"dislike":  bts(publication.Dislike),
	}

	// If there is no group indicated, do not take into account publication.Official
	if publication.Group != 0 {
		params["official"] = bts(publication.Official)
	}

	// Is there an event ? Then add it
	if publication.Event.populated() {
		params["event_title"] = publication.Event.Name
		params["event_start"] = publication.Event.Start.Format(timeLayout)
		params["event_end"] = publication.Event.End.Format(timeLayout)
	}

	// Is there a survey ? Then add it
	if publication.Survey.populated() {
		params["survey_question"] = publication.Survey.Question
		params["survey_end"] = publication.Survey.End.Format(timeLayout)
		params["survey_multiple"] = bts(publication.Survey.Multiple)
		// Answers will be added later, as there are multiple of them
	}

	// Add the key/value pairs to the multipart request
	for key, val := range params {
		err := writer.WriteField(key, val)
		if err != nil {
			return "", err
		}
	}

	// Add the survey answers
	if publication.Survey.populated() {
		for _, answer := range params {
			err := writer.WriteField("survey_answer[]", answer)
			if err != nil {
				return "", err
			}
		}
	}

	var err error

	// Add the attachments
	if publication.Attachments.Populated() {
		err = publication.Attachments.writeToMultipart(writer)
	}

	// Get content type
	contentType := writer.FormDataContentType()

	// Close the writer
	err = writer.Close()

	return contentType, err
}

// The postPublicationURI is the URI for posting a publication
const postPublicationURI string = BaseURLString + "post/add"

// Publish publication
func (s *Session) PostPublication(publication *Publication) error {
	// Check
	if s == nil {
		return fmt.Errorf("Given session pointer is nil")
	}

	// Prepare the body of the request
	body := &bytes.Buffer{}

	// Write to the body
	contentType, err := publication.write(body)
	if err != nil {
		return err
	}

	// Create the request
	req, err := http.NewRequest("POST", postPublicationURI, body)
	if err != nil {
		return err
	}

	// Add the right header
	req.Header.Set("Content-Type", contentType)

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
