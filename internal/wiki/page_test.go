package wiki

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddPages(t *testing.T) {
	var (
		p1 = &Page{Title: "Page 1"}
		p2 = &Page{Title: "Page 2"}
		p3 = &Page{Title: "Page 3"}
	)
	actual := NewPath()
	actual.AddPage(p1)
	actual.AddPage(p2)
	actual.AddPage(p3)

	expected := NewPath()
	expected.sequence = []*Page{
		&Page{Title: "Page 1"},
		&Page{Title: "Page 2"},
		&Page{Title: "Page 3"},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Mismatch result.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestPathString(t *testing.T) {
	path := NewPath()
	path.sequence = []*Page{
		&Page{Title: "Title 0"},
		&Page{Title: "Title 1"},
		&Page{Title: "Title 2"},
	}

	expected := "Title 0" + pathDelimiter + "Title 1" + pathDelimiter + "Title 2"
	if actual := fmt.Sprintf("%s", path); expected != actual {
		t.Errorf("Mismatch result. Expected %q. Actual %q", expected, actual)
	}
}

func TestClone(t *testing.T) {
	path := NewPath()
	path.sequence = []*Page{
		&Page{Title: "Title 0"},
		&Page{Title: "Title 1"},
		&Page{Title: "Title 2"},
	}

	expected := &Path{
		sequence: []*Page{&Page{Title: "Title 0"}, &Page{Title: "Title 1"}, &Page{Title: "Title 2"}},
	}
	actual := NewPath()
	actual.Clone(path)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Mismatch result. Expected %q. Actual %q", expected, actual)
	}
}
