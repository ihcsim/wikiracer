package wikiracer

import (
	"fmt"
	"testing"

	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/crawler"
	"github.com/ihcsim/wikiracer/internal/validator"
	"github.com/ihcsim/wikiracer/internal/wiki"
	"github.com/ihcsim/wikiracer/test"
)

func TestFindPath(t *testing.T) {
	mockWiki := test.NewMockWiki()

	t.Run("Existent Pages", func(t *testing.T) {
		t.Run("Single Path", func(t *testing.T) {
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
				{origin: "Mike Tyson", destination: "Diodotus I", expected: "Mike Tyson -> Alexander the Great -> Diodotus I"},
				{origin: "Mike Tyson", destination: "1984 Summer Olympics", expected: "Mike Tyson -> 1984 Summer Olympics"},
				{origin: "Mike Tyson", destination: "7-Eleven", expected: "Mike Tyson -> 1984 Summer Olympics -> 7-Eleven"},
				{origin: "Mike Tyson", destination: "Big C", expected: "Mike Tyson -> 1984 Summer Olympics -> 7-Eleven -> Big C"},
				{origin: "Mike Tyson", destination: "Calgary", expected: "Mike Tyson -> 1984 Summer Olympics -> 7-Eleven -> Calgary"},
				{origin: "Mike Tyson", destination: "Eurocash", expected: "Mike Tyson -> 1984 Summer Olympics -> 7-Eleven -> Eurocash"},
				{origin: "Mike Tyson", destination: "Małpka Express", expected: "Mike Tyson -> 1984 Summer Olympics -> 7-Eleven -> Eurocash -> Małpka Express"},
				{origin: "Mike Tyson", destination: "Tea", expected: "Mike Tyson -> 1984 Summer Olympics -> 7-Eleven -> Eurocash -> Tea"},
				{origin: "Mike Tyson", destination: "Afghanistan", expected: "Mike Tyson -> 1984 Summer Olympics -> Afghanistan"},
			}

			for id, testCase := range testCases {
				racer := &Racer{
					Crawler:   crawler.NewForward(mockWiki),
					Validator: &validator.InputValidator{mockWiki},
				}

				if actual := racer.FindPath(testCase.origin, testCase.destination); testCase.expected != actual {
					t.Errorf("Mismatch path. Test case: %d\nExpected: %s\nActual: %s", id, testCase.expected, actual)
				}
			}
		})

		t.Run("Multiple Paths", func(t *testing.T) {
			var testCases = []struct {
				origin      string
				destination string
				expected    []string
			}{
				{origin: "Mike Tyson", destination: "Vancouver",
					expected: []string{
						"Mike Tyson -> Alexander the Great -> Greek language -> Fruit anatomy -> Segment -> Vancouver",
						"Mike Tyson -> 1984 Summer Olympics -> 7-Eleven -> Big C -> Vancouver"}},
			}

			for id, testCase := range testCases {
				racer := &Racer{
					Crawler:   crawler.NewForward(mockWiki),
					Validator: &validator.InputValidator{mockWiki},
				}

				actual := racer.FindPath(testCase.origin, testCase.destination)
				passed := false
				for _, option := range testCase.expected {
					if option == actual {
						passed = true
						break
					}
				}

				if !passed {
					t.Errorf("Mismatch path. Test case: %d\nExpected either one of: %v\nActual: %s", id, testCase.expected, actual)
				}
			}
		})
	})

	t.Run("Non-Existent Pages", func(t *testing.T) {
		var testCases = []struct {
			origin      string
			destination string
			expected    error
		}{
			{origin: "", expected: errors.InvalidEmptyInput{}},
			{origin: "123456789", expected: errors.InvalidEmptyInput{Origin: "123456789"}},
			{origin: "123456789", destination: "Mike Tyson", expected: errors.PageNotFound{wiki.Page{Title: "123456789"}}},
			{origin: "Mike Tyson", destination: "123456789", expected: errors.PageNotFound{wiki.Page{Title: "123456789"}}},
			{origin: "Mike Tyson", destination: "Michael Jordan", expected: errors.DestinationUnreachable{Destination: "Michael Jordan"}},
		}

		for id, testCase := range testCases {
			racer := &Racer{
				Crawler:   crawler.NewForward(mockWiki),
				Validator: &validator.InputValidator{mockWiki},
			}

			actual := racer.FindPath(testCase.origin, testCase.destination)
			if fmt.Sprintf("%s", testCase.expected) != actual {
				t.Errorf("Mismatch error. Test case: %d\nExpected: %s\nActual: %s", id, testCase.expected, actual)
			}
		}
	})
}
