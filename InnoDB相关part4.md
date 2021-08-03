---
title: InnoDB相关part4
date: 2021-08-02 
--- 



当前MySQL数据库支持以下几种类型的分区   

+ RANGE分区:  行数据基于属于一个给定连续区间的列值被放入分区  
  
  ```sql
  CREATE TABLE t (
      id INT
  )ENGINE=INNODB
  PARTITION BY RANGE (id) (
      PARTITION p0 VALUES LESS THAN (10),
      PARTITION p1 VALUES LESS THAN (20)
  );
  ```

+ LIST分区: 和RANGE分区类似,只是LIST分区面向的是离散的值
  
  ```sql
  CREATE TABLE t (
      a INT,
      b INT)ENGINE=INNODB
      PARTITION BY LIST(b) (
          PARTITION p0 VALUES IN (1,3,5,7,9),
          PARTITION p1 VALUES IN (0,2,4,6,8)
      );
  ```

+ HASH分区: 根据用户自定义的表达式返回值来进行分区,返回值不能为负数.Hash分区的目的的时将数据均匀地分布到预先定义的各个分区中,保证各分区的数据量大致是一样的.在RANGE和LIST分区中,必须明确指定一个给定的列值或者列值集合应该保存在那个分区中,而在HASH分区中,MySQL自动完成这些工作,用户所要做得只是基于将要进行哈希分区的列值指定一个列值或表达式,以及指定被分区的表将要被分割成的分区数量.创建一个HASH分区的表t,分区按日期列b进行,分区数量为4:  PARTITION BY HASH (*expr*) expr是一个返回一个整数的表达式

    ```sql
    CREATE TABLE t_hash (
        a INT,
        b DATETIME
    )ENGINE=INNODB
    PARTITION BY HASH (YEAR(b))
    PARTITIONS 4;
    ```

+ KEY分区: 根据MySQL数据库提供的哈希函数来进行分区.在KEY分区中使用关键字LINEAR和在HASH分区中使用具有同样的效果,分区的编号是通过2的幂,在KEY分区中使用关键字LINEAR和在HASH分区中使用具有同样的效果,分区的编号是通过2的幂算法得到的,而不是通过模数算法.    

    ```sql
    CREATE TABLE t_key (
        a INT,
        b DATETIME)ENGINE=InnoDB
        PARTITION BY KEY (b)
        PARTITION 4;
    ```

不论创建何种类型的分区,如果表中存在主键或唯一索引时,分区列必须是唯一索引的一个组成部分,因此下面创建分区的SQL语句会产生错误.   
    
    ```sql
    CREATE TABLE t1 (
        col1 INT NOT NULL,
        col2 DATE NOT NULL,
        col3 INT NOT NULL,
        col4 INT NOT NULL,
        UNIQUE KEY (col1, col2)
    ) PARTITION BY HASH(col3)
    PARTITIONS 4;

    ERROR 1503 (HY000): A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).
    ``` 



MySQL数据库还支持一中称为LINEAR HASH的分区,它使用一个更加复杂的算法来确定新行插入到已经分区的表中的位置. 它的语法和HASH分区的语法相似,只是将关键字HASH改为LINEAR HASH.  

LINEAR HASH分区的优点在于,增加/删除/合并和拆分分区变得更加快捷,这有利于处理含有大量数据的表.它的缺点在于,与使用HASH分区得到的数据分布相比,各个分区间数据的分布可能不大均衡.  


+ COLUMNS分区  
  在前面介绍RANGE、LIST、HASH和KEY这四种分区中，分区的条件是：数据必须是整型(interger)，如果不是整型，需要通过函数将其转化成整型。COLUMNS分区可以用非整型的数据进行分区，分区根据类型直接比较而得到，不需要转化成整型。  

   
