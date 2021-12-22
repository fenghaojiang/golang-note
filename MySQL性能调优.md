---
title: MySQL性能调优
date: 2021-12-22   
---  

## 剖析单条查询  

```sql
set profiling = 1; //默认是金庸的
```
设置完之后在服务器上执行的所有语句都会测量耗费的时间和其他的一些查询执行状态变更的数据。   


```sql
show profiles;
```

```sql 
mysql> show profiles;
+----------+------------+-------------------+
| Query_ID | Duration   | Query             |
+----------+------------+-------------------+
|        1 | 0.05957500 | show databases    |
|        2 | 0.00034350 | SELECT DATABASE() |
|        3 | 0.02670975 | show tables       |
|        4 | 0.04739675 | select * from vod |
+----------+------------+-------------------+
4 rows in set, 1 warning (0.00 sec)
```

会给出一份报告这份报告会给出每个步骤还有花费的时间。  


+ show status可以显示某些活动读取索引的频繁程度，但是没有办法给出消耗了多少时间。 
+ show processlist可以用来观察是否有大量线程处于不正常的状态或者有其他不正常的特征。  



## 慢查询日志  
开启慢查询日志 , 配置样例：

/etc/mysql/my.cnf
[mysqld]
log-slow-queries   



```sql
show variables like 'long%';
```


```sql
+-----------------+-----------+
| Variable_name   | Value     |
+-----------------+-----------+
| long_query_time | 10.000000 |
+-----------------+-----------+
1 row in set, 1 warning (0.01 sec)
```

```sql
set long_query_time = 1;
```












