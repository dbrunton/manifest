package manifest

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
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

func Compare (m, n Manifest) bool {
	sort.Sort(ByPath{m})
	sort.Sort(ByPath{n})

	if m.String() == n.String() {
		return true
	}
	return false
}

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

func Load(s string) (m Manifest) {
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		if len(fields) > 1 {
			checksum, _ := hex.DecodeString(fields[1])
			path := fields[0]
			me := ManifestEntry{path, checksum}
			m = append(m, me)
		}
	}
	return
}

func Create(dir string) (m Manifest) {

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

	for me := range ch {
		// grab off of a channel
		m = append(m, me)
	}
	return
}
