// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"errors"
	"testing"

	"github.com/openstor/madmin-go/v4"
	"github.com/stretchr/testify/assert"
)

func TestServiceRestart(t *testing.T) {
	assert := assert.New(t)
	adminClient := AdminClientMock{}
	ctx := context.Background()
	function := "serviceRestart()"
	// Test-1 : serviceRestart() restart services no errors
	// mock function response from listGroups()
	minioServiceRestartMock = func(_ context.Context) error {
		return nil
	}
	MinioServerInfoMock = func(_ context.Context) (madmin.InfoMessage, error) {
		return madmin.InfoMessage{}, nil
	}
	if err := serviceRestart(ctx, adminClient); err != nil {
		t.Errorf("Failed on %s:, errors occurred: %s", function, err.Error())
	}

	// Test-2 : serviceRestart() returns errors on client.serviceRestart call
	// and see that the errors is handled correctly and returned
	minioServiceRestartMock = func(_ context.Context) error {
		return errors.New("error")
	}
	MinioServerInfoMock = func(_ context.Context) (madmin.InfoMessage, error) {
		return madmin.InfoMessage{}, nil
	}
	if err := serviceRestart(ctx, adminClient); assert.Error(err) {
		assert.Equal("error", err.Error())
	}

	// Test-3 : serviceRestart() returns errors on client.serverInfo() call
	// and see that the errors is handled correctly and returned
	minioServiceRestartMock = func(_ context.Context) error {
		return nil
	}
	MinioServerInfoMock = func(_ context.Context) (madmin.InfoMessage, error) {
		return madmin.InfoMessage{}, errors.New("error on server info")
	}
	if err := serviceRestart(ctx, adminClient); assert.Error(err) {
		assert.Equal("error on server info", err.Error())
	}
}
