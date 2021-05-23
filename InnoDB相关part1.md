---
title: InnoDB相关
date: 2021-04-18 
---

## InnoDB相关

事务数据库系统普遍采用了Write Ahead Log策略，当事务提交时，先写重做日志，再修改页。当由于发生宕机而导致数据丢失时，通过重做日志来完成数据的恢复。这也是事务ACID中的D(Durability)


### InnoDB的后台线程  

1. Master Thread  
Master Thread负责将缓冲池中的数据异步刷新的磁盘，保证数据的一致性，脏页的刷新、合并插入缓冲(Insert Buffer)、Undo页的回收

2. IO Thread  
负责IO请求的回调(Async IO)

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
异步进行，差不多每秒或者每10秒的速度从缓冲池的脏页列表中刷新一定比例的页回磁盘，用户查询线程不会阻塞  

+ FLUSH_LRU_LIST Checkpoint  
InnoDB存储引擎需要保证LRU列表中需要差不多100个空闲页可供使用。如果没有，InnoDB会将LRU列表尾端的页移除。如果这些页有脏页，那么需要进行Checkpoint

+ Async/Sync Flush Checkpoint  
重做日志不可用的情况，这时需要强制将一些页刷新回磁盘，此时脏页是从脏页列表中选取的。

+ Dirty Page too much Checkpoint  
脏页数目太多，导致InnoDB强制进行Checkpoint。总的来说还是为了保证缓冲池有足够可用的页。


### Master Thread
Master Thread具有最高的线程优先级别。由主循环Loop、后台循环backgroup loop、刷新循环flush loop、暂停循环suspend loop

### InnoDB关键特性
+ 插入缓冲
  插入缓冲不是缓冲池中的一个部分。在InnoDB存储引擎中，主键是行唯一的标识符。通常应用程序中行记录的插入顺序是按照主键递增的顺序进行插入的。插入聚集索引(Primary Key)一般是顺序的，不需要磁盘的随机读取。比如按下列SQL定义表：

```sql
CREATE TABLE t (
  a INT AUTO_INCREMENT,
  b VARCHAR(30),
  PRIMARY KEY(a)
);
```
a列自增长。页中的行记录按a的值进行顺序存放。一般情况下不需要读取另一个页中的记录。这类情况下的插入操作，速度是非常快的。（但并不是所所有的主键插入都是顺序的，比如uuid）  

Insert Buffer的使用需要同时满足：
1. 索引是辅助索引
2. 索引不是唯一(unique)的

+ 两次写
+ 自适应哈希索引
+ 异步IO
+ 刷新邻接页 


### 自适应哈希索引  

哈希(hash)是一种非常快的查找方法，在一般情况下这种时间复杂度O(1)，即一般仅需要一次查找就能定位数据。而B+数的查找次数，取决于B+的高度，在生产环境中，B+树的高度一般就为3~4层，所以一般需要3~4次的查询  

InnoDB存储引擎会监控表上各索引页的查询。如果观察到建立哈希索引可以带来速度提升，则建立哈希索引。这就叫自适应哈希索引(Adaptive Hash Index, AHI)  

AHI有一个要求，即对这个页的连续访问模式必须是一样  

+ WHERE a=xxx
+ WHERE a=xxx and b=xxx  
  
访问模式一样是查询条件一样，如果交替进行上述两种查询，则不会对该页构造AHI  
AHI还有如下要求:  
+ 以该模式访问了100次  
+ 页通过该模式访问了N次，其中N=页中记录*16  



## 异步IO  
为了提高磁盘操作性能，当前的数据库系统都采用异步IO(Asynchronous IO)的方式来处理磁盘操作。InnoDB存储引擎亦是如此。  
用户可以在发出一个IO请求后立即再发出另一个IO请求，当全部IO请求发送完毕后，等待所有操作的完成，这就是AIO  

AIO的另一个优势是可以进行IO Merge操作，也就是将多个IO合并成1个IO。例如用户需要访问页的(space, page_no) 为：
(8, 6)、(8, 7)、(8, 8)  

每个页的大小为16KB，那么同步IO需要进行三次IO操作，而AIO会判断到这三个页是连续的，AIO底层会发送一个IO请求，从(8, 6)读取48KB的页。














