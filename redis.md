---
title: Redis
date: 2021-09-30 
---  


### SDS的定义   

sdshar结构

```cpp
struct sdshdr {
    int len;
    int free;
    char buf[];
};
```

SDS中，buf数组的长度不一定就是字符数量加一，数组里面可以包含未使用的字节，二这些字节的数量就由SDS的free属性记录。  


通过未使用空间，SDS实现了两种优化策略  

+ 空间预分配 
+ 惰性空间释放  


+ Redis指挥使用C字符串作为字面量，在大多数情况下，Redis使用SDS作为字符串表示
+ 比起C字符串，SDS优点在于:  
  1. 常熟级复杂度获取字符串长度
  2. 杜绝缓冲区溢出
  3. 减少修改字符串长度时所需的内存重分配次数
  4. 二进制安全
  5. 兼容部分C字符串函数



### 链表   

节点数据结构
```cpp
typedef struct listNode {
    struct listNode *prev;
    struct listNode *next;

    void *value;
}listNode;
```


链表数据结构


```cpp
typedef struct list {
    //头节点
    listNode *head;
    //尾节点
    listNode *tail;
    //包含节点数量
    unsigned long len;
    //节点值复制函数
    void *(*dup) (void *ptr);
    //节点值释放函数
    void *(*free) (void *ptr);
    //节点值对比函数
    void *(*match) (void *ptr, void *key);
} list;
```


### 字典  

哈希表数据结构
```cpp
typedef struct dictht {
    //哈希表数组
    dictEntry **table;
    //哈希表大小
    unsigned long size;
    //哈希表大小掩码，用于计算索引值
    //总是等于size-1
    unsigned long sizemask;

    //该哈希表已有节点的数量
    unsigned long used;
}dictht;
```

哈希表节点
```cpp
typedef struct dictEntry {
    //键
    void *key;

    //值
    union {
        void *val;
        uint64_tu64;
        int64_ts64;
    } v;

    //指向下个哈希表节点，形成链表
    struct dictEntry *next;
}dictEntry;
```

### 字典  

```cpp
typedef struct dict {
    // 类型特定函数
    dictType *type;
    // 私有数据
    void *privdata;
    // 哈希表
    dictht ht[2];
    // rehash索引
    in trehashidx;
}dict;
```

Redis用链地址法来解决键冲突，每个哈希表节点都有一个next指针，多个哈希表节点可以用next指针构成一个单向链表，被分配到同一个索引上的多个节点可以用单向链表连接起来，这就解决了键冲突的问题。  










```shell
127.0.0.1:6379> llen integers
(integer) 0
127.0.0.1:6379> lrange integers 0 10
(empty array)
127.0.0.1:6379> hset web msg "hello"
(integer) 1
127.0.0.1:6379> hlen web
(integer) 1
127.0.0.1:6379> hgetall web
1) "msg" #键
2) "hello" #值
```












