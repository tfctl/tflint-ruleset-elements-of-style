# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0
# no-cloc

# #########
# Tests that will emit issues.

# VCS sources should use an explicit ref (tag, release, etc).
module "fail_gh_no_ref" {
  source = "github.com/eos/test.git"
}

# VCS sources should use an explicit ref (tag, release, etc).
module "fail_git_no_ref" {
  source = "git::github.com/eos/test.git"
}

# VCS sources should use an explicit ref (tag, release, etc).
module "fail_hg_no_ref" {
  source = "hg::https://example.com/eos/test"
}

# https sources should specify an explicit file that has a known
# extension.
# https://developer.hashicorp.com/terraform/language/block/module#http-urls
module "fail_https_no_extension" {
  source = "https://example.com/eos/test"
}

# https sources should specify an explicit file that has a known
# extension.
module "fail_https_no_extension2" {
  source = "https://example.com/eos/test.bad"
}

# Registry sources should specify a version constraint.
module "fail_no_version" {
  source = "eos/module/zakpxy"
}

# Registry sources should specify a version constraint.
module "fail_no_version2" {
  source = "app.terraform.io/eos/module/zakpxy"
}

# The version constraint shouldn't be open ended.
module "fail_version_gt" {
  source  = "eos/module/zakpxy"
  version = "> 1.2.0"
}

# The version constraint shouldn't be open ended.
module "fail_version_gte" {
  source  = "eos/module/zakpxy"
  version = ">= 1.2.0"
}

# The version constraint shouldn't be open ended.
module "fail_version_mixed_gte" {
  source  = "eos/module/zakpxy"
  version = "~> 1.2.0, >= 1.0"
}

# The version constraint shouldn't be open ended. This is a synonym for
# >= 1.
module "fail_version_pessimistic_short" {
  source  = "eos/module/zakpxy"
  version = "~> 1"
}

# #########
# Tests that will not emit issues.

module "pass_git_ref" {
  source = "github.com/eos/test.git?ref=v0.0.1"
}

module "pass_hg_ref" {
  source = "hg::https://example.com/eos/test#v0.0.1"
}

module "pass_https_arc" {
  source = "https://example.com/eos/test/v0.1?archive=.zakpxy"
}

module "pass_https_ext" {
  source = "https://example.com/eos/test/v0.1.tar.gz"
}

module "pass_local" {
  module = "./modules"
}

module "pass_version_eq" {
  source  = "eos/module/zakpxy"
  version = "= 1.2.0"
}

module "pass_version_explicit" {
  source  = "eos/module/zakpxy"
  version = "0.0.1"
}

module "pass_version_lt" {
  source  = "eos/module/zakpxy"
  version = "< 2.0.0"
}

module "pass_version_lte" {
  source  = "eos/module/zakpxy"
  version = "<= 2.0.0"
}

module "pass_version_pessimistic" {
  source  = "eos/module/zakpxy"
  version = "~> 1.2"
}

# This will actually pass the rule, but TF will fail this at init sincce the two
# constraints are in opposition and both can't be resolved by any one version.
module "pass_version_opposition" {
  source  = "eos/module/zakpxy"
  version = "~> 1.1, ~> 2.0"
}
