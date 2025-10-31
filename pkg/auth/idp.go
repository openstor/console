// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package auth

import (
	"context"

	"github.com/openstor/console/pkg/auth/idp/oauth2"
	"github.com/openstor/openstor-go/v7/pkg/credentials"
	xoauth2 "golang.org/x/oauth2"
)

// IdentityProviderI interface with all functions to be implemented
// by mock when testing, it should include all IdentityProvider respective api calls
// that are used within this project.
type IdentityProviderI interface {
	VerifyIdentity(ctx context.Context, code, state string) (*credentials.Credentials, error)
	VerifyIdentityForOperator(ctx context.Context, code, state string) (*xoauth2.Token, error)
	GenerateLoginURL() string
}

// Interface implementation
//
// Define the structure of a IdentityProvider with Client inside and define the functions that are used
// during the authentication flow.
type IdentityProvider struct {
	KeyFunc oauth2.StateKeyFunc
	Client  *oauth2.Provider
	RoleARN string
}

// VerifyIdentity will verify the user identity against the idp using the authorization code flow
func (c IdentityProvider) VerifyIdentity(ctx context.Context, code, state string) (*credentials.Credentials, error) {
	return c.Client.VerifyIdentity(ctx, code, state, c.RoleARN, c.KeyFunc)
}

// VerifyIdentityForOperator will verify the user identity against the idp using the authorization code flow
func (c IdentityProvider) VerifyIdentityForOperator(ctx context.Context, code, state string) (*xoauth2.Token, error) {
	return c.Client.VerifyIdentityForOperator(ctx, code, state, c.KeyFunc)
}

// GenerateLoginURL returns a new URL used by the user to login against the idp
func (c IdentityProvider) GenerateLoginURL() string {
	return c.Client.GenerateLoginURL(c.KeyFunc, c.Client.IDPName)
}
