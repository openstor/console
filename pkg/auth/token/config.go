// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"time"

	"github.com/openstor/console/pkg/auth/utils"
	"github.com/openstor/pkg/v3/env"
)

// GetConsoleSTSDuration returns the default session duration for the STS requested tokens (defaults to 12h)
func GetConsoleSTSDuration() time.Duration {
	duration, err := time.ParseDuration(env.Get(ConsoleSTSDuration, "12h"))
	if err != nil || duration <= 0 {
		duration = 12 * time.Hour
	}
	return duration
}

var defaultPBKDFPassphrase = utils.RandomCharString(64)

// GetPBKDFPassphrase returns passphrase for the pbkdf2 function used to encrypt JWT payload
func GetPBKDFPassphrase() string {
	return env.Get(ConsolePBKDFPassphrase, defaultPBKDFPassphrase)
}

var defaultPBKDFSalt = utils.RandomCharString(64)

// GetPBKDFSalt returns salt for the pbkdf2 function used to encrypt JWT payload
func GetPBKDFSalt() string {
	return env.Get(ConsolePBKDFSalt, defaultPBKDFSalt)
}
