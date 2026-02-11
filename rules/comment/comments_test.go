// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0
// no-cloc

package comment

import (
	"flag"
	"testing"

	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/rulehelper"
	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/testhelper"
)

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

func floatPtr(f float64) *float64 {
	return &f
}

func testCommentsConfig(t *testing.T) {
	cases := []testhelper.ConfigTestCase{
		{
			Name: "eos_comments",
			Want: commentsRuleConfig{
				Block:     true,
				EOL:       true,
				Jammed:    true,
				Length:    &lengthConfig{Column: 80, AllowURL: rulehelper.BoolPtr(true)},
				Threshold: floatPtr(0.2),
			},
		},
		{
			Name: "eos_comments_noblock",
			Want: commentsRuleConfig{
				Block: false,
			},
		},
		{
			Name: "eos_comments_nojammed",
			Want: commentsRuleConfig{
				Jammed: false,
			},
		},
		{
			Name: "eos_comments_nolength",
			Want: commentsRuleConfig{
				Length: nil,
			},
		},
		{
			Name: "eos_comments_nocolumn",
			Want: commentsRuleConfig{
				Length: &lengthConfig{Column: 0, AllowURL: rulehelper.BoolPtr(true)},
			},
		},
		{
			Name: "eos_comments_nourl",
			Want: commentsRuleConfig{
				Length: &lengthConfig{Column: 80, AllowURL: rulehelper.BoolPtr(false)},
			},
		},
	}

	testhelper.ConfigTestRunner(t, commentsRuleConfig{}, cases)
}
