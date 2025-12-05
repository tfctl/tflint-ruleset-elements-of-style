// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"flag"
	"os"
	"testing"

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
		Want    []string
	}{
		{
			Name: "threshold_fail",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/comments_threshold.tf")
				return string(content)
			}(),
			Want: []string{
				"Comments ratio is 14 percent (minimum threshold 50 percent)",
			},
		},
		{
			Name: "threshold_good",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/comments_threshold_good.tf")
				return string(content)
			}(),
			Want: []string{},
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

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
	}
}
