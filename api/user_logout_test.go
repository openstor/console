// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import "testing"

// mock function of Get()
func (ac consoleCredentialsMock) Expire() {
	// Do nothing
	// Implementing this method for the consoleCredentials interface
}

func TestLogout(_ *testing.T) {
	// There's nothing to test right now
}
