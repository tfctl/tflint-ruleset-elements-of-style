// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package dry

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// dryConfig represents the configuration for the DryRule.
type dryConfig struct {
	Enabled   *bool  `hclext:"enabled,optional" hcl:"enabled,optional"`
	Level     string `hclext:"level,optional" hcl:"level,optional"`
	Threshold int    `hclext:"threshold,optional" hcl:"threshold,optional"`
}

// defaultDryConfig is the default configuration for the DryRule.
var defaultDryConfig = dryConfig{
	Enabled:   rulehelper.BoolPtr(true),
	Level:     "warning",
	Threshold: 2,
}

// DryRule checks for repeated interpolations.
type DryRule struct {
	tflint.DefaultRule
	Config dryConfig
	// RuleName is the rule block name to load from the config file. If empty,
	// defaults to "eos_dry".
	RuleName string
	// ConfigFile is the path to the config file. If empty, LoadRuleConfig will
	// search CWD then $HOME for .tflint.hcl.
	ConfigFile string
}

// Check checks whether the rule conditions are met.
func (r *DryRule) Check(runner tflint.Runner) error {
	// Load config using the rule name and optional config file path.
	if err := rulehelper.LoadRuleConfig(r.Name(), &r.Config, r.ConfigFile); err != nil {
		return err
	}

	// Bail out early if the rule is not enabled. This will occur if the EOS
	// plugin is enabled, but this specific rule is not.
	if !r.Enabled() {
		return nil
	}

	// Set a floor of 2 for threshold. It doesn't make sense to use < 2, so this
	// prevents a misconfiguration.
	threshold := r.Config.Threshold
	if threshold < 2 {
		threshold = 2
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
		if len(ranges) >= threshold {
			sort.Slice(ranges, func(i, j int) bool {
				return ranges[i].Start.Byte < ranges[j].Start.Byte
			})

			msg := fmt.Sprintf("Avoid repeating value '%s' %d times.", name, len(ranges))
			if strings.HasPrefix(name, "[") {
				msg = fmt.Sprintf("Avoid repeating list %d times.", len(ranges))
			} else if strings.HasPrefix(name, "{") {
				msg = fmt.Sprintf("Avoid repeating map %d times.", len(ranges))
			}

			if err := runner.EmitIssue(r, msg, ranges[0]); err != nil {
				return err
			}
		}
	}

	if err := r.checkDupe(runner, files, threshold); err != nil {
		return err
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
// the provided candidates map. THESE FUNCTIONS ARE ALMOST ENTIRELY AI-GENERATED
// with Gemini 3 Pro, which seems to produce less slop than some of it's peers.
// Update - not true, slop pervasive.
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

// checkDupe checks for duplicate resource and data blocks.
func (r *DryRule) checkDupe(runner tflint.Runner, files map[string]*hcl.File, threshold int) error {
	blockHashes := make(map[string][]hcl.Range)

	for filename, file := range files {
		if body, ok := file.Body.(*hclsyntax.Body); ok {
			r.collectBlocks(body, filename, file.Bytes, blockHashes)
		}
	}

	for _, ranges := range blockHashes {
		if len(ranges) >= threshold {
			sort.Slice(ranges, func(i, j int) bool {
				return ranges[i].Start.Byte < ranges[j].Start.Byte
			})

			msg := fmt.Sprintf("Duplicate block found %d times.", len(ranges))
			if err := runner.EmitIssue(r, msg, ranges[0]); err != nil {
				return err
			}
		}
	}

	return nil
}

// collectBlocks recursively collects resource and data blocks and computes their hashes.
func (r *DryRule) collectBlocks(body *hclsyntax.Body, filename string, fileBytes []byte, blockHashes map[string][]hcl.Range) {
	for _, block := range body.Blocks {
		if block.Type == "resource" || block.Type == "data" {
			hash := r.hashBlock(block, fileBytes)
			rng := block.Range()
			blockHashes[hash] = append(blockHashes[hash], rng)
		}
		r.collectBlocks(block.Body, filename, fileBytes, blockHashes)
	}
}

// hashBlock normalizes and hashes a block.
func (r *DryRule) hashBlock(block *hclsyntax.Block, fileBytes []byte) string {
	// Normalize by sorting attributes
	sortedAttrs := make([]string, 0, len(block.Body.Attributes))
	for name := range block.Body.Attributes {
		sortedAttrs = append(sortedAttrs, block.Body.Attributes[name].Name)
	}
	sort.Strings(sortedAttrs)

	var normalized strings.Builder
	normalized.WriteString(block.Type)
	normalized.WriteString(" {")

	for _, name := range sortedAttrs {
		attr := block.Body.Attributes[name]
		normalized.WriteString(name)
		normalized.WriteString("=")
		// Get the source text and minimize
		start := attr.Range().Start.Byte
		end := attr.Range().End.Byte
		source := string(fileBytes[start:end])
		minimized := r.minimizeSource(source)
		normalized.WriteString(minimized)
		normalized.WriteString(";")
	}

	normalized.WriteString("}")

	hash := sha256.Sum256([]byte(normalized.String()))
	return fmt.Sprintf("%x", hash)
}

// minimizeSource removes whitespace from source.
func (r *DryRule) minimizeSource(source string) string {
	var result strings.Builder
	for _, r := range source {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Enabled returns whether the rule is enabled by default.
func (r *DryRule) Enabled() bool {
	return r.Config.Enabled == nil || *r.Config.Enabled
}

// Link returns the rule link.
func (r *DryRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_dry.md"
}

// Name returns the rule name.
func (r *DryRule) Name() string {
	if r.RuleName != "" {
		return r.RuleName
	}
	return "eos_dry"
}

// Severity returns the rule severity.
func (r *DryRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
