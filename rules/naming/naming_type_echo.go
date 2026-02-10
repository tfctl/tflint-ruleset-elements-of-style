// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package naming

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// checkTypeEcho checks if a word in type is echoed in the name.
func checkTypeEcho(runner tflint.Runner,
	rule *Rule, defRange hcl.Range,
	typ string, name string, synonym string) {

	// Assume there is no echo.
	echo := false

	lowerTyp := strings.ToLower(typ)   // aws_s3_bucket
	lowerName := strings.ToLower(name) // my_bucket
	synonymText := ""

	// For each word in type, see if it exists in name. Note that this is
	// comparing against the entire name value. The impact of that being a name of
	// "mys3widget" will match against a type of "aws_s3_bucket" (the "s3" word
	// exists in the name).
	for part := range strings.SplitSeq(lowerTyp, "_") {
		if strings.Contains(lowerName, part) {
			echo = true
			break
		}

		// Get synonyms for the word.
		var synonyms []string
		te := rule.Config.TypeEcho
		if te != nil && te.Synonyms != nil {
			synonyms = te.Synonyms[part]
			if synonym != "" {
				synonyms = append(synonyms, synonym)
			}
		}

		// Check synonyms.  This logic is different than above in that synonyms are
		// checked for on word boundaries. So "aws_s3_bucket" DOES NOT match
		// "mys3_widget", but would match "my_s3_widget".
		splitName := strings.SplitSeq(lowerName, "_-")
		for _, syn := range synonyms {
			for n := range splitName {
				if strings.Contains(n, syn) {
					echo = true
					synonymText = fmt.Sprintf(" (via synonym '%s')", syn)
					break
				}
			}

			// Don't bother checking more synonyms because we know we already have a
			// winner.
			if echo {
				break
			}
		}
	}

	if echo {
		if err := runner.EmitIssue(
			rule,
			fmt.Sprintf("Avoid echoing type \"%s\"%s in label \"%s\".", typ, synonymText, name),
			defRange,
		); err != nil {
			logger.Error(err.Error())
		}
	}
}
