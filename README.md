# TFLint Ruleset: Elements of Style

This TFLint ruleset checks for idiomatic conventions and styles for Terraform code -much like a writing style guide. Style isn’t just aesthetic; it makes code easier to read, review, and modify safely. Consistently shared conventions reduce friction and help teams move faster with fewer surprises.

## Rules

|Name|Identifies|Link|
| --- | --- | --- |
|eos_comments|Problematic comment styles and structures.|[Link](docs/rules/eos_comments.md)|
|eos_death_mask|Blocks of commented out (ie. "dead") code.|[Link](docs/rules/eos_death_mask.md)|
|eos_dry|Repeatedly used values (strings, interpolations, lists, maps).|[Link](docs/rules/eos_dry.md)|
|eos_heredoc|Confusing heredoc styles and structures.|[Link](docs/rules/eos_heredoc.md)|
|eos_hungarian|Use of Hungarian notation in variable and block names.|[Link](docs/rules/eos_hungarian.md)|
|eos_meta|Problematic meta-argument syntax and values.|[Link](docs/rules/eos_meta.md)|
|eos_naming|Awkward naming conventions.|[Link](docs/rules/eos_naming.md)|
|eos_reminder|Use of reminder tags.|[Link](docs/rules/eos_reminder.md)|

## Installation

### Pre-built binary

1. Download the zip file for your platform from the [Releases](https://github.com/tfctl/tflint-ruleset-elements-of-style/releases/latest) page.

2. Unzip it to your `${HOME}/.tflint.d/plugins` folder.

### Building the plugin from source

Building from source requires Go 1.25+.

1. Clone the repository locally:

```
$ git clone https://github.com/tfctl/tflint-ruleset-elements-of-style.git
$ cd tflint-ruleset-elements-of-style
```

2. Build the binary:
```
$ make
```

3. Install the plugin binary with the following:

```
$ make install
```

## Requirements

- TFLint v0.46+

## Configuration

The plugin can be enabled with `tflint --init` after declaring the plugin in `.tflint.hcl` as follows:

```hcl
plugin "elements-of-style" {
  enabled = true

  source  = "github.com/tfctl/tflint-ruleset-elements-of-style"
  version = "1.0.0" # replace as needed
}
```

## AI Acknowledgment

This project uses AI-assisted tools (mostly GitHub CoPilot w/Claude Opus and Gemini 3) selectively:

- **AI-assisted** — The test scaffolding, and much of the documentation were created with AI assistance and reviewed before inclusion.

- **Routine refactoring** — AI tools assisted with lint corrections and minor optimizations.
