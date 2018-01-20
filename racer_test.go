package wikiracer

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/crawler"
	"github.com/ihcsim/wikiracer/internal/validator"
	"github.com/ihcsim/wikiracer/internal/wiki"
	"github.com/ihcsim/wikiracer/log"
	"github.com/ihcsim/wikiracer/test"
)

var (
	timeout  = 500 * time.Millisecond
	mockWiki = test.NewMockWiki()
)

func TestFindPath(t *testing.T) {
	log.Instance().SetBackend(log.QuietBackend)

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
				var (
					racer           = New(crawler.NewForward(mockWiki), &validator.InputValidator{mockWiki})
					result          = make(chan *Result)
					ctx, cancelFunc = context.WithTimeout(context.Background(), timeout)
				)
				defer cancelFunc()

				go func() {
					result <- racer.FindPath(ctx, testCase.origin, testCase.destination)
				}()

				select {
				case actual := <-result:
					b, err := ioutil.ReadAll(actual.Path)
					if err != nil {
						t.Fatal("Unexpected error: ", err)
					}
					if testCase.expected != string(b) {
						t.Errorf("Mismatch path. Test case: %d\nExpected: %s\nActual: %s", id, testCase.expected, b)
					}

				case <-ctx.Done():
					t.Fatalf("Test case %d timed out")
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
				var (
					racer           = New(crawler.NewForward(mockWiki), &validator.InputValidator{mockWiki})
					ctx, cancelFunc = context.WithTimeout(context.Background(), timeout)
					result          = make(chan *Result)
				)
				defer cancelFunc()

				go func() {
					result <- racer.FindPath(ctx, testCase.origin, testCase.destination)
				}()

				select {
				case actual := <-result:
					b, err := ioutil.ReadAll(actual.Path)
					if err != nil {
						t.Fatal("Unexpected error: ", err)
					}

					passed := false
					for _, option := range testCase.expected {
						if option == string(b) {
							passed = true
							break
						}
					}

					if !passed {
						t.Errorf("Mismatch path. Test case: %d\nExpected either one of: %v\nActual: %s", id, testCase.expected, b)
					}
				case <-ctx.Done():
					t.Fatalf("Test case %d timed out")
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
			var (
				racer           = New(crawler.NewForward(mockWiki), &validator.InputValidator{mockWiki})
				result          = make(chan *Result)
				ctx, cancelFunc = context.WithTimeout(context.Background(), timeout)
			)
			defer cancelFunc()

			go func() {
				result <- racer.FindPath(ctx, testCase.origin, testCase.destination)
			}()

			select {
			case <-ctx.Done():
				// the page is considered unreachable if racer can't find it before the context timed out
				if actual := <-result; actual.Err != testCase.expected {
					t.Errorf("Error mismatch.\nExpected: %s.\nActual: %s", testCase.expected, actual.Err)
				}
			case actual := <-result:
				// handle other non-timeout errors
				switch cast := actual.Err.(type) {
				case errors.PageNotFound:
					if testCase.expected.Error() != cast.Error() {
						t.Errorf("Mismatch error. Test case: %d\nExpected: %s\nActual: %s", id, testCase.expected, cast)
					}

				case errors.DestinationUnreachable:
					if testCase.expected != cast {
						t.Errorf("Mismatch error. Test case: %d\nExpected: %s\nActual: %s", id, testCase.expected, cast)
					}

				case errors.InvalidEmptyInput:
					if testCase.expected != cast {
						t.Errorf("Mismatch error. Test case: %d\nExpected: %s\nActual: %s", id, testCase.expected, cast)
					}

				default:
					t.Fatalf("Type assertion failed: %T", cast)
				}
			}
		}
	})
}

func TestTimedFindPath(t *testing.T) {
	var (
		racer        = New(crawler.NewForward(mockWiki), &validator.InputValidator{mockWiki})
		ctx          = context.Background()
		origin       = "Mike Tyson"
		destination  = "Segment"
		expectedPath = "Mike Tyson -> Alexander the Great -> Greek language -> Fruit anatomy -> Segment"
	)

	actual := racer.TimedFindPath(ctx, origin, destination)
	if actual.Err != nil {
		t.Fatal("Unexpected error: ", actual.Err)
	}

	actualPath, err := ioutil.ReadAll(actual.Path)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	if string(actualPath) != expectedPath {
		t.Errorf("Mismatch path.\nExpected %q\nActual: %q", expectedPath, actualPath)
	}
}
