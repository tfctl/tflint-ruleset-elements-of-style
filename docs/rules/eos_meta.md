# eos_meta

Enforces Terraform meta-argument syntax conventions.

## Sub-rules

| Sub-rule | Identifies | Default |
|----------|------------|---------|
| `count_guard` | Improper count usage. | `true` |
| `source_version` | Module sources without required versioning. | `true` |

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

### source_version

Ensures that module sources specify required versioning information.

**Valid:**

```hcl
module "registry_with_version" {
  source  = "hashicorp/consul/aws"
  version = "~> 0.10"
}

module "git_with_ref" {
  source = "git::https://github.com/example/repo.git?ref=v1.0.0"
}

module "local" {
  source = "./modules/local"
}
```

**Invalid:**

```hcl
module "registry_no_version" {
  source = "hashicorp/consul/aws"
}

module "git_no_ref" {
  source = "git::https://github.com/example/repo.git"
}
```

## Configuration

This rule is enabled by default and can be disabled with:

```hcl
rule "eos_meta" {
  enabled = false
}
```

Configure sub-rules individually:

```hcl
rule "eos_meta" {
  source_version = false  # Disable source version checks
  level          = "error"  # Change severity to error
}
```

