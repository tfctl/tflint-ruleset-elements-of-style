// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package hungarian

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

var hungarianDeep = flag.Bool("hungarianDeep", false, "enable deep assert")

func boolPtr(b bool) *bool {
	return &b
}

func TestHungarian(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Rule", testHungarianRule)
	t.Run("Config", testHungarianConfig)
}

func testHungarianRule(t *testing.T) {
	var config hungarianRuleConfig
	testhelper.LoadRuleConfig(t, "hungarian", &config)

	cases := []struct {
		Name    string
		Content string
		Want    []string
	}{
		{
			Name: "hungarian_names",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/hungarian_test.tf")
				return string(content)
			}(),
			Want: []string{
				makeHungarianMessage("str_hung", "str"),
				makeHungarianMessage("hung_int", "int"),
				makeHungarianMessage("hung_bool_check", "bool"),
				makeHungarianMessage("map_hung", "map"),
				makeHungarianMessage("hung_lst", "lst"),
				makeHungarianMessage("hung_set_mod", "set"),
				makeHungarianMessage("num_hung", "num"),
				makeHungarianMessage("str_hung", "str"),
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"hungarian_test.tf": tc.Content})
		rule := NewHungarianRule()
		rule.Config = config

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
	}
}

func testHungarianConfig(t *testing.T) {
	cases := []struct {
		Name string
		Want hungarianRuleConfig
	}{
		{
			Name: "hungarian",
			Want: hungarianRuleConfig{
				Tags: []string{"str", "int", "num", "bool", "list", "lst", "set", "map", "arr", "array"},
			},
		},
		{
			Name: "hungarian_disabled",
			Want: hungarianRuleConfig{
				Enabled: boolPtr(false),
			},
		},
		{
			Name: "hungarian_custom_tags",
			Want: hungarianRuleConfig{
				Tags: []string{"foo", "bar"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var got hungarianRuleConfig
			testhelper.LoadRuleConfig(t, tc.Name, &got)

			if diff := cmp.Diff(tc.Want, got); diff != "" {
				t.Errorf("config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func makeHungarianMessage(name string, key string) string {
	return fmt.Sprintf("Avoid Hungarian notation '%s' in '%s'.", key, name)
}
