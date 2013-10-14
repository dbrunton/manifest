package manifest

import (
	"bytes"
	"io/ioutil"
	"fmt"
	"crypto/md5"
	"hash"
	"path/filepath"
	"os"
)

type ManifestEntry struct {
	path string
	checksum []byte
}

type Manifest []ManifestEntry

func (m Manifest) String() (s string) {
	buffer := bytes.NewBufferString("")
	for _, entry := range m {
		fmt.Fprintf(buffer, "%s\t%x\n", entry.path, entry.checksum)
	}
	return string(buffer.Bytes())
}

func listFiles(dir string) (files []string) {
	err := filepath.Walk("./data",
			     func(path string, f os.FileInfo, err error) error {
				     if f.IsDir() {
					     return nil
				     }
				     files = append(files, filepath.Clean(path));
				     return nil;
			     })
	if err != nil {
		panic(err)
	}
	return
}

func checksum(file string) ManifestEntry {
	contents, _ := ioutil.ReadFile(file)
	var h hash.Hash = md5.New()
	var b []byte
	h.Write([]byte(contents))
	return ManifestEntry{ file, h.Sum(b) }
}

func MakeManifest(dir string) (manifest Manifest) {

	for _, path := range listFiles(dir) {
		manifest = append(manifest, checksum(path))
	}
	return
}

