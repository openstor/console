// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package oauth2

// Environment constants for console IDP/SSO configuration
const (
	ConsoleMinIOServer           = "CONSOLE_MINIO_SERVER"
	ConsoleIDPURL                = "CONSOLE_IDP_URL"
	ConsoleIDPClientID           = "CONSOLE_IDP_CLIENT_ID"
	ConsoleIDPSecret             = "CONSOLE_IDP_SECRET"
	ConsoleIDPCallbackURL        = "CONSOLE_IDP_CALLBACK"
	ConsoleIDPCallbackURLDynamic = "CONSOLE_IDP_CALLBACK_DYNAMIC"
	ConsoleIDPHmacPassphrase     = "CONSOLE_IDP_HMAC_PASSPHRASE"
	ConsoleIDPHmacSalt           = "CONSOLE_IDP_HMAC_SALT"
	ConsoleIDPScopes             = "CONSOLE_IDP_SCOPES"
	ConsoleIDPUserInfo           = "CONSOLE_IDP_USERINFO"
	ConsoleIDPTokenExpiration    = "CONSOLE_IDP_TOKEN_EXPIRATION"
)
