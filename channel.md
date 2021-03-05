---
title: Channel的使用
date: 2021-1-25
---



# Channel  

Channel是Go核心的数据结构和Goroutine之间的通信方式，Channel是支撑Go语言高性能并发编程的重要结构。  

## 设计原理  
Go中经常被人提及的设计模式就是：不要通过共享内存的方式进行通信，而是应该通过通信的方式共享内存。在很多主流编程语言中，多个线程传递数据的方式一般都是共享内存，为了解决线程竞争，我们需要限制同一时间能够读写这些变量的线程数目，然而与Go语言鼓励的设计并不相同。  

两个独立运行的goroutine可以向Channel中发送数据，另一个会从channel中接受数据，能通过Channel间接完成通信  

如果说goroutine是Go语言程序的并发体的话，那么channels则是它们之间的通信机制。一个channel是一个通信机制，他可以让一个goroutine通过它给另一个goroutine发送值信息。每个channel都有一个特殊的类型，也就是channels可发送数据的类型。  

与map类似，channel也对应一个make创建的底层数据结构的引用。当我们复制一个channel或者用于函数参数传递时，我们只是拷贝了一个channel的引用，因此调用者和被调用者将引用同一个channel对象。和其它引用类型一样，channel的零值也是nil  


两个相同的类型的channel可以使用==运算符比较。如果两个channel引用的是相同的对象，那么比较的结果为真。一个channel也可以和nil进行比较。  

一个channel有发送和接受两个主要操作，都是通信的行为。一个发送语句将一个值从一个goroutine通过channel发送到另一个执行接收操作的goroutine。发送和接收两个操作符都使用<-运算符。在发送语句中，<-运算符分割channel和要发送的值，在接收语句中，<-运算符写在channel对象之前。一个不适用接收结果的接收操作也是合法的。


```go
ch <- x   // a send statement
x = <- ch // a receive expression in an assignment statement
<-ch      // a receive statement; result is discarded
```

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


## Channel的种类  

channel分为无缓冲channel和有缓冲channel。两者区别如下：  


+ 无缓冲: 发送与接受操作是同时发送，如果没有goroutine读取channel(<-channel)，则发送者(channel<-)会一直阻塞  


+ 缓冲: 缓冲channel类似一个有容量的队列。当队列满的时候发送者会阻塞；当队列空的时候，接收者会阻塞  


## channel的关闭  

```go
ch := make(chan int)
```

关于关闭channel有几点需要注意的是:  

+ 重复关闭channel会导致panic  
+ 向已经关闭的channel发送数据会panic  
+ 从关闭的channel读取数据不会panic，读取channel中已有的数据之后再读就是channel类似的默认值，比如chan int类型的channel关闭之后读取的值为0  

可以用map中类似的方式去判断channel是否关闭  

```go
ch := make(chan int, 10)
close(ch)

val, ok := <-ch
if ok == false {
	//channel close
}
```


## Select和Channel  

go语言中的select可以让goroutine同时等待多个Channel可读或者可写，在多个文件或者Channel状态改变之前，select会一直阻塞当前线程或者goroutine  


select是与switch相似的控制结构，与switch不同的是，select中虽然也有多个case，但是这些case中的表达式必须都是Channel的收发操作。下面的代码就展示了一个包含Channel收发操作的select结构  


```go
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	c := make(chan int)
	quit := make(chan int, 2)
	go func() {
		for i := 1; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}
```
上述控制结构会等待c <- x或者<-quit两个表达式中任意一个返回。无论哪一个表达式返回都会立刻执行case中的代码，当select中的两个case同时被触发时，会随机执行其中一个  

## select中会遇到的现象  

1. select能够在Channel上进行非阻塞的收发操作  
2. select在遇到多个Channel同时响应时，会随机执行一种情况

example:  

```go
func main() {
	ch := make(chan int)
	select {
	case i := <-ch:
		println(i)
	default:
		println("default")
	}
}
```

