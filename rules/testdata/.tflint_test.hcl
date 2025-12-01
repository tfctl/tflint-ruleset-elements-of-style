rule "comments" {
  block      = true
  eol        = true
  jammed {
    enabled = true
    tails = true
  }
  length {
    allow_url = true
    column     = 80
  }
  threshold = 0.2
}

rule "comments_noblock" {
  block = false
}

rule "comments_nojammed" {
  jammed {
    enabled = false
  }
}

rule "comments_nolength" {
  // length {
  //   column = 999
  // }
}

rule "comments_nocolumn" {
  length {
    allow_url = true
    column = 0
  }
}

rule "comments_nourl" {
  length {
    allow_url = false
    column = 80
  }
}

rule "death_mask" {
  enabled = true
}

rule "dry" {
  enabled = true
}

rule "dry_disabled" {
  enabled = false
}

rule "dry_info" {
  level = "info"
}

rule "heredoc" {
  EOF = true
}

rule "heredoc_disabled" {
  enabled = false
}

rule "heredoc_no_eof" {
  EOF = false
}

rule "hungarian" {
  tags = ["str", "int", "num", "bool", "list", "lst", "set", "map", "arr", "array"]
}

rule "hungarian_custom_tags" {
  tags = ["foo", "bar"]
}

rule "hungarian_disabled" {
  enabled = false
}

rule "naming" {
  length {
    limit = 13
  }
}

rule "naming_disabled" {
  enabled = false
}

rule "naming_negative_limit" {
  length {
    limit = -5
  }
}

rule "naming_nolength" {
  length {
    enabled = false
  }
}

rule "naming_noshout" {
  shout {
    enabled = false
  }
}

rule "meta" {
  enabled = true
}

rule "reminder" {
  tags = ["TODO", "FIXME"]
}

rule "reminder_custom_tags" {
  tags = ["HACK", "BUG"]
}

rule "reminder_disabled" {
  enabled = false
}

rule "type_echo" {
  synonyms = {
    "aws_instance" = ["vm", "box"]
  }
}

rule "type_echo_custom" {
  synonyms = {
    "aws_s3_bucket" = ["pail"]
  }
}

rule "type_echo_disabled" {
  enabled = false
}
