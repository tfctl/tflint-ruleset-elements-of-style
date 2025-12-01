// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

var hungarianDeep = flag.Bool("hungarianDeep", false, "enable deep assert")

func TestHungarian(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Rule", testHungarianRule)
	t.Run("Config", testHungarianConfig)
}

func testHungarianRule(t *testing.T) {
	var config hungarianRuleConfig
	testhelper.LoadRuleConfig(t, "hungarian", &config)

	cases := []struct {
		Name    string
		Content string
		Want    helper.Issues
	}{
		{
			Name: "hungarian_names",
			Content: func() string {
				content, _ := os.ReadFile("testdata/hungarian_test.tf")
				return string(content)
			}(),
			Want: helper.Issues{
				{
					Rule:    NewHungarianRule(),
					Message: makeHungarianMessage("str_hung", "str"),
					Range: hcl.Range{
						Filename: "hungarian_test.tf",
						Start:    hcl.Pos{Line: 7, Column: 1},
						End:      hcl.Pos{Line: 7, Column: 20},
					},
				},
				{
					Rule:    NewHungarianRule(),
					Message: makeHungarianMessage("hung_int", "int"),
					Range: hcl.Range{
						Filename: "hungarian_test.tf",
						Start:    hcl.Pos{Line: 10, Column: 3},
						End:      hcl.Pos{Line: 10, Column: 15},
					},
				},
				{
					Rule:    NewHungarianRule(),
					Message: makeHungarianMessage("hung_bool_check", "bool"),
					Range: hcl.Range{
						Filename: "hungarian_test.tf",
						Start:    hcl.Pos{Line: 13, Column: 1},
						End:      hcl.Pos{Line: 13, Column: 24},
					},
				},
				{
					Rule:    NewHungarianRule(),
					Message: makeHungarianMessage("map_hung", "map"),
					Range: hcl.Range{
						Filename: "hungarian_test.tf",
						Start:    hcl.Pos{Line: 20, Column: 1},
						End:      hcl.Pos{Line: 20, Column: 38},
					},
				},
				{
					Rule:    NewHungarianRule(),
					Message: makeHungarianMessage("hung_lst", "lst"),
					Range: hcl.Range{
						Filename: "hungarian_test.tf",
						Start:    hcl.Pos{Line: 22, Column: 1},
						End:      hcl.Pos{Line: 22, Column: 39},
					},
				},
				{
					Rule:    NewHungarianRule(),
					Message: makeHungarianMessage("hung_set_mod", "set"),
					Range: hcl.Range{
						Filename: "hungarian_test.tf",
						Start:    hcl.Pos{Line: 26, Column: 1},
						End:      hcl.Pos{Line: 26, Column: 22},
					},
				},
				{
					Rule:    NewHungarianRule(),
					Message: makeHungarianMessage("num_hung", "num"),
					Range: hcl.Range{
						Filename: "hungarian_test.tf",
						Start:    hcl.Pos{Line: 30, Column: 1},
						End:      hcl.Pos{Line: 30, Column: 18},
					},
				},
				{
					Rule:    NewHungarianRule(),
					Message: makeHungarianMessage("str_hung", "str"),
					Range: hcl.Range{
						Filename: "hungarian_test.tf",
						Start:    hcl.Pos{Line: 35, Column: 1},
						End:      hcl.Pos{Line: 35, Column: 35},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"hungarian_test.tf": tc.Content})
		rule := NewHungarianRule()
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
			if *hungarianDeep {
				helper.AssertIssues(t, tc.Want, runner.Issues)
			} else {
				helper.AssertIssuesWithoutRange(t, tc.Want, runner.Issues)
			}
		})
	}
}

func testHungarianConfig(t *testing.T) {
	cases := []struct {
		Name string
		Want hungarianRuleConfig
	}{
		{
			Name: "hungarian",
			Want: hungarianRuleConfig{
				Tags: []string{"str", "int", "num", "bool", "list", "lst", "set", "map", "arr", "array"},
			},
		},
		{
			Name: "hungarian_disabled",
			Want: hungarianRuleConfig{
				Enabled: boolPtr(false),
			},
		},
		{
			Name: "hungarian_custom_tags",
			Want: hungarianRuleConfig{
				Tags: []string{"foo", "bar"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var got hungarianRuleConfig
			testhelper.LoadRuleConfig(t, tc.Name, &got)

			if diff := cmp.Diff(tc.Want, got); diff != "" {
				t.Errorf("config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func makeHungarianMessage(name string, key string) string {
	return fmt.Sprintf("Avoid Hungarian notation '%s' in '%s'.", key, name)
}
