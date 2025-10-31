// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openstor/console/models"
)

func TestLogSearch(t *testing.T) {
	responseItem := []map[string]interface{}{
		{
			"time":                    "2006-01-02T15:04:05Z",
			"api_time":                "GetConfigKV",
			"bucket":                  "",
			"object":                  "",
			"time_to_response_ns":     float64(452546530),
			"remote_host":             "10.116.1.94",
			"request_id":              "16595A4E30CCFE79",
			"user_agent":              "MinIO (linux; amd64) madmin-go/0.0.1",
			"response_status":         "OK",
			"response_status_code":    float64(200),
			"request_content_length":  nil,
			"response_content_length": nil,
		}, {
			"time":                    "2006-01-02T15:04:05Z",
			"api_time":                "AssumeRole",
			"bucket":                  "",
			"object":                  "",
			"time_to_response_ns":     float64(307423794),
			"remote_host":             "127.0.0.1",
			"request_id":              "16595A4DA906FBA9",
			"user_agent":              "Go-http-client/1.1",
			"response_status":         "OK",
			"response_status_code":    float64(200),
			"request_content_length":  nil,
			"response_content_length": nil,
		},
	}

	type args struct {
		apiResponse     string
		apiResponseCode int
	}

	response, _ := json.Marshal(responseItem)

	successfulResponse := &models.LogSearchResponse{
		Results: responseItem,
	}

	tests := []struct {
		name             string
		args             args
		expectedResponse *models.LogSearchResponse
		wantErr          bool
	}{
		{
			name: "200 Success response",
			args: args{
				apiResponse:     string(response),
				apiResponseCode: 200,
			},
			expectedResponse: successfulResponse,
			wantErr:          false,
		},
		{
			name: "500 unsuccessful response",
			args: args{
				apiResponse:     "Some random error",
				apiResponseCode: 500,
			},
			expectedResponse: nil,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			testRequest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.args.apiResponseCode)
				fmt.Fprintln(w, tt.args.apiResponse)
			}))
			defer testRequest.Close()

			resp, err := logSearch(testRequest.URL, "127.0.0.1")

			if (err != nil) != tt.wantErr {
				t.Errorf("logSearch() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(resp, tt.expectedResponse) {
				t.Errorf("\ngot: %d \nwant: %d", resp, tt.expectedResponse)
			}
			// if tt.wantErr {
			//	assert.Equal(tt.expectedError.Code, err.Code, fmt.Sprintf("logSearch() error code: `%v`, wantErr: `%v`", err.Code, tt.expectedError))
			//	assert.Equal(tt.expectedError.Message, err.Message, fmt.Sprintf("logSearch() error message: `%v`, wantErr: `%v`", err.Message, tt.expectedError))
			// } else {
			//	assert.Nil(err, fmt.Sprintf("logSearch() error: %v, wantErr: %v", err, tt.expectedError))
			//	buf1, err1 := tt.expectedResponse.MarshalBinary()
			//	buf2, err2 := resp.MarshalBinary()
			//	if err1 != err2 {
			//		t.Errorf("logSearch() resp: %v, expectedResponse: %v", resp, tt.expectedResponse)
			//		return
			//	}
			//	h := sha256.New()
			//	h.Write(buf1)
			//	checkSum1 := fmt.Sprintf("%x\n", h.Sum(nil))
			//	h.Reset()
			//	h.Write(buf2)
			//	checkSum2 := fmt.Sprintf("%x\n", h.Sum(nil))
			//	if checkSum1 != checkSum2 {
			//		t.Errorf("logSearch() resp: %v, expectedResponse: %v", resp, tt.expectedResponse)
			//	}
			// }
		})
	}
}
