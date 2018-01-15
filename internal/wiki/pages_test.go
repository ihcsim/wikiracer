package wiki

import (
	"fmt"
	"reflect"
	"testing"
)

func TestConcat(t *testing.T) {
	actual := &Path{
		&Page{Title: "Title 0"},
		&Page{Title: "Title 1"},
		&Page{Title: "Title 2"},
	}

	more := &Path{
		&Page{Title: "Title 3"},
		&Page{Title: "Title 4"},
		&Page{Title: "Title 5"},
	}

	expected := &Path{
		&Page{Title: "Title 0"},
		&Page{Title: "Title 1"},
		&Page{Title: "Title 2"},
		&Page{Title: "Title 3"},
		&Page{Title: "Title 4"},
		&Page{Title: "Title 5"},
	}
	actual.Concat(more)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Mismatch result.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestAddPages(t *testing.T) {
	p1, p2, p3 := &Page{Title: "Page 1"}, &Page{Title: "Page 2"}, &Page{Title: "Page 3"}
	actual := &Path{}
	actual.AddPage(p1)
	actual.AddPage(p2)
	actual.AddPage(p3)

	expected := &Path{
		&Page{Title: "Page 1"},
		&Page{Title: "Page 2"},
		&Page{Title: "Page 3"},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Mismatch result.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestPathString(t *testing.T) {
	path := &Path{
		&Page{Title: "Title 0"},
		&Page{Title: "Title 1"},
		&Page{Title: "Title 2"},
	}

	expected := "Title 0" + pathDelimiter + "Title 1" + pathDelimiter + "Title 2"
	if actual := fmt.Sprintf("%s", path); expected != actual {
		t.Errorf("Mismatch result. Expected %q. Actual %q", expected, actual)
	}
}
