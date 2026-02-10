# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

rule "eos_meta" {
  enabled = true
}

rule "eos_meta_disabled" {
  enabled = false
}

rule "eos_meta_order" {
  enabled = true
  order {
    first = ["zal"]
    last  = ["kpx"]
  }
}

rule "eos_meta_source_version_disabled" {
  source_version = false
}
