# Cookie和Session


## Go设置Cookie  

Go语言中通过net/http包中的SetCookie来设置：  

```go
http.SetCookie(w ResponseWrite, cookie *Cookie)
```

Cookie对象:  

```go
type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string

// MaxAge=0 means no 'Max-Age' attribute specified.
// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}
```

example：  

```go
expiration := time.Now()
expiration = expiration.AddDate(1, 0, 0)
cookie := http.Cookie{Name:"username", Value:"fenghaojiang", Expires:expiration}
http.SetCookie(w, &cookie)
```

## Go读取Cookie

```go
cookie , _ := r.Cookie("username")
fmt.Fprintf(w, cookie)
```

## Session创建过程详细描述  

session的基本原理是有服务器为每一个会话维护一份信息数据，客户端和服务端依靠一个全局唯一的标识来访问这份数据，以达到交互的目的。当用户访问Web应用时，服务端程序会随需要创建session，这个过程可以概括为三个步骤。 

+ 生成全局唯一的标识符(sessionId)  

+ 开辟数据存储空间。一般会在内存中创建相应的数据结构，但这种情况下，系统一旦掉电、宕机，所有的会话数据就会丢失，如果这种情况发生在淘宝之类的电子商务网站，后果不堪设想。  
为了解决这种问题，可以将会话数据写到文件里或存储在数据库中，当然这样会增加I/O开销，但是它可以实现某种成都的session持久化，也更有利于session的共享。  

+ 将session的全局唯一标识符发给客户端

以上三个步骤，最关键的是如何发送这个session的唯一标识这一步上，考虑到http协议的定义，数据无非可以放到请求行，头域或Body里，所以一般来说会有两种常用的方式：cookie和重写

1. Cookie服务端通过设置Set-cookie头就可以将session标识符传给客户端，而此后客户端的每一次请求都会带上这个标识符，另外一般包含session信息的cookie会将失效时间设置为0(会话cookie)，即浏览器进程有效时间。至于浏览器怎么处理这个0，每个浏览器都有自己的方案。

2. URL重写就是在返回给用户的页面里的所有URL后面追加session标识符，这样用户在收到响应之后，无论点击响应页面里的哪个链接或者提交表单，都会自动带上session标识符，从而实现会话的保持。虽然比较麻烦但是如果客户端禁用了cookie，这种方案将会是首选。



Go标准包没有为session提供任何支持，自己动手，丰衣足食咯。

## Session管理设计  

一般session管理涉及以下因素：

+ 全局session管理器
+ 保证sessionId的全局唯一性
+ 为每个客户关联一个session
+ session的存储(内存、文件、数据库)
+ session过期处理


管理器：

```go
type Manager struct {
    cookieName string //private
    lock sync.Mutex // protects session
    provider Provider
    maxLifeTime int64
}

func NewManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
    provider, ok := provides[providerName]
    if !ok {
        return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
        return &Manager{provider: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}, nil
    }
}
```

Provider接口:  
```go
type Provider interface {
    //Session初始化，操作成果则返回新的Session变量
    SessionInit(sid string) (Session, error)

    //返回sid所代表的Session变量，不存在则用sid调用SessionInit创建并返回
    SessionRead(sid string) (Session, error)

    //销毁Session
    SessionDestroy(sid string) error

    //根据maxLifeTime来删除过期数据
    SessionGC(maxLifeTime int64)
}
```

在main包中创建全局session管理器
```go
var globalSessions *session.Manager
func init() {
    globalSessions, _ := NewManager("memory","gosessionid",3600)
}
```

Session接口：  
设置值、读取值、删除值以及获取当前的sessionID这四个操作

```go
type Session interface {
    Set(key, value interface{}) error
    Get(key interface{}) interface{}
    Delete(key interface{}) error
    SessionID() string
}
```

注册实现：

```go
var provides = make(map[string]Provider)

// Register makes a session provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.

func Register(name string, provider Provider) {
    if provider == nil {
        panic("session: Register provider is nil")
    }
    if _, dup := provides[name];dup {
        panic("session: Register called twice for provider " + name)
    }
    provider[name] = provider
}
```

全局唯一的SessionId  
必须保证它是全局唯一的(GUID)  

```go
func (manager *Manager) sessionId() string {
    b := make([]byte, 32)
    //rand.Reader是一个全局、共享的密码用强随机数生成器
    if _, err := rand.Read(b); err != nil {
        return ""
    }
    return base64.URIEncoding.EncodeToString(b) //全局唯一随机数b编码后返回字符串作为sessionId
}
```

**Session创建**  
为每个来访用户分配或者获取与他相关联的Session，以便后面根据Session信息来验证操作。  
SessionStart函数：用于检测是否已经有某个Session与当前来访用户发生关联，没有则创建之。  

```go
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
    manager.lock.Lock()
    defer manager.lock.Unlock()
    cookie, err := r.Cookie(manager.cookieName)
    //读取失败或者首次创建为空
    if err != nil || cookie.Value == "" {
        sid := manager.sessionId()
        session, _ = manager.provider.SessionInit(sid)
        cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxLifeTime)}
        http.SetCookie(w, &cookie)
    } else {
        sid, _ := url.QueryUnescape(cookie.Value)
        session, _ = manager.provider.SessionRead(sid)
    }
    return 
}
```

session应用：  
```go
func login(w http.ResponseWriter, r *http.Request) {
    sess := globalSessions.SessionStart(w, r)
    r.ParseForm()
    if r.Method == "GET" {
        t, _ := template.ParseFiles("login.gtpl")
        w.Header().Set("Content-Type","text/html")
        t.Execute(w, sess.Get("username"))
    } else {
        sess.Set("username", r.Form["username"])
        http.Redirect(w, r, "/", 302)
    }
}
```

操作值：设置、读取与删除;上面代码已经实现了基本的读取数据操作，现在让我们来看看详细的操作。

```go
func count(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	createtime := sess.Get("createtime")
	if createtime == nil {
		sess.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
		globalSessions.SessionDestroy(w, r)
		sess = globalSessions.SessionStart(w, r)
	}
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	t, _ := template.ParseFiles("count.gtpl")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("countnum"))
}
```


**Session重置**  
Web应用有用户退出就需要对用户的session数据进行销毁操作。  
上面的代码已经演示了如何使用session重置操作。  

```go
//Destroy sessionid
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
        //把过期时间设置为现在
        expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}
```

**Session的销毁**  
看到Session的销毁gc时，感觉真的是懂的都懂，都给他懂完了都，我佛了  

在main启动的时候启动:  
```go
func init() {
    go globalSession.GC()
}

func (manager *Manager) GC() {
    manager.lock.Lock()
    defer manager.lock.Unlock()
    manager.provider.SessionGC(manager.maxLifeTime)
    time.AfterFunc(time.Duration(manager.maxLifeTime), func() { manager.GC() })
}

```

