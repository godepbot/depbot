<div style="text-align: center; padding-bottom: 20px">
<img src="logo_light.png#gh-dark-mode-only" style="width: 300px; margin-left: auto; margin-right: auto;">
<img src="logo.png#gh-light-mode-only" style="width: 300px; margin-left: auto; margin-right: auto;">
</div>

# Depbot

Depbot is a CLI tool that identifies dependencies in a given source code repository. It aims to support multiple dependency management systems and multiple languages.

## Installation

To install depbot, run:

```bash
$ go install github.com/dep-bot/depbot/cmd/depbot@latest
```

Or download the appropriate binary from the releases page.

## Usage

Depbot is focuses on doing the analysis of dependencies within the source code. The default behavior of the command does that. Besides it, the 

```bash
$ depbot
```

And that will list the dependencies found in current directory and bellow.

```bash
[TODO] show how results look like
```

## Options

[TODO] Document options

## Running in development

Assuming you have Go installed in your machine you can invoke the depbot command by running:

```bash
$ go run ./cmd/depbot
```

