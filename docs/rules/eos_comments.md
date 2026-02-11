# eos_comments

Enforces comment style guidelines: no jammed comments, no block comments, maximum line length, and minimum comment ratio.

## Sub-rules

| Sub-rule | Identifies | Default |
|----------|------------|---------|
| `block` | Block comments. | `true` |
| `eol` | End-of-line comments. | `true` |
| `jammed` | Comments without space after marker. | `true` |
| `length` | Comments exceeding line length. | `true` (column 80) |
| `threshold` | Files with low comment ratio. | `0.0` (disabled) |

## Example

```hcl
#This is a jammed comment
//This is also jammed

/*
  Block comments are not allowed.
*/

# This comment is way too long and exceeds the configured column limit which defaults to 80 characters so it will trigger a warning.

resource "aws_instance" "example" {
  ami = var.ami # EOL comment
}
```

```
$ tflint
5 issue(s) found:

Warning: Comment is jammed ('#This ...'). (eos_comments)

  on main.tf line 1:
   1: #This is a jammed comment

Warning: Comment is jammed ('//Thi ...'). (eos_comments)

  on main.tf line 2:
   2: //This is also jammed

Warning: Block comments not allowed. (eos_comments)

  on main.tf line 4:
   4: /*
   5:   Block comments are not allowed.
   6: */

Warning: Comment extends beyond column 80 to 126. (eos_comments)

  on main.tf line 8:
   8: # This comment is way too long and exceeds the configured column limit which defaults to 80 characters so it will trigger a warning.

Warning: EOL comments not allowed. (eos_comments)

  on main.tf line 11:
   11:   ami = var.ami # EOL comment

Reference: https://github.com/tfctl/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_comments.md
```

## Why

Readable comments improve code maintainability. "Jammed" comments (without a space after the `#` or `//` marker) are harder to read. Block comments (`/* ... */`) are generally discouraged in favor of line comments (`#`) for consistency and better diffs. Nothing disrupts readability more than a comment that disappears off the right side of the editor pane or wraps unnaturally. End-of-line comments can clutter code. Enforcing a comment threshold ensures that code is adequately documented.

## How To Fix

Add a space after the comment marker, convert block comments to line comments, break long comments across multiple lines, or move end-of-line comments to their own line.

```hcl
# This is a properly spaced comment.

# This comment is broken across multiple lines so it does not
# exceed the configured column limit.

resource "aws_instance" "example" {
  # This comment is on its own line.
  ami = var.ami
}
```

The rule can be ignored with:

```hcl
# tflint-ignore: eos_comments
#This jammed comment is intentional
resource "aws_instance" "example" {
  ami = var.ami
}
```

## Configuration

This rule is enabled by default and can be disabled with:

```hcl
rule "eos_comments" {
  enabled = false
}
```

Configure sub-rules individually:

```hcl
rule "eos_comments" {
  block   = false  # Allow block comments
  eol     = false  # Allow EOL comments
  jammed  = false  # Allow jammed comments
  length {
    column    = 100    # Set max column to 100
    allow_url = false  # Don't allow URLs to exceed limit
  }
  threshold = 0.1    # Require 10% comment ratio
  level     = "error"  # Change severity to error
}
```
