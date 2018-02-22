package wikipedia

import (
	"reflect"
	"testing"

	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/wiki"
	"github.com/sadbox/mediawiki"
)

func TestFindPage(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	client.api = mockAPI

	t.Run("One Batch", func(t *testing.T) {
		var (
			title     = "Mike Tyson"
			nextBatch = ""
		)
		actual, err := client.FindPage(title, nextBatch)
		if err != nil {
			t.Fatal(err)
		}

		expected := &wiki.Page{
			ID:        39027,
			Title:     title,
			Namespace: 0,
			Links:     []string{"1984 Summer Olympics", "20/20 (US television show)", "Aaron Pryor", "Abdullah the Butcher"},
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Mismatch page. Expected %+v\nActual %+v\n", expected, actual)
		}
	})

	t.Run("Multiple Batches", func(t *testing.T) {
		var (
			title     = "Alexander the Great"
			nextBatch = ""
		)
		actual, err := client.FindPage(title, nextBatch)
		if err != nil {
			t.Fatal(err)
		}

		expected := &wiki.Page{
			ID:        783,
			Title:     title,
			Namespace: 0,
			Links:     []string{"Apepi", "Aahotepre", "Abbasid Caliphate", "Abdalonymus", "Dutch Empire", "Dynamis (Bosporan queen)", "Dynasty", "Early Dynastic Period (Egypt)", "Menandar", "Menes", "Mental health", "Mentuhotep I"},
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Mismatch page. Expected %+v\nActual %+v\n", expected, actual)
		}
	})

	t.Run("Missing Page", func(t *testing.T) {
		var (
			title     = "Missing Page"
			nextBatch = ""
		)
		_, actual := client.FindPage(title, nextBatch)

		expected := errors.PageNotFound{wiki.Page{Title: title}}
		if expected.Error() != actual.Error() {
			t.Fatal(err)
		}
	})
}

func mockAPI(api *mediawiki.MWApi, values ...map[string]string) ([]byte, error) {
	var json []byte
	switch values[0]["titles"] {
	case "Mike Tyson":
		json = []byte(`
{
	"batchcomplete": true,
  "query": {
    "pages": [
      {
        "pageid": 39027,
        "ns": 0,
        "title": "Mike Tyson",
        "links": [
          {"ns": 0, "title": "1984 Summer Olympics"},
          {"ns": 0, "title": "20\/20 (US television show)"},
          {"ns": 0, "title": "Aaron Pryor"},
          {"ns": 0, "title": "Abdullah the Butcher"}
				]
      }
    ]
  },
  "limits": {
    "links": 500
  }
}`)

	case "Alexander the Great":
		switch values[0]["plcontinue"] {
		case "783|0|Dutch_Empire":
			json = []byte(`
{
  "continue": {
    "plcontinue": "783|0|Menander",
    "continue": "||"
  },
  "query": {
    "pages": [
      {
        "pageid": 783,
        "ns": 0,
        "title": "Alexander the Great",
        "links": [
          {"ns": 0, "title": "Dutch Empire"},
          {"ns": 0, "title": "Dynamis (Bosporan queen)"},
          {"ns": 0, "title": "Dynasty"},
          {"ns": 0, "title": "Early Dynastic Period (Egypt)"}
				]
      }
    ]
  },
  "limits": {
    "links": 500
  }
}`)

		case "783|0|Menander":
			json = []byte(`
{
	"batchcomplete": true,
  "query": {
    "pages": [
      {
        "pageid": 783,
        "ns": 0,
        "title": "Alexander the Great",
        "links": [
          {"ns": 0, "title": "Menandar"},
          {"ns": 0, "title": "Menes"},
          {"ns": 0, "title": "Mental health"},
          {"ns": 0, "title": "Mentuhotep I"}
				]
      }
    ]
  },
  "limits": {
    "links": 500
  }
}`)

		default:
			json = []byte(`
{
  "continue": {
    "plcontinue": "783|0|Dutch_Empire",
    "continue": "||"
  },
  "query": {
    "pages": [
      {
        "pageid": 783,
        "ns": 0,
        "title": "Alexander the Great",
        "links": [
          {"ns": 0, "title": "Apepi"},
          {"ns": 0, "title": "Aahotepre"},
          {"ns": 0, "title": "Abbasid Caliphate"},
          {"ns": 0, "title": "Abdalonymus"}
        ]
      }
    ]
  },
  "limits": {
    "links": 500
  }
}`)
		}

	case "Missing Page":
		json = []byte(`
{
  "batchcomplete": true,
  "query": {
    "normalized": [{"fromencoded": false, "from": "llll", "to": "Llll"}],
    "pages": [{"ns": 0, "title": "Llll", "missing": true}]
  },
  "limits": {
    "links": 500
  }
}`)
	}

	return json, nil
}
