---
title: Go语言处理并发
date: 2020-10-29
---

# 并发的概念

Erlang之父Armstrong用一张很经典的图来描述并发与并行  
![concurrency](concurrency.jpg)  
与并发相对的概念并不是并行，而是顺序  
顺序执行，则必须有先有后，而并发执行，强调的是"同时出发"，无论另一个任务是否完成都可以进行  

而并行相对的概念是串行，类似于电路中的串联跟并联  
* 串行：有一个任务执行单元，从物理上就只能一个任务一个任务地执行
* 并行：有多个任务执行单元，从物理上就可以多个任务一起执行  

并发与并行并不是互斥的，前者着重于任务的调度，后者着重于任务实际执行情况  
  
* 单核CPU多任务：并发(不必等上个任务完成就可以开始下一个任务)、串行(只有一个CPU核心)

* 多线程：  

  + 并发、串行(所有线程都在同一个核上执行)
  + 并发、并行(不同线程在不同的核上执行)

---

# Go语言实现并发

Golang是通过协程的方式实现并发，goroutine与Java中的Thread不同，goroutine的创建与销毁都不需要太多的资源消耗，Golang中的协程比线程更加“轻量级”  
Golang内部实现了goroutine之间的内存共享，与此同时，执行goroutine只需要4~5kb的内存空间，正因为轻量的特性，Golang很容易实现成千上万级别的并发任务。    

---

goroutine的关键字：go

```go
go func() {

}
```

---
**Concurrency**  
```go
package main

import (
    "fmt"
    "runtime"
)

func say(s string) {
    for i := 0; i < 5; i++ {
        runtime.Gosched()
        fmt.Println(s)
    }
}

func main() {
    go say("world")
    say("hello")
}
```
**Output**
```shell
hello
world
hello
world
hello
world
world
world
hello
hello
```

当编译器遇到go关键字时，会跳到下一行同时执行语句