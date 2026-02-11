// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package reminder

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/rulehelper"
	"github.com/tfctl/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestReminder(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("Config", testReminderConfig)
	t.Run("Rule", testReminderRule)
}

func testReminderConfig(t *testing.T) {
	cases := []testhelper.ConfigTestCase{
		{
			Name: "eos_reminder",
			Want: defaultReminderConfig,
		},
		{
			Name: "eos_reminder_disabled",
			Want: func() reminderConfig {
				cfg := defaultReminderConfig
				cfg.Enabled = rulehelper.BoolPtr(false)
				return cfg
			}(),
		},
		{
			Name: "eos_reminder_extras",
			Want: func() reminderConfig {
				cfg := defaultReminderConfig
				cfg.Extras = []string{"NOTGOOD", "REALBAD"}
				return cfg
			}(),
		},
		{
			Name: "eos_reminder_many_tags",
			Want: func() reminderConfig {
				cfg := defaultReminderConfig
				cfg.Tags = []string{"BUG", "FIXME", "HORROR", "TODO"}
				return cfg
			}(),
		},
	}

	testhelper.ConfigTestRunner(t, defaultReminderConfig, cases)
}

func testReminderRule(t *testing.T) {
	content, _ := os.ReadFile("./testdata/reminder_test.tf")
	testContent := string(content)

	// eos_reminder_many_tags has tags = ["BUG", "FIXME", "HORROR", "TODO"]
	manyTagsWant := []string{
		makeReminderMessage("// TODO Reminder found."),
		makeReminderMessage("# TODO Reminder found."),
		makeReminderMessage("# TODO Reminder found."),
		makeReminderMessage("# HORROR This is so bad."),
		makeReminderMessage("# FIXME Reminder found."),
	}

	// eos_reminder_extras has default tags (BUG, TODO) plus extras (NOTGOOD, REALBAD)
	extrasWant := []string{
		makeReminderMessage("// TODO Reminder found."),
		makeReminderMessage("# TODO Reminder found."),
		makeReminderMessage("# TODO Reminder found."),
		makeReminderMessage("// NOTGOOD This is really not a good idea."),
		makeReminderMessage("// REALBAD This is a really bad idea."),
	}

	cases := []testhelper.RuleTestCase{
		{
			Name:    "eos_reminder_many_tags",
			Content: testContent,
			Want:    manyTagsWant,
		},
		{
			Name:    "eos_reminder_extras",
			Content: testContent,
			Want:    extrasWant,
		},
	}

	ruleFactory := func() tflint.Rule { return NewReminderRule() }
	testhelper.RuleTestRunner(t, ruleFactory, "testdata/.tflint_test.hcl", cases, "reminder_test.tf")
}

func makeReminderMessage(line string) string {
	return fmt.Sprintf("Resolve reminder: '%s'.", line)
}
