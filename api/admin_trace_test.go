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

func TestAdminTrace(t *testing.T) {
	assert := assert.New(t)
	adminClient := AdminClientMock{}
	mockWSConn := mockConn{}
	function := "startTraceInfo(ctx, )"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testReceiver := make(chan shortTraceMsg, 5)
	textToReceive := "test"
	testStreamSize := 5
	isClosed := false // testReceiver is closed?

	// Test-1: Serve Trace with no errors until trace finishes sending
	// define mock function behavior for minio server Trace
	minioServiceTraceMock = func(_ context.Context, _ int64, _, _, _, _, _ bool) <-chan madmin.ServiceTraceInfo {
		ch := make(chan madmin.ServiceTraceInfo)
		// Only success, start a routine to start reading line by line.
		go func(ch chan<- madmin.ServiceTraceInfo) {
			defer close(ch)
			lines := make([]int, testStreamSize)
			// mocking sending 5 lines of info
			for range lines {
				info := madmin.TraceInfo{
					FuncName: textToReceive,
				}
				ch <- madmin.ServiceTraceInfo{Trace: info}
			}
		}(ch)
		return ch
	}
	writesCount := 1
	// mock connection WriteMessage() no error
	connWriteMessageMock = func(_ int, data []byte) error {
		// emulate that receiver gets the message written
		var t shortTraceMsg
		_ = json.Unmarshal(data, &t)
		if writesCount == testStreamSize {
			// for testing we need to close the receiver channel
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
	if err := startTraceInfo(ctx, mockWSConn, adminClient, TraceRequest{s3: true, internal: true, storage: true, os: true, onlyErrors: false}); err != nil {
		t.Errorf("Failed on %s:, error occurred: %s", function, err.Error())
	}
	// check that the TestReceiver got the same number of data from trace.
	for i := range testReceiver {
		assert.Equal(textToReceive, i.FuncName)
	}

	// Test-2: if error happens while writing, return error
	connWriteMessageMock = func(_ int, _ []byte) error {
		return fmt.Errorf("error on write")
	}
	if err := startTraceInfo(ctx, mockWSConn, adminClient, TraceRequest{}); assert.Error(err) {
		assert.Equal("error on write", err.Error())
	}

	// Test-3: error happens on serviceTrace Minio, trace should stop
	// and error shall be returned.
	minioServiceTraceMock = func(_ context.Context, _ int64, _, _, _, _, _ bool) <-chan madmin.ServiceTraceInfo {
		ch := make(chan madmin.ServiceTraceInfo)
		// Only success, start a routine to start reading line by line.
		go func(ch chan<- madmin.ServiceTraceInfo) {
			defer close(ch)
			lines := make([]int, 2)
			// mocking sending 5 lines of info
			for range lines {
				info := madmin.TraceInfo{
					NodeName: "test",
				}
				ch <- madmin.ServiceTraceInfo{Trace: info}
			}
			ch <- madmin.ServiceTraceInfo{Err: fmt.Errorf("error on trace")}
		}(ch)
		return ch
	}
	connWriteMessageMock = func(_ int, _ []byte) error {
		return nil
	}
	if err := startTraceInfo(ctx, mockWSConn, adminClient, TraceRequest{}); assert.Error(err) {
		assert.Equal("error on trace", err.Error())
	}
}
