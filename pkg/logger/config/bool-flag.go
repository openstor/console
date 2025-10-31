// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// BoolFlag - wrapper bool type.
type BoolFlag bool

// String - returns string of BoolFlag.
func (bf BoolFlag) String() string {
	if bf {
		return "on"
	}

	return "off"
}

// MarshalJSON - converts BoolFlag into JSON data.
func (bf BoolFlag) MarshalJSON() ([]byte, error) {
	return json.Marshal(bf.String())
}

// UnmarshalJSON - parses given data into BoolFlag.
func (bf *BoolFlag) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err == nil {
		b := BoolFlag(true)
		if s == "" {
			// Empty string is treated as valid.
			*bf = b
		} else if b, err = ParseBoolFlag(s); err == nil {
			*bf = b
		}
	}

	return err
}

// ParseBool returns the boolean value represented by the string.
func ParseBool(str string) (bool, error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "on", "ON", "On":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "off", "OFF", "Off":
		return false, nil
	}
	if strings.EqualFold(str, "enabled") {
		return true, nil
	}
	if strings.EqualFold(str, "disabled") {
		return false, nil
	}
	return false, fmt.Errorf("ParseBool: parsing '%s': %w", str, strconv.ErrSyntax)
}

// ParseBoolFlag - parses string into BoolFlag.
func ParseBoolFlag(s string) (bf BoolFlag, err error) {
	b, err := ParseBool(s)
	return BoolFlag(b), err
}
