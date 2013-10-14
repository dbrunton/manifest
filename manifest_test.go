package manifest

import(
	"testing"
)


func TestMakeManifest(t *testing.T) {
	correct := "data/bar.txt\td41d8cd98f00b204e9800998ecf8427e\ndata/foo.txt\td41d8cd98f00b204e9800998ecf8427e\n"
	result := MakeManifest("data")
	if result.String() != correct {
		t.Errorf("Expected: '%s'\nGot: '%s'\n", correct, result)
	}
}
