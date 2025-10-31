// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	accountApi "github.com/openstor/console/api/operations/account"
	"github.com/openstor/console/models"
	"github.com/stretchr/testify/assert"
)

func Test_getChangePasswordResponse(t *testing.T) {
	assert := assert.New(t)
	session := &models.Principal{
		AccountAccessKey: "TESTTEST",
	}
	CurrentSecretKey := "string"
	NewSecretKey := "string"
	changePasswordParameters := accountApi.AccountChangePasswordParams{
		HTTPRequest: &http.Request{},
		Body: &models.AccountChangePasswordRequest{
			CurrentSecretKey: &CurrentSecretKey,
			NewSecretKey:     &NewSecretKey,
		},
	}
	loginResponse, actualError := getChangePasswordResponse(session, changePasswordParameters)
	expected := (*models.LoginResponse)(nil)
	assert.Equal(expected, loginResponse)
	expectedError := "error please check your current password" // errChangePassword
	assert.Equal(expectedError, actualError.APIError.DetailedMessage)
}

func Test_changePassword(t *testing.T) {
	client := AdminClientMock{}
	type args struct {
		ctx              context.Context
		client           AdminClientMock
		session          *models.Principal
		currentSecretKey string
		newSecretKey     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "password changed successfully",
			args: args{
				client: client,
				ctx:    context.Background(),
				session: &models.Principal{
					AccountAccessKey: "TESTTEST",
				},
				currentSecretKey: "TESTTEST",
				newSecretKey:     "TESTTEST2",
			},
			mock: func() {
				minioChangePasswordMock = func(_ context.Context, _, _ string) error {
					return nil
				}
			},
		},
		{
			name: "error when changing password",
			args: args{
				client: client,
				ctx:    context.Background(),
				session: &models.Principal{
					AccountAccessKey: "TESTTEST",
				},
				currentSecretKey: "TESTTEST",
				newSecretKey:     "TESTTEST2",
			},
			mock: func() {
				minioChangePasswordMock = func(_ context.Context, _, _ string) error {
					return errors.New("there was an error, please try again")
				}
			},
			wantErr: true,
		},
		{
			name: "error because current password doesn't match",
			args: args{
				client: client,
				ctx:    context.Background(),
				session: &models.Principal{
					AccountAccessKey: "TESTTEST",
				},
				currentSecretKey: "TESTTEST",
				newSecretKey:     "TESTTEST2",
			},
			mock: func() {
				minioChangePasswordMock = func(_ context.Context, _, _ string) error {
					return errors.New("there was an error, please try again")
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}
			if err := changePassword(tt.args.ctx, tt.args.client, tt.args.session, tt.args.newSecretKey); (err != nil) != tt.wantErr {
				t.Errorf("changePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
