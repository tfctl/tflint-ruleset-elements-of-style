// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

var heredocDeep = flag.Bool("heredocDeep", false, "enable deep assert")

func TestHeredoc(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Rule", testHeredocRule)
	t.Run("Config", testHeredocConfig)
}

func testHeredocRule(t *testing.T) {

	var config heredocRuleConfig
	testhelper.LoadRuleConfig(t, "heredoc", &config)

	cases := []struct {
		Name    string
		Content string
		Want    helper.Issues
	}{
		{
			Name: "heredoc",
			Content: func() string {
				content, _ := os.ReadFile("testdata/heredoc_test.tf")
				return string(content)
			}(),
			Want: helper.Issues{
				{
					Rule:    NewHeredocRule(),
					Message: "Avoid standard heredoc (<<). Use indented (<<-) instead.",
					Range: hcl.Range{
						Filename: "heredoc.tf",
						Start:    hcl.Pos{Line: 8, Column: 11},
						End:      hcl.Pos{Line: 9, Column: 1},
					},
				},
				{
					Rule:    NewHeredocRule(),
					Message: "Avoid using 'EOF' as the heredoc delimiter.",
					Range: hcl.Range{
						Filename: "heredoc.tf",
						Start:    hcl.Pos{Line: 8, Column: 11},
						End:      hcl.Pos{Line: 9, Column: 1},
					},
				},
				{
					Rule:    NewHeredocRule(),
					Message: "Avoid standard heredoc (<<). Use indented (<<-) instead.",
					Range: hcl.Range{
						Filename: "heredoc.tf",
						Start:    hcl.Pos{Line: 14, Column: 13},
						End:      hcl.Pos{Line: 15, Column: 1},
					},
				},
				{
					Rule:    NewHeredocRule(),
					Message: "Avoid using 'EOF' as the heredoc delimiter.",
					Range: hcl.Range{
						Filename: "heredoc.tf",
						Start:    hcl.Pos{Line: 14, Column: 13},
						End:      hcl.Pos{Line: 15, Column: 1},
					},
				},
				{
					Rule:    NewHeredocRule(),
					Message: "Avoid using 'EOF' as the heredoc delimiter.",
					Range: hcl.Range{
						Filename: "heredoc.tf",
						Start:    hcl.Pos{Line: 21, Column: 38},
						End:      hcl.Pos{Line: 22, Column: 1},
					},
				},
				{
					Rule:    NewHeredocRule(),
					Message: "Avoid using 'EOF' as the heredoc delimiter.",
					Range: hcl.Range{
						Filename: "heredoc.tf",
						Start:    hcl.Pos{Line: 29, Column: 11},
						End:      hcl.Pos{Line: 30, Column: 1},
					},
				},
				{
					Rule:    NewHeredocRule(),
					Message: "Avoid standard heredoc (<<). Use indented (<<-) instead.",
					Range: hcl.Range{
						Filename: "heredoc.tf",
						Start:    hcl.Pos{Line: 35, Column: 12},
						End:      hcl.Pos{Line: 36, Column: 1},
					},
				},
				{
					Rule:    NewHeredocRule(),
					Message: "Avoid using 'EOF' as the heredoc delimiter.",
					Range: hcl.Range{
						Filename: "heredoc.tf",
						Start:    hcl.Pos{Line: 35, Column: 12},
						End:      hcl.Pos{Line: 36, Column: 1},
					},
				},
				{
					Rule:    NewHeredocRule(),
					Message: "Avoid using 'EOF' as the heredoc delimiter.",
					Range: hcl.Range{
						Filename: "heredoc.tf",
						Start:    hcl.Pos{Line: 41, Column: 11},
						End:      hcl.Pos{Line: 42, Column: 1},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"heredoc.tf": tc.Content})
		rule := NewHeredocRule()
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
			if *heredocDeep {
				helper.AssertIssues(t, tc.Want, runner.Issues)
			} else {
				helper.AssertIssuesWithoutRange(t, tc.Want, runner.Issues)
			}
		})
	}
}

func testHeredocConfig(t *testing.T) {
	cases := []struct {
		Name string
		Want heredocRuleConfig
	}{
		{
			Name: "heredoc",
			Want: heredocRuleConfig{
				EOF: true,
			},
		},
		{
			Name: "heredoc_disabled",
			Want: heredocRuleConfig{
				Enabled: boolPtr(false),
			},
		},
		{
			Name: "heredoc_no_eof",
			Want: heredocRuleConfig{
				EOF: false,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var got heredocRuleConfig
			testhelper.LoadRuleConfig(t, tc.Name, &got)

			if diff := cmp.Diff(tc.Want, got); diff != "" {
				t.Errorf("config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
