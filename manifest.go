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
	"sort"
	"sync"
)

type ManifestEntry struct {
	path     string
	checksum []byte
}

type Manifest []ManifestEntry

func (m Manifest) Len() int { return len(m) }

func (m Manifest) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

type ByPath struct{ Manifest }

func (s ByPath) Less(i, j int) bool { return s.Manifest[i].path < s.Manifest[j].path }

func (m Manifest) String() (s string) {
	sort.Sort(ByPath{m})
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

func Create(dir string) (manifest Manifest) {

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
