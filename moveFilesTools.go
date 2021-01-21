package main

import (
	"io/ioutil"
	"strings"
)

func main() {
	suffix := ".md"
	fileDir := "./"
	files, _ := ioutil.ReadDir(fileDir)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffix) {

		}
	}
}
