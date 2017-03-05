// Testing the face book scrapper

package goil

import "testing"

func Test_GetStudentList(t *testing.T) {
	skipIfNoSession(t)
	
	_, err := session.GetStudentList()
	if err != nil {
		t.Error(err)
	}
}
