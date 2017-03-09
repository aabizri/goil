package goil

import (
	"fmt"
	"strconv"
	"time"
)

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

func (c Category) Check() error {
	var err error
	if c < 1 || c > 10 {
		err = fmt.Errorf("Given category (%d) not in acceptable range [1;10]", c)
	}
	return err
}

type Group uint8

// Groups ID
const (
	NoGroup    Group = 0
	HustleISEP       = 46
	// Please contact me if you have the other codes, or just do a pull request !
)

// format() stringifies a group
func (g Group) format() string {
	return strconv.FormatUint(uint64(g), 10)
}

// timeLayout is the necessary layout to publish dates on the endpoint
const timeLayout string = "02/01/2006 à 15:04"

// An event
// Currently only filled in / used when publishing
type Event struct {
	// The name of the event, mandatory
	Name string

	// Start and end times
	Start time.Time
	End   time.Time

	// Start and end times will be formatted as such: "DD/MM/YYYY à HH:MM"
}

// populated() returns true if an event is populated
func (e Event) populated() bool {
	return e.Name != ""
}

// Check() checks an event's validity
func (e Event) Check() error {
	switch {
	case e.Name == "":
		return fmt.Errorf("No event name indicated")

	case e.Start.Before(time.Now()):
		return fmt.Errorf("Event start time (%s) is before(!) current time (%s)", e.Start, time.Now())

	case e.End.Before(time.Now()):
		return fmt.Errorf("Event end time (%s) is before(!) current time (%s)", e.End, time.Now())

	case e.End.Before(e.Start):
		return fmt.Errorf("Evend end time (%s) is before(!) start time (%s)", e.End, e.Start)
	}

	return nil
}

// A survey
// Currently only filled in / used when publishing
type Survey struct {
	// Mandatory
	Question string

	End time.Time // Format in "DD/MM/YYYY à HH:MM"

	Answers []string

	// Whether the survey accepts multiple answers
	Multiple bool
}

// populated() returns true if an event is populated
func (s Survey) populated() bool {
	return s.Question != ""
}

// Check validity
func (s Survey) Check() error {
	switch {
	case s.Question == "":
		return fmt.Errorf("No question included in survey")

	case s.End.Before(time.Now()):
		return fmt.Errorf("End time (%s) for survey is before(!) current time (%s)", s.End, time.Now())

	case len(s.Answers) < 2:
		return fmt.Errorf("A survey needs at least two answers to be a survey, duh !")
	}

	return nil
}

// A publication's ID
type PublicationID uint

// A publication on Iseplive
// Not all fields are filled in when getting publications
type Publication struct {
	// Publication's ID
	// Not populated when posting, currently unused
	ID PublicationID

	// Message is the text body of a publication
	// Mandatory for publishing
	Message string

	// The category of the publication meant to be published
	// Only filled in on publishing
	// Mandatory for publishing
	Category Category

	// The Group responsible for the publication
	// Only filled in on publishing
	Group Group

	// Whether that messages is official in reference to the group
	// Only used when publishing
	// If Group == 0 then it has no effect
	Official bool

	// Whether the publication should be private
	Private bool

	// Whether the dislike button should be activated
	Dislike bool

	// Event
	Event Event

	// Survey
	Survey Survey

	// Attachments paths for upload
	Attachments Attachments
}

func (p *Publication) Check() error {

	switch {
	case len(p.Message) == 0:
		return fmt.Errorf("Publication message is empty. It shouldn't")

	case p.Category.Check() != nil:
		return fmt.Errorf("Publication's category (%d) is invalid: %s", p.Category, p.Category.Check())

	case p.Event.populated() && p.Event.Check() != nil:
		return fmt.Errorf("Publication's event is invalid: %s", p.Event.Check())

	case p.Survey.populated() && p.Survey.Check() != nil:
		return fmt.Errorf("Publication's survey is invalid: %s", p.Survey.Check())

		/*case err := p.Attachments.Check(); err != nil {
		return fmt.Errorf("Publication's attachments are invalid: %s",err.Error())*/
	}

	return nil
}
