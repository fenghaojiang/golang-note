---
title: unsafe包
date: 2021-1-18
---

## Golang的指针  
Go语言的指针相比C语言的指针有很多限制。这当然是为了安全考虑，要知道Java/Python这些现代语言直接把指针扬了嗷。  
而C/Cpp这些语言又要程序员自己手动去清理"垃圾"  
为什么Go要有指针类型呢？

有个懂哥举了个例子:

```go
package main

import "fmt"

func double(x *int) {
	*x += *x
	x = nil
}

func main() {
	var a = 3
	double(&a)
	fmt.Println(a) // 6
    
	p := &a
	double(p)
	fmt.Println(a, p == nil) // 12 false
}
```

这些代码都很常见，不用多解释，唯一可能有些疑惑的在这一句：

```go
x = nil
```
因为是值传递(Golang没有引用传递)，所以x也只是&a的一个副本

```go
*x += *x
```

这一句把x指向的值(也就是&a指向的值，即变量a)变为原来的2倍。但是对于x本身(一个指针)的操作却不会影响到外层的a，所以x=nil不会有任何改变  
相比于C/Cpp自由灵活的指针，Go的指针多了一些限制，既可以享受指针带来的便利，又避免了指针的危险性。

```go
a := 5
p := &a
p++

p = &a + 3
```

以上代码会编译错误，**invalid operation**  
不能对指针进行数学运算，而且

**不同类型的指针不能使用==或!=比较**  

**不同类型的指针变量不能相互赋值**  

```go
func main() {
	a := int(100)
	var f *float64
	
	f = &a
}
```

编译错误：  

```terminal
cannot use &a (type *int) as type *float64 in assignment
```



## 什么是unsafe  

前面说的指针是类型安全的，但是有了许多的限制，Go还有非类型安全的指针，这就是unsafe包提供的unsafe.Pointer  
某种情况下会使代码变得更高效，但是同时也会变得更危险。   

unsafe包用于Go编译器，在编译阶段使用。从名字就可以看出，它是不安全的，官方并不建议使用。  

那么为什么要用unsafe包呢，因为它可以绕过Go语言的类型系统，直接操作内存。例如：一般我们不能操作一个结构体的未导出成员，但是通过unsafe包就能做到。unsafe包让我可以直接读写内存。

## 为什么有unsafe  

Go语言类型系统是为了安全和效率设计的，有时，安全会导致效率底下。有了unsafe包，高阶的程序员就可以利用它绕过类型系统的低效。因此，它就有了存在的意义，阅读Go源码，会发现有大量使用unsafe包的例子  


### unsafe实现原理  

源码:  
```go
type ArbitraryType int 
type Point *ArbitraryType
```

这里普及一个生活小常识，Arbitrary是任意的意思，也就是说Pointer可以指向任意类型，实际上它类似于C语言的void*   


```go
func Sizeof(x ArbitraryType) uintptr
func Offsetof(x ArbitraryType) uintptr
func Alignof(x ArbitraryType) uintptr
```

Sizeof返回类型x所占据的字节数，但不包含x所指向的内容大小。例如，对一个指针，函数返回的大小为8字节(64位机上)，一个slice的大小则为slice header的大小   


Offsetof返回结构体成员在内存中的位置离结构体起始处的字节数，所传参数必须是结构体的成员   

Alignof返回m，m是指当类型进行内存对齐时，它分配到的内存地址能整除m  

上述3个函数返回的结构都是uintptr类型，代表着可以和unsafe.Pointer可以相互转换。三个函数都是在编译期间执行，它们的结果可以直接赋值给const型变量。另外，因为三个函数执行的结果和操作系统、编译器相关，所以是不可以移植的。  

综上，unsafe包提供了2点重要的能力:

*1.任何类型的指针和unsafe.Pointer可以相互转换*
*2.uintptr类型和unsafe.Pointer可以相互转换*

Pointer不能直接进行数学运算，但可以把它转换成uintptr，对uintptr类型进行数学运算，再转换成pointer类型  

uintptr并没有指针的语义，所以uintptr所指的对象会被gc无情地回收。而unsafe.Pointer有指针语义，可以保护它所指向的对象在“有用”的时候不会被垃圾回收。 

总结一下：  

1. 任何类型的指针都可以被转化成Pointer  
2. Pointer可以被转化成任何类型的指针  
3. uintptr可以被转化成Pointer  
4. Pointer可以被转化成uintptr  

uintptr是golang的内置类型，是能够存储指针的整型  

**一个unsafe.Pointer指针也可以被转化为uintptr类型，然后保存到指针型数值变量中（注：这只是和当前指针相同的一个数字值，并不是一个指针），然后用以做必要的指针数值运算。（uintptr是一个无符号的整型数，足以保存一个地址）这种转换虽然也是可逆的，但是将uintptr转为unsafe.Pointer指针可能会破坏类型系统，因为并不是所有的数字都是有效的内存地址。**



