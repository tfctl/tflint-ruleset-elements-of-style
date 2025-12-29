// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package death_mask

import (
	"flag"
	"os"
	"testing"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestDeathMask(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Config", testDeathMaskConfig)
	t.Run("Rule", testDeathMaskRule)
}

func testDeathMaskConfig(t *testing.T) {
	cases := []testhelper.ConfigTestCase{
		{
			Name: "eos_death_mask",
			Want: defaultDeathMaskConfig,
		},
		{
			Name: "eos_death_mask_disabled",
			Want: func() deathMaskConfig {
				cfg := defaultDeathMaskConfig
				cfg.Enabled = rulehelper.BoolPtr(false)
				return cfg
			}(),
		},
		{
			Name: "eos_death_mask_error",
			Want: func() deathMaskConfig {
				cfg := defaultDeathMaskConfig
				cfg.Enabled = rulehelper.BoolPtr(false)
				cfg.Level = "error"
				return cfg
			}(),
		},
	}

	testhelper.ConfigTestRunner(t, defaultDeathMaskConfig, cases)
}

func testDeathMaskRule(t *testing.T) {
	cases := []testhelper.RuleTestCase{
		{
			Name: "eos_death_mask",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/death_mask.tf")
				return string(content)
			}(),
			Want: FillWantMessages(4, AvoidDeathMaskMessage),
		},
		{
			Name: "eos_death_mask_disabled",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/death_mask.tf")
				return string(content)
			}(),
			Want: []string{},
		},
	}

	ruleFactory := func() tflint.Rule { return NewDeathMaskRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "death_mask.tf")
}

func FillWantMessages(count int, message string) []string {
	want := make([]string, count)
	for i := range count {
		want[i] = message
	}
	return want
}
