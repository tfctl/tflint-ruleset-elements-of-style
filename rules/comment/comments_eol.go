// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const avoidEOLCommentsMessage = "Avoid EOL comments."

// checkEOL checks if EOL comments are used.
func checkEOL(r *Rule, _ string, runner tflint.Runner, token hclsyntax.Token, prevToken *hclsyntax.Token) {
	if !r.Config.EOL {
		return
	}

	if prevToken != nil {
		if prevToken.Type == hclsyntax.TokenNewline {
			return
		}

		if prevToken.Range.End.Line == token.Range.Start.Line {
			message := avoidEOLCommentsMessage
			if err := runner.EmitIssue(r, message, token.Range); err != nil {
				logger.Error(err.Error())
			}
			logger.Debug(message)
		}
	}
}
