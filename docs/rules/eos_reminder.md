# eos_reminder

Identify comments containing reminder tags.

## Example

```hcl
# TODO: Fix this later
resource "terraform_data" "foo" {
  # ...
}
```

```
$ tflint
1 issue(s) found:

Warning: '# TODO: Fix this later' has reminder tag. (eos_reminder)

  on config.tf line 1:
  1: # TODO: Fix this later

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_reminder.md

```

## Why

Reminders (TODOs, FIXMEs, etc.) in code often get ignored and accumulate over time. It is generally better to track these tasks in an issue tracker where they can be prioritized and assigned. Keeping the codebase clean of these tags ensures that technical debt is visible and managed properly.

## Configuration

This rule is enabled by default and can be disabled with:

```hcl
rule "eos_reminder" {
  enabled = false
}
```

By default, the following `tags` are checked: `BUG`, `TODO`. Use the `tags`
argument to replace the default list, or use `extras` to add additional tags:

```hcl
rule "eos_reminder" {
  tags   = ["FIXME", "HACK"]        # Replace default tags
  extras = ["NOTGOOD", "REALBAD"]   # Add additional tags
  level  = "error"                  # Change severity to error
}
```

## How To Fix

Address the reminder and remove the comment, or move the task to an issue tracker.

```hcl
resource "terraform_data" "foo" {
  # ...
}
```

The rule can be ignored with:

```hcl
# tflint-ignore: eos_reminder
# TODO: Fix this later
resource "terraform_data" "foo" {
  # ...
}
```
