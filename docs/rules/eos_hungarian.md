# eos_hungarian

Identify [Hungarian notation](https://en.wikipedia.org/wiki/Hungarian_notation) in names. It is quite common to unnecessarily include the data type of a variable in it's name. Terraform is not a language that lends itself to that and the type is often times abstracted away.

## Examples

```hcl
variable "length_int" {
  type = number
}

locals {
  str_username = "ubuntu"
}

$ tflint
2 issue(s) found:

Warning: 'length_int' uses Hungarian notation with 'int' (eos_hungarian)

  on main.tf line 1:
  1: variable "length_int" {

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_hungarian.md

Warning: 'str_username' uses Hungarian notation with 'str' (eos_hungarian)

  on main.tf line 6:
  6: str_username = "ubuntu"

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_hungarian.md
```

## Why

Hungarian notation (encoding type information in variable names) is generally redundant in strongly typed languages or declarative configurations like Terraform, where the type is usually obvious from the name of the variable. Avoiding Hungarian notation leads to cleaner, more readable code that is easier to refactor.

## How To Fix

Rename the Hungarian component of the variable name.

```hcl
variable "length" {
  type = number
}

locals {
  username = "ubuntu"
}
```

The rule can be ignored with:

```hcl
locals {
  # tflint-ignore: eos_hungarian
  username = "ubuntu"
}

```

## Configuration

This rule is enabled by default and can be disabled with:

```hcl
rule "eos_hungarian" {
  enabled = false
}
```

By default, the following `tags` are considered as Hungarian indicators - 	arr, array, bool, int, list, lst, str, map, num, set. Additional tags can be added in the `.tflint.hcl`:

```hcl
rule "eos_hungarian" {
  enabled = false
  tags = ["foo", "bar"]
}
```



The `tags` argument accepts a list of prefixes (for example `str`, `int`, `num`, or `bool`) that the rule should treat as Hungarian notation indicators.

