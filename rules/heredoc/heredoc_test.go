// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package heredoc

import (
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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
		Want    []string
	}{
		{
			Name: "heredoc",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/heredoc_test.tf")
				return string(content)
			}(),
			Want: []string{
				AvoidStandardHeredocMessage,
				AvoidEOFHeredocMessage,
				AvoidStandardHeredocMessage,
				AvoidEOFHeredocMessage,
				AvoidEOFHeredocMessage,
				AvoidEOFHeredocMessage,
				AvoidStandardHeredocMessage,
				AvoidEOFHeredocMessage,
				AvoidEOFHeredocMessage,
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

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
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
				Enabled: testhelper.BoolPtr(false),
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
