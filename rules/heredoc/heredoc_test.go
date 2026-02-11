// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package heredoc

import (
	"flag"
	"os"
	"testing"

	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/rulehelper"
	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestHeredoc(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Config", testHeredocConfig)
	t.Run("Rule", testHeredocRule)
}

func testHeredocConfig(t *testing.T) {
	cases := []testhelper.ConfigTestCase{
		{
			Name: "eos_heredoc",
			Want: defaultHeredocConfig,
		},
		{
			Name: "eos_heredoc_disabled",
			Want: func() heredocConfig {
				cfg := defaultHeredocConfig
				cfg.Enabled = rulehelper.BoolPtr(false)
				return cfg
			}(),
		},
		{
			Name: "eos_heredoc_no_eof",
			Want: func() heredocConfig {
				cfg := defaultHeredocConfig
				cfg.EOF = false
				return cfg
			}(),
		},
	}

	testhelper.ConfigTestRunner(t, defaultHeredocConfig, cases)
}

func testHeredocRule(t *testing.T) {
	cases := []testhelper.RuleTestCase{
		{
			Name: "eos_heredoc",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/heredoc_test.tf")
				return string(content)
			}(),
			Want: testhelper.MakeMessageList(
				AvoidEOFHeredocMessage, 6,
				AvoidStandardHeredocMessage, 3,
			),
		},
	}

	ruleFactory := func() tflint.Rule { return NewHeredocRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "heredoc_test.tf")
}
