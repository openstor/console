// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package certs

const (
	// Default minio configuration directory where below configuration files/directories are stored.
	DefaultConsoleConfigDir = ".console"

	// Directory contains below files/directories for HTTPS configuration.
	CertsDir = "certs"

	// Directory contains all CA certificates other than system defaults for HTTPS.
	CertsCADir = "CAs"

	// Public certificate file for HTTPS.
	PublicCertFile = "public.crt"

	// Private key file for HTTPS.
	PrivateKeyFile = "private.key"
)
