# THINK This is fragile af. Order dependent.  See comment at bottom.
# TODO Figure out how to use .tflint.hcl in a test.

# #########
# Tests that will emit issues.

locals {
  # TEST
  # Both literal1 and literal2 have the same, static string value.
  literal1 = "zakpxy"
  literal2 = "zakpxy"

  # TEST
  # Both interpolation1 and interpolation2 have the same string value which
  # contains an interpolation.
  interpolation1 = "1${local.literal1}"
  interpolation2 = "1${local.literal1}"

  # TEST
  # Both list1 and list2 have the exact same list value. Only test for the
  # entire, ordered contents of a list. The individual elements of a list should
  # never be tested by the dry rule.
  list1 = ["z", "a", "k", "p", "x", "y"]
  list2 = ["z", "a", "k", "p", "x", "y"]

  # TEST
  # Both map1 and map2 the exact same map values.
  map1 = {
    za = "x"
    kp = "y"
  }
  map2 = {
    za = "x"
    kp = "y"
  }

  # TEST
  # Both expr1 and expr2 have the exact same map expression.
  expr1 = { for k, v in local.map1 : v => k }
  expr2 = { for k, v in local.map1 : v => k }

  # TEST
  # Both expr3 and expr4 have the exact same map expression.
  expr3 = [for i, v in local.map1 : v]
  expr4 = [for i, v in local.map1 : v]
}

# TEST
# Both resources have the same trigger block.
resource "null_resource" "r1" {
  triggers = {
    source = "zakpxy"
  }
}
resource "null_resource" "r2" {
  triggers = {
    source = "zakpxy"
  }
}

# #########
# Tests that will not emit issues.

variable "good" {
  default = "zakpxy1"
}

locals {
  good1 = var.good
  good2 = var.good


  xDirectLiteral1 = local.literal1
  xDirectLiteral2 = local.literal1

  xliteral1 = "1zakpxy"
  xliteral2 = "2zakpxy"

  xinterpolation1 = "1${local.xliteral1}"
  xinterpolation2 = "1${local.xliteral2}"

  xlist1 = ["1", "z", "a", "k", "p", "x", "y"]
  xlist2 = ["2", "z", "a", "k", "p", "x", "y"]

  xl2m1 = { for i, v in local.xlist1 : i => v }
  xl2m2 = { for i, v in local.xlist2 : i => v }

  xmap1 = {
    za = "1x"
    kp = "y"
  }
  xmap2 = {
    za = "x"
    kp = "2y"
  }
}

resource "terraform_data" "each1" {
  for_each = local.xmap1
  input    = each.key
}

resource "terraform_data" "each2" {
  for_each = local.xmap1
  input    = each.key
}

resource "terraform_data" "count1" {
  count = length(local.xlist1)
  input = count.index
}

resource "terraform_data" "count2" {
  count = length(local.xlist1)
  input = count.index
}

module "module1" {
  source = "github.com/org/repo"
}
module "module2" {
  source = "github.com/org/repo"
}
