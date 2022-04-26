## Go泛型

[demo](./demo/generics/main.go)


Interfaces used as constraints may be given names (such as Ordered), or they may be literal interfaces inlined in a type parameter list. For example:

```go
[S interface{~[]E}, E interface{}]
```
Here S must be a slice type whose element type can be any type.

Because this is a common case, the enclosing interface{} may be omitted for interfaces in constraint position, and we can simply write:

```go
[S ~[]E, E interface{}]
```


链接： 
+ https://go.dev/blog/intro-generics
+ https://go.dev/doc/tutorial/generics
+ https://tonybai.com/2022/04/20/some-changes-in-go-1-18/
+ 

