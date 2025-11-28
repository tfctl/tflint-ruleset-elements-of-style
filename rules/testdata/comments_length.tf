# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0
# no-cloc

# #########
# Tests that will emit issues.

# This comment is way too long and it will definitely extend beyond the eighty character limit that we have set for this rule.

resource "foo" "bar" {
  # Indented comment that is also way too long and should trigger the rule because it goes past column 80.
}

# #########
# Tests that will not emit issues.

# Good comment
// Good comment

# Short comment.

# This comment is very long but it contains a url http://example.com/very/long/url/that/makes/this/line/exceed/the/limit so it should be ignored.
