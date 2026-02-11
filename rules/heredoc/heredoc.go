// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package heredoc

import (
	"regexp"

	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const AvoidEOFHeredocMessage = "Avoid using 'EOF' as the heredoc delimiter."
const AvoidStandardHeredocMessage = "Avoid standard heredoc (<<). Use indented (<<-) instead."

// heredocPattern is a regex to match heredoc delimiters.
var heredocPattern = regexp.MustCompile(`<<(-?)([a-zA-Z0-9]+)\s*$`)

// heredocConfig represents the configuration for the HeredocRule.
type heredocConfig struct {
	Enabled *bool  `hclext:"enabled,optional" hcl:"enabled,optional"`
	EOF     bool   `hclext:"EOF,optional" hcl:"EOF,optional"`
	Level   string `hclext:"level,optional" hcl:"level,optional"`
}

// defaultHeredocConfig is the default configuration for the HeredocRule.
var defaultHeredocConfig = heredocConfig{
	EOF:   true,
	Level: "warning",
}

// Rule checks for standard heredoc usage.
type Rule struct {
	tflint.DefaultRule
	Config heredocConfig
	// RuleName is the rule block name to load from the config file. If empty,
	// defaults to "eos_heredoc".
	RuleName string
	// ConfigFile is the path to the config file. If empty, LoadRuleConfig will
	// search CWD then $HOME for .tflint.hcl.
	ConfigFile string
}

// Check checks whether the rule conditions are met.
func (r *Rule) Check(runner tflint.Runner) error {
	// Load config using the rule name and optional config file path.
	if err := rulehelper.LoadRuleConfig(r.Name(), &r.Config, r.ConfigFile); err != nil {
		return err
	}

	// Bail out early if the rule is not enabled. This will occur if the EOS
	// plugin is enabled, but this specific rule is not.
	if !r.Enabled() {
		return nil
	}
	return rulehelper.WalkTokens(runner, r, checkHeredocToken)
}

// checkHeredocToken checks for heredoc style violations in a token.
func checkHeredocToken(runner tflint.Runner, r *Rule, token hclsyntax.Token) {
	if token.Type == hclsyntax.TokenOHeredoc {
		text := string(token.Bytes)
		if matches := heredocPattern.FindStringSubmatch(text); matches != nil {
			indentMarker := matches[1]
			heredocLabel := matches[2]

			if indentMarker == "" {
				message := AvoidStandardHeredocMessage
				if err := runner.EmitIssue(r, message, token.Range); err != nil {
					logger.Error(err.Error())
				}
			}

			if r.Config.EOF {
				// Also check for EOF usage.
				if heredocLabel == "EOF" {
					eofMessage := AvoidEOFHeredocMessage
					if err := runner.EmitIssue(r, eofMessage, token.Range); err != nil {
						logger.Error(err.Error())
					}
				}
			}
		}
	}
}

// NewHeredocRule returns a new rule.
func NewHeredocRule() *Rule {
	rule := &Rule{}
	rule.Config = defaultHeredocConfig
	return rule
}

// Enabled returns whether the rule is enabled by default.
func (r *Rule) Enabled() bool {
	return r.Config.Enabled == nil || *r.Config.Enabled
}

// Link returns the rule link.
func (r *Rule) Link() string {
	return "https://github.com/tfctl/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_heredoc.md"
}

// Name returns the rule name.
func (r *Rule) Name() string {
	if r.RuleName != "" {
		return r.RuleName
	}
	return "eos_heredoc"
}

// Severity returns the rule severity.
func (r *Rule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
