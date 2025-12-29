# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

# THINK This is fragile af. Order dependent.  See comment at bottom.
# TODO Figure out how to use .tflint.hcl in a test.

# #########
# Tests that will emit issues.

variable "str_hung" {}

locals {
  hung_int = 1
}

check "hung_bool_check" {
  assert {
    condition     = local.hung_int > 0
    error_message = "Must be > 0."
  }
}

data "aws_caller_identity" "map_hung" {}

ephemeral "random_password" "hung_lst" {
  length = 8
}

module "hung_set_mod" {
  source = "./modules/"
}

output "num_hung" {
  value = local.hung_int
}

# tflint-ignore: eos_hungarian
resource "aws_instance" "str_hung" {
  ami = "ami-abc123"
}

output "foo_hung" {
  value = local.hung_int
}

output "hung_bar_hung" {
  value = local.hung_int
}

# #########
# Tests that will not emit issues.

variable "no_hung" {}

locals {
  no_hung = 1
}

resource "aws_instance" "no_hung" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}
