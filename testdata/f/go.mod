module github.com/ldez/gomoddirectives/testdata/f

go 1.16

require (
	github.com/gorilla/mux v1.7.3
	github.com/ldez/grignotin v0.4.1
)

replace (
	github.com/gorilla/mux v1.7.3 => github.com/gorilla/mux v1.7.3
	github.com/ldez/grignotin => ../b
)
