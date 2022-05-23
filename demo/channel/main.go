package main

import (
	"fmt"
	"time"
)

func send() {
	var c chan struct{}
	c <- struct{}{}
}

func receive() {
	var c chan struct{}
	<-c
}

func cl() {
	var c chan struct{}
	close(c)
}

func control() {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminate := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(terminate)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminate
	}

	done := make(chan interface{})
	terminate := doWork(done, nil)
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("canceling dowork goroutine...")
		close(done)
	}()
	<-terminate
	fmt.Println("done")

}

func main1() {
	// runtime.GOMAXPROCS(runtime.NumCPU())
	//send()
	//receive()
	// cl()
	changeOwner := func() <-chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("received: %d\n", result)
		}
		fmt.Println("done")
	}

	results := changeOwner()
	consumer(results)

}

func main() {
	control()
}
