# implgen

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
implgen basic --src example/in/interface.go \
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

## Inspiration and References

- [gomock](https://github.com/golang/mock) - mock generation