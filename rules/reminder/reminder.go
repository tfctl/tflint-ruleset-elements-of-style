// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package reminder

import (
	"fmt"
	"strings"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// defaultReminderTags is the default list of tags to check for.
var defaultReminderTags = []string{
	"BUG",
	"FIXME",
	"HACK",
	"TODO",
}

// defaultReminderConfig is the default configuration for the ReminderRule.
var defaultReminderConfig = reminderRuleConfig{
	Tags:  defaultReminderTags,
	Level: "warning",
}

// reminderRuleConfig represents the configuration for the ReminderRule.
type reminderRuleConfig struct {
	Enabled *bool    `hclext:"enabled,optional" hcl:"enabled,optional"`
	Tags    []string `hclext:"tags,optional" hcl:"tags,optional"`
	Level   string   `hclext:"level,optional" hcl:"level,optional"`
}

// ReminderRule checks for reminders.
type ReminderRule struct {
	tflint.DefaultRule
	Config reminderRuleConfig
}

// Check checks whether the rule conditions are met.
func (r *ReminderRule) Check(runner tflint.Runner) error {
	if err := runner.DecodeRuleConfig(r.Name(), &r.Config); err != nil {
		return err
	}

	return rulehelper.WalkTokens(runner, r, checkReminder)
}

func checkReminder(runner tflint.Runner, r *ReminderRule, token hclsyntax.Token) {
	if token.Type != hclsyntax.TokenComment {
		return
	}
	text := string(token.Bytes)
	tags := r.Config.Tags

	tokens := strings.SplitAfterN(strings.ToUpper(text), " ", 2)
	if len(tokens) < 2 {
		return
	}

	for _, t := range tags {
		if strings.HasSuffix(strings.TrimSpace(tokens[0]), t) || strings.HasPrefix(tokens[1], t) {
			message := fmt.Sprintf("Resolve reminder: '%s'.", strings.TrimSpace(text))
			if err := runner.EmitIssue(r, message, token.Range); err != nil {
				logger.Error(err.Error())
			}
			logger.Debug(message)
		}
	}
}

// NewReminderRule returns a new rule.
func NewReminderRule() *ReminderRule {
	rule := &ReminderRule{}
	rule.Config = defaultReminderConfig
	return rule
}

// Enabled returns whether the rule is enabled by default.
func (r *ReminderRule) Enabled() bool {
	return true
}

// Link returns the rule reference link.
func (r *ReminderRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_reminder.md"
}

// Name returns the rule name.
func (r *ReminderRule) Name() string {
	return "eos_reminder"
}

// Severity returns the rule severity.
func (r *ReminderRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
