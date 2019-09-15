package listen

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

const contentType = "Content-Type: application/json"

var filePath string = "./env.json"
var timeout int = 5

type EnvFile struct {
	Unique  string `json:"uniq_name"`
	Server  string `json:"serv_name"`
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

func SetTimeout(sec int) {
	if sec > 0 {
		timeout = sec
	}
}

// Query server status all in env file
func StatusAll() ([]*ServerStatus, error) {
	rows, err := readAll()
	if err != nil {
		return nil, err
	}
	return getState(rows), nil
}

// Query server status by serv_name
func StatusByServer(name string) ([]*ServerStatus, error) {
	rows, err := readSer(name)
	if err != nil {
		return nil, err
	}
	return getState(rows), nil
}

// Query server status by uniq_name
func StatusByName(name string) (*ServerStatus, error) {
	row, err := readOne(name)
	if err != nil {
		return nil, err
	}

	result := getState([]*EnvFile{row})
	return result[0], nil
}

func getState(rows []*EnvFile) []*ServerStatus {
	var c = &WaitChannel{
		data: make([]*ServerStatus, len(rows)),
	}

	for i, row := range rows {
		c.wg.Add(1)
		go c.post(i, row, "", CommandStatus)
	}
	c.wg.Wait()
	return c.data
}

func (c *WaitChannel) post(num int, env *EnvFile, action string, mode ListenCommand) {
	defer c.wg.Done()
	body, err := QueryToByte(env.Process, env.User, action, mode)
	if err != nil {
		env.ErrorStatus(err)
		return
	}

	bts, err := PostQuery(env.Url, body)
	if err != nil {
		c.data[num] = env.ErrorStatus(err)
		return
	}

	var result ServerStatus
	err = json.Unmarshal(bts, &result)
	if err != nil {
		c.data[num] = env.ErrorStatus(err)
	} else {
		result.Name = env.Unique
		c.data[num] = &result
	}
	return
}

func (e *EnvFile) ErrorStatus(err error) *ServerStatus {
	return &ServerStatus{
		Name: e.Unique,
		Code: CodeFAIL,
		Err:  err.Error(),
	}
}

func control(name string, mode ListenCommand) (*ServerStatus, error) {
	if name == "" {
		return nil, errors.New("Name Can't be Null")
	} else if mode == CommandStatus {
		return nil, errors.New("Command Mode is Wrong")
	}

	env, err := readOne(name)
	if err != nil {
		return nil, err
	}

	action := env.GetAction(mode)
	if action == "" {
		return nil, errors.New("Mode or Env Action Error")
	}

	body, err := QueryToByte(env.Process, env.User, action, mode)
	if err != nil {
		return nil, err
	}

	bts, err := PostQuery(env.Url, body)
	if err != nil {
		return nil, err
	}

	var result ServerStatus
	err = json.Unmarshal(bts, &result)
	result.Name = env.Unique

	return &result, err
}

// Control process status by uniq_name
func ControlByServer(name string, mode ListenCommand) ([]*ServerStatus, error) {
	if name == "" {
		return nil, errors.New("Name Can't be Null")
	} else if mode == CommandStatus {
		return nil, errors.New("Command Mode is Wrong")
	}

	rows, err := readSer(name)
	if err != nil {
		return nil, err
	}

	var c = &WaitChannel{
		data: make([]*ServerStatus, len(rows)),
	}

	for i, row := range rows {
		action := row.GetAction(mode)
		if action == "" {
			return nil, errors.New("Mode or Env Action Error")
		}

		c.wg.Add(1)
		go c.post(i, row, action, mode)
	}
	c.wg.Wait()

	return c.data, nil
}

func readOne(name string) (*EnvFile, error) {
	if rows, err := readAll(); err != nil {
		return nil, err
	} else {
		for _, row := range rows {
			if row.Unique == name {
				return row, nil
			}
		}
	}
	return nil, errors.New("No this uniq_name in env file")
}

func readSer(name string) ([]*EnvFile, error) {
	if rows, err := readAll(); err != nil {
		return nil, err
	} else {
		var result []*EnvFile
		for _, row := range rows {
			if row.Server == name {
				result = append(result, row)
			}
		}
		if len(result) == 0 {
			return nil, errors.New("No this serv_name in env file")
		}
		return result, nil
	}
}

func readAll() ([]*EnvFile, error) {
	var result []*EnvFile
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	} else if len(result) == 0 {
		return nil, errors.New("Env File is Null")
	} else {
		return result, nil
	}
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

// Control process status by serv_name
func ControlByUnique(name string, mode ListenCommand) (*ServerStatus, error) {
	return control(name, mode)
}

// Control process status by uniq_name to start
func ControlStart(name string) (*ServerStatus, error) {
	return control(name, CommandStart)
}

// Control process status by uniq_name to stop
func ControlStop(name string) (*ServerStatus, error) {
	return control(name, CommandStop)
}

// Control process status by uniq_name to restart
func ControlRestart(name string) (*ServerStatus, error) {
	return control(name, CommandRestart)
}

func PostQuery(url string, body []byte) ([]byte, error) {
	client := &http.Client{Timeout: time.Second * time.Duration(timeout)}
	resp, err := client.Post(url, contentType, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	} else {
		return ioutil.ReadAll(resp.Body)
	}
}
