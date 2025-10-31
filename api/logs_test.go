// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"flag"
	"fmt"
	"testing"

	"github.com/minio/cli"
	"github.com/stretchr/testify/assert"
)

func TestContext_Load(t *testing.T) {
	type fields struct {
		Host           string
		HTTPPort       int
		HTTPSPort      int
		TLSRedirect    string
		TLSCertificate string
		TLSKey         string
		TLSca          string
	}
	type args struct {
		values map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid args",
			args: args{
				values: map[string]string{
					"tls-redirect": "on",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid args",
			args: args{
				values: map[string]string{
					"tls-redirect": "aaaa",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid port http",
			args: args{
				values: map[string]string{
					"tls-redirect": "on",
					"port":         "65536",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid port https",
			args: args{
				values: map[string]string{
					"tls-redirect": "on",
					"port":         "65534",
					"tls-port":     "65536",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			c := &Context{}

			fs := flag.NewFlagSet("flags", flag.ContinueOnError)
			for k, v := range tt.args.values {
				fs.String(k, v, "ok")
			}

			ctx := cli.NewContext(nil, fs, &cli.Context{})

			err := c.Load(ctx)
			if tt.wantErr {
				assert.NotNilf(t, err, fmt.Sprintf("Load(%v)", err))
			} else {
				assert.Nilf(t, err, fmt.Sprintf("Load(%v)", err))
			}
		})
	}
}

func Test_logInfo(_ *testing.T) {
	logInfo("message", nil)
}
