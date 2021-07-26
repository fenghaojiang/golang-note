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

