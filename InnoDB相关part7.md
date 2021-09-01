---
title: InnoDB相关part7 事务  
date: 2021-08-30 
---    

InnoDB存储引擎当中的事务完全符合ACID的特性。  

+ 原子性(atomicity)
+ 一致性(consistency)
+ 隔离性(isolation)
+ 持久性(durability)


+ 事务(transaction)指一组SQL语句
+ 回滚(rollback)指撤销指定SQL语句的过程
+ 提交(commit)指将未存储的SQL语句结果写入数据库表
+ 保存点(savepoint)指事务处理中设置的临时占位符(placeholder)，你可以对它发布回退  

COMMIT
```sql
START TRANSACTION;



DELETE FROM ......
COMMIT;
```


保存点
```sql
SAVEPOINT delete1;
ROLLBACK TO delete1;
```


### 事务的实现  
事务隔离性由锁来实现。原子性、一致性、持久性通过数据库的redo log和undo log来完成。redo log称为重做日志，用来保证事务的原子性和持久性。undo log用来保证事务的一致性。   


redo恢复提交事务修改的页操作，undo回滚行记录到某个特定版本。    

在MySQL数据库中还有一种二进制日志(binlog)，其用来进行POINT-IN-TIME(PIT)恢复及主从复制环境的建立。从表面上看其和重做日志非常相似，都是记录了对于数据库操作的日志，然而本质上有很大的区别   

1. 重做日志是在InnoDB存储引擎层产生，而二进制日志是在MySQL数据库上层产生的，并且二进制日志不仅仅针对InnoDB存储引擎，MySQL数据库中的任何存储引擎对于数据库的更改都会产生二进制日志。   
2. 两种日志的记录内容形式不同。MySQL数据库上层的二进制日志是一种逻辑日志，记录的是对应的SQL语句，InnoDB的重做日志是物理格式日志，记录每个页的修改  
3. 两种日志的写入磁盘时间点不同，二进制日志只在事务提交完成后进行一次写入。而InnoDB存储引擎的重做日志在事务进行中不断被写入。   


二进制日志不是幂等的，重复执行可能插入多条重复的记录，而insert操作的重做日志时幂等的  



### LSN  

LSN是Log Sequence Number的缩写，其代表的是日志序列号。在InnoDB存储引擎中，LSN占用8字节，并且单调递增。LSN表示的含义有:   

+ 重做日志写入的总量
+ checkpoint的位置  
+ 页的版本   

InnoDB存储引擎在启动时不管上次数据库运行时是否正常关闭，都会尝试进行恢复操作
页1的LSN为10000，而数据库启动时，InnoDB检测到写入重做日志中的LSN为13000，并且该事务已经提交，那么数据库需要进行恢复操作，将重做日志应用到1页中   

### undo   

undo是逻辑日志，因此只是将数据库逻辑地恢复到原来的样子。   
当InnoDB存储引擎进行回滚时，它实际上做的是与先前相反的工作。对于每个INSERT，InnoDB存储引擎会完成一个DELETE；对于每个DELETE存储引擎会执行一个INSERT...   

除了回滚，undo的另一个作用是MVCC，InnoDB存储引擎中的MVCC是通过undo来完成。当用户读取一行记录时，若该记录已经被其他事务占用，当前事务可以通过undo来读取之前的行版本信息，一次实现非锁定读取   

undo log 会产生redo log，这也是undo log也需要持久性的保护    




