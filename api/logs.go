// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/minio/cli"
)

var (
	infoLog  = log.New(os.Stdout, "I: ", log.LstdFlags)
	errorLog = log.New(os.Stdout, "E: ", log.LstdFlags)
)

func logInfo(msg string, data ...interface{}) {
	infoLog.Printf(msg+"\n", data...)
}

func logError(msg string, data ...interface{}) {
	errorLog.Printf(msg+"\n", data...)
}

func logIf(_ context.Context, _ error, _ ...interface{}) {
}

// globally changeable logger styles
var (
	LogInfo  = logInfo
	LogError = logError
	LogIf    = logIf
)

// Context captures all command line flags values
type Context struct {
	Host                string
	HTTPPort, HTTPSPort int
	TLSRedirect         string
	// Legacy options, TODO: remove in future
	TLSCertificate, TLSKey, TLSca string
}

// Load loads api Context from command line context.
func (c *Context) Load(ctx *cli.Context) error {
	*c = Context{
		Host:        ctx.String("host"),
		HTTPPort:    ctx.Int("port"),
		HTTPSPort:   ctx.Int("tls-port"),
		TLSRedirect: ctx.String("tls-redirect"),
		// Legacy options to be removed.
		TLSCertificate: ctx.String("tls-certificate"),
		TLSKey:         ctx.String("tls-key"),
		TLSca:          ctx.String("tls-ca"),
	}
	if c.HTTPPort > 65535 {
		return errors.New("invalid argument --port out of range - ports can range from 1-65535")
	}
	if c.HTTPSPort > 65535 {
		return errors.New("invalid argument --tls-port out of range - ports can range from 1-65535")
	}
	if c.TLSRedirect != "on" && c.TLSRedirect != "off" {
		return errors.New("invalid argument --tls-redirect only accepts either 'on' or 'off'")
	}
	return nil
}
