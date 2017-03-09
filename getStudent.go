package goil

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

// Get a Student given a relative link to iseplive.fr: "/student/alexandrebezri"
// endpoint: http://iseplive.fr/student/{studentIdentifier}
func (s *Session) GetStudentByLink(link string) (Student, error) {
	// Query the student page
	resp, err := s.Client.Get(BaseURLString + link)
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
