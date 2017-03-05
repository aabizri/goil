// Testing the face book scrapper

package goil

import "testing"

var slist StudentList

func Test_GetStudentList(t *testing.T) {
	skipIfNoSession(t)

	var err error
	slist, err = session.GetStudentList()
	if err != nil {
		t.Error(err)
	}
}

/*
func Test_ExportToCSV(t *testing.T) {
	sess
}*/
