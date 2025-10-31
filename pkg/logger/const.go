// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package logger

import (
	"context"

	"github.com/openstor/console/pkg/logger/target/http"
)

// Audit/Logger constants
const (
	EnvLoggerJSONEnable      = "CONSOLE_LOGGER_JSON_ENABLE"
	EnvLoggerAnonymousEnable = "CONSOLE_LOGGER_ANONYMOUS_ENABLE"
	EnvLoggerQuietEnable     = "CONSOLE_LOGGER_QUIET_ENABLE"

	EnvGlobalDeploymentID      = "CONSOLE_GLOBAL_DEPLOYMENT_ID"
	EnvLoggerWebhookEnable     = "CONSOLE_LOGGER_WEBHOOK_ENABLE"
	EnvLoggerWebhookEndpoint   = "CONSOLE_LOGGER_WEBHOOK_ENDPOINT"
	EnvLoggerWebhookAuthToken  = "CONSOLE_LOGGER_WEBHOOK_AUTH_TOKEN"
	EnvLoggerWebhookClientCert = "CONSOLE_LOGGER_WEBHOOK_CLIENT_CERT"
	EnvLoggerWebhookClientKey  = "CONSOLE_LOGGER_WEBHOOK_CLIENT_KEY"
	EnvLoggerWebhookQueueSize  = "CONSOLE_LOGGER_WEBHOOK_QUEUE_SIZE"

	EnvAuditWebhookEnable     = "CONSOLE_AUDIT_WEBHOOK_ENABLE"
	EnvAuditWebhookEndpoint   = "CONSOLE_AUDIT_WEBHOOK_ENDPOINT"
	EnvAuditWebhookAuthToken  = "CONSOLE_AUDIT_WEBHOOK_AUTH_TOKEN"
	EnvAuditWebhookClientCert = "CONSOLE_AUDIT_WEBHOOK_CLIENT_CERT"
	EnvAuditWebhookClientKey  = "CONSOLE_AUDIT_WEBHOOK_CLIENT_KEY"
	EnvAuditWebhookQueueSize  = "CONSOLE_AUDIT_WEBHOOK_QUEUE_SIZE"
)

// Config console and http logger targets
type Config struct {
	HTTP         map[string]http.Config `json:"http"`
	AuditWebhook map[string]http.Config `json:"audit"`
}

var (
	globalDeploymentID string
	GlobalContext      context.Context
)
