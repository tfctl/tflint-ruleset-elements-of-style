rule "hungarian" {
  tags = ["str", "int", "num", "bool", "list", "lst", "set", "map", "arr", "array"]
}

rule "hungarian_custom_tags" {
  tags = ["foo", "bar"]
}

rule "hungarian_disabled" {
  enabled = false
}
