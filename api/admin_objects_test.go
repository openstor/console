// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"testing"
	"time"

	mc "github.com/openstor/mc/cmd"
	"github.com/openstor/openstor-go/v7"
	"github.com/stretchr/testify/assert"
)

func TestWSRewindObjects(t *testing.T) {
	assert := assert.New(t)
	client := s3ClientMock{}

	tests := []struct {
		name         string
		testOptions  objectsListOpts
		testMessages []*mc.ClientContent
	}{
		{
			name: "Get list with multiple elements",
			testOptions: objectsListOpts{
				BucketName: "buckettest",
				Prefix:     "/",
				Date:       time.Now(),
			},
			testMessages: []*mc.ClientContent{
				{
					BucketName: "buckettest",
					URL:        mc.ClientURL{Path: "/file1.txt"},
				},
				{
					BucketName: "buckettest",
					URL:        mc.ClientURL{Path: "/file2.txt"},
				},
				{
					BucketName: "buckettest",
					URL:        mc.ClientURL{Path: "/path1"},
				},
			},
		},
		{
			name: "Empty list of elements",
			testOptions: objectsListOpts{
				BucketName: "emptybucket",
				Prefix:     "/",
				Date:       time.Now(),
			},
			testMessages: []*mc.ClientContent{},
		},
		{
			name: "Get list with one element",
			testOptions: objectsListOpts{
				BucketName: "buckettest",
				Prefix:     "/",
				Date:       time.Now(),
			},
			testMessages: []*mc.ClientContent{
				{
					BucketName: "buckettestsingle",
					URL:        mc.ClientURL{Path: "/file12.txt"},
				},
			},
		},
		{
			name: "Get data from subpaths",
			testOptions: objectsListOpts{
				BucketName: "buckettest",
				Prefix:     "/path1/path2",
				Date:       time.Now(),
			},
			testMessages: []*mc.ClientContent{
				{
					BucketName: "buckettestsingle",
					URL:        mc.ClientURL{Path: "/path1/path2/file12.txt"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mcListMock = func(_ context.Context, _ mc.ListOptions) <-chan *mc.ClientContent {
				ch := make(chan *mc.ClientContent)
				go func() {
					defer close(ch)
					for _, m := range tt.testMessages {
						ch <- m
					}
				}()
				return ch
			}

			rewindList := startRewindListing(ctx, client, &tt.testOptions)

			// check that the rewindList got the same number of data from Console.

			totalItems := 0
			for data := range rewindList {
				// Compare elements as we are defining the channel responses
				assert.Equal(tt.testMessages[totalItems].URL.Path, data.URL.Path)
				totalItems++
			}
			assert.Equal(len(tt.testMessages), totalItems)
		})
	}
}

func TestWSListObjects(t *testing.T) {
	assert := assert.New(t)
	client := minioClientMock{}

	tests := []struct {
		name         string
		wantErr      bool
		testOptions  objectsListOpts
		testMessages []openstor.ObjectInfo
	}{
		{
			name:    "Get list with multiple elements",
			wantErr: false,
			testOptions: objectsListOpts{
				BucketName: "buckettest",
				Prefix:     "/",
			},
			testMessages: []openstor.ObjectInfo{
				{
					Key:          "/file1.txt",
					Size:         500,
					IsLatest:     true,
					LastModified: time.Now(),
				},
				{
					Key:          "/file2.txt",
					Size:         500,
					IsLatest:     true,
					LastModified: time.Now(),
				},
				{
					Key: "/path1",
				},
			},
		},
		{
			name:    "Empty list of elements",
			wantErr: false,
			testOptions: objectsListOpts{
				BucketName: "emptybucket",
				Prefix:     "/",
			},
			testMessages: []openstor.ObjectInfo{},
		},
		{
			name:    "Get list with one element",
			wantErr: false,
			testOptions: objectsListOpts{
				BucketName: "buckettest",
				Prefix:     "/",
			},
			testMessages: []openstor.ObjectInfo{
				{
					Key:          "/file2.txt",
					Size:         500,
					IsLatest:     true,
					LastModified: time.Now(),
				},
			},
		},
		{
			name:    "Get data from subpaths",
			wantErr: false,
			testOptions: objectsListOpts{
				BucketName: "buckettest",
				Prefix:     "/path1/path2",
			},
			testMessages: []openstor.ObjectInfo{
				{
					Key:          "/path1/path2/file1.txt",
					Size:         500,
					IsLatest:     true,
					LastModified: time.Now(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			minioListObjectsMock = func(_ context.Context, _ string, _ openstor.ListObjectsOptions) <-chan openstor.ObjectInfo {
				ch := make(chan openstor.ObjectInfo)
				go func() {
					defer close(ch)
					for _, m := range tt.testMessages {
						ch <- m
					}
				}()
				return ch
			}

			objectsListing := startObjectsListing(ctx, client, &tt.testOptions)

			// check that the TestReceiver got the same number of data from Console
			totalItems := 0
			for data := range objectsListing {
				// Compare elements as we are defining the channel responses
				assert.Equal(tt.testMessages[totalItems].Key, data.Key)
				totalItems++
			}
			assert.Equal(len(tt.testMessages), totalItems)
		})
	}
}
