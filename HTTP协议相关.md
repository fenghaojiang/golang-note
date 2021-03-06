---
title: HTTP网络协议相关知识笔记
date: 2020-11-2
---


# HTTP网络协议相关

## URI  
URI包含URL跟URN  
URI全称Uniform Resource Identifier 统一资源标识符  
URL全称Uniform Resource Locator 统一资源定位符  
URN全称Uniform Resource Name 统一资源名称  
看这中文名跟英文名就知道URI包含URL跟URN，懂的都懂，真的是给他懂完了。  


## HTTP URL  
格式：```http://host[:port][abs_path]```    
host: 合法的Internet主机域名或者IP地址  
port: 指定的一个端口号，拥有被请求资源的服务器主机监听该端口的TCP连接  
abs_path指定请求资源的URI,如果URL中没有给出abs_path，那么当它作为请求URI时，必须以"/"的形式给出。通常这个工作浏览器已经完成  

URL结构  
```go
type URL struct {
    Scheme string
    Opaque string
    User   *UserInfo
    Host   string
    Path   string
    RawQuery string
    Fragment string
}
```
scheme://[userinfo@]host/path[?query][#fragment]  
那些在scheme之后不带斜线的URL则会被解释为：  
scheme:opaque[?query][#fragment]

## 请求报文和响应报文  
请求报文的第一行是请求行，包含了方法字段。  
+ 请求报文 

![request](./img/request.png)  

一个标准的HTTP请求由以下部分
+ 请求行
+ 请求头  
+ 请求体  

### 请求行  

请求行以一个方法符号开头，以空格分开，后面跟着请求的URI和协议的版本，格式如下：  

```shell
Method Request-URI HTTP-Version CRLF
```
*Method*表示请求方法  
*Request-URI*是一个统一资源标识符  
*HTTP-Version*表示请求的HTTP协议版本  
*CRLF*表示回车和换行(除了作为结尾的CRLF外，不允许出现单独的CR或LF字符)  

### 请求头  
HTTP消息报头包括: 普通报头、请求报头、响应报头、实体报头，每一个报头域都是由  
**名字+ : + 空格 + 值**组成，消息报头域的名字是大小写无关的  

#### 普通报头  
在普通报头中，有少数报头域用于所有的请求和响应消息，但并不用于被传输的实体，只用于传输的消息  
| 字段名 | 说明 |
|---|---|
|Cache-Control|控制缓存行为|
|Connection|连接的管理|
|Date|普通报头域表示消息产生的日期和时间|
|Pragma|http 1.0中的保温指令控制|


#### 请求报头  
请求报头允许客户端向服务端传递请求的附加信息以及客户端自身的信息。常见的请求报头包括:  

|字段名|说明|
|---|---|
|Accept|客户端可处理的媒体类型：Accept：image/gif|
|Accept-Charset|客户端可处理的字符集|
|Accept-Encoding|客户端的编码方式|
|Accept-Langulage|客户端指定的语言类型|
|Authrization|web认证信息|
|Expect|期待服务器的特定行为|
|Host|请求报头域主要用于指定被请求资源的 Internet 主机和端口号|
|User-Agent|请求报头域允许客户端将它的操作系统、浏览器和其它属性|
|Referer|请求中的 url 上一跳地址|

#### 响应报头  

响应报头允许服务器传递不能放在状态行中的附加响应信息，以及关于服务器的信息和对Request-URI所标识的资源进行下一步访问的信息。  

常见的响应报头包括：  
|字段名|说明|
|---|---|
|Age|资源的创建时间|
|Location|	客户端重定向至指定的URL|
|Retry-After|再次发送请求的时机|
|www-Authenticate|服务器对客户端的认证|

#### 实体报头  
请求和响应消息都可以传送一个实体。一个实体由实体报头域和实体正文组成，但并不是说实体报头域和实体正文要在一起发送，可以只发送实体报头域。实体报头定义了关于实体正文和请求所标识的资源的元信息。  

|字段名|说明|
|---|---|
|Allow|资源所支持的HTTP请求类型|
|Content-Encoding|数据编码方式|
|Content-Language|数据的语言类型|
|Content-Length|实体的内容大小|
|Content-Location|替代对应资源的URI|
|Content-Type|实体报头域用语指明发送给接收者的实体正文的媒体类型|
|Expires|数据过期时间|
|Last-Modified|资源的最后修改时间|


+ 响应报文

![response](./img/response.png)  

## HTTP方法简单介绍  
请求方法字段  
+ GET  
请求指定的页面信息，并返回实体主体
+ HEAD  
与GET请求类似，但不返回报文实体主体部分。  
主要用于确认URL的有效性以及资源的更新的日期时间等操作。
+ POST  
向指定资源提交数据进行处理请求（例如提交表单或者上传文件）。数据被包含在请求体中，POST请求可能会导致新的资源的建立和/或已有资源的修改。  
POST主要用来传输数据，而GET主要用来获取资源。
+ PUT  
从客户端向服务器传送的数据取代指定的文档的内容，自身并不带验证的机制，任何人都可以上传文件，存在安全性问题，一般不使用该方法。
+ DELETE  
请求服务器删除指定的页面，与PUT请求相反，同样不带验证机制。
+ PATCH  
对PUT方法的补充，对已知资源的进行局部更新。PUT只能完全替代，PATCH允许部分修改。
+ OPTIONS  
查询指定URL能够支持的方法  
返回Allow: GET, POST, HEAD, OPTIONS这样的内容
+ CONNECT  
HTTP/1.1 协议中预留给能够将连接改为管道方式的代理服务器。要求与代理服务器通信时建立隧道。  
使用SSL（Secure Socket Layer，安全套接层）和TLS（Transport Layer Security，传输层安全）协议把通信内容加密后经网络隧道传输。  
![connect](./img/connect.jpg)  
CONNECT www.example.com:443 HTTP/1.1  
看图解有点像webRtc?  

+ TRACE  
回显服务器收到的请求，主要用于测试或诊断。  
服务器把通信路径返回给客户端，通常不会用到TRACE，并且容易受到XST攻击(Cross-Site Tracing，跨站追踪)  

  
## GET和POST的区别  
  
GET请求——从指定的URL中获取资源与数据
+ GET请求可以被缓存，可以保留在浏览器的历史记录中，可被收藏为书签  
+ GET请求时的结果幂等的（多次进行该运算的结果跟第一次进行的结果是一样的，就称为幂等）
+ GET有长度限制，在HTTP协议的定义中，没有对GET请求的大小进行限制，不过因为浏览器不一样，所以一般限制在2~8K  
+ GET请求的所有参数都被包装在URL中，服务器的访问日志会记录，不要传递敏感信息。  

POST请求——向指定的资源提交要被处理的数据
+ POST请求不会被缓存，不会保留在浏览器的历史历史记录中，不能被收藏为书签。
+ POST向服务器中发送数据，也可以获得服务器处理之后的结果，效率不如GET
+ POST提交的数据比较大，大小由服务器的设定值限制，PHP通常限定2M
+ POST提交的参数包装成二进制的数据体，格式与GET基本一致
+ URL中只有资源路径，但是不包含参数，服务器日志不会记录，相对更安全
+ 涉及用户隐私的数据（密码、银行卡号、身份证）一定要用POST方式传递。

## HTTP状态响应码

服务器返回的   **响应报文**   中第一行为状态行，包含了状态码以及原因短语，用来告知客户端请求的结果。

| 状态码 | 类别 | 含义 |
| :---: | :---: | :---: |
| 1XX | Informational（信息性状态码） | 接收的请求正在处理 |
| 2XX | Success（成功状态码） | 请求正常处理完毕 |
| 3XX | Redirection（重定向状态码） | 需要进行附加操作以完成请求 |
| 4XX | Client Error（客户端错误状态码） | 服务器无法处理请求 |
| 5XX | Server Error（服务器错误状态码） | 服务器处理请求出错 |

**1XX 信息**

-   **100 Continue**  ：表明到目前为止都很正常，客户端可以继续发送请求或者忽略这个响应。

**2XX 成功**

-   **200 OK**  

-   **204 No Content**  ：请求已经成功处理，但是返回的响应报文不包含实体的主体部分。一般在只需要从客户端往服务器发送信息，而不需要返回数据时使用。

-   **206 Partial Content**  ：表示客户端进行了范围请求，响应报文包含由 Content-Range 指定范围的实体内容。

**3XX 重定向**

-   **301 Moved Permanently**  ：永久性重定向

-   **302 Found**  ：临时性重定向

-   **303 See Other**  ：和 302 有着相同的功能，但是 303 明确要求客户端应该采用 GET 方法获取资源。

- 注：虽然 HTTP 协议规定 301、302 状态下重定向时不允许把 POST 方法改成 GET 方法，但是大多数浏览器都会在 301、302 和 303 状态下的重定向把 POST 方法改成 GET 方法。

-   **304 Not Modified**  ：如果请求报文首部包含一些条件，例如：If-Match，If-Modified-Since，If-None-Match，If-Range，If-Unmodified-Since，如果不满足条件，则服务器会返回 304 状态码。

-   **307 Temporary Redirect**  ：临时重定向，与 302 的含义类似，但是 307 要求浏览器不会把重定向请求的 POST 方法改成 GET 方法。

**4XX 客户端错误**

-   **400 Bad Request**  ：请求报文中存在语法错误。

-   **401 Unauthorized**  ：该状态码表示发送的请求需要有认证信息（BASIC 认证、DIGEST 认证）。如果之前已进行过一次请求，则表示用户认证失败。

-   **403 Forbidden**  ：请求被拒绝。

-   **404 Not Found**  

**5XX 服务器错误**

-   **500 Internal Server Error**  ：服务器正在执行请求时发生错误。

-   **503 Service Unavailable**  ：服务器暂时处于超负载或正在进行停机维护，现在无法处理请求。  


## Cookie  
HTTP协议是无状态的，为了让HTTP协议尽可能的简单并能够支持大量的事务，HTTP/1.1引入Cookie来保存状态信息。  
Cookie是服务器发送到用户浏览器并保存在本地的一小块数据，它会在浏览器之后向同一服务器再次发起请求时被携带上，用于告知服务器两个请求是否来自同一个浏览器，由于之后每次请求都会需要携带Cookie数据，因此也带来额外的性能开销。  
Cookie曾一度用于客户端数据的存储，因为当时并没有其他合适存储方式，现在各种现代的浏览器开始支持其他的存储方式，Cookie逐渐被淘汰，新的浏览器API已经允许开发者直接将数据存储到本地，如使用Web storage API或indexedDB  
Cookie的设计本意是要客服HTTP的无状态性，虽然cookie并不是完成这一目的的唯一方法。  

## Cookie在go语言中的结构  

```go
type Cookie struct {
    Name       string
    Value      string
    Path       string
    Domain     string
    Expires    time.Time
    RawExpires string
    MaxAge     int
    Secure     bool
    HttpOnly   bool
    Raw        string
    Unparsed   []string
}
```
<br>
没有设置Expires字段的cookie通常称为会话cookie或者临时cookie，这种cookie在浏览器关闭就会被移除。  
相对而言，设置了Expires字段的cookie通常称为持久cookie，这种cookie会一直存在，直到指定的过期时间或者被手动删除为止。  
Expires字段用于明确地指定cookie应该在什么时候过期，而MaxAge字段则指明了cookie在被浏览器创建出来能存活多少秒。  
之所以会出现两种方式完全是因为不同浏览器使用了各不相同的cookie机制，跟Go语言本身的设计没有关系。  
<br>
HTTP1.1中废弃了Expires，推荐使用MaxAge，但是几乎所有的浏览器都支持Expires。微软的IE6、IE7、IE8都不支持MaxAge  
为了让cookie在所有浏览器上都能正常运行，一个实际的方法是只使用Expires，或者同时使用MaxAge跟Expires  



**Cookie主要用途**  
* 会话状态管理（登陆状态、购物车、游戏分数或者其他需要记录的信息）
* 个性化设置（用户自定义设置、主题等）
* 浏览器行为跟踪（跟踪分析用户行为）  

**Cookie创建过程**  
+ 服务端发送的响应报文包含Set-Cookie的首部字段，客户端得到响应报文后把Cookie内容保存到浏览器中。  
+ 客户端之后对同一个服务器发送请求时，会从浏览器中取出Cookie信息并通过Cookie请求首部字段发送给服务器。

## Session  
除了可以将用户信息通过Cookie存储在用户浏览器，也可以通过Session存储在服务器中，存储在服务器更加安全。  
Session可以存储在服务器上的文件、数据库、或者内存中。也可以将Session存储在Redis这种内存型数据库中，效率会更高。  
  
**Session维护用户登录状态过程**
* 用户进行登录时，用户提交包含用户名和密码的表单，放入 HTTP 请求报文中  
* 服务器验证该用户名和密码，如果正确则把用户信息存储到 Redis 中，它在 Redis 中的 Key 称为 Session ID；  
* 服务器返回的响应报文的 Set-Cookie 首部字段包含了这个 Session ID，客户端收到响应报文之后将该 Cookie 值存入浏览器中；
* 客户端之后对同一个服务器进行请求时会包含该 Cookie 值，服务器收到之后提取出 Session ID，从 Redis 中取出用户信息，继续之前的业务操作。

从上述描述中就可以简单得知，SessionKey被别人猜出就相当于免去登陆的过程，所以不能生成一个被别人轻易猜到的Session ID，与此同时还要经常重新生成Session ID。  
除了使用Session管理用户状态之外，还需要对用户进行重新验证，比如重新输入密码，或者使用短信验证码的方式验证。  

## Cookie与Session选择  
+ Cookie只能存储ASCII码字符串，而Session则可以存储任何类型的数据，考虑数据的复杂性时首选Session  
+ Cookie存储在浏览器中，容易被恶意查看，可以将Cookie值加密，在服务器端再进行解密。
+ 一般不会将用户的所有信息放在Session中，如果把所有信息放在Session中，开销将会非常巨大。

## ServeMux

+ 程序绑定的URL如果不是以/结尾，那么它只会与完全相同的URL匹配；
+ 但如果绑定的URL以/结尾，即使URL只有前缀部分与被绑定的URL相同，ServerMux也会认为两个URL匹配，比如绑定/hello/，浏览器请求/hello/there的时候，服务器在找不到与之完全匹配的处理器时，就会退而求其次，找到与/hello/匹配的处理器。


## Request字段的对比  
**对比Form、PostForm和MultipartForm**  


| 字段 | 需要调用的方法或需要访问的字段 | 键值对的来源 | 内容类型 |
| :---: | :---: | :---: | :---: | :---: |
| - | - |URL/表单|URL编码/Multipart编码|
|Form|ParseForm|√/√|√/-|
|PostForm|Form|-|√|√/-|
|MultipartForm|ParseMultipartForm|-/√|-/√|
|FormValue|无|√/√|√/-|
|PostFormValue|无|-/√|√/-|

<br>
<br>
<br>

## 为什么要以传值的方式将ResponseWriter传递给ServeHTTP  

ServeHTTP为什么要接受ResponseWriter接口和一个指向Request结构的指针作为参数呢？  
接受Request结构指针的原因很简单： 为了让服务器能够察觉到处理器对Request结构的修改，我们必须以传引用(pass by reference)而不是传值(pass by value)  
而ResponseWriter实际上就是response这个非导出结构的接口，ResponseWriter在使用response结构时，传递的也是只想response结构的指针，这也就是说，ResponseWriter是以传引用而不是传值的方式在使用response结构  
总结，实际上ServeHTTP函数两个参数传递的都是引用而不是值


