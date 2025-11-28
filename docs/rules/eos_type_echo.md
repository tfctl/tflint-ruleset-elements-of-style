# eos_type_echo

Identify type echoing in names.

## Example

```hcl
resource "aws_s3_bucket" "logging-bucket" {
  # ...
}

```

```
$ tflint
1 issue(s) found:

Warning: The type "aws_s3_bucket" is echoed in the label "logging-bucket" (eos_type_echo)

  on config.tf line 1:
  1: resource "aws_s3_bucket" "logging-bucket" {

Reference: https://github.com/staranto/tflint-ruleset-elements-of-style/blob/main/docs/rules/eos_type_echo.md

```

## Why

Type echoing (aka. type jittering or, sometimes, [Hungarian Notation](https://en.wikipedia.org/wiki/Hungarian_notation)) is considered a bad practice when writing Terraform.  In *all* cases, the Terraform and OpenTofu tooling displays the type (`aws_s3_bucket`) immediately adjacent to the label, or name, (`logging-bucket`) of the occurence.

In the HCL language itself, the syntax is, for example -

```hcl
resource "aws_s3_bucket" "logging-bucket" {
  # ...
}
```
not -

```hcl
resource "aws_s3_bucket"
# A whole bunch of comments describing
# what this resource is about
  "logging-bucket" {
  # ...
}
```

When listing the contents of a state file (with `terraform state list` or `tfctl sq`), or executing a `plan/apply`, the output is -

```
aws_s3_bucket.logging-bucket
```

In *all* cases, you would "jitter" as you pronounced this - "aws S3 bucket logging bucket". Or, even more pronounced - "bucket logging bucket".  Neither "flows" as well as simply saying "S3 bucket logging".

Since HCL is a verbose language this can also quickly spin out of control if you were to write something like -

```hcl
resource "aws_security_group" "primary_security_group" {
  # ...
}

resource "aws_vpc_security_group_ingress_rule" "security_group_ingress_rule" {
  security_group_id = aws_security_group.primary_security_group.id
  # ...
}
```

It's much more readable and, thus, maintanable to write -

```hcl
resource "aws_security_group" "primary" {
  # ...
}

resource "aws_vpc_security_group_ingress_rule" "ingress" {
  security_group_id = aws_security_group.primary.id
  # ...
}
```

## Configuration

| Name | Default | Description |
| --- | --- | --- |
| `level` | `"warning"` | TFLint alert level. |
| `synonyms` | `{}` | Map of type synonyms to identify. |

```hcl
rule "eos_type_echo" {
  enabled = true
  level = "warning"
  synonyms = {
    "group"    = ["sg"],
    "bucket"   = ["bkt"]
    "variable" = ["var"]
  }
}
```

## How To Fix

Rename the resource block to remove the repetitive jitter. The rule can be ignored with -

```tf
# tflint-ignore: eos_type_echo
resource "aws_s3_bucket" "logging-bucket" {
  # ...
}
```
