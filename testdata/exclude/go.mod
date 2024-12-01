module github.com/ldez/gomoddirectives/testdata/exclude

go 1.16

require (
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/ldez/grignotin v0.4.1
)

exclude (
    golang.org/x/crypto v1.4.5
    golang.org/x/text v1.6.7
)
