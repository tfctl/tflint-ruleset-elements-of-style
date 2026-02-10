# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

# FAIL - First argument (for_each) appears after non-First argument (input).
resource "terraform_data" "first_after_other" {
  input    = each.value
  for_each = { za = { l = 1 }, kp = { x = 2 } }
}

# FAIL - Last argument (depends_on) appears before non-Last argument (input).
resource "terraform_data" "last_before_other" {
  depends_on = [terraform_data.first_after_other]
  input      = "test"
}

# FAIL - First argument (for_each) is correct, but Last argument (depends_on)
# appears before non-Last argument (input).
resource "terraform_data" "mixed_violation" {
  for_each   = { za = { l = 1 }, kp = { x = 2 } }
  depends_on = [terraform_data.first_after_other]
  input      = each.value
}

# PASS - No First or Last arguments.
resource "terraform_data" "pass_no_meta" {
  input = "1"
}

# PASS - First argument (for_each) before other arguments.
resource "terraform_data" "pass_first_correct" {
  for_each = { za = { l = 1 }, kp = { x = 2 } }
  input    = each.value
}

# PASS - Last argument (depends_on) after other arguments.
resource "terraform_data" "pass_last_correct" {
  input      = "test"
  depends_on = [terraform_data.pass_no_meta]
}

# PASS - First (for_each) before others, Last (depends_on) after others.
resource "terraform_data" "pass_both_correct" {
  for_each   = { za = { l = 1 }, kp = { x = 2 } }
  input      = each.value
  depends_on = [terraform_data.pass_no_meta]
}

# PASS - Multiple First arguments (for_each, count) in any order, before others.
resource "terraform_data" "pass_multiple_first" {
  count    = 1
  for_each = { za = { l = 1 }, kp = { x = 2 } }
  input    = each.value
}

# PASS - Multiple Last arguments (depends_on, provider) in any order, after
# others.
resource "terraform_data" "pass_multiple_last" {
  input      = "test"
  provider   = aws.west
  depends_on = [terraform_data.pass_no_meta]
}
