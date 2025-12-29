// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// metaConfig represents the configuration for the MetaRule.
type metaConfig struct {
	Enabled       *bool  `hclext:"enabled,optional" hcl:"enabled,optional"`
	Level         string `hclext:"level,optional" hcl:"level,optional"`
	SourceVersion *bool  `hcl:"source_version,optional"`
}

// defaultMetaConfig is the default configuration for the MetaRule.
var defaultMetaConfig = metaConfig{
	Enabled:       rulehelper.BoolPtr(true),
	Level:         "warning",
	SourceVersion: rulehelper.BoolPtr(true),
}

// MetaRule checks for meta-argument style violations.
type MetaRule struct {
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
func (r *MetaRule) Check(runner tflint.Runner) error {
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

func (r *MetaRule) emitIssue(runner tflint.Runner, message string, rng hcl.Range) {
	if err := runner.EmitIssue(r, message, rng); err != nil {
		logger.Error(err.Error())
	}
}

// NewMetaRule returns a new rule.
func NewMetaRule() *MetaRule {
	rule := &MetaRule{}
	rule.Config = defaultMetaConfig
	return rule
}

// Enabled returns whether the rule is enabled by default.
func (r *MetaRule) Enabled() bool {
	return r.Config.Enabled == nil || *r.Config.Enabled
}

// Link returns the rule reference link.
func (r *MetaRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_meta.md"
}

// Name returns the rule name.
func (r *MetaRule) Name() string {
	if r.RuleName != "" {
		return r.RuleName
	}

	return "eos_meta"
}

// Severity returns the rule severity.
func (r *MetaRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
