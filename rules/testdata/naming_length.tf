# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0
# no-cloc

# #########
# Tests that will emit issues.

variable "very_long_variable_name" {}

locals {
  very_long_local_name = 1
}

check "very_long_check_name" {
  assert {
    condition     = true
    error_message = "Must be true."
  }
}

data "aws_get_caller_identity" "very_long_data_name" {}

ephemeral "random_password" "very_long_ephemeral_name" {
  length = 8
}

module "very_long_module_name" {
  source = "./modules/"
}

output "very_long_output_name" {
  value = "test"
}

resource "aws_instance" "very_long_instance_name" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

resource "aws_instance" "very_long_instance_name_disabled" {
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
