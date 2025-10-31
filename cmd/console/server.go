// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/minio/cli"
	"github.com/openstor/console/api"
	"github.com/openstor/console/api/operations"
	"github.com/openstor/console/pkg/certs"
)

// starts the server
var serverCmd = cli.Command{
	Name:    "server",
	Aliases: []string{"srv"},
	Usage:   "Start OpenStor Console server",
	Action:  StartServer,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: api.GetHostname(),
			Usage: "bind to a specific HOST, HOST can be an IP or hostname",
		},
		cli.IntFlag{
			Name:  "port",
			Value: api.GetPort(),
			Usage: "bind to specific HTTP port",
		},
		// This is kept here for backward compatibility,
		// hostname's do not have HTTP or HTTPs
		// hostnames are opaque so using --host
		// works for both HTTP and HTTPS setup.
		cli.StringFlag{
			Name:   "tls-host",
			Value:  api.GetHostname(),
			Hidden: true,
		},
		cli.StringFlag{
			Name:  "certs-dir",
			Value: certs.GlobalCertsCADir.Get(),
			Usage: "path to certs directory",
		},
		cli.IntFlag{
			Name:  "tls-port",
			Value: api.GetTLSPort(),
			Usage: "bind to specific HTTPS port",
		},
		cli.StringFlag{
			Name:  "tls-redirect",
			Value: api.GetTLSRedirect(),
			Usage: "toggle HTTP->HTTPS redirect",
		},
		cli.StringFlag{
			Name:   "tls-certificate",
			Value:  "",
			Usage:  "path to TLS public certificate",
			Hidden: true,
		},
		cli.StringFlag{
			Name:   "tls-key",
			Value:  "",
			Usage:  "path to TLS private key",
			Hidden: true,
		},
		cli.StringFlag{
			Name:   "tls-ca",
			Value:  "",
			Usage:  "path to TLS Certificate Authority",
			Hidden: true,
		},
	},
}

func buildServer() (*api.Server, error) {
	swaggerSpec, err := loads.Embedded(api.SwaggerJSON, api.FlatSwaggerJSON)
	if err != nil {
		return nil, err
	}

	consoleAPI := operations.NewConsoleAPI(swaggerSpec)
	consoleAPI.Logger = api.LogInfo
	server := api.NewServer(consoleAPI)

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "OpenStor Console Server"
	parser.LongDescription = swaggerSpec.Spec().Info.Description

	server.ConfigureFlags()

	// register all APIs
	server.ConfigureAPI()

	for _, optsGroup := range consoleAPI.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			return nil, err
		}
	}

	if _, err := parser.Parse(); err != nil {
		return nil, err
	}

	return server, nil
}

func loadAllCerts(ctx *cli.Context) error {
	var err error
	// Set all certs and CAs directories path
	certs.GlobalCertsDir, _, err = certs.NewConfigDirFromCtx(ctx, "certs-dir", certs.DefaultCertsDir.Get)
	if err != nil {
		return err
	}

	certs.GlobalCertsCADir = &certs.ConfigDir{Path: filepath.Join(certs.GlobalCertsDir.Get(), certs.CertsCADir)}
	// check if certs and CAs directories exists or can be created
	if err = certs.MkdirAllIgnorePerm(certs.GlobalCertsCADir.Get()); err != nil {
		return fmt.Errorf("unable to create certs CA directory at %s: failed with %w", certs.GlobalCertsCADir.Get(), err)
	}

	// load the certificates and the CAs
	api.GlobalRootCAs, api.GlobalPublicCerts, api.GlobalTLSCertsManager, err = certs.GetAllCertificatesAndCAs()
	if err != nil {
		return fmt.Errorf("unable to load certificates at %s: failed with %w", certs.GlobalCertsDir.Get(), err)
	}

	{
		// TLS flags from swagger server, used to support VMware vsphere operator version.
		swaggerServerCertificate := ctx.String("tls-certificate")
		swaggerServerCertificateKey := ctx.String("tls-key")
		swaggerServerCACertificate := ctx.String("tls-ca")
		// load tls cert and key from swagger server tls-certificate and tls-key flags
		if swaggerServerCertificate != "" && swaggerServerCertificateKey != "" {
			if err = api.GlobalTLSCertsManager.AddCertificate(swaggerServerCertificate, swaggerServerCertificateKey); err != nil {
				return err
			}
			x509Certs, err := certs.ParsePublicCertFile(swaggerServerCertificate)
			if err == nil {
				api.GlobalPublicCerts = append(api.GlobalPublicCerts, x509Certs...)
			}
		}

		// load ca cert from swagger server tls-ca flag
		if swaggerServerCACertificate != "" {
			caCert, caCertErr := os.ReadFile(swaggerServerCACertificate)
			if caCertErr == nil {
				api.GlobalRootCAs.AppendCertsFromPEM(caCert)
			}
		}
	}

	if api.GlobalTLSCertsManager != nil {
		api.GlobalTLSCertsManager.ReloadOnSignal(syscall.SIGHUP)
	}

	return nil
}
