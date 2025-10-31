// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/openstor/console/api/operations"
	release "github.com/openstor/console/api/operations/release"
	"github.com/openstor/console/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReleasesTestSuite struct {
	suite.Suite
	assert        *assert.Assertions
	currentServer string
	isServerSet   bool
	getServer     *httptest.Server
	withError     bool
}

func (suite *ReleasesTestSuite) SetupSuite() {
	suite.assert = assert.New(suite.T())
	suite.getServer = httptest.NewServer(http.HandlerFunc(suite.getHandler))
	suite.currentServer, suite.isServerSet = os.LookupEnv(releaseServiceHostEnvVar)
	os.Setenv(releaseServiceHostEnvVar, suite.getServer.URL)
}

func (suite *ReleasesTestSuite) TearDownSuite() {
	if suite.isServerSet {
		os.Setenv(releaseServiceHostEnvVar, suite.currentServer)
	} else {
		os.Unsetenv(releaseServiceHostEnvVar)
	}
}

func (suite *ReleasesTestSuite) getHandler(
	w http.ResponseWriter, _ *http.Request,
) {
	if suite.withError {
		w.WriteHeader(400)
	} else {
		w.WriteHeader(200)
		response := &models.ReleaseListResponse{}
		bytes, _ := json.Marshal(response)
		fmt.Fprint(w, string(bytes))
	}
}

func (suite *ReleasesTestSuite) TestRegisterReleasesHandlers() {
	api := &operations.ConsoleAPI{}
	suite.assert.Nil(api.ReleaseListReleasesHandler)
	registerReleasesHandlers(api)
	suite.assert.NotNil(api.ReleaseListReleasesHandler)
}

func (suite *ReleasesTestSuite) TestGetReleasesWithError() {
	api := &operations.ConsoleAPI{}
	current := "mock"
	registerReleasesHandlers(api)
	params := release.NewListReleasesParams()
	params.Current = &current
	params.HTTPRequest = &http.Request{}
	suite.withError = true
	response := api.ReleaseListReleasesHandler.Handle(params, &models.Principal{})
	_, ok := response.(*release.ListReleasesDefault)
	suite.assert.True(ok)
}

func (suite *ReleasesTestSuite) TestGetReleasesWithoutError() {
	api := &operations.ConsoleAPI{}
	registerReleasesHandlers(api)
	params := release.NewListReleasesParams()
	params.HTTPRequest = &http.Request{}
	suite.withError = false
	response := api.ReleaseListReleasesHandler.Handle(params, &models.Principal{})
	_, ok := response.(*release.ListReleasesOK)
	suite.assert.True(ok)
}

func TestReleases(t *testing.T) {
	suite.Run(t, new(ReleasesTestSuite))
}
