// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"

	systemApi "github.com/openstor/console/api/operations/system"

	"github.com/go-openapi/runtime/middleware"
	"github.com/openstor/console/api/operations"
	"github.com/openstor/console/models"
)

func registerAdminArnsHandlers(api *operations.ConsoleAPI) {
	// return a list of arns
	api.SystemArnListHandler = systemApi.ArnListHandlerFunc(func(params systemApi.ArnListParams, session *models.Principal) middleware.Responder {
		arnsResp, err := getArnsResponse(session, params)
		if err != nil {
			return systemApi.NewArnListDefault(err.Code).WithPayload(err.APIError)
		}
		return systemApi.NewArnListOK().WithPayload(arnsResp)
	})
}

// getArns invokes admin info and returns a list of arns
func getArns(ctx context.Context, client MinioAdmin) (*models.ArnsResponse, error) {
	serverInfo, err := client.serverInfo(ctx)
	if err != nil {
		return nil, err
	}
	// build response
	return &models.ArnsResponse{
		Arns: serverInfo.SQSARN,
	}, nil
}

// getArnsResponse returns a list of active arns in the instance
func getArnsResponse(session *models.Principal, params systemApi.ArnListParams) (*models.ArnsResponse, *CodedAPIError) {
	ctx, cancel := context.WithCancel(params.HTTPRequest.Context())
	defer cancel()
	mAdmin, err := NewMinioAdminClient(params.HTTPRequest.Context(), session)
	if err != nil {
		return nil, ErrorWithContext(ctx, err)
	}
	// create a minioClient interface implementation
	// defining the client to be used
	adminClient := AdminClient{Client: mAdmin}

	// serialize output
	arnsList, err := getArns(ctx, adminClient)
	if err != nil {
		return nil, ErrorWithContext(ctx, err)
	}
	return arnsList, nil
}
