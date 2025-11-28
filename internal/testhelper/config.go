// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0
// no-cloc

package testhelper

import (
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
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
