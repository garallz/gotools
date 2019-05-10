package listen

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

const contentType = "Content-Type: application/json"

var filePath string = "./env.json"

type EnvFile struct {
	Name    string `json:"serv_name"`
	Process string `json:"proc_name"`
	User    string `json:"user_name"`
	Url     string `json:"url"`
	Action  struct {
		Start   string `json:"start"`
		Stop    string `json:"stop"`
		Restart string `json:"restart"`
	} `json:"action"`
}

type WaitChannel struct {
	wg   sync.WaitGroup
	data []*ServerStatus
}

// Default env file path : ./env.json
func SetEnvFilePath(path string) {
	if path == "" {
		panic("Listen Server Env File Path Set Null")
	} else {
		if _, err := os.Open(path); os.IsNotExist(err) {
			panic("Env file path not exist")
		} else {
			filePath = path
		}
	}
}

func StatusAll() ([]*ServerStatus, error) {
	rows, err := readAll()
	if err != nil {
		return nil, err
	} else if len(rows) == 0 {
		return nil, errors.New("Env File is Null")
	}

	var c = &WaitChannel{
		data: make([]*ServerStatus, len(rows)),
	}

	for i, row := range rows {
		c.wg.Add(1)
		go c.post(i, row)
	}
	c.wg.Wait()
	return c.data, nil
}

// Query server status by server_name
// default path: ./env.json
func StatusByName(name string) (*ServerStatus, error) {
	row, err := readOne(name)
	if err != nil {
		return nil, err
	} else if row == nil {
		return nil, errors.New("Not Server Name in Env File")
	}

	var c = &WaitChannel{
		data: make([]*ServerStatus, 1),
	}
	c.wg.Add(1)
	go c.post(0, row)
	c.wg.Wait()

	return c.data[0], nil
}

func (c *WaitChannel) post(num int, env *EnvFile) {
	defer c.wg.Done()
	body, err := QueryToByte(env.Process, env.User, "", CommandStatus)
	if err != nil {
		env.ErrorStatus(err)
		return
	}

	resp, err := http.Post(env.Url, contentType, bytes.NewBuffer([]byte(body)))
	if err != nil {
		c.data[num] = env.ErrorStatus(err)
		return
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.data[num] = env.ErrorStatus(err)
		return
	}

	var result ServerStatus
	err = json.Unmarshal(bts, &result)
	if err != nil {
		c.data[num] = env.ErrorStatus(err)
	} else {
		result.Name = env.Name
		c.data[num] = &result
	}
	return
}

func (e *EnvFile) ErrorStatus(err error) *ServerStatus {
	return &ServerStatus{
		Name: e.Name,
		Code: CodeFAIL,
		Err:  err.Error(),
	}
}

func Control(name string, mode ListenCommand) (*ServerStatus, error) {
	if name == "" {
		return nil, errors.New("Name Can't be Null")
	} else if mode == CommandStatus {
		return nil, errors.New("Command Mode is Wrong")
	}

	env, err := readOne(name)
	if err != nil {
		return nil, err
	} else if env == nil {
		return nil, errors.New("Don't Have This Server")
	}

	action := env.GetAction(mode)
	if action == "" {
		return nil, errors.New("Mode or Env Action Error")
	}

	body, err := QueryToByte(env.Process, env.User, action, mode)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(env.Url, contentType, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ServerStatus
	err = json.Unmarshal(bts, &result)

	return &result, err
}

func readOne(name string) (*EnvFile, error) {
	rows, err := readAll()
	if err != nil {
		return nil, err
	} else {
		for _, row := range rows {
			if row.Name == name {
				return row, nil
			}
		}
	}
	return nil, nil
}

func readAll() ([]*EnvFile, error) {
	var result []*EnvFile
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &result)
	return result, err
}

func (e *EnvFile) GetAction(m ListenCommand) string {
	switch m {
	case CommandStart:
		return e.Action.Start
	case CommandStop:
		return e.Action.Stop
	case CommandRestart:
		return e.Action.Restart
	default:
		return ""
	}
}

func ControlStart(name string) (*ServerStatus, error) {
	return Control(name, CommandStart)
}

func ControlStop(name string) (*ServerStatus, error) {
	return Control(name, CommandStop)
}

func ControlRestart(name string) (*ServerStatus, error) {
	return Control(name, CommandRestart)
}
