package cli

import (
	"flag"
	"github.com/Kalitsune/treeintoreality/types"
	"github.com/charmbracelet/huh"
	"os"
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
	flag.BoolVar(&args.TrustParser, "t", false, "Do not ask for confirmation after parsing.")
	flag.Parse()

	return args
}
