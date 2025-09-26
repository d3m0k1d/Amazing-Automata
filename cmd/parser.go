package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Dependency struct {
	DependencyFile   string   `json:"dependency_file"`
	InstallCommand   string   `json:"install_command"`
	BuildCommand     string   `json:"build_command"`
	SourceExtensions []string `json:"source_extensions"`
}

type LangDeps struct {
	Name         string
	Dependencies []Dependency
}

func CollectFiles() ([]string, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var files []string
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			name := filepath.Base(path)
			files = append(files, name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func PathCollectFiles() ([]string, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var files []string
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func ParseLangDeps(path string) ([]LangDeps, error) {
	rawData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var tmp map[string][]Dependency
	if err := json.Unmarshal(rawData, &tmp); err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}

	result := make([]LangDeps, 0, len(tmp))
	for name, deps := range tmp {
		result = append(result, LangDeps{
			Name:         name,
			Dependencies: deps,
		})
	}

	return result, nil
}
