# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

# THINK This is fragile af. Order dependent.  See comment at bottom.
# TODO Figure out how to use .tflint.hcl in a test.

# #########
# Tests that will emit issues.

variable "variable_echo" {}

locals {
  local_echo = 1
}

check "check_echo" {
  assert {
    condition     = local.local_echo > 0
    error_message = "Must be > 0."
  }
}

data "aws_caller_identity" "caller_echo" {}

ephemeral "random_password" "password_echo" {
  length = 8
}

module "module_echo" {
  source = "./modules/"
}

output "output_echo" {
  value = local.local_echo
}

# tflint-ignore: eos_type_echo
resource "aws_instance" "instance_echo" {
  ami = "ami-12345678"
}

# #########
# Tests that will not emit issues.

variable "clean_var" {}

locals {
  clean_val = 1
}

resource "aws_instance" "clean_vm" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}


# Input for type_echo_test.go. Each of the blocks in the first ("emit issues")
# section has a matching test case in type_echo_test.go. What's more, the order
# of the blocks here must be synced with the order of the test cases. This is
# super fragile, but I don't have the will to think of a better way right now.
# The blocks align with common.go/allLintableBlocks.
