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

func testCommentsLengthRule(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	var config commentsRuleConfig
	testhelper.LoadRuleConfig(t, "comments", &config)

	cases := []struct {
		Name    string
		Content string
		Want    []string
	}{
		{
			Name: "length_comments",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/comments_length.tf")
				return string(content)
			}(),
			Want: []string{
				"Wrap comment at column 80 (currently 126).",
				"Wrap comment at column 80 (currently 106).",
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

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
	}
}
