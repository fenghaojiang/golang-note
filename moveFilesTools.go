package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	suffix := ".md"
	fileDir := "./"
	outputDir := "./output/"
	files, _ := ioutil.ReadDir(fileDir)
	var cnt int32 = 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffix) {
			fmt.Println(file.Name())
			fileObj, err := os.Open(fileDir + file.Name())
			defer fileObj.Close()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			content, err := ioutil.ReadAll(fileObj)
			if err != nil {
				fmt.Println(err.Error())
			}
			result := modify(string(content))
			newFile, err := os.OpenFile(outputDir+file.Name(), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			newFile.WriteString(result)
			cnt++
		}
	}
	fmt.Println("Modify Files Number == ", cnt)
}

func modify(s string) string { //add string "./img/" after "![]("
	subs := "]("
	result := ""
	article := strings.SplitAfter(s, subs)
	for i, v := range article {
		if i%2 == 0 { //start from zero
			result += v
		} else {
			result = (result + "./img/") + v
		}
	}
	// fmt.Print(result)
	return result
}
