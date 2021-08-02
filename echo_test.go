/*
 * Copyright 2021 Nikki Nikkhoui <nnikkhoui@wikimedia.org> and Wikimedia Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const prefixURI = "/v0"

func TestInvalidMethod(t *testing.T) {

	req, err := http.NewRequest("GET", "/v0/echo", nil)

	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()
	handler := EchoHandler{}
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Incorrect status code returned")

}

func TestEmptyBody(t *testing.T) {

	req, err := http.NewRequest("POST", "/v0/echo", strings.NewReader(""))

	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()
	handler := EchoHandler{}

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Incorrect status code returned")
}

func TestInvalidJSON(t *testing.T) {

	req, err := http.NewRequest("POST", "/v0/echo", strings.NewReader("bad req body"))

	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()
	handler := EchoHandler{}

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Incorrect status code returned")
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestFailedBody(t *testing.T) {

	req, err := http.NewRequest("POST", "/v0/echo", errReader(0))

	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()

	handler := EchoHandler{}

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Incorrect status code returned")
}

type Body struct {
	message string
}

func TestTimeStamp(t *testing.T) {
	msg, _ := json.Marshal(Echo{Message: "echo"})

	req, err := http.NewRequest("POST", "/v0/echo", bytes.NewBuffer(msg))

	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()
	handler := EchoHandler{}

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Incorrect status code returned")
}
