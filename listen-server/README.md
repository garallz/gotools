# Listen and Control Server Status

## env file set

```json
[
	{
		"url":"http://127.0.0.1:8090", 			# server url
		"serv_name":"server_one", 				# server process name
		"proc_name":"go", 						# process name
		"user_name":"user_name", 				# server user name
		"action":{								# server control process command 	
			"start":"systemctl start listen",	# eg
			"stop":"",
			"restart":""
		}
	}
]
```

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