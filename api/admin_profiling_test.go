// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/openstor/madmin-go/v4"
	"github.com/stretchr/testify/assert"
)

// Implementing fake closingBuffer to mock stopProfiling() (io.ReadCloser, error)
type ClosingBuffer struct {
	*bytes.Buffer
}

// Implementing a fake Close function for io.ReadCloser
func (cb *ClosingBuffer) Close() error {
	return nil
}

func TestStartProfiling(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert := assert.New(t)
	adminClient := AdminClientMock{}
	mockWSConn := mockConn{}
	function := "startProfiling()"
	testOptions := &profileOptions{
		Types: "cpu",
	}

	// Test-1 : startProfiling() Get response from MinIO server with one profiling object without errors
	// mock function response from startProfiling()
	minioStartProfiling = func(_ madmin.ProfilerType, _ time.Duration) (io.ReadCloser, error) {
		return &ClosingBuffer{bytes.NewBufferString("In memory string eaeae")}, nil
	}
	// mock function response from mockConn.writeMessage()
	connWriteMessageMock = func(_ int, _ []byte) error {
		return nil
	}
	err := startProfiling(ctx, mockWSConn, adminClient, testOptions)
	if err != nil {
		t.Errorf("Failed on %s:, error occurred: %s", function, err.Error())
	}
	assert.Equal(err, nil)

	// Test-2 : startProfiling() Correctly handles errors returned by MinIO
	// mock function response from startProfiling()
	minioStartProfiling = func(_ madmin.ProfilerType, _ time.Duration) (io.ReadCloser, error) {
		return nil, errors.New("error")
	}
	err = startProfiling(ctx, mockWSConn, adminClient, testOptions)
	if assert.Error(err) {
		assert.Equal("error", err.Error())
	}

	// Test-3: getProfileOptionsFromReq() correctly returns profile options from request
	u, _ := url.Parse("ws://localhost/ws/profile?types=cpu,mem,block,mutex,trace,threads,goroutines")
	req := &http.Request{
		URL: u,
	}
	opts, err := getProfileOptionsFromReq(req)
	if assert.NoError(err) {
		expectedOptions := profileOptions{
			Types: "cpu,mem,block,mutex,trace,threads,goroutines",
		}
		assert.Equal(expectedOptions.Types, opts.Types)
	}
}
