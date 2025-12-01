// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"flag"
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func testNamingShoutRule(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	var config namingRuleConfig
	testhelper.LoadRuleConfig(t, "naming", &config)

	cases := []struct {
		Name    string
		Content string
		Want    helper.Issues
	}{
		{
			Name: "shouted_names",
			Content: func() string {
				content, _ := os.ReadFile("testdata/naming_shout.tf")
				return string(content)
			}(),
			Want: helper.Issues{
				{
					Rule:    NewNamingRule(),
					Message: "Avoid SHOUTED names (SHOUT)",
					Range: hcl.Range{
						Filename: "naming_shout.tf",
						Start:    hcl.Pos{Line: 7, Column: 10},
						End:      hcl.Pos{Line: 7, Column: 17},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: "Avoid SHOUTED names (SHOUT🤡)",
					Range: hcl.Range{
						Filename: "naming_shout.tf",
						Start:    hcl.Pos{Line: 8, Column: 10},
						End:      hcl.Pos{Line: 8, Column: 18},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: "Avoid SHOUTED names (SHOUT)",
					Range: hcl.Range{
						Filename: "naming_shout.tf",
						Start:    hcl.Pos{Line: 10, Column: 3},
						End:      hcl.Pos{Line: 10, Column: 8},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: "Avoid SHOUTED names (SHOUT)",
					Range: hcl.Range{
						Filename: "naming_shout.tf",
						Start:    hcl.Pos{Line: 13, Column: 7},
						End:      hcl.Pos{Line: 13, Column: 14},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: "Avoid SHOUTED names (SHOUT)",
					Range: hcl.Range{
						Filename: "naming_shout.tf",
						Start:    hcl.Pos{Line: 20, Column: 28},
						End:      hcl.Pos{Line: 20, Column: 35},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: "Avoid SHOUTED names (SHOUT)",
					Range: hcl.Range{
						Filename: "naming_shout.tf",
						Start:    hcl.Pos{Line: 22, Column: 29},
						End:      hcl.Pos{Line: 22, Column: 36},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: "Avoid SHOUTED names (SHOUT)",
					Range: hcl.Range{
						Filename: "naming_shout.tf",
						Start:    hcl.Pos{Line: 26, Column: 8},
						End:      hcl.Pos{Line: 26, Column: 15},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: "Avoid SHOUTED names (SHOUT)",
					Range: hcl.Range{
						Filename: "naming_shout.tf",
						Start:    hcl.Pos{Line: 30, Column: 8},
						End:      hcl.Pos{Line: 30, Column: 15},
					},
				},
				{
					Rule:    NewNamingRule(),
					Message: "Avoid SHOUTED names (SHOUT)",
					Range: hcl.Range{
						Filename: "naming_shout.tf",
						Start:    hcl.Pos{Line: 35, Column: 25},
						End:      hcl.Pos{Line: 35, Column: 32},
					},
				},
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
