module github.com/ldez/gomoddirectives/testdata/ignore_defaults

go 1.24

ignore (
	./.git
	./_data
	./vendor
	./dir/testdata/
)
