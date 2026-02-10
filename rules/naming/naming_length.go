// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package naming

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// checkNameLength checks if the name is too long.
func checkNameLength(runner tflint.Runner, r *Rule, defRange hcl.Range, _ string, name string, _ string) {
	limit := defaultLimit
	if r.Config.Length != nil {
		limit = *r.Config.Length
	}

	if len(name) > limit {
		message := fmt.Sprintf("Avoid names longer than %d ('%s' is %d).", limit, name, len(name))
		if err := runner.EmitIssue(r, message, defRange); err != nil {
			logger.Error(err.Error())
		}
		logger.Debug(message)
	}
}
