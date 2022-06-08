<div style="text-align: center; padding-bottom: 20px">
<img src="logo_light.png#gh-dark-mode-only" style="width: 300px; margin-left: auto; margin-right: auto;">
<img src="logo.png#gh-light-mode-only" style="width: 300px; margin-left: auto; margin-right: auto;">
</div>

# Depbot

Depbot is a CLI tool that identifies dependencies in a given source code repository. It aims to support multiple dependency management systems and multiple languages.

## Installation

To install the Depbot agent you can download to download a precompiled binary from [GitHub](github.com/godepbot/depbot/releases).

### MacOS

#### M1 Chip

```
$ curl -OL https://github.com/godepbot/depbot/releases/latest/download/depbot_darwin_arm64.tar.gz
$ tar -xvzf depbot_darwin_arm64.tar.gz
$ sudo mv depbot /usr/local/bin/depbot

# or if you have ~/bin folder setup in the environment PATH variable
$ mv depbot ~/bin/depbot
```

#### Intel / 64-bit

```sh
$ curl -OL https://github.com/godepbot/depbot/releases/latest/download/depbot_darwin_amd64.tar.gz
$ tar -xvzf depbot_darwin_amd64.tar.gz
$ sudo mv depbot /usr/local/bin/depbot

# or if you have ~/bin folder setup in the environment PATH variable
$ mv depbot ~/bin/depbot
```

### Linux

```
$ wget https://github.com/godepbot/depbot/releases/latest/download/depbot_linux_arm64.tar.gz
$ tar -xvzf depbot_linux_arm64.tar.gz
$ sudo mv depbot /usr/local/bin/depbot
```

### Installing from source

Alternative, if you have the Go toolkit installed you can use the Go command to install from source.

```bash
$ go install github.com/godepbot/depbot/cmd/depbot@latest
```

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

The output format of the `list` command can be modified by specifying the `--output` flag.

`--output` flag supports 3 differents options, `csv`, `json`, `plain`, by default will be set to `plain`.

```sh
# Print with csv option

> $ depbot list --output=csv

"Name","Version","File","Direct"
"Go","1.18","go.mod","true"
"golang.org/x/mod","v0.5.1","go.mod","true"
"golang.org/x/xerrors","v0.0.0-20191011141410-1b5146add898","go.mod","false"


# Print with json option

> $ depbot list --output=json

[{"File":"go.mod","Name":"Go","Version":"1.18","License":"","Kind":"language","Direct":true},{"File":"go.mod","Name":"golang.org/x/mod","Version":"v0.5.1","License":"","Kind":"library","Direct":true},{"File":"go.mod","Name":"golang.org/x/xerrors","Version":"v0.0.0-20191011141410-1b5146add898","License":"","Kind":"library","Direct":false}]
```

## Sync

The sync command POST's dependencies to the server running at the `DEPBOT_SERVER_ADDR` address. It requires the `DEPBOT_API_KEY` or `--api-key` flag to run, otherwise it errors.

Here is an example specifying the key as an environment variable:

```sh
$ DEPBOT_API_KEY=AAAA depbot sync

Success! 34 Dependencies synchronized.
```

It can also be specified by the `--api-key` flag.

```sh
# Key specified with the --api-key flag
$ depbot sync --api-key=AAAA

Success! 34 Dependencies synchronized.
```

And the command errors if the key is not specified.

```sh
# No API Key specified
$ depbot sync

Error: No API Key specified.
```

The server address can be modified by using the `--server-address` flag.

```sh
$ depbot sync --api-key=AAAA --server-address=yourserver.com

Success! 50 Dependencies synchronized.
```

The HTTPS schema will be prepend to the address if it is missing.
