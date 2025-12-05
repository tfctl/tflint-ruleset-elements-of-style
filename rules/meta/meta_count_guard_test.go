// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	"flag"
	"os"
	"testing"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func testMetaCountGuardRule(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	var config metaRuleConfig
	testhelper.LoadRuleConfig(t, "meta", &config)
	// Disable source_version for this test
	sourceVersionDisabled := false
	config.SourceVersion = &sourceVersionDisabled

	cases := []struct {
		Name    string
		Content string
		Want    []string
	}{
		{
			Name: "count_guard",
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

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"meta.tf": tc.Content})
		rule := NewMetaRule()
		rule.Config = config

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
	}
}
