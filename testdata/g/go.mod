module github.com/ldez/gomoddirectives/testdata/g

go 1.16

require (
	github.com/gorilla/mux v1.7.3
	github.com/ldez/grignotin v0.4.1
)

replace github.com/ldez/grignotin => ../b/

replace (
	github.com/gorilla/mux => example.com/gorilla/mux v1.7.3
	github.com/gorilla/mux v1.2.3 => example.com/gorilla/mux v1.4.5
)
