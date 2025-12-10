// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package reminder

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func testReminderConfig(t *testing.T) {
	cases := []struct {
		Name string
		Want reminderRuleConfig
	}{
		{
			Name: "reminder",
			Want: reminderRuleConfig{
				Tags: []string{"TODO"},
			},
		},
		{
			Name: "reminder_many_tags",
			Want: reminderRuleConfig{
				Tags: []string{"BUG", "FIXME", "HORROR", "TODO"},
			},
		},
		{
			Name: "reminder_disabled",
			Want: reminderRuleConfig{
				Enabled: testhelper.BoolPtr(false),
				Tags:    []string{},
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

func testReminderRule(t *testing.T) {
	var config reminderRuleConfig
	testhelper.LoadRuleConfig(t, "reminder_many_tags", &config)

	cases := []struct {
		Name    string
		Content string
		Want    []string
	}{
		{
			Name: "reminders",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/reminder_test.tf")
				return string(content)
			}(),
			Want: []string{
				makeReminderMessage("// TODO Reminder found."),
				makeReminderMessage("# TODO Reminder found."),
				makeReminderMessage("# TODO Reminder found."),
				makeReminderMessage("# HORROR This is so bad."),
				makeReminderMessage("# FIXME Reminder found."),
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

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
	}
}

func makeReminderMessage(line string) string {
	return fmt.Sprintf("Resolve reminder: '%s'.", line)
}
