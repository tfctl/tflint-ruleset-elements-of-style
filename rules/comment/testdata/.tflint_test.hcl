# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

rule "eos_comments" {
  block      = true
  eol        = true
  jammed = true
  length {
    allow_url = true
    column     = 80
  }
  threshold = 0.2
}

rule "eos_comments_noblock" {
  block = false
}

rule "eos_comments_nojammed" {
  jammed = false
}

rule "eos_comments_nolength" {
  // length {
  //   column = 999
  // }
}

rule "eos_comments_nocolumn" {
  length {
    allow_url = true
    column = 0
  }
}

rule "eos_comments_nourl" {
  length {
    allow_url = false
    column = 80
  }
}

rule "eos_comments_threshold_fail" {
  threshold = 0.5
}

rule "eos_comments_threshold_good" {
  threshold = 0.1
}
