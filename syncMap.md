---
title: Golang中sync.Map
date: 2020-12-4
---

参考链接： 

1. https://colobu.com/2017/07/11/dive-into-sync-Map/    
2. https://eddycjy.com/posts/go/map/2019-03-05-map-access/  




Golang中的map不支持并发读写，Golang团队为了解决这个问题实现了sync.Map  

## 在此前的解决方案  

使用嵌入struct的方式增加一个读写锁  

例如：  

```go
var counter = struct {
    sync.RWMutex
    m map[string]int
}{m: make(map[string]int)}
```

读取时:  

```go
counter.RLock()
n := counter.m["some_key"]
counter.RUnlock()
```

写入时：  
```go
counter.Lock()
counter.m["some_key"]++
counter.Unlock()
```

可以说，上面的解决方案已经是懂哥的解决思路，相当简洁漂亮美观  


