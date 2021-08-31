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
