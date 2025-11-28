# eos_reminder

Identify comments containing reminder tags.

## Example

```hcl
# TODO: Fix this later
resource "aws_instance" "foo" {
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

| Name | Default | Description |
| --- | --- | --- |
| `level` | `"warning"` | TFLint alert level. |
| `tags` | `["BUG", "FIXME", "HACK", "TODO"]` | List of reminder tags to identify. |

```hcl
rule "eos_reminder" {
  enabled = true
  level = "warning"
  tags = ["BUG", "FIXME", "HACK", "TODO"]
}
```

## How To Fix

Address the reminder and remove the comment, or move the task to an issue tracker.

```hcl
resource "aws_instance" "foo" {
  # ...
}
```

The rule can be ignored with:

```hcl
# tflint-ignore: eos_reminder
# TODO: Fix this later
resource "aws_instance" "foo" {
  # ...
}
```
