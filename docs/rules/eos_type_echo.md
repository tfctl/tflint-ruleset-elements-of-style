# eos_type_echo

Similar to Hungarian notation, type echoing, or jittering, is the practice of repeating parts of the block type in it's name. Terraform is already a quite verbose language.  Type echoing add no value as the full type and name are *always* presented adjacent to each other:

## Examples

```hcl
resource "aws_s3_bucket" "log_bucket" {
}

resource "aws_security_group" "inbound_group" {
}

$ tflint
2 issue(s) found:

Warning: The type "aws_s3_bucket" is echoed in the label "log_bucket" (eos_type_echo)

  on main.tf line 1:
  1: resource "aws_s3_bucket" "log_bucket" {

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_type_echo.md

Warning: The type "aws_secuirty_group" is echoed in the label "inbound_group" (eos_type_echo)

  on main.tf line 1:
  1: resource "aws_security_group" "inbound_group" {

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_type_echo.md
```

## Why

Type echoing is considered a bad practice when writing Terraform.  In *all* cases, the Terraform and OpenTofu tooling displays the type (`aws_s3_bucket`) immediately adjacent to the label, or name, (`log_bucket`) of the occurrence.

In the HCL language itself, the syntax is, for example:

```hcl
resource "aws_s3_bucket" "log_bucket" {
  # ...
}
```
And not:

```hcl
resource "aws_s3_bucket" {
# A whole bunch of comments describing
# what this resource is about
  "log_bucket" {
  # ...
}
```

When listing the contents of a state file (with `terraform state list` or `tfctl sq`), or executing a `plan/apply`, the output is -

```
aws_s3_bucket.logging_bucket
```

In *all* cases, you would "jitter" as you pronounced this - "s3 bucket logging bucket". Neither "flows" as well as simply saying "logging".

Since HCL is a verbose language this can also quickly spin out of control if you were to write something like -

```hcl
resource "aws_security_group" "inbound_group" {
}

resource "aws_security_group_ingress_rule" "inbound_group_ingress_rule" {
  security_group_id = aws_security_group.inbound_group.id
}

# The output echos, too!
output "inbound_rule_id_output" {
  value = aws_security_group_ingress_rule.inbound_security_ingress_rule.id
}
```

It's much more readable and, thus, maintainable to write:

```hcl
resource "aws_security_group" "inbound" {
}

resource "aws_security_group_ingress_rule" "inbound" {
  security_group_id = aws_security_group.inbound.id
}

output "inbound_rule" {
  value = aws_security_group_ingress_rule.inbound.id
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
