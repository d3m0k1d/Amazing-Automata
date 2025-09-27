package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	// "github.com/alperdrsnn/clime"
	// "github.com/samber/lo"
)

const baseTpl = `{{define "base"}}name: {{ .WorkflowName }}
on:
  push:
    branches:
      - master

jobs:
{{end}}`

const ciTpl = `{{define "build_steps"}}    steps:
      - name: Clone repository
        uses: actions/checkout@v4{{ range .Projects }}
      - name: setup env
        uses: {{ or .Type.Setup "actions/setup-* # actions/setup-node@v5 actions/setup-python@v6 actions/setup-java@v5 actions/setup-go@v6 actions/setup-dotnet@v5 ruby/setup-ruby@v1"}}
      - name: Install project dependencies {{.Root}}
        run: {{.Type.InstallCommand}}
        working-directory: {{.Root}}
      - name: Build project {{.Root}}
        run: {{.Type.BuildCommand}}
        working-directory: {{.Root}}{{ end }}{{ end }}{{define "ci"}}{{template "base"}}  build:
    runs-on: ubuntu-latest
{{template "build_steps" .}}
{{ end }}`

const cdTpl = `{{define "cd"}}{{template "base"}}  build:
    strategy: 
            matrix:
                os: [ubuntu-latest, windows-latest, macos-latest]
                include:
                    - os: ubuntu-latest
                    artifacts-path: artifacts/linux
                    - os: windows-latest
                    artifacts-path: artifacts/windows
                    - os: macos-latest
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
                GITHUB_TOKEN: ${{"{{"}} secrets.GITHUB_TOKEN {{"}}"}}{{ end }}
`

type Project struct {
	Type ProjectType
	Root string
}

func walkproj(dir string, types []ProjectType) ([]Project, error) {
	// var err error
	projects := make([]Project, 0)
	var rec func(dir string) error
	// root := dir
	rec = func(dir string) error {

		projects_curr := make([]Project, 0)
		fs, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, f := range fs {
			// if f.Name() == "." || f.Name() == ".." {
			// 	continue
			// }
			if f.IsDir() {
				err := rec(filepath.Join(dir, f.Name()))
				if err != nil {
					return err
				}
				// projects = append(projects, projects2...)
				// f.Name()
			} else {
				for _, m := range types {
					if m.DependencyFileGlob.Match(f.Name()) {
						// r, err := filepath.Rel(root, dir)
						// if err != nil {
						// 	return err
						// }
						r := dir
						projects_curr = append(projects_curr, Project{Type: m, Root: r})
					}
				}
			}
		}
		if len(projects_curr) > 1 {
			// choices, err := clime.AskMultiChoice(fmt.Sprintf("Multiple project lockfiles found in %s:", dir), lo.Map(projects_curr, func(item Project, index int) string { return item.Type.Name })...)
			// if err != nil {
			// 	return err
			// }
			// choosen_projects := make([]Project, 0)
			// for _, k := range choices {
			// 	choosen_projects = append(choosen_projects, projects_curr[k])
			// }

			projects = append(projects, projects_curr...)
		} else {
			projects = append(projects, projects_curr...)
		}

		return nil
	}
	if err := rec(dir); err != nil {
		return nil, err
	}
	return projects, nil
}

func YamlGenerator(filename, projectPath string, ci, cd, dryRun, appendM bool) error {
	// Открытие файла или stdout
	var f *os.File
	var err erro
	if dryRun {
		f = os.Stdout
	} else {
		flags := os.O_CREATE | os.O_WRONLY
		if appendM {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}
		f, err = os.OpenFile(filename, flags, 0o644)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
	}

	types, err := ParseLangDeps()
	if err != nil {
		return fmt.Errorf("failed to parse language deps: %w", err)
	}
	projects, err := walkproj(projectPath, types)
	if err != nil {
		return fmt.Errorf("failed to scan projects: %w", err)
	}

	// По умолчанию оба блока, если ни ci ни cd не заданы
	// if !ci && !cd {
	// 	ci, cd = true, true
	// }

	// if !appendM {
	// 	baseTmpl, err := template.New("base").Parse(baseTpl)
	// 	if err != nil {
	// 		return fmt.Errorf("parse base template: %w", err)
	// 	}
	// 	if err := baseTmpl.Execute(f, map[string]string{"WorkflowName": "Amazing-Automata CI/CD"}); err != nil {
	// 		return fmt.Errorf("execute base template: %w", err)
	// 	}
	// }

	// Общий шаблон с именованными частями
	t := template.New("pipeline")
	t = template.Must(t.Parse(ciTpl))
	t = template.Must(t.Parse(cdTpl))
	t = template.Must(t.Parse(baseTpl))

	data := map[string]interface{}{"Projects": projects}

	if ci {
		if err := t.ExecuteTemplate(f, "ci", data); err != nil {
			return fmt.Errorf("execute ci template: %w", err)
		}
	}

	if cd {
		if err := t.ExecuteTemplate(f, "cd", data); err != nil {
			return fmt.Errorf("execute cd template: %w", err)
		}
	}

	if dryRun {
		fmt.Fprintln(os.Stderr, "\n# Dry-run complete: pipeline preview above")
	}

	return nil
}
