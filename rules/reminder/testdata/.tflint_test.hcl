rule "reminder" {
  tags = ["TODO", "FIXME"]
}

rule "reminder_custom_tags" {
  tags = ["HACK", "BUG"]
}

rule "reminder_disabled" {
  enabled = false
}
