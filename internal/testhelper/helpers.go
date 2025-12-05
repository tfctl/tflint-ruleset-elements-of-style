// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0
// no-cloc

package testhelper

import (
	"os"
	"sort"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// configFile is the path to the test configuration file.
const configFile = "testdata/.tflint_test.hcl"

// LoadRuleConfig loads the configuration for a specific rule from the test config file.
func LoadRuleConfig(t *testing.T, ruleName string, targetConfig interface{}) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %s", err)
	}

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, configFile)
	if diags.HasErrors() {
		t.Fatalf("Failed to parse config: %s", diags)
	}

	var config struct {
		Rules []struct {
			Name string   `hcl:"name,label"`
			Body hcl.Body `hcl:",remain"`
		} `hcl:"rule,block"`
	}

	if diags := gohcl.DecodeBody(file.Body, nil, &config); diags.HasErrors() {
		t.Fatalf("Failed to decode config: %s", diags)
	}

	found := false
	for _, r := range config.Rules {
		if r.Name == ruleName {
			if diags := gohcl.DecodeBody(r.Body, nil, targetConfig); diags.HasErrors() {
				t.Fatalf("Failed to decode %s config: %s", ruleName, diags)
			}
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("%s rule config not found", ruleName)
	}
}

type IssueSummary struct {
	Message string
	Line    int
}

func BoolPtr(b bool) *bool {
	return &b
}

func FloatPtr(f float64) *float64 {
	return &f
}

func AssertIssuesMessages(t *testing.T, expected []string, issues []*helper.Issue) {
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
