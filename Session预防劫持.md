---
title: Session预防劫持
date: 2020-11-4
---


# Session预防劫持  

session劫持是一种存在比较严重的安全威胁，session标识符在客户端与服务端的会话中很容易被嗅探到  


## 劫持例子  
example:  
```go
func count(w http.ResponseWriter, r *http.Request) {
    sess := globalSessions.SessionStart(w, r)
    ct := sess.Get("countnum")
    if ct == nil {
        sess.Set("countnum", 1)
    } else {
        sess.Set("countnum",(ct.(int) + 1))
    }
    t, _ := template.ParseFiles("count.gtpl")
    w.Header().Set("Content-Type", "text/html")
    t.Execute(w, sess.Get("countnum"))
}
```

count.gtpl:  
```go
Hi. Now count:{{.}}
```

执行上述代码，不停刷新，数字不断增长。  
但复制链接，用另一个firefox浏览器打开，用cookie模拟插件新建一个cookie复制一份。  
从Chrome浏览器换了firefox浏览器但是获得了sessionID，模拟了cookie的存储的过程。就算是用不同一台计算机来做结果也会完全一样，当你交替点击两个不同的浏览器的链接你会发现操控的是同一个计数器。  
firefox盗用了Chrome和goserver之间的维持会话的钥匙，即gosessionid，这是一种类型的“会话劫持”。  

由于http协议的无状态性，goserver无法得知sessionId是firefox从Chrome哪里“劫持”来的，它依然会去查找对应的session，并执行相关计算。  
与此同时，Chrome也没有办法得知自己保持的会话已经被“劫持”  

很喜欢牛皮钊钊的一句话，“那咋办嘛？”  

## session劫持防范  
**cookieonly和token**  

1. SessionId只允许cookie设置，而不是通过URL重置方式设置。同时设置cookie的httponly为true，这个属性是设置是否可通过客户端脚本访问这个设置的cookie，第一可以防止cookie被XSS读取从而引起session劫持，第二cookie设置不会像URL重置方式那么容易获取sessionID  

2.  每个请求都加上token，实现之前写的防止form重复递交类似的功能，在每个请求里面加上一个隐藏的token，然后每次验证这个token，保证用户请求都是唯一的。  

```go
h := md5.New()
salt := "fenghaojiang%635509290"
io.WriteString(h, salt+time.Now().String())
token := fmt.Sprintf("%x", h.Sum(nil))
if r.Form["token"] != token {
    //提示登陆
}
sess.Set("token",token)
```  

**间隔生成新的SID**  
还有一个解决方案就是，我们给session额外设置一个创建时间的值，一旦过了一定的时间就销毁sessionID重新生成新的session  

```go
createtime := sess.Get("createtime")
if createtime == nil {
	sess.Set("createtime", time.Now().Unix())
} else if (createtime.(int64) + 60) < (time.Now().Unix()) {
	globalSessions.SessionDestroy(w, r)
	sess = globalSessions.SessionStart(w, r)
}
```

初始赋值，每次请求与上一次比较看看是否过期，60s后重新生成新的ID  
上述两种手段组合可以在实践中消除session劫持的风险，另一方面，由于sessionID频繁改变，使得攻击者很难有机会获取有效的sessionID，另一方面，因为sessionID只能在cookie中传递，然后设置了httponly，基于URL攻击的可能性为0，同时被XSS获取sessionID也不可能。最后设置MaxAge=0，这样就相当于session cookie不会留在浏览器的历史记录里面。