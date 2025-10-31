// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openstor/console/models"
)

func (wsc *wsMinioClient) objectManager(session *models.Principal) {
	// Storage of Cancel Contexts for this connection
	var cancelContexts sync.Map
	// Initial goroutine
	defer func() {
		// We close socket at the end of requests
		wsc.conn.close()
		cancelContexts.Range(func(key, value interface{}) bool {
			cancelFunc := value.(context.CancelFunc)
			cancelFunc()

			cancelContexts.Delete(key)
			return true
		})
	}()

	writeChannel := make(chan WSResponse)
	done := make(chan interface{})

	sendWSResponse := func(r WSResponse) {
		select {
		case writeChannel <- r:
		case <-done:
		}
	}

	// Read goroutine
	go func() {
		defer close(writeChannel)

		for {
			select {
			case <-done:
				return
			default:
			}

			mType, message, err := wsc.conn.readMessage()
			if err != nil {
				LogInfo("Error while reading objectManager message: %s", err)
				return
			}

			if mType == websocket.TextMessage {
				// We get request data & review information
				var messageRequest ObjectsRequest

				if err := json.Unmarshal(message, &messageRequest); err != nil {
					LogInfo("Error on message request unmarshal: %s", err)
					continue
				}

				// new message, new context
				ctx, cancel := context.WithCancel(context.Background())

				// We store the cancel func associated with this request
				cancelContexts.Store(messageRequest.RequestID, cancel)

				switch messageRequest.Mode {
				case "objects", "rewind":
					// cancel all previous open objects requests for listing
					cancelContexts.Range(func(key, value interface{}) bool {
						rid := key.(int64)
						if rid < messageRequest.RequestID {
							cancelFunc := value.(context.CancelFunc)
							cancelFunc()

							cancelContexts.Delete(key)
						}
						return true
					})
				}

				const itemsPerBatch = 1000
				switch messageRequest.Mode {
				case "close":
					return
				case "cancel":
					// if we have that request id, cancel it
					if cancelFunc, ok := cancelContexts.Load(messageRequest.RequestID); ok {
						cancelFunc.(context.CancelFunc)()
						cancelContexts.Delete(messageRequest.RequestID)
					}
				case "objects":
					// start listing and writing to web socket
					objectRqConfigs, err := getObjectsOptionsFromReq(messageRequest)
					if err != nil {
						LogInfo(fmt.Sprintf("Error during Objects OptionsParse %s", err.Error()))

						sendWSResponse(WSResponse{
							RequestID:  messageRequest.RequestID,
							Error:      ErrorWithContext(ctx, err),
							Prefix:     messageRequest.Prefix,
							BucketName: messageRequest.BucketName,
						})

						return
					}

					var buffer []ObjectResponse
					for lsObj := range startObjectsListing(ctx, wsc.client, objectRqConfigs) {
						if lsObj.Err != nil {
							sendWSResponse(WSResponse{
								RequestID:  messageRequest.RequestID,
								Error:      ErrorWithContext(ctx, lsObj.Err),
								Prefix:     messageRequest.Prefix,
								BucketName: messageRequest.BucketName,
							})

							continue
						}
						// if the key is same as requested prefix it would be nested directory object, so skip
						// and show only objects under the prefix
						// E.g:
						// bucket/prefix1/prefix2/ -- this should be skipped from list item.
						// bucket/prefix1/prefix2/an-object
						// bucket/prefix1/prefix2/another-object
						if messageRequest.Prefix != lsObj.Key {
							objItem := ObjectResponse{
								Name:         lsObj.Key,
								Size:         lsObj.Size,
								LastModified: lsObj.LastModified.Format(time.RFC3339),
								VersionID:    lsObj.VersionID,
								IsLatest:     lsObj.IsLatest,
								DeleteMarker: lsObj.IsDeleteMarker,
							}
							buffer = append(buffer, objItem)
						}

						if len(buffer) >= itemsPerBatch {
							sendWSResponse(WSResponse{
								RequestID: messageRequest.RequestID,
								Data:      buffer,
							})
							buffer = nil
						}
					}
					if len(buffer) > 0 {
						sendWSResponse(WSResponse{
							RequestID: messageRequest.RequestID,
							Data:      buffer,
						})
					}

					sendWSResponse(WSResponse{
						RequestID:  messageRequest.RequestID,
						RequestEnd: true,
					})

					// if we have that request id, cancel it
					if cancelFunc, ok := cancelContexts.Load(messageRequest.RequestID); ok {
						cancelFunc.(context.CancelFunc)()
						cancelContexts.Delete(messageRequest.RequestID)
					}
				case "rewind":
					// start listing and writing to web socket
					objectRqConfigs, err := getObjectsOptionsFromReq(messageRequest)
					if err != nil {
						LogInfo(fmt.Sprintf("Error during Objects OptionsParse %s", err.Error()))
						sendWSResponse(WSResponse{
							RequestID:  messageRequest.RequestID,
							Error:      ErrorWithContext(ctx, err),
							Prefix:     messageRequest.Prefix,
							BucketName: messageRequest.BucketName,
						})

						return
					}

					clientIP := wsc.conn.remoteAddress()

					s3Client, err := newS3BucketClient(session, objectRqConfigs.BucketName, objectRqConfigs.Prefix, clientIP)
					if err != nil {
						sendWSResponse(WSResponse{
							RequestID:  messageRequest.RequestID,
							Error:      ErrorWithContext(ctx, err),
							Prefix:     messageRequest.Prefix,
							BucketName: messageRequest.BucketName,
						})

						return
					}

					mcS3C := mcClient{client: s3Client}

					var buffer []ObjectResponse

					for lsObj := range startRewindListing(ctx, mcS3C, objectRqConfigs) {
						if lsObj.Err != nil {
							sendWSResponse(WSResponse{
								RequestID:  messageRequest.RequestID,
								Error:      ErrorWithContext(ctx, lsObj.Err.ToGoError()),
								Prefix:     messageRequest.Prefix,
								BucketName: messageRequest.BucketName,
							})

							continue
						}

						name := strings.Replace(lsObj.URL.Path, fmt.Sprintf("/%s/", objectRqConfigs.BucketName), "", 1)

						objItem := ObjectResponse{
							Name:         name,
							Size:         lsObj.Size,
							LastModified: lsObj.Time.Format(time.RFC3339),
							VersionID:    lsObj.VersionID,
							IsLatest:     lsObj.IsLatest,
							DeleteMarker: lsObj.IsDeleteMarker,
						}
						buffer = append(buffer, objItem)

						if len(buffer) >= itemsPerBatch {
							sendWSResponse(WSResponse{
								RequestID: messageRequest.RequestID,
								Data:      buffer,
							})
							buffer = nil
						}

					}
					if len(buffer) > 0 {
						sendWSResponse(WSResponse{
							RequestID: messageRequest.RequestID,
							Data:      buffer,
						})
					}

					sendWSResponse(WSResponse{
						RequestID:  messageRequest.RequestID,
						RequestEnd: true,
					})

					// if we have that request id, cancel it
					if cancelFunc, ok := cancelContexts.Load(messageRequest.RequestID); ok {
						cancelFunc.(context.CancelFunc)()
						cancelContexts.Delete(messageRequest.RequestID)
					}
				}
			}
		}
	}()

	defer close(done)

	for writeM := range writeChannel {
		jsonData, err := json.Marshal(writeM)
		if err != nil {
			LogInfo("Error while marshaling the response: %s", err)
			return
		}

		err = wsc.conn.writeMessage(websocket.TextMessage, jsonData)
		if err != nil {
			LogInfo("Error while writing the message: %s", err)
			return
		}
	}
}
