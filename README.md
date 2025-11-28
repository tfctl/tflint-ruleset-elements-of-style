# TFLint Ruleset: Elements of Style

This is a custom TFLint ruleset for both idiomatic conventions and opinionated styles and consistent baseline for Terraform code. Much like a style guide for writing, these rules are prescriptive by design. Style isn’t just aesthetics — it makes code easier to read, review, and change safely. Shared conventions reduce friction and help teams move faster with fewer surprises.

## Rules

|Name|Description|Link|
| --- | --- | --- |
|eos_comments|Identify non-standard comment styles.|[Link](docs/rules/eos_comments.md)|
|eos_death_mask|Identify commented-out Terraform code.|[Link](docs/rules/eos_death_mask.md)|
|eos_dry|Identify repeated values (strings, interpolations, lists, maps).|[Link](docs/rules/eos_dry.md)|
|eos_heredoc|Identify standard heredoc syntax and EOF delimiter.|[Link](docs/rules/eos_heredoc.md)|
|eos_hungarian|Identify Hungarian notation in names.|[Link](docs/rules/eos_hungarian.md)|
|eos_meta|Enforce Terraform meta-argument syntax conventions.|[Link](docs/rules/eos_meta.md)|
|eos_naming|Enforce naming conventions (length, casing).|[Link](docs/rules/eos_naming.md)|
|eos_reminder|Identify comments containing reminder tags.|[Link](docs/rules/eos_reminder.md)|
|eos_type_echo|Identify type echoing in names.|[Link](docs/rules/eos_type_echo.md)|

## Installation

### Pre-built binary

1. Download the zip file for your platform from [Release](https://github.com/staranto/tflint-ruleset-elements-of-style/releases/latest).

2. Unzip it to your `${HOME}/.tflint.d/plugins` folder.

### Building the plugin from source

Building from source requires Go 1.25+.

1. Clone the repository locally and then build the binary:

```
$ make
```

2. Install the plugin binary with the following:

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

  version = "1.2.0"
  source  = "github.com/staranto/tflint-ruleset-elements-of-style"
}
```
