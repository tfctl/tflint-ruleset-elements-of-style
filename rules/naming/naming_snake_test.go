// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package naming

import (
	"os"
	"testing"

	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func testNamingSnakeRule(t *testing.T) {
	content, _ := os.ReadFile("./testdata/naming_snake.tf")
	testContent := string(content)

	cases := []testhelper.RuleTestCase{
		{
			Name:    "eos_naming",
			Content: testContent,
			Want: []string{
				"Names should be snake_case (CamelCase).",
				"Names should be snake_case (kebab-case).",
				"Names should be snake_case (with.dots).",
				"Names should be snake_case (CamelCase).",
				"Names should be snake_case (CamelCase).",
				"Names should be snake_case (kebab-case).",
				"Names should be snake_case (with.dots).",
			},
		},
	}

	ruleFactory := func() tflint.Rule { return NewNamingRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "naming_snake.tf")
}
