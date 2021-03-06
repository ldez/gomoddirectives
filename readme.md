# gomoddirectives

A linter that handle [`replace`](https://golang.org/ref/mod#go-mod-file-replace), [`retract`](https://golang.org/ref/mod#go-mod-file-retract), [`exclude`](https://golang.org/ref/mod#go-mod-file-exclude) directives into `go.mod`.

Features:

- ban all [`replace`](https://golang.org/ref/mod#go-mod-file-replace) directives
- allow only local [`replace`](https://golang.org/ref/mod#go-mod-file-replace) directives
- allow only some [`replace`](https://golang.org/ref/mod#go-mod-file-replace) directives
- force explanation for [`retract`](https://golang.org/ref/mod#go-mod-file-retract) directives
- ban all [`exclude`](https://golang.org/ref/mod#go-mod-file-exclude) directives
- detect duplicated [`replace`](https://golang.org/ref/mod#go-mod-file-replace) directives
- detect identical [`replace`](https://golang.org/ref/mod#go-mod-file-replace) directives