输出
```terminal
default
```

运行上述代码的时候就不会阻塞当前的Goroutine，它会直接执行default中的代码。  
只要稍微想一下就知道，这样设计很合理，select的作用是同时监听多个case是否可以执行，如果多个Channel都不是很彳亍，那么就默认运行default  
兄弟们把合理打到公屏上  

非阻塞的Channel发送和接收操作还是很有不要的，很多场景下我们不希望Channel操作阻塞当前的Goroutine，只是想看看Channel的可读或者可写状态，如下代码所示:  

```go
errCh := make(chan error, len(tasks))
wg := sync.WaitGroup{}
wg.Add(len(tasks))
for i := range tasks {
    go func() {
        defer wg.Done()
        if err := tasks[i].Run(); err != nil {
            errCh <- err
        }
    }()
}
wg.Wait()

select {
case err := <-errCh:
    return err
default:
    return nil
}
```


在上面这段代码中，我们不关心到底多少个任务执行失败了，只关心是否存在返回错误的任务，最后的select语句能够很好地完成这个任务


## 随机执行  

另一个使用select遇到的情况是同时有多个case就绪时，select会选择哪个case执行的问题，我们通过下面的代码可以简单了解一下：  

```go
func main() {
	ch := make(chan int)
	go func() {
		for range time.Tick(1 * time.Second) {
			ch <- 0
		}
	}()

	for {
		select {
		case <-ch:
			println("case1")
		case <-ch:
			println("case2")
		}
	}
}
```

```terminal
case1
case2
case2
case2
case1
case1
case2
case2
case2
```

可以看到select在遇到多个<-ch同时满足可读或者可写条件时会随机选择一个case执行其中的代码  

其实是伪随机。如果没有case就会选择default去处理，如果default case存在的话，如果不存在则select会阻塞，知道某个case需要处理  

在上面的代码中，两个case都是同时满足执行条件的，如果我们按照顺序依次判断，那么后面的条件永远都会得不到执行，而随机的引入就是为了避免饥饿问题的发生  



## Range  

for ... range 语句可以处理Channel  

some examples:  
```go
func main() {
	go func() {
		time.Sleep(1 * time.Second)
	}()
	c := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			c <- i
		}
		close(c)
	}()

	for i := range c {
		fmt.Println("Finished ", i)
	}
}
```

range c 产生的迭代值为Channel中发送的值，它会一直迭代知道channel被关闭。如果把上述的**close(c)**注释掉，程序会发生阻塞  


## timeout  

select有很重要的一个应用就是超时处理，如果没有case的话select就会一直阻塞，在应用中肯定不被允许。这时候我们可能就需要一个超时操作用来处理超时的情况。  

如下述例子:  

```go
func main() {
	c1 := make(chan string, 1)
	go func() {
		for {
			time.Sleep(time.Second * 2)
			c1 <- "result1"
		}
	}()

	for {
		select {
		case res := <-c1:
			fmt.Println(res)
		case <-time.After(time.Second * 1):
			fmt.Println("timeout 1")
		}
	}
}
```

## Timer和Ticker

接下来看看关于时间的两个Channel。
timer是一个定时器，代表未来的一个单一事情，你可以告诉timer你要等待多长时间，它提供一个Channel，在将来的那个时间那个Channel提供了一个时间值。下面的例子中第二行会阻塞两秒钟的时间，然后才会执行下面的代码


```go
func main() {
	timer1 := time.NewTimer(time.Second * 2)
	<-timer1.C
	fmt.Println("Timer 1 expired")
}
```

当然，还可以使用time.Sleep实现，time.Stop也可以  

而Ticker是一个定时触发的计时器，它会以一个interval往Channel发送一个事件(当前时间)，而Channel的接收者可以以固定的时间间隔从Channel中读取事件。

```go
func main() {
	ticker := time.NewTicker(time.Millisecond * 500)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
		wg.Done()
	}()
	wg.Wait()
}
```
