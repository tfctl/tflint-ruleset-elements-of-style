// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package naming

import (
	"fmt"
	"os"
	"testing"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func testNamingLengthRule(t *testing.T) {
	content, _ := os.ReadFile("./testdata/naming_length.tf")
	testContent := string(content)

	limit := defaultLimit
	cases := []testhelper.RuleTestCase{
		{
			Name:    "eos_naming",
			Content: testContent,
			Want: []string{
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name' is 23).", limit),
				fmt.Sprintf("Avoid names longer than %d ('really_a_very_long_name_disabled' is 32).", limit),
			},
		},
	}

	ruleFactory := func() tflint.Rule { return NewNamingRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "naming_length.tf")
}
