// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"regexp"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// defaultHeredocConfig is the default configuration for the HeredocRule.
var defaultHeredocConfig = heredocRuleConfig{
	EOF:   true,
	Level: "warning",
}

// heredocPattern is a regex to match heredoc delimiters.
var heredocPattern = regexp.MustCompile(`<<(-?)([a-zA-Z0-9]+)\s*$`)

// heredocRuleConfig represents the configuration for the HeredocRule.
type heredocRuleConfig struct {
	Enabled *bool  `hclext:"enabled,optional" hcl:"enabled,optional"`
	EOF     bool   `hclext:"EOF,optional" hcl:"EOF,optional"`
	Level   string `hclext:"level,optional" hcl:"level,optional"`
}

// HeredocRule checks for standard heredoc usage.
type HeredocRule struct {
	tflint.DefaultRule
	Config heredocRuleConfig
}

// Check checks whether the rule conditions are met.
func (r *HeredocRule) Check(runner tflint.Runner) error {
	if err := runner.DecodeRuleConfig(r.Name(), &r.Config); err != nil {
		return err
	}

	return rulehelper.WalkTokens(runner, r, checkHeredocToken)
}

// checkHeredocToken checks for heredoc style violations in a token.
func checkHeredocToken(runner tflint.Runner, r *HeredocRule, token hclsyntax.Token) {
	if token.Type == hclsyntax.TokenOHeredoc {
		text := string(token.Bytes)
		if matches := heredocPattern.FindStringSubmatch(text); matches != nil {
			indentMarker := matches[1]
			heredocLabel := matches[2]

			if indentMarker == "" {
				message := "Avoid standard heredoc (<<). Use indented (<<-) instead."
				if err := runner.EmitIssue(r, message, token.Range); err != nil {
					logger.Error(err.Error())
				}
			}

			if r.Config.EOF {
				// Also check for EOF usage.
				if heredocLabel == "EOF" {
					eofMessage := "Avoid using 'EOF' as the heredoc delimiter."
					if err := runner.EmitIssue(r, eofMessage, token.Range); err != nil {
						logger.Error(err.Error())
					}
				}
			}
		}
	}
}

// NewHeredocRule returns a new rule.
func NewHeredocRule() *HeredocRule {
	rule := &HeredocRule{}
	rule.Config = defaultHeredocConfig
	return rule
}

// Enabled returns whether the rule is enabled by default.
func (r *HeredocRule) Enabled() bool {
	return true
}

// Link returns the rule link.
func (r *HeredocRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_heredoc.md"
}

// Name returns the rule name.
func (r *HeredocRule) Name() string {
	return "eos_heredoc"
}

// Severity returns the rule severity.
func (r *HeredocRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
