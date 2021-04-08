---
title: gin相关
date: 2021-4-8
---


+ gin.H{...}: 就是一个map[string]interface{}
+ gin.Context: Context是gin的上下文，它允许我们在中间件之间传递变量、管理流、验证JSON请求、响应JSON请求，在gin中包含大量Context的方法，例如我们常用的```DefaultQuery```, ```Query```, ```DefaultPostForm```, ```Form```



1. http.Server:
```go
type Server struct {
    Addr    string
    Handler Handler
    TLSConfig *tls.Config
    ReadTimeout time.Duration
    ReadHeaderTimeout time.Duration
    WriteTimeout time.Duration
    IdleTimeout time.Duration
    MaxHeaderBytes int
    ConnState func(net.Conn, ConnState)
    ErrorLog *log.Logger
}
```

+ Addr: 监听的TCP地址，格式为:8000
+ Handler: http句柄，实质为ServeHTTP，用于处理程序响应HTTP请求
+ TLSConfig: 安全传输层协议(TLS)的配置
+ ReadTimeout: 允许读取的最大时间
+ ReadHeaderTimeout: 允许读取请求头的最大时间
+ WriterTimeout: 允许写入的最大时间
+ IdleTimeout: 等待的最大时间
+ MaxHeaderBytes: 请求头的最大字节数
+ ConnState: 指定一个可选的回调函数，当客户端连接发生变化时调用
+ ErrorLog: 指定一个可选的日志记录器, 用于接收程序的意外行为和底层系统错误 如果未设置或为nil则默认以日志包的标准日志记录器完成  


2. ListenAndServe:

```go
func (srv *Server) ListenAndServe() error {
    addr := srv.Addr
    if addr == "" {
        addr := ":http"
    }
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}
```

开始监听服务，监听TCP网络地址，Addr和调用应用程序处理连接上的请求


3. http.ListenAndServe 和r.Run()有区别吗?

```go
func (engine *Engine) Run(addr ...string) (err error) {
    defer func() { debugPrintError(err) }()

    address := resolveAddress(addr)
    debugPrint("Listening and serving HTTP on %s\n", address)
    err = http.ListenAndServe(address, engine)
    return
}
```

没区别


