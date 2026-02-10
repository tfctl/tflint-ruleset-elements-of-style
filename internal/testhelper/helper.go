// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0
// no-cloc

package testhelper

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/rulehelper"
)

// TestConfigFile is the standard path to the test config file relative to each
// rule's test directory.
const TestConfigFile = "testdata/.tflint_test.hcl"

// ConfigTestCase represents a test case for parsing rule configurations from
// TFLint HCL config files.
type ConfigTestCase struct {
	Name string
	// The want is generic so that we can use it for any rule config struct.
	Want any
}

// IssueSummary represents a summary of an issue's message and line number.
type IssueSummary struct {
	Message string
	Line    int
}

// RuleTestCase represents a test case for rule checking. For our rule tests,
// we're only going to check the issue message returned as opposed to doing a
// "deep" want vs got comparison as we're trusting the TFLint SDK to get the
// source range/etc identification correct.
type RuleTestCase struct {
	Name    string
	Content string
	Want    []string
}

// assertRuleIssueMessages tests that the issues collected by the rule test
// match in both length and values.
func assertRuleIssueMessages(t *testing.T, expected []string, issues []*helper.Issue) {
	// Quick check to make sure the number of issues we found matches the number
	// expected.
	if len(issues) != len(expected) {
		t.Logf("Expected %d issues, got %d", len(expected), len(issues))
		for i, issue := range issues {
			t.Logf("Issue %d: %s at %s", i, issue.Message, issue.Range)
		}
		t.Fatalf("Number of issues mismatch: got %d, want %d", len(issues), len(expected))
	}

	// Build a slice of messages and starting line numbers from the issues.
	var actualWithLines []IssueSummary
	for _, issue := range issues {
		actualWithLines = append(actualWithLines, IssueSummary{Message: issue.Message, Line: issue.Range.Start.Line})
	}

	// Sort both actual and expected so we can iterate and compare.
	sort.Slice(actualWithLines, func(i, j int) bool {
		return actualWithLines[i].Message < actualWithLines[j].Message
	})

	expectedMessages := make([]string, len(expected))
	copy(expectedMessages, expected)
	sort.Strings(expectedMessages)

	// Iterate over each and make sure their message matches.
	for i := range actualWithLines {
		if actualWithLines[i].Message != expectedMessages[i] {
			t.Errorf("Message mismatch at index %d, line %d:\nexpected: %s,\ngot:      %s", i, actualWithLines[i].Line, expectedMessages[i], actualWithLines[i].Message)
		}
	}
}

// ConfigTestRunner runs each of the config test cases.
func ConfigTestRunner[D any](t *testing.T, defaultConfig D, cases []ConfigTestCase) {
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {

			// We have to deep copy here so that the got doesn't clobber the default
			// config and, therefore, screw up subsequent tests. This is not needed
			// for the mainline as we'll only be loading the config once and, thus,
			// nothing to clobber down the road.
			got, err := deepCopy(defaultConfig)
			if err != nil {
				t.Fatalf("DeepCopy failed: %v", err)
			}

			if err := rulehelper.LoadRuleConfig(tc.Name, &got, TestConfigFile); err != nil {
				t.Fatalf("LoadRuleConfig failed: %v", err)
			}

			if diff := cmp.Diff(tc.Want, got); diff != "" {
				t.Errorf("config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// deepCopy performs a deep copy of a generic type T using JSON. This is needed
// in the tests so that the parsed config (the got) doesn't clobber the want.
func deepCopy[T any](src T) (T, error) {
	var dst T
	b, err := json.Marshal(src)
	if err != nil {
		return dst, err
	}
	err = json.Unmarshal(b, &dst)
	return dst, err
}

// MakeMessageList creates a slice of repeated messages. Arguments are provided
// in pairs: message string followed by count. For example:
//
//	MakeMessageList(MsgX, 2)                      // [MsgX, MsgX]
//	MakeMessageList(MsgX, 3, MsgY, 1, MsgZ, 2)    // [MsgX, MsgX, MsgX, MsgY, MsgZ, MsgZ]
func MakeMessageList(args ...any) []string {
	var messages []string
	for i := 0; i < len(args)-1; i += 2 {
		message, ok := args[i].(string)
		if !ok {
			continue
		}
		count, ok := args[i+1].(int)
		if !ok {
			continue
		}
		for j := 0; j < count; j++ {
			messages = append(messages, message)
		}
	}
	return messages
}

// RuleTestRunner runs each of the rule test cases. The rules are collectively
// asserted as opposed to individually. This allows a much less porcelain
// definition of test cases. If we were to evaluate them individually or deeply
// the connectivity between the test case definitions and source TF test files
// would be much more brittle.
func RuleTestRunner(t *testing.T, ruleFactory func() tflint.Rule, configFile string, cases []RuleTestCase, sourceFilename string) {
	for _, cv := range cases {
		// Capture loop variable for closure, thus dealing with parallel tests.
		c := cv
		t.Run(c.Name, func(t *testing.T) {
			rule := ruleFactory()

			// Use reflection to set RuleName and ConfigFile on the rule struct.
			// This tells Check() which config to load and from where.
			rVal := reflect.ValueOf(rule)
			if rVal.Kind() == reflect.Ptr {
				rVal = rVal.Elem()
			}

			ruleNameField := rVal.FieldByName("RuleName")
			if ruleNameField.IsValid() && ruleNameField.CanSet() {
				ruleNameField.SetString(c.Name)
			}

			configFileField := rVal.FieldByName("ConfigFile")
			if configFileField.IsValid() && configFileField.CanSet() {
				configFileField.SetString(configFile)
			}

			runner := helper.TestRunner(t, map[string]string{sourceFilename: c.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			assertRuleIssueMessages(t, c.Want, runner.Issues)
		})
	}
}
