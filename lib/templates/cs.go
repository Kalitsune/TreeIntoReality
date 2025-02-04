package templates

import (
	"fmt"
	"github.com/Kalitsune/treeintoreality/types"
	"strings"
)

func CsTemplate(node *types.Node, prefix string) string {
	const res = "namespace %s;\n\npublic class %s\n{\n	\n}"

	path := strings.Replace(strings.Trim(prefix, "/"), "/", ".", -1)
	name := strings.TrimRight(node.Name, ".cs")

	return fmt.Sprintf(res, path, name)
}
