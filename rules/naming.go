// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"fmt"
	"unicode"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// defaultNameLimit is the default maximum length for names.
const defaultNameLimit = 16

// defaultNamingConfig is the default configuration for the NamingRule.
var defaultNamingConfig = namingRuleConfig{
	Level: "warning",
}

// namingLengthConfig represents the configuration for name length limits.
type namingLengthConfig struct {
	Enabled *bool `hclext:"enabled,optional" hcl:"enabled,optional"`
	Limit   int   `hclext:"limit,optional" hcl:"limit,optional"`
}

// namingShoutConfig represents the configuration for uppercase name checks.
type namingShoutConfig struct {
	Enabled *bool `hclext:"enabled,optional" hcl:"enabled,optional"`
}

// namingRuleConfig represents the configuration for the NamingRule.
type namingRuleConfig struct {
	Enabled *bool               `hclext:"enabled,optional" hcl:"enabled,optional"`
	Level   string              `hclext:"level,optional" hcl:"level,optional"`
	Length  *namingLengthConfig `hclext:"length,optional" hcl:"length,block"`
	Shout   *namingShoutConfig  `hclext:"shout,optional" hcl:"shout,block"`
}

// NamingRule checks whether a block's name is excessively long.
type NamingRule struct {
	tflint.DefaultRule
	Config namingRuleConfig
}

// Check checks whether the rule conditions are met.
func (rule *NamingRule) Check(runner tflint.Runner) error {
	if err := runner.DecodeRuleConfig(rule.Name(), &rule.Config); err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("rule.Config=%v", rule.Config))

	return rulehelper.WalkBlocks(runner, rulehelper.AllLintableBlocks, rule,
		checkNameLength, checkShout)
}

// checkNameLength checks if the name is too long.
func checkNameLength(runner tflint.Runner, r *NamingRule, defRange hcl.Range, _ string, name string, _ string) {
	limit := defaultNameLimit
	if r.Config.Length != nil && r.Config.Length.Limit > 0 {
		limit = r.Config.Length.Limit
	}

	if len(name) > limit {
		message := fmt.Sprintf("Avoid names longer than %d ('%s' is %d).", limit, name, len(name))
		if err := runner.EmitIssue(r, message, defRange); err != nil {
			logger.Error(err.Error())
		}
		logger.Debug(message)
	}
}

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

// NewNamingRule returns a new rule.
func NewNamingRule() *NamingRule {
	rule := &NamingRule{}
	rule.Config = defaultNamingConfig
	return rule
}

// Enabled returns whether the rule is enabled by default
func (rule *NamingRule) Enabled() bool {
	return true
}

// Link returns the rule reference link
func (rule *NamingRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_naming.md"
}

// Name returns the rule name.
func (rule *NamingRule) Name() string {
	return "eos_naming"
}

// Severity returns the rule severity
func (rule *NamingRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(rule.Config.Level)
}
