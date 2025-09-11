[简体中文](README_CN.md) | **[English](README.md)**

---

![GitHub License](https://img.shields.io/github/license/yhlooo/gcovgo)
[![GitHub Release](https://img.shields.io/github/v/release/yhlooo/gcovgo)](https://github.com/yhlooo/gcovgo/releases/latest)
[![release](https://github.com/yhlooo/gcovgo/actions/workflows/release.yaml/badge.svg)](https://github.com/yhlooo/gcovgo/actions/workflows/release.yaml)

# gcovgo

该项目是一个纯 go 实现的 [gcov](https://gcc.gnu.org/onlinedocs/gcc/Gcov.html) 覆盖率数据解析工具。二进制工具 `gcovgo` 具有与 `gcov` 、 `gcov-dump` 和 `gcov-tool` 命令具有相似的功能。同时该项目可作为 go 库被集成。

## 安装

### Docker

直接使用镜像 [`ghcr.io/yhlooo/gcovgo`](https://github.com/yhlooo/gcovgo/pkgs/container/gcovgo) docker run 即可：

```bash
docker run -it --rm ghcr.io/yhlooo/gcovgo:latest --help
```

### 通过二进制安装

通过 [Releases](https://github.com/yhlooo/gcovgo/releases) 页面下载可执行二进制，解压并将其中 `gcovgo` 文件放置到任意 `$PATH` 目录下。

### 从源码编译

要求 Go 1.24 ，执行以下命令下载源码并构建：

```bash
go install github.com/yhlooo/gcovgo/cmd/gcovgo@latest
```

构建的二进制默认将在 `${GOPATH}/bin` 目录下，需要确保该目录包含在 `$PATH` 中。

## 使用

### 解析覆盖率数据

与 `gcov` 命令作用类似。输入 gcov 插桩编译后生成的 `.gcno` 文件和插桩编译的程序运行时产生的 `.gcda` 文件，输出解析计算后的覆盖率信息。

```bash
gcovgo path/to/file.gcno
```

### 查看覆盖率数据内容

与 `gcov-dump` 命令作用类似。输入 gcov 插桩编译后生成的 `.gcno` 文件或插桩编译的程序运行时产生的 `.gcda` 文件，以 JSON 等易于处理或人类可读的形式输出该文件内容。

```bash
gcovgo dump path/to/file.gcno
```
