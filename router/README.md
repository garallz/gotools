# Golang Router (version: 0.1)

## Base on Golang net/http

### deal with response body to json and xml

### Make request uid
    - uid = IP + Time + Rand (28 bytes)

### Make request log and response log

```go
	// see router_test.go

	s := NewRouter(&http.Server{
		Addr: ":9080",
	})
	s.Post("/sandbox", Sandbox)
	s.Run()
```

## Downloa and update use package
```sh
	git clone github.com/garallz/Go/tree/master/router
```
- You can download the package on your program and update to do
- Like change the log write or request (response) deal with
- When you do it, should be declared to use


## Need to do with next version

- log write system
- safety to close server
- find router path by tree