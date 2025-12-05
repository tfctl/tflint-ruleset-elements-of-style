// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0
// no-cloc

package naming

import (
	"flag"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
)

var namingDeep = flag.Bool("namingDeep", false, "enable deep assert")

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
	cases := []struct {
		Name string
		Want namingRuleConfig
	}{
		{
			Name: "naming",
			Want: namingRuleConfig{},
		},
		{
			Name: "naming_noshout",
			Want: namingRuleConfig{
				Shout: testhelper.BoolPtr(false),
			},
		},
		{
			Name: "naming_disabled",
			Want: namingRuleConfig{
				Enabled: testhelper.BoolPtr(false),
			},
		},
		{
			Name: "naming_negative_length",
			Want: namingRuleConfig{
				Length: func() *int { v := -5; return &v }(),
			},
		},
		{
			Name: "naming_type_echo",
			Want: namingRuleConfig{
				TypeEcho: &typeEchoConfig{
					Enabled: testhelper.BoolPtr(true),
					Synonyms: map[string][]string{
						"aws_instance": {"vm", "box"},
					},
				},
			},
		},
		{
			Name: "naming_type_echo_disabled",
			Want: namingRuleConfig{
				TypeEcho: &typeEchoConfig{
					Enabled: testhelper.BoolPtr(false),
				},
			},
		},
		{
			Name: "naming_type_echo_custom",
			Want: namingRuleConfig{
				TypeEcho: &typeEchoConfig{
					Synonyms: map[string][]string{
						"aws_s3_bucket": {"pail"},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var got namingRuleConfig
			testhelper.LoadRuleConfig(t, tc.Name, &got)

			if diff := cmp.Diff(tc.Want, got); diff != "" {
				t.Errorf("config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
