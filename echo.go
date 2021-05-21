package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"schneider.vip/problem"
)

type Echo struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type EchoHandler struct{}

func (s *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var echoBody []byte
	var echoResponse = Echo{}
	var err error
	var response []byte

	if r.Method != "POST" {
		problem.New(
			problem.Status(http.StatusMethodNotAllowed),
			problem.Detail("Method Not Allowed"),
		).WriteTo(w)
		return
	}

	if echoBody, err = ioutil.ReadAll(r.Body); err != nil {
		problem.New(
			problem.Status(http.StatusBadRequest),
			problem.Detail("Unable to read response body"),
		).WriteTo(w)
		return
	}

	if len(echoBody) < 1 {
		problem.New(
			problem.Status(http.StatusBadRequest),
			problem.Detail("Body must include `message`"),
		).WriteTo(w)
		return
	}

	if err = json.Unmarshal(echoBody, &echoResponse); err != nil {
		problem.New(
			problem.Status(http.StatusBadRequest),
			problem.Detail("Invalid JSON"),
		).WriteTo(w)
		return
	}

	echoResponse.Timestamp = time.Now().Format(time.RFC3339)

	if response, err = json.MarshalIndent(echoResponse, "", "  "); err != nil {
		problem.New(
			problem.Status(http.StatusInternalServerError),
			problem.Detail("Unable to echo"),
		).WriteTo(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
