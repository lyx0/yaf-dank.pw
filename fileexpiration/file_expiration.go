package fileexpiration

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

var (
	deleteIgnoreRegexp = regexp.MustCompile("index\\.html|favicon\\.ico")
)

func DeleteExpired(fd string, maxAge time.Duration) {
	files, err := ioutil.ReadDir(fd)
	if err != nil {
		return
	}

	for _, file := range files {
		fname := file.Name()

		if file.IsDir() || deleteIgnoreRegexp.MatchString(fname) {
			continue
		}

		if time.Since(file.ModTime()) > maxAge {
			err := os.Remove(fd + fname)

			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("Removed %s \n", fname)
		}
	}
}
