# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

# THINK This is fragile af. Order dependent.  See comment at bottom.
# TODO Figure out how to use .tflint.hcl in a test.

# #########
# Tests that will emit issues.

variable "SHOUT" {}
variable "SHOUT🤡" {}

locals {
  SHOUT = 1
}

check "SHOUT" {
  assert {
    condition     = local.SHOUT > 0
    error_message = "Must be > 0."
  }
}

data "aws_caller_identity" "SHOUT" {}

ephemeral "random_password" "SHOUT" {
  length = 8
}

module "SHOUT" {
  source = "./modules/"
}

output "SHOUT" {
  value = local.SHOUT
}


# #########
# Tests that will not emit issues.

# tflint-ignore: eos_shout
resource "aws_instance" "SHOUT" {
  ami = "ami-12345678"
}

variable "no_shout" {}
variable "no_shout🤡" {}


locals {
  no_shout = 1
}

resource "aws_instance" "no_shout" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}


# Input for shout_test.go. Each of the blocks in the first ("emit issues")
# section has a matching test case in shout_test.go. What's more, the order of
# the blocks here must be synced with the order of the test cases. This is super
# fragile, but I don't have the will to think of a better way right now.  The
# blocks align with common.go/allLintableBlocks.
