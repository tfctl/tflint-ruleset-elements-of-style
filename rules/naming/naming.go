// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package naming

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// defaultLimit is the default maximum length for names.
const defaultLimit = 16

// defaultLength is the default length value.
var defaultLength = defaultLimit

// defaultConfig is the default configuration for the NamingRule.
var defaultConfig = namingRuleConfig{
	Enabled: func() *bool { b := true; return &b }(),
	Level:   "warning",
	Length:  &defaultLength,
}

// typeEchoConfig represents the configuration for type echo checks.
type typeEchoConfig struct {
	Enabled  *bool               `hclext:"enabled,optional" hcl:"enabled,optional"`
	Level    string              `hclext:"level,optional" hcl:"level,optional"`
	Synonyms map[string][]string `hclext:"synonyms,optional" hcl:"synonyms,optional"`
}

// namingRuleConfig represents the configuration for the NamingRule.
type namingRuleConfig struct {
	Enabled  *bool           `hclext:"enabled,optional" hcl:"enabled,optional"`
	Level    string          `hclext:"level,optional" hcl:"level,optional"`
	Length   *int            `hclext:"length,optional" hcl:"length,optional"`
	Shout    *bool           `hclext:"shout,optional" hcl:"shout,optional"`
	Snake    *bool           `hclext:"snake,optional" hcl:"snake,optional"`
	TypeEcho *typeEchoConfig `hclext:"type_echo,optional" hcl:"type_echo,block"`
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

	var checks []func(tflint.Runner, *NamingRule, hcl.Range, string, string, string)
	length := defaultLimit
	if rule.Config.Length != nil {
		length = *rule.Config.Length
	}
	if length > 0 {
		checks = append(checks, checkNameLength)
	}

	if rule.Config.Shout == nil || *rule.Config.Shout {
		checks = append(checks, checkShout)
	}

	if rule.Config.Snake == nil || *rule.Config.Snake {
		checks = append(checks, checkSnake)
	}

	te := rule.Config.TypeEcho
	if te == nil || te.Enabled == nil || *te.Enabled {
		checks = append(checks, checkTypeEcho)
	}

	return rulehelper.WalkBlocks(runner, rulehelper.AllLintableBlocks, rule, checks...)
}

// NewNamingRule returns a new rule.
func NewNamingRule() *NamingRule {
	rule := &NamingRule{}
	rule.Config = defaultConfig
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
