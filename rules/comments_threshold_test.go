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

func testCommentsThresholdRule(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	var config commentsRuleConfig
	testhelper.LoadRuleConfig(t, "comments", &config)

	threshold := 0.5
	cases := []struct {
		Name    string
		Content string
		Want    helper.Issues
	}{
		{
			Name: "threshold_fail",
			Content: func() string {
				content, _ := os.ReadFile("testdata/comments_threshold.tf")
				return string(content)
			}(),
			Want: helper.Issues{
				{
					Rule:    NewCommentsRule(),
					Message: "Comments ratio is 14 percent (minimum threshold 50 percent)",
					Range: hcl.Range{
						Filename: "comments_threshold.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
		{
			Name: "threshold_good",
			Content: func() string {
				content, _ := os.ReadFile("testdata/comments_threshold_good.tf")
				return string(content)
			}(),
			Want: helper.Issues{},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"comments_threshold.tf": tc.Content})
		rule := NewCommentsRule()
		rule.Config = config
		rule.Config.Threshold = &threshold

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
