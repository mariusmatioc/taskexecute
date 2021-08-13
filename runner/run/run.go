package run

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mariusmatioc/taskexecute/controller/task"
)

// a map of all currently running tasks
var mtx sync.Mutex

// TODO: perhaps add timestamps to lastResult and errors
// The result of the last run
var lastResult map[uuid.UUID][]byte

const maxErrors = 10 // older errors are discarded
// Errors encountered so far
var errors map[uuid.UUID][]string

func init() {
	lastResult = make(map[uuid.UUID][]byte)
	errors = make(map[uuid.UUID][]string)
}

func Submit(td task.TaskDesc) error {
	mtx.Lock()
	defer mtx.Unlock()
	_, found := lastResult[td.Id]
	if found {
		return fmt.Errorf("task already exists")
	}
	lastResult[td.Id] = nil
	errors[td.Id] = make([]string, 0)
	go runPeriodically(td)
	return nil
}

// GetLast returns the result and error of the last run
func GetLast(id uuid.UUID) (result task.TaskResult, err error) {
	mtx.Lock()
	lastRes, found := lastResult[id]
	if !found {
		err = fmt.Errorf("taks does not exist")
	}
	lastErr := errors[id][len(errors[id])-1]
	mtx.Unlock()
	result.LastOutput = lastRes
	result.LastError = lastErr
	return
}

// runPeriodically runs the given task periodically
func runPeriodically(td task.TaskDesc) {
	// TODO: add kill/terminate
	for {
		select {
		case <-time.After(time.Second * time.Duration(td.PeriodInSeconds)):
			output, err := exec.Command(td.Cmd, td.Args...).Output()
			mtx.Lock()
			if err != nil {
				errMsg := ""
				exitError := err.(*exec.ExitError)
				if exitError != nil {
					errMsg = exitError.Error()
				} else {
					errMsg = err.Error()
				}
				errors[td.Id] = append(errors[td.Id], errMsg)
				l := len(errors[td.Id])
				if l > maxErrors {
					errors[td.Id] = errors[td.Id][l-maxErrors:]
				}
			} else {
				lastResult[td.Id] = output
			}
			mtx.Unlock()
		}
	}
}
