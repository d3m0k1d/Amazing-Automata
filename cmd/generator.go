package cmd

import (
	"os"
	"path/filepath"
	"text/template"
)

const baseTpl = `name: {{ .WorkflowName }}
on:
  push:
    branches:
      - master

jobs:
`

const ciTpl = `  build:
    runs-on: ubuntu-latest
{{define "build_steps"}}    steps:
      - name: Clone repository
        uses: actions/checkout@v4{{ range .Projects }}
      - name: setup env
        uses: {{ or .Type.Setup "actions/setup-* # actions/setup-node@v5 actions/setup-python@v6 actions/setup-java@v5 actions/setup-go@v6 actions/setup-dotnet@v5 ruby/setup-ruby@v1"}}
      - name: Install project dependencies {{.Root}}
        run: {{.Type.InstallCommand}}
      - name: Build project {{.Root}}
        run: {{.Type.BuildCommand}}{{ end }}{{ end }}`

const cdTpl = `  build:
    strategy: 
            matrix:
                os: [ubuntu-latest, windows-latest, macos-latest]
                include:
                    - os: ubuntu-latest
                    artifacts-path: artifacts/linux
                    - os: windows-latest
                    artifacts-path: artifacts/windows
                    -os: macos-latest
                    artifacts-path: artifacts/macos
                
    runs-on: ubuntu-latest
{{template "build_steps" .}}
  release:
    needs: build
    runs-on: ubuntu-latest
        steps:
          - name: Download artifacts
            uses: actions/download-artifact@v4
            with:
                path: release-artifacts
          - name: Create Github release
            uses: actions/create-github-release@v4
            with: 
                files: |
                    release-artifacts/**/*
                draft: false
                prerelease: false
            env:
                GITHUB_TOKEN: ${{"{{"}} secrets.GITHUB_TOKEN {{"}}"}}
`

type Project struct {
	Type ProjectType
	Root string
}

func walkproj(dir string, types []ProjectType) ([]Project, error) {
	// var err error
	projects := make([]Project, 0)
	matchfurther := true
	fs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		if f.IsDir() {
			projects2, err := walkproj(filepath.Join(dir, f.Name()), types)
			if err != nil {
				return nil, err
			}
			projects = append(projects, projects2...)
			// f.Name()
		} else {
			if matchfurther {
				for _, m := range types {
					if m.DependencyFileGlob.Match(f.Name()) {
						// println("yea")
						// println(dir)
						projects = append(projects, Project{Type: m, Root: dir})
						matchfurther = false
						// break
					}
				}
			}
		}
	}
	return projects, nil
}

// YamlGenerator создаёт YAML и условно добавляет шаг Checkout для CI
func YamlGenerator(filename string, projectPath string, ci, cd, dryRun, appendM bool) error {
	// 1. Открываем или создаём файл, перезаписываем содержимое
	f, err := os.OpenFile(filename,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0o644,
	)
	if err != nil {
		return err
	}
	defer f.Close()

	// root, err := os.Getwd()
	// if err != nil {
	// 	return nil, err
	// }
	types, err := ParseLangDeps()
	if err != nil {
		return err
	}
	projects, err := walkproj(projectPath, types)
	if err != nil {
		return err
	}
	// 2. Парсим и рендерим базовый шаблон
	baseTmpl, err := template.New("base").Parse(baseTpl)
	if err != nil {
		return err
	}
	if err := baseTmpl.Execute(f, map[string]string{
		"WorkflowName": "Amazing-Automata CI/CD",
	}); err != nil {
		return err
	}

	// fmt.Print(len(projects))
	t := template.New("whocares")
	t = template.Must(t.Parse(ciTpl))
	t = template.Must(t.Parse(cdTpl))
	if ci {
		if err := t.Execute(f, map[string]interface{}{"Projects": projects}); err != nil {
			return err
		}
	}
	if cd {
		if err := t.Execute(f, map[string]interface{}{"Projects": projects}); err != nil {
			return err
		}
	}

	return nil
}
