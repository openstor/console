// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ConfigAPI(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name           string
		expectedStatus int
		expectedError  error
	}{
		{
			name:           "Config - Valid",
			expectedStatus: 200,
			expectedError:  nil,
		},
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			request, err := http.NewRequest("GET", "http://localhost:9090/api/v1/configs", nil)
			if err != nil {
				log.Println(err)
				return
			}
			request.Header.Add("Cookie", fmt.Sprintf("token=%s", token))
			request.Header.Add("Content-Type", "application/json")
			response, err := client.Do(request)
			if err != nil {
				log.Println(err)
				return
			}
			if response != nil {
				assert.Equal(tt.expectedStatus, response.StatusCode, tt.name+" Failed")
			}
		})
	}
}

func Test_GetConfigAPI(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		name string
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedError  error
	}{
		{
			name: "Get Config - Valid",
			args: args{
				name: "storage_class",
			},
			expectedStatus: 200,
			expectedError:  nil,
		},
		{
			name: "Get Config - Invalid",
			args: args{
				name: "asdf",
			},
			expectedStatus: 404,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			client := &http.Client{
				Timeout: 3 * time.Second,
			}

			request, err := http.NewRequest(
				"GET", fmt.Sprintf("http://localhost:9090/api/v1/configs/%s", tt.args.name), nil)
			if err != nil {
				log.Println(err)
				return
			}
			request.Header.Add("Cookie", fmt.Sprintf("token=%s", token))
			request.Header.Add("Content-Type", "application/json")
			response, err := client.Do(request)
			if err != nil {
				log.Println(err)
				return
			}
			if response != nil {
				assert.Equal(tt.expectedStatus, response.StatusCode, tt.name+" Failed")
			}
		})
	}
}

func Test_SetConfigAPI(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		name      string
		keyValues []map[string]interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedError  error
	}{
		{
			name: "Set Config - Valid",
			args: args{
				name:      "region",
				keyValues: []map[string]interface{}{{"key": "name", "value": "testServer"}, {"key": "region", "value": "us-west-1"}},
			},
			expectedStatus: 200,
			expectedError:  nil,
		},
		{
			name: "Set Config - Invalid",
			args: args{
				name:      "regiontest",
				keyValues: []map[string]interface{}{{"key": "name", "value": "testServer"}, {"key": "region", "value": "us-west-1"}},
			},
			expectedStatus: 500,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			client := &http.Client{
				Timeout: 3 * time.Second,
			}

			requestDataPolicy := map[string]interface{}{}

			requestDataPolicy["key_values"] = tt.args.keyValues

			requestDataJSON, _ := json.Marshal(requestDataPolicy)
			requestDataBody := bytes.NewReader(requestDataJSON)
			request, err := http.NewRequest(
				"PUT", fmt.Sprintf("http://localhost:9090/api/v1/configs/%s", tt.args.name), requestDataBody)
			if err != nil {
				log.Println(err)
				return
			}
			request.Header.Add("Cookie", fmt.Sprintf("token=%s", token))
			request.Header.Add("Content-Type", "application/json")
			response, err := client.Do(request)
			if err != nil {
				log.Println(err)
				return
			}
			if response != nil {
				assert.Equal(tt.expectedStatus, response.StatusCode, tt.name+" Failed")
			}
		})
	}
}

func Test_ResetConfigAPI(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		name string
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedError  error
	}{
		{
			name: "Reset Config - Valid",
			args: args{
				name: "region",
			},
			expectedStatus: 200,
			expectedError:  nil,
		},
		{
			name: "Reset Config - Invalid",
			args: args{
				name: "regiontest",
			},
			expectedStatus: 500,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			client := &http.Client{
				Timeout: 3 * time.Second,
			}

			requestDataPolicy := map[string]interface{}{}

			requestDataJSON, _ := json.Marshal(requestDataPolicy)
			requestDataBody := bytes.NewReader(requestDataJSON)
			request, err := http.NewRequest(
				"POST", fmt.Sprintf("http://localhost:9090/api/v1/configs/%s/reset", tt.args.name), requestDataBody)
			if err != nil {
				log.Println(err)
				return
			}
			request.Header.Add("Cookie", fmt.Sprintf("token=%s", token))
			request.Header.Add("Content-Type", "application/json")
			response, err := client.Do(request)
			if err != nil {
				log.Println(err)
				return
			}
			if response != nil {
				assert.Equal(tt.expectedStatus, response.StatusCode, tt.name+" Failed")
			}
		})
	}
}
