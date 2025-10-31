// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"net/http"

	authApi "github.com/openstor/console/api/operations/auth"

	"github.com/openstor/console/pkg/auth"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/openstor/console/api/operations"
	accountApi "github.com/openstor/console/api/operations/account"
	"github.com/openstor/console/models"
)

func registerAccountHandlers(api *operations.ConsoleAPI) {
	// change user password
	api.AccountAccountChangePasswordHandler = accountApi.AccountChangePasswordHandlerFunc(func(params accountApi.AccountChangePasswordParams, session *models.Principal) middleware.Responder {
		changePasswordResponse, err := getChangePasswordResponse(session, params)
		if err != nil {
			return accountApi.NewAccountChangePasswordDefault(err.Code).WithPayload(err.APIError)
		}
		// Custom response writer to update the session cookies
		return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
			cookie := NewSessionCookieForConsole(changePasswordResponse.SessionID)
			http.SetCookie(w, &cookie)
			authApi.NewLoginNoContent().WriteResponse(w, p)
		})
	})
}

// changePassword validate current current user password and if it's correct set the new password
func changePassword(ctx context.Context, client MinioAdmin, session *models.Principal, newSecretKey string) error {
	return client.changePassword(ctx, session.AccountAccessKey, newSecretKey)
}

// getChangePasswordResponse will validate user knows what is the current password (avoid account hijacking), update user account password
// and authenticate the user generating a new session token/cookie
func getChangePasswordResponse(session *models.Principal, params accountApi.AccountChangePasswordParams) (*models.LoginResponse, *CodedAPIError) {
	ctx, cancel := context.WithCancel(params.HTTPRequest.Context())
	defer cancel()
	clientIP := getClientIP(params.HTTPRequest)
	client := GetConsoleHTTPClient(clientIP)

	// changePassword operations requires an AdminClient initialized with parent account credentials not
	// STS credentials
	parentAccountClient, err := NewMinioAdminClient(params.HTTPRequest.Context(), &models.Principal{
		STSAccessKeyID:     session.AccountAccessKey,
		STSSecretAccessKey: *params.Body.CurrentSecretKey,
	})
	if err != nil {
		return nil, ErrorWithContext(ctx, err)
	}
	// parentAccountClient will contain access and secret key credentials for the user
	userClient := AdminClient{Client: parentAccountClient}
	accessKey := session.AccountAccessKey
	newSecretKey := *params.Body.NewSecretKey

	// currentSecretKey will compare currentSecretKey against the stored secret key inside the encrypted session
	if err := changePassword(ctx, userClient, session, newSecretKey); err != nil {
		return nil, ErrorWithContext(ctx, ErrChangePassword, nil, err)
	}
	// user credentials are updated at this point, we need to generate a new admin client and authenticate using
	// the new credentials
	credentials, err := getConsoleCredentials(accessKey, newSecretKey, client)
	if err != nil {
		return nil, ErrorWithContext(ctx, ErrInvalidLogin, nil, err)
	}
	// authenticate user and generate new session token
	sessionID, err := login(credentials, &auth.SessionFeatures{HideMenu: session.Hm})
	if err != nil {
		return nil, ErrorWithContext(ctx, ErrInvalidLogin, nil, err)
	}
	// serialize output
	loginResponse := &models.LoginResponse{
		SessionID: *sessionID,
	}
	return loginResponse, nil
}
