// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package hungarian

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

// hungarianConfig represents the configuration for the HungarianRule.
type hungarianConfig struct {
	Enabled *bool    `hclext:"enabled,optional" hcl:"enabled,optional"`
	Level   string   `hclext:"level,optional" hcl:"level,optional"`
	Tags    []string `hclext:"tags,optional" hcl:"tags,optional"`
}

// defaultHungarianConfig is the default configuration for the HungarianRule.
var defaultHungarianConfig = hungarianConfig{
	Enabled: rulehelper.BoolPtr(true),
	Tags:    defaultHungarianTags,
	Level:   "warning",
}

// HungarianRule checks whether a block's type is echoed in its name.
type HungarianRule struct {
	tflint.DefaultRule
	Config hungarianConfig
	// RuleName is the rule block name to load from the config file. If empty,
	// defaults to "eos_hungarian".
	RuleName string
	// ConfigFile is the path to the config file. If empty, LoadRuleConfig will
	// search CWD then $HOME for .tflint.hcl.
	ConfigFile string
}

// Check checks whether the rule conditions are met.
func (r *HungarianRule) Check(runner tflint.Runner) error {
	// Load config using the rule name and optional config file path.
	if err := rulehelper.LoadRuleConfig(r.Name(), &r.Config, r.ConfigFile); err != nil {
		return err
	}

	// Bail out early if the rule is not enabled. This will occur if the EOS
	// plugin is enabled, but this specific rule is not.
	if !r.Enabled() {
		return nil
	}

	return rulehelper.WalkBlocks(runner, rulehelper.AllLintableBlocks, r, checkForHungarian)
}

// checkForHungarian checks if the name uses Hungarian notation.
func checkForHungarian(runner tflint.Runner, r *HungarianRule, defRange hcl.Range, typ string, name string, _ string) {

	// Combine the built-in defaults with extras defined in config.
	tags := make([]string, 0, len(defaultHungarianTags)+len(r.Config.Tags))
	tags = append(tags, defaultHungarianTags...)
	tags = append(tags, r.Config.Tags...)

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
	return r.Config.Enabled == nil || *r.Config.Enabled
}

// Link returns the rule reference link.
func (r *HungarianRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_hungarian.md"
}

// Name returns the rule name.
func (r *HungarianRule) Name() string {
	if r.RuleName != "" {
		return r.RuleName
	}

	return "eos_hungarian"
}

// Severity returns the rule severity.
func (r *HungarianRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
