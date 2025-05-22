# Code Generation

While Mockio is designed to work without code generation, it also provides a code generation tool that can help with creating mock implementations for your interfaces. This can be useful in cases where you want to have more control over the mock implementation, or don't want to rely on an unsafe features mockio uses.

## Installation

Since mockio relies on mockery generation tool, mockery should be installed.

```sh
go install github.com/vektra/mockery/v3@v3.2.5
```

Please refer to official mockery [installation guide](https://vektra.github.io/mockery/latest/installation/)

## Configuration

To configure mockery tool, provide a link to mockio tempate:

```yaml
all: false
dir: '{{.InterfaceDir}}'
filename: mocks_test.go
force-file-write: true
formatter: goimports
log-level: info
structname: '{{.Mock}}{{.InterfaceName}}'
pkgname: '{{.SrcPackageName}}'
recursive: false
require-template-schema-exists: true
template: https://raw.githubusercontent.com/ovechkin-dm/mockio/refs/heads/main/templates/mockery.tmpl
template-schema: '{{.Template}}.schema.json'
packages:
  your-package-name:
    config:
      all: true
```

The main difference here is the `template` parameter.

Please refer to official mockery [configuration guide](https://vektra.github.io/mockery/latest/configuration/) for more details.

## Code generation

Now that mockery is installed and configured, we can run mockery command in project root:

```sh
mockery
```

It will generate necessary mocks

## Usage

To use the generated mock, you can create them via generated contructor call `New{YourMockName}`.
Here is an example:

```
func TestSimple(t *testing.T) {
    ctrl := NewMockController(t)
    m := NewMockUserService(ctrl)
}
```

Other than that, API stays the same as if runtime mocks were used. 

## Full example

You can check full example [here](https://github.com/ovechkin-dm/mockio-codegen-example)





