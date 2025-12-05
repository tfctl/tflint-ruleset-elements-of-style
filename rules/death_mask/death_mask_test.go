// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package death_mask

import (
	"flag"
	"os"
	"testing"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

var deathMaskDeep = flag.Bool("deathMaskDeep", false, "enable deep assert")

func TestDeathMask(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Rule", testDeathMaskRule)
}

func testDeathMaskRule(t *testing.T) {
	var config deathMaskRuleConfig
	testhelper.LoadRuleConfig(t, "death_mask", &config)

	cases := []struct {
		Name    string
		Content string
		Want    []string
	}{
		{
			Name: "death_mask",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/death_mask.tf")
				return string(content)
			}(),
			Want: []string{
				AvoidDeathMaskMessage,
				AvoidDeathMaskMessage,
				AvoidDeathMaskMessage,
				AvoidDeathMaskMessage,
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"death_mask.tf": tc.Content})
		rule := NewDeathMaskRule()
		rule.Config = config

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
	}
}
