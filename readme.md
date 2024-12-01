# gomoddirectives

[![Sponsor](https://img.shields.io/badge/Sponsor%20me-%E2%9D%A4%EF%B8%8F-pink)](https://github.com/sponsors/ldez)
[![Build Status](https://github.com/ldez/gomoddirectives/workflows/Main/badge.svg?branch=master)](https://github.com/ldez/gomoddirectives/actions)

A linter that handle directives into `go.mod`.

## [`retract`](https://golang.org/ref/mod#go-mod-file-retract) directives

- Force explanation for `retract` directives.

```go
module example.com/foo

go 1.22

require (
	github.com/ldez/grignotin v0.4.1
)

retract (
    v1.0.0 // Explanation
)
```

## [`replace`](https://golang.org/ref/mod#go-mod-file-replace) directives

- Ban all `replace` directives.
- Allow only local `replace` directives.
- Allow only some `replace` directives.
- Detect duplicated `replace` directives.
- Detect identical `replace` directives.

```go
module example.com/foo

go 1.22

require (
	github.com/ldez/grignotin v0.4.1
)

replace github.com/ldez/grignotin => ../grignotin/
```

## [`exclude`](https://golang.org/ref/mod#go-mod-file-exclude) directives

- Ban all `exclude` directives.

```go
module example.com/foo

go 1.22

require (
	github.com/ldez/grignotin v0.4.1
)

exclude (
    golang.org/x/crypto v1.4.5
    golang.org/x/text v1.6.7
)
```

## [`tool`](https://golang.org/ref/mod#go-mod-file-tool) directives

- Ban all `tool` directives.

```go
module example.com/foo

go 1.24

tool (
    example.com/module/cmd/a
    example.com/module/cmd/b
)
```

## [`toolchain`](https://golang.org/ref/mod#go-mod-file-toolchain) directive

- Ban `toolchain` directive.

```go
module example.com/foo

go 1.22

toolchain go1.23.3
```

## [`godebug`](https://go.dev/ref/mod#go-mod-file-godebug) directives

- Ban `godebug` directive.

```go
module example.com/foo

go 1.22

godebug default=go1.21
godebug (
    panicnil=1
    asynctimerchan=0
)
```

## [`go`](https://go.dev/ref/mod#go-mod-file-go) directive

- Use a regular expression to constraint the Go minimum version.

```go
module example.com/foo

go 1.22.0
```
