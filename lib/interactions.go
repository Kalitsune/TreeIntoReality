package lib

import (
	"github.com/charmbracelet/huh"
	"treeintoreality/types"
)

func confirmTree(node *types.Node) bool {

	success := true
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Parsing complete.").
				Description(printTree(node, "")).
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
