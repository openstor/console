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

	"github.com/openstor/console/api/operations"
	kmsAPI "github.com/openstor/console/api/operations/k_m_s"
	"github.com/openstor/console/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type KMSTestSuite struct {
	suite.Suite
	assert        *assert.Assertions
	currentServer string
	isServerSet   bool
	server        *httptest.Server
	adminClient   AdminClientMock
}

func (suite *KMSTestSuite) SetupSuite() {
	suite.assert = assert.New(suite.T())
	suite.adminClient = AdminClientMock{}
}

func (suite *KMSTestSuite) SetupTest() {
	suite.server = httptest.NewServer(http.HandlerFunc(suite.serverHandler))
	suite.currentServer, suite.isServerSet = os.LookupEnv(ConsoleMinIOServer)
	os.Setenv(ConsoleMinIOServer, suite.server.URL)
}

func (suite *KMSTestSuite) serverHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(400)
}

func (suite *KMSTestSuite) TearDownSuite() {
}

func (suite *KMSTestSuite) TearDownTest() {
	if suite.isServerSet {
		os.Setenv(ConsoleMinIOServer, suite.currentServer)
	} else {
		os.Unsetenv(ConsoleMinIOServer)
	}
}

func (suite *KMSTestSuite) TestRegisterKMSHandlers() {
	api := &operations.ConsoleAPI{}
	suite.assertHandlersAreNil(api)
	registerKMSHandlers(api)
	suite.assertHandlersAreNotNil(api)
}

func (suite *KMSTestSuite) assertHandlersAreNil(api *operations.ConsoleAPI) {
	suite.assert.Nil(api.KmsKMSStatusHandler)
	suite.assert.Nil(api.KmsKMSMetricsHandler)
	suite.assert.Nil(api.KmsKMSAPIsHandler)
	suite.assert.Nil(api.KmsKMSVersionHandler)
	suite.assert.Nil(api.KmsKMSCreateKeyHandler)
	suite.assert.Nil(api.KmsKMSListKeysHandler)
	suite.assert.Nil(api.KmsKMSKeyStatusHandler)
}

func (suite *KMSTestSuite) assertHandlersAreNotNil(api *operations.ConsoleAPI) {
	suite.assert.NotNil(api.KmsKMSStatusHandler)
	suite.assert.NotNil(api.KmsKMSMetricsHandler)
	suite.assert.NotNil(api.KmsKMSAPIsHandler)
	suite.assert.NotNil(api.KmsKMSVersionHandler)
	suite.assert.NotNil(api.KmsKMSCreateKeyHandler)
	suite.assert.NotNil(api.KmsKMSListKeysHandler)
	suite.assert.NotNil(api.KmsKMSKeyStatusHandler)
}

func (suite *KMSTestSuite) TestKMSStatusHandlerWithError() {
	params, api := suite.initKMSStatusRequest()
	response := api.KmsKMSStatusHandler.Handle(params, &models.Principal{})
	_, ok := response.(*kmsAPI.KMSStatusDefault)
	suite.assert.True(ok)
}

func (suite *KMSTestSuite) initKMSStatusRequest() (params kmsAPI.KMSStatusParams, api operations.ConsoleAPI) {
	registerKMSHandlers(&api)
	params.HTTPRequest = &http.Request{}
	return params, api
}

func (suite *KMSTestSuite) TestKMSStatusWithoutError() {
	ctx := context.Background()
	res, err := kmsStatus(ctx, suite.adminClient)
	suite.assert.NotNil(res)
	suite.assert.Nil(err)
}

func (suite *KMSTestSuite) TestKMSMetricsHandlerWithError() {
	params, api := suite.initKMSMetricsRequest()
	response := api.KmsKMSMetricsHandler.Handle(params, &models.Principal{})
	_, ok := response.(*kmsAPI.KMSMetricsDefault)
	suite.assert.True(ok)
}

func (suite *KMSTestSuite) initKMSMetricsRequest() (params kmsAPI.KMSMetricsParams, api operations.ConsoleAPI) {
	registerKMSHandlers(&api)
	params.HTTPRequest = &http.Request{}
	return params, api
}

func (suite *KMSTestSuite) TestKMSMetricsWithoutError() {
	ctx := context.Background()
	res, err := kmsMetrics(ctx, suite.adminClient)
	suite.assert.NotNil(res)
	suite.assert.Nil(err)
}

