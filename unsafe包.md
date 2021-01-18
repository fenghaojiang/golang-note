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



