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

var commentsDeep = flag.Bool("commentsDeep", false, "enable deep assert")

func TestComments(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Block", testCommentsBlockRule)
	t.Run("EOL", testCommentsEOLRule)
	t.Run("Jammed", testCommentsJammedRule)
	t.Run("Length", testCommentsLengthRule)
	t.Run("Threshold", testCommentsThresholdRule)
	t.Run("Config", testCommentsConfig)
}

func testCommentsConfig(t *testing.T) {
	cases := []struct {
		Name string
		Want commentsRuleConfig
	}{
		{
			Name: "comments",
			Want: commentsRuleConfig{
				Block:     true,
				EOL:       true,
				Jammed:    true,
				Length:    &lengthConfig{Column: 80, AllowURL: boolPtr(true)},
				Threshold: floatPtr(0.2),
			},
		},
		{
			Name: "comments_noblock",
			Want: commentsRuleConfig{
				Block: false,
			},
		},
		{
			Name: "comments_nojammed",
			Want: commentsRuleConfig{
				Jammed: false,
			},
		},
		{
			Name: "comments_nolength",
			Want: commentsRuleConfig{
				Length: nil,
			},
		},
		{
			Name: "comments_nocolumn",
			Want: commentsRuleConfig{
				Length: &lengthConfig{Column: 0, AllowURL: boolPtr(true)},
			},
		},
		{
			Name: "comments_nourl",
			Want: commentsRuleConfig{
				Length: &lengthConfig{Column: 80, AllowURL: boolPtr(false)},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var got commentsRuleConfig
			testhelper.LoadRuleConfig(t, tc.Name, &got)

			if diff := cmp.Diff(tc.Want, got); diff != "" {
				t.Errorf("config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func floatPtr(f float64) *float64 {
	return &f
}
