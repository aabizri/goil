package goil

import (
	"fmt"
	"strconv"
)

const (
	delBasePath           string = "/ajax/post/"
	delPublicationEndpath string = BaseURLString + delBasePath
	delEndPath            string = "/delete"
)

// GET http://iseplive.fr/ajax/post/PUBLICATIONID/delete
// TODO: More info can be retrieved before redirect, notably the success/failure of the delete. retrieve it.
func (s *Session) DeletePublication(publicationID PublicationID) error {
	path := delBasePath + strconv.FormatUint(uint64(publicationID), 10) + delEndPath
	resp, err := s.Client.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Status Code of response isn't right: %d instead of 200", resp.StatusCode)
	}

	fmt.Println(resp.Body)
	return nil
}
