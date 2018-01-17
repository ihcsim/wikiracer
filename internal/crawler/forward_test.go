package crawler

import (
	"testing"

	"github.com/ihcsim/wikiracer/test"
)

func TestDiscover(t *testing.T) {
	t.Run("Visited Pages", func(t *testing.T) {
		// these test cases only verify pages from origin to destination are marked as visited.
		// other goroutines might have been created for to discover other paths and aren't included in these tests.

		var testCases = []struct {
			origin      string
			destination string
			expected    []string
		}{
			{origin: "Mike Tyson", destination: "Alexander the Great", expected: []string{"Mike Tyson", "Alexander the Great"}},
			{origin: "Mike Tyson", destination: "1984 Summer Olympics", expected: []string{"Mike Tyson", "1984 Summer Olympics"}},
			{origin: "Mike Tyson", destination: "Apepi", expected: []string{"Mike Tyson", "Alexander the Great", "Apepi"}},
			{origin: "Mike Tyson", destination: "Greek language", expected: []string{"Mike Tyson", "Alexander the Great", "Greek language"}},
			{origin: "Mike Tyson", destination: "Fruit anatomy", expected: []string{"Mike Tyson", "Alexander the Great", "Greek language", "Fruit anatomy"}},
			{origin: "Mike Tyson", destination: "Segment", expected: []string{"Mike Tyson", "Alexander the Great", "Greek language", "Fruit anatomy", "Segment"}},
			{origin: "Mike Tyson", destination: "Diodotus I", expected: []string{"Mike Tyson", "Alexander the Great", "Diodotus I"}},
			{origin: "Mike Tyson", destination: "1984 Summer Olympics", expected: []string{"Mike Tyson", "1984 Summer Olympics"}},
			{origin: "Mike Tyson", destination: "7-Eleven", expected: []string{"Mike Tyson", "1984 Summer Olympics", "7-Eleven"}},
			{origin: "Mike Tyson", destination: "Big C", expected: []string{"Mike Tyson", "1984 Summer Olympics", "7-Eleven", "Big C"}},
			{origin: "Mike Tyson", destination: "Calgary", expected: []string{"Mike Tyson", "1984 Summer Olympics", "7-Eleven", "Calgary"}},
			{origin: "Mike Tyson", destination: "Eurocash", expected: []string{"Mike Tyson", "1984 Summer Olympics", "7-Eleven", "Eurocash"}},
			{origin: "Mike Tyson", destination: "Małpka Express", expected: []string{"Mike Tyson", "1984 Summer Olympics", "7-Eleven", "Eurocash", "Małpka Express"}},
			{origin: "Mike Tyson", destination: "Tea", expected: []string{"Mike Tyson", "1984 Summer Olympics", "7-Eleven", "Eurocash", "Tea"}},
			{origin: "Mike Tyson", destination: "Afghanistan", expected: []string{"Mike Tyson", "1984 Summer Olympics", "Afghanistan"}},
		}

		for id, testCase := range testCases {
			crawler := NewForward(test.NewMockWiki())
			go crawler.Run(testCase.origin, testCase.destination)

			// wait for either path result is received or the crawler finishes all its activities
			select {
			case <-crawler.Path():
			case <-crawler.Done():
			}

			for _, title := range testCase.expected {
				if !crawler.visited(title) {
					t.Errorf("Test case %d failed.\nExpected page %q to be included in crawler's visited map: %q", id, title)
				}
			}
		}
	})
}
