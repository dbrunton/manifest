package manifest
/* package manifest */

/* create a type to have the path (e.g. Name) and checksum, return that,
   print it out at the end.  Then do it all over in goroutines.  */

import (
	"io/ioutil"
	"fmt"
	"crypto/md5"
	"hash"
	"bytes"
)

func Manifest(dir string) string {

	buffer := bytes.NewBufferString("")

	directoryContents, _ := ioutil.ReadDir(dir)

	for _, entry := range directoryContents {

		if entry.IsDir() {
			Manifest(fmt.Sprintf("%s/%s", dir, entry.Name()))

		} else {
			contents, _ := ioutil.ReadFile(entry.Name())
			var h hash.Hash = md5.New()
			var b []byte
			h.Write([]byte(contents))
			fmt.Fprintf(buffer, "%s/%s\t%x\n", dir, entry.Name(), h.Sum(b))
		}
	}

	return string(buffer.Bytes())
}

