// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	"flag"
	"testing"
)

func TestMeta(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Run("CountGuard", testMetaCountGuardRule)
	t.Run("SourceVersion", testMetaSourceVersionRule)
}
