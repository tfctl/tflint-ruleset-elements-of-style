# This example will not emit and issue
# because there are a ton of comment
# lines in this source file. Thus there
# is nothing to worry about.  And this
# test should pass.
locals {
  x = 1
}

resource "foo" "bar" {
  count = 1
}
# One comment
