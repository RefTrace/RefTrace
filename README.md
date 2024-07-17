# RefTrace (Go Version)

```
go generate ./...
go build
```

https://issues.apache.org/jira/browse/GROOVY-9232

https://go.dev/doc/effective_go#embedding
https://eli.thegreenplace.net/2020/embedding-in-go-part-3-interfaces-in-structs
https://gobyexample.com/struct-embedding
https://preslav.me/2023/02/06/golang-do-we-need-struct-pointers-everywhere/

```
type T struct{}
var _ I = T{}       // Verify that T implements I.
var _ I = (*T)(nil) // Verify that *T implements I.
```
