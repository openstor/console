// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/openstor/console/api/operations"
	"github.com/openstor/console/models"

	svcApi "github.com/openstor/console/api/operations/service"
)

func registerServiceHandlers(api *operations.ConsoleAPI) {
	// Restart Service
	api.ServiceRestartServiceHandler = svcApi.RestartServiceHandlerFunc(func(params svcApi.RestartServiceParams, session *models.Principal) middleware.Responder {
		if err := getRestartServiceResponse(session, params); err != nil {
			return svcApi.NewRestartServiceDefault(err.Code).WithPayload(err.APIError)
		}
		return svcApi.NewRestartServiceNoContent()
	})
}

// serviceRestart - restarts the MinIO cluster
func serviceRestart(ctx context.Context, client MinioAdmin) error {
	if err := client.serviceRestart(ctx); err != nil {
		return err
	}
	// copy behavior from minio/mc mainAdminServiceRestart()
	//
	// Max. time taken by the server to shutdown is 5 seconds.
	// This can happen when there are lot of s3 requests pending when the server
	// receives a restart command.
	// Sleep for 6 seconds and then check if the server is online.
	time.Sleep(6 * time.Second)

	// Fetch the service status of the specified MinIO server
	_, err := client.serverInfo(ctx)
	if err != nil {
		return err
	}
	return nil
}

// getRestartServiceResponse performs serviceRestart()
func getRestartServiceResponse(session *models.Principal, params svcApi.RestartServiceParams) *CodedAPIError {
	ctx, cancel := context.WithCancel(params.HTTPRequest.Context())
	defer cancel()
	mAdmin, err := NewMinioAdminClient(params.HTTPRequest.Context(), session)
	if err != nil {
		return ErrorWithContext(ctx, err)
	}
	// create a MinIO Admin Client interface implementation
	// defining the client to be used
	adminClient := AdminClient{Client: mAdmin}

	if err := serviceRestart(ctx, adminClient); err != nil {
		return ErrorWithContext(ctx, err)
	}
	return nil
}
