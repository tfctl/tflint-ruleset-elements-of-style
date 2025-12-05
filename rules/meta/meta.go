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

// defaultMetaConfig is the default configuration for the MetaRule.
var defaultMetaConfig = metaRuleConfig{
	Level:         "warning",
	SourceVersion: func() *bool { b := true; return &b }(),
}

// metaRuleConfig represents the configuration for the MetaRule.
type metaRuleConfig struct {
	Enabled       *bool  `hclext:"enabled,optional" hcl:"enabled,optional"`
	Level         string `hclext:"level,optional" hcl:"level,optional"`
	SourceVersion *bool  `hcl:"source_version,optional"`
}

// MetaRule checks for meta-argument style violations.
type MetaRule struct {
	tflint.DefaultRule
	Config metaRuleConfig
}

// NewMetaRule returns a new rule.
func NewMetaRule() *MetaRule {
	rule := &MetaRule{}
	rule.Config = defaultMetaConfig
	return rule
}

// Check checks whether the rule conditions are met.
func (r *MetaRule) Check(runner tflint.Runner) error {
	if err := runner.DecodeRuleConfig(r.Name(), &r.Config); err != nil {
		return err
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

// Name returns the rule name.
func (r *MetaRule) Name() string {
	return "eos_meta"
}

// Enabled returns whether the rule is enabled by default.
func (r *MetaRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *MetaRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}

// Link returns the rule reference link.
func (r *MetaRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_meta.md"
}
