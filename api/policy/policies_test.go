// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package policy

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/openstor/madmin-go/v4"
	minioIAMPolicy "github.com/openstor/pkg/v3/policy"
)

func TestReplacePolicyVariables(t *testing.T) {
	type args struct {
		claims      map[string]interface{}
		accountInfo *madmin.AccountInfo
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Bad Policy",
			args: args{
				claims: nil,
				accountInfo: &madmin.AccountInfo{
					AccountName: "test",
					Server:      madmin.BackendInfo{},
					Policy:      []byte(""),
					Buckets:     nil,
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Replace basic AWS",
			args: args{
				claims: nil,
				accountInfo: &madmin.AccountInfo{
					AccountName: "test",
					Server:      madmin.BackendInfo{},
					Policy: []byte(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::${aws:username}",
        "arn:aws:s3:::${aws:userid}"
      ]
    }
  ]
}`),
					Buckets: nil,
				},
			},
			want: `{
          "Version": "2012-10-17",
          "Statement": [
            {
              "Effect": "Allow",
              "Action": [
                "s3:ListBucket"
              ],
              "Resource": [
                "arn:aws:s3:::test",
                "arn:aws:s3:::test"
              ]
            }
          ]
        }`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			got := ReplacePolicyVariables(tt.args.claims, tt.args.accountInfo)
			policy, err := minioIAMPolicy.ParseConfig(bytes.NewReader(got))
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplacePolicyVariables() error = %v, wantErr %v", err, tt.wantErr)
			}
			wantPolicy, err := minioIAMPolicy.ParseConfig(bytes.NewReader([]byte(tt.want)))
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplacePolicyVariables() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(policy, wantPolicy) {
				t.Errorf("ReplacePolicyVariables() = %s, want %v", got, tt.want)
			}
		})
	}
}
