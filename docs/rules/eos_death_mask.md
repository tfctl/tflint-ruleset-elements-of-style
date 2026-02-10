# eos_death_mask

Identify death masks - blocks of code commented out and left behind after its demise. Like preserving a corpse in case we ever need it again. We won’t.

## Example

```hcl
# resource "terraform_data" "example" {
#   name = "example"
# }
```

```
$ tflint
1 issue(s) found:

Warning: Avoid commented-out code. (eos_death_mask)

  on main.tf line 1:
   1: # resource "terraform_data" "example" {
   2: #   name = "example"
   3: # }

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_death_mask.md
```

## Why

Commented-out code creates confusion and clutter. It is often unclear why the code was commented out, whether it is still relevant, or if it should be deleted. Version control systems (like Git) are the appropriate place to store history of deleted code.

## How To Fix

Delete the commented-out code. If you need to preserve it for reference, commit it to version control before removing it.

```hcl
resource "terraform_data" "example" {
  name = "example"
}
```

The rule can be ignored with:

```hcl
# tflint-ignore: eos_death_mask
# resource "terraform_data" "example" {
#   name = "example"
# }
```

## Configuration

This rule is enabled by default and can be disabled with:

```hcl
rule "eos_death_mask" {
  enabled = false
}
```

Configure the severity:

```hcl
rule "eos_death_mask" {
  level = "error"  # Change severity to error
}
```
