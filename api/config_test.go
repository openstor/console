// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostname(t *testing.T) {
	os.Setenv(ConsoleHostname, "x")
	defer os.Unsetenv(ConsoleHostname)
	assert.Equalf(t, "x", GetHostname(), "GetHostname()")
}

func TestGetPort(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "valid port",
			args: args{
				env: "9091",
			},
			want: 9091,
		},
		{
			name: "invalid port",
			args: args{
				env: "duck",
			},
			want: 9090,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			os.Setenv(ConsolePort, tt.args.env)
			assert.Equalf(t, tt.want, GetPort(), "GetPort()")
			os.Unsetenv(ConsolePort)
		})
	}
}

func TestGetTLSPort(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "valid port",
			args: args{
				env: "9444",
			},
			want: 9444,
		},
		{
			name: "invalid port",
			args: args{
				env: "duck",
			},
			want: 9443,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			os.Setenv(ConsoleTLSPort, tt.args.env)
			assert.Equalf(t, tt.want, GetTLSPort(), "GetTLSPort()")
			os.Unsetenv(ConsoleTLSPort)
		})
	}
}

func TestGetSecureAllowedHosts(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid hosts",
			args: args{
				env: "host1,host2",
			},
			want: []string{"host1", "host2"},
		},
		{
			name: "empty hosts",
			args: args{
				env: "",
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			os.Setenv(ConsoleSecureAllowedHosts, tt.args.env)
			assert.Equalf(t, tt.want, GetSecureAllowedHosts(), "GetSecureAllowedHosts()")
			os.Unsetenv(ConsoleSecureAllowedHosts)
		})
	}
}

func TestGetSecureHostsProxyHeaders(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid headers",
			args: args{
				env: "header1,header2",
			},
			want: []string{"header1", "header2"},
		},
		{
			name: "empty headers",
			args: args{
				env: "",
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			os.Setenv(ConsoleSecureHostsProxyHeaders, tt.args.env)
			assert.Equalf(t, tt.want, GetSecureHostsProxyHeaders(), "GetSecureHostsProxyHeaders()")
			os.Unsetenv(ConsoleSecureHostsProxyHeaders)
		})
	}
}

func TestGetSecureSTSSeconds(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "valid",
			args: args{
				env: "1",
			},
			want: 1,
		},
		{
			name: "invalid",
			args: args{
				env: "duck",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			os.Setenv(ConsoleSecureSTSSeconds, tt.args.env)
			assert.Equalf(t, tt.want, GetSecureSTSSeconds(), "GetSecureSTSSeconds()")
			os.Unsetenv(ConsoleSecureSTSSeconds)
		})
	}
}

func Test_getLogSearchAPIToken(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "env set",
			args: args{
				env: "value",
			},
			want: "value",
		},
		{
			name: "env not set",
			args: args{
				env: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			os.Setenv(ConsoleLogQueryAuthToken, tt.args.env)
			assert.Equalf(t, tt.want, getLogSearchAPIToken(), "getLogSearchAPIToken()")
			os.Setenv(ConsoleLogQueryAuthToken, tt.args.env)
		})
	}
}

func Test_getPrometheusURL(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "env set",
			args: args{
				env: "value",
			},
			want: "value",
		},
		{
			name: "env not set",
			args: args{
				env: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			os.Setenv(PrometheusURL, tt.args.env)
			assert.Equalf(t, tt.want, getPrometheusURL(), "getPrometheusURL()")
			os.Setenv(PrometheusURL, tt.args.env)
		})
	}
}

func Test_getPrometheusJobID(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "env set",
			args: args{
				env: "value",
			},
			want: "value",
		},
		{
			name: "env not set",
			args: args{
				env: "",
			},
			want: "minio-job",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			os.Setenv(PrometheusJobID, tt.args.env)
			assert.Equalf(t, tt.want, getPrometheusJobID(), "getPrometheusJobID()")
			os.Setenv(PrometheusJobID, tt.args.env)
		})
	}
}

func Test_getMaxConcurrentUploadsLimit(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "valid",
			args: args{
				env: "1",
			},
			want: 1,
		},
		{
			name: "invalid",
			args: args{
				env: "duck",
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			os.Setenv(ConsoleMaxConcurrentUploads, tt.args.env)
			assert.Equalf(t, tt.want, getMaxConcurrentUploadsLimit(), "getMaxConcurrentUploadsLimit()")
			os.Unsetenv(ConsoleMaxConcurrentUploads)
		})
	}
}

func Test_getMaxConcurrentDownloadsLimit(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "valid",
			args: args{
				env: "1",
			},
			want: 1,
		},
		{
			name: "invalid",
			args: args{
				env: "duck",
			},
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			os.Setenv(ConsoleMaxConcurrentDownloads, tt.args.env)
			assert.Equalf(t, tt.want, getMaxConcurrentDownloadsLimit(), "getMaxConcurrentDownloadsLimit()")
			os.Unsetenv(ConsoleMaxConcurrentDownloads)
		})
	}
}

func Test_getConsoleDevMode(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "value set",
			args: args{
				env: "on",
			},
			want: true,
		},
		{
			name: "value not set",
			args: args{
				env: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			os.Setenv(ConsoleDevMode, tt.args.env)
			assert.Equalf(t, tt.want, getConsoleDevMode(), "getConsoleDevMode()")
			os.Unsetenv(ConsoleDevMode)
		})
	}
}
