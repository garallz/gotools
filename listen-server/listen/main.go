package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/garallz/Go/ctrl-server"
)

var errResult []byte

func init() {
	errResult = []byte(`{"code":"400","error":"Json Marshal Error"}`)
}

var post = flag.String("p", "9080", "Listen Post")

func main() {
	flag.Parse()

	http.HandleFunc("/", server)

	err := http.ListenAndServe(":"+*post, nil)
	if err != nil {
		panic(err)
	}
}

func server(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		bts := MakeErrResult("", err)
		w.Write(bts)
		return
	}

	var query listen.ServerQuery
	if err := json.Unmarshal(data, &query); err != nil {
		bts := MakeErrResult("", err)
		w.Write(bts)
	} else {
		var result interface{}
		if query.Command == listen.CommandStatus {
			result = status(&query)
		} else {
			result = control(&query)
		}

		bts, _ := json.Marshal(result)
		w.Write(bts)
	}
}

// computer status : %cpu, %cache, db(int)
// process status : %cpu, %cache, status(start, running, sleep)
func status(data *listen.ServerQuery) *listen.ServerStatus {
	pid := getpid(data.User, data.Name)
	if pid == 0 {
		return &listen.ServerStatus{
			Code: listen.CodeFAIL,
			Err:  "Process maybe not run, use Restart to run",
		}
	}

	proc := getProcessStatus(pid)
	comp := getcomputer()

	return &listen.ServerStatus{
		Pid:      pid,
		Name:     data.Name,
		Process:  proc,
		Computer: comp,
		Status:   convertState(proc["state"]),
		Code:     listen.CodeSuccess,
	}
}

func convertState(state string) string {
	switch state {
	case "S":
		return listen.ModeSleep
	case "R":
		return listen.ModeRunning
	case "Z":
		return listen.ModeZombie
	case "T":
		return listen.ModeTraced
	case "X":
		return listen.ModeExiting
	case "D":
		return listen.ModeDoSleep

	default:
		return listen.ModeDefault
	}
}

// control process status : [start, stop, restart]
func control(data *listen.ServerQuery) *listen.ServerStatus {
	err := exec.Command("sh", "-c", data.Action).Run()
	if err != nil {
		return &listen.ServerStatus{
			Code: listen.CodeFAIL,
			Err:  err.Error(),
		}
	} else {
		return &listen.ServerStatus{
			Code:     listen.CodeSuccess,
			Computer: getcomputer(),
		}
	}
}

func MakeErrResult(name string, err error) []byte {
	result, err := json.Marshal(&listen.ServerStatus{
		Name: name,
		Code: listen.CodeFAIL,
		Err:  err.Error(),
	})

	if err != nil {
		return errResult
	} else {
		return result
	}
}

func getpid(user, name string) int {
	query := fmt.Sprintf("ps -u %s -o 'pid,comm' | grep %s", user, name)
	result, err := exec.Command("sh", "-c", query).CombinedOutput()
	if err != nil || len(result) == 0 {
		return 0
	} else {
		rows := strings.Split(string(result), " ")
		for _, row := range rows {
			if row != "" {
				if pid, err := strconv.Atoi(row); err == nil && pid != 0 {
					return pid
				}
				return 0
			}
		}
		return 0
	}
}

var ComputerMemValues = []string{"type", "mem_total", "mem_used", "mem_free", "mem_shared", "mem_buffers", "mem_cached"}
var ComputerCpuValues = []string{"cpu_us", "cpu_sy", "cpu_ni", "cpu_id", "cpu_wa", "cpu_hi", "cpu_si", "cpu_st"}

// cpu, cache, rom
/*
	pid 进程ID
	comm 进程名
	pcpu 占用CPU百分比
	pmem 占用内存百分比
	rsz 占用物理内存大小
	vsz 占用虚拟内存大小
	stime 进程启动时间
*/

/*
	%us：表示用户空间程序的cpu使用率（没有通过nice调度）
	%sy：表示系统空间的cpu使用率，主要是内核程序。
	%ni：表示用户空间且通过nice调度过的程序的cpu使用率。
	%id：空闲cp
	%wa：cpu运行时在等待io的时间
	%hi：cpu处理硬中断的数量
	%si：cpu处理软中断的数量
	%st：被虚拟机偷走的cpu
*/

func getcomputer() map[string]string {
	// get mem data
	query := "free -m | grep 'Mem'"
	mem := filterValues(ComputerMemValues, query)
	delete(mem, "type")

	// get cpu data
	query = "top -bn 1 -i -c  | grep Cpu"
	cpu := filterFloat(ComputerCpuValues, query)

	for k, v := range cpu {
		mem[k] = v
	}
	return mem
}

var ProcStatusValues = []string{"pid", "comm", "state", "pcpu", "pmem", "rsz", "vsz", "stime"}

func getProcessStatus(pid int) map[string]string {
	str := "ps -p %d -o '%s' | grep %d"
	query := fmt.Sprintf(str, pid, strings.Join(ProcStatusValues, ","), pid)
	return filterValues(ProcStatusValues, query)
}

func filterValues(values []string, query string) map[string]string {
	resp, err := exec.Command("sh", "-c", query).CombinedOutput()
	if err != nil {
		return nil
	}

	rows := strings.Split(string(resp), " ")

	var result = make(map[string]string)
	var data []string

	for _, row := range rows {
		if row != "" {
			data = append(data, row)
		}
	}

	if len(values) == len(data) {
		for i, v := range values {
			result[v] = data[i]
		}
	}
	return result
}

func filterFloat(values []string, query string) map[string]string {
	resp, err := exec.Command("sh", "-c", query).CombinedOutput()
	if err != nil {
		return nil
	}

	body := strings.Replace(string(resp), "\n", " ", -1)
	rows := strings.Split(body, " ")

	var result = make(map[string]string)
	var data []string

	for _, row := range rows {
		if row != "" {
			if _, err := strconv.ParseFloat(row, 64); err == nil {
				data = append(data, row)
			}
		}
	}

	if len(values) == len(data) {
		for i, v := range values {
			result[v] = data[i]
		}
	}
	return result
}
