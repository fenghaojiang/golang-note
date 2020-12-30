---
title: GPM模型
date: 2020-12-28 
---

# GPM的含义

1. M:操作系统的主线程(物理线程)，操作器，用于将一个G搬到线程上去执行Machine
2. P:协程执行需要的上下文，一个装满G的队列Processor
3. G:协程, goroutine我需要分担出去的任务


## Goroutine

Goroutine就是代码中使用go关键字创建的执行单位，也就是协程，上下文切换不需要经过内核态，加上协程所占用的内存空间极小，所以有着非常大的发展潜力。  


在Go语言中，Goroutine由一个名为runtime.go的结构体表示，该结构体非常复杂，有40多个成员变量，主要用于存储执行栈、状态、当前占用的线程、调度相关的数据。  


```go
type g struct {
    stack struct {
        lo uintptr
        hi uintptr
    }
    stackgurad0 uintptr
    stackguard1 uintptr
    _panic       *_panic
	_defer       *_defer
	m            *m				// 当前的 m
	sched        gobuf
	stktopsp     uintptr		// 期望 sp 位于栈顶，用于回溯检查
	param        unsafe.Pointer // wakeup 唤醒时候传递的参数
	atomicstatus uint32
	goid         int64
	preempt      bool       	// 抢占信号，stackguard0 = stackpreempt 的副本
	timer        *timer         // 为 time.Sleep 缓存的计时器
}
```


## Machine

M就是对应操作系统的线程，最多会有GOMAXPROCS个活跃线程能够正常运行，默认情况下GOMAXPROCS被设置为内核数，假如有四个内核，那么默认就创建四个线程，每一个线程对应一个runtime.m结构体。线程数等于CPU个数的原因是，每个线程分配到一个CPU上就不至于出现线程的上下文切换，可以保证系统开销降到最低。  

```go
type m struct {
    g0 *g
    curg *g
    ...

    ...
    p puintptr
    nextp puintptr
    oldp puintptr
}
```

M里面存储了两个比较重要的东西，一个是g0，另一个是curg  

+ g0: 会深度参与运行时的调度过程，比如goroutine的创建、内存分配等，0号选手懂的都懂
+ curg(current goroutine): 代表当前正在线程上执行的goroutine  

M中还要存储与P相关的数据

+ p: 正在运行代码的处理器
+ nextp: 暂存的处理器
+ old: 系统调用之前的线程处理器


## Processor
Processor负责Machine与Goroutine的连接，它能提供线程需要的上下文环境，也能分配G到它应该去的线程上执行，有了它，每个G都能得到合理的调用。
同样，处理器的数量默认也是按照GOMAXPROCS来设置的，与线程的数量一一对应


```go
type p struct {
    m muintptr

    runqhead uint32
    runqtail uint32
    runq [256]guintptr
    runnext guintptr
    ...
}
```

结构体P中存储了性能追踪、垃圾回收、计时器等相关的字段外，还存储了处理器的待运行队列，队列中存储的是待执行的Goroutine列表


## 三者的关系

首先，默认启动四个线程四个处理器，然后互相绑定  
这个时候，一个Goroutine结构体被创建，在进行函数体地址、参数起始地址、参数长度等信息以及调度相关属性更新之后，它就要进到一个处理器的队列等待发车。  

当你又创建一个G时就轮流往P中放，假如有很多G，都塞满了可以把G塞到全局队列里(候车大厅)  

M首先会从私有队列中取G执行，如果取完的话就去全局队列中取，如果全局队列也没有的话，就去其他处理器队列里偷

如果都没找到要执行的G，M就会与P断开关系，然后去睡觉(idle)  

如果两个Goroutine正在通过Channel管道传输阻塞住了M会去寻找另外的G去执行  

## 系统调用

如果G进行了系统调用syscall，M也会跟着进入系统调用状态，那么这个P留在这里就浪费了，P不会傻傻的等待G和M系统调用完成，而会去找其他比较闲的M执行其他的G  

当G完成了系统调用，因为要往下执行，所以必须再找一个空闲的处理器发车  

如果没有空闲的处理器了，那就只能把G放回全局队列中等待分配   

sysmom是我们的保洁阿姨，它是一个M，不需要P就可以独立运行，没20μs~10ms会被唤醒一次出来打扫卫生，主要工作就是回收垃圾，回收长时间系统调度阻塞的P、向长时间运行的G发出抢占调度等等














