# eos_dry

Identify repeated values (strings, interpolations, lists, maps, sets, expressions) to encourage DRY (Don't Repeat Yourself) principles.

**Note:** Numeric and boolean values are not checked by this rule.

## Example

```hcl
locals {
  # Literals
  val1 = "some-value"
  val2 = "some-value" # Repeated

  # Interpolations
  interp1 = "prefix-${var.suffix}"
  interp2 = "prefix-${var.suffix}" # Repeated

  # Lists
  list1 = ["a", "b"]
  list2 = ["a", "b"] # Repeated

  # Sets (checked as lists)
  set1 = toset(["x", "y"])
  set2 = toset(["x", "y"]) # Repeated

  # Maps
  map1 = {
    key = "value"
  }
  map2 = {
    key = "value"
  } # Repeated

  # Expressions
  expr1 = { for k, v in var.map : k => v }
  expr2 = { for k, v in var.map : k => v } # Repeated
}
```

```
$ tflint
6 issue(s) found:

Warning: Value '"some-value"' is repeated 2 times. (eos_dry)

  on main.tf line 4:
   4:   val2 = "some-value"

Warning: Value '"prefix-${var.suffix}"' is repeated 2 times. (eos_dry)

  on main.tf line 8:
   8:   interp2 = "prefix-${var.suffix}"

Warning: List is repeated 2 times. (eos_dry)

  on main.tf line 12:
   12:   list2 = ["a", "b"]

Warning: List is repeated 2 times. (eos_dry)

  on main.tf line 16:
   16:   set2 = toset(["x", "y"])

Warning: Map is repeated 2 times. (eos_dry)

  on main.tf line 22:
   22:   map2 = {

Warning: Map is repeated 2 times. (eos_dry)

  on main.tf line 28:
   28:   expr2 = { for k, v in var.map : k => v }
```

## Why

Repeating values can lead to maintenance issues. If a value needs to change, it must be updated in multiple places. Using a local value or variable ensures consistency and easier updates.

## Configuration

This rule has no specific configuration.

```hcl
rule "eos_dry" {
  enabled = true
}
```
