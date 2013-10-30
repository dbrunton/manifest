package manifest

import(
	"testing"
)

func TestCompare(t *testing.T) {
	content1:= "data/foo.txt\td41d8cd98f00b204e9800998ecf8427e\ndata/bar.txt\td41d8cd98f00b204e9800998ecf8427e\n"
	content2:= "data/bar.txt\td41d8cd98f00b204e9800998ecf8427e\ndata/foo.txt\td41d8cd98f00b204e9800998ecf8427e\n"
	m := Load(content1)
	n := Load(content2)
	if Compare(m, n) != true {
		t.Errorf("Unsorted manifests don't match.")
	}
}

func TestLoad(t *testing.T) {
	content := "data/bar.txt\td41d8cd98f00b204e9800998ecf8427e\ndata/foo.txt\td41d8cd98f00b204e9800998ecf8427e\n"
	m := Load(content)
	if m.String() != content {
		t.Errorf("Loaded manifest is not correct.")
	}
}

func TestCreate(t *testing.T) {
	correct := "data/bar.txt\td41d8cd98f00b204e9800998ecf8427e\ndata/foo.txt\td41d8cd98f00b204e9800998ecf8427e\n"
	result := Create("data")
	if result.String() != correct {
		t.Errorf("Expected: '%s'\nGot: '%s'\n", correct, result)
	}
}
