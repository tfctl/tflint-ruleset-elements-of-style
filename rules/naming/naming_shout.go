// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package naming

import (
	"fmt"
	"unicode"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// checkShout checks if the name is all uppercase.
func checkShout(runner tflint.Runner, r *NamingRule, defRange hcl.Range, _ string, name string, _ string) {
	hasAlpha := false
	allUpper := true

	for _, ch := range name {
		if unicode.IsLetter(ch) {
			hasAlpha = true
			if !unicode.IsUpper(ch) {
				allUpper = false
			}
		}
	}

	if hasAlpha && allUpper {
		message := fmt.Sprintf("Avoid SHOUTED names (%s)", name)
		if err := runner.EmitIssue(r, message, defRange); err != nil {
			logger.Error(err.Error())
		}
		logger.Debug(message)
	}
}
