// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/go-openapi/runtime/middleware"
	"github.com/openstor/console/api/operations/system"
	"github.com/openstor/console/models"

	"github.com/go-openapi/loads"
	"github.com/openstor/console/api/operations"
	"github.com/openstor/madmin-go/v4"

	asrt "github.com/stretchr/testify/assert"
)

func TestArnsList(t *testing.T) {
	assert := asrt.New(t)
	adminClient := AdminClientMock{}
	// Test-1 : getArns() returns proper arn list
	MinioServerInfoMock = func(_ context.Context) (madmin.InfoMessage, error) {
		return madmin.InfoMessage{
			SQSARN: []string{"uno"},
		}, nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	arnsList, err := getArns(ctx, adminClient)
	assert.NotNil(arnsList, "arn list was returned nil")
	if arnsList != nil {
		assert.Equal(len(arnsList.Arns), 1, "Incorrect arns count")
	}
	assert.Nil(err, "Error should have been nil")

	// Test-2 : getArns(ctx) fails for whatever reason
	MinioServerInfoMock = func(_ context.Context) (madmin.InfoMessage, error) {
		return madmin.InfoMessage{}, errors.New("some reason")
	}

	arnsList, err = getArns(ctx, adminClient)
	assert.Nil(arnsList, "arn list was not returned nil")
	assert.NotNil(err, "An error should have been returned")
}

func TestRegisterAdminArnsHandlers(t *testing.T) {
	assert := asrt.New(t)
	swaggerSpec, err := loads.Embedded(SwaggerJSON, FlatSwaggerJSON)
	if err != nil {
		assert.Fail("Error")
	}
	api := operations.NewConsoleAPI(swaggerSpec)
	api.SystemArnListHandler = nil
	registerAdminArnsHandlers(api)
	if api.SystemArnListHandler == nil {
		assert.Fail("Assignment should happen")
	} else {
		fmt.Println("Function got assigned: ", api.SystemArnListHandler)
	}

	// To test error case in registerAdminArnsHandlers
	request, _ := http.NewRequest(
		"GET",
		"http://localhost:9090/api/v1/buckets/",
		nil,
	)
	ArnListParamsStruct := system.ArnListParams{
		HTTPRequest: request,
	}
	modelsPrincipal := models.Principal{
		STSAccessKeyID: "accesskey",
	}
	var value middleware.Responder = api.SystemArnListHandler.Handle(ArnListParamsStruct, &modelsPrincipal)
	str := fmt.Sprintf("%#v", value)
	fmt.Println("value: ", str)
	assert.Equal(strings.Contains(str, "_statusCode:500"), true)
}
