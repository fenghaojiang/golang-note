---
title: MySQL笔记
date: 2021-04-18 
---


事务数据库系统普遍采用了Write Ahead Log策略，当事务提交时，先写重做日志，再修改页。当由于发生宕机而导致数据丢失时，通过重做日志来完成数据的恢复。这也是事务ACID中的D(Durability)


### InnoDB的后台线程  

1. Master Thread  
Master Thread负责将缓冲池中的数据异步刷新的磁盘，保证数据的一致性，脏页的刷新、合并插入缓冲(Insert Buffer)、Undo页的回收

2. IO Thread  
负责IO请求的回调

3. Purge Thread
事务被提交后，其所使用的undolog可能不再需要，因此需要PurgeThread来回收已经使用并分配的undo页


### Checkpoint技术
+ 缩短数据库的恢复时间
+ 缓冲池不够用时，将脏页刷新到磁盘
+ 重做日志不可用时，刷新脏页

当数据库发生宕机时，数据库不需要重做所有的日志，因为Checkpoint之前的页都已经刷新回磁盘。故数据库只需对Checkpoint后的重做日志进行恢复，这样就大大缩短了回复的时间。

此外，当缓冲池不够用时，根据LRU算法会溢出最近最少使用的页，若此页为脏页，那么需要强制执行Checkpoint，将脏页也就是页的新版本刷回磁盘。  

在InnoDB存储引擎内部，有两种Checkpoint，分别为:  
+ Sharp Checkpoint
+ Fuzzy Checkpoint

Sharp checkpoint---数据库关闭时将所有脏页都刷新回磁盘，但是若在数据库运行时也使用Sharp checkpoint 那么数据库的可用性就会收到很大的影响。所以在InnoDB存储引擎内部使用Fuzzy Checkpoint进行页的刷新，即只刷新一部分脏页，而不是刷新所有的脏页回磁盘

可能发生Fuzzy Checkpoint的情况:  
+ Master Thread Checkpoint  
异步进行，差不多每秒或者每10秒的速度从缓冲池的脏页列表中刷新一定比例的页回磁盘。  

+ FLUSH_LRU_LIST Checkpoint  


+ Async/Sync Flush Checkpoint  


+ Dirty Page too much Checkpoint  










