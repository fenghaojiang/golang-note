---
title: WebSocket、HTTP/2与gRPC
date: 2021-2-1
---


## WebSocket  

WebSocket是一个双向通信协议，它在握手阶段采用HTTP/1.1协议(暂时不支持HTTP/2)  

握手过程如下:  

1. 客户端向服务端发起一个特殊的HTTP请求，其消息头如下:  

```message
GET /chat HTTP/1.1   //请求行
Host: server.example.com  
Upgrade: websocket  //required  
Connection: Upgrade //required
Sec-Websocket-Key: dGhlIHNhbXBsZSBub25jZQ== //required，一个 16bits 编码得到的 base64 串
Origin: http://example.com // 用于防止未认证的跨域脚本使用浏览器 websocket api 与服务端进行通信
Sec-WebSocket-Protocol: chat, superchat //optional, 子协议协商字段  
Sec-WebSocket-Version: 13  
```

2. 如果服务端支持该版本的WebSocket，会返回101响应，响应标头如下：  


```message
HTTP/1.1 101 Switching Protocols  // 状态行
Upgrade: websocket   // required
Connection: Upgrade  // required
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo= // required，加密后的 Sec-WebSocket-Key
Sec-WebSocket-Protocol: chat // 表明选择的子协议
```


握手完成后，接下来的 TCP 数据包就都是 WebSocket 协议的帧了。

可以看到，这里的握手不是 TCP 的握手，而是在 TCP 连接内部，从 HTTP/1.1 upgrade 到 WebSocket 的握手。

WebSocket 提供两种协议：不加密的 ws:// 和 加密的 wss://. 因为是用 HTTP 握手，它和 HTTP 使用同样的端口：ws 是 80（HTTP），wss 是 443（HTTPS）   



## gRPC协议  

gRPC是一个远程过程调用框架，默认使用protobuf3进行数据的高效序列化与service定义，使用HTTP/2进行数据传输  

目前gRPC主要被用在微服务通信中，但是因为其优越的性能，它也很契合游戏、IoT等需要高性能低延迟的场景  

其实光从协议先进程度上讲，gRPC基本上全面超越REST:  

1. 使用二进制进行数据序列化，比json更节约流量、序列化与反序列化也更快
2. protobuf3要求api被完全清晰的定义好，而REST api只能靠程序员自觉定义  
3. gRPC官方就支持从api定义生成代码，而REST api需要借助aopenapi-codegen等第三方工具  
4. 支持4种通信模式：一对一(unary)、客户端流、服务端流、双端流。更灵活   

只是目前gRPC对broswer的支持并不是很好，如果需要通过浏览器访问api，那么gRPC可能不是你的菜。  

如果你的产品只打算考虑面向APP等可控的客户端，可以考虑上gRPC  

对同时需要为浏览器和APP提供服务应用而言，也可以考虑APP使用gRPC协议，而浏览器使用API网关提供的HTTP接口，在API网关上进行HTTP-gRPC协议转换  













