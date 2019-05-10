package listen

import (
	"encoding/json"
)

type ListenCommand int

const (
	CommandStatus  ListenCommand = 1
	CommandStart   ListenCommand = 2
	CommandStop    ListenCommand = 3
	CommandRestart ListenCommand = 4
)

const (
	ModeStart   = "START"
	ModeStop    = "STOP"
	ModeRestart = "RESTART"
	ModeRunning = "RUNNING"                     // R (TASK_RUNNING)，可执行状态。
	ModeSleep   = "SLEEPING"                    // S (TASK_INTERRUPTIBLE)，可中断的睡眠状态。
	ModeDoSleep = "TASK_UNINTERRUPTIBLE"        // D (TASK_UNINTERRUPTIBLE)，不可中断的睡眠状态。
	ModeExiting = "TASK_DEAD - EXIT_DEAD"       // X (TASK_DEAD - EXIT_DEAD)，退出状态，进程即将被销毁。
	ModeTraced  = "TASK_STOPPED or TASK_TRACED" // T (TASK_STOPPED or TASK_TRACED)，暂停状态或跟踪状态。
	ModeZombie  = "TASK_DEAD - EXIT_ZOMBIE"     // Z (TASK_DEAD - EXIT_ZOMBIE)，退出状态，进程成为僵尸进程。
	ModeDefault = "CHECK-PROC-STATE"
)

const (
	CodeSuccess = "SUCCESS"
	CodeFAIL    = "FAIL"
)

type ServerQuery struct {
	Name    string        `json:"name"`              // process name
	User    string        `json:"user,omitempty"`    // server running process by user
	Action  string        `json:"action,omitempty"`  // when control process not null
	Command ListenCommand `json:"command,omitempty"` // command: [status, start, stop, restart]
}

type ServerStatus struct {
	Name     string            `json:"name,omitempty"`
	Computer map[string]string `json:"computer,omitempty"`
	Pid      int               `json:"pid,omitempty"`
	Status   string            `json:"status,omitempty"`
	Process  map[string]string `json:"process,omitempty"`
	Code     string            `json:"code"`
	Err      string            `json:"error,omitempty"`
}

func QueryToByte(name, user, action string, mode ListenCommand) ([]byte, error) {
	return json.Marshal(&ServerQuery{
		Name:    name,
		Action:  action,
		User:    user,
		Command: mode,
	})
}
