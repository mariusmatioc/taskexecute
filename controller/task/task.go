package task

import (
	"github.com/google/uuid"
)

type TaskDesc struct {
	Id              uuid.UUID `json:"id"`   // unique id over all tasks
	Cmd             string    `json:"cmd"`  // the command to be executed
	Args            []string  `json:"args"` // the arguments for the command
	PeriodInSeconds int       `json:"period_in_seconds"`
}

func NewTaskDesc(cmd string, args []string, periodInSeconds int) TaskDesc {
	return TaskDesc{
		Id:              uuid.New(),
		Cmd:             cmd,
		Args:            args,
		PeriodInSeconds: periodInSeconds,
	}
}

type TaskResult struct {
	LastOutput []byte `json:"last_output"`
	LastError  string `json:"last_error"`
	// TODO: add timestamp
}

type Task struct {
	TaskDesc
	TaskResult
}

type GetLastRequest struct {
	ID string `json:"id"`
}

type GetLastReponse struct {
	// map from server URL to the response from the server
	ServerResp map[string]TaskResult `json:"server_response"`
}
