package wikiracer

import (
	"fmt"
	"testing"
)

func TestPageNotFound(t *testing.T) {
	title := "non-existent page"
	p := Page{Title: title}
	expected := fmt.Errorf("%s: %s", ErrPrefixPageNotFound, title)
	actual := PageNotFound{p}

	if fmt.Sprintf("%s", expected) != fmt.Sprintf("%s", actual) {
		t.Errorf("Mismatch errors.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestInvalidEmptyInput(t *testing.T) {
	origin, destination := "Mike Tyson", "Donald Duck"
	expected := fmt.Errorf("%s: (%s, %s)", ErrPrefixInvalidEmptyInputs, origin, destination)
	actual := InvalidEmptyInput{origin: origin, destination: destination}

	if fmt.Sprintf("%s", actual) != fmt.Sprintf("%s", expected) {
		t.Errorf("Mismatch errors.\nExpected: %v\nActual: %v", expected, actual)
	}
}
