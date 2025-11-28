# eos_type_echo

Identify type echoing in names.

## Example

```hcl
resource "terraform_data" "terraform_data_logging" {
  # ...
}

```

```
$ tflint
1 issue(s) found:

Warning: The type "terraform_data" is echoed in the label "terraform_data_logging" (eos_type_echo)

  on config.tf line 1:
  1: resource "terraform_data" "terraform_data_logging" {

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_type_echo.md

```

## Why

Type echoing (aka. type jittering or, sometimes, [Hungarian Notation](https://en.wikipedia.org/wiki/Hungarian_notation)) is considered a bad practice when writing Terraform.  In *all* cases, the Terraform and OpenTofu tooling displays the type (`terraform_data`) immediately adjacent to the label, or name, (`terraform_data_logging`) of the occurrence.

In the HCL language itself, the syntax is, for example -

```hcl
resource "terraform_data" "terraform_data_logging" {
  # ...
}
```
not -

```hcl
resource "terraform_data"
# A whole bunch of comments describing
# what this resource is about
  "terraform_data_logging" {
  # ...
}
```

When listing the contents of a state file (with `terraform state list` or `tfctl sq`), or executing a `plan/apply`, the output is -

```
terraform_data.terraform_data_logging
```

In *all* cases, you would "jitter" as you pronounced this - "terraform data logging". Neither "flows" as well as simply saying "logging".

Since HCL is a verbose language this can also quickly spin out of control if you were to write something like -

```hcl
resource "terraform_data" "terraform_data_primary_security_group" {
  # ...
}

resource "terraform_data" "terraform_data_security_group_ingress_rule" {
  security_group_id = terraform_data.terraform_data_primary_security_group.id
  # ...
}
```

It's much more readable and, thus, maintainable to write -

```hcl
resource "terraform_data" "primary" {
  # ...
}

resource "terraform_data" "ingress" {
  security_group_id = terraform_data.primary.id
  # ...
}
```

## Configuration

This rule is enabled by default and can be disabled with:

```hcl
rule "eos_type_echo" {
  enabled = false
}
```

Use the `synonyms` map to provide alternate type prefixes, such as `group`, `bucket`, or `variable`, for the rule to recognize.

## How To Fix

Rename the resource block to remove the repetitive jitter. The rule can be ignored with -

```tf
# tflint-ignore: eos_type_echo
resource "terraform_data" "terraform_data_logging" {
  # ...
}
```
