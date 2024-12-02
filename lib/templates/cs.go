package templates

import (
	"fmt"
	"strings"
	"treeintoreality/types"
)

func CsTemplate(node *types.Node, prefix string) string {
	const res = "namespace %s;\n\npublic class %s\n{\n	\n}"

	path := strings.Replace(strings.Trim(prefix[2:], "/"), "/", ".", -1)
	name := strings.TrimRight(node.Name, ".cs")

	return fmt.Sprintf(res, path, name)
}
