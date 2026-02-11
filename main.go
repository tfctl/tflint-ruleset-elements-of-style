// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"

	"github.com/tfctl/tflint-ruleset-elements-of-style/rules/comment"
	deathmask "github.com/tfctl/tflint-ruleset-elements-of-style/rules/death_mask"
	"github.com/tfctl/tflint-ruleset-elements-of-style/rules/dry"
	"github.com/tfctl/tflint-ruleset-elements-of-style/rules/heredoc"
	"github.com/tfctl/tflint-ruleset-elements-of-style/rules/hungarian"
	"github.com/tfctl/tflint-ruleset-elements-of-style/rules/meta"
	"github.com/tfctl/tflint-ruleset-elements-of-style/rules/naming"
	"github.com/tfctl/tflint-ruleset-elements-of-style/rules/reminder"

	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// main is the entry point for the plugin.
func main() {
	log.SetFlags(0)
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "elements-of-style",
			Version: "1.0.0",
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
