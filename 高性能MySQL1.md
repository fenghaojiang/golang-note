---
title: 高性能MySQL
date: 2021-11-29
---



# MyISAM存储引擎  


对于只读的数据，或者表比较小、可以忍受修复（repair）操作，则仍然可以继续使用MyISAM   

MyISAM是非聚簇索引

可以转换引擎的语句：  

+ ALTER TABLE
  ```sql
  ALTER TABLE mytable ENGINE = InnoDB;
  ```
+ mysqldump

+ CREATE和SELECT
  ```sql
  CREATE table innodb_table like myisam_table;
  alter table innodb_table ENGINE=InnoDB;
  insert into innodb_table select * from mysiam_table;
  ```

  




