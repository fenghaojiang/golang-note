package main

import "fmt"

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

func main() {
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
