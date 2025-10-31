// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package integration

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestStartProfiling(t *testing.T) {
	testAssert := assert.New(t)

	tests := []struct {
		name string
	}{
		{
			name: "start/stop profiling",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			files := map[string]bool{
				"profile-127.0.0.1:9000-goroutines.txt":                false,
				"profile-127.0.0.1:9000-goroutines-before.txt":         false,
				"profile-127.0.0.1:9000-goroutines-before,debug=2.txt": false,
				"profile-127.0.0.1:9000-threads-before.pprof":          false,
				"profile-127.0.0.1:9000-mem.pprof":                     false,
				"profile-127.0.0.1:9000-threads.pprof":                 false,
				"profile-127.0.0.1:9000-cpu.pprof":                     false,
				"profile-127.0.0.1:9000-mem-before.pprof":              false,
				"profile-127.0.0.1:9000-block.pprof":                   false,
				"profile-127.0.0.1:9000-trace.trace":                   false,
				"profile-127.0.0.1:9000-mutex.pprof":                   false,
				"profile-127.0.0.1:9000-mutex-before.pprof":            false,
			}

			wsDestination := "/ws/profile?types=cpu,mem,block,mutex,trace,threads,goroutines"
			wsFinalURL := fmt.Sprintf("ws://localhost:9090%s", wsDestination)

			ws, _, err := websocket.DefaultDialer.Dial(wsFinalURL, nil)
			if err != nil {
				log.Println(err)
				return
			}
			defer ws.Close()

			_, zipFileBytes, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			filetype := http.DetectContentType(zipFileBytes)
			testAssert.Equal("application/zip", filetype)

			zipReader, err := zip.NewReader(bytes.NewReader(zipFileBytes), int64(len(zipFileBytes)))
			if err != nil {
				testAssert.Nil(err, fmt.Sprintf("%s returned an error: %v", tt.name, err))
			}

			// Read all the files from zip archive
			for _, zipFile := range zipReader.File {
				files[zipFile.Name] = true
			}

			for k, v := range files {
				testAssert.Equal(true, v, fmt.Sprintf("%s : compressed file expected to have %v file inside", tt.name, k))
			}
		})
	}
}
