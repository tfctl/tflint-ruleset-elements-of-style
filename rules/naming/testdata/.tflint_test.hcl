rule "naming" {}

rule "naming_disabled" {
  enabled = false
}

rule "naming_negative_length" {
  length = -5
}

rule "naming_noshout" {
  shout = false
}

rule "naming_type_echo" {
  type_echo {
    enabled = true
    synonyms = {
      "aws_instance" = ["vm", "box"]
    }
  }
}

rule "naming_type_echo_disabled" {
  type_echo {
    enabled = false
  }
}

rule "naming_type_echo_custom" {
  type_echo {
    synonyms = {
      "aws_s3_bucket" = ["pail"]
    }
  }
}
