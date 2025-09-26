package cmd

import (
	"os"
	"text/template"
)

const baseTpl = `name: {{ .WorkflowName }}
on:
  push:
    branches:
      - master

jobs:
`

const buildTpl = `  build:
    runs-on: ubuntu-latest
    steps:
      - name: Clone repository
        uses: actions/checkout@v4
      {{ range .Projects }}
      - name: setup
        with: {{.setup}}
      - name: Build project {{.name}}
        run: {{.cmd}}
      {{ end }}
`

// YamlGenerator создаёт YAML и условно добавляет шаг Checkout для CI
func YamlGenerator(filename string, ci, cd, dryRun, appendM bool) error {
	// 1. Открываем или создаём файл, перезаписываем содержимое
	f, err := os.OpenFile(filename,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0o644,
	)
	if err != nil {
		return err
	}
	defer f.Close()

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

	if ci {
		checkoutTmpl, err := template.New("checkout").Parse(checkoutTpl)
		if err != nil {
			return err
		}
		if err := checkoutTmpl.Execute(f, nil); err != nil {
			return err
		}
	}

	return nil
}
