# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

rule "eos_naming" {}

rule "eos_naming_disabled" {
  enabled = false
}

rule "eos_naming_negative_length" {
  length = -5
}

rule "eos_naming_noshout" {
  shout = false
}

rule "eos_naming_nosnake" {
  snake = false
}

rule "eos_naming_type_echo" {
  type_echo {
    enabled = true
    synonyms = {
      "aws_instance" = ["vm", "box"]
    }
  }
}

rule "eos_naming_type_echo_disabled" {
  type_echo {
    enabled = false
  }
}

rule "eos_naming_type_echo_custom" {
  type_echo {
    synonyms = {
      "aws_s3_bucket" = ["pail"]
    }
  }
}
