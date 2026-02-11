// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0
// no-cloc

package naming

import (
	"flag"
	"testing"

	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/rulehelper"
	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/testhelper"
)

func TestNaming(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Length", testNamingLengthRule)
	t.Run("Shout", testNamingShoutRule)
	t.Run("Snake", testNamingSnakeRule)
	t.Run("TypeEcho", testNamingTypeEchoRule)
	t.Run("Config", testNamingConfig)
}

func testNamingConfig(t *testing.T) {
	cases := []testhelper.ConfigTestCase{
		{
			Name: "eos_naming",
			Want: namingRuleConfig{},
		},
		{
			Name: "eos_naming_noshout",
			Want: namingRuleConfig{
				Shout: rulehelper.BoolPtr(false),
			},
		},
		{
			Name: "eos_naming_disabled",
			Want: namingRuleConfig{
				Enabled: rulehelper.BoolPtr(false),
			},
		},
		{
			Name: "eos_naming_negative_length",
			Want: namingRuleConfig{
				Length: func() *int { v := -5; return &v }(),
			},
		},
		{
			Name: "eos_naming_type_echo",
			Want: namingRuleConfig{
				TypeEcho: &typeEchoConfig{
					Enabled: rulehelper.BoolPtr(true),
					Synonyms: map[string][]string{
						"aws_instance": {"vm", "box"},
					},
				},
			},
		},
		{
			Name: "eos_naming_type_echo_disabled",
			Want: namingRuleConfig{
				TypeEcho: &typeEchoConfig{
					Enabled: rulehelper.BoolPtr(false),
				},
			},
		},
		{
			Name: "eos_naming_type_echo_custom",
			Want: namingRuleConfig{
				TypeEcho: &typeEchoConfig{
					Synonyms: map[string][]string{
						"aws_s3_bucket": {"pail"},
					},
				},
			},
		},
	}

	testhelper.ConfigTestRunner(t, namingRuleConfig{}, cases)
}
