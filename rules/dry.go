// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// defaultDryConfig is the default configuration for the DryRule.
var defaultDryConfig = dryRuleConfig{
	Level: "warning",
}

// dryRuleConfig represents the configuration for the rule.
type dryRuleConfig struct {
	Enabled *bool  `hclext:"enabled,optional" hcl:"enabled,optional"`
	Level   string `hclext:"level,optional" hcl:"level,optional"`
}

// DryRule checks for repeated interpolations.
type DryRule struct {
	tflint.DefaultRule
	Config dryRuleConfig
}

// Check checks whether the rule conditions are met.
func (r *DryRule) Check(runner tflint.Runner) error {
	if err := runner.DecodeRuleConfig(r.Name(), &r.Config); err != nil {
		return err
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	candidates := make(map[string][]hcl.Range)

	for filename, file := range files {
		if body, ok := file.Body.(*hclsyntax.Body); ok {
			r.checkDry(body, filename, file.Bytes, candidates, false)
		}
	}

	for name, ranges := range candidates {
		if len(ranges) > 1 {
			sort.Slice(ranges, func(i, j int) bool {
				return ranges[i].Start.Byte < ranges[j].Start.Byte
			})

			msg := fmt.Sprintf("Avoid repeating value '%s' %d times.", name, len(ranges))
			if strings.HasPrefix(name, "[") {
				msg = fmt.Sprintf("Avoid repeating list %d times.", len(ranges))
			} else if strings.HasPrefix(name, "{") {
				msg = fmt.Sprintf("Avoid repeating map %d times.", len(ranges))
			}

			if err := runner.EmitIssue(
				r,
				msg,
				ranges[0],
			); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkDry recursively checks for repeated expressions in the body.
func (r *DryRule) checkDry(body *hclsyntax.Body, filename string, fileBytes []byte, candidates map[string][]hcl.Range, isModule bool) {
	for name, attr := range body.Attributes {
		if isModule && name == "source" {
			continue
		}
		r.walkExpression(attr.Expr, filename, fileBytes, candidates)
	}
	for _, block := range body.Blocks {
		r.checkDry(block.Body, filename, fileBytes, candidates, block.Type == "module")
	}
}

// walkExpression recursively traverses an hclsyntax.Expression tree and
// collects candidate expression source snippets and their byte ranges into
// the provided candidates map.
// THIS FUNCTION IS ALMOST ENTIRELY AI-GENERATED with Gemini 3 Pro, which seems
// to produce far less slop than some of it's peers.
func (r *DryRule) walkExpression(expr hclsyntax.Expression, filename string, fileBytes []byte, candidates map[string][]hcl.Range) {
	// HCL does not provide a generic "GetChildren()" method for expressions.
	// Instead, each expression type stores its sub-expressions in different
	// fields. We have to check the type to know which fields to walk
	// recursively.
	switch t := expr.(type) {
	case *hclsyntax.TemplateWrapExpr:
		r.walkExpression(t.Wrapped, filename, fileBytes, candidates)
	case *hclsyntax.TemplateExpr:
		if !r.hasCountOrEach(t) {
			rng := t.Range()
			if rng.Start.Byte < len(fileBytes) && rng.End.Byte <= len(fileBytes) {
				source := string(fileBytes[rng.Start.Byte:rng.End.Byte])
				candidates[source] = append(candidates[source], rng)
			}
		}

		// Recurse into the parts of the template. For example, in
		// "foo ${func(var.a)}", we need to walk into the function call.
		for _, part := range t.Parts {
			r.walkExpression(part, filename, fileBytes, candidates)
		}
	case *hclsyntax.ScopeTraversalExpr:
		// No longer track bare traversals (e.g. var.foo).
		// Only care about interpolations inside strings.
	case *hclsyntax.FunctionCallExpr:
		// Recurse into function arguments.
		for _, arg := range t.Args {
			r.walkExpression(arg, filename, fileBytes, candidates)
		}
	case *hclsyntax.TupleConsExpr:
		if !r.hasCountOrEach(t) {
			rng := t.Range()
			if rng.Start.Byte < len(fileBytes) && rng.End.Byte <= len(fileBytes) {
				source := string(fileBytes[rng.Start.Byte:rng.End.Byte])
				candidates[source] = append(candidates[source], rng)
			}
		}
	case *hclsyntax.ObjectConsExpr:
		if !r.hasCountOrEach(t) {
			rng := t.Range()
			if rng.Start.Byte < len(fileBytes) && rng.End.Byte <= len(fileBytes) {
				source := string(fileBytes[rng.Start.Byte:rng.End.Byte])
				candidates[source] = append(candidates[source], rng)
			}
		}
	case *hclsyntax.ConditionalExpr:
		// Recurse into condition, true result, and false result.
		r.walkExpression(t.Condition, filename, fileBytes, candidates)
		r.walkExpression(t.TrueResult, filename, fileBytes, candidates)
		r.walkExpression(t.FalseResult, filename, fileBytes, candidates)
	case *hclsyntax.ForExpr:
		if !r.hasCountOrEach(t) {
			rng := t.Range()
			if rng.Start.Byte < len(fileBytes) && rng.End.Byte <= len(fileBytes) {
				source := string(fileBytes[rng.Start.Byte:rng.End.Byte])
				candidates[source] = append(candidates[source], rng)
			}
		}

		// Recurse into for loop components.
		r.walkExpression(t.CollExpr, filename, fileBytes, candidates)
		r.walkExpression(t.KeyExpr, filename, fileBytes, candidates)
		r.walkExpression(t.ValExpr, filename, fileBytes, candidates)
		if t.CondExpr != nil {
			r.walkExpression(t.CondExpr, filename, fileBytes, candidates)
		}
	case *hclsyntax.SplatExpr:
		r.walkExpression(t.Source, filename, fileBytes, candidates)
	case *hclsyntax.ParenthesesExpr:
		r.walkExpression(t.Expression, filename, fileBytes, candidates)
	case *hclsyntax.UnaryOpExpr:
		r.walkExpression(t.Val, filename, fileBytes, candidates)
	case *hclsyntax.BinaryOpExpr:
		r.walkExpression(t.LHS, filename, fileBytes, candidates)
		r.walkExpression(t.RHS, filename, fileBytes, candidates)
	case *hclsyntax.IndexExpr:
		r.walkExpression(t.Collection, filename, fileBytes, candidates)
		r.walkExpression(t.Key, filename, fileBytes, candidates)
	}
}

// hasCountOrEach reports whether the given HCL expression contains a scope traversal
// whose root identifier is "count" or "each".
func (r *DryRule) hasCountOrEach(expr hclsyntax.Expression) bool {
	found := false
	var check func(hclsyntax.Expression)
	check = func(e hclsyntax.Expression) {
		if found {
			return
		}
		switch t := e.(type) {
		case *hclsyntax.ScopeTraversalExpr:
			if len(t.Traversal) > 0 {
				if root, ok := t.Traversal[0].(hcl.TraverseRoot); ok {
					if root.Name == "count" || root.Name == "each" {
						found = true
					}
				}
			}
		case *hclsyntax.TemplateExpr:
			for _, part := range t.Parts {
				check(part)
			}
		case *hclsyntax.TemplateWrapExpr:
			check(t.Wrapped)
		case *hclsyntax.FunctionCallExpr:
			for _, arg := range t.Args {
				check(arg)
			}
		case *hclsyntax.TupleConsExpr:
			for _, item := range t.Exprs {
				check(item)
			}
		case *hclsyntax.ObjectConsExpr:
			for _, item := range t.Items {
				check(item.KeyExpr)
				check(item.ValueExpr)
			}
		case *hclsyntax.ConditionalExpr:
			check(t.Condition)
			check(t.TrueResult)
			check(t.FalseResult)
		case *hclsyntax.ForExpr:
			check(t.CollExpr)
			check(t.KeyExpr)
			check(t.ValExpr)
			if t.CondExpr != nil {
				check(t.CondExpr)
			}
		case *hclsyntax.SplatExpr:
			check(t.Source)
		case *hclsyntax.ParenthesesExpr:
			check(t.Expression)
		case *hclsyntax.UnaryOpExpr:
			check(t.Val)
		case *hclsyntax.BinaryOpExpr:
			check(t.LHS)
			check(t.RHS)
		case *hclsyntax.IndexExpr:
			check(t.Collection)
			check(t.Key)
		}
	}
	check(expr)
	return found
}

// NewDryRule returns a new rule.
func NewDryRule() *DryRule {
	rule := &DryRule{}
	rule.Config = defaultDryConfig
	return rule
}

// Name returns the rule name.
func (r *DryRule) Name() string {
	return "eos_dry"
}

// Enabled returns whether the rule is enabled by default.
func (r *DryRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *DryRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}

// Link returns the rule link.
func (r *DryRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_dry.md"
}
