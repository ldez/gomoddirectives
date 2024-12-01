module github.com/ldez/gomoddirectives/testdata/tool_multiple

go 1.24

tool golang.org/x/tools/cmd/stringer

tool (
    example.com/module/cmd/a
    example.com/module/cmd/b
)
