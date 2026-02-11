// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"os"
	"testing"

	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func testCommentsThresholdRule(t *testing.T) {
	thresholdContent, _ := os.ReadFile("./testdata/comments_threshold.tf")
	goodContent, _ := os.ReadFile("./testdata/comments_threshold_good.tf")

	cases := []testhelper.RuleTestCase{
		{
			Name:    "eos_comments_threshold_fail",
			Content: string(thresholdContent),
			Want: []string{
				"Comments ratio is 33 percent (minimum threshold 50 percent)",
			},
		},
		{
			Name:    "eos_comments_threshold_good",
			Content: string(goodContent),
			Want:    []string{},
		},
	}

	ruleFactory := func() tflint.Rule { return NewCommentsRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "comments_threshold.tf")
}
