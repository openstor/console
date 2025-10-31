// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package types

// TargetType indicates type of the target e.g. console, http, kafka
type TargetType uint8

// Constants for target types
const (
	_ TargetType = iota
	TargetConsole
	TargetHTTP
)
