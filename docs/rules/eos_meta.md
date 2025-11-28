# eos_meta

Enforces Terraform meta-argument syntax conventions.

## Sub-rules

### count_guard

Ensures that `count` is only used for dynamic guarding (conditional creation) and not for loops or static values.

**Valid:**

```hcl
resource "terraform_data" "example" {
  count = var.enabled ? 1 : 0
}
```

**Invalid:**

```hcl
resource "terraform_data" "example" {
  count = var.item_count
}

resource "terraform_data" "example" {
  count = length(var.items)
}
```

## Configuration

This rule is enabled by default and can be disabled with:

```hcl
rule "eos_meta" {
  enabled = false
}
```

