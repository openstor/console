// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/openstor/madmin-go/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminConsoleLog(t *testing.T) {
	assert := assert.New(t)
	adminClient := AdminClientMock{}
	mockWSConn := mockConn{}
	function := "startConsoleLog(ctx, )"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testReceiver := make(chan madmin.LogInfo, 5)
	textToReceive := "test message"
	testStreamSize := 5
	isClosed := false // testReceiver is closed?

	// Test-1: Serve Console with no errors until Console finishes sending
	// define mock function behavior for minio server Console
	minioGetLogsMock = func(_ context.Context, _ string, _ int, _ string) <-chan madmin.LogInfo {
		ch := make(chan madmin.LogInfo)
		// Only success, start a routine to start reading line by line.
		go func(ch chan<- madmin.LogInfo) {
			defer close(ch)
			lines := make([]int, testStreamSize)
			// mocking sending 5 lines of info
			for range lines {
				info := madmin.LogInfo{}
				info.Message = textToReceive
				ch <- info
			}
		}(ch)
		return ch
	}
	writesCount := 1
	// mock connection WriteMessage() no error
	connWriteMessageMock = func(_ int, data []byte) error {
		// emulate that receiver gets the message written
		var t madmin.LogInfo
		_ = json.Unmarshal(data, &t)
		if writesCount == testStreamSize {
			if !isClosed {
				close(testReceiver)
				isClosed = true
			}
			return nil
		}
		testReceiver <- t
		writesCount++
		return nil
	}
	if err := startConsoleLog(ctx, mockWSConn, adminClient, LogRequest{node: "", logType: "all"}); err != nil {
		t.Errorf("Failed on %s:, error occurred: %s", function, err.Error())
	}
	// check that the TestReceiver got the same number of data from Console.
	for i := range testReceiver {
		assert.Equal(textToReceive, i.Message)
	}

	// Test-2: if error happens while writing, return error
	connWriteMessageMock = func(_ int, _ []byte) error {
		return fmt.Errorf("error on write")
	}
	if err := startConsoleLog(ctx, mockWSConn, adminClient, LogRequest{node: "", logType: "all"}); assert.Error(err) {
		assert.Equal("error on write", err.Error())
	}

	// Test-3: error happens on GetLogs Minio, Console should stop
	// and error shall be returned.
	minioGetLogsMock = func(_ context.Context, _ string, _ int, _ string) <-chan madmin.LogInfo {
		ch := make(chan madmin.LogInfo)
		// Only success, start a routine to start reading line by line.
		go func(ch chan<- madmin.LogInfo) {
			defer close(ch)
			lines := make([]int, 2)
			// mocking sending 5 lines of info
			for range lines {
				info := madmin.LogInfo{}
				info.Message = textToReceive
				ch <- info
			}
			ch <- madmin.LogInfo{Err: fmt.Errorf("error on Console")}
		}(ch)
		return ch
	}
	connWriteMessageMock = func(_ int, _ []byte) error {
		return nil
	}
	if err := startConsoleLog(ctx, mockWSConn, adminClient, LogRequest{node: "", logType: "all"}); assert.Error(err) {
		assert.Equal("error on Console", err.Error())
	}
}
