// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/openstor/console/api/operations"
	inspectApi "github.com/openstor/console/api/operations/inspect"
	"github.com/openstor/console/models"
	"github.com/openstor/madmin-go/v4"
	"github.com/secure-io/sio-go"
)

func registerInspectHandler(api *operations.ConsoleAPI) {
	api.InspectInspectHandler = inspectApi.InspectHandlerFunc(func(params inspectApi.InspectParams, principal *models.Principal) middleware.Responder {
		k, r, err := getInspectResult(principal, &params)
		if err != nil {
			return inspectApi.NewInspectDefault(err.Code).WithPayload(err.APIError)
		}

		return middleware.ResponderFunc(processInspectResponse(&params, k, r))
	})
}

func getInspectResult(session *models.Principal, params *inspectApi.InspectParams) ([]byte, io.ReadCloser, *CodedAPIError) {
	ctx := params.HTTPRequest.Context()
	mAdmin, err := NewMinioAdminClient(params.HTTPRequest.Context(), session)
	if err != nil {
		return nil, nil, ErrorWithContext(ctx, err)
	}

	cfg := madmin.InspectOptions{
		File:   params.File,
		Volume: params.Volume,
	}

	// TODO: Remove encryption option and always encrypt.
	// Maybe also add public key field.
	if params.Encrypt != nil && *params.Encrypt {
		cfg.PublicKey, _ = base64.StdEncoding.DecodeString("MIIBCgKCAQEAs/128UFS9A8YSJY1XqYKt06dLVQQCGDee69T+0Tip/1jGAB4z0/3QMpH0MiS8Wjs4BRWV51qvkfAHzwwdU7y6jxU05ctb/H/WzRj3FYdhhHKdzear9TLJftlTs+xwj2XaADjbLXCV1jGLS889A7f7z5DgABlVZMQd9BjVAR8ED3xRJ2/ZCNuQVJ+A8r7TYPGMY3wWvhhPgPk3Lx4WDZxDiDNlFs4GQSaESSsiVTb9vyGe/94CsCTM6Cw9QG6ifHKCa/rFszPYdKCabAfHcS3eTr0GM+TThSsxO7KfuscbmLJkfQev1srfL2Ii2RbnysqIJVWKEwdW05ID8ryPkuTuwIDAQAB")
	}

	// create a MinIO Admin Client interface implementation
	// defining the client to be used
	adminClient := AdminClient{Client: mAdmin}

	k, r, err := adminClient.inspect(ctx, cfg)
	if err != nil {
		return nil, nil, ErrorWithContext(ctx, err)
	}
	return k, r, nil
}

// borrowed from mc cli
func decryptInspectV1(key [32]byte, r io.Reader) io.ReadCloser {
	stream, err := sio.AES_256_GCM.Stream(key[:])
	if err != nil {
		return nil
	}
	nonce := make([]byte, stream.NonceSize())
	return io.NopCloser(stream.DecryptReader(r, nonce, nil))
}

func processInspectResponse(params *inspectApi.InspectParams, k []byte, r io.ReadCloser) func(w http.ResponseWriter, _ runtime.Producer) {
	isEnc := params.Encrypt != nil && *params.Encrypt
	return func(w http.ResponseWriter, _ runtime.Producer) {
		ext := "enc"
		if len(k) == 32 && !isEnc {
			ext = "zip"
			r = decryptInspectV1(*(*[32]byte)(k), r)
		}
		fileName := fmt.Sprintf("inspect-%s-%s.%s", params.Volume, params.File, ext)
		fileName = strings.Map(func(r rune) rune {
			switch {
			case r >= 'A' && r <= 'Z':
				return r
			case r >= 'a' && r <= 'z':
				return r
			case r >= '0' && r <= '9':
				return r
			default:
				if strings.ContainsAny(string(r), "-+._") {
					return r
				}
				return '_'
			}
		}, fileName)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

		_, err := io.Copy(w, r)
		if err != nil {
			LogError("unable to write all the data: %v", err)
		}
	}
}
