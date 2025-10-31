// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"errors"
)

// EnsureCertAndKey checks if both client certificate and key paths are provided
func EnsureCertAndKey(clientCert, clientKey string) error {
	if (clientCert != "" && clientKey == "") ||
		(clientCert == "" && clientKey != "") {
		return errors.New("cert and key must be specified as a pair")
	}
	return nil
}
