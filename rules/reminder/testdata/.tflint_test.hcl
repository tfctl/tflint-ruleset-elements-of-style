# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

rule "eos_reminder" {}

rule "eos_reminder_disabled" {
  enabled = false
}

rule "eos_reminder_extras" {
  extras = ["NOTGOOD", "REALBAD"]
}

rule "eos_reminder_many_tags" {
  tags = ["BUG", "FIXME", "HORROR", "TODO"]
}
