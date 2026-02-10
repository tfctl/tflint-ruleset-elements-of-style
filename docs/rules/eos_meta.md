# eos_meta

Enforces Terraform meta-argument syntax conventions.

## Sub-rules

| Sub-rule | Identifies | Default |
|----------|------------|---------|
| `count_guard` | Improper count usage. | Always enabled |
| `order` | Meta-argument ordering. | `for_each`/`count` first, `depends_on`/`provider`/`lifecycle` last |
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

### order

Enforces consistent ordering of meta-arguments within blocks. By default,
`for_each` and `count` must appear before other arguments, while `depends_on`,
`provider`, and `lifecycle` must appear last.

**Valid:**

```hcl
resource "aws_instance" "example" {
  count = var.enabled ? 1 : 0

  ami           = var.ami
  instance_type = "t3.micro"

  lifecycle {
    create_before_destroy = true
  }
}
```

**Invalid:**

```hcl
resource "aws_instance" "example" {
  ami           = var.ami
  count         = var.enabled ? 1 : 0  # count should appear first
  instance_type = "t3.micro"
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

## Why

**count_guard**: Using `count` for anything other than conditional creation (e.g., `count = var.enabled ? 1 : 0`) can lead to confusing state changes when the count value changes. Terraform may destroy and recreate resources unexpectedly. Use `for_each` for iterating over collections.

**order**: Consistent ordering of meta-arguments makes configurations easier to scan and review. Placing `count` or `for_each` at the top immediately signals conditional or iterated resources. Placing lifecycle-related arguments at the bottom keeps them together and out of the way.

**source_version**: Unversioned module sources can lead to unexpected changes when upstream modules are updated. Pinning versions ensures reproducible infrastructure.

## How To Fix

For `count_guard`, use a ternary expression that returns 1 or 0:

```hcl
resource "terraform_data" "example" {
  count = var.enabled ? 1 : 0
}
```

For `order`, move meta-arguments to their expected positions:

```hcl
resource "aws_instance" "example" {
  for_each = var.instances  # First

  ami           = each.value.ami
  instance_type = each.value.type

  depends_on = [aws_vpc.main]  # Last
}
```

For `source_version`, add explicit version constraints:

```hcl
module "consul" {
  source  = "hashicorp/consul/aws"
  version = "~> 0.10"
}
```

The rule can be ignored with:

```hcl
# tflint-ignore: eos_meta
module "unversioned" {
  source = "hashicorp/consul/aws"
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

### Order Configuration

Customize which arguments must appear first and last in blocks:

```hcl
rule "eos_meta" {
  order {
    first = ["count", "for_each"]
    last  = ["depends_on", "lifecycle"]
  }
}
```

To disable ordering checks entirely, provide an empty order block:

```hcl
rule "eos_meta" {
  order {
    first = []
    last  = []
  }
}
```

### Source Version Configuration

The `source_version` sub-rule has no additional configuration options beyond a simple
boolean toggle. It does not accept a configuration block (for example,
`source_version { ... }` is not supported). To disable source version checks, set
`source_version = false` in the rule configuration. Use the top-level `level`
parameter to adjust severity.

