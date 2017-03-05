package goil

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

// Find out what groups are available to the session
func (s *Session) AvailableGroups() ([]Group,error) {
	// Retrive main page
	resp, err := s.Client.Get("http://iseplive.fr/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// Prepare goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	
	// Select groups
	sel := doc.Find("publish-group-select").ChildrenFiltered("option")
	
	// Prepare return
	// NB: There is always at least one option: 0 - None
	groups := make([]Group, sel.Length())
	
	// For each option, add to list
	sel.Each(func(i int, s *goquery.Selection) {
		rawgroup, _ := s.Attr("value")
		group,_ := strconv.ParseUint(rawgroup,10,64)
		groups[i]=Group(group)
	})
	
	// Return
	return groups,err
}
