// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2"
	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

var reminderDeep = flag.Bool("reminderDeep", false, "enable deep assert")

func TestReminder(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Rule", testReminderRule)
	t.Run("Config", testReminderConfig)
}

func testReminderRule(t *testing.T) {
	var config reminderRuleConfig
	testhelper.LoadRuleConfig(t, "reminder", &config)

	cases := []struct {
		Name    string
		Content string
		Want    helper.Issues
	}{
		{
			Name: "reminders",
			Content: func() string {
				content, _ := os.ReadFile("testdata/reminder_test.tf")
				return string(content)
			}(),
			Want: helper.Issues{
				{
					Rule:    NewReminderRule(),
					Message: makeReminderMessage("// TODO Reminder found."),
					Range: hcl.Range{
						Filename: "reminder_test.tf",
						Start:    hcl.Pos{Line: 4, Column: 1},
						End:      hcl.Pos{Line: 4, Column: 23},
					},
				},
				{
					Rule:    NewReminderRule(),
					Message: makeReminderMessage("# TODO Reminder found."),
					Range: hcl.Range{
						Filename: "reminder_test.tf",
						Start:    hcl.Pos{Line: 5, Column: 1},
						End:      hcl.Pos{Line: 5, Column: 22},
					},
				},
				{
					Rule:    NewReminderRule(),
					Message: makeReminderMessage("# TODO Reminder found."),
					Range: hcl.Range{
						Filename: "reminder_test.tf",
						Start:    hcl.Pos{Line: 8, Column: 11},
						End:      hcl.Pos{Line: 8, Column: 32},
					},
				},
				{
					Rule:    NewReminderRule(),
					Message: makeReminderMessage("# FIXME Reminder found."),
					Range: hcl.Range{
						Filename: "reminder_test.tf",
						Start:    hcl.Pos{Line: 14, Column: 37},
						End:      hcl.Pos{Line: 14, Column: 59},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"reminder_test.tf": tc.Content})
		rule := NewReminderRule()
		rule.Config = config

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		if len(runner.Issues) != len(tc.Want) {
			t.Logf("Expected %d issues, got %d", len(tc.Want), len(runner.Issues))
			for i, issue := range runner.Issues {
				t.Logf("Issue %d: %s at %s", i, issue.Message, issue.Range)
			}
			t.Fatalf("Number of issues mismatch: got %d, want %d", len(runner.Issues), len(tc.Want))
		}

		t.Run(tc.Name, func(t *testing.T) {
			if *reminderDeep {
				helper.AssertIssues(t, tc.Want, runner.Issues)
			} else {
				helper.AssertIssuesWithoutRange(t, tc.Want, runner.Issues)
			}
		})
	}
}

func testReminderConfig(t *testing.T) {
	cases := []struct {
		Name string
		Want reminderRuleConfig
	}{
		{
			Name: "reminder",
			Want: reminderRuleConfig{
				Tags: []string{"TODO", "FIXME"},
			},
		},
		{
			Name: "reminder_disabled",
			Want: reminderRuleConfig{
				Enabled: boolPtr(false),
			},
		},
		{
			Name: "reminder_custom_tags",
			Want: reminderRuleConfig{
				Tags: []string{"HACK", "BUG"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var got reminderRuleConfig
			testhelper.LoadRuleConfig(t, tc.Name, &got)

			if diff := cmp.Diff(tc.Want, got); diff != "" {
				t.Errorf("config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func makeReminderMessage(line string) string {
	return fmt.Sprintf("Resolve reminder: '%s'.", line)
}
