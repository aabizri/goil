package goilscrap

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strings"
)

// A Student
type Student struct {
	ID        string
	Promo     string
	Name      string
	ISEPEmail string
	Email     string
	Cell      string
	Birthday  string
	Quote     string
}

// The list of all students
type StudentList []Student

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
		student, err := s.scrapStudentPage(k)
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

// Scrap data from each student page
// from http://iseplive.fr/student/example
func (s *Session) scrapStudentPage(link string) (Student, error) {
	// Query the student page
	resp, err := s.Client.Get("http://iseplive.fr" + link)
	if err != nil {
		return Student{}, err
	}
	defer resp.Body.Close()

	// Parse the response's body
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Student{}, err
	}

	// Selection object
	sel := doc.Find(".profile-info")

	// Parse it
	student, err := parseProfileInfo(sel)

	return student, err

}

// Parse the profile information scrapped into the Student{} struct
func parseProfileInfo(sel *goquery.Selection) (student Student, err error) {
	// Name is in a <h1> tag, so we can extract it
	student.Name = sel.ChildrenFiltered("h1").Text()

	// Well, didn't think it would work but it does
	// It works this way: if an element, as dictated by goquery, is the name of a field, we know the next one is the value of that field, so we mark it for the next iteration of the loop.
	// If it is marked, then we match the marking to the value and voila.
	var nextIs string
	sel.Contents().Each(func(i int, s *goquery.Selection) {
		if nextIs == "" {
			// If there is no next indicated for the current element, check if the current element indicates a next
			fields := [7]string{"Promo", "Numéro", "Adresse ISEP", "Adresse email", "Tél", "Date de naissance", "Ma citation"}
			for _, k := range fields {
				if strings.Contains(s.Text(), k) {
					nextIs = k
				}
			}
		} else {
			switch nextIs {
			case "Promo":
				student.Promo = s.Text()
			case "Numéro":
				student.ID = s.Text()
			case "Adresse ISEP":
				student.ISEPEmail = s.Text()
			case "Adresse email":
				student.Email = s.Text()
			case "Tél":
				student.Cell = s.Text()
			case "Date de naissance":
				student.Birthday = s.Text()
			case "Ma citation":
				student.Quote = s.Text()
			}
			nextIs = ""
		}
	})

	return
}
