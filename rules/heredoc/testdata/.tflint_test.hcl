rule "heredoc" {
  EOF = true
}

rule "heredoc_disabled" {
  enabled = false
}

rule "heredoc_no_eof" {
  EOF = false
}
