// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const avoidBlockCommentsMessage = "Avoid block comments."

// checkBlock checks if block comments are used.
func checkBlock(r *Rule, text string, runner tflint.Runner, token hclsyntax.Token, _ *hclsyntax.Token) {
	if r.Config.Block {
		if strings.HasPrefix(text, "/*") {
			message := avoidBlockCommentsMessage
			if err := runner.EmitIssue(r, message, token.Range); err != nil {
				logger.Error(err.Error())
			}
			logger.Debug(message)
		}
	}
}
