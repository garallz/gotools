# Golang Router

## Base on Golang net/http

### Auto deal with response body to json and xml

### Make request uid
    - uid = IP + Time + Rand

### Make request log and response log

```go

	s := NewRouter(&http.Server{
		Addr: ":9080",
	})
	s.Post("/sandbox", Sandbox)
	s.Run()
```