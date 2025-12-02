// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// defaultCommentsConfig is the default configuration for the CommentsRule.
var defaultCommentsConfig = commentsRuleConfig{
	Block:  true,
	EOL:    true,
	Jammed: true,
	Length: &lengthConfig{
		AllowURL: func() *bool { b := true; return &b }(),
		Column:   80,
	},
	Level: "warning",
}

// commentsRuleConfig represents the configuration for the CommentsRule.
type commentsRuleConfig struct {
	// Enable block /* */ comment check.
	Block bool `hclext:"block,optional" hcl:"block,optional"`
	// Enable EOL comment check.
	EOL    bool          `hclext:"eol,optional" hcl:"eol,optional"`
	Jammed bool          `hclext:"jammed,optional" hcl:"jammed,optional"`
	Length *lengthConfig `hclext:"length,block" hcl:"length,block"`
	// Issue level.
	Level string `hclext:"level,optional" hcl:"level,optional"`
	// Minimum ration threshold of comments to code PER SOURCE FILE.
	Threshold *float64 `hclext:"threshold,optional" hcl:"threshold,optional"`
}

type lengthConfig struct {
	// Allow comments with URL to bust the Column limit.
	AllowURL *bool `hclext:"allow_url,optional" hcl:"allow_url,optional"`
	// Maximum allowed column for comments. <=0 >= 99 effectively disable check.
	Column int `hclext:"column,optional" hcl:"column,optional"`
}

// CommentsRule checks for comment style.
type CommentsRule struct {
	tflint.DefaultRule
	Config commentsRuleConfig
}

// Check checks whether the rule conditions are met.
func (r *CommentsRule) Check(runner tflint.Runner) error {
	if err := runner.DecodeRuleConfig(r.Name(), &r.Config); err != nil {
		return err
	}

	// The threshold check is done first as it's parsing and checking an entire
	// source file as opposed to parsing an entire set of source files and then
	// checking each comment.
	if err := checkThreshold(r, runner); err != nil {
		return err
	}

	return checkCommentsWithContext(runner, r,
		checkBlock,
		checkEOL,
		checkJammed,
		checkLength)
}

// checkCommentsWithContext iterates over all files in the root module, parses
// them, and applies the check function to each comment token, providing the
// previous token for context.
func checkCommentsWithContext(
	runner tflint.Runner,
	rule *CommentsRule,
	checkFuncs ...func(*CommentsRule, string, tflint.Runner, hclsyntax.Token, *hclsyntax.Token),
) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for filename, file := range files {
		tokens, diags := hclsyntax.LexConfig(file.Bytes, filename, hcl.InitialPos)
		if diags.HasErrors() {
			return diags
		}

		for i, token := range tokens {
			if token.Type != hclsyntax.TokenComment {
				continue
			}

			var prevToken *hclsyntax.Token
			for j := i - 1; j >= 0; j-- {
				if tokens[j].Type != hclsyntax.TokenComment {
					prevToken = &tokens[j]
					break
				}
			}

			for _, checkFunc := range checkFuncs {
				checkFunc(rule, string(token.Bytes), runner, token, prevToken)
			}
		}
	}

	return nil
}

// NewCommentsRule returns a new rule.
func NewCommentsRule() *CommentsRule {
	rule := &CommentsRule{}
	rule.Config = defaultCommentsConfig

	return rule
}

// Enabled returns whether the rule is enabled by default.
func (r *CommentsRule) Enabled() bool {
	return true
}

// Link returns the rule link.
func (r *CommentsRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_comments.md"
}

// Name returns the rule name.
func (r *CommentsRule) Name() string {
	return "eos_comments"
}

// Severity returns the rule severity.
func (r *CommentsRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
