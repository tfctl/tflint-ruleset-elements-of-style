# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

rule "eos_hungarian" {
  tags = ["str", "int", "num", "bool", "list", "lst", "set", "map", "arr", "array"]
}

rule "eos_hungarian_custom_tags" {
  tags = ["foo", "bar"]
}

rule "eos_hungarian_disabled" {
  enabled = false
}
