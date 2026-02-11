// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package death_mask

import (
	"strings"

	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AvoidDeathMaskMessage is the message emitted when commented-out code is
// detected.
const AvoidDeathMaskMessage = "Avoid commented-out code."

// deathMaskConfig represents the configuration for the DeathMaskRule.
type deathMaskConfig struct {
	Enabled *bool  `hclext:"enabled,optional" hcl:"enabled,optional"`
	Level   string `hclext:"level,optional" hcl:"level,optional"`
}

// defaultDeathMaskConfig is the default configuration for the DeathMaskRule.
var defaultDeathMaskConfig = deathMaskConfig{
	Enabled: rulehelper.BoolPtr(true),
	Level:   "warning",
}

// Rule checks for commented-out code.
type Rule struct {
	tflint.DefaultRule
	Config deathMaskConfig
	// RuleName is the rule block name to load from the config file. If empty,
	// defaults to "eos_death_mask".
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

	for name, file := range files {
		if err := r.checkDeathMask(runner, name, file); err != nil {
			return err
		}
	}

	return nil
}

// checkDeathMask checks for commented-out code in a file.
func (r *Rule) checkDeathMask(runner tflint.Runner, filename string, file *hcl.File) error {
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
func (r *Rule) processCommentBlock(runner tflint.Runner, tokens []hclsyntax.Token) {
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

				message := AvoidDeathMaskMessage
				if err := runner.EmitIssue(r, message, issueRange); err != nil {
					logger.Error(err.Error())
				}
				return // We found a match, so we stop checking this block.
			}
		}
	}
}

// NewDeathMaskRule returns a new rule whose config is set to the default.
func NewDeathMaskRule() *Rule {
	rule := &Rule{}
	rule.Config = defaultDeathMaskConfig
	return rule
}

// Enabled returns whether the rule is enabled by default.
func (r *Rule) Enabled() bool {
	return r.Config.Enabled == nil || *r.Config.Enabled
}

// Link returns the rule reference link.
func (r *Rule) Link() string {
	return "https://github.com/tfctl/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_death_mask.md"
}

// Name returns the rule name.
func (r *Rule) Name() string {
	if r.RuleName != "" {
		return r.RuleName
	}
	return "eos_death_mask"
}

// Severity returns the rule severity.
func (r *Rule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
