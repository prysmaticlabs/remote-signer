# Eth2 Remote Signer

![Go](https://github.com/prysmaticlabs/remote-signer/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/prysmaticlabs/remote-signer)](https://goreportcard.com/report/github.com/prysmaticlabs/remote-signer)
[![GoDoc](https://godoc.org/github.com/prysmaticlabs/remote-signer?status.svg)](https://godoc.org/github.com/prysmaticlabs/remote-signer)


Remote signing gRPC server reference implementation for the Go Ethereum 2.0 client [prysmaticlabs/prysm](https://github.com/prysmaticlabs/prysm)

## Overview

This is a simple, remote signing reference implementation to be used with the [Prysm](https://github.com/prysmaticlabs/prysm) project. It is **not** meant to be used in production deployments, but instead as an example of how to create a minimal remote-signer for eth2 validator keys in Go.

- reads and parses configuration structure from the file
- reads and overwrites configuration structure from environment variables
- writes a detailed variable list to help output

## Content

- [Installation](#installation)
- [Usage](#usage)
    - [Read Configuration](#read-configuration)
    - [Read Environment Variables Only](#read-environment-variables-only)
    - [Update Environment Variables](#update-environment-variables)
    - [Description](#description)
- [Extending the Implementation](#model-format)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install the package run

```bash
go get -u github.com/ilyakaznacheev/cleanenv
```
