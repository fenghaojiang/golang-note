---
title: 对Go语言的理解
date: 2020-10-29
---
  
# 对Go语言的理解

**Java与Go的区别**  

学习Golang已有三个月时间，在工作这段时间一直在使用Go进行游戏开发。  
有师弟问我Java跟Go哪个简单一点，我不知道怎么回答。在我看来Golang更容易上手一点，因为其无与伦比的简洁性和处理高并发时带来的便利都是Java给不了开发者的体验。  
Java本来就属于一种面向对象开发的语言，Java开发中很多东西都是约定俗成的规矩，要我打个比喻，Java就像是一个戴眼镜的学院派的学霸，做事情规矩很多，但是轮子工具也会很多。  
Go更像一个年轻的不拘一格的艺术家，很多时候开发者只是用了面向对象开发的思想去应用到Golang上，感觉其实Golang并不是一种面向对象开发的语言，感觉是很像C语言但是不同的是Golang又有GC，有GC这一点就跟C语言有很大不用，这就离谱儿。  
Golang中没有规定很多的划分方式，因为Golang的设计思路就是简洁，所以github上Golang的开源工具遵守了Golang简洁的设计方式使得Golang的代码通俗易懂很容易上手，对比Java，体验不要好太多了。Java的代码本身就长得要命，用久了Go每次读Java的代码本身就是一种折磨，WDNMD。  
Golang跟Python、JavaScript又不一样，是一种静态的语言  
什么是静态语言，那就是在运行前(比如编译过程中)就去做类型的检查，比如Java、C#、CPP/C都是静态语言，俄日Python、Ruby这些属于动态语言  

* 动态类型语言：在运行期间才去做数据类型检查的语言，在用动态语言编程时，不用给变量指定数据类型，该语言会在你第一次赋值给变量时，在内部将数据类型记录下来  

* 静态类型语言：数据类型检查发生在在编译阶段，也就是说在写程序时要声明变量的数据类型

此外，还有强类型语言跟弱类型语言的区别  
+ 强类型语言：使之强制数据类型定义的语言。没有强制类型转化前，不允许两种不同类型的变量相互操作。强类型定义语言是类型安全的语言，如Java、C# 和 Python，比如Java中“int i = 0.0;”是无法通过编译的  

+ 弱类型语言：数据类型可以被忽略的语言。与强类型语言相反, 一个变量可以赋不同数据类型的值，允许将一块内存看做多种类型，比如直接将整型变量与字符变量相加。C/C++、PHP都是弱类型语言，比如C++中“int i = 0.0;”是可以编译运行的

---

**Golang面向对象开发**  
面向对象的三大特点：封装、继承、多态。  

**封装**：先来说说Golang的封装，Golang的封装用了变量、函数的开头大小写来区分，小写是private，大写就是public，我愿称之为企业级好活，当传入参数的时候，要修改值就传引用(注:map本身就是引用类型，传map变量就是传引用)，不修改就传值。  

example:  
```go
func(p *type) Modify(value int) {
    p.value = value //穿址修改
}

func(p type) Modify(value int) {
    p.value = value //传值,没有改变
}
```

**继承**: Golang中没有继承的关键字，是通过struct来实现的，感觉struct相当于Java的Class，里面可以放入方法、变量。Go语言很多人说不推荐继承的概念，但在开发过程中或多或少利用好这些概念完全可以是办公版，与其他语言不同的是，Golang支持用接收者的方式来声明方法。  

example: 
```go
func(接收者) saySomething() {
    //...
}
```

最简单的继承：  
```go
type Rectangle struct {
    width, height float64
}

type Circle struct {
    radius float64
}

func (r Rectangle) area() float64 {
    return r.width*r.height
}

func (c Circle) area() float64 {
    return c.radius*c.radius*math.Pi
}
```  

1. method的名字一模一样，但是只要接收者的名字不一样，就是不同的method
  

通过接口封装方法：某个对象实现了interface中的的所有方法，对象就实现了该接口  
example:  
```go
type Human struct {
	name string
	age int
	phone string
}

type Student struct {
	Human //匿名字段  继承了Human
	school string
	loan float32
}

type Employee struct {
	Human //匿名字段  继承了Human
	company string
	money float32
}

//Human实现SayHi方法
func (h Human) SayHi() {
	fmt.Printf("Hi, I am %s you can call me on %s\n", h.name, h.phone)
}

//Human实现Sing方法
func (h Human) Sing(lyrics string) {
	fmt.Println("La la la la...", lyrics)
}

//Employee重载Human的SayHi方法
func (e Employee) SayHi() {
	fmt.Printf("Hi, I am %s, I work at %s. Call me on %s\n", e.name,
		e.company, e.phone)
	}

// Interface Men被Human,Student和Employee实现
// 因为这三个类型都实现了这两个方法
type Men interface {
	SayHi()
	Sing(lyrics string)
}

func main() {
	mike := Student{Human{"Mike", 25, "222-222-XXX"}, "MIT", 0.00}
	paul := Student{Human{"Paul", 26, "111-222-XXX"}, "Harvard", 100}
	sam := Employee{Human{"Sam", 36, "444-222-XXX"}, "Golang Inc.", 1000}
	tom := Employee{Human{"Tom", 37, "222-444-XXX"}, "Things Ltd.", 5000}

	//定义Men类型的变量i
	var i Men

	//i能存储Student
	i = mike
	fmt.Println("This is Mike, a Student:")
	i.SayHi()
	i.Sing("November rain")

	//i也能存储Employee
	i = tom
	fmt.Println("This is tom, an Employee:")
	i.SayHi()
	i.Sing("Born to be wild")

	//定义了slice Men
	fmt.Println("Let's use a slice of Men and see what happens")
	x := make([]Men, 3)
	//这三个都是不同类型的元素，但是他们实现了interface同一个接口
	x[0], x[1], x[2] = paul, sam, mike

	for _, value := range x{
		value.SayHi()
	}
}
```
**Output:**  
```shell
This is Mike, a Student:
Hi, I am Mike you can call me on 222-222-XXX
La la la la... November rain
This is tom, an Employee:
Hi, I am Tom, I work at Things Ltd.. Call me on 222-444-XXX
La la la la... Born to be wild
Let's use a slice of Men and see what happens
Hi, I am Paul you can call me on 111-222-XXX
Hi, I am Sam, I work at Golang Inc.. Call me on 444-222-XXX
Hi, I am Mike you can call me on 222-222-XXX
```



