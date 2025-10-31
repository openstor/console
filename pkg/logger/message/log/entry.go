// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package log

import (
	"time"
)

// ObjectVersion object version key/versionId
type ObjectVersion struct {
	ObjectName string `json:"objectName"`
	VersionID  string `json:"versionId,omitempty"`
}

// Args - defines the arguments for the API.
type Args struct {
	Bucket    string            `json:"bucket,omitempty"`
	Object    string            `json:"object,omitempty"`
	VersionID string            `json:"versionId,omitempty"`
	Objects   []ObjectVersion   `json:"objects,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Trace - defines the trace.
type Trace struct {
	Message   string                 `json:"message,omitempty"`
	Source    []string               `json:"source,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// API - defines the api type and its args.
type API struct {
	Name string `json:"name,omitempty"`
}

// Entry - defines fields and values of each log entry.
type Entry struct {
	DeploymentID string    `json:"deploymentid,omitempty"`
	Level        string    `json:"level"`
	LogKind      string    `json:"errKind"`
	Time         time.Time `json:"time"`
	API          *API      `json:"api,omitempty"`
	RemoteHost   string    `json:"remotehost,omitempty"`
	Host         string    `json:"host,omitempty"`
	RequestID    string    `json:"requestID,omitempty"`
	SessionID    string    `json:"sessionID,omitempty"`
	UserAgent    string    `json:"userAgent,omitempty"`
	Message      string    `json:"message,omitempty"`
	Trace        *Trace    `json:"errors,omitempty"`
}
