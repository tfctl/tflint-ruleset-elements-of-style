// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// MisOrderedMessage is the message emitted when meta arguments are not in the
// configured order.
const MisOrderedMessage = "Meta arguments should be ordered consistently"

// checkOrder verifies that arguments in a block respect the Order config.
// Arguments in Order.First must appear before all other arguments. Arguments
// in Order.Last must appear after all other arguments. Arguments not in either
// list can appear anywhere in between.
func checkOrder(runner tflint.Runner, r *Rule, block *hclsyntax.Block) {
	if len(r.Config.Order) == 0 {
		return
	}
	order := r.Config.Order[0]
	if len(order.First) == 0 && len(order.Last) == 0 {
		return
	}

	// Build sets for First and Last arguments.
	firstSet := make(map[string]bool, len(order.First))
	for _, name := range order.First {
		firstSet[name] = true
	}
	lastSet := make(map[string]bool, len(order.Last))
	for _, name := range order.Last {
		lastSet[name] = true
	}

	// Collect all attributes with their line numbers and classification.
	type attrInfo struct {
		name       string
		lineNumber int
		isFirst    bool
		isLast     bool
	}

	var attrs []attrInfo
	for name, attr := range block.Body.Attributes {
		info := attrInfo{
			name:       name,
			lineNumber: attr.SrcRange.Start.Line,
			isFirst:    firstSet[name],
			isLast:     lastSet[name],
		}
		attrs = append(attrs, info)
	}

	if len(attrs) == 0 {
		return
	}

	// Sort by line number to get the actual order in the source.
	for i := 0; i < len(attrs)-1; i++ {
		for j := i + 1; j < len(attrs); j++ {
			if attrs[i].lineNumber > attrs[j].lineNumber {
				attrs[i], attrs[j] = attrs[j], attrs[i]
			}
		}
	}

	// Check that First arguments appear before all non-First arguments. Once
	// we see a non-First argument, any subsequent First argument is out of
	// order.
	seenNonFirst := false
	for _, attr := range attrs {
		if !attr.isFirst {
			seenNonFirst = true
		} else if seenNonFirst {
			a := block.Body.Attributes[attr.name]
			r.emitIssue(runner, MisOrderedMessage, a.SrcRange)
			return
		}
	}

	// Check that Last arguments appear after all non-Last arguments. Once we
	// see a Last argument, any subsequent non-Last argument is out of order.
	seenLast := false
	for _, attr := range attrs {
		if attr.isLast {
			seenLast = true
		} else if seenLast {
			// A non-Last argument appears after a Last argument.
			a := block.Body.Attributes[attr.name]
			r.emitIssue(runner, MisOrderedMessage, a.SrcRange)
			return
		}
	}
}
