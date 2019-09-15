# Listen and Control Server Status

## env file set

```json
[
	{
		"url":"http://127.0.0.1:8090",
		"uniq_name":"unqi_one"
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

+ url:       listen server url
+ uniq_name: unique name with server process
+ serv_name: server name or process clump
+ proc_name: process name in server
+ user_name: run in server by user
+ action:    server control process by command

## Useing

### Interface

```go
// Set env file path
func SetEnvFilePath(path string)

// Set post server timeout, unin is second(s)
func SetTimeout(sec int)

// Get all process status by env file setting
func StatusAll(path string) ([]*ServerStatus, error)

// Get server clump process status by serv_name
func StatusByServer(path string) ([]*ServerStatus, error)

// Get process clump status with serv_name by env file setting
func StatusByName(path string, name string) (*ServerStatus, error)

// Control process status
func Control(name string, mode ListenCommand) (*ServerResult, error)

// Control process clump by serv_name
func ControlByServer(name string, mode ListenCommand) ([]*ServerResult, error)

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
	Name    string              // serv_name
	Server  map[string]string   // server detail
	Pid     int                 // process pid
	Status  string              // process status
	Process map[string]string   // process detail
	Code    string              // return status: [SUCCESS, FAIL]
	Err     string              // if Code == FAIL, Err != Null
}
```

```sh
# process state
	pid  ：进程ID
	comm ：进程名
	pcpu ：占用CPU百分比
	pmem ：占用内存百分比
	rsz  ：占用物理内存大小
	vsz  ：占用虚拟内存大小
	stime：进程启动时间

# server cpu
	cpu_us：表示用户空间程序的cpu使用率（没有通过nice调度）
	cpu_sy：表示系统空间的cpu使用率，主要是内核程序
	cpu_ni：表示用户空间且通过nice调度过的程序的cpu使用率
	cpu_id：空闲cpu
	cpu_wa：cpu运行时在等待io的时间
	cpu_hi：cpu处理硬中断的数量
	cpu_si：cpu处理软中断的数量
	cpu_st：被虚拟机偷走的cpu

# server cache (-m)
	mem_total  ：系统总的可用物理内存和交换空间大小
	mem_used   ：已经被使用的物理内存和交换空间
	mem_free   ：空闲的物理内存和交换空间
	mem_shared ：被共享使用的物理内存大小
	mem_buffers：Buffer缓冲区大小
	mem_cached ：Cache缓冲区大小
```