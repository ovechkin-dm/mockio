# Limitations

## Architecture

Because library uses assembly code to generate mocks, it is not possible to use it in pure Go code. 
This means that you cannot use it in a project that is intended to be cross-compiled to multiple platforms.

For now supported platforms are:

- AMD64
- ARM64

This list may be extended in the future.

## Backwards compatibility and new Go versions

This library is tested for GO 1.18 up to 1.23

Caution: there is no guarantee that it will work with future versions of Go. 
However there is not much that can break the library, so it should be easy to fix it if it stops working. As of latest mockio version, almost all of dependencies on golang internal runtime features were removed.

 Please refer to [go-dyno documentation](https://ovechkin-dm.github.io/go-dyno/latest/) for more information on compatibility.