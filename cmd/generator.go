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
    steps:
	  - name: setup env
        uses: actions/setup-* # actions/setup-node@v5 actions/setup-python@v6 actions/setup-java@v5 actions/setup-go@v6 actions/setup-dotnet@v5 ruby/setup-ruby@v1
      - name: Clone repository
        uses: actions/checkout@v4
      {{ range .Projects }}
      - name: setup
        with: {{.setup}}
      - name: Build project {{.name}}
        run: {{.cmd}}
      {{ end }}
`

const cdTpl = `  build:
    runs-on: ubuntu-latest
    steps:
      - name: Clone repository
        uses: actions/checkout@v4
      {{ range .Projects }}
      - name: setup
        with:
      - name: Install deps for project {{.Root}}
        run: {{.Type.InstallCommand}}
      - name: Build project {{.Root}}
        run: {{.Type.BuildCommand}}
      - name: Build project {{.Root}}
        run: {{.Type.BuildCommand}}{{ end }}`

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
	if ci {
		if err := template.Must(template.New("ci").Parse(ciTpl)).Execute(f, map[string]interface{}{"Projects": projects}); err != nil {
			return err
		}
	}
	if cd {
		if err := template.Must(template.New("cd").Parse(cdTpl)).Execute(f, map[string]interface{}{"Projects": projects}); err != nil {
			return err
		}
	}

	return nil
}
