// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openstor/madmin-go/v4"
)

type profileOptions struct {
	Types    string
	Duration time.Duration
}

func getProfileOptionsFromReq(req *http.Request) (*profileOptions, error) {
	pOptions := profileOptions{}
	pOptions.Types = req.FormValue("types")
	pOptions.Duration = 10 * time.Second // TODO: make this configurable
	return &pOptions, nil
}

func startProfiling(ctx context.Context, conn WSConn, client MinioAdmin, pOpts *profileOptions) error {
	data, err := client.startProfiling(ctx, madmin.ProfilerType(pOpts.Types), pOpts.Duration)
	if err != nil {
		return err
	}
	message, err := io.ReadAll(data)
	if err != nil {
		return err
	}
	return conn.writeMessage(websocket.BinaryMessage, message)
}
