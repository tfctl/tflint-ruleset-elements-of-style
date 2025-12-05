# THINK This is fragile af. Order dependent.  See comment at bottom.
# TODO Figure out how to use .tflint.hcl in a test.

# #########
# Tests that will emit issues.

variable "heredoc" {
  value = <<EOF
This is the heredoc.
EOF
}

locals {
  heredoc = <<EOF
This is the heredoc.
EOF
}

check "heredoc" {
  assert {
    condition     = local.heredoc == <<-EOF
This is the heredoc.
EOF
    error_message = "Wrong heredoc syntax."
  }
}

ephemeral "random_password" "heredoc" {
  lower = <<-EOF
This is the heredoc.
EOF
}

module "heredoc" {
  source = <<EOF
This is the heredoc.
EOF
}

output "heredoc" {
  value = <<-EOF
This is the heredoc.
EOF
}


# #########
# Tests that will not emit issues.

variable "no_heredoc" {
  value = <<-VAR
    This is the heredoc.
    VAR
}

locals {
  heredoc = <<-LOC
    This is the heredoc.
    LOC
}

check "no_heredoc" {
  assert {
    condition     = local.heredoc == <<-CHECK
      This is the heredoc.
      CHECK
    error_message = "Wrong heredoc syntax."
  }
}

ephemeral "random_password" "no_heredoc" {
  lower = <<-PW
    This is the heredoc.
    PW
}

module "no_heredoc" {
  source = <<-MODULE
    This is the heredoc.
    MODULE
}

output "no_heredoc" {
  value = <<-OUT
    This is the heredoc.
    OUT
}

# This looks like a heredoc, but it's not <<EOF
