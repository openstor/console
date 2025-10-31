// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package replication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-openapi/loads"
	"github.com/openstor/console/api"
	"github.com/openstor/console/api/operations"
)

var token string

func initConsoleServer() (*api.Server, error) {
	// os.Setenv("CONSOLE_MINIO_SERVER", "localhost:9000")

	swaggerSpec, err := loads.Embedded(api.SwaggerJSON, api.FlatSwaggerJSON)
	if err != nil {
		return nil, err
	}

	noLog := func(string, ...interface{}) {
		// nothing to log
	}

	// Initialize MinIO loggers
	api.LogInfo = noLog
	api.LogError = noLog

	consoleAPI := operations.NewConsoleAPI(swaggerSpec)
	consoleAPI.Logger = noLog

	server := api.NewServer(consoleAPI)
	// register all APIs
	server.ConfigureAPI()

	// api.GlobalRootCAs, api.GlobalPublicCerts, api.GlobalTLSCertsManager = globalRootCAs, globalPublicCerts, globalTLSCerts

	consolePort, _ := strconv.Atoi("9090")

	server.Host = "0.0.0.0"
	server.Port = consolePort
	api.Port = "9090"
	api.Hostname = "0.0.0.0"

	return server, nil
}

func TestMain(m *testing.M) {
	// start console server
	go func() {
		fmt.Println("start server")
		srv, err := initConsoleServer()
		if err != nil {
			log.Println(err)
			log.Println("init fail")
			return
		}
		srv.Serve()
	}()

	fmt.Println("sleeping")
	time.Sleep(2 * time.Second)

	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	// get login credentials

	requestData := map[string]string{
		"accessKey": "minioadmin",
		"secretKey": "minioadmin",
	}

	requestDataJSON, _ := json.Marshal(requestData)

	requestDataBody := bytes.NewReader(requestDataJSON)

	request, err := http.NewRequest("POST", "http://localhost:9090/api/v1/login", requestDataBody)
	if err != nil {
		log.Println(err)
		return
	}

	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}

	if response != nil {
		for _, cookie := range response.Cookies() {
			if cookie.Name == "token" {
				token = cookie.Value
				break
			}
		}
	}

	if token == "" {
		log.Println("authentication token not found in cookies response")
		return
	}

	code := m.Run()
	os.Exit(code)
}
