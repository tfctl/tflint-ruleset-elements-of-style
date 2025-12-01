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

func testCommentsJammedRule(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	var config commentsRuleConfig
	testhelper.LoadRuleConfig(t, "comments", &config)

	cases := []struct {
		Name    string
		Content string
		Want    helper.Issues
	}{
		{
			Name: "jammed_comments",
			Content: func() string {
				content, _ := os.ReadFile("testdata/comments_jammed.tf")
				return string(content)
			}(),
			Want: helper.Issues{
				{
					Rule:    NewCommentsRule(),
					Message: "Avoid jammed comment ('#Jamm ...').",
					Range: hcl.Range{
						Filename: "comments_jammed.tf",
						Start:    hcl.Pos{Line: 8, Column: 1},
						End:      hcl.Pos{Line: 8, Column: 16},
					},
				},
				{
					Rule:    NewCommentsRule(),
					Message: "Avoid jammed comment ('##Jam ...').",
					Range: hcl.Range{
						Filename: "comments_jammed.tf",
						Start:    hcl.Pos{Line: 9, Column: 1},
						End:      hcl.Pos{Line: 9, Column: 17},
					},
				},
				{
					Rule:    NewCommentsRule(),
					Message: "Avoid jammed comment ('//Jam ...').",
					Range: hcl.Range{
						Filename: "comments_jammed.tf",
						Start:    hcl.Pos{Line: 10, Column: 1},
						End:      hcl.Pos{Line: 10, Column: 17},
					},
				},
				{
					Rule:    NewCommentsRule(),
					Message: "Avoid jammed comment ('///Ja ...').",
					Range: hcl.Range{
						Filename: "comments_jammed.tf",
						Start:    hcl.Pos{Line: 11, Column: 1},
						End:      hcl.Pos{Line: 11, Column: 18},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"comments_jammed.tf": tc.Content})
		rule := NewCommentsRule()
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
			if *commentsDeep {
				helper.AssertIssues(t, tc.Want, runner.Issues)
			} else {
				helper.AssertIssuesWithoutRange(t, tc.Want, runner.Issues)
			}
		})
	}
}
