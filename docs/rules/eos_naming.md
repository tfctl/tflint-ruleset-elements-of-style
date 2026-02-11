# eos_naming

Enforces naming conventions on Terraform blocks and locals. Checks for excessive
length, "shouting" (all-uppercase names), snake_case enforcement, and type
echoing.

## Sub-rules

| Sub-rule | Description | Default |
|----------|-------------|---------|
| `length` | Checks that names do not exceed a configurable length limit. | `16` |
| `shout` | Checks that names are not all-uppercase (shouting). | `true` |
| `snake` | Checks that names consist only of lowercase letters, digits, and underscores. | `true` |
| `type_echo` | Checks that type names are not echoed in labels. | `true` |

## Example

```hcl
resource "terraform_data" "very_long_instance_name" {
  # ...
}

variable "MY_VAR" {
  # ...
}

variable "CamelCase" {
  # ...
}

resource "aws_s3_bucket" "log_bucket" {
  # ...
}
```

```
$ tflint
4 issue(s) found:

Warning: Avoid names longer than 16 ('very_long_instance_name' is 24). (eos_naming)

  on config.tf line 1:
  1: resource "terraform_data" "very_long_instance_name" {

Warning: Avoid SHOUTED names (MY_VAR) (eos_naming)

  on config.tf line 5:
  5: variable "MY_VAR" {

Warning: Names should be snake_case (CamelCase). (eos_naming)

  on config.tf line 9:
  9: variable "CamelCase" {

Warning: Avoid echoing type "aws_s3_bucket" in label "log_bucket". (eos_naming)

  on config.tf line 13:
  13: resource "aws_s3_bucket" "log_bucket" {

Reference: https://github.com/tfctl/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_naming.md
```

## Why

**Length**: Long names can make Terraform configurations harder to read and
maintain. They can also cause issues with tools like `tfctl` or `terraform` by
causing content to be pushed way past the right edge of the terminal. Keeping
names concise encourages better naming practices and improves overall code
quality.

**Shout**: All-uppercase names (shouting) can be harder to read and may imply a
significance, such as constants or macros, that doesn't exist. Using snake_case
improves readability and aligns with common naming conventions.

**Snake**: Names that include uppercase letters, hyphens, spaces, or other
special characters can be inconsistent and harder to work with in scripts or
automation. Restricting to lowercase alphanumeric and underscores ensures
consistency and compatibility.

**Type Echo**: Repeating parts of the block type in its name (e.g., `bucket` in
`log_bucket` for `aws_s3_bucket`) is redundant since the type and name are
always displayed adjacent to each other. This leads to unnecessary verbosity
like "s3 bucket log bucket".

## How To Fix

Rename the block to a shorter, more descriptive name, or use snake_case instead
of all-uppercase or mixed case. The rule can be ignored with:

```hcl
# tflint-ignore: eos_naming
resource "terraform_data" "very_long_instance_name" {
  # ...
}

# tflint-ignore: eos_naming
variable "MY_VAR" {
  # ...
}

# tflint-ignore: eos_naming
variable "CamelCase" {
  # ...
}
```

## Configuration

This rule is enabled by default and can be disabled with:

```hcl
rule "eos_naming" {
  enabled = false
}
```

Configure sub-rules individually using simple values:

```hcl
rule "eos_naming" {
  length = 24      # Set max name length (default: 16, use -1 to disable)
  shout  = false   # Disable shout check
  snake  = false   # Disable snake_case check
  level  = "error" # Change severity to error
}
```

### Type Echo Configuration

The `type_echo` sub-rule can be configured with a block to specify synonyms:

```hcl
rule "eos_naming" {
  type_echo {
    enabled  = true
    synonyms = {
      bucket = ["container", "store"]
      group  = ["sg", "secgroup"]
    }
  }
}
```