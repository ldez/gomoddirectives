module github.com/ldez/gomoddirectives/testdata/godebug

go 1.22

godebug default=go1.21
godebug (
    panicnil=1
    asynctimerchan=0
)
