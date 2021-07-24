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





