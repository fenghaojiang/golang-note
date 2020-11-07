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

![request](request.png)  

+ 响应报文

![response](response.png)  

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
![connect](connect.jpg)  
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

|字段|需要调用的方法或需要访问的字段|键值对的来源|内容类型|
|:---:|:---:|:---:|:---:|:---:|
|||URL/表单|URL编码/Multipart编码|
|Form|ParseForm|√/√|√/-|
|PostForm|Form|-|√|√/-|
|MultipartForm|ParseMultipartForm|-/√|-/√|
|FormValue|无|√/√|√/-
|PostFormValue|无|-/√|√/-|

<br>
<br>
<br>

## 为什么要以传值的方式将ResponseWriter传递给ServeHTTP  

ServeHTTP为什么要接受ResponseWriter接口和一个指向Request结构的指针作为参数呢？  
接受Request结构指针的原因很简单： 为了让服务器能够察觉到处理器对Request结构的修改，我们必须以传引用(pass by reference)而不是传值(pass by value)  
而ResponseWriter实际上就是response这个非导出结构的接口，ResponseWriter在使用response结构时，传递的也是只想response结构的指针，这也就是说，ResponseWriter是以传引用而不是传值的方式在使用response结构  
总结，实际上ServeHTTP函数两个参数传递的都是引用而不是值