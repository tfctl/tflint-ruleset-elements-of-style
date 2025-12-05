# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0
# no-cloc

# #########
# Tests that will emit issues.

variable "really_a_very_long_name" {}

locals {
  really_a_very_long_name = 1
}

check "really_a_very_long_name" {
  assert {
    condition     = true
    error_message = "Must be true."
  }
}

data "aws_get_caller_identity" "really_a_very_long_name" {}

ephemeral "random_password" "really_a_very_long_name" {
  length = 8
}

module "really_a_very_long_name" {
  source = "./modules/"
}

output "really_a_very_long_name" {
  value = "test"
}

resource "aws_instance" "really_a_very_long_name" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

resource "aws_instance" "really_a_very_long_name_disabled" {
  count         = 0
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

# #########
# Tests that will not emit issues.

variable "short" {}

locals {
  short = 1
}

check "short" {
  assert {
    condition     = true
    error_message = "Must be true."
  }
}

data "aws_get_caller_identity" "short" {}

ephemeral "random_password" "short" {
  length = 8
}

module "short" {
  source = "./modules/"
}

output "short" {
  value = "test"
}

resource "aws_instance" "short" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}
