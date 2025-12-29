// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rulehelper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// FindConfigFile searches for .tflint.hcl in CWD first, then $HOME. Returns the
// path if found, or an error describing where it looked.
func FindConfigFile() (string, error) {
	var searched []string

	// Check CWD.
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get CWD: %w", err)
	}
	cwdConfig := filepath.Join(cwd, ".tflint.hcl")
	searched = append(searched, cwdConfig)
	if _, err := os.Stat(cwdConfig); err == nil {
		return cwdConfig, nil
	}

	// Check HOME.
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}
	homeConfig := filepath.Join(home, ".tflint.hcl")
	searched = append(searched, homeConfig)
	if _, err := os.Stat(homeConfig); err == nil {
		return homeConfig, nil
	}

	return "", fmt.Errorf("no .tflint.hcl found (searched: %v)", searched)
}

// LoadRuleConfig loads the configuration for a specific rule from a config
// file. If configFile is empty, it searches CWD then $HOME for .tflint.hcl.
func LoadRuleConfig(ruleName string, targetConfig interface{}, configFile string) error {
	// Determine config file path.
	configPath := configFile
	if configPath == "" {
		var err error
		configPath, err = FindConfigFile()
		if err != nil {
			return err
		}
	}

	// Read and parse the config file.
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, configPath)
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse config %s: %s", configPath, diags)
	}

	// Decode the top-level structure.
	var config struct {
		Plugins []struct {
			Name string   `hcl:"name,label"`
			Body hcl.Body `hcl:",remain"`
		} `hcl:"plugin,block"`
		Rules []struct {
			Name string   `hcl:"name,label"`
			Body hcl.Body `hcl:",remain"`
		} `hcl:"rule,block"`
	}

	if diags := gohcl.DecodeBody(file.Body, nil, &config); diags.HasErrors() {
		return fmt.Errorf("failed to decode config %s: %s", configPath, diags)
	}

	// Find and decode the specific rule. If not found, the caller's config
	// struct retains its default values.
	for _, r := range config.Rules {
		if r.Name == ruleName {
			if diags := gohcl.DecodeBody(r.Body, nil, targetConfig); diags.HasErrors() {
				return fmt.Errorf("failed to decode rule %s in %s: %s", ruleName, configPath, diags)
			}
			break
		}
	}

	return nil
}
