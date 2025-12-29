# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

plugin "elements-of-style" {
  enabled = true
}

rule "eos_death_mask" {
  enabled = true
}

rule "eos_death_mask_disabled" {
  enabled = false
}

rule "eos_death_mask_error" {
  enabled = false
  level = "error"
}
