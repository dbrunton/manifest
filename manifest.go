package manifest

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"hash"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"sync"
)

type ManifestEntry struct {
	path     string
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

func listFiles(dir string) (files chan string) {
	files = make(chan string, 65536)

	go func() {
		err := filepath.Walk("./data",
			func(path string, f os.FileInfo, err error) error {
				if f.IsDir() {
					return nil
				}
				files <- filepath.Clean(path)
				return nil
			})
		if err != nil {
			panic(err)
		}
		close(files)
	}()
	return
}

func checksum(file string) ManifestEntry {
	contents, _ := ioutil.ReadFile(file)
	var h hash.Hash = md5.New()
	var b []byte
	h.Write([]byte(contents))
	return ManifestEntry{file, h.Sum(b)}
}

func MakeManifest(dir string) (manifest Manifest) {

	runtime.GOMAXPROCS(runtime.NumCPU())

	files := listFiles(dir)

	ch := make(chan ManifestEntry, 65536)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range files {
				ch <- checksum(path)
			}
		}()
	}
	wg.Wait()
	close(ch)

	for e := range ch {
		// grab off of a channel
		manifest = append(manifest, e)
	}
	return
}
