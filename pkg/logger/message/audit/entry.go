// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package audit

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v4"

	"github.com/openstor/console/pkg/utils"

	xhttp "github.com/openstor/console/pkg/http"
)

// Version - represents the current version of audit log structure.
const Version = "1"

// ObjectVersion object version key/versionId
type ObjectVersion struct {
	ObjectName string `json:"objectName"`
	VersionID  string `json:"versionId,omitempty"`
}

// Entry - audit entry logs.
type Entry struct {
	Version      string    `json:"version"`
	DeploymentID string    `json:"deploymentid,omitempty"`
	Time         time.Time `json:"time"`
	Trigger      string    `json:"trigger"`
	API          struct {
		Path            string `json:"path,omitempty"`
		Status          string `json:"status,omitempty"`
		Method          string `json:"method"`
		StatusCode      int    `json:"statusCode,omitempty"`
		InputBytes      int64  `json:"rx"`
		OutputBytes     int64  `json:"tx"`
		TimeToFirstByte string `json:"timeToFirstByte,omitempty"`
		TimeToResponse  string `json:"timeToResponse,omitempty"`
	} `json:"api"`
	RemoteHost string                 `json:"remotehost,omitempty"`
	RequestID  string                 `json:"requestID,omitempty"`
	SessionID  string                 `json:"sessionID,omitempty"`
	UserAgent  string                 `json:"userAgent,omitempty"`
	ReqClaims  map[string]interface{} `json:"requestClaims,omitempty"`
	ReqQuery   map[string]string      `json:"requestQuery,omitempty"`
	ReqHeader  map[string]string      `json:"requestHeader,omitempty"`
	RespHeader map[string]string      `json:"responseHeader,omitempty"`
	Tags       map[string]interface{} `json:"tags,omitempty"`
}

// NewEntry - constructs an audit entry object with some fields filled
func NewEntry(deploymentID string) Entry {
	return Entry{
		Version:      Version,
		DeploymentID: deploymentID,
		Time:         time.Now().UTC(),
	}
}

// ToEntry - constructs an audit entry from a http request
func ToEntry(w http.ResponseWriter, r *http.Request, reqClaims map[string]interface{}, deploymentID string) Entry {
	entry := NewEntry(deploymentID)

	entry.RemoteHost = r.RemoteAddr
	entry.UserAgent = r.UserAgent()
	entry.ReqClaims = reqClaims

	q := r.URL.Query()
	reqQuery := make(map[string]string, len(q))
	for k, v := range q {
		reqQuery[k] = strings.Join(v, ",")
	}
	entry.ReqQuery = reqQuery

	reqHeader := make(map[string]string, len(r.Header))
	for k, v := range r.Header {
		reqHeader[k] = strings.Join(v, ",")
	}
	entry.ReqHeader = reqHeader

	wh := w.Header()

	var requestID interface{}
	requestID = r.Context().Value(utils.ContextRequestID)
	if requestID == nil {
		requestID = uuid.NewString()
	}
	entry.RequestID = requestID.(string)

	if val := r.Context().Value(utils.ContextRequestUserID); val != nil {
		sessionID := val.(string)
		if os.Getenv("CONSOLE_OPERATOR_MODE") != "" && os.Getenv("CONSOLE_OPERATOR_MODE") == "on" {
			claims := jwt.MapClaims{}
			_, _ = jwt.ParseWithClaims(sessionID, claims, nil)
			if sub, ok := claims["sub"]; ok {
				sessionID = sub.(string)
			}
		}
		entry.SessionID = sessionID
	}

	respHeader := make(map[string]string, len(wh))
	for k, v := range wh {
		respHeader[k] = strings.Join(v, ",")
	}
	entry.RespHeader = respHeader

	if etag := respHeader[xhttp.ETag]; etag != "" {
		respHeader[xhttp.ETag] = strings.Trim(etag, `"`)
	}

	return entry
}
