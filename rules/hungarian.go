// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"fmt"
	"strings"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// defaultHungarianTags is the default list of Hungarian notation tags.
var defaultHungarianTags = []string{
	"str",
	"int", "num",
	"bool",
	"list", "lst", "set", "map", "arr", "array",
}

// defaultHungarianConfig is the default configuration for the HungarianRule.
var defaultHungarianConfig = hungarianRuleConfig{
	Tags:  defaultHungarianTags,
	Level: "warning",
}

// hungarianRuleConfig represents the configuration for the HungarianRule.
type hungarianRuleConfig struct {
	Enabled *bool    `hclext:"enabled,optional" hcl:"enabled,optional"`
	Tags    []string `hclext:"tags,optional" hcl:"tags,optional"`
	Level   string   `hclext:"level,optional" hcl:"level,optional"`
}

// HungarianRule checks whether a block's type is echoed in its name.
type HungarianRule struct {
	tflint.DefaultRule
	Config hungarianRuleConfig
}

// Check checks whether the rule conditions are met.
func (r *HungarianRule) Check(runner tflint.Runner) error {
	if err := runner.DecodeRuleConfig(r.Name(), &r.Config); err != nil {
		return err
	}

	return rulehelper.WalkBlocks(runner, rulehelper.AllLintableBlocks, r, checkForHungarian)
}

// checkForHungarian checks if the name uses Hungarian notation.
func checkForHungarian(runner tflint.Runner, r *HungarianRule, defRange hcl.Range, typ string, name string, _ string) {
	tags := r.Config.Tags

	for _, t := range tags {
		if strings.HasPrefix(name, t) || strings.HasSuffix(name, t) || strings.Contains(name, "_"+t) {
			if err := runner.EmitIssue(r, fmt.Sprintf("Avoid Hungarian notation '%s' in '%s'.", t, name),
				defRange); err != nil {
				logger.Error(err.Error())
			}
			return
		}
	}
}

// NewHungarianRule returns a new rule.
func NewHungarianRule() *HungarianRule {
	rule := &HungarianRule{}
	rule.Config = defaultHungarianConfig
	return rule
}

// Enabled returns whether the rule is enabled by default.
func (r *HungarianRule) Enabled() bool {
	return true
}

// Link returns the rule reference link.
func (r *HungarianRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_hungarian.md"
}

// Name returns the rule name.
func (r *HungarianRule) Name() string {
	return "eos_hungarian"
}

// Severity returns the rule severity.
func (r *HungarianRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
