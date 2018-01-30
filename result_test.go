package wikiracer

import (
	"fmt"
	"testing"
	"time"

	"github.com/ihcsim/wikiracer/errors"
)

func TestResultString(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		var (
			destination = "Example Docs"
			err         = errors.DestinationUnreachable{Destination: destination}
		)

		actual := fmt.Sprintf("%s", Result{Err: err})
		expected := fmt.Sprintf("%s", err)
		if actual != expected {
			t.Errorf("Mismatch result. Expected %q, but got %q", expected, actual)
		}
	})

	t.Run("Path and duration", func(t *testing.T) {
		var (
			path     = "Mike Tyson -> Alexander the Great -> Greek language -> Fruit Anatomy -> Segment"
			duration = 6 * time.Microsecond
		)

		actual := fmt.Sprintf("%s", Result{Path: []byte(path), Duration: duration})
		expected := fmt.Sprintf("Path: %q, Duration: %s", path, duration)
		if actual != expected {
			t.Errorf("Mismatch result. Expected %q, but got %q", expected, actual)
		}
	})

	t.Run("Path, duration and error", func(t *testing.T) {
		var (
			path        = `"Mike Tyson -> Alexander the Great -> Greek language -> Fruit Anatomy -> Segment"`
			duration    = 6 * time.Microsecond
			destination = "Example Docs"
			err         = errors.DestinationUnreachable{Destination: destination}
		)

		actual := fmt.Sprintf("%s", Result{Path: []byte(path), Duration: duration, Err: err})
		expected := fmt.Sprintf("%s", err)
		if actual != expected {
			t.Errorf("Mismatch result. Expected %q, but got %q", expected, actual)
		}
	})
}
