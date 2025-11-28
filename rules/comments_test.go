// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0
// no-cloc

package rules

import (
	"flag"
	"testing"

	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

var commentsDeep = flag.Bool("commentsDeep", false, "enable deep assert")

func TestCommentsRule(t *testing.T) {
	flag.Parse()

	var config commentsRuleConfig
	testhelper.LoadRuleConfig(t, "comments", &config)

	t.Run("block", func(t *testing.T) {
		cases := []struct {
			Name    string
			Content string
			Want    helper.Issues
		}{
			{
				Name: "block_comments",
				Content: func() string {
					content, _ := os.ReadFile("testdata/comments_block.tf")
					return string(content)
				}(),
				Want: helper.Issues{
					{
						Rule:    NewCommentsRule(),
						Message: "Avoid block comments.",
						Range: hcl.Range{
							Filename: "comments_block.tf",
							Start:    hcl.Pos{Line: 7, Column: 1},
							End:      hcl.Pos{Line: 10, Column: 3},
						},
					},
				},
			},
		}

		for _, tc := range cases {
			runner := helper.TestRunner(t, map[string]string{"comments_block.tf": tc.Content})
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
	})

	t.Run("eol", func(t *testing.T) {
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
	})

	t.Run("jammed", func(t *testing.T) {
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
	})

	t.Run("length", func(t *testing.T) {
		cases := []struct {
			Name    string
			Content string
			Want    helper.Issues
		}{
			{
				Name: "length_comments",
				Content: func() string {
					content, _ := os.ReadFile("testdata/comments_length.tf")
					return string(content)
				}(),
				Want: helper.Issues{
					{
						Rule:    NewCommentsRule(),
						Message: "Wrap comment at column 80 (currently 126).",
						Range: hcl.Range{
							Filename: "comments_length.tf",
							Start:    hcl.Pos{Line: 7, Column: 1},
							End:      hcl.Pos{Line: 7, Column: 114},
						},
					},
					{
						Rule:    NewCommentsRule(),
						Message: "Wrap comment at column 80 (currently 106).",
						Range: hcl.Range{
							Filename: "comments_length.tf",
							Start:    hcl.Pos{Line: 11, Column: 3},
							End:      hcl.Pos{Line: 11, Column: 108},
						},
					},
				},
			},
		}

		for _, tc := range cases {
			runner := helper.TestRunner(t, map[string]string{"comments_length.tf": tc.Content})
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
	})

	t.Run("threshold", func(t *testing.T) {
		threshold := 0.5
		cases := []struct {
			Name    string
			Content string
			Want    helper.Issues
		}{
			{
				Name: "threshold_violation",
				Content: func() string {
					content, _ := os.ReadFile("testdata/comments_threshold.tf")
					return string(content)
				}(),
				Want: helper.Issues{
					{
						Rule:    NewCommentsRule(),
						Message: "Comments ratio is 14 (threshold 50)",
						Range: hcl.Range{
							Filename: "comments_threshold.tf",
							Start:    hcl.Pos{Line: 1, Column: 1},
							End:      hcl.Pos{Line: 1, Column: 1},
						},
					},
				},
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
	})
}

func TestCommentsConfig(t *testing.T) {
	cases := []struct {
		Name string
		Want commentsRuleConfig
	}{
		{
			Name: "comments",
			Want: commentsRuleConfig{
				Block:     true,
				Column:    80,
				EOL:       true,
				Jammed:    &jammedConfig{Enabled: boolPtr(true), Tails: boolPtr(true)},
				Threshold: floatPtr(0.2),
				URLBypass: true,
			},
		},
		{
			Name: "comments_noblock",
			Want: commentsRuleConfig{
				Block: false,
			},
		},
		{
			Name: "comments_nojammed",
			Want: commentsRuleConfig{
				Jammed: &jammedConfig{Enabled: boolPtr(false)},
			},
		},
		{
			Name: "comments_nocolumn",
			Want: commentsRuleConfig{
				Column: 0,
			},
		},
		{
			Name: "comments_nourlbypass",
			Want: commentsRuleConfig{
				URLBypass: false,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var got commentsRuleConfig
			testhelper.LoadRuleConfig(t, tc.Name, &got)

			if diff := cmp.Diff(tc.Want, got); diff != "" {
				t.Errorf("config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func floatPtr(f float64) *float64 {
	return &f
}
