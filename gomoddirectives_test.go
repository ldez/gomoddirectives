package gomoddirectives

import (
	"cmp"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
)

func TestAnalyze(t *testing.T) {
	t.Chdir("./testdata/replace/")

	results, err := Analyze(Options{})
	require.NoError(t, err)

	assert.Len(t, results, 2)
}

func TestAnalyzeFile(t *testing.T) {
	testCases := []struct {
		desc       string
		modulePath string
		opts       Options
		expected   []Result
	}{
		{
			desc:       "replace: allow nothing",
			modulePath: "replace/go.mod",
			opts:       Options{},
			expected: []Result{
				{
					Reason: "replacement are not allowed: github.com/gorilla/mux",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 12, Column: 2},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 12, Column: 88},
				},
				{
					Reason: "local replacement are not allowed: github.com/ldez/grignotin",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 13, Column: 2},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 13, Column: 36},
				},
			},
		},
		{
			desc:       "replace: allow an element",
			modulePath: "replace/go.mod",
			opts: Options{
				ReplaceAllowList: []string{
					"github.com/gorilla/mux",
				},
			},
			expected: []Result{{
				Reason: "local replacement are not allowed: github.com/ldez/grignotin",
				Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 13, Column: 2},
				End:    token.Position{Filename: "go.mod", Offset: 0, Line: 13, Column: 36},
			}},
		},
		{
			desc:       "replace: allow local",
			modulePath: "replace/go.mod",
			opts: Options{
				ReplaceAllowLocal: true,
			},
			expected: []Result{{
				Reason: "replacement are not allowed: github.com/gorilla/mux",
				Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 12, Column: 2},
				End:    token.Position{Filename: "go.mod", Offset: 0, Line: 12, Column: 88},
			}},
		},
		{
			desc:       "replace: exclude all",
			modulePath: "replace/go.mod",
			opts: Options{
				ReplaceAllowLocal: true,
				ReplaceAllowList: []string{
					"github.com/ldez/grignotin",
					"github.com/gorilla/mux",
				},
			},
		},
		{
			desc:       "replace: allow list doesn't override allow local",
			modulePath: "replace/go.mod",
			opts: Options{
				ReplaceAllowLocal: false,
				ReplaceAllowList: []string{
					"github.com/ldez/grignotin",
				},
			},
			expected: []Result{
				{
					Reason: "replacement are not allowed: github.com/gorilla/mux",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 12, Column: 2},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 12, Column: 88},
				},
				{
					Reason: "local replacement are not allowed: github.com/ldez/grignotin",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 13, Column: 2},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 13, Column: 36},
				},
			},
		},
		{
			desc:       "replace: duplicate replacement",
			modulePath: "replace_duplicate/go.mod",
			opts: Options{
				ReplaceAllowLocal: true,
				ReplaceAllowList: []string{
					"github.com/gorilla/mux",
					"github.com/ldez/grignotin",
				},
			},
			expected: []Result{
				{
					Reason: "multiple replacement of the same module",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 17, Column: 2},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 17, Column: 88},
				},
				{
					Reason: "multiple replacement of the same module",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 18, Column: 2},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 18, Column: 36},
				},
			},
		},
		{
			desc:       "replace: replaced with identical element",
			modulePath: "replace_identical/go.mod",
			opts: Options{
				ReplaceAllowLocal: true,
				ReplaceAllowList: []string{
					"github.com/gorilla/mux",
					"github.com/ldez/grignotin",
				},
			},
			expected: []Result{{
				Reason: "the original module and the replacement are identical",
				Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 11, Column: 2},
				End:    token.Position{Filename: "go.mod", Offset: 0, Line: 11, Column: 64},
			}},
		},
		{
			desc:       "replace: duplicate replacement but for the different versions",
			modulePath: "replace_duplicate_versions/go.mod",
			opts: Options{
				ReplaceAllowLocal: true,
				ReplaceAllowList: []string{
					"github.com/gorilla/mux",
					"github.com/ldez/grignotin",
				},
			},
		},
		{
			desc:       "retract: allow no explanation",
			modulePath: "retract/go.mod",
			opts: Options{
				RetractAllowNoExplanation: true,
			},
		},
		{
			desc:       "retract: explanation is require",
			modulePath: "retract/go.mod",
			opts: Options{
				RetractAllowNoExplanation: false,
			},
			expected: []Result{{
				Reason: "a comment is mandatory to explain why the version has been retracted",
				Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 13, Column: 5},
				End:    token.Position{Filename: "go.mod", Offset: 0, Line: 13, Column: 21},
			}},
		},
		{
			desc:       "exclude: don't allow",
			modulePath: "exclude/go.mod",
			opts: Options{
				ExcludeForbidden: true,
			},
			expected: []Result{
				{
					Reason: "exclude directive is not allowed",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 12, Column: 5},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 12, Column: 31},
				},
				{
					Reason: "exclude directive is not allowed",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 13, Column: 5},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 13, Column: 29},
				},
			},
		},
		{
			desc:       "exclude: allow",
			modulePath: "exclude/go.mod",
			opts: Options{
				ExcludeForbidden: false,
			},
		},
		{
			desc:       "ignore: don't allow",
			modulePath: "ignore/go.mod",
			opts: Options{
				IgnoreForbidden: true,
			},
			expected: []Result{
				{
					Reason: "ignore directive is not allowed",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 6, Column: 2},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 6, Column: 16},
				},
				{
					Reason: "ignore directive is not allowed",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 7, Column: 2},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 7, Column: 9},
				},
			},
		},
		{
			desc:       "ignore: allow",
			modulePath: "ignore/go.mod",
			opts: Options{
				IgnoreForbidden: false,
			},
		},
		{
			desc:       "tool: don't allow",
			modulePath: "tool/go.mod",
			opts: Options{
				ToolForbidden: true,
			},
			expected: []Result{{
				Reason: "tool directive is not allowed",
				Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 5, Column: 1},
				End:    token.Position{Filename: "go.mod", Offset: 0, Line: 5, Column: 45},
			}},
		},
		{
			desc:       "tool: don't allow (multiple)",
			modulePath: "tool_multiple/go.mod",
			opts: Options{
				ToolForbidden: true,
			},
			expected: []Result{
				{
					Reason: "tool directive is not allowed",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 5, Column: 1},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 5, Column: 37},
				},
				{
					Reason: "tool directive is not allowed",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 8, Column: 5},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 8, Column: 29},
				},
				{
					Reason: "tool directive is not allowed",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 9, Column: 5},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 9, Column: 29},
				},
			},
		},
		{
			desc:       "tool: allow",
			modulePath: "tool/go.mod",
			opts: Options{
				ToolForbidden: false,
			},
		},
		{
			desc:       "godebug: don't allow",
			modulePath: "godebug/go.mod",
			opts: Options{
				GoDebugForbidden: true,
			},
			expected: []Result{
				{
					Reason: "godebug directive is not allowed",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 5, Column: 1},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 5, Column: 23},
				},
				{
					Reason: "godebug directive is not allowed",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 7, Column: 5},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 7, Column: 15},
				},
				{
					Reason: "godebug directive is not allowed",
					Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 8, Column: 5},
					End:    token.Position{Filename: "go.mod", Offset: 0, Line: 8, Column: 21},
				},
			},
		},
		{
			desc:       "godebug: allow",
			modulePath: "godebug/go.mod",
			opts: Options{
				GoDebugForbidden: false,
			},
		},
		{
			desc:       "goversion: don't check",
			modulePath: "goversion_family/go.mod",
			opts: Options{
				GoVersionPattern: nil,
			},
		},
		{
			desc:       "goversion: pattern match",
			modulePath: "goversion_family/go.mod",
			opts: Options{
				GoVersionPattern: regexp.MustCompile(`\d\.\d+(\.0)?`),
			},
		},
		{
			desc:       "goversion: pattern not matched",
			modulePath: "goversion_family/go.mod",
			opts: Options{
				GoVersionPattern: regexp.MustCompile(`\d\.\d+\.0$`),
			},
			expected: []Result{{
				Reason: "go directive (1.22) doesn't match the pattern '\\d\\.\\d+\\.0$'",
				Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 3, Column: 1},
				End:    token.Position{Filename: "go.mod", Offset: 0, Line: 3, Column: 8},
			}},
		},
		{
			desc:       "goversion: no Go version",
			modulePath: "empty/go.mod",
			opts: Options{
				GoVersionPattern: regexp.MustCompile(`\d\.\d+(\.0)?$`),
			},
		},
		{
			desc:       "toolchain: don't allow",
			modulePath: "toolchain/go.mod",
			opts: Options{
				ToolchainForbidden: true,
				ToolchainPattern:   regexp.MustCompile(`go\d\.\d+\.\d+$`),
			},
			expected: []Result{{
				Reason: "toolchain directive is not allowed",
				Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 5, Column: 1},
				End:    token.Position{Filename: "go.mod", Offset: 0, Line: 5, Column: 19},
			}},
		},
		{
			desc:       "toolchain: allow",
			modulePath: "toolchain/go.mod",
			opts: Options{
				ToolchainForbidden: false,
			},
		},
		{
			desc:       "toolchain: pattern match",
			modulePath: "toolchain/go.mod",
			opts: Options{
				ToolchainPattern: regexp.MustCompile(`go\d\.\d+\.\d+$`),
			},
		},
		{
			desc:       "toolchain: pattern not matched",
			modulePath: "toolchain/go.mod",
			opts: Options{
				ToolchainPattern: regexp.MustCompile(`go\d\.22\.\d+$`),
			},
			expected: []Result{{
				Reason: "toolchain directive (go1.23.3) doesn't match the pattern 'go\\d\\.22\\.\\d+$'",
				Start:  token.Position{Filename: "go.mod", Offset: 0, Line: 5, Column: 1},
				End:    token.Position{Filename: "go.mod", Offset: 0, Line: 5, Column: 19},
			}},
		},
		{
			desc:       "toolchain: no Go version",
			modulePath: "empty/go.mod",
			opts: Options{
				ToolchainPattern: regexp.MustCompile(`go\d\.22\.\d+$`),
			},
		},
		{
			desc:       "all: empty go.mod",
			modulePath: "empty/go.mod",
			opts: Options{
				ReplaceAllowLocal: true,
				ReplaceAllowList: []string{
					"github.com/gorilla/mux",
					"github.com/ldez/grignotin",
				},
				ExcludeForbidden:          true,
				RetractAllowNoExplanation: false,
				ToolchainForbidden:        true,
				ToolchainPattern:          regexp.MustCompile(`go\d\.\d+\.\d+$`),
				ToolForbidden:             true,
				GoDebugForbidden:          true,
				GoVersionPattern:          regexp.MustCompile(`\d\.\d+(\.0)?$`),
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			raw, err := os.ReadFile(filepath.FromSlash("./testdata/" + test.modulePath))
			require.NoError(t, err)

			file, err := modfile.Parse("go.mod", raw, nil)
			require.NoError(t, err)

			results := AnalyzeFile(file, test.opts)

			slices.SortFunc(results, func(a, b Result) int {
				return cmp.Or(cmp.Compare(a.Start.Line, b.Start.Line), cmp.Compare(a.End.Line, b.End.Line))
			})

			assert.Equal(t, test.expected, results)
		})
	}
}
