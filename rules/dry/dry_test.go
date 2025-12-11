// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package dry

import (
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestDry(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Config", testDryConfig)
	t.Run("Rule", testDryRule)
}

func testDryConfig(t *testing.T) {
	cases := []struct {
		Name string
		Want dryRuleConfig
	}{
		{
			Name: "dry",
			Want: dryRuleConfig{
				Enabled: testhelper.BoolPtr(true),
			},
		},
		{
			Name: "dry_disabled",
			Want: dryRuleConfig{
				Enabled: testhelper.BoolPtr(false),
			},
		},
		{
			Name: "dry_info",
			Want: dryRuleConfig{
				Level: "info",
			},
		},
		{
			Name: "dry_threshold",
			Want: dryRuleConfig{
				Threshold: 5,
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

func testDryRule(t *testing.T) {
	cases := []struct {
		Name       string
		ConfigName string
		Content    string
		Want       []string
	}{
		{
			Name: "dry",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/dry_test.tf")
				return string(content)
			}(),
			Want: []string{
				`Avoid repeating value '"zakpxy"' 2 times.`,
				`Avoid repeating value '"1${local.literal1}"' 2 times.`,
				"Avoid repeating list 2 times.",
				"Avoid repeating map 2 times.",
				"Avoid repeating map 2 times.",
				"Avoid repeating list 2 times.",
				"Avoid repeating map 2 times.",
				"Duplicate block found 2 times.",
				"Duplicate block found 2 times.",
				"Duplicate block found 3 times.",
			},
		},
		{
			Name:       "dry_threshold",
			ConfigName: "dry_threshold",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/dry_test.tf")
				return string(content)
			}(),
			Want: []string{},
		},
	}

	for _, tc := range cases {
		var config dryRuleConfig
		configName := "dry"
		if tc.ConfigName != "" {
			configName = tc.ConfigName
		}
		testhelper.LoadRuleConfig(t, configName, &config)

		runner := helper.TestRunner(t, map[string]string{"dry_test.tf": tc.Content})
		rule := NewDryRule()
		rule.Config = config

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
	}
}
