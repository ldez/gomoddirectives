module github.com/ldez/gomoddirectives/testdata/ignore_defaults

go 1.24

ignore (
	.git
	_data
	vendor
	dir/testdata/
	dir/vendor/
)

ignore (
	./.git
	./_data
	./vendor
	./dir/testdata/
	./dir/vendor/
)

ignore (
	.
	./foo.foo
	./foo_foo
	./testdata_foo
	./foo_testdata
	./vendor_foo
	./foo_vendor
)
