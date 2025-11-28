# eos_meta

Enforces Terraform meta-argument syntax conventions.

## Sub-rules

### count_guard

Ensures that `count` is only used for dynamic guarding (conditional creation) and not for loops or static values.

**Valid:**

```hcl
resource "aws_instance" "example" {
  count = var.enabled ? 1 : 0
}
```

**Invalid:**

```hcl
resource "aws_instance" "example" {
  count = var.item_count
}

resource "aws_instance" "example" {
  count = length(var.items)
}
```

## Configuration

This rule is enabled by default.

```hcl
rule "eos_meta" {
  enabled = true
}
```
