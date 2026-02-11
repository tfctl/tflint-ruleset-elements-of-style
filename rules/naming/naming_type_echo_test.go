// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package naming

import (
	"fmt"
	"os"
	"testing"

	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func testNamingTypeEchoRule(t *testing.T) {
	content, _ := os.ReadFile("./testdata/naming_type_echo.tf")
	testContent := string(content)

	cases := []testhelper.RuleTestCase{
		{
			Name:    "eos_naming_type_echo",
			Content: testContent,
			Want: []string{
				makeTypeEchoMessage("variable", "variable_echo"),
				makeTypeEchoMessage("local", "local_echo"),
				makeTypeEchoMessage("check", "check_echo"),
				makeTypeEchoMessage("aws_caller_identity", "caller_echo"),
				makeTypeEchoMessage("random_password", "password_echo"),
				makeTypeEchoMessage("module", "module_echo"),
				makeTypeEchoMessage("output", "output_echo"),
				makeTypeEchoMessage("aws_instance", "instance_echo"),
			},
		},
	}

	ruleFactory := func() tflint.Rule { return NewNamingRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "naming_type_echo.tf")
}

func makeTypeEchoMessage(typ string, name string) string {
	return fmt.Sprintf("Avoid echoing type \"%s\" in label \"%s\".", typ, name)
}
