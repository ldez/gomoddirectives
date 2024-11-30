package gomoddirectives

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

type modInfo struct {
	Path      string `json:"Path"`
	Dir       string `json:"Dir"`
	GoMod     string `json:"GoMod"`
	GoVersion string `json:"GoVersion"`
	Main      bool   `json:"Main"`
}

// GetModuleFile gets module file.
// It's better to use [GetGoModFile] instead of this function.
func GetModuleFile() (*modfile.File, error) {
	goMod, err := getFirstModulePath()
	if err != nil {
		return nil, err
	}

	if goMod == "" {
		return nil, errors.New("working directory is not part of a module")
	}

	return parseGoMod(goMod)
}

func getFirstModulePath() (string, error) {
	// https://github.com/golang/go/issues/44753#issuecomment-790089020
	cmd := exec.Command("go", "list", "-m", "-json")

	raw, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command go list: %w: %s", err, string(raw))
	}

	var v modInfo
	err = json.NewDecoder(bytes.NewBuffer(raw)).Decode(&v)
	if err != nil {
		return "", fmt.Errorf("unmarshaling error: %w: %s", err, string(raw))
	}

	return v.GoMod, nil
}

func parseGoMod(goMod string) (*modfile.File, error) {
	raw, err := os.ReadFile(filepath.Clean(goMod))
	if err != nil {
		return nil, fmt.Errorf("reading go.mod file: %w", err)
	}

	return modfile.Parse("go.mod", raw, nil)
}

func getModuleInfo() ([]modInfo, error) {
	// https://github.com/golang/go/issues/44753#issuecomment-790089020
	cmd := exec.Command("go", "list", "-m", "-json")

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command go list: %w: %s", err, string(out))
	}

	var infos []modInfo

	for dec := json.NewDecoder(bytes.NewBuffer(out)); dec.More(); {
		var v modInfo
		if err := dec.Decode(&v); err != nil {
			return nil, fmt.Errorf("unmarshaling error: %w: %s", err, string(out))
		}

		if v.GoMod == "" {
			return nil, errors.New("working directory is not part of a module")
		}

		if !v.Main || v.Dir == "" {
			continue
		}

		infos = append(infos, v)
	}

	if len(infos) == 0 {
		return nil, errors.New("go.mod file not found")
	}

	return infos, nil
}
