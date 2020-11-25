---
title: Golang中的reflect
date: 2020-11-21
---


# 什么是reflect?  
很多人来问我说:离老师发生甚么事辣?  
我一看，哦，原来是reflect是甚么?  
reflect翻译过来就是反射，实现了运行时的反射能力，能够让程序操作不同类型的对象。  
reflect.TypeOf能获取类型信息，reflect.ValueOf能够获取数据的运行时表示  

反射包中的所有方法基本都是围绕着Type和Value这两个类型设计的。我们通过reflect.TypeOf、reflect.ValueOf可以将一个普通的变量转换成反射包提供的Type和Value    


类型```Type```是反射包定义的一个接口，我们可以使用`reflect.TypeOf`函数获取任意变量的类型, `Type`接口中定义了一些有趣的方法，`MethodByName`可以获取当前类型对应方法的引用、`Implements`可以判断当前类型是否实现了某个接口。  

```go
type Type interface {
    Align() int
    FieldAlign() int
    Method(int) Method
    MethodByName(string) (Method, bool)
    NumMethod() int
    ...
    Implements(u Type) bool
}
```


## 反射的三大法则  
反射带来的灵活性是一把双刃剑，反射作为一种元编程方式可以减少重复代码，但过量的使用会使我们的程序逻辑变得难以理解并且运行缓慢。在下面介绍Go语言反射中的三大法则  

+ 从interface{} 变量可以反射出反射对象
+ 从反射对象可以获取interface{}对象
+ 要修改反射对象，其值必修可设置  


**第一法则**
反射的第一法则是我们能将Go语言的interface{}变量转换成反射对象，朋友问，马老师，为什么是从interface{} 变量到反射对象? 当我们执行reflect.ValueOf(1)时，虽然看起来是获取了基本类型int对应的反射类型  
但是由于reflect.TypeOf、reflect.ValueOf两个方法的入参都是interface{} 类型，所以在方法执行的过程中发生了类型转换  
在函数调用一节曾经介绍过，Go语言的函数调用都是值传递的，变量会在函数调用的过程中进行类型转换，所以在方法执行的过程中发生了类型转换  基本类型int会转换成interface{}类型，这也就是为什么第一条法则『从接口到反射对象』  
上面提到的`reflect.TypeOf`和`reflect.ValueOf`函数就能完成这里的转换，如果我们认为Go语言的类型和反射类型处于两个不同的世界，那么两个函数就是连接这两个世界的桥梁。

源码：
```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	author := "draven"
	fmt.Println("TypeOf author:", reflect.TypeOf(author))
	fmt.Println("ValueOf author:", reflect.ValueOf(author))
}
```

输出:  
```go
$ go run main.go
TypeOf author: string
ValueOf author: draven
```


**第二法则**  
反射的第二法则是我们可以从反射对象获取interface{}变量  
既然能够将接口类型变量转换成反射对象，那么一定需要其他方法将反射对象还原成接口类型的变量， `reflect`中的reflect.Value.Interface方法就能完成这项工作：  

不过调用 reflect.Value.Interface 方法只能获得 interface{} 类型的变量，如果想要将其还原成最原始的状态还需要经过如下所示的显式类型转换：

```go
v := reflect.ValueOf(1)
v.Interface().(int)
```
从反射对象到接口值的过程就是从接口值到反射对象的镜面过程，两个过程都需要经历两次转换：

+ 从接口值到反射对象：
  + 从基本类型到接口类型的类型转换；
  + 从接口类型到反射对象的转换；
+ 从反射对象到接口值：
  + 反射对象转换成接口类型；
  + 通过显式类型转换变成原始类型；


当然不是所有的变量都需要类型转换这一过程。如果变量本身就是 interface{} 类型，那么它不需要类型转换，因为类型转换这一过程一般都是隐式的，所以我不太需要关心它，只有在我们需要将反射对象转换回基本类型时才需要显式的转换操作。  

**第三法则**  

反射的最后一条法则是 与值是否可以被更改有关，如果我们想要更新一个`reflect.Value`，那么它持有的值一定是可以被更新的，假设我们有以下代码:  
  
```go
func main() {
	i := 1
	v := reflect.ValueOf(i)
	v.SetInt(10)
	fmt.Println(i)
}
```

Output:
```go
$ go run reflect.go
panic: reflect: reflect.flag.mustBeAssignable using unaddressable value

goroutine 1 [running]:
reflect.flag.mustBeAssignableSlow(0x82, 0x1014c0)
	/usr/local/go/src/reflect/value.go:247 +0x180
reflect.flag.mustBeAssignable(...)
	/usr/local/go/src/reflect/value.go:234
reflect.Value.SetInt(0x100dc0, 0x414020, 0x82, 0x1840, 0xa, 0x0)
	/usr/local/go/src/reflect/value.go:1606 +0x40
main.main()
	/tmp/sandbox590309925/prog.go:11 +0xe0
```


运行上述代码会导致程序崩溃并报出 reflect: reflect.flag.mustBeAssignable using unaddressable value 错误，仔细思考一下就能够发现出错的原因，Go 语言的函数调用都是传值的，所以我们得到的反射对象跟最开始的变量没有任何关系，所以直接对它修改会导致崩溃。

想要修改原有的变量只能通过如下的方法：  

```go
func main() {
	i := 1
	v := reflect.ValueOf(&i)
	v.Elem().SetInt(10)
	fmt.Println(i)
}
```

Output:  
```go
$ go run reflect.go
10
```

go语言的函数都是用值传递，所以只能先获取指针对应的reflect.Value，再通过reflect.Value.Elem方法迂回的方式得到可以被设置的变量  
可以通过以下代码理解：  


```go
func main() {
	i := 1
	v := &i
	*v = 10
}
```



### 方法调用  
因为Golang是一门静态语言，想要通过reflect包利用反射在运行期间执行方法不是一件容易的事情  
下面的十几行代码就使用反射来执行Add(0, 1)函数

```go
func Add(a, b int) int { return a + b}

func main() {
	v := reflect.ValueOf(Add)
	if v.Kind() != reflect.Func {
		return 
	}
	t := v.Type()
	argv := make([]reflect.Value, t.NumIn())
	for i := range argv {
		if t.In(i).Kind() != reflect.Int {
			return
		}
		argv[i] = reflect.ValueOf(i)
	}
	result := v.Call(argv)
	if len(result) != 1 || result[0].Kind() != reflect.Int {
		return 
	}
	fmt.Println(result[0].Int)
}
```


1. 通过reflect.ValueOf获取函数Add对应的反射对象  
2. 根据反射对象reflect.rtype.NumIn方法返回的参数个数创建argv数组
3. 多次调用reflect.ValueOf函数逐一设置argv数组中的各个参数
4. 调用反射对象Add的reflect.Value.Call方法并传入参数列表
5. 获取返回值数组、验证数组长度以及类型并打印其中的数据
---


reflect.Value.Call 方法是运行时调用方法的入口，它通过两个 MustBe 开头的方法确定了当前反射对象的类型是函数以及可见性，随后调用 reflect.Value.call 完成方法调用，这个私有方法的执行过程会分成以下的几个部分：  

1. 检查输入参数以及类型的合法性；
2. 将传入的 reflect.Value 参数数组设置到栈上；
3. 通过函数指针和输入参数调用函数；
4. 从栈上获取函数的返回值；

