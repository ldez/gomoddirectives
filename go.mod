module github.com/ldez/gomoddirectives

go 1.23.0

require (
	github.com/ldez/grignotin v0.9.0
	github.com/stretchr/testify v1.10.0
	golang.org/x/mod v0.25.0
	golang.org/x/tools v0.32.0
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
