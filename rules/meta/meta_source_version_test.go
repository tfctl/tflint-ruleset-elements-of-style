// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	"flag"
	"os"
	"testing"

	"github.com/staranto/tflint-ruleset-elements-of-style/internal/testhelper"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func testMetaSourceVersionRule(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	cases := []struct {
		Name    string
		Content string
		Want    []string
	}{
		{
			Name: "source_version",
			Content: func() string {
				content, _ := os.ReadFile("./testdata/meta_source_version_test.tf")
				return string(content)
			}(),
			Want: []string{
				"Git module source should specify ref parameter.",
				"Git module source should specify ref parameter.",
				"Mercurial module source should specify #revision.",
				"https module source should specify a valid archive extension.",
				"https module source should specify a valid archive extension.",
				"Module from registry should specify version.",
				"Module from registry should specify version.",
			},
		},
	}

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"meta.tf": tc.Content})
		rule := NewMetaRule()
		// source_version enabled by default

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		testhelper.AssertIssuesMessages(t, tc.Want, runner.Issues)
	}
}
