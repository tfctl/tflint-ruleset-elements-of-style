// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package rulehelper

import (
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestWalking(t *testing.T) {
	t.Run("WalkBlocks", testWalkBlocks)
	t.Run("WalkTokens", testWalkTokens)
}

type BlockInfo struct {
	Typ     string
	Name    string
	Synonym string
}

func testWalkBlocks(t *testing.T) {
	cases := []struct {
		Name      string
		Content   string
		MyBlocks  []BlockDef
		CheckFunc func(runner tflint.Runner, rule interface{}, rng hcl.Range, typ, name, synonym string)
		Setup     func() interface{}
		Assert    func(t *testing.T, collected interface{})
	}{
		{
			Name:     "simple_variable",
			Content:  `variable "test" { type = string }`,
			MyBlocks: []BlockDef{{Typ: "variable", Labels: []string{"name"}}},
			Setup: func() interface{} {
				return &[]BlockInfo{}
			},
			CheckFunc: func(runner tflint.Runner, rule interface{}, rng hcl.Range, typ, name, synonym string) {
				blocks := rule.(*[]BlockInfo)
				*blocks = append(*blocks, BlockInfo{Typ: typ, Name: name, Synonym: synonym})
			},
			Assert: func(t *testing.T, collected interface{}) {
				blocks := collected.(*[]BlockInfo)
				if len(*blocks) != 1 {
					t.Errorf("Expected 1 block, got %d", len(*blocks))
				}
				if (*blocks)[0].Typ != "variable" || (*blocks)[0].Name != "test" {
					t.Errorf("Unexpected block: %+v", (*blocks)[0])
				}
			},
		},
		{
			Name: "resource_and_output",
			Content: `resource "aws_instance" "example" { ami = "ami-123" }
output "id" { value = aws_instance.example.id }`,
			MyBlocks: AllLintableBlocks,
			Setup: func() interface{} {
				return &[]BlockInfo{}
			},
			CheckFunc: func(runner tflint.Runner, rule interface{}, rng hcl.Range, typ, name, synonym string) {
				blocks := rule.(*[]BlockInfo)
				*blocks = append(*blocks, BlockInfo{Typ: typ, Name: name, Synonym: synonym})
			},
			Assert: func(t *testing.T, collected interface{}) {
				blocks := collected.(*[]BlockInfo)
				if len(*blocks) != 2 {
					t.Errorf("Expected 2 blocks, got %d", len(*blocks))
				}
				foundResource := false
				foundOutput := false
				for _, b := range *blocks {
					if b.Typ == "aws_instance" && b.Name == "example" {
						foundResource = true
					}
					if b.Typ == "output" && b.Name == "id" {
						foundOutput = true
					}
				}
				if !foundResource || !foundOutput {
					t.Errorf("Missing expected blocks: resource=%v, output=%v", foundResource, foundOutput)
				}
			},
		},
		{
			Name: "complex_file",
			Content: func() string {
				content, _ := os.ReadFile("testdata/walk_test.tf")
				return string(content)
			}(),
			MyBlocks: AllLintableBlocks,
			Setup: func() interface{} {
				return &[]BlockInfo{}
			},
			CheckFunc: func(runner tflint.Runner, rule interface{}, rng hcl.Range, typ, name, synonym string) {
				blocks := rule.(*[]BlockInfo)
				*blocks = append(*blocks, BlockInfo{Typ: typ, Name: name, Synonym: synonym})
			},
			Assert: func(t *testing.T, collected interface{}) {
				blocks := collected.(*[]BlockInfo)
				expectedTypes := map[string]int{
					"variable":           2,
					"aws_ami":            1,
					"aws_instance":       1,
					"aws_security_group": 1,
					"output":             2,
					"module":             1,
					"check":              1,
				}
				actual := make(map[string]int)
				for _, b := range *blocks {
					actual[b.Typ]++
				}
				for typ, exp := range expectedTypes {
					if actual[typ] != exp {
						t.Errorf("For type %s, expected %d, got %d", typ, exp, actual[typ])
					}
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"test.tf": tc.Content})
			collected := tc.Setup()
			checkFuncs := []func(tflint.Runner, interface{}, hcl.Range, string, string, string){tc.CheckFunc}
			err := WalkBlocks(runner, tc.MyBlocks, collected, checkFuncs...)
			if err != nil {
				t.Fatalf("WalkBlocks failed: %v", err)
			}
			tc.Assert(t, collected)
		})
	}
}

func testWalkTokens(t *testing.T) {
	cases := []struct {
		Name      string
		Content   string
		CheckFunc func(runner tflint.Runner, rule interface{}, token hclsyntax.Token)
		Setup     func() interface{}
		Assert    func(t *testing.T, collected interface{})
	}{
		{
			Name:    "count_tokens",
			Content: `variable "test" { default = "value" }`,
			Setup: func() interface{} {
				return &[]hclsyntax.Token{}
			},
			CheckFunc: func(runner tflint.Runner, rule interface{}, token hclsyntax.Token) {
				tokens := rule.(*[]hclsyntax.Token)
				*tokens = append(*tokens, token)
			},
			Assert: func(t *testing.T, collected interface{}) {
				tokens := collected.(*[]hclsyntax.Token)
				if len(*tokens) == 0 {
					t.Errorf("Expected tokens, got none")
				}
			},
		},
		{
			Name:    "count_ident_tokens",
			Content: `variable "test" { default = "value" }`,
			Setup: func() interface{} {
				return &map[hclsyntax.TokenType]int{}
			},
			CheckFunc: func(runner tflint.Runner, rule interface{}, token hclsyntax.Token) {
				counts := rule.(*map[hclsyntax.TokenType]int)
				(*counts)[token.Type]++
			},
			Assert: func(t *testing.T, collected interface{}) {
				counts := collected.(*map[hclsyntax.TokenType]int)
				if (*counts)[hclsyntax.TokenIdent] < 1 {
					t.Errorf("Expected at least 1 IDENT token, got %d", (*counts)[hclsyntax.TokenIdent])
				}
			},
		},
		{
			Name: "complex_file",
			Content: func() string {
				content, _ := os.ReadFile("testdata/walk_test.tf")
				return string(content)
			}(),
			Setup: func() interface{} {
				return &map[hclsyntax.TokenType]int{}
			},
			CheckFunc: func(runner tflint.Runner, rule interface{}, token hclsyntax.Token) {
				counts := rule.(*map[hclsyntax.TokenType]int)
				(*counts)[token.Type]++
			},
			Assert: func(t *testing.T, collected interface{}) {
				counts := collected.(*map[hclsyntax.TokenType]int)
				total := 0
				for _, count := range *counts {
					total += count
				}
				if total == 0 {
					t.Errorf("Expected tokens in complex file, got 0")
				}
				if (*counts)[hclsyntax.TokenIdent] == 0 {
					t.Errorf("Expected IDENT tokens")
				}
				if (*counts)[hclsyntax.TokenOBrace] == 0 {
					t.Errorf("Expected OBrace tokens")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"test.tf": tc.Content})
			collected := tc.Setup()
			checkFunc := func(runner tflint.Runner, rule interface{}, token hclsyntax.Token) {
				tc.CheckFunc(runner, collected, token)
			}
			err := WalkTokens(runner, collected, checkFunc)
			if err != nil {
				t.Fatalf("WalkTokens failed: %v", err)
			}
			tc.Assert(t, collected)
		})
	}
}
