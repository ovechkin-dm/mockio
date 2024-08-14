# Limitations

## Architecture

Because library uses assembly code to generate mocks, it is not possible to use it in pure Go code. 
This means that you cannot use it in a project that is intended to be cross-compiled to multiple platforms.

For now supported platforms are:
* AMD64
* ARM64

This list may be extended in the future.

## Backwards compatibility and new Go versions

This library is tested for GO 1.18 up to  