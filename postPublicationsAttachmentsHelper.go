package goil

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"sync"
)

// An attachment
type Attachment struct {
	// Filename sent to the server, optional
	BasePath string

	// File to be sent, be sure to check max length
	Reader io.Reader

	mutex *sync.Mutex
}

const (
	MB                         uint = 1 << (10 * 2)
	AttachmentMaxSizePictures       = 5 * MB
	AttachmentMaxSizeVideos         = 800 * MB
	AttachmentMaxSizeAudio          = 20 * MB
	AttachmentMaxSizeDocuments      = 10 * MB
)

// Attachments filepath
type Attachments struct {
	/* Restrictions
	- can't be more than 5Mo in size each
	- must be of PNG, JPEG or GIF format
	*/
	Pictures []Attachment

	/* Restrictions
	- can't be more than 800Mo in size each
	- must have an audio track
	*/
	Videos []Attachment

	/* Restrictions
	- can't be more than 20Mo in size each
	- must have an .mp3 extension (and probably must be mp3 themselves)
	*/
	Audio []Attachment

	/* Restrictions
	- can't be more than 10Mo in size each
	- must have an extension, (i.e must match regex "\.[a-z0-9]{2,4}")
	- but can't have a .jpg, .png, .gif, .mp3, .flv extension (i.e must not match regex "\.(jpg|png|gif|mp3|flv)")
	*/
	Documents []Attachment
}

// Checks if there are any attachments
func (as Attachments) Populated() bool {
	return len(as.Pictures)+len(as.Videos)+len(as.Audio)+len(as.Documents) != 0
}

// Attach a picture given a filepath and a io.Reader
func (as Attachments) AttachPicture(filename string, r io.Reader) {
	as.Pictures = append(as.Pictures, Attachment{BasePath: filename, Reader: r})
}

// Attach a video given a filepath and a io.Reader
func (as Attachments) AttachVideo(filename string, r io.Reader) {
	as.Videos = append(as.Videos, Attachment{BasePath: filename, Reader: r})
}

// Attach an audio file given a filepath and a io.Reader
func (as Attachments) AttachAudio(filename string, r io.Reader) {
	as.Audio = append(as.Audio, Attachment{BasePath: filename, Reader: r})
}

// Attach a document given a filepath and a io.Reader
func (as Attachments) AttachDocument(filename string, r io.Reader) {
	as.Documents = append(as.Documents, Attachment{BasePath: filename, Reader: r})
}

// Add a single attachment to the writer
func (a Attachment) writeToMultipart(w *multipart.Writer, key string, maxSize uint) error {
	// Lock the attachment
	a.mutex.Lock()

	// Create the form field
	part, err := w.CreateFormFile(key, filepath.Base(a.BasePath))
	if err != nil {
		return err
	}

	// Copy the file to the form
	reader := io.LimitReader(a.Reader, int64(maxSize))
	_, err = io.Copy(part, reader)
	return err
}

// Write attachments to the multipart.Writer
func (as Attachments) writeToMultipart(w *multipart.Writer) error {

	// Helper func that adds each attachment to its corresponding fields
	addToForm := func(key string, at []Attachment, maxSize uint) error {
		// Now add the files for the pictures
		for _, a := range at {
			a.writeToMultipart(w, key, maxSize)
		}
		return nil
	}

	// Add the files
	err := addToForm("attachment_photo[]", as.Pictures, AttachmentMaxSizePictures)
	if err != nil {
		return err
	}
	err = addToForm("attachment_video[]", as.Videos, AttachmentMaxSizeVideos)
	if err != nil {
		return err
	}
	err = addToForm("attachment_audio[]", as.Audio, AttachmentMaxSizeAudio)
	if err != nil {
		return err
	}
	err = addToForm("attachment_file[]", as.Documents, AttachmentMaxSizeDocuments)
	return err
}
