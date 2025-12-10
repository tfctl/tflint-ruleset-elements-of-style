rule "reminder" {
  tags = ["TODO"]
}

rule "reminder_many_tags" {
  tags = ["BUG", "FIXME", "HORROR", "TODO"]
}

rule "reminder_disabled" {
  enabled = false
  tags = []
}
