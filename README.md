<div style="text-align: center; padding-bottom: 20px">
<img src="logo_light.png#gh-dark-mode-only" style="width: 300px; margin-left: auto; margin-right: auto;">
<img src="logo.png#gh-light-mode-only" style="width: 300px; margin-left: auto; margin-right: auto;">
</div>

# Depbot

Depbot is a CLI tool that identifies dependencies in a given source code repository. It aims to support multiple dependency management systems and multiple languages.

## Installation

The first option to install the Depbot agent is to download a precompiled binary from [GitHub](github.com/godepbot/depbot/releases).

### GNU / Linux
```sh
$ wget https://github.com/godepbot/depbot/releases/latest/download/depbot_Linux_x86_64.tar.gz
$ tar -xvzf depbot_Linux_x86_64.tar.gz
$ sudo mv depbot /usr/local/bin/depbot
```

### MacOS
```sh
$ curl -OL https://github.com/godepbot/depbot/releases/latest/download/depbot_Darwin_x86_64.tar.gz
$ tar -xvzf depbot_Darwin_x86_64.tar.gz
$ sudo mv depbot /usr/local/bin/depbot

# or if you have ~/bin folder setup in the environment PATH variable
$ mv depbot ~/bin/depbot
```

Also, if you have the Go toolkit installed you can run.

```bash
$ go install github.com/godepbot/depbot/cmd/depbot@latest
```

Or download the appropriate binary from the releases page.

## Usage

Depbot CLI is in charge of determining dependencies on the current source code folder, It has the list and the sync command.

### List

List analyzes the current source code folder and prints a list of dependencies it has.

```sh
> $ depbot list

Total dependencies found: 3

Name                    Version                                 File    Direct
----                    -------                                 ----    -------
Go                      1.18                                    go.mod  false
golang.org/x/mod        v0.5.1                                  go.mod  false
golang.org/x/xerrors    v0.0.0-20191011141410-1b5146add898      go.mod  false
```

## sync

The sync command synchronizes the dependencies to a given Depbot server. The command looks for the DEPBOT_API_KEY DEPBOT_SERVER_URL environment variables when synchronizing. If the DEPBOT_API_KEY is not set the command errors.

```sh
$ depbot sync

Success! 34 Dependencies synchronized.
```

## Running in development

Assuming you have Go installed in your machine you can invoke the Depbot command by running:

```bash
$ go run ./cmd/depbot
```

