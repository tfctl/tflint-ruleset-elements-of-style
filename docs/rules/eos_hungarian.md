# eos_hungarian

Identify Hungarian notation in names.

## Example

```hcl
resource "aws_instance" "str_instance" {
  # ...
}
```

```
$ tflint
1 issue(s) found:

Warning: 'str_instance' uses Hungarian notation with 'str' (eos_hungarian)

  on config.tf line 1:
  1: resource "aws_instance" "str_instance" {

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_hungarian.md

```

## Why

Hungarian notation (encoding type information in variable names) is generally considered redundant in strongly typed languages or declarative configurations like Terraform where the type is often evident from the context (e.g., `resource "aws_instance"` clearly defines an instance). Avoiding it leads to cleaner and more readable code.

## Configuration

| Name | Default | Description |
| --- | --- | --- |
| `level` | `"warning"` | TFLint alert level. |
| `tags` | `["str", "int", "num", "bool", "list", "lst", "set", "map", "arr", "array"]` | List of Hungarian notation tags to identify. |

```hcl
rule "eos_hungarian" {
  enabled = true
  level = "warning"
  tags = ["str", "int", "num", "bool", "list", "lst", "set", "map", "arr", "array"]
}
```

## How To Fix

Rename the block to remove the Hungarian notation.

```hcl
resource "aws_instance" "web" {
  # ...
}
```

The rule can be ignored with:

```hcl
# tflint-ignore: eos_hungarian
resource "aws_instance" "str_instance" {
  # ...
}
```
