// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"net/http"
)

type ConsoleTransport struct {
	Transport http.RoundTripper
	ClientIP  string
}

func (t *ConsoleTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.ClientIP != "" {
		// Do not set an empty x-forwarded-for
		req.Header.Add(xForwardedFor, t.ClientIP)
	}
	return t.Transport.RoundTrip(req)
}

// PrepareSTSClientTransport :
func PrepareSTSClientTransport(clientIP string) *ConsoleTransport {
	return &ConsoleTransport{
		Transport: GlobalTransport,
		ClientIP:  clientIP,
	}
}

// PrepareConsoleHTTPClient returns an http.Client with custom configurations need it by *credentials.STSAssumeRole
// custom configurations include the use of CA certificates
func PrepareConsoleHTTPClient(clientIP string) *http.Client {
	// Return http client with default configuration
	return &http.Client{
		Transport: PrepareSTSClientTransport(clientIP),
	}
}
