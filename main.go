// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"

	comment "github.com/staranto/tflint-ruleset-elements-of-style/rules/comment"
	deathmask "github.com/staranto/tflint-ruleset-elements-of-style/rules/death_mask"
	dry "github.com/staranto/tflint-ruleset-elements-of-style/rules/dry"
	heredoc "github.com/staranto/tflint-ruleset-elements-of-style/rules/heredoc"
	hungarian "github.com/staranto/tflint-ruleset-elements-of-style/rules/hungarian"
	meta "github.com/staranto/tflint-ruleset-elements-of-style/rules/meta"
	naming "github.com/staranto/tflint-ruleset-elements-of-style/rules/naming"
	reminder "github.com/staranto/tflint-ruleset-elements-of-style/rules/reminder"

	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// main is the entry point for the plugin.
func main() {
	log.SetFlags(0)
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "elements-of-style",
			Version: "0.0.2",
			Rules: []tflint.Rule{
				comment.NewCommentsRule(),
				deathmask.NewDeathMaskRule(),
				dry.NewDryRule(),
				heredoc.NewHeredocRule(),
				hungarian.NewHungarianRule(),
				meta.NewMetaRule(),
				naming.NewNamingRule(),
				reminder.NewReminderRule(),
			},
		},
	})
}
