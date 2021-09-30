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





