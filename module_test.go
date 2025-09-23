package gomoddirectives

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetModuleFile(t *testing.T) {
	t.Chdir("./testdata/replace/")

	file, err := GetModuleFile()
	require.NoError(t, err)

	assert.Equal(t, "github.com/ldez/gomoddirectives/testdata/replace", file.Module.Mod.Path)
}

func TestGetModuleFile_here(t *testing.T) {
	file, err := GetModuleFile()
	require.NoError(t, err)

	assert.Equal(t, "github.com/ldez/gomoddirectives", file.Module.Mod.Path)
}
