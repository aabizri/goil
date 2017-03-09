package goil

import "testing"

const studentListMinimumLength int = 0

func Test_GetStudentList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	skipIfNoSession(t)

	studentList, err := session.GetStudentList()
	if err != nil {
		t.Fatal(err)
	}
	if len(studentList) <= studentListMinimumLength {
		t.Errorf("studentList length should at least be of %d but it is of %d", studentListMinimumLength, len(studentList))
	}
}
