---
title: InnoDB相关part5 索引与算法  
date: 2021-08-09 
---  


### B+树的插入操作  
| Leaf Page 满 | Index Page 满| 操作| 
|:--:|:--:|:--:|
|No|No|直接将记录插入到叶子节点|
|Yes|No|1)拆分Leaf Page <br> 2)将中间的节点放入到Index Page中 <br> 3)小于中间节点的记录放左边 <br> 4) 大于或等于中间节点的记录放右边
|Yes|Yes|1)拆分Leaf Page <br> 2)小于中间节点的记录放左边 <br> 3)大于或等于中间节点的记录放右边 <br> 4)拆分Index Page <br> 5)小于中间节点的记录放左边 <br> 6)大于中间节点的记录放右边 <br> 7)中间节点放入上一层Index Page <br>|








