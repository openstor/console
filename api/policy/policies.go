// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package policy

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/openstor/madmin-go/v4"
)

// ReplacePolicyVariables replaces known variables from policies with known values
func ReplacePolicyVariables(claims map[string]interface{}, accountInfo *madmin.AccountInfo) json.RawMessage {
	// AWS Variables
	rawPolicy := bytes.ReplaceAll(accountInfo.Policy, []byte("${aws:username}"), []byte(accountInfo.AccountName))
	rawPolicy = bytes.ReplaceAll(rawPolicy, []byte("${aws:userid}"), []byte(accountInfo.AccountName))
	// JWT Variables
	rawPolicy = replaceJwtVariables(rawPolicy, claims)
	// LDAP Variables
	rawPolicy = replaceLDAPVariables(rawPolicy, claims)
	return rawPolicy
}

func replaceJwtVariables(rawPolicy []byte, claims map[string]interface{}) json.RawMessage {
	// list of valid JWT fields we will replace in policy if they are in the response
	jwtFields := []string{
		"sub",
		"iss",
		"aud",
		"jti",
		"upn",
		"name",
		"groups",
		"given_name",
		"family_name",
		"middle_name",
		"nickname",
		"preferred_username",
		"profile",
		"picture",
		"website",
		"email",
		"gender",
		"birthdate",
		"phone_number",
		"address",
		"scope",
		"client_id",
	}
	// check which fields are in the claims and replace as variable by casting the value to string
	for _, field := range jwtFields {
		if val, ok := claims[field]; ok {
			variable := fmt.Sprintf("${jwt:%s}", field)
			rawPolicy = bytes.ReplaceAll(rawPolicy, []byte(variable), []byte(fmt.Sprintf("%v", val)))
		}
	}
	return rawPolicy
}

// ReplacePolicyVariables replaces known variables from policies with known values
func replaceLDAPVariables(rawPolicy []byte, claims map[string]interface{}) json.RawMessage {
	// replace ${ldap:user}
	if val, ok := claims["ldapUser"]; ok {
		rawPolicy = bytes.ReplaceAll(rawPolicy, []byte("${ldap:user}"), []byte(fmt.Sprintf("%v", val)))
	}
	// replace ${ldap:username}
	if val, ok := claims["ldapUsername"]; ok {
		rawPolicy = bytes.ReplaceAll(rawPolicy, []byte("${ldap:username}"), []byte(fmt.Sprintf("%v", val)))
	}
	return rawPolicy
}