## 如何使用unsafe?  

```go
// runtime/slice.go
type slice struct {
	array unsafe.Pointer //元素指针
	len int // 长度
	cap int // 容量
}
```

调用make函数新建一个slice，底层调用的是makeslice函数，返回的是slice结构体   
```go 
func makeslice(et *_type, len, cap int) slice 
```

我们可以用过unsafe.Pointer和uintptr进行转换，得到slice的字段值   

demo:  

```go
func main() {
	s := make([]int, 9, 20)
	var Len = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(8)))
	fmt.Println(Len, len(s))

	var Cap = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(16)))

	fmt.Println(Cap, cap(s))
}
```

Len、Cap的转换过程如下：

```go
Len: &s => pointer => uintptr => pointer => *int => int
Cap: &s => pointer => uintptr => pointer => *int => int
```

## 获取map长度  

```go
type hmap struct {
	count     int
	flags     uint8
	B         uint8
	noverflow uint16
	hash0     uint32

	buckets    unsafe.Pointer
	oldbuckets unsafe.Pointer
	nevacuate  uintptr

	extra 	   *mapextra
}
```

和slice不同，makemap函数返回的是hmap的指针，注意是指针:  
```go
func makemap(t *maptype, hint int64, h *hmap, bucket unsafe.Pointer) *hmap
```

我们依然能通过unsafe.Pointer和uintptr进行转换，得到hamp字段的值，只不过，现在count变成二级指针了：  

```go
func main() {
	mp := make(map[string]int) 
	mp["nmsl"] = 100
	mp["sunxiaochuan"] = 258
	count := **(**int)(unsafe.Pointer(&mp)) //返回的是指针
}
```

count 转换过程:
```go
&mp => pointer => **int => int
```

## Offsetof获取成员偏移量  
对于一个结构体，通过offset函数可以获取结构体成员的偏移量，进而获取成员的地址，读写该地址的内存，就可以达到改变成员值的目的  

这里有一个生活小常识，结构体会被分配一块连续的内存，结构体的地址也代表了第一个成员的地址。  



eg:  

```go
package main

import (
	"fmt"
	"unsafe"
)

type SunxiaochuanNMSL struct {
	name     string
	language string
}

func main() {
	p := SunxiaochuanNMSL{"SunXiaoChuan", "Chouxianghua"}
	fmt.Println(p)

	name := (*string)(unsafe.Pointer(&p))

	*name = "daidaidashixiong"

	lang := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + unsafe.Offsetof(p.language)))
	*lang = "nmsl"

	fmt.Println(p)
}
```

输出:  

```terminal
{SunXiaoChuan Chouxianghua}
{daidaidashixiong nmsl}
```

name是结构体的第一个成员，因此可以直接将&p解析成*string。这一点，在前面获取map的count成员时，用的是同样的原理。   

对于结构体的私有成员，现在有方法可以通过unsafe.Pointer改变它的值了。   

把上面的结构体多加一个字段  

```go
type SunxiaochuanNMSL struct {
	name     string
	age      int
	language string
}
```

并且放在其他包，这样在main函数中，他的三个字段都是私有成员变量，不能直接修改。但可以通过unsafe.Sizeof()获取到成员的大小，进而计算出成员的地址，直接修改内存。   

```go
package main

import (
	"fmt"
	"unsafe"
)

type SunxiaochuanNMSL struct {
	name     string
	age      int
	language string
}

func main() {
	p := SunxiaochuanNMSL{"SunXiaoChuan", 257, "Chouxianghua"}
	fmt.Println(p)

	name := (*string)(unsafe.Pointer(&p))
	*name = "daidaidashixiong"

	lang := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + 8 + unsafe.Sizeof(string(""))))
	*lang = "nmsl"

	fmt.Println(p)
}
```

## string和slice的相互转换  

接下来的是一个经典经典典中典的例子。实现字符串和bytes切片之间的转换，要求是zero-copy。想一下，一般的做法，都需要遍历字符串或者bytes切片，再挨个赋值。  

首先我们要了解slice和string底层的数据结构：  

```go
type StringHeader struct {
	Data uintptr
	Len  int
}

type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
```  

上面是反射包下的结构体，我们只需要贡献底层的[]byte数组就可以实现zero-copy  


```go
type StringHeader struct {
	Data uintptr
	Len  int
}

type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

func string2bytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func bytes2string(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))
}
```

uintptr可以和unsafe.Pointer进行相互转换，uintptr可以进行数学运算，通过两者结合就解决了Go指针不能进行数学运算的限制








