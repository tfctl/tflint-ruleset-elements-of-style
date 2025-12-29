# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

locals {
  x = 1
}

resource "foo" "bar" {
  count = 1
}
# One comment
