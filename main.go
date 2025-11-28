// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"

	"github.com/staranto/tflint-ruleset-elements-of-style/rules"

	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// main is the entry point for the plugin.
func main() {
	log.SetFlags(0)
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "elements-of-style",
			Version: "0.0.1",
			Rules: []tflint.Rule{
				rules.NewCommentsRule(),
				rules.NewDeathMaskRule(),
				rules.NewDryRule(),
				rules.NewHeredocRule(),
				rules.NewHungarianRule(),
				rules.NewMetaRule(),
				rules.NewNamingRule(),
				rules.NewReminderRule(),
				rules.NewTypeEchoRule(),
			},
		},
	})
}
