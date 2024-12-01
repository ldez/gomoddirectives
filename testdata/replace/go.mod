module github.com/ldez/gomoddirectives/testdata/replace

go 1.16

require (
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/ldez/grignotin v0.4.1
)

replace (
	github.com/gorilla/mux => github.com/containous/mux v0.0.0-20181024131434-c33f32e26898
	github.com/ldez/grignotin => ../b/
)
