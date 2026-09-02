module github.com/ldez/gomoddirectives

go 1.25.0

require (
	github.com/ldez/grignotin v0.10.1
	github.com/stretchr/testify v1.12.1
	golang.org/x/mod v0.40.0
	golang.org/x/tools v0.49.0
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect

retract (
	v0.4.0 // invalid implementation of AnalyzePass
	v0.3.0 // invalid functions
)
