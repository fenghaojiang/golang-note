---
title: RPC模型
date: 2020-11-5
---


# RPC模型  
RPC全程Remote Procedure Call远程过程调用，简单的理解是一个节点请求另一个节点提供服务  

通过学习我们了解到Socket和HTTP连接时采用的是类似“信息交换”的模式，即客户端发送一条信息到服务端，然后服务端都会返回一定的信息以表示响应。客户端和服务端之间约定了交互信息的格式，以便双方都能清楚解析交互所产生的信息。但是很多独立的应用并没有采用这种模式，而是采用类似常规的函数调用方式来完成想要的功能。  

RPC就是想实现函数调用模式的网络化。客户端就像调用本地函数一样，然后客户端把这些参数打包之后通过网络传递到服务端，服务端解包到处理过程中执行，然后执行的结果反馈给客户端。  

Remote Procedure Call远程过程调用协议，是一种网络从远程计算机程序上请求服务，而不需要了解底层网络技术的协议。它假定某些传输协议的存在，如TCP或UDP，以便为通信程序之间携带信息数据，通过它可以使函数调用模式网络化。在OSI网络通信模型中，RPC跨越了传输层和应用层。RPC使得开发包括网络分布式多程序在内的应用程序更加容易。  

一次客户端对服务器RPC的调用，大致有如下十步：  

+ 调用客户端句柄；执行传送参数
+ 调用本地系统内核发送网络消息
+ 消息传送到远程主机
+ 服务器句柄得到消息并取得参数
+ 执行远程过程
+ 执行的过程将结果返回服务器句柄
+ 服务器句柄返回结果，调用远程系统内核
+ 消息传回本地主机
+ 客户句柄由内核接收消息
+ 客户接收句柄返回的数据

---

## Go RPC
Go标准包中已经提供了对RPC的支持，并且支持三个级别的RPC：TCP、HTTP、JSONRPC。但Go的RPC包是独一无二的RPC，它和传统的RPC系统不同，它只支持Go开发的服务器与客户端之间的交互，因为在内部使用了Gob来编码。  

Go RPC函数必须满足以下条件才能被远程访问，不然会被忽略。详细的要求如下：  
+ 函数必须是导出的(首字母大写，相当于public)
+ 必须有两个导出类型的参数
+ 第一个参数是接收的参数，第二个参数是返回给客户端的参数，第二个参数必须是指针类型的
+ 函数还要有一个返回值error

举个例子，正确的RPC函数格式如下:  
```go
func (t *T) MethodName(argType T1, replyType *T2) error
```

PS: T、T1、T2类型必须能被``encoding/gob``包编解码  

任何的RPC都需要通过网络来传递数据，Go RPC可以利用HTTP和TCP来传递数据，利用HTTP的好处是可以直接复用net/http里面的一些函数。详细的例子请看下面的实现



HTTP RPC实现:  

```go
package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {

	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
```






