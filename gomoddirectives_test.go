package gomoddirectives

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
)

func TestAnalyze(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(wd)
	}()

	err = os.Chdir("./testdata/replace/")
	require.NoError(t, err)

	results, err := Analyze(Options{})
	require.NoError(t, err)

	assert.Len(t, results, 2)
}

func TestAnalyzeFile(t *testing.T) {
	testCases := []struct {
		desc       string
		modulePath string
		opts       Options
		expected   int
	}{
		{
			desc:       "replace: allow nothing",
			modulePath: "replace/go.mod",
			opts:       Options{},
			expected:   2,
		},
		{
			desc:       "replace: allow an element",
			modulePath: "replace/go.mod",
			opts: Options{
				ReplaceAllowList: []string{
					"github.com/gorilla/mux",
				},
			},
			expected: 1,
		},
		{
			desc:       "replace: allow local",
			modulePath: "replace/go.mod",
			opts: Options{
				ReplaceAllowLocal: true,
			},
			expected: 1,
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
			expected: 0,
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
			expected: 2,
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
			expected: 2,
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
			expected: 1,
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
			expected: 0,
		},
		{
			desc:       "retract: allow no explanation",
			modulePath: "retract/go.mod",
			opts: Options{
				RetractAllowNoExplanation: true,
			},
			expected: 0,
		},
		{
			desc:       "retract: explanation is require",
			modulePath: "retract/go.mod",
			opts: Options{
				RetractAllowNoExplanation: false,
			},
			expected: 1,
		},
		{
			desc:       "exclude: don't allow",
			modulePath: "exclude/go.mod",
			opts: Options{
				ExcludeForbidden: true,
			},
			expected: 2,
		},
		{
			desc:       "exclude: allow",
			modulePath: "exclude/go.mod",
			opts: Options{
				ExcludeForbidden: false,
			},
			expected: 0,
		},
		{
			desc:       "tool: don't allow",
			modulePath: "tool/go.mod",
			opts: Options{
				ToolForbidden: true,
			},
			expected: 1,
		},
		{
			desc:       "tool: don't allow (multiple)",
			modulePath: "tool_multiple/go.mod",
			opts: Options{
				ToolForbidden: true,
			},
			expected: 3,
		},
		{
			desc:       "tool: allow",
			modulePath: "tool/go.mod",
			opts: Options{
				ToolForbidden: false,
			},
			expected: 0,
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

			assert.Len(t, results, test.expected)
		})
	}
}
