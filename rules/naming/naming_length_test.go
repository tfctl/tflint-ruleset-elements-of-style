// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package naming

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func testNamingLengthRule(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	var config namingRuleConfig
	testhelper.LoadRuleConfig(t, "naming", &config)

	limit := defaultLimit
	if config.Length != nil {
		limit = *config.Length
	}
	cases := []struct {
		Name    string
		Content string
		Want    []string
	}{
		{
			Name: "long_names",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/naming_length.tf")
				return string(content)
			}(),
			Want: []string{
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name_disabled' is 32).", limit),
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"naming_length.tf": tc.Content})
		rule := NewNamingRule()
		rule.Config = config

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
	}
}
