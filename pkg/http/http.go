// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package http

import (
	"io"
	"net/http"
)

// ClientI interface with all functions to be implemented
// by mock when testing, it should include all HttpClient respective api calls
// that are used within this project.
type ClientI interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	Do(req *http.Request) (*http.Response, error)
}

// Client is an HTTP Interface implementation
//
// Define the structure of a http client and define the functions that are actually used
type Client struct {
	Client *http.Client
}

// Get implements http.Client.Get()
func (c *Client) Get(url string) (resp *http.Response, err error) {
	return c.Client.Get(url)
}

// Post implements http.Client.Post()
func (c *Client) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return c.Client.Post(url, contentType, body)
}

// Do implement http.Client.Do()
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

// DrainBody close non nil response with any response Body.
// convenient wrapper to drain any remaining data on response body.
//
// Subsequently this allows golang http RoundTripper
// to re-use the same connection for future requests.
func DrainBody(respBody io.ReadCloser) {
	// Callers should close resp.Body when done reading from it.
	// If resp.Body is not closed, the Client's underlying RoundTripper
	// (typically Transport) may not be able to re-use a persistent TCP
	// connection to the server for a subsequent "keep-alive" request.
	if respBody != nil {
		// Drain any remaining Body and then close the connection.
		// Without this closing connection would disallow re-using
		// the same connection for future uses.
		//  - http://stackoverflow.com/a/17961593/4465767
		defer respBody.Close()
		io.Copy(io.Discard, respBody)
	}
}
