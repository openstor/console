// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/openstor/console/pkg/logger"

	"github.com/minio/cli"
	"github.com/openstor/console/api"
)

var appCmds = []cli.Command{
	serverCmd,
	updateCmd,
}

// StartServer starts the console service
func StartServer(ctx *cli.Context) error {
	if err := loadAllCerts(ctx); err != nil {
		// Log this as a warning and continue running console without TLS certificates
		api.LogError("Unable to load certs: %v", err)
	}

	xctx := context.Background()

	transport := api.PrepareSTSClientTransport(api.LocalAddress).Transport.(*http.Transport)
	if err := logger.InitializeLogger(xctx, transport); err != nil {
		fmt.Println("error InitializeLogger", err)
		logger.CriticalIf(xctx, err)
	}
	// custom error configuration
	api.LogInfo = logger.Info
	api.LogError = logger.Error
	api.LogIf = logger.LogIf

	var rctx api.Context
	if err := rctx.Load(ctx); err != nil {
		api.LogError("argument validation failed: %v", err)
		return err
	}

	server, err := buildServer()
	if err != nil {
		api.LogError("Unable to initialize console server: %v", err)
		return err
	}

	server.Host = rctx.Host
	server.Port = rctx.HTTPPort
	// set conservative timesout for uploads
	server.ReadTimeout = 1 * time.Hour
	// no timeouts for response for downloads
	server.WriteTimeout = 0
	api.Port = strconv.Itoa(server.Port)
	api.Hostname = server.Host

	if len(api.GlobalPublicCerts) > 0 {
		// If TLS certificates are provided enforce the HTTPS schema, meaning console will redirect
		// plain HTTP connections to HTTPS server
		server.EnabledListeners = []string{"http", "https"}
		server.TLSPort = rctx.HTTPSPort
		// Need to store tls-port, tls-host un config variables so secure.middleware can read from there
		api.TLSPort = strconv.Itoa(server.TLSPort)
		api.Hostname = rctx.Host
		api.TLSRedirect = rctx.TLSRedirect
	}

	defer server.Shutdown()

	if err = server.Serve(); err != nil {
		server.Logf("error serving API: %v", err)
		return err
	}

	return nil
}
