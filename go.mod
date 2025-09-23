module github.com/ldez/gomoddirectives

go 1.24.0

require (
	github.com/ldez/grignotin v0.10.1
	github.com/stretchr/testify v1.11.1
	golang.org/x/mod v0.28.0
	golang.org/x/tools v0.37.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v0.4.0 // invalid implementation of AnalyzePass
	v0.3.0 // invalid functions
)
