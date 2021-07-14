package gomoddirectives

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/mod/modfile"
)

// GetModuleFile gets module file.
func GetModuleFile() (*modfile.File, error) {
	// https://github.com/golang/go/issues/44753#issuecomment-790089020
	cmd := exec.Command("go", "list", "-m", "-f", "{{.GoMod}}")

	raw, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("command go list: %w: %s", err, string(raw))
	}
	raw = bytes.TrimSpace(raw)

	if len(raw) == 0 {
		return nil, errors.New("working directory is not part of a module")
	}

	raw, err = os.ReadFile(string(raw))
	if err != nil {
		return nil, fmt.Errorf("reading go.mod file: %w", err)
	}

	return modfile.Parse("go.mod", raw, nil)
}
