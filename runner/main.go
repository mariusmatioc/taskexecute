package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/mariusmatioc/taskexecute/controller/task"
	"github.com/mariusmatioc/taskexecute/runner/run"
)

func main() {
	http.HandleFunc("/submit", submit)
	http.HandleFunc("/getlast", getlast)
	http.HandleFunc("/geterrors", geterrors)
	http.HandleFunc("/kill", kill)

	http.ListenAndServe(":5001", nil)
}

// submit receives a task.TaskDesc and starts running it
func submit(w http.ResponseWriter, req *http.Request) {
	var taskDesk task.TaskDesc
	err := json.NewDecoder(req.Body).Decode(&taskDesk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := run.Submit(taskDesk); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getlast(w http.ResponseWriter, req *http.Request) {
	var input task.GetLastRequest
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, _ := uuid.FromBytes([]byte(input.ID))
	result, err := run.GetLast(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	output, _ := json.Marshal(&result)
	w.Write(output)
}

func geterrors(w http.ResponseWriter, req *http.Request) {
}

func kill(w http.ResponseWriter, req *http.Request) {
}
