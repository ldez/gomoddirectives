// Package gomoddirectives a linter that handle `replace`, `retract`, `exclude` directives into `go.mod`.
package gomoddirectives

import (
	"fmt"
	"go/token"
	"regexp"
	"strings"

	"github.com/ldez/grignotin/gomod"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/analysis"
)

const (
	reasonRetract          = "a comment is mandatory to explain why the version has been retracted"
	reasonExclude          = "exclude directive is not allowed"
	reasonToolchain        = "toolchain directive is not allowed"
	reasonTool             = "tool directive is not allowed"
	reasonGoDebug          = "godebug directive is not allowed"
	reasonGoVersion        = "go directive (%s) doesn't match the pattern '%s'"
	reasonReplaceLocal     = "local replacement are not allowed"
	reasonReplace          = "replacement are not allowed"
	reasonReplaceIdentical = "the original module and the replacement are identical"
	reasonReplaceDuplicate = "multiple replacement of the same module"
)

// Result the analysis result.
type Result struct {
	Reason string
	Start  token.Position
	End    token.Position
}

// NewResult creates a new Result.
func NewResult(file *modfile.File, line *modfile.Line, reason string) Result {
	return Result{
		Start:  token.Position{Filename: file.Syntax.Name, Line: line.Start.Line, Column: line.Start.LineRune},
		End:    token.Position{Filename: file.Syntax.Name, Line: line.End.Line, Column: line.End.LineRune},
		Reason: reason,
	}
}

func (r Result) String() string {
	return fmt.Sprintf("%s: %s", r.Start, r.Reason)
}

// Options the analyzer options.
type Options struct {
	ReplaceAllowList          []string
	ReplaceAllowLocal         bool
	ExcludeForbidden          bool
	RetractAllowNoExplanation bool
	ToolchainForbidden        bool
	ToolForbidden             bool
	GoDebugForbidden          bool
	GoVersionPattern          *regexp.Regexp
}

// AnalyzePass analyzes a pass.
func AnalyzePass(pass *analysis.Pass, opts Options) ([]Result, error) {
	info, err := gomod.GetModuleInfo()
	if err != nil {
		return nil, fmt.Errorf("get information about modules: %w", err)
	}

	goMod := info[0].GoMod
	if pass.Module != nil && pass.Module.Path != "" {
		for _, m := range info {
			if m.Path == pass.Module.Path {
				goMod = m.GoMod
				break
			}
		}
	}

	f, err := parseGoMod(goMod)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", goMod, err)
	}

	return AnalyzeFile(f, opts), nil
}

// Analyze analyzes a project.
func Analyze(opts Options) ([]Result, error) {
	f, err := GetModuleFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get module file: %w", err)
	}

	return AnalyzeFile(f, opts), nil
}

// AnalyzeFile analyzes a mod file.
func AnalyzeFile(file *modfile.File, opts Options) []Result {
	checks := []func(file *modfile.File, opts Options) []Result{
		checkRetractDirectives,
		checkExcludeDirectives,
		checkToolDirectives,
		checkReplaceDirectives,
		checkToolchainDirective,
		checkGoDebugDirectives,
		checkGoVersionDirectives,
	}

	var results []Result
	for _, check := range checks {
		results = append(results, check(file, opts)...)
	}

	return results
}

func checkGoVersionDirectives(file *modfile.File, opts Options) []Result {
	var results []Result

	if file.Go != nil && opts.GoVersionPattern != nil && !opts.GoVersionPattern.MatchString(file.Go.Version) {
		results = append(results, NewResult(file, file.Go.Syntax, fmt.Sprintf(reasonGoVersion, file.Go.Version, opts.GoVersionPattern.String())))
	}

	return results
}

func checkRetractDirectives(file *modfile.File, opts Options) []Result {
	var results []Result

	if !opts.RetractAllowNoExplanation {
		for _, r := range file.Retract {
			if r.Rationale != "" {
				continue
			}

			results = append(results, NewResult(file, r.Syntax, reasonRetract))
		}
	}

	return results
}

func checkExcludeDirectives(file *modfile.File, opts Options) []Result {
	var results []Result

	if opts.ExcludeForbidden {
		for _, e := range file.Exclude {
			results = append(results, NewResult(file, e.Syntax, reasonExclude))
		}
	}

	return results
}

func checkToolDirectives(file *modfile.File, opts Options) []Result {
	var results []Result

	if opts.ToolForbidden {
		for _, e := range file.Tool {
			results = append(results, NewResult(file, e.Syntax, reasonTool))
		}
	}

	return results
}

func checkReplaceDirectives(file *modfile.File, opts Options) []Result {
	var results []Result

	uniqReplace := map[string]struct{}{}

	for _, r := range file.Replace {
		reason := checkReplaceDirective(opts, r)
		if reason != "" {
			results = append(results, NewResult(file, r.Syntax, reason))
			continue
		}

		if r.Old.Path == r.New.Path && r.Old.Version == r.New.Version {
			results = append(results, NewResult(file, r.Syntax, reasonReplaceIdentical))
			continue
		}

		if _, ok := uniqReplace[r.Old.Path+r.Old.Version]; ok {
			results = append(results, NewResult(file, r.Syntax, reasonReplaceDuplicate))
		}

		uniqReplace[r.Old.Path+r.Old.Version] = struct{}{}
	}

	return results
}

func checkReplaceDirective(o Options, r *modfile.Replace) string {
	if isLocal(r) {
		if o.ReplaceAllowLocal {
			return ""
		}

		return fmt.Sprintf("%s: %s", reasonReplaceLocal, r.Old.Path)
	}

	for _, v := range o.ReplaceAllowList {
		if r.Old.Path == v {
			return ""
		}
	}

	return fmt.Sprintf("%s: %s", reasonReplace, r.Old.Path)
}

func checkToolchainDirective(file *modfile.File, opts Options) []Result {
	var results []Result

	if opts.ToolchainForbidden && file.Toolchain != nil {
		results = append(results, NewResult(file, file.Toolchain.Syntax, reasonToolchain))
	}

	return results
}

func checkGoDebugDirectives(file *modfile.File, opts Options) []Result {
	var results []Result

	if opts.GoDebugForbidden {
		for _, e := range file.Godebug {
			results = append(results, NewResult(file, e.Syntax, reasonGoDebug))
		}
	}

	return results
}

// Filesystem paths found in "replace" directives are represented by a path with an empty version.
// https://github.com/golang/mod/blob/bc388b264a244501debfb9caea700c6dcaff10e2/module/module.go#L122-L124
func isLocal(r *modfile.Replace) bool {
	return strings.TrimSpace(r.New.Version) == ""
}
