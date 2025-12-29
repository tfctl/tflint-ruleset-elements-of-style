# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

# Tests for camel rule: names must be lowercase alphanumeric and underscores.

# #########
# Tests that will emit issues (invalid names).

variable "CamelCase" {}
variable "kebab-case" {}
variable "with.dots" {}

locals {
  CamelCase = 1
}

check "CamelCase" {
  assert {
    condition     = local.CamelCase > 0
    error_message = "Must be > 0."
  }
}

ephemeral "random_password" "kebab-case" {
  length = 8
}

module "with.dots" {
  source = "./modules/"
}

# #########
# Tests that will not emit issues (valid names).

variable "snake_case" {}
variable "lower123" {}
variable "_underscore" {}
variable "a" {}

locals {
  valid_name    = 1
  another_valid = 2
}

resource "aws_instance" "valid" {
  ami = "ami-12345678"
}
