# eos_naming

Enforces naming conventions on Terraform blocks and locals. Currently checks for excessive length and "shouting" (all-uppercase names).

## Example

```hcl
resource "aws_instance" "very_long_instance_name" {
  # ...
}

variable "MY_VAR" {
  # ...
}
```

```
$ tflint
2 issue(s) found:

Warning: 'very_long_instance_name' is 22 characters and should not be longer than 16 (eos_naming)

  on config.tf line 1:
  1: resource "aws_instance" "very_long_instance_name" {

Warning: 'MY_VAR' should not be all uppercase (eos_naming)

  on config.tf line 5:
  5: variable "MY_VAR" {

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_naming.md

```

## Why

**Length**: Long names can make Terraform configurations harder to read and maintain. They can also cause issues with tools like `tfctl` or `terraform` by causing content to be pushed way past the right edge of the terminal. Keeping names concise encourages better naming practices and improves overall code quality.

**Shout**: All-uppercase names (shouting) can be harder to read and may imply a significance, such as constants or macros, that doesn't exist. Using snake_case, mixedCase, or lowercase names improves readability and aligns with common naming conventions.

## Configuration

| Name | Default | Description |
| --- | --- | --- |
| `level` | `"warning"` | TFLint alert level. |
| `length` | `{}` | Configuration for the length sub-rule. |
| `shout` | `{}` | Configuration for the shout sub-rule. |

### Length Sub-rule

| Name | Default | Description |
| --- | --- | --- |
| `enable` | `true` | Enable the length check. |
| `length` | `16` | Maximum length of names. |

### Shout Sub-rule

| Name | Default | Description |
| --- | --- | --- |
| `enable` | `true` | Enable the shout check. |

```hcl
rule "eos_naming" {
  enabled = true
  level = "warning"

  length = {
    enable = true
    length = 16
  }

  shout = {
    enable = true
  }
}
```

## How To Fix

Rename the block to a shorter, more descriptive name, or use snake_case/mixedCase instead of all-uppercase. The rule can be ignored with -

```hcl
# tflint-ignore: eos_naming
resource "aws_instance" "very_long_instance_name" {
  # ...
}

# tflint-ignore: eos_naming
variable "MY_VAR" {
  # ...
}
```