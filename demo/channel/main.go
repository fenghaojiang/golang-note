package main

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
	//send()
	//receive()
	cl()
}
