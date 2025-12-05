# eos_naming

Enforces naming conventions on Terraform blocks and locals. Checks for excessive length, "shouting" (all-uppercase names), and ensures names are lowercase alphanumeric with underscores only.

## Sub-rules

| Sub-rule | Description | Default |
|----------|-------------|---------|
| `length` | Checks that names do not exceed a configurable length limit. | Enabled |
| `shout` | Checks that names are not all-uppercase (shouting). | Enabled |
| `camel` | Checks that names consist only of lowercase letters, digits, and underscores. | Disabled |

Note: If both `shout` and `camel` are enabled, only `camel` will be enforced as it is more restrictive.

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
```

```
$ tflint
3 issue(s) found:

Warning: 'very_long_instance_name' is 22 characters and should not be longer than 16 (eos_naming)

  on config.tf line 1:
  1: resource "terraform_data" "very_long_instance_name" {

Warning: 'MY_VAR' should not be all uppercase (eos_naming)

  on config.tf line 5:
  5: variable "MY_VAR" {

Warning: 'CamelCase' must be lowercase alphanumeric and underscores only (eos_naming)

  on config.tf line 9:
  9: variable "CamelCase" {

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_naming.md

```

## Why

**Length**: Long names can make Terraform configurations harder to read and maintain. They can also cause issues with tools like `tfctl` or `terraform` by causing content to be pushed way past the right edge of the terminal. Keeping names concise encourages better naming practices and improves overall code quality.

**Shout**: All-uppercase names (shouting) can be harder to read and may imply a significance, such as constants or macros, that doesn't exist. Using snake_case, mixedCase, or lowercase names improves readability and aligns with common naming conventions.

**Camel**: Names that include uppercase letters, hyphens, spaces, or other special characters can be inconsistent and harder to work with in scripts or automation. Restricting to lowercase alphanumeric and underscores ensures consistency and compatibility.

## Configuration

This rule is enabled by default and can be disabled with:

```hcl
rule "eos_naming" {
  enabled = false
}
```

Use the `length`, `shout`, and `camel` configuration blocks to adjust settings or disable sub-rules individually.

### Length Configuration

```hcl
rule "eos_naming" {
  length {
    enabled = true
    limit = 16
  }
}
```

### Shout Configuration

```hcl
rule "eos_naming" {
  shout {
    enabled = true
  }
}
```

### Camel Configuration

```hcl
rule "eos_naming" {
  camel {
    enabled = true
  }
}
```

## How To Fix

Rename the block to a shorter, more descriptive name, or use snake_case instead of all-uppercase or mixed case. The rule can be ignored with -

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