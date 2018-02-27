package wikipedia

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/wiki"
	"github.com/sadbox/mediawiki"
)

const (
	invalidAction = "invalidAction"
	invalidParam  = "invalidParam"
)

func TestFindPage(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		client, err := NewClient()
		if err != nil {
			t.Fatal(err)
		}
		client.api = mockAPI

		t.Run("Single Batch Result", func(t *testing.T) {
			t.Run("One Title", func(t *testing.T) {
				var (
					title     = "Mike Tyson"
					nextBatch = ""
				)
				actual, err := client.FindPages(title, nextBatch)
				if err != nil {
					t.Fatal(err)
				}

				expected := &wiki.Page{
					ID:        39027,
					Title:     title,
					Namespace: 0,
					Links:     []string{"1984 Summer Olympics", "20/20 (US television show)", "Aaron Pryor", "Abdullah the Butcher"},
				}

				if !reflect.DeepEqual(actual[0], expected) {
					t.Errorf("Mismatch page.\nExpected %+v\nActual %+v\n", expected, actual[0])
				}
			})

			t.Run("Multiple Titles", func(t *testing.T) {
				var (
					firstTitle  = "Mike Tyson"
					secondTitle = "Alexander the Great"
					titles      = firstTitle + "|" + secondTitle
					nextBatch   = ""
				)
				actual, err := client.FindPages(titles, nextBatch)
				if err != nil {
					t.Fatal(err)
				}

				expected := []*wiki.Page{
					&wiki.Page{
						ID:        39027,
						Title:     firstTitle,
						Namespace: 0,
						Links:     []string{"1984 Summer Olympics", "20/20 (US television show)", "Aaron Pryor", "Abdullah the Butcher"},
					},
					&wiki.Page{
						ID:        783,
						Title:     secondTitle,
						Namespace: 0,
						Links:     []string{"Dutch Empire", "Dynamis (Bosporan queen)", "Dynasty", "Early Dynastic Period (Egypt)"},
					},
				}

				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("Mismatch page.\nExpected %s\n  Actual %s\n", expected, actual)
				}
			})
		})

		t.Run("Multi-Batch Result", func(t *testing.T) {
			t.Run("One Title", func(t *testing.T) {
				var (
					title     = "Alexander the Great"
					nextBatch = ""
				)
				actual, err := client.FindPages(title, nextBatch)
				if err != nil {
					t.Fatal(err)
				}

				expected := &wiki.Page{
					ID:        783,
					Title:     title,
					Namespace: 0,
					Links:     []string{"Apepi", "Aahotepre", "Abbasid Caliphate", "Abdalonymus", "Dutch Empire", "Dynamis (Bosporan queen)", "Dynasty", "Early Dynastic Period (Egypt)", "Menandar", "Menes", "Mental health", "Mentuhotep I"},
				}

				if !reflect.DeepEqual(actual[0], expected) {
					t.Errorf("Mismatch page.\nExpected %+v\nActual %+v\n", expected, actual[0])
				}
			})

			t.Run("Multiple Titles", func(t *testing.T) {
				var (
					firstTitle  = "Mike Tyson"
					secondTitle = "Vancouver"
					titles      = firstTitle + "|" + secondTitle
					nextBatch   = ""
				)
				actual, err := client.FindPages(titles, nextBatch)
				if err != nil {
					t.Fatal(err)
				}

				expected := []*wiki.Page{
					&wiki.Page{
						ID:        32706,
						Title:     secondTitle,
						Namespace: 0,
						Links:     []string{"1918 Vancouver general strike", "1924 Winter Olympics", "1928 Winter Olympics", "1930 British Empire Games", "Hudson's Bay Company", "Humidex", "Hydroponics", "Ice hockey"},
					},
					&wiki.Page{
						ID:        39027,
						Title:     firstTitle,
						Namespace: 0,
						Links:     []string{"1984 Summer Olympics", "20/20 (US television show)", "Aaron Pryor", "Abdullah the Butcher"},
					},
				}

				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("Mismatch page.\nExpected %s\n  Actual %s\n", expected, actual)
				}
			})
		})

		t.Run("Missing Page", func(t *testing.T) {
			var (
				title     = "Missing Page"
				nextBatch = ""
			)
			_, actual := client.FindPages(title, nextBatch)

			expected := errors.PageNotFound{wiki.Page{Title: title}}
			if expected.Error() != actual.Error() {
				t.Errorf("Mismatch result.\nExpected error: %v\nActual error: %v", expected, actual)
			}
		})
	})

	t.Run("Error", func(t *testing.T) {
		client, err := NewClient()
		if err != nil {
			t.Fatal(err)
		}

		client.api = mockAPIError

		t.Run("Invalid action", func(t *testing.T) {
			expected := &serverError{
				msg: fmt.Sprintf("Unrecognized value for parameter \"action\": %s.\n", invalidAction),
			}
			_, actual := client.FindPages(invalidAction, "")
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("Expected error didn't occur. Got %q", actual)
			}
		})

		t.Run("Invalid param", func(t *testing.T) {
			expected := &serverError{
				msg: fmt.Sprintf("Unrecognized parameter: %s.Unrecognized value for parameter \"list\": raremodule.", invalidParam),
			}
			_, actual := client.FindPages(invalidParam, "")
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("Expected error didn't occur. Got %q", actual)
			}
		})
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

	case "Mike Tyson|Alexander the Great":
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
			},
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

	case "Mike Tyson|Vancouver":
		switch values[0]["plcontinue"] {
		case "32706|0|Hudson's_Bay_Company":
			json = []byte(`
{
  "continue": {
    "plcontinue": "32706|0|The_Killing_(U.S._TV_series)",
    "continue": "||"
  },
  "query": {
    "pages": [
      {
        "pageid": 32706,
        "ns": 0,
        "title": "Vancouver",
        "links": [
          {"ns": 0, "title": "Hudson's Bay Company"},
          {"ns": 0, "title": "Humidex"},
          {"ns": 0, "title": "Hydroponics"},
          {"ns": 0, "title": "Ice hockey"}
        ]
      },
      {
        "pageid": 39027,
        "ns": 0,
        "title": "Mike Tyson"
      }
    ]
  },
  "limits": {
    "links": 500
  }
}
`)

		case "32706|0|The_Killing_(U.S._TV_series)":
			json = []byte(`
{
  "batchcomplete": true,
  "query": {
    "pages": [
      {
        "pageid": 32706,
        "ns": 0,
        "title": "Vancouver"
      },
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

		default:
			json = []byte(`
{
  "continue": {
    "plcontinue": "32706|0|Hudson's_Bay_Company",
    "continue": "||"
  },
  "query": {
    "pages": [
      {
        "pageid": 32706,
        "ns": 0,
        "title": "Vancouver",
        "links": [
          {"ns": 0, "title": "1918 Vancouver general strike"},
          {"ns": 0, "title": "1924 Winter Olympics"},
          {"ns": 0, "title": "1928 Winter Olympics"},
          {"ns": 0, "title": "1930 British Empire Games"}
        ]
      },
      {
        "pageid": 39027,
        "ns": 0,
        "title": "Mike Tyson"
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
    "normalized": [{"fromencoded": false, "from": "Missing Page", "to": "Missing Page"}],
    "pages": [{"ns": 0, "title": "Missing Page", "missing": true}]
  },
  "limits": {
    "links": 500
  }
}`)
	}

	return json, nil
}

func mockAPIError(api *mediawiki.MWApi, values ...map[string]string) ([]byte, error) {
	var json []byte

	switch values[0]["titles"] {
	case invalidAction:
		json = []byte(fmt.Sprintf(`
{
  "errors": [
    {
      "code": "unknown_action",
      "text": "Unrecognized value for parameter \"action\": %s.",
      "module": "main"
    }
  ],
  "docref": "See https://en.wikipedia.org/w/api.php for API usage. Subscribe to the mediawiki-api-announce mailing list at &lt;https://lists.wikimedia.org/mailman/listinfo/mediawiki-api-announce&gt; for notice of API deprecations and breaking changes.",
  "servedby": "mw1282"
}`, invalidAction))

	case invalidParam:
		json = []byte(fmt.Sprintf(`
{
  "batchcomplete": true,
  "warnings": {
    "main": {
      "warnings": "Unrecognized parameter: %s."
    },
    "query": {
      "warnings": "Unrecognized value for parameter \"list\": raremodule."
    }
  }
}`, invalidParam))
	}

	return json, nil
}
