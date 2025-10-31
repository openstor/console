// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package ldap

import (
	"strings"

	"github.com/openstor/pkg/v3/env"
)

func GetLDAPEnabled() bool {
	return strings.ToLower(env.Get(ConsoleLDAPEnabled, "off")) == "on"
}
