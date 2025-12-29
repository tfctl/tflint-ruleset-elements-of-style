// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rulehelper

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// BoolPtr returns a pointer to the given bool value.
func BoolPtr(b bool) *bool {
	return &b
}

// ToSeverity converts a string level to a tflint.Severity.
func ToSeverity(level string) tflint.Severity {
	switch strings.ToLower(level) {
	case "notice":
		return tflint.NOTICE
	case "warning":
		return tflint.WARNING
	}

	return tflint.ERROR
}
