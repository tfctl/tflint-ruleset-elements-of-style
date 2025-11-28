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

var dryDeep = flag.Bool("dryDeep", false, "enable deep assert")

func TestDryRule(t *testing.T) {
	flag.Parse()

	var config dryRuleConfig
	testhelper.LoadRuleConfig(t, "dry", &config)

	cases := []struct {
		Name    string
		Content string
		Want    helper.Issues
	}{
		{
			Name: "dry",
			Content: func() string {
				content, _ := os.ReadFile("testdata/dry_test.tf")
				return string(content)
			}(),
			Want: helper.Issues{
				{
					Rule:    NewDryRule(),
					Message: `Avoid repeating value '"zakpxy"' 2 times.`,
					Range: hcl.Range{
						Filename: "dry.tf",
						Start:    hcl.Pos{Line: 10, Column: 14},
						End:      hcl.Pos{Line: 10, Column: 22},
					},
				},
				{
					Rule:    NewDryRule(),
					Message: `Avoid repeating value '"1${local.literal1}"' 2 times.`,
					Range: hcl.Range{
						Filename: "dry.tf",
						Start:    hcl.Pos{Line: 16, Column: 20},
						End:      hcl.Pos{Line: 16, Column: 40},
					},
				},
				{
					Rule:    NewDryRule(),
					Message: `Avoid repeating list 2 times.`,
					Range: hcl.Range{
						Filename: "dry.tf",
						Start:    hcl.Pos{Line: 23, Column: 11},
						End:      hcl.Pos{Line: 23, Column: 41},
					},
				},
				{
					Rule:    NewDryRule(),
					Message: `Avoid repeating map 2 times.`,
					Range: hcl.Range{
						Filename: "dry.tf",
						Start:    hcl.Pos{Line: 28, Column: 10},
						End:      hcl.Pos{Line: 31, Column: 4},
					},
				},
				{
					Rule:    NewDryRule(),
					Message: `Avoid repeating map 2 times.`,
					Range: hcl.Range{
						Filename: "dry.tf",
						Start:    hcl.Pos{Line: 39, Column: 11},
						End:      hcl.Pos{Line: 39, Column: 46},
					},
				},
				{
					Rule:    NewDryRule(),
					Message: `Avoid repeating list 2 times.`,
					Range: hcl.Range{
						Filename: "dry.tf",
						Start:    hcl.Pos{Line: 44, Column: 11},
						End:      hcl.Pos{Line: 44, Column: 39},
					},
				},
				{
					Rule:    NewDryRule(),
					Message: `Avoid repeating map 2 times.`,
					Range: hcl.Range{
						Filename: "dry.tf",
						Start:    hcl.Pos{Line: 51, Column: 14},
						End:      hcl.Pos{Line: 53, Column: 4},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"dry.tf": tc.Content})
		rule := NewDryRule()
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
			if *dryDeep {
				helper.AssertIssues(t, tc.Want, runner.Issues)
			} else {
				helper.AssertIssuesWithoutRange(t, tc.Want, runner.Issues)
			}
		})
	}
}

func TestDryConfig(t *testing.T) {
	cases := []struct {
		Name string
		Want dryRuleConfig
	}{
		{
			Name: "dry",
			Want: dryRuleConfig{
				Enabled: boolPtr(true),
			},
		},
		{
			Name: "dry_disabled",
			Want: dryRuleConfig{
				Enabled: boolPtr(false),
			},
		},
		{
			Name: "dry_info",
			Want: dryRuleConfig{
				Level: "info",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var got dryRuleConfig
			testhelper.LoadRuleConfig(t, tc.Name, &got)

			if diff := cmp.Diff(tc.Want, got); diff != "" {
				t.Errorf("config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
