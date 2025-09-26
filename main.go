package main

import (
	"amazing-automata/cmd"
	"fmt"
)

func main() {
	files, err := cmd.CollectFiles()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(files)
}
