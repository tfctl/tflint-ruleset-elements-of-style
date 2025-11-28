# eos_heredoc

Identifies usage of standard heredoc syntax (`<<`) and suggest indented heredoc syntax (`<<-`). Also, optionally, identifies usage of "EOF" as a heredoc delimiter.

## Examples

```hcl
# Bad.
resource "terraform_data" "example" {
  input = <<EOF
#!/bin/bash
echo "hello"
EOF
}

# Good.
resource "terraform_data" "example" {
  input = <<-SHELL
    #!/bin/bash
    echo "hello"
  SHELL
}

resource "terraform_data" "config" {
  input = <<-CFG
    {
      "enabled": "true",
      "options": {
        "read_only": "true"
      }
    }
  CFG
}
```

```
$ tflint
2 issue(s) found:

Warning: Use indented heredoc (<<-) instead of standard heredoc (<<). (eos_heredoc)

  on main.tf line 10:
   10:   input = <<EOF

Warning: 'EOF' is used as a heredoc delimiter. (eos_heredoc)

  on main.tf line 10:
   10:   user_data = <<EOF

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_heredoc.md
```

## Why

Standard heredocs (`<<`) require the content to be left-aligned, which breaks the visual indentation hierarchy of the Terraform code. Indented heredocs (`<<-`) allow the content to be indented relative to the surrounding code, improving readability and maintainability.

Using "EOF" as a delimiter is generic and doesn't convey the content type. Plus, multiple identical delimiters in one source file will likely be confusing. Using descriptive delimiters like `SHELL`, `JSON`, `YAML`, etc., improves readability.

## How To Fix

Change `<<` to `<<-` and indent the content to match the surrounding code. Choose a delimiter that conveys the type of content or purpose.

```hcl
resource "terraform_data" "example" {
  input = <<-SHELL
    #!/bin/bash
    echo "hello"
  SHELL
}
```

The rule can be ignored with:

```hcl
# tflint-ignore: eos_heredoc
resource "terraform_data" "example" {
  input = <<EOF
#!/bin/bash
echo "hello"
EOF
}
```

## Configuration

This rule is enabled by default and can be disabled with:

```hcl
rule "eos_heredoc" {
  enabled = false
}
```
