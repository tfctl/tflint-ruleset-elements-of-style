// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"os"
	"testing"

	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func testCommentsLengthRule(t *testing.T) {
	content, _ := os.ReadFile("./testdata/comments_length.tf")
	testContent := string(content)

	cases := []testhelper.RuleTestCase{
		{
			Name:    "eos_comments",
			Content: testContent,
			Want: []string{
				"Wrap comment at column 80 (currently 126).",
				"Wrap comment at column 80 (currently 106).",
			},
		},
	}

	ruleFactory := func() tflint.Rule { return NewCommentsRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "comments_length.tf")
}
