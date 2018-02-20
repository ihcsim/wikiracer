package wikipedia

import "testing"

func TestFindPage(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	page, err := client.FindPage("Mike Tyson")
	if err != nil {
		t.Fatal(err)
	}

	if page == nil {
	}
}
