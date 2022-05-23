---
title: Concurrency In Go
date: 2022-5-21
---  



link: https://geektutu.com/post/hpg-sync-cond.html  //写法有待加强

代码在demo目录  

## sync.Cond  

Signal 提供同志goroutine阻塞调用Wait，另一种叫Broadcast，运行时内部维护一个FIFO列表，等待接受信号。
Signal发现等待最长时间的goroutine并通知它，而Broadcast向所有的goroutine发送信号。  
Broadcast提供了一种同时与多个goroutine通信的方法。  

## sync.Pool

用Pool来尽可能快递将预先分配的对象缓存加载启动。  
在这种情况下，我们不是试图通过限制创建的对象的数量来节省主机的内存，而是通过提前加载获取饮用道另一个对象所需的时间。  


## Channel  

Closing a channel is also one of the ways you can signal multiple goroutines simultaneously  

Closing a channel is both cheaper and faster than performing `n` write  


