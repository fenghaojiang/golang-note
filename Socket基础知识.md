---
title: Socket相关知识
date: 2020-11-9
---

# Socket编程  
现在的网络编程几乎都是用Socket进行编程的。  
进程与服务器进行通信，都是靠Socket来进行通信的。  

## 什么是Socket?  
Socket起源于Unix，Unix的基本哲学之一就是“一切皆文件”，都可以用“打开open->读写write/read->关闭close”模式来操作。  
Socket就是该模式的一个实现，网络的Socket数据传输就是一种特殊的I/O，Socket也是一种文件描述符。Socket也具有一个类似打开文件的函数调用: ``Socket()``, 该函数返回一个整形的Socket描述符，随后的建立、数据传输等操作都是通过Socket实现的。  

常用的Socket类型有两种：流式Socket（SOCK_STREAM）和数据报式Socket（SOCK_DGRAM）。流式是一种面向连接的Socket，针对于面向连接的TCP服务应用；数据报式Socket是一种无连接的Socket，对应于无连接的UDP服务应用  

<br>

## Socket如何通信
网络中的进程之间如何通过Socket通信呢？首要解决的问题是如何唯一标识一个进程，否则通信无从谈起！在本地可以通过进程PID来唯一标识一个进程，但是在网络中这是行不通的。其实TCP/IP协议族已经帮我们解决了这个问题，网络层的“ip地址”可以唯一标识网络中的主机，而传输层的“协议+端口”可以唯一标识主机中的应用程序（进程）。这样利用三元组（ip地址，协议，端口）就可以标识网络的进程了，网络中需要互相通信的进程，就可以利用这个标志在他们之间进行交互。

使用TCP/IP协议的应用程序通常采用应用编程接口：UNIX BSD的套接字（socket）和UNIX System V的TLI（已经被淘汰），来实现网络进程之间的通信。就目前而言，几乎所有的应用程序都是采用socket，而现在又是网络时代，网络中进程通信是无处不在，这就是为什么说“一切皆Socket”。  

<br>
<br>  

## Socket基础知识  

通过上面的介绍我们知道Socket有两种：TCP Socket和UDP Socket，TCP和UDP是协议，而要确定一个进程的需要三元组，需要IP地址和端口。  

## IPv4地址
目前的全球因特网所采用的协议族是TCP/IP协议。IP是TCP/IP协议中网络层的协议，是TCP/IP协议族的核心协议。目前主要采用的IP协议的版本号是4(简称为IPv4)，发展至今已经使用了30多年。

IPv4的地址位数为32位，也就是最多有2的32次方的网络设备可以联到Internet上。近十年来由于互联网的蓬勃发展，IP位址的需求量愈来愈大，使得IP位址的发放愈趋紧张，前一段时间，据报道IPV4的地址已经发放完毕，我们公司目前很多服务器的IP都是一个宝贵的资源。

地址格式类似这样：127.0.0.1 172.122.121.111

<br>
<br> 

## IPv6地址
IPv6是下一版本的互联网协议，也可以说是下一代互联网的协议，它是为了解决IPv4在实施过程中遇到的各种问题而被提出的，IPv6采用128位地址长度，几乎可以不受限制地提供地址。按保守方法估算IPv6实际可分配的地址，整个地球的每平方米面积上仍可分配1000多个地址。在IPv6的设计过程中除了一劳永逸地解决了地址短缺问题以外，还考虑了在IPv4中解决不好的其它问题，主要有端到端IP连接、服务质量（QoS）、安全性、多播、移动性、即插即用等。

地址格式类似这样：2002:c0e8:82e7:0:0:0:c0e8:82e7  

<br>
<br> 

## Go支持的IP类型
在Go的net包中，有IP的定义如下：  
```go
type IP []byte
```

net包中ParseIP(s string) IP函数会把一个IPv4或者IPv6的地址转化成IP类型，请看下面的例子:  

```go
package main
import (
    "net"
    "os"
    "fmt"
)
func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
        os.Exit(1)
    }
    name := os.Args[1]
    addr := net.ParseIP(name)
    if addr == nil {
        fmt.Println("Invalid address")
    } else {
        fmt.Println("The address is ", addr.String())
    }
    os.Exit(0)
}
```

## TCP Socket  
当我们知道如何通过网络端口访问一个服务时，那么我们能够做什么呢？  
作为客户端来说，我们可以通过向远端某台机器的的某个网络端口发送一个请求，然后得到在机器的此端口上监听的服务反馈的信息。  
作为服务端，我们需要把服务绑定到某个指定端口，并且在此端口上监听，当有客户端来访问时能够读取信息并且写入反馈信息。  

