// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

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

// jammedConfig represents the configuration for jammed comments.
type jammedConfig struct {
	Enabled *bool `hclext:"enabled,optional" hcl:"enabled,optional"`
	Tails   *bool `hclext:"tails,optional" hcl:"tails,optional"`
}

// checkJammed checks if comments are jammed (no space after delimiter).
func checkJammed(r *CommentsRule, text string, runner tflint.Runner, token hclsyntax.Token, _ *hclsyntax.Token) {
	enabled := true
	if r.Config.Jammed != nil && r.Config.Jammed.Enabled != nil {
		enabled = *r.Config.Jammed.Enabled
	}

	tails := false
	if r.Config.Jammed != nil && r.Config.Jammed.Tails != nil {
		tails = *r.Config.Jammed.Tails
	}

	if !tails {
		// https://regex101.com/r/GCMXXL/1
		jammedCommentParser = regexp.MustCompile(`^\s*(//|#|/\*)\S`)
	}

	if enabled {
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
