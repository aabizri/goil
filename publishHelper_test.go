package goil

import "testing"

func Test_AvailableGroups(t *testing.T) {
	skipIfNoSession(t)

	groups, err := session.AvailableGroups()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("These are the available groups")
	for id, name := range groups {
		t.Logf("%d: %s\n", id, name)
	}
}
