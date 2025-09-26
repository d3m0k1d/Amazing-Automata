package main

import (
	"amazing-automata/cmd"
	_ "fmt"
)

func main() {
	cmd.YamlGenerator("workflow.yml", true, false, false, false)
}
