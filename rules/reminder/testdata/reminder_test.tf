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

# #########
# Tests that will not emit issues.

# No reminder.
# Reminder BUG in middle of text.

#

locals {
  good = 1 # No reminder.
}
