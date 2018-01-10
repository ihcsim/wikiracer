package wikiracer

import (
	"fmt"
	"testing"
)

func TestPageNotFound(t *testing.T) {
	title := "non-existent page"
	p := &Page{Title: title}
	expected := fmt.Errorf("%s: %s", pageNotFoundErrorPrefix, title)
	actual := PageNotFound{p}

	if fmt.Sprintf("%s", expected) != fmt.Sprintf("%s", actual) {
		t.Errorf("Mismatch errors.\nExpected: %v\nActual: %v", expected, actual)
	}
}
