# gomoddirectives

A linter that handle `replace`, `retract`, `exclude` directives into `go.mod`.

Features:

- ban all `replace` directives
- allow only local `replace` directives
- allow only some `replace` directives
- force explanation for `retract` directives
- forbid `exclude` directives
