// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package hungarian

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestHungarian(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Config", testHungarianConfig)
	t.Run("Rule", testHungarianRule)
}

func testHungarianConfig(t *testing.T) {
	cases := []testhelper.ConfigTestCase{
		{
			Name: "eos_hungarian",
			Want: defaultHungarianConfig,
			// Tags: []string{"str", "int", "num", "bool", "list", "lst", "set", "map", "arr", "array"},
		},
		{
			Name: "eos_hungarian_disabled",
			Want: func() hungarianConfig {
				cfg := defaultHungarianConfig
				cfg.Enabled = rulehelper.BoolPtr(false)
				return cfg
			}(),
		},
		{
			Name: "eos_hungarian_custom_tags",
			Want: func() hungarianConfig {
				cfg := defaultHungarianConfig
				cfg.Tags = []string{"foo", "bar"}
				return cfg
			}(),
		},
	}

	testhelper.ConfigTestRunner(t, defaultHungarianConfig, cases)
}

func testHungarianRule(t *testing.T) {
	content, _ := os.ReadFile("./testdata/hungarian_test.tf")
	testContent := string(content)

	baseWant := []string{
		makeMessage("str_hung", "str"),
		makeMessage("hung_int", "int"),
		makeMessage("hung_bool_check", "bool"),
		makeMessage("map_hung", "map"),
		makeMessage("hung_lst", "lst"),
		makeMessage("hung_set_mod", "set"),
		makeMessage("num_hung", "num"),
		makeMessage("str_hung", "str"),
	}

	cases := []testhelper.RuleTestCase{
		{
			Name:    "eos_hungarian",
			Content: testContent,
			Want:    baseWant,
		},
		{
			Name:    "eos_hungarian_custom_tags",
			Content: testContent,
			Want: append(baseWant,
				makeMessage("foo_hung", "foo"),
				makeMessage("hung_bar_hung", "bar"),
			),
		},
	}

	ruleFactory := func() tflint.Rule { return NewHungarianRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "hungarian_test.tf")
}

func makeMessage(name string, key string) string {
	return fmt.Sprintf("Avoid Hungarian notation '%s' in '%s'.", key, name)
}
