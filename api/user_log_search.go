// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/openstor/console/api/operations"
	logApi "github.com/openstor/console/api/operations/logging"
	"github.com/openstor/console/models"
	iampolicy "github.com/openstor/pkg/v3/policy"
)

func registerLogSearchHandlers(api *operations.ConsoleAPI) {
	// log search
	api.LoggingLogSearchHandler = logApi.LogSearchHandlerFunc(func(params logApi.LogSearchParams, session *models.Principal) middleware.Responder {
		searchResp, err := getLogSearchResponse(session, params)
		if err != nil {
			return logApi.NewLogSearchDefault(err.Code).WithPayload(err.APIError)
		}
		return logApi.NewLogSearchOK().WithPayload(searchResp)
	})
}

// getLogSearchResponse performs a query to Log Search if Enabled
func getLogSearchResponse(session *models.Principal, params logApi.LogSearchParams) (*models.LogSearchResponse, *CodedAPIError) {
	ctx, cancel := context.WithCancel(params.HTTPRequest.Context())
	defer cancel()
	sessionResp, err := getSessionResponse(ctx, session)
	if err != nil {
		return nil, err
	}
	var allowedToQueryLogSearchAPI bool
	if permissions, ok := sessionResp.Permissions[ConsoleResourceName]; ok {
		for _, permission := range permissions {
			if permission == iampolicy.HealthInfoAdminAction {
				allowedToQueryLogSearchAPI = true
				break
			}
		}
	}

	if !allowedToQueryLogSearchAPI {
		return nil, &CodedAPIError{
			Code: 403,
			APIError: &models.APIError{
				Message:         "Forbidden",
				DetailedMessage: "The Log Search API not available.",
			},
		}
	}

	token := getLogSearchAPIToken()
	endpoint := fmt.Sprintf("%s/api/query?token=%s&q=reqinfo", getLogSearchURL(), token)
	for _, fp := range params.Fp {
		endpoint = fmt.Sprintf("%s&fp=%s", endpoint, fp)
	}

	endpoint = fmt.Sprintf("%s&%s=ok", endpoint, *params.Order)

	// timeStart
	if params.TimeStart != nil && *params.TimeStart != "" {
		endpoint = fmt.Sprintf("%s&timeStart=%s", endpoint, *params.TimeStart)
	}

	// timeEnd
	if params.TimeEnd != nil && *params.TimeEnd != "" {
		endpoint = fmt.Sprintf("%s&timeEnd=%s", endpoint, *params.TimeEnd)
	}

	// page size and page number
	endpoint = fmt.Sprintf("%s&pageSize=%d", endpoint, *params.PageSize)
	endpoint = fmt.Sprintf("%s&pageNo=%d", endpoint, *params.PageNo)

	response, errLogSearch := logSearch(endpoint, getClientIP(params.HTTPRequest))
	if errLogSearch != nil {
		return nil, ErrorWithContext(ctx, errLogSearch)
	}
	return response, nil
}

func logSearch(endpoint string, clientIP string) (*models.LogSearchResponse, error) {
	httpClnt := GetConsoleHTTPClient(clientIP)
	resp, err := httpClnt.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("the Log Search API cannot be reached. Please review the URL and try again %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error retrieving logs: %s", http.StatusText(resp.StatusCode))
	}

	var results []map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	return &models.LogSearchResponse{
		Results: results,
	}, nil
}
