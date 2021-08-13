package main

// controller implements a REST API to submit tasks and retrieve results

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/mariusmatioc/taskexecute/controller/task"
)

// This should come from a config file
var runners = []string{"http:localhose:5001", "http:localhose:5002"}

func main() {
	http.HandleFunc("/submit", submit)
	http.HandleFunc("/getlast", getlast)
	http.HandleFunc("/geterrors", geterrors)
	http.HandleFunc("/kill", kill)

	http.ListenAndServe(":8001", nil)
}

// submitRequest is the body of the "submit" request
type submitRequest struct {
	Cmd             string   `json:"cmd"`
	Args            []string `json:"args"`
	PeriodInSeconds int      `json:"period_in_seconds"`
}

// submit receives a JSON task request and returns the id of the task as a string
func submit(w http.ResponseWriter, req *http.Request) {
	var subReq submitRequest
	err := json.NewDecoder(req.Body).Decode(&subReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tskd := task.NewTaskDesc(subReq.Cmd, subReq.Args, subReq.PeriodInSeconds)
	body, err := json.Marshal(tskd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Send to all runners
	_, err = submitToAll("submit", body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// At least one runner got it
	io.WriteString(w, tskd.Id.String())
}

// getlast returns the output (and errors) from every runner
func getlast(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fromServers, err := submitToAll("getlast", body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := task.GetLastReponse{ServerResp: make(map[string]task.TaskResult)}
	for runner, body := range fromServers {
		var taskRes task.TaskResult
		_ = json.Unmarshal(body, &taskRes)
		resp.ServerResp[runner] = taskRes
	}
	output, _ := json.Marshal(&resp)
	w.Write(output)
}

func geterrors(w http.ResponseWriter, req *http.Request) {
	// TODO: same logic as getLast
}

func kill(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
}

// submitToAll returns the bodies of the responses if there is at least one returned
func submitToAll(uri string, body []byte) (response map[string][]byte, err error) {
	response = make(map[string][]byte)
	atLeastOne := false
	for _, runner := range runners {
		req, err := http.NewRequest("POST", path.Join(runner, uri), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed: %d", resp.StatusCode)
			continue
		}
		atLeastOne = true
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		response[runner] = body
	}
	if atLeastOne {
		return
	}
	err = fmt.Errorf("no runners available")
	return
}
