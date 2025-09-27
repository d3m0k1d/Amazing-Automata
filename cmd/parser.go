package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	_ "embed"
)

type Project struct {
	DependencyFile   string   `json:"project_root_file"`
	InstallCommand   string   `json:"dependency_install_command"`
	BuildCommand     string   `json:"build_command"`
	Name string `json:"name"`
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

//go:embed deps.json
var depfilesjson []byte
func ParseLangDeps() ([]Project, error) {

	var result []Project
	if err := json.Unmarshal(depfilesjson, &result); err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}

	return result, nil
}
