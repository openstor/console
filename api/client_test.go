// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import "testing"

func Test_computeObjectURLWithoutEncode(t *testing.T) {
	type args struct {
		bucketName string
		prefix     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "http://localhost:9000/bucket-1/小飼弾小飼弾小飼弾.jp",
			args: args{
				bucketName: "bucket-1",
				prefix:     "小飼弾小飼弾小飼弾.jpg",
			},
			want:    "http://localhost:9000/bucket-1/小飼弾小飼弾小飼弾.jpg",
			wantErr: false,
		},
		{
			name: "http://localhost:9000/bucket-1/a a - a a & a a - a a a.jpg",
			args: args{
				bucketName: "bucket-1",
				prefix:     "a a - a a & a a - a a a.jpg",
			},
			want:    "http://localhost:9000/bucket-1/a a - a a & a a - a a a.jpg",
			wantErr: false,
		},
		{
			name: "http://localhost:9000/bucket-1/02%20-%20FLY%20ME%20TO%20THE%20MOON%20.jpg",
			args: args{
				bucketName: "bucket-1",
				prefix:     "02%20-%20FLY%20ME%20TO%20THE%20MOON%20.jpg",
			},
			want:    "http://localhost:9000/bucket-1/02%20-%20FLY%20ME%20TO%20THE%20MOON%20.jpg",
			wantErr: false,
		},
		{
			name: "http://localhost:9000/bucket-1/!@#$%^&*()_+.jpg",
			args: args{
				bucketName: "bucket-1",
				prefix:     "!@#$%^&*()_+.jpg",
			},
			want:    "http://localhost:9000/bucket-1/!@#$%^&*()_+.jpg",
			wantErr: false,
		},
		{
			name: "http://localhost:9000/bucket-1/test/test2/小飼弾小飼弾小飼弾.jpg",
			args: args{
				bucketName: "bucket-1",
				prefix:     "test/test2/小飼弾小飼弾小飼弾.jpg",
			},
			want:    "http://localhost:9000/bucket-1/test/test2/小飼弾小飼弾小飼弾.jpg",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			got, err := computeObjectURLWithoutEncode(tt.args.bucketName, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeObjectURLWithoutEncode() errors = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if got != tt.want {
					t.Errorf("computeObjectURLWithoutEncode() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
