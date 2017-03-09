package goil

import "testing"

func Test_AvailableGroups(t *testing.T) {
	skipIfNoSession(t)

	groups, err := session.GetAvailableGroups()
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) == 0 {
		t.Fatal("No groups returned when there should at least be one")
	}

	t.Log("These are the available groups")
	for id, name := range groups {
		t.Logf("%d: %s\n", id, name)
	}
}
