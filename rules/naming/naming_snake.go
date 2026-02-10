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

// checkSnake checks if the name consists only of lowercase alphanumeric and underscores.
func checkSnake(runner tflint.Runner, rule *Rule, defRange hcl.Range, _ string, name string, _ string) {
	valid := true
	for _, ch := range name {
		if !(unicode.IsLower(ch) || unicode.IsDigit(ch) || ch == '_') {
			valid = false
			break
		}
	}

	if !valid {
		message := fmt.Sprintf("Names should be snake_case (%s).", name)
		if err := runner.EmitIssue(rule, message, defRange); err != nil {
			logger.Error(err.Error())
		}
		logger.Debug(message)
	}
}
