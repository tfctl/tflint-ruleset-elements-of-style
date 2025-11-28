// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"strings"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// defaultDeathMaskConfig is the default configuration for the DeathMaskRule.
var defaultDeathMaskConfig = deathMaskRuleConfig{
	Level: "warning",
}

// deathMaskRuleConfig represents the configuration for the DeathMaskRule.
type deathMaskRuleConfig struct {
	Enabled *bool  `hclext:"enabled,optional" hcl:"enabled,optional"`
	Level   string `hclext:"level,optional" hcl:"level,optional"`
}

// DeathMaskRule checks for commented-out code.
type DeathMaskRule struct {
	tflint.DefaultRule
	Config deathMaskRuleConfig
}

// Check checks whether the rule conditions are met.
func (r *DeathMaskRule) Check(runner tflint.Runner) error {
	if err := runner.DecodeRuleConfig(r.Name(), &r.Config); err != nil {
		return err
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for name, file := range files {
		if err := r.checkDeathMask(runner, name, file); err != nil {
			return err
		}
	}

	return nil
}

// checkDeathMask checks for commented-out code in a file.
func (r *DeathMaskRule) checkDeathMask(runner tflint.Runner, filename string, file *hcl.File) error {
	tokens, diags := hclsyntax.LexConfig(file.Bytes, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return diags
	}

	var commentBlock []hclsyntax.Token

	for _, token := range tokens {
		switch token.Type {
		case hclsyntax.TokenComment:
			// If this is the first comment or it follows the previous one, we
			// add it.
			// Check adjacency.
			if len(commentBlock) > 0 {
				last := commentBlock[len(commentBlock)-1]
				// Check if this token is on the next line or same line.
				if token.Range.Start.Line > last.Range.End.Line {
					// Detected a gap, so flush the previous block.
					r.processCommentBlock(runner, commentBlock)
					commentBlock = nil
				}
			}
			commentBlock = append(commentBlock, token)
		case hclsyntax.TokenNewline, hclsyntax.TokenEOF:
			// Continue and let newlines pass.
		default:
			// A non-comment, non-newline token breaks the block.
			if len(commentBlock) > 0 {
				r.processCommentBlock(runner, commentBlock)
				commentBlock = nil
			}
		}
	}

	// Flush the remaining tokens.
	if len(commentBlock) > 0 {
		r.processCommentBlock(runner, commentBlock)
	}

	return nil
}

// processCommentBlock unwraps and validates a block of comments.
func (r *DeathMaskRule) processCommentBlock(runner tflint.Runner, tokens []hclsyntax.Token) {
	var lines []string
	for _, token := range tokens {
		text := string(token.Bytes)

		if s, cut := strings.CutPrefix(text, "//"); cut {
			s = strings.TrimPrefix(s, " ")
			lines = append(lines, s)
			continue
		}

		if s, cut := strings.CutPrefix(text, "#"); cut {
			s = strings.TrimPrefix(s, " ")
			lines = append(lines, s)
			continue
		}

		if s, cut := strings.CutPrefix(text, "/*"); cut {
			s = strings.TrimPrefix(s, "/*")
			s = strings.TrimSuffix(s, "*/")
			// Split the block comment into lines.
			blockLines := strings.Split(s, "\n")
			lines = append(lines, blockLines...)
		}
	}

	// Try to parse subsets of lines to handle header text.
	for i := 0; i < len(lines); i++ {
		candidate := strings.Join(lines[i:], "\n")
		// Parse the candidate string as HCL.
		file, diags := hclsyntax.ParseConfig([]byte(candidate), "candidate.tf", hcl.InitialPos)
		if diags.HasErrors() {
			continue
		}

		// Check if it actually contains code (Attributes or Blocks).
		if body, ok := file.Body.(*hclsyntax.Body); ok {
			if len(body.Attributes) > 0 || len(body.Blocks) > 0 {
				// It is valid code. Flag the whole block.
				start := tokens[0].Range.Start
				end := tokens[len(tokens)-1].Range.End
				issueRange := hcl.Range{
					Filename: tokens[0].Range.Filename,
					Start:    start,
					End:      end,
				}

				message := "Avoid commented-out code."
				if err := runner.EmitIssue(r, message, issueRange); err != nil {
					logger.Error(err.Error())
				}
				return // We found a match, so we stop checking this block.
			}
		}
	}
}

// NewDeathMaskRule returns a new rule.
func NewDeathMaskRule() *DeathMaskRule {
	rule := &DeathMaskRule{}
	rule.Config = defaultDeathMaskConfig
	return rule
}

// Enabled returns whether the rule is enabled by default.
func (r *DeathMaskRule) Enabled() bool {
	return true
}

// Link returns the rule reference link.
func (r *DeathMaskRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_death_mask.md"
}

// Name returns the rule name.
func (r *DeathMaskRule) Name() string {
	return "eos_death_mask"
}

// Severity returns the rule severity.
func (r *DeathMaskRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
