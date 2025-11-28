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
  ami = "ami-12345678"
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


# Input for hungarian_test.go. Each of the blocks in the first ("emit issues")
# section has a matching test case in hungarian_test.go. What's more, the order of
# the blocks here must be synced with the order of the test cases. This is super
# fragile, but I don't have the will to think of a better way right now.  The
# blocks align with common.go/allLintableBlocks.
