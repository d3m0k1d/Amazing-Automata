package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/gobwas/glob"
	"github.com/samber/lo"
)

type ProjectTypeDTO struct {
	DependencyFile string `json:"project_root_file"`
	InstallCommand string `json:"dependency_install_command"`
	BuildCommand   string `json:"build_command"`
	Name           string `json:"name"`
}
type ProjectType struct {
	ProjectTypeDTO
	DependencyFileGlob glob.Glob
}

//go:embed deps.json
var depfilesjson []byte

func ParseLangDeps() ([]ProjectType, error) {

	var result []ProjectTypeDTO
	if err := json.Unmarshal(depfilesjson, &result); err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}

	return lo.Map(result, func(item ProjectTypeDTO, index int) ProjectType {
		return ProjectType{item, glob.MustCompile(item.DependencyFile)}
	}), nil
}
