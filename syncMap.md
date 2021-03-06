---
title: Golang中sync.Map
date: 2020-12-4
---

参考链接： 

1. https://colobu.com/2017/07/11/dive-into-sync-Map/    
2. https://eddycjy.com/posts/go/map/2019-03-05-map-access/  
3. https://mp.weixin.qq.com/s/kblDTqKlUaTITIppigq9yA




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

## 为什么大家宁愿使用Mutex + Map, 也不愿使用sync.Map:  

1. sync.Map本身就很难用，使用起来根本就不像一个Map那样简单。失去了map应有的特权语法。  
2. sync.Map方法较多。让一个简单的Map使用起来有了较高的学习成本。  

用法demo

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var syncMap sync.Map
	syncMap.Store("11", 11)
	syncMap.Store("22", 22)

	fmt.Println(syncMap.Load("11")) // 11
	fmt.Println(syncMap.Load("33")) // 空

	fmt.Println(syncMap.LoadOrStore("33", 33)) // 33
	fmt.Println(syncMap.Load("33"))            // 33

	syncMap.Delete("33")            // 33
	fmt.Println(syncMap.Load("33")) // 空

	syncMap.Range(func(key, value interface{}) bool {
		fmt.Printf("key:%v value:%v\n", key, value)
		return true
	})
	// key:22 value:22
	// key:11 value:11
}
```

其实 sync.Map 并不复杂，只是将普通 map 的相关操作转成对应函数而已。   

<br>


||普通map|sync.Map|
|:--:|:--:|:--:|
|map获取某个元素|map[1]|sync.Load(1)|
|map添加元素|map[1] = 10|sync.Store(1,10)|
|map删除一个key|delete(map, 1)|sync.Delete(1)|
|遍历map|for...range|sync.Range()|  

sync.Map 两个特有的函数，不过从字面就能理解是什么意思了。LoadOrStore：sync.Map 存在就返回，不存在就插入 LoadAndDelete：sync.Map 获取某个 key，如果存在的话，同时删除这个 key  

## 源码
```go
type Map struct {
    mu Mutex
    read atomic.Value // readOnly read map
    dirty map[interface{}]*entry
    misses int
}
```

![syncMap](./img/syncMap.webp)  


![syncMap2](./img/syncMap2.webp)    

![syncMap3](./img/syncMap3.webp)   


### read map的值是什么时间更新的？  
1. Load/LoadOrStore/LoadAndDelete时，misses数量大于等于dirty map元素个数时，会整体复制dirty map到read map
2. Store/LoadOrStore时，当read map中存在这个key，则更新  
3. Delete/LoadAndDelete时，如果read map中存在这个key，则设置这个值为nil  

### dirty map的值是什么时间更新的？   
1. 完全是一个新key，第一次插入sync.Map，必先插入dirty map  
2. Store/LoadOrStore时，当read map中不存在这个key，在dirty map存在这个key，则更新  
3. Delete/LoadAndDelete时，如果read map中不存在这个key，在dirty map存在这个key，则从dirty map中删除这个key  
4. 当misses数量大于等于dirty map的元素个数时，会整体复制dirty map到read map，同时设置dirty map为nil  

疑问：当 dirty map 复制到 read map 后，将 dirty map 设置为 nil，也就是 dirty map 中就不存在这个 key 了。如果又新插入某个 key，多次访问后达到了 dirty map 往 read map 复制的条件，如果直接用 read map 覆盖 dirty map，那岂不是就丢了之前在 read map 但不在 dirty map 的 key ?  


答：其实并不会。当 dirty map 向 read map 复制后，readOnly.amended 等于了 false。当新插入了一个值时，会将 read map 中的值，重新给 dirty map 赋值一遍，也就是 read map 也会向 dirty map 中复制  

```go
func (m *Map) dirtyLocked() {
    if m.dirty != nil {
        return 
    }
    read, _ := m.read.Load().(readOnly)
    m.dirty = make(map[interface{}]*entry, len(read.m))
    for k, e := range read.m {
        if !e.tryExpungeLocked() {
            m.dirty[k] = e
        }
    }
}
```

### read map和dirty map是什么时间删除的?  

+ 当read map中存在某个key时候，这个时候只会删除read map，并不会删除dirty map(因为dirty map不存在这个值)  
+ 当read map中不存在，才会去删除dirty map里面的值  

疑问：如果按照这个删除方式，那岂不是 dirty map 中会有残余的 key，导致没删除掉？

答：其实并不会。当 misses 数量大于等于 dirty map 的元素个数时，会整体复制 read map 到 dirty map。这个过程中还附带了另外一个操作：将 dirty map 置为 nil。  


```go
func (m *Map) missLocked() {
    m.misses++
    if m.misses < len(m.dirty) {
        return 
    }
    m.read.Store(readOnly{m: m.dirty})
    m.dirty = nil
    m.misses = 0
}
```

### read map 与 dirty map 的关系 ？  

+ 在 read map 中存在的值，在 dirty map 中可能不存在。
+ 在 dirty map 中存在的值，在 read map 中也可能存在。
+ 当访问多次，发现 dirty map 中存在，read map  中不存在，导致 misses 数量大于等于 dirty map 的元素个数时，会整体复制 dirty map 到 read map。
+ 当出现 dirty map 向 read map 复制后，dirty map 会被置成 nil。
+ 当出现 dirty map 向 read map 复制后，readOnly.amended 等于了 false。当新插入了一个值时，会将 read map 中的值，重新给 dirty map 赋值一遍

### read/dirty map 中的值一定是有效的吗？  

+ nil: 如果获取的value是nil，那说明这个key是已经删除过的。既不在read map，也不在dirty map
+ expunged: 这个key在dirty map中是不存在的
+ valid: 存在于两者其中一个  

### sync.Map 是如何提高性能的？  

通过源码解析，我们知道sync.Map里面有两个普通的map, read map主要是负责读，dirty map是负责读和写(加锁)。在读多写少的场景下，read map的值基本不发生变化，可以让read map做到无锁操作，就减少了使用Mutex+Map必须加锁/解锁的环节，因此也就提高了性能。  

不过也能够看出来，read map 也是会发生变化的，如果某些 key 写操作特别频繁的话，sync.Map 基本也就退化成了 Mutex + Map（有可能性能还不如 Mutex + Map）。

所以，不是说使用了 sync.Map 就一定能提高程序性能，我们日常使用中尽量注意拆分粒度来使用 sync.Map。

关于如何分析 sync.Map 是否优化了程序性能，同样可以使用 pprof。具体过程可以参考 《这可能是最容易理解的 Go Mutex 源码剖析》  


### 应用场景  

1. 读多写少  
2. 写操作也多，但是修改的key和读取的key特别不重合。  

关于第二点我觉得挺扯的，毕竟我们很难把控这一点，不过由于是官方的注释还是放在这里。

实际开发中我们要注意使用场景和擅用 pprof 来分析程序性能。  


### sync.Map 使用注意点  

和 Mutex 一样， sync.Map 也同样不能被复制，因为 atomic.Value 是不能被复制的。  








