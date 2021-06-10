package traversal_service

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestSymbolicLink(t *testing.T) {
	files, err := ioutil.ReadDir(`/home/saber/SRC/katmai-1.24/lyrics_osd`)

	if err != nil {
		panic(err.Error())
	}

	for _, file := range files {
		if file.Mode().IsRegular() {

			fmt.Println(file.Name(), " ---")
		} else if file.Mode().IsDir() {
			fmt.Println(file.Name(), "<dir>")
		} else {
			fmt.Println(file.Name(), "???")
		}
	}
}
