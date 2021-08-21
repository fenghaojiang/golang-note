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



### B+树索引
前面讨论都是B+树数据结构及其一般操作，B+树索引的本质就是B+树在数据库中的实现。但是B+索引在数据库中有一个特点是高扇出性，因此在数据库中，B+树的高度一般在2~4层，这也就是说查找某一键值的行记录只需要2~4次的IO。  


聚集索引与辅助索引不同的是，叶子节点存放的是否是一整行的信息。    

### 辅助索引   

对于辅助索引，叶子节点并不包含行记录的全部数据。叶子节点出了包含键值外，每个叶子节点中的索引行还包含一个书签。该书签用来告诉InnoDB存储引擎哪里可以找到与索引相对应的行数据。    

由于InnoDB存储引擎表是索引组织表，因此InnoDB存储引擎的辅助索引的书签就是相应行数据的索引聚集索引键。   



如果在一棵高度为3的辅助索引树中查找数据，那需要对这棵辅助索引树遍历3次找到指定主键，如果聚集索引树的高度同样为3，那么还需要对聚集索引树进行3次查找，最终找到一个完整的行数据所在的页，因此一共需要6次逻辑IO访问以得到最终的一个数据页。  

### B+树索引的管理  

1. 索引管理  
   ALTER TABLE tbl_name | ADD {INDEX|KEY} [index_name] [index_type] (index_col_name,...) [index_option] ...

   ALTER TABLE tbl_name DROP RRIMARY KEY | DROP {INDEX|KEY} index_name

   ```sql
   ALTER TABLE t ADD KEY idx_b (b(100));
   ```

   查看

   ```sql
   SHOW INDEX FROM t\G;
   ```

使用USE INDEX的索引提示可以来使用a这个索引，如： 

```sql
SELECT * FROM t USE INDEX(a) WHERE a = 1 and b = 2;
```

虽然指定使用a索引，但是优化器实际选择的是通过表扫描的方式。因此，USE INDEX只是告诉优化器可以选择该索引，实际上优化器还是会再根据自己判断进行选择。  


如果使用FORCE INDEX索引提示，优化器最终选择的和用户指定的索引是一致的。  


```sql
SELECT * FROM t FORCE INDEX(a) WHERE a = 1 AND b = 2;
```


### InnoDB存储引擎中的哈希算法  

InnoDB存储引擎使用哈希算法来对字典进行查找，其冲突机制采用链表方式，哈希函数采用除法散列方式。对于缓冲池页的哈希表来说，在缓冲池中的page页都有一个chain指针，它指向相同哈希函数值得页。而对于除法散列，h(k) = k mod m  
m的取值为略大于2倍的缓冲池页数量的质数。  


### 倒排索引  

+ inverted file index 仅存储文档ID
+ full inverted index 存储的是对(pair),还存储了单词所在的位置信息，占用更多空间，但是能更好的定位数据。  

### InnoDB全文索引  

InnoDB从1.2版本开始采用full inverted index方式。在InnoDB存储引擎中，将(DocumentID, Position)视为一个“ilist”. 因此在全文检索的表中，有两个列，一个是word，另一个是ilist。  


由于在ilist字段中存放了Position，所以可以进行Proximity Search  


倒排索引需要将word存放到一张表中，这个表称为Auxiliary Table(辅助表)。在InnoDB存储引擎中，为了提高全文检索的并行性能，共有6张Axillary Table 根据word的Latin编码进行分区。  

Auxiliary Table是存在磁盘上的持久的表。  

FTS Index Cache(全文检索索引缓存)(底层是红黑树)，其用来提高全文检索的性能。   


