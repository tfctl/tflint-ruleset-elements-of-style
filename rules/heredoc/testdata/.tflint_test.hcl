# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

rule "eos_heredoc" {
  EOF = true
}

rule "eos_heredoc_disabled" {
  enabled = false
}

rule "eos_heredoc_no_eof" {
  EOF = false
}
