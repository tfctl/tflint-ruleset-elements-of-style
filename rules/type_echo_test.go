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

var typeEchoDeep = flag.Bool("typeEchoDeep", false, "enable deep assert")

func TestTypeEcho(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Rule", testTypeEchoRule)
	t.Run("Config", testTypeEchoConfig)
}

func testTypeEchoRule(t *testing.T) {
	var config typeEchoRuleConfig
	testhelper.LoadRuleConfig(t, "type_echo", &config)

	cases := []struct {
		Name    string
		Content string
		Want    helper.Issues
	}{
		{
			Name: "echoed_names",
			Content: func() string {
				content, _ := os.ReadFile("testdata/type_echo_test.tf")
				return string(content)
			}(),
			Want: helper.Issues{
				{
					Rule:    NewTypeEchoRule(),
					Message: makeTypeEchoMessage("variable", "variable_echo"),
					Range: hcl.Range{
						Filename: "type_echo_test.tf",
						Start:    hcl.Pos{Line: 7, Column: 1},
						End:      hcl.Pos{Line: 7, Column: 25},
					},
				},
				{
					Rule:    NewTypeEchoRule(),
					Message: makeTypeEchoMessage("local", "local_echo"),
					Range: hcl.Range{
						Filename: "type_echo_test.tf",
						Start:    hcl.Pos{Line: 10, Column: 3},
						End:      hcl.Pos{Line: 10, Column: 17},
					},
				},
				{
					Rule:    NewTypeEchoRule(),
					Message: makeTypeEchoMessage("check", "check_echo"),
					Range: hcl.Range{
						Filename: "type_echo_test.tf",
						Start:    hcl.Pos{Line: 13, Column: 1},
						End:      hcl.Pos{Line: 13, Column: 19},
					},
				},
				{
					Rule:    NewTypeEchoRule(),
					Message: makeTypeEchoMessage("aws_caller_identity", "caller_echo"),
					Range: hcl.Range{
						Filename: "type_echo_test.tf",
						Start:    hcl.Pos{Line: 20, Column: 1},
						End:      hcl.Pos{Line: 20, Column: 41},
					},
				},
				{
					Rule:    NewTypeEchoRule(),
					Message: makeTypeEchoMessage("random_password", "password_echo"),
					Range: hcl.Range{
						Filename: "type_echo_test.tf",
						Start:    hcl.Pos{Line: 22, Column: 1},
						End:      hcl.Pos{Line: 22, Column: 44},
					},
				},
				{
					Rule:    NewTypeEchoRule(),
					Message: makeTypeEchoMessage("module", "module_echo"),
					Range: hcl.Range{
						Filename: "type_echo_test.tf",
						Start:    hcl.Pos{Line: 26, Column: 1},
						End:      hcl.Pos{Line: 26, Column: 21},
					},
				},
				{
					Rule:    NewTypeEchoRule(),
					Message: makeTypeEchoMessage("output", "output_echo"),
					Range: hcl.Range{
						Filename: "type_echo_test.tf",
						Start:    hcl.Pos{Line: 30, Column: 1},
						End:      hcl.Pos{Line: 30, Column: 21},
					},
				},
				{
					Rule:    NewTypeEchoRule(),
					Message: makeTypeEchoMessage("aws_instance", "instance_echo"),
					Range: hcl.Range{
						Filename: "type_echo_test.tf",
						Start:    hcl.Pos{Line: 35, Column: 1},
						End:      hcl.Pos{Line: 35, Column: 40},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"type_echo_test.tf": tc.Content})
		rule := NewTypeEchoRule()
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
			if *typeEchoDeep {
				helper.AssertIssues(t, tc.Want, runner.Issues)
			} else {
				helper.AssertIssuesWithoutRange(t, tc.Want, runner.Issues)
			}
		})
	}
}

func testTypeEchoConfig(t *testing.T) {
	cases := []struct {
		Name string
		Want typeEchoRuleConfig
	}{
		{
			Name: "type_echo",
			Want: typeEchoRuleConfig{
				Synonyms: map[string][]string{
					"aws_instance": {"vm", "box"},
				},
			},
		},
		{
			Name: "type_echo_disabled",
			Want: typeEchoRuleConfig{
				Enabled: boolPtr(false),
			},
		},
		{
			Name: "type_echo_custom",
			Want: typeEchoRuleConfig{
				Synonyms: map[string][]string{
					"aws_s3_bucket": {"pail"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var got typeEchoRuleConfig
			testhelper.LoadRuleConfig(t, tc.Name, &got)

			if diff := cmp.Diff(tc.Want, got); diff != "" {
				t.Errorf("config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func makeTypeEchoMessage(typ string, name string) string {
	return fmt.Sprintf("Avoid echoing type \"%s\" in label \"%s\".", typ, name)
}
