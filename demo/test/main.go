package main

import (
	"fmt"
	"time"
)

var numChan chan uint64

func init() {
	numChan = make(chan uint64)
}

func main() {
	go listen()
	go listen()
	go listen()

	var i uint64
	for {
		numChan <- i
		i++
		time.Sleep(time.Second)
	}
}

func listen() {
	for {
		select {
		case num := <-numChan:
			fmt.Println(num)
		}
	}
}
