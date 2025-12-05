rule "comments" {
  block      = true
  eol        = true
  jammed = true
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
  jammed = false
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
