// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package naming

import (
	"flag"
	"os"
	"testing"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func testNamingShoutRule(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	var config namingRuleConfig
	testhelper.LoadRuleConfig(t, "naming", &config)

	// Disable snake rule for shout test
	snakeDisabled := false
	config.Snake = &snakeDisabled

	cases := []struct {
		Name    string
		Content string
		Want    []string
	}{
		{
			Name: "shouted_names",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/naming_shout.tf")
				return string(content)
			}(),
			Want: []string{
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT🤡)",
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT)",
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"naming_shout.tf": tc.Content})
		rule := NewNamingRule()
		rule.Config = config

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
	}
}
