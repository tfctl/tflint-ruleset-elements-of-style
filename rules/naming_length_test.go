// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func testNamingLengthRule(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	var config namingRuleConfig
	testhelper.LoadRuleConfig(t, "naming", &config)

	limit := config.Length.Limit
	cases := []struct {
		Name    string
		Content string
		Want    helper.Issues
	}{
		{
			Name: "long_names",
			Content: func() string {
				content, _ := os.ReadFile("testdata/naming_length.tf")
				return string(content)
			}(),
			Want: helper.Issues{
				{
					Rule:    NewNamingRule(),
					Message: fmt.Sprintf("Avoid names longer than %d ('very_long_variable_name' is 23).", limit),
					Range: hcl.Range{
						Filename: "naming_length.tf",
						Start:    hcl.Pos{Line: 7, Column: 10},
						End:      hcl.Pos{Line: 7, Column: 35},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: fmt.Sprintf("Avoid names longer than %d ('very_long_local_name' is 20).", limit),
					Range: hcl.Range{
						Filename: "naming_length.tf",
						Start:    hcl.Pos{Line: 10, Column: 3},
						End:      hcl.Pos{Line: 10, Column: 27},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: fmt.Sprintf("Avoid names longer than %d ('very_long_check_name' is 20).", limit),
					Range: hcl.Range{
						Filename: "naming_length.tf",
						Start:    hcl.Pos{Line: 14, Column: 7},
						End:      hcl.Pos{Line: 14, Column: 29},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: fmt.Sprintf("Avoid names longer than %d ('very_long_data_name' is 19).", limit),
					Range: hcl.Range{
						Filename: "naming_length.tf",
						Start:    hcl.Pos{Line: 21, Column: 32},
						End:      hcl.Pos{Line: 21, Column: 53},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: fmt.Sprintf("Avoid names longer than %d ('very_long_ephemeral_name' is 24).", limit),
					Range: hcl.Range{
						Filename: "naming_length.tf",
						Start:    hcl.Pos{Line: 23, Column: 29},
						End:      hcl.Pos{Line: 23, Column: 55},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: fmt.Sprintf("Avoid names longer than %d ('very_long_module_name' is 21).", limit),
					Range: hcl.Range{
						Filename: "naming_length.tf",
						Start:    hcl.Pos{Line: 27, Column: 8},
						End:      hcl.Pos{Line: 27, Column: 31},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: fmt.Sprintf("Avoid names longer than %d ('very_long_output_name' is 21).", limit),
					Range: hcl.Range{
						Filename: "naming_length.tf",
						Start:    hcl.Pos{Line: 31, Column: 8},
						End:      hcl.Pos{Line: 31, Column: 31},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: fmt.Sprintf("Avoid names longer than %d ('very_long_instance_name' is 23).", limit),
					Range: hcl.Range{
						Filename: "naming_length.tf",
						Start:    hcl.Pos{Line: 35, Column: 25},
						End:      hcl.Pos{Line: 35, Column: 50},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: fmt.Sprintf("Avoid names longer than %d ('very_long_instance_name_disabled' is 32).", limit),
					Range: hcl.Range{
						Filename: "naming_length.tf",
						Start:    hcl.Pos{Line: 40, Column: 25},
						End:      hcl.Pos{Line: 40, Column: 57},
					},
				},
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

		if len(runner.Issues) != len(tc.Want) {
			t.Logf("Expected %d issues, got %d", len(tc.Want), len(runner.Issues))
			for i, issue := range runner.Issues {
				t.Logf("Issue %d: %s at %s", i, issue.Message, issue.Range)
			}
			t.Fatalf("Number of issues mismatch: got %d, want %d", len(runner.Issues), len(tc.Want))
		}

		t.Run(tc.Name, func(t *testing.T) {
			if *namingDeep {
				helper.AssertIssues(t, tc.Want, runner.Issues)
			} else {
				helper.AssertIssuesWithoutRange(t, tc.Want, runner.Issues)
			}
		})
	}
}
