// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
	mc "github.com/openstor/mc/cmd"
)

type watchOptions struct {
	BucketName string
	mc.WatchOptions
}

func startWatch(ctx context.Context, conn WSConn, wsc MCClient, options *watchOptions) error {
	wo, pErr := wsc.watch(ctx, options.WatchOptions)
	if pErr != nil {
		LogError("error initializing watch: %v", pErr.Cause)
		return pErr.Cause
	}
	for {
		select {
		case <-ctx.Done():
			close(wo.DoneChan)
			return nil
		case events, ok := <-wo.Events():
			// zero value returned because the channel is closed and empty
			if !ok {
				return nil
			}
			for _, event := range events {
				// Serialize message to be sent
				bytes, err := json.Marshal(event)
				if err != nil {
					LogError("error on json.Marshal: %v", err)
					return err
				}
				// Send Message through websocket connection
				err = conn.writeMessage(websocket.TextMessage, bytes)
				if err != nil {
					LogError("error writeMessage: %v", err)
					return err
				}
			}
		case pErr, ok := <-wo.Errors():
			// zero value returned because the channel is closed and empty
			if !ok {
				return nil
			}
			if pErr != nil {
				LogError("error on watch: %v", pErr.Cause)
				return pErr.Cause

			}
		}
	}
}

// getWatchOptionsFromReq gets bucket name, events, prefix, suffix from a websocket
// watch path if defined.
// path come as : `/watch/bucket1` and query
// params come on request form
func getWatchOptionsFromReq(req *http.Request) (*watchOptions, error) {
	wOptions := watchOptions{}
	// Default Events if not defined
	wOptions.Events = []string{"put", "get", "delete"}

	re := regexp.MustCompile(`(/watch/)(.*?$)`)
	matches := re.FindAllSubmatch([]byte(req.URL.Path), -1)
	// matches comes as e.g.
	// [["...", "/watch/", "bucket1"]]
	// [["/watch/" "/watch/" ""]]

	if len(matches) == 0 || len(matches[0]) < 3 {
		return nil, fmt.Errorf("invalid url: %s", req.URL.Path)
	}

	wOptions.BucketName = strings.TrimSpace(string(matches[0][2]))

	events := req.FormValue("events")
	if strings.TrimSpace(events) != "" {
		wOptions.Events = strings.Split(events, ",")
	}
	wOptions.Prefix = req.FormValue("prefix")
	wOptions.Suffix = req.FormValue("suffix")
	return &wOptions, nil
}
