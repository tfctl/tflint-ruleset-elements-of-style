# Copyright (c) 2025 Steve Taranto <eos@gmail.com>.
# SPDX-License-Identifier: Apache-2.0
# no-cloc

# #########
# Tests that will emit issues.

# TEST

module "fail_gh_noref" {
  source = "github.com/eos/test.git"
}

module "fail_git_noref" {
  source = "git::github.com/eos/test.git"
}

module "fail_hg_noref" {
  source = "hg::https://example.com/eos/test"
}

module "fail_https_ext1" {
  source = "https://example.com/eos/test"
}

module "fail_https_ext2" {
  source = "https://example.com/eos/test.bad"
}

module "fail_reg_ver1" {
  source = "eos/module/zakpxy"
}

module "fail_reg_ver1" {
  source = "app.terraform.io/eos/module/zakpxy"
}

# #########
# Tests that will not emit issues.

module "pass_local" {
  module = "./modules"
}

module "pass_git_ref" {
  source = "github.com/eos/test.git?ref=v0.0.1"
}

module "pass_hg_ref" {
  source = "hg::https://example.com/eos/test#v0.0.1"
}

module "pass_https_ext" {
  source = "https://example.com/eos/test/v0.1.tar.gz"
}

module "pass_https_arc" {
  source = "https://example.com/eos/test/v0.1?archive=.zakpxy"
}

module "pass_reg_ver1" {
  source  = "eos/module/zakpxy"
  version = "0.0.1"
}
