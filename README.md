[简体中文](README_CN.md) | **[English](README.md)**

---

![GitHub License](https://img.shields.io/github/license/yhlooo/gcovgo)
[![GitHub Release](https://img.shields.io/github/v/release/yhlooo/gcovgo)](https://github.com/yhlooo/gcovgo/releases/latest)
[![release](https://github.com/yhlooo/gcovgo/actions/workflows/release.yaml/badge.svg)](https://github.com/yhlooo/gcovgo/actions/workflows/release.yaml)

# gcovgo

This project is a pure Go implementation of a [gcov](https://gcc.gnu.org/onlinedocs/gcc/Gcov.html) coverage data parsing tool. The binary tool `gcovgo` provides functionality similar to the `gcov`, `gcov-dump` and `gcov-tool` commands. Additionally, the project can be integrated as a Go library.

## Installation

### Docker

docker run with image [`ghcr.io/yhlooo/gcovgo`](https://github.com/yhlooo/gcovgo/pkgs/container/gcovgo):

```bash
docker run -it --rm ghcr.io/yhlooo/gcovgo:latest --help
```

### Binaries

Download the executable binary from the [Releases](https://github.com/yhlooo/gcovgo/releases) page, extract it, and place the `gcovgo` file into any `$PATH` directory.

### From Sources

Requires Go 1.24. Execute the following command to download the source code and build it:

```bash
go install github.com/yhlooo/gcovgo/cmd/gcovgo@latest
```

The built binary will be located in `${GOPATH}/bin` by default. Make sure this directory is included in your `$PATH`.

## Usage

### Parsing Coverage Data

Similar to the `gcov` command, this function takes `.gcno` and `.gcda` files, then outputs parsed coverage information.

```bash
gcovgo path/to/file.gcno
```

### Print Coverage Data Content

Similar to the `gcov-dump` command, this function accepts either `.gcno` or `.gcda` files. It outputs the file content in a human-readable or easily processable format (e.g. JSON).

```bash
gcovgo dump path/to/file.gcno
```
