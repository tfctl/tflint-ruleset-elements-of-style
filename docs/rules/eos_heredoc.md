# eos_heredoc

Identify usage of standard heredoc syntax (`<<`) and suggest indented heredoc syntax (`<<-`).
Also optionally identifies usage of "EOF" as a heredoc delimiter.

## Example

```hcl
resource "aws_instance" "web" {
  # Bad
  user_data = <<EOF
#!/bin/bash
echo "hello"
EOF

  # Good
  user_data = <<-SHELL
    #!/bin/bash
    echo "hello"
  SHELL
}
```

```
$ tflint
2 issue(s) found:

Warning: Use indented heredoc (<<-) instead of standard heredoc (<<). (eos_heredoc)

  on main.tf line 3:
   3:   user_data = <<EOF

Warning: 'EOF' is used as a heredoc delimiter. (eos_heredoc)

  on main.tf line 3:
   3:   user_data = <<EOF

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_heredoc.md
```

## Why

Standard heredocs (`<<`) require the content to be left-aligned, which breaks the visual indentation hierarchy of the HCL code. Indented heredocs (`<<-`) allow the content to be indented relative to the surrounding code, improving readability and maintainability.

Using "EOF" as a delimiter is generic and doesn't convey the content type. Using descriptive delimiters like `SHELL`, `JSON`, `YAML`, etc., improves readability.

## Configuration

| Name | Default | Description |
| --- | --- | --- |
| `eof` | `true` | Identify usage of "EOF" as a heredoc delimiter. |
| `level` | `"warning"` | TFLint alert level. |

```hcl
rule "eos_heredoc" {
  enabled = true
  eof = true
  level = "warning"
}
```

## How To Fix

Change `<<` to `<<-` and indent the content to match the surrounding code.

```hcl
resource "aws_instance" "web" {
  user_data = <<-EOF
    #!/bin/bash
    echo "hello"
  EOF
}
```

The rule can be ignored with:

```hcl
# tflint-ignore: eos_heredoc
resource "aws_instance" "web" {
  user_data = <<EOF
#!/bin/bash
echo "hello"
EOF
}
```
