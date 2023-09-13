*A systemd inspired system configuration and deployment tool*

[![Go](https://github.com/khulnasoft-lab/system-deploy/workflows/Go/badge.svg)](https://github.com/khulnasoft-lab/system-deploy/actions)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/khulnasoft-lab/system-deploy?style=flat-square)
![Stable: Beta](https://img.shields.io/badge/Stable-BETA-yellowgreen?style=flat-square)
![GitHub](https://img.shields.io/github/license/khulnasoft-lab/system-deploy?style=flat-square)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/khulnasoft-lab/system-deploy?label=Release&style=flat-square)

`system-deploy` is my personal server management and deployment tool. It's inspired by systemd's unit files and the deployment tool/script used @safing. It currently supports copying files/directories, installing packages and installing/enabling systemd unit files. `system-deploy` is meant to be executed periodically and supports running different actions and tasks when changes are detected.

The compiled binary itself includes help and documentation for almost all supported operations and even some examples.

**Checkout [system-conf](https://github.com/khulnasoft-lab/system-conf) for a systemd inspired configuration system for Go projects.**

## Getting Started 

Please checkout the documentation at [khulnasoft-lab.github.io/system-deploy](https://khulnasoft-lab.github.io/system-deploy).

## Installation from Source

`system-deploy` is written in <img src="https://golang.org/lib/godoc/images/go-logo-blue.svg" alt="Go" height="13"> and thus requires a working go installation to be compiled from source.

```bash
go install github.com/khulnasoft-lab/system-deploy/cmd/system-deploy
```
This will install the `deploy` command into your `$GOBIN` or `$GOPATH/bin` if the former is not set.


## Contributing

Any contributions to the `system-deploy` project are welcome! Just fork the repository and create a PR with your changes. It's recommended to discuss planned changes in an [issue](https://github.com/khulnasoft-lab/system-deploy/issues) first.

## License

`system-deploy` itself is licensed under a BSD 3-clause license. See [LICENSE](LICENSE) for more information.

Note that the bineries distributed via the release page are licensed under GPL-3 because the `EditFile` action is compiled against a GPL-3 licensed library.
