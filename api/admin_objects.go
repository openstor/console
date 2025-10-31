// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"time"

	"github.com/openstor/mc/cmd"
	"github.com/openstor/openstor-go/v7"
)

type objectsListOpts struct {
	BucketName string
	Prefix     string
	Date       time.Time
}

type ObjectsRequest struct {
	Mode       string `json:"mode,omitempty"`
	BucketName string `json:"bucket_name"`
	Prefix     string `json:"prefix"`
	Date       string `json:"date"`
	RequestID  int64  `json:"request_id"`
}

type WSResponse struct {
	RequestID  int64            `json:"request_id,omitempty"`
	Error      *CodedAPIError   `json:"error,omitempty"`
	RequestEnd bool             `json:"request_end,omitempty"`
	Prefix     string           `json:"prefix,omitempty"`
	BucketName string           `json:"bucketName,omitempty"`
	Data       []ObjectResponse `json:"data,omitempty"`
}

type ObjectResponse struct {
	Name         string `json:"name,omitempty"`
	LastModified string `json:"last_modified,omitempty"`
	Size         int64  `json:"size,omitempty"`
	VersionID    string `json:"version_id,omitempty"`
	DeleteMarker bool   `json:"delete_flag,omitempty"`
	IsLatest     bool   `json:"is_latest,omitempty"`
}

func getObjectsOptionsFromReq(request ObjectsRequest) (*objectsListOpts, error) {
	pOptions := objectsListOpts{
		BucketName: request.BucketName,
		Prefix:     request.Prefix,
	}

	if request.Mode == "rewind" {
		parsedDate, errDate := time.Parse(time.RFC3339, request.Date)

		if errDate != nil {
			return nil, errDate
		}

		pOptions.Date = parsedDate
	}

	return &pOptions, nil
}

func startObjectsListing(ctx context.Context, client MinioClient, objOpts *objectsListOpts) <-chan openstor.ObjectInfo {
	opts := openstor.ListObjectsOptions{
		Prefix: objOpts.Prefix,
	}

	return client.listObjects(ctx, objOpts.BucketName, opts)
}

func startRewindListing(ctx context.Context, client MCClient, objOpts *objectsListOpts) <-chan *cmd.ClientContent {
	lsRewind := client.list(ctx, cmd.ListOptions{TimeRef: objOpts.Date, WithDeleteMarkers: true})

	return lsRewind
}
