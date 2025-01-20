package cli

import (
	"fmt"
	"github.com/Kalitsune/treeintoreality/lib"
	"github.com/Kalitsune/treeintoreality/types"
	"github.com/charmbracelet/huh"
)

func MakeTree(rootDir string, treeOutput string, args *types.Args) {
	root, err := lib.ParseTree(treeOutput)
	if err != nil {
		fmt.Println("Error parsing tree output:", err)
		return
	}

	if !args.TrustParser && !confirmTree(root) {
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
