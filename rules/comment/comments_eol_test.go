// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"os"
	"testing"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func testCommentsEOLRule(t *testing.T) {
	content, _ := os.ReadFile("./testdata/comments_eol.tf")
	testContent := string(content)

	cases := []testhelper.RuleTestCase{
		{
			Name:    "eos_comments",
			Content: testContent,
			Want:    testhelper.MakeMessageList(avoidBlockCommentsMessage, 1, avoidEOLCommentsMessage, 4),
		},
	}

	ruleFactory := func() tflint.Rule { return NewCommentsRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "comments_eol.tf")
}
