// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// checkThreshold checks if the comment ratio is below the threshold.
func checkThreshold(r *CommentsRule, runner tflint.Runner) error {
	if r.Config.Threshold == nil {
		return nil
	}
	threshold := *r.Config.Threshold

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for filename, file := range files {
		tokens, diags := hclsyntax.LexConfig(file.Bytes, filename, hcl.InitialPos)
		if diags.HasErrors() {
			return diags
		}

		linesWithContent := make(map[int]bool)
		linesWithComment := make(map[int]bool)

		for _, token := range tokens {
			if token.Type == hclsyntax.TokenEOF || token.Type == hclsyntax.TokenNewline {
				continue
			}

			startLine := token.Range.Start.Line
			endLine := token.Range.End.Line

			// Adjust endLine if the token ends at the start of the next line (e.g. # comments)
			if token.Range.End.Column == 1 && endLine > startLine {
				endLine--
			}

			for i := startLine; i <= endLine; i++ {
				if token.Type == hclsyntax.TokenComment {
					linesWithComment[i] = true
				}
				linesWithContent[i] = true
			}
		}

		totalLines := len(linesWithContent)
		commentLines := len(linesWithComment)

		if totalLines == 0 {
			continue
		}

		ratio := float64(commentLines) / float64(totalLines)
		if ratio < threshold {
			rng := hcl.Range{
				Filename: filename,
				Start:    hcl.Pos{Line: 1, Column: 1},
				End:      hcl.Pos{Line: 1, Column: 1},
			}

			message := fmt.Sprintf("Comments ratio is %.0f percent (minimum threshold %.0f percent)", ratio*100, threshold*100)
			if err := runner.EmitIssue(r, message, rng); err != nil {
				logger.Error(err.Error())
			}
			logger.Debug(message)
		}
	}
	return nil
}
