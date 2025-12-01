// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0
// no-cloc

package rules

import (
	"flag"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
)

var namingDeep = flag.Bool("namingDeep", false, "enable deep assert")

func boolPtr(b bool) *bool {
	return &b
}

func TestNaming(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Length", testNamingLengthRule)
	t.Run("Shout", testNamingShoutRule)
	t.Run("Config", testNamingConfig)
}

func testNamingConfig(t *testing.T) {
	cases := []struct {
		Name string
		Want namingRuleConfig
	}{
		{
			Name: "naming",
			Want: namingRuleConfig{
				Length: &namingLengthConfig{
					Limit: 13,
				},
			},
		},
		{
			Name: "naming_noshout",
			Want: namingRuleConfig{
				Shout: &namingShoutConfig{
					Enabled: boolPtr(false),
				},
			},
		},
		{
			Name: "naming_nolength",
			Want: namingRuleConfig{
				Length: &namingLengthConfig{
					Enabled: boolPtr(false),
				},
			},
		},
		{
			Name: "naming_disabled",
			Want: namingRuleConfig{
				Enabled: boolPtr(false),
			},
		},
		{
			Name: "naming_negative_limit",
			Want: namingRuleConfig{
				Length: &namingLengthConfig{
					Limit: -5,
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
