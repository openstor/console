// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package utils

import (
	"context"
)

// Key used for Get/SetReqInfo
type key string

const (
	ContextLogKey            = key("console-log")
	ContextRequestID         = key("request-id")
	ContextRequestUserID     = key("request-user-id")
	ContextRequestUserAgent  = key("request-user-agent")
	ContextRequestHost       = key("request-host")
	ContextRequestRemoteAddr = key("request-remote-addr")
	ContextAuditKey          = key("request-audit-entry")
	ContextClientIP          = key("client-ip")
)

// ClientIPFromContext attempts to get the Client IP from a context, if it's not present, it returns
// 127.0.0.1
func ClientIPFromContext(ctx context.Context) string {
	val := ctx.Value(ContextClientIP)
	if val != nil {
		return val.(string)
	}
	return "127.0.0.1"
}
