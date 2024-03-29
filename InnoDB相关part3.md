---
title: InnoDB相关
date: 2021-07-24 
--- 

ENUM进行约束规定域范围  

```sql
CREATE TABLE a (
    id INT,
    sex ENUM('male', 'female'));
```  

```sql
INSERT INTO a
SELECT 2, 'bi';
``` 
会显示如下警告  

```shell script
ERROR 1265 (01000): Data truncated for column 'sex' at row 1
``` 

### 触发器与约束  


触发器的作用是在INSERT、DELETE和UPDATE命令之前或之后自动调用SQL命令或存储过程。  

创建触发器的命令是CREATE TRIGGER，只有具备Super权限的MySQL数据库用户才可以执行这条命令:  

```sql
UPDATE usercash
SET cash=cash-(-20) WHERE userid=1;
```  

通过触发器约束上述逻辑行为，可以如下设置：  

```sql
CREATE TABLE usercash_err_log (
    userid INT NOT NULL,
    old_cash INT UNSIGNED NOT NULL,
    new_cash INT UNSIGNED NOT NULL,
    user VARCHAR(30),
    time DATETIME
);
```

```sql
CREATE TRIGGER tgr_usercash_update BEFORE UPDATE ON usercash
FOR EACH ROW
BEGIN
if new.cash-old.cash > 0 THEN
INSERT INTO usercash_err_log
SELECT old.userid, old.cash, new.cash, USER(), NOW();
SET new.cash = old.cash;
END IF;
END;
$$
``` 

```sql
DELIMITER $$
```  


### 4.6.7 外键约束

外键用来保证参照完整性，MySQL数据库的MyISAM存储引擎本身并不支持外键，对于外键的定义只是起到一个注释的作用。而InnoDB存储引擎则完整外键约束。  



```sql
CREATE TABLE parent (
    id INT NOT NULL,
    PRIMARY KEY (id)
)ENGINE=INNODB;

CREATE TABLE child (
    id INT, parent_id INT,
    FOREIGN KEY (parent_id) REFERENCES parent(id)
) ENGINE=INNODB;
``` 

被引用的是父表，引用的是子表， 外键定义时的ON DELETE和ON UPDATE表示在对父表进行DELETE和UPDATE操作时，对子表所做的操作，可定义的子表操作有：  


+ CASCADE
+ SET NULL
+ NO ACTION
+ RESTRICT

而InnoDB存储引擎会在外键建立时会自动对该列加一个索引


```sql
SHOW CREATE TABLE child\G;
```


```sql
CREATE TABLE `child` (
  `id` int DEFAULT NULL,
  `parent_id` int DEFAULT NULL,
  KEY `parent_id` (`parent_id`),
  CONSTRAINT `child_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `parent` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
```  


