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

func testNamingSnakeRule(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	var config namingRuleConfig
	testhelper.LoadRuleConfig(t, "naming", &config)
	// Enable snake for this test
	config.Snake = testhelper.BoolPtr(true)

	cases := []struct {
		Name    string
		Content string
		Want    []string
	}{
		{
			Name: "invalid_names",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/naming_snake.tf")
				return string(content)
			}(),
			Want: []string{
				"Names should be snake_case (CamelCase).",
				"Names should be snake_case (kebab-case).",
				"Names should be snake_case (with.dots).",
				"Names should be snake_case (CamelCase).",
				"Names should be snake_case (CamelCase).",
				"Names should be snake_case (kebab-case).",
				"Names should be snake_case (with.dots).",
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"naming_snake.tf": tc.Content})
		rule := NewNamingRule()
		rule.Config = config

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
	}
}
