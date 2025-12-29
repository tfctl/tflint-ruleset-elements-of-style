// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package naming

import (
	"os"
	"testing"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func testNamingShoutRule(t *testing.T) {
	content, _ := os.ReadFile("./testdata/naming_shout.tf")
	testContent := string(content)

	cases := []testhelper.RuleTestCase{
		{
			Name:    "eos_naming_nosnake",
			Content: testContent,
			Want: []string{
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT🤡)",
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT)",
				"Avoid SHOUTED names (SHOUT)",
			},
		},
	}

	ruleFactory := func() tflint.Rule { return NewNamingRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "naming_shout.tf")
}
