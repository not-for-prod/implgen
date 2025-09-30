# implgen

![Go](https://img.shields.io/badge/Go-1.25.1-blue)
![License](https://img.shields.io/github/license/not-for-prod/implgen)

`implgen` is a developer utility designed to **automatically generate boilerplate implementations** for Go interfaces. 
It helps kick-start the development of service layers, adapters, and stubs while saving time and reducing manual effort.

---
## Overview

`implgen` provides some fetures:
- 🛠 Generate empty implementations of interfaces
- 📂 Supports single-file or per-method file output
- 📦 Customizable output package and struct name
- 🎯 Target a specific interface or process all in the source file
- 🧭 Supports OpenTelemetry span instrumentation for methods with `context.Context`
- 🐫 Automatic file/folder naming via `kebab-case` and `snake_case` converters

---

## 🚀 Installation

```bash
go install github.com/not-for-prod/implgen@latest
```

## Usage example 

```shell
implgen --src=./service --dst=./serviceimpl --interface-name=Greeter
```

Flags (required):

- `src` - source file path
- `dst` - destination dir path

Flags (optional):

- `interface-name` - source `interface` name
- `impl-name` - generated implementation `struct` name
- `impl-package` - generated implementation `package` name, can be used only if `interface-name` set
- `enable-trace` - enables writing `otel.Traсer(...).Start(...)` in methods, 
where first argument type is `context.Context` 
- `tracer-name` - name used in `otel.Traсer(<tracer-name>)`
- `single-file` - indicates whether methods will be generated into single file or
file per method

Assume you have an [interface](./example/in/interface.go):

```go
type TestInterface interface {
    A(ctx context.Context, req dto.GoRequest) error
    B(ctx context.Context, req map[dto.GoRequest]dto.GoRequest) error
    C(ctx context.Context, req []dto.GoRequest) error
    D(ctx context.Context, req int, opts ...dto.GoRequest) error
}
```

Run:

```shell
implgen --src example/in/interface.go \
		--dst example/out/ \
		--interface-name TestInterface \
		--impl-name Test \
		--impl-package test \
		--enable-trace \
		--tracer-name my-brilliant-tracer
```

This will generate:

- A struct with name `Test` with method stubs
- Traced methods using OpenTelemetry if `context.Context` is the first method param

See [dst example](example/out) for more details
