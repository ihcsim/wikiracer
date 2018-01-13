package wikiracer

import (
	"fmt"
	"testing"
)

func TestFindPath(t *testing.T) {
	wiki = NewMockWiki()

	t.Run("Pages Exist", func(t *testing.T) {
		var testCases = []struct {
			origin      string
			destination string
			expected    string
		}{
			{origin: "Mike Tyson", destination: "Mike Tyson", expected: "Mike Tyson"},
			{origin: "Mike Tyson", destination: "Alexander the Great", expected: "Mike Tyson -> Alexander the Great"},
			{origin: "Mike Tyson", destination: "Apepi", expected: "Mike Tyson -> Alexander the Great -> Apepi"},
			{origin: "Mike Tyson", destination: "Greek language", expected: "Mike Tyson -> Alexander the Great -> Greek language"},
			{origin: "Mike Tyson", destination: "Fruit anatomy", expected: "Mike Tyson -> Alexander the Great -> Greek language -> Fruit anatomy"},
			{origin: "Mike Tyson", destination: "Segment", expected: "Mike Tyson -> Alexander the Great -> Greek language -> Fruit anatomy -> Segment"},
		}

		for id, testCase := range testCases {
			if actual := FindPath(testCase.origin, testCase.destination); testCase.expected != actual {
				t.Errorf("Mismatch path. Test case: %d\nExpected: %s\nActual: %s", id, testCase.expected, actual)
			}
		}
	})

	t.Run("Pages Don't Exist", func(t *testing.T) {
		var testCases = []struct {
			origin      string
			destination string
			expected    error
		}{
			{origin: "", expected: InvalidEmptyInput{}},
			{origin: "123456789", expected: InvalidEmptyInput{origin: "123456789"}},
			{origin: "123456789", destination: "Mike Tyson", expected: PageNotFound{Page{Title: "123456789"}}},
			{origin: "Mike Tyson", destination: "123456789", expected: PageNotFound{Page{Title: "123456789"}}},
			{origin: "Mike Tyson", destination: "Michael Jordan", expected: DestinationUnreachable{destination: "Michael Jordan"}},
		}

		for id, testCase := range testCases {
			actual := FindPath(testCase.origin, testCase.destination)
			if fmt.Sprintf("%s", testCase.expected) != actual {
				t.Errorf("Mismatch error. Test case: %d\nExpected: %s\nActual: %s", id, testCase.expected, actual)
			}
		}
	})

}
