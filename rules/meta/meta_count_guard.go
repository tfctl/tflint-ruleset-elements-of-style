// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

const OnlyDynamicGuardMessage = "Avoid using count for anything other than dynamic guarding (condition ? 1 : 0)."
const GuardMustReturn10Message = "Count guard must return 1 or 0."

// checkCountGuard checks for proper count guard usage.
func checkCountGuard(runner tflint.Runner, r *MetaRule, attr *hclsyntax.Attribute) {
	expr := attr.Expr
	condExpr, ok := expr.(*hclsyntax.ConditionalExpr)

	// We want to check if it is a conditional expression:
	// condition ? true_val : false_val
	if !ok {
		// Allow literal 0 or 1.
		if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
			val := getLiteralValue(lit)
			if val == 0 || val == 1 {
				return
			}
		}

		r.emitIssue(runner, OnlyDynamicGuardMessage, attr.Range())
		return
	}

	// Check true/false results.
	if !isValidGuardResult(condExpr.TrueResult) || !isValidGuardResult(condExpr.FalseResult) {
		r.emitIssue(runner, GuardMustReturn10Message, attr.Range())
	}
}

// isValidGuardResult checks if the expression is a valid guard result (0 or 1).
func isValidGuardResult(expr hclsyntax.Expression) bool {
	val := getLiteralValue(expr)
	return val == 0 || val == 1
}

// getLiteralValue extracts the integer value from a literal expression.
func getLiteralValue(expr hclsyntax.Expression) int {
	if lit, ok := expr.(*hclsyntax.LiteralValueExpr); ok {
		if lit.Val.Type() == cty.Number {
			f, _ := lit.Val.AsBigFloat().Float64()
			if f == 0 {
				return 0
			}
			if f == 1 {
				return 1
			}
		}
	}
	return -1
}
