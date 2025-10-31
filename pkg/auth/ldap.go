// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package auth

import (
	"net/http"

	"github.com/openstor/openstor-go/v7/pkg/credentials"
)

// GetCredentialsFromLDAP authenticates the user against MinIO when the LDAP integration is enabled
// if the authentication succeed *credentials.Login object is returned and we continue with the normal STSAssumeRole flow
func GetCredentialsFromLDAP(client *http.Client, endpoint, ldapUser, ldapPassword string) (*credentials.Credentials, error) {
	creds := credentials.New(&credentials.LDAPIdentity{
		Client:       client,
		STSEndpoint:  endpoint,
		LDAPUsername: ldapUser,
		LDAPPassword: ldapPassword,
	})
	return creds, nil
}
