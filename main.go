package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/huh"
	"os"
	"treeintoreality/lib"
	"treeintoreality/types"
)

func main() {
	args := parseArgs()

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

	MakeTree(prefix+"/", tree, &args)
}

func parseArgs() types.Args {
	args := types.Args{}

	flag.BoolVar(&args.Overwrite, "o", false, "will overwrite files if they exist when recreating the tree.")
	flag.Parse()

	return args
}

func MakeTree(rootDir string, treeOutput string, args *types.Args) {
	root, err := lib.ParseTree(treeOutput)
	if err != nil {
		fmt.Println("Error parsing tree output:", err)
		return
	}

	if !confirmTree(root) {
		fmt.Println("Oops! Cancelling...")
		return
	}

	err = lib.CreateTree(root, rootDir, "", args)
	if err != nil {
		return
	}
}

func confirmTree(node *types.Node) bool {
	success := true
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Parsing complete.").
				Description(lib.PrintTree(node, "")).
				Affirmative("Sounds about right!").
				Negative("Hold on...").
				Value(&success),
		),
	).Run()
	if err != nil {
		return false
	}

	return success
}
