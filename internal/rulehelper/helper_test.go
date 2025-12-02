// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rulehelper

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestRuleHelper(t *testing.T) {
	t.Run("ToSeverity", testToSeverity)
}

func testToSeverity(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected tflint.Severity
	}{
		{
			Name:     "notice",
			Input:    "notice",
			Expected: tflint.NOTICE,
		},
		{
			Name:     "NOTICE",
			Input:    "NOTICE",
			Expected: tflint.NOTICE,
		},
		{
			Name:     "warning",
			Input:    "warning",
			Expected: tflint.WARNING,
		},
		{
			Name:     "WARNING",
			Input:    "warning",
			Expected: tflint.WARNING,
		},
		{
			Name:     "error",
			Input:    "error",
			Expected: tflint.ERROR,
		},
		{
			Name:     "unknown",
			Input:    "unknown",
			Expected: tflint.ERROR,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			result := ToSeverity(tc.Input)
			if result != tc.Expected {
				t.Errorf("ToSeverity(%q) = %v, expected %v", tc.Input, result, tc.Expected)
			}
		})
	}
}
