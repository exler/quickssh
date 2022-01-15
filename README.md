<p align="center">
    <img src="logo.png" width="168">
    <p align="center">🖥 SSH/SFTP profile manager and client</p>
    <p align="center">
      <img src="https://img.shields.io/github/v/release/exler/quickssh?label=Release">
      <img src="https://github.com/exler/quickssh/actions/workflows/tests.yml/badge.svg">
      <img src="https://img.shields.io/github/go-mod/go-version/exler/quickssh">
      <img alt="MIT License" src="https://img.shields.io/github/license/exler/quickssh?color=lightblue">
    </p>
</p>

## Overview

QuickSSH manages your SSH server profiles and provides utility functions for performing common operations.

QuickSSH supports password authentication, key-based authentication and SSH agent forwarding (Unix only). 

## Installation

### Go

**Requires**: Go >= 1.17

You can use Go to compile and install the package to your `$GOBIN` directory.

```bash
$ go install github.com/exler/quickssh
```

### Releases

Binaries for Windows, Linux and MacOS systems can be found on [GitHub releases](https://github.com/exler/quickssh/releases).

## Usage

```bash
$ quickssh -h

QuickSSH allows for easy management of SSH profiles and simplifies
working with the SSH and SFTP protocols.

Usage:
  quickssh [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  connect     Connect to a SSH server
  exec        Execute a command on a SSH server
  file        Download or upload files using SFTP
  help        Help about any command
  profile     Manage server profiles
  version     Show current program version

Flags:
  -h, --help   help for quickssh

Use "quickssh [command] --help" for more information about a command.
```

Configuration files are saved in the user directory:
* Windows: `%APPDATA%\quickssh`
* Unix: `~/.config/quickssh`

## License

`QuickSSH` is under the terms of the [MIT License](https://www.tldrlegal.com/l/mit), following all clarifications stated in the [license file](LICENSE).

