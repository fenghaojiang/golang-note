---
title: 多版本并发控制
date: 2021-11-11 
---  

MySQL大多数事务性存储引擎实现的都不是简单的行级锁。基于提升并发性能的考虑，一般都会实现多版本并发控制MVVC。（Oracle、PostgreSQL）   




