// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// jammedCommentParser is a regex to detect jammed comments.
// https://regex101.com/r/5HRrLc/1
var jammedCommentParser = regexp.MustCompile(`^\s*(///*|##*|/\*\**)([^\s/#])`)

// checkJammed checks if comments are jammed (no space after delimiter).
func checkJammed(r *CommentsRule, text string, runner tflint.Runner, token hclsyntax.Token, _ *hclsyntax.Token) {
	if r.Config.Jammed {
		if jammedCommentParser.MatchString(text) {
			trimmed := strings.TrimSpace(text)
			rns := []rune(trimmed)
			snippet := trimmed
			if len(rns) > 5 {
				snippet = string(rns[:5])
			}
			message := fmt.Sprintf("Avoid jammed comment ('%s ...').", snippet)
			if err := runner.EmitIssue(r, message, token.Range); err != nil {
				logger.Error(err.Error())
			}
			logger.Debug(message)
		}
	}
}