在Go语言的net包中有一个类型```````TCPConn```````，这个类型可以用来作为客户端和服务器端交互的通道，他有两个主要的函数：  
```go
func (c *TCPConn) Write(b []byte) (int, error)
func (c *TCPConn) Read(b []byte) (int,error)
```

``TCPConn``可以用在客户端和服务端来读写数据  
TCP地址信息：  
```go
type TCPAddr struct {
    IP   IP
    Port int
    Zone string
}
```

在Go语言中通过ResolveTCPAddr获取一个TCPAddr

```go
func ResolveTCPAddr(net, addr string) (*TCPAddr, os.Error)
```

net参数是"tcp4"、"tcp6"、"tcp"中的任意一个，分别表示TCP(IPv4-only), TCP(IPv6-only)或者TCP(IPv4, IPv6的任意一个)。
addr表示域名或者IP地址，例如"www.google.com:80" 或者"127.0.0.1:22"。  


## TCP Client  
Go中通过net包中的DialTCP函数来建立一个TCP连接，并返回一个TCPConn类型的对象，当连接建立时服务器端也创建一个同类型的对象，此时客户端和服务端通过各自拥有的TCPConn来进行数据交换。  
客户端通过TCPConn对象将请求信息发送到服务端，读取服务端响应的信息。服务端读取并解析来自客户端的请求，并返回应答信息，这个链接只有当任一端关闭了连接之后才失效，不然这连接可以一直使用。

建立连接的函数定义：  

```go
func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
```

+ network参数是"tcp4"、"tcp6"、"tcp"中的任意一个，分别表示TCP(IPv4-only)、TCP(IPv6-only)或者TCP(IPv4,IPv6的任意一个)
+ laddr表示本机地址(local address)，一般设置为nil
+ raddr表示远程的服务地址(remote address)  

接下来我们写一个简单的例子，模拟一个基于HTTP协议的客户端请求去连接一个Web服务端。我们要写一个简单的http请求头，格式类似如下：  

```shell
"HEAD / HTTP/1.0\r\n\r\n"
```

从服务端接收到的响应信息格式可能如下：  
```HTTP
HTTP/1.0 200 OK
ETag: "-9985996"
Last-Modified: Thu, 25 Mar 2010 17:51:10 GMT
Content-Length: 18074
Connection: close
Date: Sat, 28 Aug 2010 00:43:48 GMT
Server: lighttpd/1.4.23
```
客户端代码:  
```go
package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	checkError(err)
	// result, err := ioutil.ReadAll(conn)
	result := make([]byte, 256)
	_, err = conn.Read(result)
	checkError(err)
	fmt.Println(string(result))
	os.Exit(0)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
```

通过上面的代码我们可以看出：首先程序将用户的输入作为参数service传入net.ResolveTCPAddr获取一个tcpAddr,然后把tcpAddr传入DialTCP后创建了一个TCP连接conn，通过conn来发送请求信息，最后通过ioutil.ReadAll从conn中读取全部的文本，也就是服务端响应反馈的信息。  


## TCP Server  
当我们编写了一个TCP的客户端程序，也可以通过net包来创建一个服务器程序，在服务端我们需要绑定服务到指定的非激活端口，并监听此端口，当有客户端到达的时候可以接收到来自客户端连接的请求。net包中有相应的功能函数，函数：  
```go
func ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)
func (l *TCPListener) Accept() (Conn, error)
```

下面我们实现一个简单的时间同步服务，监听7777端口:  
```go
package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
    service := "7777"
    tcpAddr, err := net.ResolveTCPAddr("tcp4",service)
    checkError(err)
    listener, err := net.ListenTCP("tcp")
    checkError(err)
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        daytime := time.Now().String()
        conn.Write([]byte(daytime))
        conn.Close()
    }
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error :%s", err.Error())
        os.Exit(1)
    }
}
```

上面的服务跑起来之后，它将会一直在那里等待，直到有新的客户端请求到达。当有新的客户端请求到达并接受同意接受Accept该请求的时候他会反馈当前的时间信息。值得注意的是，在代码中for循环里面，当有错误发生时，直接continue而不是退出，是因为在服务器跑代码的时候，当有错误发生的情况下最好是有服务端记录错误，然后当前连接的客户端直接报错而退出，从而不会影响到当前服务端运行的整个服务。  

