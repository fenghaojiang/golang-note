---
title: 多版本并发控制
date: 2021-11-11 
---  

设置隔离级别SQL
```sql
set session transaction isolation level read committed;
```

MySQL大多数事务性存储引擎实现的都不是简单的行级锁。基于提升并发性能的考虑，一般都会实现多版本并发控制MVCC。（Oracle、PostgreSQL）   

可以认为MVCC是行级锁的一个变种，但是很多情况下避免了加锁操作，因此开销更低。虽然实现机制有所不同，但大多实现了非阻塞的读操作。写操作也只锁定必要的行。   

MVCC的实现，是通过保存数据在某个时间点的快照来实现的。不同存储引擎的MVCC实现是不同的，典型的有乐观并发控制和悲观并发控制。   

悲观锁：对同一条数据修改，直接加锁。具有强烈的独占和排他性质。   
乐观锁：假设数据一般不会造成冲突，在数据提交时的才会正式对数据的冲突与否进行检测，如果冲突，则返回给用户异常信息，让用户决定如何去做。   
乐观锁适用于读多写少的场景，可以提高程序吞吐量。   


InnoDB通过每行后面保存两个隐藏的列来实现的，一列保存行的创建时间，另一列保存行的过期时间，当然存储的并不是实际的时间值，而是系统版本号。   

每开始一个新的事务，系统版本号都会自动递增。  

Repeatable read mvcc是如何操作的

+ Select
  InnoDB会根据以下两个条件检查每行记录： 
  1. InnoDB只查找版本早于当前事务版本的数据行（也就是行系统版本号小于等于事务的系统版本号），这样可以确保事务读取的行，要么是在事务开始前已经存在的，要么是事务自身插入或修改的
  2. 行的删除版本要么未定义，要么大于当前事务版本号。这可以确保事务读取到的行，在事务开始之前未被删除。

+ Insert
  InnoDB为新插入的每一行保存当前的系统版本号作为行版本号

+ Delete
  InnoDB为删除的每一行保存当前的系统版本作为行删除标识

+ Update
  InnoDB为插入一行新纪录，保存当前的系统版本号为行版本号，同时保存当前系统版本号到原来的行作为行删除标识。  



MVCC只在repeatable read和read committed两个级别下工作，其他两个隔离版本都与mvcc不兼容。  













