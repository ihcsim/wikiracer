package wikiracer

import "testing"

func TestFindPath(t *testing.T) {
	wiki = NewMockWiki()

	expected := "Mike Tyson -> Alexander the Great -> Greek language -> Fruit anatomy -> Segment"
	if actual := FindPath("Mike Tyson", "Segment"); expected != actual {
		t.Errorf("Mismatch path.\nExpected: %s\n.Actual: %s", expected, actual)
	}
}
