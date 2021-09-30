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








