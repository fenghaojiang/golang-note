---
title: Channel的使用
date: 2021-1-25
---



# Channel  

Channel是Go核心的数据结构和Goroutine之间的通信方式，Channel是支撑Go语言高性能并发编程的重要结构。  

## 设计原理  
Go中经常被人提及的设计模式就是：不要通过共享内存的方式进行通信，而是应该通过通信的方式共享内存。在很多主流编程语言中，多个线程传递数据的方式一般都是共享内存，为了解决线程竞争，我们需要限制同一时间能够读写这些变量的线程数目，然而与Go语言鼓励的设计并不相同。  

两个独立运行的goroutine可以向Channel中发送数据，另一个会从channel中接受数据，能通过Channel间接完成通信  

## FIFO  

目前的Channel收发操作都遵守了先进先出的设计，具体规则如下：  
+ 先从Channel读取数据的Goroutine会先接收到数据  
+ 先向Channel发送数据的Goroutine会得到先发送数据的权利  


usage:  

```go 
package main

import "fmt"

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

func main() {
	s := []int{1, 2, 3, 4, 5, 6}
	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c
	fmt.Println(x, y, x+y)
}
```

output:  

```terminal
6 15 21
```

