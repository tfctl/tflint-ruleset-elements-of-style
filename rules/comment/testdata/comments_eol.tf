# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0
# no-cloc

# #########
# Tests that will emit issues.

locals {
  eol1 = 1 # This is an eol comment.
  eol2 = 2 // This is an eol comment.
  eol3 = {
    eola = 1 # This is an eol comment.
  }
  eol4 = "4" /* This is an eol comment. */
}

# #########
# Tests that will not emit issues.

# This is not an eol comment.
// This is not an eol comment.

locals {
  # This is not an eol comment.
  eol1 = 1
  # This is not an eol comment.
  eol2 = 2
  eol3 = {
    # This is not an eol comment.
    eola = 1
  }
}
