# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0
# no-cloc

# #########
# Tests that will emit issues.

# TEST
# Although invalid syntax (no parent block), this should still emit an issue.
# x = 1

# TEST
# Successive commented lines that collectively represent a valid expression.
# resource "resource" "dead" {
#   x = 1
# }

# TEST
# A single dead line embedded in a larger, "live" block.
resource "resource" "dead" {
  # x = 1
}

# TEST
# Mixed content, all of which is dead.
# resource "foo" "bar" {
#   # This is a nested comment
#   name = "baz"
# }

# #########
# Tests that will not emit issues.

resource "resource" "dead" {
  # A comment x = 1
}
