package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ldez/gomoddirectives"
)

//nolint:recvcheck // required for the marshaling.
type flagSlice []string

func (f flagSlice) String() string {
	return strings.Join(f, ":")
}

func (f *flagSlice) Set(s string) error {
	*f = append(*f, strings.Split(s, ",")...)
	return nil
}

type config struct {
	ReplaceAllowList          flagSlice
	ReplaceAllowLocal         bool
	ExcludeForbidden          bool
	RetractAllowNoExplanation bool
	ToolchainForbidden        bool
	ToolForbidden             bool
}

func main() {
	cfg := config{}

	flag.BoolVar(&cfg.ExcludeForbidden, "exclude", false, "Forbid the use of exclude directives")
	flag.Var(&cfg.ReplaceAllowList, "list", "List of allowed replace directives")
	flag.BoolVar(&cfg.ReplaceAllowLocal, "local", false, "Allow local replace directives")
	flag.BoolVar(&cfg.RetractAllowNoExplanation, "retract-no-explanation", false, "Allow to use retract directives without explanation")
	flag.BoolVar(&cfg.ToolchainForbidden, "toolchain", false, "Forbid the use of toolchain directive")
	flag.BoolVar(&cfg.ToolForbidden, "tool", false, "Forbid the use of tool directives")

	help := flag.Bool("h", false, "Show this help.")

	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
	}

	results, err := gomoddirectives.Analyze(gomoddirectives.Options{
		ExcludeForbidden:          cfg.ExcludeForbidden,
		ReplaceAllowList:          cfg.ReplaceAllowList,
		ReplaceAllowLocal:         cfg.ReplaceAllowLocal,
		RetractAllowNoExplanation: cfg.RetractAllowNoExplanation,
		ToolchainForbidden:        cfg.ToolchainForbidden,
		ToolForbidden:             cfg.ToolForbidden,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range results {
		fmt.Println(e)
	}

	if len(results) > 0 {
		os.Exit(1)
	}
}

func usage() {
	_, _ = os.Stderr.WriteString(`GoModDirectives

gomoddirectives [flags]

Flags:
`)
	flag.PrintDefaults()
	os.Exit(2)
}