func (suite *KMSTestSuite) TestKMSAPIsHandlerWithError() {
	params, api := suite.initKMSAPIsRequest()
	response := api.KmsKMSAPIsHandler.Handle(params, &models.Principal{})
	_, ok := response.(*kmsAPI.KMSAPIsDefault)
	suite.assert.True(ok)
}

func (suite *KMSTestSuite) initKMSAPIsRequest() (params kmsAPI.KMSAPIsParams, api operations.ConsoleAPI) {
	registerKMSHandlers(&api)
	params.HTTPRequest = &http.Request{}
	return params, api
}

func (suite *KMSTestSuite) TestKMSAPIsWithoutError() {
	ctx := context.Background()
	res, err := kmsAPIs(ctx, suite.adminClient)
	suite.assert.NotNil(res)
	suite.assert.Nil(err)
}

func (suite *KMSTestSuite) TestKMSVersionHandlerWithError() {
	params, api := suite.initKMSVersionRequest()
	response := api.KmsKMSVersionHandler.Handle(params, &models.Principal{})
	_, ok := response.(*kmsAPI.KMSVersionDefault)
	suite.assert.True(ok)
}

func (suite *KMSTestSuite) initKMSVersionRequest() (params kmsAPI.KMSVersionParams, api operations.ConsoleAPI) {
	registerKMSHandlers(&api)
	params.HTTPRequest = &http.Request{}
	return params, api
}

func (suite *KMSTestSuite) TestKMSVersionWithoutError() {
	ctx := context.Background()
	res, err := kmsVersion(ctx, suite.adminClient)
	suite.assert.NotNil(res)
	suite.assert.Nil(err)
}

func (suite *KMSTestSuite) TestKMSCreateKeyHandlerWithError() {
	params, api := suite.initKMSCreateKeyRequest()
	response := api.KmsKMSCreateKeyHandler.Handle(params, &models.Principal{})
	_, ok := response.(*kmsAPI.KMSCreateKeyDefault)
	suite.assert.True(ok)
}

func (suite *KMSTestSuite) initKMSCreateKeyRequest() (params kmsAPI.KMSCreateKeyParams, api operations.ConsoleAPI) {
	registerKMSHandlers(&api)
	params.HTTPRequest = &http.Request{}
	key := "key"
	params.Body = &models.KmsCreateKeyRequest{Key: &key}
	return params, api
}

func (suite *KMSTestSuite) TestKMSCreateKeyWithoutError() {
	ctx := context.Background()
	err := createKey(ctx, "key", suite.adminClient)
	suite.assert.Nil(err)
}

func (suite *KMSTestSuite) TestKMSListKeysHandlerWithError() {
	params, api := suite.initKMSListKeysRequest()
	response := api.KmsKMSListKeysHandler.Handle(params, &models.Principal{})
	_, ok := response.(*kmsAPI.KMSListKeysDefault)
	suite.assert.True(ok)
}

func (suite *KMSTestSuite) initKMSListKeysRequest() (params kmsAPI.KMSListKeysParams, api operations.ConsoleAPI) {
	registerKMSHandlers(&api)
	params.HTTPRequest = &http.Request{}
	return params, api
}

func (suite *KMSTestSuite) TestKMSListKeysWithoutError() {
	ctx := context.Background()
	res, err := listKeys(ctx, "", suite.adminClient)
	suite.assert.NotNil(res)
	suite.assert.Nil(err)
}

func (suite *KMSTestSuite) TestKMSKeyStatusHandlerWithError() {
	params, api := suite.initKMSKeyStatusRequest()
	response := api.KmsKMSKeyStatusHandler.Handle(params, &models.Principal{})
	_, ok := response.(*kmsAPI.KMSKeyStatusDefault)
	suite.assert.True(ok)
}

func (suite *KMSTestSuite) initKMSKeyStatusRequest() (params kmsAPI.KMSKeyStatusParams, api operations.ConsoleAPI) {
	registerKMSHandlers(&api)
	params.HTTPRequest = &http.Request{}
	return params, api
}

func (suite *KMSTestSuite) TestKMSKeyStatusWithoutError() {
	ctx := context.Background()
	res, err := keyStatus(ctx, "key", suite.adminClient)
	suite.assert.NotNil(res)
	suite.assert.Nil(err)
}

func TestKMS(t *testing.T) {
	suite.Run(t, new(KMSTestSuite))
}
