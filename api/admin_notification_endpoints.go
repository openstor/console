// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"errors"

	"github.com/go-openapi/runtime/middleware"
	"github.com/openstor/console/api/operations"
	configurationApi "github.com/openstor/console/api/operations/configuration"
	"github.com/openstor/console/models"
)

func registerAdminNotificationEndpointsHandlers(api *operations.ConsoleAPI) {
	// return a list of notification endpoints
	api.ConfigurationNotificationEndpointListHandler = configurationApi.NotificationEndpointListHandlerFunc(func(params configurationApi.NotificationEndpointListParams, session *models.Principal) middleware.Responder {
		notifEndpoints, err := getNotificationEndpointsResponse(session, params)
		if err != nil {
			return configurationApi.NewNotificationEndpointListDefault(err.Code).WithPayload(err.APIError)
		}
		return configurationApi.NewNotificationEndpointListOK().WithPayload(notifEndpoints)
	})
	// add a new notification endpoints
	api.ConfigurationAddNotificationEndpointHandler = configurationApi.AddNotificationEndpointHandlerFunc(func(params configurationApi.AddNotificationEndpointParams, session *models.Principal) middleware.Responder {
		notifEndpoints, err := getAddNotificationEndpointResponse(session, params)
		if err != nil {
			return configurationApi.NewAddNotificationEndpointDefault(err.Code).WithPayload(err.APIError)
		}
		return configurationApi.NewAddNotificationEndpointCreated().WithPayload(notifEndpoints)
	})
}

// getNotificationEndpoints invokes admin info and returns a list of notification endpoints
func getNotificationEndpoints(ctx context.Context, client MinioAdmin) (*models.NotifEndpointResponse, error) {
	serverInfo, err := client.serverInfo(ctx)
	if err != nil {
		return nil, err
	}
	var listEndpoints []*models.NotificationEndpointItem
	for i := range serverInfo.Services.Notifications {
		for service, endpointStatus := range serverInfo.Services.Notifications[i] {
			for j := range endpointStatus {
				for account, status := range endpointStatus[j] {
					listEndpoints = append(listEndpoints, &models.NotificationEndpointItem{
						Service:   models.NofiticationService(service),
						AccountID: account,
						Status:    status.Status,
					})
				}
			}
		}
	}

	// build response
	return &models.NotifEndpointResponse{
		NotificationEndpoints: listEndpoints,
	}, nil
}

// getNotificationEndpointsResponse returns a list of notification endpoints in the instance
func getNotificationEndpointsResponse(session *models.Principal, params configurationApi.NotificationEndpointListParams) (*models.NotifEndpointResponse, *CodedAPIError) {
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
	notfEndpointResp, err := getNotificationEndpoints(ctx, adminClient)
	if err != nil {
		return nil, ErrorWithContext(ctx, err)
	}
	return notfEndpointResp, nil
}

func addNotificationEndpoint(ctx context.Context, client MinioAdmin, params *configurationApi.AddNotificationEndpointParams) (*models.SetNotificationEndpointResponse, error) {
	configs := []*models.ConfigurationKV{}
	var configName string

	// we have different add validations for each service
	switch *params.Body.Service {
	case models.NofiticationServiceAmqp:
		configName = "notify_amqp"
	case models.NofiticationServiceMqtt:
		configName = "notify_mqtt"
	case models.NofiticationServiceElasticsearch:
		configName = "notify_elasticsearch"
	case models.NofiticationServiceRedis:
		configName = "notify_redis"
	case models.NofiticationServiceNats:
		configName = "notify_nats"
	case models.NofiticationServicePostgres:
		configName = "notify_postgres"
	case models.NofiticationServiceMysql:
		configName = "notify_mysql"
	case models.NofiticationServiceKafka:
		configName = "notify_kafka"
	case models.NofiticationServiceWebhook:
		configName = "notify_webhook"
	case models.NofiticationServiceNsq:
		configName = "notify_nsq"
	default:
		return nil, errors.New("provided service is not supported")
	}

	// set all the config values if found on the param.Body.Properties
	for k, val := range params.Body.Properties {
		configs = append(configs, &models.ConfigurationKV{
			Key:   k,
			Value: val,
		})
	}

	needsRestart, err := setConfigWithARNAccountID(ctx, client, &configName, configs, *params.Body.AccountID)
	if err != nil {
		return nil, err
	}

	return &models.SetNotificationEndpointResponse{
		AccountID:  params.Body.AccountID,
		Properties: params.Body.Properties,
		Service:    params.Body.Service,
		Restart:    needsRestart,
	}, nil
}

// getNotificationEndpointsResponse returns a list of notification endpoints in the instance
func getAddNotificationEndpointResponse(session *models.Principal, params configurationApi.AddNotificationEndpointParams) (*models.SetNotificationEndpointResponse, *CodedAPIError) {
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
	notfEndpointResp, err := addNotificationEndpoint(ctx, adminClient, &params)
	if err != nil {
		return nil, ErrorWithContext(ctx, err)
	}
	return notfEndpointResp, nil
}
