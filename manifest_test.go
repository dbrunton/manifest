package manifest

import(
	"testing"
)


func TestManifest(t *testing.T) {
	correct := "data/bar.txt\td41d8cd98f00b204e9800998ecf8427e\ndata/foo.txt\td41d8cd98f00b204e9800998ecf8427e\n"
	result := Manifest("data")
	if result != correct {
		t.Errorf("Expected: '%s'\nGot: '%s'\n", correct, result)
	}
}
