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
var defaultReminderTags = []string{"BUG", "TODO"}

// reminderConfig represents the configuration for the ReminderRule.
type reminderConfig struct {
	Enabled *bool    `hclext:"enabled,optional" hcl:"enabled,optional"`
	Extras  []string `hclext:"extras,optional" hcl:"extras,optional"`
	Tags    []string `hclext:"tags,optional" hcl:"tags,optional"`
	Level   string   `hclext:"level,optional" hcl:"level,optional"`
}

// defaultReminderConfig is the default configuration for the ReminderRule.
var defaultReminderConfig = reminderConfig{
	Enabled: rulehelper.BoolPtr(true),
	Extras:  []string{},
	Tags:    defaultReminderTags,
	Level:   "warning",
}

// ReminderRule checks for reminders.
type ReminderRule struct {
	tflint.DefaultRule
	Config reminderConfig
	// RuleName is the rule block name to load from the config file. If empty,
	// defaults to "eos_reminder".
	RuleName string
	// ConfigFile is the path to the config file. If empty, LoadRuleConfig will
	// search CWD then $HOME for .tflint.hcl.
	ConfigFile string
}

// Check checks whether the rule conditions are met.
func (r *ReminderRule) Check(runner tflint.Runner) error {
	// Load config using the rule name and optional config file path.
	if err := rulehelper.LoadRuleConfig(r.Name(), &r.Config, r.ConfigFile); err != nil {
		return err
	}

	// Bail out early if the rule is not enabled. This will occur if the EOS
	// plugin is enabled, but this specific rule is not.
	if !r.Enabled() {
		return nil
	}

	return rulehelper.WalkTokens(runner, r, checkReminder)
}

func checkReminder(runner tflint.Runner, r *ReminderRule, token hclsyntax.Token) {
	if token.Type != hclsyntax.TokenComment {
		return
	}
	text := string(token.Bytes)

	tags := r.Config.Tags
	extras := r.Config.Extras
	if len(extras) > 0 {
		tags = append(tags, extras...)
	}
	// TODO Dedup tags.

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
	return r.Config.Enabled == nil || *r.Config.Enabled
}

// Link returns the rule reference link.
func (r *ReminderRule) Link() string {
	return "https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_reminder.md"
}

// Name returns the rule name.
func (r *ReminderRule) Name() string {
	if r.RuleName != "" {
		return r.RuleName
	}
	return "eos_reminder"
}

// Severity returns the rule severity.
func (r *ReminderRule) Severity() tflint.Severity {
	return rulehelper.ToSeverity(r.Config.Level)
}
