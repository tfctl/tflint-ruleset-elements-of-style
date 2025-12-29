// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"os"
	"testing"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func testCommentsBlockRule(t *testing.T) {
	content, _ := os.ReadFile("./testdata/comments_block.tf")
	testContent := string(content)

	cases := []testhelper.RuleTestCase{
		{
			Name:    "eos_comments",
			Content: testContent,
			Want: []string{
				AvoidBlockCommentsMessage,
			},
		},
	}

	ruleFactory := func() tflint.Rule { return NewCommentsRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "comments_block.tf")
}
