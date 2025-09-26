package cmd

import (
	_ "fmt"
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

func YamlGenerator(filename string, cd bool, ci bool, dryRun bool, appendM bool) error {

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	tpl, err := template.New("base").Parse(baseTpl)
	if err != nil {
		return err
	}

	data := map[string]string{
		"WorkflowName": "Amazing-Automata CI/CD",
	}

	return tpl.Execute(f, data)

	if ci {
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		defer f.Close()
	}
	if cd {

	}
	if dryRun {

	}
	if appendM {

	}
	return nil
}
