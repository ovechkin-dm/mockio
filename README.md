# Mockio 

[![Build Status](https://github.com/ovechkin-dm/mockio/actions/workflows/build.yml/badge.svg)](https://github.com/ovechkin-dm/mockio/actions)
[![Codecov](https://codecov.io/gh/ovechkin-dm/mockio/branch/main/graph/badge.svg)](https://app.codecov.io/gh/ovechkin-dm/mockio)
[![Go Report Card](https://goreportcard.com/badge/github.com/ovechkin-dm/mockio)](https://goreportcard.com/report/github.com/ovechkin-dm/mockio)
[![Documentation](https://pkg.go.dev/badge/github.com/ovechkin-dm/mockio.svg)](https://pkg.go.dev/github.com/ovechkin-dm/mockio)
[![Release](https://img.shields.io/github/release/ovechkin-dm/mockio.svg)](https://github.com/ovechkin-dm/mockio/releases)
[![License](https://img.shields.io/github/license/ovechkin-dm/mockio.svg)](https://github.com/ovechkin-dm/mockio/blob/main/LICENSE)

# Mock library for golang without code generation
Mockio is a Golang library that provides functionality for mocking and stubbing functions and methods in tests inspired by mockito. The library is designed to simplify the testing process by allowing developers to easily create test doubles for their code, which can then be used to simulate different scenarios.

# Documentation

Latest documentation is available [here](https://ovechkin-dm.github.io/mockio/latest/)

# Quick start

Install latest version of the library using go get command:

```bash
go get -u github.com/ovechkin-dm/mockio
```

Create a simple mock and test:
```go
package main

import (
    . "github.com/ovechkin-dm/mockio/mock"
    "testing"
)

type Greeter interface {
    Greet(name string) string
}

func TestGreet(t *testing.T) {
    ctrl := NewMockController(t)
    m := Mock[Greeter](ctrl)
    WhenSingle(m.Greet("John")).ThenReturn("Hello, John!")
    if m.Greet("John") != "Hello, John!" {
        t.Fail()
    }
}
```