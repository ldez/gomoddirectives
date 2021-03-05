module github.com/ldez/gomoddirectives/testdata/c

go 1.16

require (
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/ldez/grignotin v0.4.1
)

retract (
    v1.0.0 // foo
    [v1.0.0, v1.9.9]
)
