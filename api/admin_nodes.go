// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/openstor/console/api/operations"
	systemApi "github.com/openstor/console/api/operations/system"
	"github.com/openstor/console/models"
)

func registerNodesHandler(api *operations.ConsoleAPI) {
	api.SystemListNodesHandler = systemApi.ListNodesHandlerFunc(func(params systemApi.ListNodesParams, session *models.Principal) middleware.Responder {
		listNodesResponse, err := getListNodesResponse(session, params)
		if err != nil {
			return systemApi.NewListNodesDefault(err.Code).WithPayload(err.APIError)
		}
		return systemApi.NewListNodesOK().WithPayload(listNodesResponse)
	})
}

// getListNodesResponse returns a list of available node endpoints .
func getListNodesResponse(session *models.Principal, params systemApi.ListNodesParams) ([]string, *CodedAPIError) {
	ctx, cancel := context.WithCancel(params.HTTPRequest.Context())
	defer cancel()
	mAdmin, err := NewMinioAdminClient(params.HTTPRequest.Context(), session)
	if err != nil {
		return nil, ErrorWithContext(ctx, err)
	}
	var nodeList []string

	adminResources, _ := mAdmin.ServerInfo(ctx)

	for _, n := range adminResources.Servers {
		nodeList = append(nodeList, n.Endpoint)
	}

	return nodeList, nil
}
