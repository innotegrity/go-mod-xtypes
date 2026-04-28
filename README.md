<div align="center">
  <img src="./assets/logo.svg" width=150 alt="go-mod-xtypes logo" /><br />
  <h1 align="center">go-mod-xtypes</h1>
  <p align="center">
    A module that provides extra types for Go applications and modules.
  </p>
</div>
<hr />
<div align="center">
  <p>
    <b><em>Built using</em></b><br />
    <a href="https://cursor.com/" target="_blank"><img src="https://img.shields.io/badge/cursor-000000?style=for-the-badge&logo=cursor&logoColor=white" alt="Cursor" /></a>
    <a href="https://go.dev" target="_blank"><img src="https://img.shields.io/badge/Go-02A8EF?style=for-the-badge&logo=go&logoColor=white" alt="Go" /></a>
  </p>
  <p>
    <b><em>Supported on</em></b><br />
    <img src="https://img.shields.io/badge/Linux-yellow?style=for-the-badge&logo=linux&logoColor=black" alt="Linux" />
    <img src="https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=apple&logoColor=white" alt="MacOS" />
    <img src="https://img.shields.io/badge/windows-0078D7?style=for-the-badge&logo=windows&logoColor=white" alt="Windows" />
  </p>
  <p>
    <img src="https://img.shields.io/badge/stability-alpha-red?style=for-the-badge" alt="alpha Stability" />
    <a href="https://en.wikipedia.org/wiki/MIT_License" target="_blank"><img src="https://img.shields.io/badge/license-MIT-blue?style=for-the-badge" alt="MIT License" /></a>
    <img src="https://img.shields.io/badge/support-community-darkgreen?style=for-the-badge" alt="Community Supported" />
    <br />
    <a href="https://pkg.go.dev/go.innotegrity.dev/mod/xtypes" target="_blank"><img src="https://img.shields.io/badge/go-reference-2a7d98?style=for-the-badge" alt="Go Reference" /></a>
    <a href="https://goreportcard.com/report/go.innotegrity.dev/mod/xtypes" target="_blank"><img src="https://goreportcard.com/badge/go.innotegrity.dev/mod/xtypes?style=for-the-badge" alt="Go Report Card" /></a>
    <a href="https://sonarcloud.io/project/overview?id=go-mod-xtypes" target="_blank"><img src="https://img.shields.io/sonar/coverage/go-mod-xtypes?server=https://sonarcloud.io&style=for-the-badge&label=Code%20Coverage" alt="Code Coverage" /></a>
    <a href="https://sonarcloud.io/project/overview?id=go-mod-xtypes" target="_blank"><img src="https://img.shields.io/sonar/quality_gate/go-mod-xtypes?server=https://sonarcloud.io&style=for-the-badge&label=Code%20Quality" alt="Code Quality" /></a>
  </p>
</div>

<!-- omit in toc -->
## Table of Contents

- [👁️ Overview](#️-overview)
- [✅ Requirements](#-requirements)
- [👨‍💻 Developer Notes](#-developer-notes)
- [📃 License](#-license)
- [❓ Questions, Issues and Feature Requests](#-questions-issues-and-feature-requests)

## 👁️ Overview

`go-mod-xtypes` is a small Go library of extra types you can drop into applications and libraries when you need friendlier parsing and serialization than raw primitives: extended durations (months, weeks, days, years), byte sizes with SI and binary suffixes, POSIX user and group identifiers that round-trip through JSON and resolve to names, file modes as octal in config, a LocalPath helper for mkdir/open/write with optional chmod/chown and structured errors, generic sets, UUID generation, and a few helpers like converting slices to `[]any`.

It is intended for configuration files, APIs, and services that touch the filesystem or need human-readable quantities in structured data.

Please review the [project documentation](https://pkg.go.dev/go.innotegrity.dev/mod/xtypes) for additional details, examples, and API reference.

## ✅ Requirements

This module is supported for Go v1.25.9 and later running on Linux, MacOS and Windows operating systems.

## 👨‍💻 Developer Notes

For consistency, security and best practices, the maintainers of this repository utilize the following toolset:

- [Cursor IDE](https://cursor.com/product) with the following extensions:
  - [Go](https://marketplace.cursorapi.com/items/?itemName=golang.Go)
  - [Markdown All in One](https://marketplace.cursorapi.com/items/?itemName=yzhang.markdown-all-in-one)
  - [markdownlint](https://marketplace.cursorapi.com/items/?itemName=DavidAnson.vscode-markdownlint)
  - [SonarQube for IDE](https://marketplace.cursorapi.com/items/?itemName=SonarSource.sonarlint-vscode)
  - [YAML](https://marketplace.cursorapi.com/items/?itemName=redhat.vscode-yaml)
- [golangci-lint](https://github.com/golangci/golangci-lint) for `.go` files
- [markdownlint](https://github.com/davidanson/markdownlint) for `.md` files
- [pre-commit](https://pre-commit.com/) for checks prior to commits
- [SonarQube CLI](https://github.com/SonarSource/sonarqube-cli) for secrets checking during pre-commit

## 📃 License

This module is distributed under the MIT License.

## ❓ Questions, Issues and Feature Requests

If you have questions about this project, find a bug or wish to submit a feature request, please [submit an issue](https://github.com/innotegrity/go-mod-xtypes/issues).
