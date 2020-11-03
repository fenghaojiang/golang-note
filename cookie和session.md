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
