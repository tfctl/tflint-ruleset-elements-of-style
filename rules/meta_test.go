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

var metaDeep = flag.Bool("metaDeep", false, "enable deep assert")

func TestMetaRule(t *testing.T) {
	flag.Parse()

	var config metaRuleConfig
	testhelper.LoadRuleConfig(t, "meta", &config)

	cases := []struct {
		Name    string
		Content string
		Want    helper.Issues
	}{
		{
			Name: "meta",
			Content: func() string {
				content, _ := os.ReadFile("testdata/meta_count_test.tf")
				return string(content)
			}(),
			Want: helper.Issues{
				{
					Rule:    NewMetaRule(),
					Message: "Avoid using count for anything other than dynamic guarding (condition ? 1 : 0)",
					Range: hcl.Range{
						Filename: "meta.tf",
						Start:    hcl.Pos{Line: 21, Column: 3},
						End:      hcl.Pos{Line: 21, Column: 12},
					},
				},
				{
					Rule:    NewMetaRule(),
					Message: "Count guard must return 1 or 0",
					Range: hcl.Range{
						Filename: "meta.tf",
						Start:    hcl.Pos{Line: 30, Column: 3},
						End:      hcl.Pos{Line: 30, Column: 34},
					},
				},
				{
					Rule:    NewMetaRule(),
					Message: "Count guard must return 1 or 0",
					Range: hcl.Range{
						Filename: "meta.tf",
						Start:    hcl.Pos{Line: 37, Column: 3},
						End:      hcl.Pos{Line: 37, Column: 34},
					},
				},
				{
					Rule:    NewMetaRule(),
					Message: "Avoid using count for anything other than dynamic guarding (condition ? 1 : 0)",
					Range: hcl.Range{
						Filename: "meta.tf",
						Start:    hcl.Pos{Line: 44, Column: 3},
						End:      hcl.Pos{Line: 44, Column: 25},
					},
				},
				{
					Rule:    NewMetaRule(),
					Message: "Avoid using count for anything other than dynamic guarding (condition ? 1 : 0)",
					Range: hcl.Range{
						Filename: "meta.tf",
						Start:    hcl.Pos{Line: 50, Column: 3},
						End:      hcl.Pos{Line: 50, Column: 25},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"meta.tf": tc.Content})
		rule := NewMetaRule()
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
			if *metaDeep {
				helper.AssertIssues(t, tc.Want, runner.Issues)
			} else {
				helper.AssertIssuesWithoutRange(t, tc.Want, runner.Issues)
			}
		})
	}
}
