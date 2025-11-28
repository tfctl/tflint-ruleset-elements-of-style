// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// checkLength checks if comments exceed the column limit.
func checkLength(r *CommentsRule, text string, runner tflint.Runner, token hclsyntax.Token, _ *hclsyntax.Token) {
	if r.Config.Column > 0 {
		trimmedText := strings.TrimRight(text, "\r\n")
		end := token.Range.Start.Column + len(trimmedText) - 1

		if r.Config.URLBypass {
			// Simple URL detection.
			if strings.Contains(trimmedText, "http://") || strings.Contains(trimmedText, "https://") {
				return
			}
		}

		if end > r.Config.Column {
			message := fmt.Sprintf("Wrap comment at column %d (currently %d).", r.Config.Column, end)
			if err := runner.EmitIssue(r, message, token.Range); err != nil {
				logger.Error(err.Error())
			}
			logger.Debug(message)

		}
	}
}
