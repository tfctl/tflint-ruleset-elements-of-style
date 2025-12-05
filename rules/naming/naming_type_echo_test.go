// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package naming

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func testNamingTypeEchoRule(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	var config namingRuleConfig
	testhelper.LoadRuleConfig(t, "naming_type_echo", &config)

	cases := []struct {
		Name    string
		Content string
		Want    []string
	}{
		{
			Name: "echoed_names",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/naming_type_echo.tf")
				return string(content)
			}(),
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

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"naming_type_echo.tf": tc.Content})
		rule := NewNamingRule()
		rule.Config = config

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
	}
}

func makeTypeEchoMessage(typ string, name string) string {
	return fmt.Sprintf("Avoid echoing type \"%s\" in label \"%s\".", typ, name)
}
