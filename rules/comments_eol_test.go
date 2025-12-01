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

func testCommentsEOLRule(t *testing.T) {
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
			Name: "eol_comments",
			Content: func() string {
				content, _ := os.ReadFile("testdata/comments_eol.tf")
				return string(content)
			}(),
			Want: helper.Issues{
				{
					Rule:    NewCommentsRule(),
					Message: "Avoid EOL comments.",
					Range: hcl.Range{
						Filename: "comments_eol.tf",
						Start:    hcl.Pos{Line: 9, Column: 12},
						End:      hcl.Pos{Line: 9, Column: 37},
					},
				},
				{
					Rule:    NewCommentsRule(),
					Message: "Avoid EOL comments.",
					Range: hcl.Range{
						Filename: "comments_eol.tf",
						Start:    hcl.Pos{Line: 10, Column: 12},
						End:      hcl.Pos{Line: 10, Column: 38},
					},
				},
				{
					Rule:    NewCommentsRule(),
					Message: "Avoid EOL comments.",
					Range: hcl.Range{
						Filename: "comments_eol.tf",
						Start:    hcl.Pos{Line: 12, Column: 14},
						End:      hcl.Pos{Line: 12, Column: 39},
					},
				},
				// We need to mix in this block rule because both are emitted by
				// the EOL test in comments_eol.tf.
				{
					Rule:    NewCommentsRule(),
					Message: "Avoid block comments.",
					Range: hcl.Range{
						Filename: "comments_eol.tf",
						Start:    hcl.Pos{Line: 14, Column: 14},
						End:      hcl.Pos{Line: 14, Column: 43},
					},
				},
				{
					Rule:    NewCommentsRule(),
					Message: "Avoid EOL comments.",
					Range: hcl.Range{
						Filename: "comments_eol.tf",
						Start:    hcl.Pos{Line: 14, Column: 14},
						End:      hcl.Pos{Line: 14, Column: 43},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"comments_eol.tf": tc.Content})
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
