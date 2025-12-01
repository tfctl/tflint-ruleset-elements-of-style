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

// typeEchoRuleConfig represents the configuration for the TypeEchoRule.
type typeEchoRuleConfig struct {
	Enabled  *bool               `hclext:"enabled,optional" hcl:"enabled,optional"`
	Synonyms map[string][]string `hclext:"synonyms,optional" hcl:"synonyms,optional"`
	Level    string              `hclext:"level,optional" hcl:"level,optional"`
}

// defaultTypeEchoConfig is the default configuration for the TypeEchoRule.
var defaultTypeEchoConfig = typeEchoRuleConfig{
	Level: "warning",
}

// TypeEchoRule checks whether a block's type is echoed in its name.
type TypeEchoRule struct {
	tflint.DefaultRule
	Config typeEchoRuleConfig
}

// Check checks whether the rule conditions are met.
func (r *TypeEchoRule) Check(runner tflint.Runner) error {
	if err := runner.DecodeRuleConfig(r.Name(), &r.Config); err != nil {
		return err
	}

	return rulehelper.WalkBlocks(runner, rulehelper.AllLintableBlocks, r, checkForEcho)
}

// checkForEcho checks if a word in type is echoed in the name.
func checkForEcho(runner tflint.Runner,
	r *TypeEchoRule, defRange hcl.Range,
	typ string, name string, synonym string) {

	// Assume there is no echo.
	echo := false

	lowerTyp := strings.ToLower(typ)   // aws_s3_bucket
	lowerName := strings.ToLower(name) // my_bucket
	synonymText := ""

	// For each word in type, see if it exists in name. Note that this is
	// comparing against the entire name value. The impact of that being a name of
	// "mys3widget" will match against a type of "aws_s3_bucket" (the "s3" word
	// exists in the name).
	for part := range strings.SplitSeq(lowerTyp, "_") {
		if strings.Contains(lowerName, part) {
			echo = true
			break
		}

		// Get synonyms for the word.
		synonyms := r.Config.Synonyms[part]
		if synonym != "" {
			synonyms = append(synonyms, synonym)
		}

		// Check synonyms.  This logic is different than above in that synonyms are
		// checked for on word boundaries. So "aws_s3_bucket" DOES NOT match
		// "mys3_widget", but would match "my_s3_widget".
		splitName := strings.SplitSeq(lowerName, "_-")
		for _, syn := range synonyms {
			for n := range splitName {
				if strings.Contains(n, syn) {
					echo = true
					synonymText = fmt.Sprintf(" (via synonym '%s')", syn)
					break
				}
			}

			// Don't bother checking more synonyms because we know we already have a
			// winner.
			if echo {
				break
			}
		}
	}

	if echo {
		if err := runner.EmitIssue(
			r,
			fmt.Sprintf("Avoid echoing type \"%s\"%s in label \"%s\".", typ, synonymText, name),
			defRange,
		); err != nil {
			logger.Error(err.Error())
		}
	}
}

// NewTypeEchoRule returns a new rule.
func NewTypeEchoRule() *TypeEchoRule {
	rule := &TypeEchoRule{}
	rule.Config = defaultTypeEchoConfig
	return rule
}

// Enabled returns whether the rule is enabled by default.
func (r *TypeEchoRule) Enabled() bool {
	return true
}

// Link returns the rule reference link.
func (r *TypeEchoRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_type_echo.md"
}

// Name returns the rule name.
func (r *TypeEchoRule) Name() string {
	return "eos_type_echo"
}

// Severity returns the rule severity.
func (r *TypeEchoRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
