# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

# #########
# Tests that will emit issues.

// TODO Reminder found.
# TODO Reminder found.

locals {
  bad = 1 # TODO Reminder found.

  # HORROR This is so bad.
  password = "abc123!"
}

resource "terraform_data" "fixme" { # FIXME Reminder found.
}

/*
  FIXME
  Reminder found.
*/

// NOTGOOD This is really not a good idea.
// REALBAD This is a really bad idea.

# #########
# Tests that will not emit issues.

# No reminder.
# Reminder BUG in middle of text.

#

locals {
  good = 1 # No reminder.
}
