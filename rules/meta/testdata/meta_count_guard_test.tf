# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0
# no-cloc

variable "prod_bool" {
  default = true
}

locals {
  prod_bool  = true
  prod_count = 1
  two        = 2
}

# #########
# Tests that will emit issues.

# # Only a conditional 0 : 1 or 1 : 0 is allowed.
resource "null_resource" "loop" {
  count = 2
}

# Only a conditional 0 : 1 or 1 : 0 is allowed.
# This issue will not be emitted when run via `tflint`. I believe this is
# because tflint is short-circuiting this because it knows the true condition
# will be taken. Whereas the go test runner does not.
resource "terraform_data" "no01guard" {
  count = local.prod_bool ? 0 : 2
  input = count.index
}

# TEST
# Only a conditional 0 : 1 or 1 : 0 is allowed.
resource "terraform_data" "emit2" {
  count = local.prod_bool ? 2 : 0
  input = count.index
}

# TEST
# Only a conditional 0 : 1 or 1 : 0 is allowed.
resource "bad" "emit" {
  count = local.prod_count
}

# TEST
# Only a conditional 0 : 1 or 1 : 0 is allowed.
resource "bad" "emit" {
  count = length([1, 2])
}

# #########
# Tests that will not emit issues.

# Same deal as above. Confusing af.
resource "good" "no_emit" {
  count = 0
}

resource "good" "no_emit" {
  count = 1
}

resource "good" "no_emit" {
  count = local.prod_bool ? 0 : 1
}

resource "good" "no_emit" {
  count = local.prod_bool ? 1 : 0
}

resource "good" "no_emit" {
  count = local.prod_count > 0 ? 0 : 1
}
