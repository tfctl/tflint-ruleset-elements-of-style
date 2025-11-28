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

var deathMaskDeep = flag.Bool("deathMaskDeep", false, "enable deep assert")

func TestDeathMaskRule(t *testing.T) {
	flag.Parse()

	var config deathMaskRuleConfig
	testhelper.LoadRuleConfig(t, "death_mask", &config)

	cases := []struct {
		Name    string
		Content string
		Want    helper.Issues
	}{
		{
			Name: "death_mask",
			Content: func() string {
				content, _ := os.ReadFile("testdata/death_mask.tf")
				return string(content)
			}(),
			Want: helper.Issues{
				{
					Rule:    NewDeathMaskRule(),
					Message: "Avoid commented-out code.",
					Range: hcl.Range{
						Filename: "death_mask.tf",
						Start:    hcl.Pos{Line: 10, Column: 1},
						End:      hcl.Pos{Line: 10, Column: 8},
					},
				},
				{
					Rule:    NewDeathMaskRule(),
					Message: "Avoid commented-out code.",
					Range: hcl.Range{
						Filename: "death_mask.tf",
						Start:    hcl.Pos{Line: 14, Column: 1},
						End:      hcl.Pos{Line: 16, Column: 4},
					},
				},
				{
					Rule:    NewDeathMaskRule(),
					Message: "Avoid commented-out code.",
					Range: hcl.Range{
						Filename: "death_mask.tf",
						Start:    hcl.Pos{Line: 21, Column: 3},
						End:      hcl.Pos{Line: 21, Column: 10},
					},
				},
				{
					Rule:    NewDeathMaskRule(),
					Message: "Avoid commented-out code.",
					Range: hcl.Range{
						Filename: "death_mask.tf",
						Start:    hcl.Pos{Line: 26, Column: 1},
						End:      hcl.Pos{Line: 29, Column: 4},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"death_mask.tf": tc.Content})
		rule := NewDeathMaskRule()
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
			if *deathMaskDeep {
				helper.AssertIssues(t, tc.Want, runner.Issues)
			} else {
				helper.AssertIssuesWithoutRange(t, tc.Want, runner.Issues)
			}
		})
	}
}
