# Listen and Control Server Status

## env file set

```json
[
	{
		"url":"http://127.0.0.1:8090",
		"serv_name":"server_one",
		"proc_name":"go",
		"user_name":"user_name",
		"action":{
			"start":"systemctl start listen",
			"stop":"",
			"restart":""
		}
	}
]
```

- url: 		 listen server url
- serv_name: specify the unique process name of server
- proc_name: process name in server
- user_name: run in server by user
- action:	 server control process by command

## Useing

### Interface

```go
// Set env file path
func SetEnvFilePath(path string)

// Get all process status by env file setting
func StatusAll(path string) ([]*ServerStatus, error)

// Get process status with serv_name by env file setting
func StatusByName(path string, name string) (*ServerStatus, error)

// Control process status
func Control(name string, mode ListenCommand) (*ServerResult, error)

func ControlStart(name string) (*ServerResult, error)
func ControlStop(name string) (*ServerResult, error)
func ControlRestart(name string) (*ServerResult, error)
```

### Listen main

```go
// program in -/listen/main.go
$ go build -o listen

// run in server backgroud
$ nohup ./listen -p 8090 & 	
```

### Return struct

```go
type ServerStatus struct {
	Name     string            `json:"name,omitempty"`		// serv_name
	Server map[string]string `json:"server,omitempty"`	// server detail
	Pid      int               `json:"pid,omitempty"`		// process pid
	Status   string            `json:"status,omitempty"`	// process status
	Process  map[string]string `json:"process,omitempty"`	// process detail
	Code     string            `json:"code"`				// return status: [SUCCESS, FAIL]
	Err      string            `json:"error,omitempty"`		// if Code == FAIL, Err != Null
}
```