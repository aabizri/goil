package goil

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocarina/gocsv"
	"io"
)

// Return list of every student with data
func (s *Session) GetStudentList() (StudentList, error) {
	body, err := s.getStudentsListPage()
	if err != nil {
		return StudentList{}, err
	}

	links, err := extractLinksToStudentPages(body)
	if err != nil {
		return StudentList{}, err
	}
	body.Close() // Closing the open body

	studentList := make(StudentList, len(links))
	for i, k := range links {
		student, err := s.GetStudentByLink(k)
		if err != nil {
			fmt.Printf("ERROR: %v", err.Error())
		}
		studentList[i] = student
	}

	return studentList, err
}

// Get the page with student links
// http://iseplive.fr/students
func (s *Session) getStudentsListPage() (io.ReadCloser, error) {
	resp, err := s.Client.Get("http://iseplive.fr/students")
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// Extract links to each student page
// --> http://iseplive.fr/student/example
func extractLinksToStudentPages(body io.ReadCloser) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	// Selection object
	sel := doc.Find(".students-promo").ChildrenFiltered("a")

	// Make the slice the length of the selection
	links := make([]string, sel.Length())

	// For each, append the link to the slice
	sel.Each(func(i int, s *goquery.Selection) {
		//link := s.Find
		link, _ := s.Attr("href")
		links[i] = link
	})

	return links, nil
}

// Export StudentList to CSV
func (sl StudentList) WriteCSV(writer io.Writer) error {
	return gocsv.Marshal(sl, writer)
}

// Export StudentList to JSON
func (sl StudentList) WriteJSON(writer io.Writer) error {
	enc := json.NewEncoder(writer)
	return enc.Encode(&sl)
}
