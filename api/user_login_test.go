// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"errors"
	"reflect"
	"testing"

	xoauth2 "golang.org/x/oauth2"

	"github.com/openstor/madmin-go/v4"

	iampolicy "github.com/openstor/pkg/v3/policy"

	"github.com/openstor/console/pkg/auth"

	"github.com/openstor/openstor-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
)

// Define a mock struct of ConsoleCredentialsI interface implementation
type consoleCredentialsMock struct{}

func (ac consoleCredentialsMock) GetActions() []string {
	return []string{}
}

func (ac consoleCredentialsMock) GetAccountAccessKey() string {
	return ""
}

// Common mocks
var consoleCredentialsGetMock func() (credentials.Value, error)

// mock function of Get()
func (ac consoleCredentialsMock) Get() (credentials.Value, error) {
	return consoleCredentialsGetMock()
}

func TestLogin(t *testing.T) {
	funcAssert := assert.New(t)
	consoleCredentials := consoleCredentialsMock{}
	// Test Case 1: Valid consoleCredentials
	consoleCredentialsGetMock = func() (credentials.Value, error) {
		return credentials.Value{
			AccessKeyID:     "fakeAccessKeyID",
			SecretAccessKey: "fakeSecretAccessKey",
			SessionToken:    "fakeSessionToken",
			SignerType:      0,
		}, nil
	}
	token, err := login(consoleCredentials, nil)
	funcAssert.NotEmpty(token, "Token was returned empty")
	funcAssert.Nil(err, "error creating a session")

	// Test Case 2: Invalid credentials
	consoleCredentialsGetMock = func() (credentials.Value, error) {
		return credentials.Value{}, errors.New("")
	}
	_, err = login(consoleCredentials, nil)
	funcAssert.NotNil(err, "not error returned creating a session")
}

type IdentityProviderMock struct{}

var (
	idpVerifyIdentityMock            func(ctx context.Context, code, state string) (*credentials.Credentials, error)
	idpVerifyIdentityForOperatorMock func(ctx context.Context, code, state string) (*xoauth2.Token, error)
	idpGenerateLoginURLMock          func() string
)

func (ac IdentityProviderMock) VerifyIdentity(ctx context.Context, code, state string) (*credentials.Credentials, error) {
	return idpVerifyIdentityMock(ctx, code, state)
}

func (ac IdentityProviderMock) VerifyIdentityForOperator(ctx context.Context, code, state string) (*xoauth2.Token, error) {
	return idpVerifyIdentityForOperatorMock(ctx, code, state)
}

func (ac IdentityProviderMock) GenerateLoginURL() string {
	return idpGenerateLoginURLMock()
}

func Test_validateUserAgainstIDP(t *testing.T) {
	provider := IdentityProviderMock{}
	mockCode := "EAEAEAE"
	mockState := "HUEHUEHUE"
	type args struct {
		ctx      context.Context
		provider auth.IdentityProviderI
		code     string
		state    string
	}
	tests := []struct {
		name     string
		args     args
		want     *credentials.Credentials
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "failed to verify user identity with idp",
			args: args{
				ctx:      context.Background(),
				provider: provider,
				code:     mockCode,
				state:    mockState,
			},
			want:    nil,
			wantErr: true,
			mockFunc: func() {
				idpVerifyIdentityMock = func(_ context.Context, _, _ string) (*credentials.Credentials, error) {
					return nil, errors.New("something went wrong")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc()
			}
			got, err := verifyUserAgainstIDP(tt.args.ctx, tt.args.provider, tt.args.code, tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyUserAgainstIDP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("verifyUserAgainstIDP() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAccountInfo(t *testing.T) {
	client := AdminClientMock{}
	type args struct {
		ctx    context.Context
		client MinioAdmin
	}
	tests := []struct {
		name     string
		args     args
		want     *iampolicy.Policy
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "error getting account info",
			args: args{
				ctx:    context.Background(),
				client: client,
			},
			want:    nil,
			wantErr: true,
			mockFunc: func() {
				minioAccountInfoMock = func(_ context.Context) (madmin.AccountInfo, error) {
					return madmin.AccountInfo{}, errors.New("something went wrong")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc()
			}
			got, err := getAccountInfo(tt.args.ctx, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAccountInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("getAccountInfo() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
