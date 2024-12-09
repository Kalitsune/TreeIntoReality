package main

import (
	"github.com/charmbracelet/huh"
	"os"
	"treeintoreality/lib"
)

func main() {
	tree := ""

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Insert your tree output:").
				Value(&tree).
				CharLimit(-1),
		),
	).Run()
	if err != nil {
		return
	}

	prefix, err := os.Getwd()
	if err != nil {
		prefix = "."
	}

	lib.MakeTree(prefix+"/", tree)
}
