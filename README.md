# ajson
--
    import "."

Package ajson provides a way to marshal and unmarshal JSON with unknown fields.

[![Build & Test Action
Status](https://github.com/devnw/ajson/actions/workflows/build.yml/badge.svg)](https://github.com/devnw/ajson/actions)
[![Go Report
Card](https://goreportcard.com/badge/go.devnw.com/ajson)](https://goreportcard.com/report/go.devnw.com/ajson)
[![codecov](https://codecov.io/gh/devnw/ajson/branch/main/graph/badge.svg)](https://codecov.io/gh/devnw/ajson)
[![Go
Reference](https://pkg.go.dev/badge/go.devnw.com/ajson.svg)](https://pkg.go.dev/go.devnw.com/ajson)
[![License: Apache
2.0](https://img.shields.io/badge/license-Apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![PRs
Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

## Installation

```bash

go get -u go.devnw.com/ajson

```

## Usage

#### func  Marshal

```go
func Marshal[T comparable](t T, mm MMap) ([]byte, error)
```

#### type MMap

```go
type MMap tcontainer.MarshalMap
```


#### func  Unmarshal

```go
func Unmarshal[T comparable](data []byte) (T, MMap, error)
```