<br>
上面的代码有个缺点，执行的时候是单任务的，不能接受多个请求，用goroutine机制改善：  


```go
package main 

import (
    "fmt"
    "net"
    "os"
    "time"
)

func main() {
    service := "1200"
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkError(err)
    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()
    daytime := time.Now().String()
    conn.Write([]byte(daytime))
} 

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error %s", err.Error())
        os.Exit(1)
    }
}
```

通过go就能实现并发，是不是都给他懂完了？  
有一位细心的读者DannyD问：服务端都没有处理客户端的实际请求的内容，这是个锤子。  

问的好！如果我们需要从客户端发送不同的请求来获取不同的时间格式并且需要一个长连接。  

```go
package main

import (
    "fmt"
	"net"
	"os"
	"time"
	"strconv"
	"strings"
)

func main() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
	request := make([]byte, 128) // set maxium request length to 128B to prevent flood attack
	defer conn.Close()  // close connection before exit
	for {
		read_len, err := conn.Read(request)

		if err != nil {
			fmt.Println(err)
			break
		}

    		if read_len == 0 {
    			break // connection already closed by client
    		} else if strings.TrimSpace(string(request[:read_len])) == "timestamp" {
    			daytime := strconv.FormatInt(time.Now().Unix(), 10)
    			conn.Write([]byte(daytime))
    		} else {
    			daytime := time.Now().String()
    			conn.Write([]byte(daytime))
    		}

    		request = make([]byte, 128) // clear last read content
	}
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
```


在上面这个例子中，我们使用conn.Read()不断读取客户端发来的请求。由于我们需要保持与客户端的长连接，所以不能在读取完一次请求后就关闭连接。由于conn.SetReadDeadline()设置了超时，当一定时间内客户端无请求发送，conn便会自动关闭，下面的for循环即会因为连接已关闭而跳出。需要注意的是，request在创建时需要指定一个最大长度以防止flood attack；每次读取到请求处理完毕后，需要清理request，因为conn.Read()会将新读取到的内容append到原内容之后。


## 控制TCP连接  

TCP有很多连接控制函数，我们平常用到比较多的有如下几个函数：  
```go
func DialTimeout(net, addr string, timeout time.Duration) (Conn, error)
```
设置建立连接的超时时间，客户端和服务器端都适用，当超过设置时间时，连接自动关闭。

```go
func (c *TCPConn) SetReadDeadline(t time.Time) error
func (c *TCPConn) SetWriteDeadline(t time.Time) error
```

用来设置写入/读取一个连接的超时时间。当超过设置时间时，连接自动关闭。  
```go
func (c *TCPConn) SetKeepAlive(keepalive bool) os.Error
```
设置keepAlive属性。操作系统层在tcp上没有数据和ACK的时候，会间隔性的发送keepalive包，操作系统可以通过该包来判断一个tcp连接是否已经断开，在windows上默认2个小时没有收到数据和keepalive包的时候认为tcp连接已经断开，这个功能和我们通常在应用层加的心跳包的功能类似。

更多的内容请查看net包的文档。


# UDP Socket

Go语言包中处理UDP Socket和TCP Socket不同的地方就是在服务器端处理多个客户端请求数据包的方式不同,UDP缺少了对客户端连接请求的Accept函数。其他基本几乎一模一样，只有TCP换成了UDP而已。UDP的几个主要函数如下所示：  

```go
func ResolveUDPAddr(net, addr string) (*UDPAddr, os.Error)
func DialUDP(net string, laddr, raddr *UDPAddr) (c *UDPConn, err os.Error)
func ListenUDP(net string, laddr *UDPAddr) (c *UDPConn, err os.Error)
func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err os.Error)
func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (n int, err os.Error)
```

UDP的客户端代码如下所示,我们可以看到不同的就是TCP换成了UDP而已：  

```go
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)
	_, err = conn.Write([]byte("anything"))
	checkError(err)
	var buf [512]byte
	n, err := conn.Read(buf[0:])
	checkError(err)
	fmt.Println(string(buf[0:n]))
	os.Exit(0)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s", err.Error())
		os.Exit(1)
	}
}
```

UDP服务端如何来处理：  
```go
package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	for {
		handleClient(conn)
	}
}
func handleClient(conn *net.UDPConn) {
	var buf [512]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	daytime := time.Now().String()
	conn.WriteToUDP([]byte(daytime), addr)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s", err.Error())
		os.Exit(1)
	}
}
```
