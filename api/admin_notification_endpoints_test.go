// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/go-openapi/swag"

	cfgApi "github.com/openstor/console/api/operations/configuration"
	"github.com/openstor/console/models"
)

func Test_addNotificationEndpoint(t *testing.T) {
	client := AdminClientMock{}

	type args struct {
		ctx    context.Context
		client MinioAdmin
		params *cfgApi.AddNotificationEndpointParams
	}
	tests := []struct {
		name          string
		args          args
		mockSetConfig func(kv string) (restart bool, err error)
		want          *models.SetNotificationEndpointResponse
		wantErr       bool
	}{
		{
			name: "valid postgres",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"host":     "localhost",
							"user":     "user",
							"password": "passwrd",
						},
						Service: models.NewNofiticationService("postgres"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return false, nil
			},
			want: &models.SetNotificationEndpointResponse{
				AccountID: swag.String("1"),
				Properties: map[string]string{
					"host":     "localhost",
					"user":     "user",
					"password": "passwrd",
				},
				Service: models.NewNofiticationService("postgres"),
				Restart: false,
			},
			wantErr: false,
		},
		{
			name: "set config returns error",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"host":     "localhost",
							"user":     "user",
							"password": "passwrd",
						},
						Service: models.NewNofiticationService("postgres"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return false, errors.New("error")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid mysql",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"host":     "localhost",
							"user":     "user",
							"password": "passwrd",
						},
						Service: models.NewNofiticationService("mysql"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return false, nil
			},
			want: &models.SetNotificationEndpointResponse{
				AccountID: swag.String("1"),
				Properties: map[string]string{
					"host":     "localhost",
					"user":     "user",
					"password": "passwrd",
				},
				Service: models.NewNofiticationService("mysql"),
				Restart: false,
			},
			wantErr: false,
		},
		{
			name: "valid kafka",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"brokers": "http://localhost:8080/broker1",
						},
						Service: models.NewNofiticationService("kafka"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return false, nil
			},
			want: &models.SetNotificationEndpointResponse{
				AccountID: swag.String("1"),
				Properties: map[string]string{
					"brokers": "http://localhost:8080/broker1",
				},
				Service: models.NewNofiticationService("kafka"),
				Restart: false,
			},
			wantErr: false,
		},
		{
			name: "valid amqp",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"url": "http://localhost:8080/broker1",
						},
						Service: models.NewNofiticationService("amqp"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return false, nil
			},
			want: &models.SetNotificationEndpointResponse{
				AccountID: swag.String("1"),
				Properties: map[string]string{
					"url": "http://localhost:8080/broker1",
				},
				Service: models.NewNofiticationService("amqp"),
				Restart: false,
			},
			wantErr: false,
		},
		{
			name: "valid mqtt",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"broker": "http://localhost:8080/broker1",
							"topic":  "minio",
						},
						Service: models.NewNofiticationService("mqtt"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return false, nil
			},
			want: &models.SetNotificationEndpointResponse{
				AccountID: swag.String("1"),
				Properties: map[string]string{
					"broker": "http://localhost:8080/broker1",
					"topic":  "minio",
				},
				Service: models.NewNofiticationService("mqtt"),
				Restart: false,
			},
			wantErr: false,
		},
		{
			name: "valid elasticsearch",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"url":    "http://localhost:8080/broker1",
							"index":  "minio",
							"format": "namespace",
						},
						Service: models.NewNofiticationService("elasticsearch"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return false, nil
			},
			want: &models.SetNotificationEndpointResponse{
				AccountID: swag.String("1"),
				Properties: map[string]string{
					"url":    "http://localhost:8080/broker1",
					"index":  "minio",
					"format": "namespace",
				},
				Service: models.NewNofiticationService("elasticsearch"),
				Restart: false,
			},
			wantErr: false,
		},
		{
			name: "valid redis",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"address": "http://localhost:8080/broker1",
							"key":     "minio",
							"format":  "namespace",
						},
						Service: models.NewNofiticationService("redis"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return false, nil
			},
			want: &models.SetNotificationEndpointResponse{
				AccountID: swag.String("1"),
				Properties: map[string]string{
					"address": "http://localhost:8080/broker1",
					"key":     "minio",
					"format":  "namespace",
				},
				Service: models.NewNofiticationService("redis"),
				Restart: false,
			},
			wantErr: false,
		},
		{
			name: "valid nats",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"address": "http://localhost:8080/broker1",
							"subject": "minio",
						},
						Service: models.NewNofiticationService("nats"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return false, nil
			},
			want: &models.SetNotificationEndpointResponse{
				AccountID: swag.String("1"),
				Properties: map[string]string{
					"address": "http://localhost:8080/broker1",
					"subject": "minio",
				},
				Service: models.NewNofiticationService("nats"),
				Restart: false,
			},
			wantErr: false,
		},
		{
			name: "valid webhook",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"endpoint": "http://localhost:8080/broker1",
						},
						Service: models.NewNofiticationService("webhook"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return false, nil
			},
			want: &models.SetNotificationEndpointResponse{
				AccountID: swag.String("1"),
				Properties: map[string]string{
					"endpoint": "http://localhost:8080/broker1",
				},
				Service: models.NewNofiticationService("webhook"),
				Restart: false,
			},
			wantErr: false,
		},
		{
			name: "valid nsq",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"nsqd_address": "http://localhost:8080/broker1",
							"topic":        "minio",
						},
						Service: models.NewNofiticationService("nsq"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return false, nil
			},
			want: &models.SetNotificationEndpointResponse{
				AccountID: swag.String("1"),
				Properties: map[string]string{
					"nsqd_address": "http://localhost:8080/broker1",
					"topic":        "minio",
				},
				Service: models.NewNofiticationService("nsq"),
				Restart: false,
			},
			wantErr: false,
		},
		{
			name: "invalid service",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"host":     "localhost",
							"user":     "user",
							"password": "passwrd",
						},
						Service: models.NewNofiticationService("oorgle"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return false, errors.New("invalid config")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid config, restart required",
			args: args{
				ctx:    context.Background(),
				client: client,
				params: &cfgApi.AddNotificationEndpointParams{
					HTTPRequest: nil,
					Body: &models.NotificationEndpoint{
						AccountID: swag.String("1"),
						Properties: map[string]string{
							"host":     "localhost",
							"user":     "user",
							"password": "passwrd",
						},
						Service: models.NewNofiticationService("postgres"),
					},
				},
			},
			mockSetConfig: func(_ string) (restart bool, err error) {
				return true, nil
			},
			want: &models.SetNotificationEndpointResponse{
				AccountID: swag.String("1"),
				Properties: map[string]string{
					"host":     "localhost",
					"user":     "user",
					"password": "passwrd",
				},
				Service: models.NewNofiticationService("postgres"),
				Restart: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			// mock function response from setConfig()
			minioSetConfigKVMock = tt.mockSetConfig
			got, err := addNotificationEndpoint(tt.args.ctx, tt.args.client, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("addNotificationEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addNotificationEndpoint() got = %v, want %v", got, tt.want)
			}
		})
	}
}
