# eos_comments

Identify non-standard comment styles: space after comment marker, maximum line length, and no block comments.

## Example

```hcl
#This is a jammed comment
//This is also jammed

/*
  Block comments are not allowed.
*/

# This comment is way too long and exceeds the configured column limit which defaults to 80 characters so it will trigger a warning.
```

```
$ tflint
4 issue(s) found:

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

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_comments.md
```

## Why

Readable comments improve code maintainability. "Jammed" comments (without a space after the `#` or `//` marker) are harder to read. Nothing disrupts readability more than a comment that disappears off the right side of the editor pane or wraps unnaturally. Block comments (`/* ... */`) are generally discouraged in favor of line comments (`#`) for consistency and better diffs. Enforcing a comment threshold ensures that code is adequately documented.

## Configuration

| Name | Default | Description |
| --- | --- | --- |
| `block` | `true` | Identify block comments (`/* ... */`). |
| `column` | `80` | Maximum column for a comment. Set to `0` to disable. |
| `jammed` | `true` | Identify jammed comments (e.g. `#comment`). |
| `level` | `"warning"` | TFLint alert level. |
| `threshold` | `null` | Minimum ratio of comment lines to code lines (0.0 - 1.0). |
| `url_bypass` | `true` | Allow comments with URLs to exceed the column limit. |

```hcl
rule "eos_comments" {
  enabled = true
  block = true
  column = 80
  jammed = true
  level = "warning"
  threshold = 0.1
  url_bypass = true
}
```
