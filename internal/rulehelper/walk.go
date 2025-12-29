// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rulehelper

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// BlockDef represents a block definition for schema generation.
type BlockDef struct {
	Typ     string
	Labels  []string
	Synonym string
}

// AllLintableBlocks defines all block types and their label structures to
// check. The order is not important here, but we try to sync up with the order
// found in the *_test.tf sources.
var AllLintableBlocks = []BlockDef{
	{Typ: "variable", Labels: []string{"name"}},
	// {Typ: "locals", Labels: []string{}},
	{Typ: "check", Labels: []string{"name"}},
	{Typ: "data", Labels: []string{"type", "name"}},
	{Typ: "ephemeral", Labels: []string{"type", "name"}},
	{Typ: "module", Labels: []string{"name"}},
	{Typ: "output", Labels: []string{"name"}},
	{Typ: "resource", Labels: []string{"type", "name"}},
}

// normalizeBlock extracts the type, name, and synonym from a block. For
// two-label blocks (e.g., resource, data), the first label is the type and
// the second is the name. For single-label blocks (e.g., variable, output),
// the block type itself serves as the type and the label is the name.
func normalizeBlock(block *hclsyntax.Block, myBlocks []BlockDef) (string, string, string) {
	var name string
	var typ string

	switch len(block.Labels) {
	case 2:
		typ = block.Labels[0]
		name = block.Labels[1]
	case 1:
		typ = block.Type
		name = block.Labels[0]
	default:
		typ = block.Type
		name = ""
	}

	synonym := ""
	for _, def := range myBlocks {
		if def.Typ == typ && def.Synonym != "" {
			synonym = def.Synonym
			break
		}
	}
	return typ, name, synonym
}

// WalkBlocks iterates over blocks and locals using AST and applies the check
// function.
func WalkBlocks[T any](
	runner tflint.Runner,
	myBlocks []BlockDef,
	rule T,
	checkFuncs ...func(tflint.Runner, T, hcl.Range, string, string, string),
) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		for _, block := range body.Blocks {
			// Handle locals specifically.
			if block.Type == "locals" {
				for _, attr := range block.Body.Attributes {
					for _, checkFunc := range checkFuncs {
						checkFunc(runner, rule, attr.Range(), "local", attr.Name, "")
					}
				}
				continue
			}

			// Handle other blocks.
			typ, name, synonym := normalizeBlock(block, myBlocks)

			// Filter by myBlocks to ensure we only lint what we expect.
			found := false
			for _, def := range myBlocks {
				if def.Typ == block.Type {
					found = true
					break
				}
			}
			if !found {
				continue
			}

			// Calculate a reasonable range for the definition. Ideally we want
			// the range of the type and labels. hclsyntax.Block doesn't have a
			// single DefRange, but we can construct one or just use the block's
			// range (which includes the body). For style linting, usually the
			// header is what we care about. Here we construct it from TypeRange
			// and LabelRanges.
			defRange := block.TypeRange
			if len(block.LabelRanges) > 0 {
				lastLabel := block.LabelRanges[len(block.LabelRanges)-1]
				defRange = hcl.Range{
					Filename: defRange.Filename,
					Start:    defRange.Start,
					End:      lastLabel.End,
				}
			}

			for _, checkFunc := range checkFuncs {
				checkFunc(runner, rule, defRange, typ, name, synonym)
			}
		}
	}

	return nil
}

// WalkTokens iterates over all files in the root module, lexes them, and
// applies the check function to each token.
func WalkTokens[T any](
	runner tflint.Runner,
	rule T,
	checkFunc func(tflint.Runner, T, hclsyntax.Token),
) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		return nil
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for filename, file := range files {
		tokens, diags := hclsyntax.LexConfig(file.Bytes, filename, hcl.InitialPos)
		if diags.HasErrors() {
			return diags
		}

		for _, token := range tokens {
			checkFunc(runner, rule, token)
		}
	}

	return nil
}
