# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

rule "eos_dry" {
  enabled = true
}

rule "eos_dry_disabled" {
  enabled = false
}

rule "eos_dry_info" {
  level = "info"
}

rule "eos_dry_threshold" {
  threshold = 5
}
