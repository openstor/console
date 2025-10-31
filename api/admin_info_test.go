// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/openstor/console/pkg/utils"

	"github.com/openstor/console/api/operations"
	systemApi "github.com/openstor/console/api/operations/system"
	"github.com/openstor/console/models"
	"github.com/openstor/madmin-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AdminInfoTestSuite struct {
	suite.Suite
	assert              *assert.Assertions
	currentServer       string
	isServerSet         bool
	isPrometheusRequest bool
	server              *httptest.Server
	adminClient         AdminClientMock
}

func (suite *AdminInfoTestSuite) SetupSuite() {
	suite.assert = assert.New(suite.T())
	suite.adminClient = AdminClientMock{}
	MinioServerInfoMock = func(_ context.Context) (madmin.InfoMessage, error) {
		return madmin.InfoMessage{
			Servers: []madmin.ServerProperties{{
				Disks: []madmin.Disk{{}},
			}},
			Backend: madmin.ErasureBackend{Type: "mock"},
		}, nil
	}
}

func (suite *AdminInfoTestSuite) SetupTest() {
	suite.server = httptest.NewServer(http.HandlerFunc(suite.serverHandler))
	suite.currentServer, suite.isServerSet = os.LookupEnv(ConsoleMinIOServer)
	os.Setenv(ConsoleMinIOServer, suite.server.URL)
}

func (suite *AdminInfoTestSuite) serverHandler(w http.ResponseWriter, _ *http.Request) {
	if suite.isPrometheusRequest {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
}

func (suite *AdminInfoTestSuite) TearDownSuite() {
}

func (suite *AdminInfoTestSuite) TearDownTest() {
	if suite.isServerSet {
		os.Setenv(ConsoleMinIOServer, suite.currentServer)
	} else {
		os.Unsetenv(ConsoleMinIOServer)
	}
}

func (suite *AdminInfoTestSuite) TestRegisterAdminInfoHandlers() {
	api := &operations.ConsoleAPI{}
	suite.assertHandlersAreNil(api)
	registerAdminInfoHandlers(api)
	suite.assertHandlersAreNotNil(api)
}

func (suite *AdminInfoTestSuite) assertHandlersAreNil(api *operations.ConsoleAPI) {
	suite.assert.Nil(api.SystemAdminInfoHandler)
	suite.assert.Nil(api.SystemDashboardWidgetDetailsHandler)
}

func (suite *AdminInfoTestSuite) assertHandlersAreNotNil(api *operations.ConsoleAPI) {
	suite.assert.NotNil(api.SystemAdminInfoHandler)
	suite.assert.NotNil(api.SystemDashboardWidgetDetailsHandler)
}

func (suite *AdminInfoTestSuite) TestSystemAdminInfoHandlerWithError() {
	params, api := suite.initSystemAdminInfoRequest()
	response := api.SystemAdminInfoHandler.Handle(params, &models.Principal{})
	_, ok := response.(*systemApi.AdminInfoDefault)
	suite.assert.True(ok)
}

func (suite *AdminInfoTestSuite) initSystemAdminInfoRequest() (params systemApi.AdminInfoParams, api operations.ConsoleAPI) {
	registerAdminInfoHandlers(&api)
	params.HTTPRequest = &http.Request{}
	defaultOnly := false
	params.DefaultOnly = &defaultOnly
	return params, api
}

func (suite *AdminInfoTestSuite) TestSystemDashboardWidgetDetailsHandlerWithError() {
	params, api := suite.initSystemDashboardWidgetDetailsRequest()
	response := api.SystemDashboardWidgetDetailsHandler.Handle(params, &models.Principal{})
	_, ok := response.(*systemApi.DashboardWidgetDetailsDefault)
	suite.assert.True(ok)
}

func (suite *AdminInfoTestSuite) initSystemDashboardWidgetDetailsRequest() (params systemApi.DashboardWidgetDetailsParams, api operations.ConsoleAPI) {
	registerAdminInfoHandlers(&api)
	params.HTTPRequest = &http.Request{}
	return params, api
}

func (suite *AdminInfoTestSuite) TestGetUsageWidgetsForDeploymentWithoutError() {
	ctx := context.WithValue(context.Background(), utils.ContextClientIP, "127.0.0.1")
	suite.isPrometheusRequest = true
	res, err := getUsageWidgetsForDeployment(ctx, suite.server.URL, suite.adminClient)
	suite.assert.Nil(err)
	suite.assert.NotNil(res)
	suite.isPrometheusRequest = false
}

func (suite *AdminInfoTestSuite) TestGetWidgetDetailsWithoutError() {
	ctx := context.WithValue(context.Background(), utils.ContextClientIP, "127.0.0.1")
	suite.isPrometheusRequest = true
	var step int32 = 1
	var start int64
	var end int64 = 1
	res, err := getWidgetDetails(ctx, suite.server.URL, "mock", 1, &step, &start, &end)
	suite.assert.Nil(err)
	suite.assert.NotNil(res)
	suite.isPrometheusRequest = false
}

func TestAdminInfo(t *testing.T) {
	suite.Run(t, new(AdminInfoTestSuite))
}
