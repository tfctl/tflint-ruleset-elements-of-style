// Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// https://developer.hashicorp.com/terraform/language/block/module#http-urls
var validHTTPSExtensions = []string{
	".zip",
	".bz2", ".tar.bz2", ".tar.tbz2", ".tbz2",
	".gz", ".tar.gz", ".tgz",
	".xz", ".tar.xz", ".txz",
}

// checkModuleSourceVersion checks for proper module source versioning.
func checkModuleSourceVersion(runner tflint.Runner, r *MetaRule, block *hclsyntax.Block) {
	sourceAttr, exists := block.Body.Attributes["source"]
	if !exists {
		return
	}

	sourceExpr := sourceAttr.Expr
	sourceVal, diags := sourceExpr.Value(&hcl.EvalContext{})
	if diags.HasErrors() {
		return
	}

	if sourceVal.Type() != cty.String {
		return
	}

	source := sourceVal.AsString()

	if isLocalSource(source) {
		return
	}

	if isGitSource(source) {
		if !strings.Contains(source, "ref=") {
			r.emitIssue(runner, "Git module source should specify ref parameter.", block.Range())
		}
		return
	}

	if isHTTPSSource(source) {
		found := false
		for _, x := range validHTTPSExtensions {
			if strings.HasSuffix(source, x) {
				found = true
				break
			}
		}
		if !found {
			// https://regex101.com/r/1lV7Lx/1
			matched, _ := regexp.MatchString(`^https://.*\?.*archive=\..*&?`, source)
			if matched {
				found = true
			}
		}

		if !found {
			r.emitIssue(runner, "https module source should specify a valid archive extension.", block.Range())
		}
		return
	}

	if isMercurialSource(source) {
		if !strings.Contains(source, "#") {
			r.emitIssue(runner, "Mercurial module source should specify #revision.", block.Range())
		}
		return
	}

	if isRegistrySource(source) {
		versionAttr, exists := block.Body.Attributes["version"]
		if !exists {
			r.emitIssue(runner, "Module from registry should specify version.", block.Range())
			return
		}

		versionExpr := versionAttr.Expr
		versionVal, diags := versionExpr.Value(&hcl.EvalContext{})
		if diags.HasErrors() {
			return
		}

		if versionVal.Type() != cty.String {
			return
		}

		version := versionVal.AsString()
		constraints := strings.Split(version, ",")
		for _, c := range constraints {
			c = strings.TrimSpace(c)
			if strings.HasPrefix(c, "~>") {
				ver := strings.TrimSpace(strings.TrimPrefix(c, "~>"))
				if !strings.Contains(ver, ".") {
					r.emitIssue(runner, "Pessimistic version constraint should specify at least major and minor version.", block.Range())
					return
				}
				continue
			}

			if strings.HasPrefix(c, ">") {
				r.emitIssue(runner, "Version constraint > or >= should not be used. Use ~> or exact version.", block.Range())
				return
			}
		}
		return
	}

}

// isGitSource checks if the source is a Git repository.
func isGitSource(source string) bool {
	return strings.HasPrefix(source, "git::") ||
		strings.HasPrefix(source, "git@") ||
		strings.HasPrefix(source, "github.com") ||
		strings.HasPrefix(source, "bitbucket.org")
}

func isHTTPSSource(source string) bool {
	return strings.HasPrefix(source, "https://")
}

// isLocalSource checks if the source is a local path.
func isLocalSource(source string) bool {
	return strings.HasPrefix(source, "./") || strings.HasPrefix(source, "../") || strings.HasPrefix(source, "/")
}

// isMercurialSource checks if the source is a Mercurial repository.
func isMercurialSource(source string) bool {
	return strings.HasPrefix(source, "hg::")
}

// isRegistrySource checks if the source is from Terraform registry.
func isRegistrySource(source string) bool {
	// https://regex101.com/r/pBie8i/2
	registryRegex := regexp.MustCompile(`^(?:([a-zA-Z0-9.-]+)/)?([^/]+)/([^/]+)/([^/]+)$`)
	return registryRegex.MatchString(source)
}
