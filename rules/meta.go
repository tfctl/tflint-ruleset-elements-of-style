// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// defaultMetaConfig is the default configuration for the MetaRule.
var defaultMetaConfig = metaRuleConfig{
	Level: "warning",
}

// metaRuleConfig represents the configuration for the MetaRule.
type metaRuleConfig struct {
	Enabled *bool  `hcl:"enabled,optional"`
	Level   string `hcl:"level,optional"`
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
			}
		}
	}

	return nil
}

func checkCountGuard(runner tflint.Runner, r *MetaRule, attr *hclsyntax.Attribute) {

	expr := attr.Expr
	condExpr, ok := expr.(*hclsyntax.ConditionalExpr)

	// We want to check if it is a conditional expression:
	// condition ? true_val : false_val
	if !ok {
		// Allow literal 0 or 1.
		if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
			val := getLiteralValue(lit)
			if val == 0 || val == 1 {
				return
			}
		}

		r.emitIssue(runner, "Avoid using count for anything other than dynamic guarding (condition ? 1 : 0)", attr.Range())
		return
	}

	// Check true/false results.
	if !isValidGuardResult(condExpr.TrueResult) || !isValidGuardResult(condExpr.FalseResult) {
		r.emitIssue(runner, "Count guard must return 1 or 0", attr.Range())
	}
}

func isValidGuardResult(expr hclsyntax.Expression) bool {
	val := getLiteralValue(expr)
	return val == 0 || val == 1
}

func getLiteralValue(expr hclsyntax.Expression) int {
	if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
		if lit.Val.Type() == cty.Number {
			f, _ := lit.Val.AsBigFloat().Float64()
			if f == 0 {
				return 0
			}
			if f == 1 {
				return 1
			}
		}
	}
	return -1
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
