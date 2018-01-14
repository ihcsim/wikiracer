package wikiracer

import (
	"fmt"
	"testing"

	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/crawler"
	"github.com/ihcsim/wikiracer/internal/validator"
	"github.com/ihcsim/wikiracer/internal/wiki"
)

func TestFindPath(t *testing.T) {
	mockWiki := NewMockWiki()
	racer := &Racer{
		Crawler:   crawler.NewForward(mockWiki),
		Validator: &validator.InputValidator{mockWiki},
	}

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
			if actual := racer.FindPath(testCase.origin, testCase.destination); testCase.expected != actual {
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
			{origin: "", expected: errors.InvalidEmptyInput{}},
			{origin: "123456789", expected: errors.InvalidEmptyInput{Origin: "123456789"}},
			{origin: "123456789", destination: "Mike Tyson", expected: errors.PageNotFound{wiki.Page{Title: "123456789"}}},
			{origin: "Mike Tyson", destination: "123456789", expected: errors.PageNotFound{wiki.Page{Title: "123456789"}}},
			{origin: "Mike Tyson", destination: "Michael Jordan", expected: errors.DestinationUnreachable{Destination: "Michael Jordan"}},
		}

		for id, testCase := range testCases {
			actual := racer.FindPath(testCase.origin, testCase.destination)
			if fmt.Sprintf("%s", testCase.expected) != actual {
				t.Errorf("Mismatch error. Test case: %d\nExpected: %s\nActual: %s", id, testCase.expected, actual)
			}
		}
	})
}

type MockWiki struct {
	pages map[string]*wiki.Page
}

func NewMockWiki() *MockWiki {
	testData := map[string]*wiki.Page{
		"Alexander the Great": &wiki.Page{ID: 1000, Title: "Alexander the Great", Namespace: 0, Links: []string{"Apepi", "Greek language", "Diodotus I"}},
		"Apepi":               &wiki.Page{ID: 1005, Title: "Apepi", Namespace: 0},
		"Diodotus I":          &wiki.Page{ID: 1007, Title: "Diodotus I", Namespace: 0},
		"Fruit anatomy":       &wiki.Page{ID: 1001, Title: "Fruit anatomy", Namespace: 0, Links: []string{"Segment"}},
		"Greek language":      &wiki.Page{ID: 1002, Title: "Greek language", Namespace: 0, Links: []string{"Fruit anatomy"}},
		"Mike Tyson":          &wiki.Page{ID: 1003, Title: "Mike Tyson", Namespace: 0, Links: []string{"Alexander the Great"}},
		"Segment":             &wiki.Page{ID: 1004, Title: "Segment", Namespace: 0},
		"Michael Jordan":      &wiki.Page{ID: 1006, Title: "Michael Jordan", Namespace: 0},
	}
	return &MockWiki{pages: testData}
}

func (m *MockWiki) FindPage(title string) (*wiki.Page, error) {
	page, exist := m.pages[title]
	if !exist {
		return nil, errors.PageNotFound{wiki.Page{Title: title}}
	}

	return page, nil
}
