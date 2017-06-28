package debug

import (
	"fmt"
	"os"

	"log"

	"strconv"

	m "github.com/getgauge/spider/gauge_messages"
	"github.com/shirou/gopsutil/process"
)

const (
	localhost = "localhost"
)

type processInfo struct {
	Port string
	Cwd  string
	Pid  int
}

func getPInfos() []processInfo {
	pids, err := process.Pids()
	if err != nil {
		log.Fatalf(err.Error())
	}
	processes, errors := getProcesses(pids)
	for _, e := range errors {
		log.Println(e)
	}
	return getProcessInfo(processes)
}

func getProcessInfo(processes []*process.Process) (infos []processInfo) {
	for _, p := range processes {
		port, err := getPort(p)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		cwd, err := p.Cwd()
		if err != nil {
			api, err := newAPI(localhost, port)
			if err != nil {
				cwd = "N/A"
			} else {
				msg, err := api.getResponse(&m.APIMessage{MessageType: m.APIMessage_GetProjectRootRequest, ProjectRootRequest: &m.GetProjectRootRequest{}})
				if err != nil {
					cwd = "N/A"
				} else {
					cwd = msg.ProjectRootResponse.ProjectRoot
				}
				api.close()
			}
		}
		infos = append(infos, processInfo{Port: port, Cwd: cwd, Pid: int(p.Pid)})
	}
	return infos
}

func getPort(p *process.Process) (string, error) {
	args, err := p.CmdlineSlice()
	if err == nil {
		if len(args) > 2 {
			if _, err := strconv.Atoi(args[2]); err == nil {
				return args[2], nil
			}
		}
	}
	conns, err := p.Connections()
	if err != nil {
		return "", fmt.Errorf("Cannot find a port for the process %d", p.Pid)
	}
	var port uint32 = 0
	for _, c := range conns {
		p := fmt.Sprintf("%d", c.Laddr.Port)
		api, err := newAPI(localhost, p)
		if err == nil {
			port = c.Laddr.Port
			api.close()
			break
		}
	}
	if port == 0 {
		return "", fmt.Errorf("Cannot find a port for the process %d", p.Pid)
	}
	return fmt.Sprintf("%d", port), nil
}

func getProcesses(pids []int32) (processes []*process.Process, errors []string) {
	for _, pid := range pids {
		if int(pid) == os.Getpid() {
			continue
		}
		p, err := process.NewProcess(pid)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Cannot get process for pid: %d. Error: %s", pid, err.Error()))
			continue
		}
		name, err := p.Name()
		if err != nil {
			errors = append(errors, fmt.Sprintf("Process name error for pid: %v. Error: %v", pid, err))
			continue
		}
		if name == "gauge" || name == "gauge.exe" {
			_, err = p.CmdlineSlice()
			if err != nil {
				errors = append(errors, fmt.Sprintf("Process Cmd line slice error for pid: %v. Error: %v", p.Pid, err))
				continue
			}
			processes = append(processes, p)
		}
	}
	return processes, errors
}
