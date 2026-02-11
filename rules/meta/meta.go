// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// OrderConfig defines which arguments must appear first and last in a block.
type OrderConfig struct {
	First []string `hclext:"first,optional" hcl:"first,optional"`
	Last  []string `hclext:"last,optional" hcl:"last,optional"`
}

// metaConfig represents the configuration for the MetaRule.
type metaConfig struct {
	Enabled       *bool         `hclext:"enabled,optional" hcl:"enabled,optional"`
	Level         string        `hclext:"level,optional" hcl:"level,optional"`
	Order         []OrderConfig `hclext:"order,block" hcl:"order,block"`
	SourceVersion *bool         `hcl:"source_version,optional"`
}

// defaultMetaConfig is the default configuration for the MetaRule.
var defaultMetaConfig = metaConfig{
	Enabled: rulehelper.BoolPtr(true),
	Level:   "warning",
	Order: []OrderConfig{{
		First: []string{"for_each", "count"},
		Last:  []string{"depends_on", "provider", "lifecycle"},
	}},
	SourceVersion: rulehelper.BoolPtr(true),
}

// Rule checks for meta-argument style violations.
type Rule struct {
	tflint.DefaultRule
	Config metaConfig
	// RuleName is the rule block name to load from the config file. If empty,
	// defaults to "eos_meta".
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

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		if body, ok := file.Body.(*hclsyntax.Body); ok {
			for _, block := range body.Blocks {
				checkOrder(runner, r, block)
				if attr, exists := block.Body.Attributes["count"]; exists {
					checkCountGuard(runner, r, attr)
				}
				if block.Type == "module" && (r.Config.SourceVersion == nil || *r.Config.SourceVersion) {
					checkModuleSourceVersion(runner, r, block)
				}
			}
		}
	}

	return nil
}

func (r *Rule) emitIssue(runner tflint.Runner, message string, rng hcl.Range) {
	if err := runner.EmitIssue(r, message, rng); err != nil {
		logger.Error(err.Error())
	}
}

// NewMetaRule returns a new rule.
func NewMetaRule() *Rule {
	rule := &Rule{}
	rule.Config = defaultMetaConfig
	return rule
}

// Enabled returns whether the rule is enabled by default.
func (r *Rule) Enabled() bool {
	return r.Config.Enabled == nil || *r.Config.Enabled
}

// Link returns the rule reference link.
func (r *Rule) Link() string {
	return "https://github.com/tfctl/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_meta.md"
}

// Name returns the rule name.
func (r *Rule) Name() string {
	if r.RuleName != "" {
		return r.RuleName
	}

	return "eos_meta"
}

// Severity returns the rule severity.
func (r *Rule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
