// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package dry

import (
	"flag"
	"os"
	"testing"

	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/rulehelper"
	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestDry(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Config", testDryConfig)
	t.Run("Rule", testDryRule)
}

func testDryConfig(t *testing.T) {
	cases := []testhelper.ConfigTestCase{
		{
			Name: "eos_dry",
			Want: defaultDryConfig,
		},
		{
			Name: "eos_dry_disabled",
			Want: func() dryConfig {
				cfg := defaultDryConfig
				cfg.Enabled = rulehelper.BoolPtr(false)
				return cfg
			}(),
		},
		{
			Name: "eos_dry_threshold",
			Want: func() dryConfig {
				cfg := defaultDryConfig
				cfg.Threshold = 5
				return cfg
			}(),
		},
	}

	testhelper.ConfigTestRunner(t, defaultDryConfig, cases)

}

func testDryRule(t *testing.T) {
	cases := []testhelper.RuleTestCase{
		{
			Name: "eos_dry",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/dry_test.tf")
				return string(content)
			}(),
			Want: []string{
				`Avoid repeating value '"zakpxy"' 2 times.`,
				`Avoid repeating value '"1${local.literal1}"' 2 times.`,
				"Avoid repeating list 2 times.",
				"Avoid repeating map 2 times.",
				"Avoid repeating map 2 times.",
				"Avoid repeating list 2 times.",
				"Avoid repeating map 2 times.",
				"Duplicate block found 2 times.",
				"Duplicate block found 2 times.",
				"Duplicate block found 3 times.",
			},
		},
		{
			Name: "eos_dry_threshold",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/dry_test.tf")
				return string(content)
			}(),
			Want: []string{},
		},
	}

	ruleFactory := func() tflint.Rule { return NewDryRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "dry_test.tf")
}
