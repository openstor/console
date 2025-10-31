// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package portalui

import "embed"

//go:embed build/*
var fs embed.FS

func GetStaticAssets() embed.FS {
	return fs
}
