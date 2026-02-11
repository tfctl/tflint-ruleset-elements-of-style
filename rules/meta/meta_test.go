// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	"flag"
	"os"
	"testing"

	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/rulehelper"
	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestMeta(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Config", testMetaConfig)

	t.Run("CountGuard", testMetaCountGuardRule)
	t.Run("Order", testMetaOrderRule)
	t.Run("SourceVersion", testMetaSourceVersionRule)
}

func testMetaConfig(t *testing.T) {
	cases := []testhelper.ConfigTestCase{
		{
			Name: "eos_meta",
			Want: defaultMetaConfig,
		},
		{
			Name: "eos_meta_disabled",
			Want: func() metaConfig {
				cfg := defaultMetaConfig
				cfg.Enabled = rulehelper.BoolPtr(false)
				return cfg
			}(),
		},
		{
			Name: "eos_meta_order",
			Want: func() metaConfig {
				cfg := defaultMetaConfig
				cfg.Enabled = rulehelper.BoolPtr(true)
				cfg.Order = []OrderConfig{{
					First: []string{"zal"},
					Last:  []string{"kpx"},
				}}
				return cfg
			}(),
		},
		{
			Name: "eos_meta_source_version_disabled",
			Want: func() metaConfig {
				cfg := defaultMetaConfig
				cfg.SourceVersion = rulehelper.BoolPtr(false)
				return cfg
			}(),
		},
	}

	testhelper.ConfigTestRunner(t, defaultMetaConfig, cases)
}

func testMetaCountGuardRule(t *testing.T) {
	cases := []testhelper.RuleTestCase{
		{
			Name: "eos_meta",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/meta_count_guard_test.tf")
				return string(content)
			}(),
			Want: []string{
				OnlyDynamicGuardMessage,
				GuardMustReturn10Message,
				GuardMustReturn10Message,
				OnlyDynamicGuardMessage,
				OnlyDynamicGuardMessage,
			},
		},
	}

	ruleFactory := func() tflint.Rule { return NewMetaRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "meta_count_guard_test.tf")
}

func testMetaOrderRule(t *testing.T) {
	cases := []testhelper.RuleTestCase{
		{
			Name: "eos_meta",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/meta_order_test.tf")
				return string(content)
			}(),
			Want: testhelper.MakeMessageList(MisOrderedMessage, 3),
		},
	}

	ruleFactory := func() tflint.Rule { return NewMetaRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "meta_order_test.tf")
}

func testMetaSourceVersionRule(t *testing.T) {
	cases := []testhelper.RuleTestCase{
		{
			Name: "eos_meta",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/meta_source_version_test.tf")
				return string(content)
			}(),
			Want: []string{
				"Git module source should specify ref parameter.",
				"Git module source should specify ref parameter.",
				"https module source should specify a valid archive extension.",
				"https module source should specify a valid archive extension.",
				"Mercurial module source should specify #revision.",
				"Module from registry should specify version.",
				"Module from registry should specify version.",
				"Pessimistic version constraint should specify at least major and minor version.",
				"Version constraint > or >= should not be used. Use ~> or exact version.",
				"Version constraint > or >= should not be used. Use ~> or exact version.",
				"Version constraint > or >= should not be used. Use ~> or exact version.",
			},
		},
	}

	ruleFactory := func() tflint.Rule { return NewMetaRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "meta_source_version_test.tf")
}
