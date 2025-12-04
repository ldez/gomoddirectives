package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
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
	IgnoreForbidden           bool
	RetractAllowNoExplanation bool
	ToolchainForbidden        bool
	ToolForbidden             bool
	GoDebugForbidden          bool
	GoVersionPattern          string
	ToolchainPattern          string
	CheckModulePath           bool
}

func main() {
	cfg := config{}

	flag.BoolVar(&cfg.ExcludeForbidden, "exclude", false, "Forbid the use of exclude directives")
	flag.BoolVar(&cfg.IgnoreForbidden, "ignore", false, "Forbid the use of ignore directives")
	flag.Var(&cfg.ReplaceAllowList, "list", "List of allowed replace directives")
	flag.BoolVar(&cfg.ReplaceAllowLocal, "local", false, "Allow local replace directives")
	flag.BoolVar(&cfg.RetractAllowNoExplanation, "retract-no-explanation", false, "Allow to use retract directives without explanation")
	flag.BoolVar(&cfg.ToolchainForbidden, "toolchain", false, "Forbid the use of toolchain directive")
	flag.StringVar(&cfg.ToolchainPattern, "toolchain-pattern", "", "Pattern to validate toolchain directive")
	flag.BoolVar(&cfg.ToolForbidden, "tool", false, "Forbid the use of tool directives")
	flag.BoolVar(&cfg.GoDebugForbidden, "godebug", false, "Forbid the use of godebug directives")
	flag.StringVar(&cfg.GoVersionPattern, "goversion", "", "Pattern to validate go min version directive")
	flag.BoolVar(&cfg.CheckModulePath, "check-module-path", false, "Check module path validity")

	help := flag.Bool("h", false, "Show this help.")

	flag.Usage = usage

	flag.Parse()

	if *help {
		usage()
	}

	opts := gomoddirectives.Options{
		ReplaceAllowList:          cfg.ReplaceAllowList,
		ReplaceAllowLocal:         cfg.ReplaceAllowLocal,
		ExcludeForbidden:          cfg.ExcludeForbidden,
		IgnoreForbidden:           cfg.IgnoreForbidden,
		RetractAllowNoExplanation: cfg.RetractAllowNoExplanation,
		ToolchainForbidden:        cfg.ToolchainForbidden,
		ToolForbidden:             cfg.ToolForbidden,
		GoDebugForbidden:          cfg.GoDebugForbidden,
		CheckModulePath:           cfg.CheckModulePath,
	}

	if cfg.GoVersionPattern != "" {
		var err error

		opts.GoVersionPattern, err = regexp.Compile(cfg.GoVersionPattern)
		if err != nil {
			log.Fatal(err)
		}
	}

	if cfg.ToolchainPattern != "" {
		var err error

		opts.ToolchainPattern, err = regexp.Compile(cfg.ToolchainPattern)
		if err != nil {
			log.Fatal(err)
		}
	}

	results, err := gomoddirectives.Analyze(opts)
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
