package lib

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"treeintoreality/lib/templates"
	"treeintoreality/types"
)

func MakeTree(rootDir string, treeOutput string, args *types.Args) {
	root, err := parseTree(treeOutput)
	if err != nil {
		fmt.Println("Error parsing tree output:", err)
		return
	}

	if !confirmTree(root) {
		fmt.Println("Oops! Cancelling...")
		return
	}

	err = CreateTree(root, rootDir, "", args)
	if err != nil {
		return
	}
}

// parseTree parses the output of a `tree` command and returns the root node of the tree structure.
func parseTree(treeOutput string) (*types.Node, error) {
	scanner := bufio.NewScanner(strings.NewReader(treeOutput))
	root := &types.Node{Name: ".", IsDir: true} // Root node
	nodeStack := []*types.Node{root}            // Stack to manage hierarchy
	lastDepth := -1                             // Track depth of the previous line

	for scanner.Scan() {
		text := scanner.Text()
		line := strings.TrimSpace(text)

		// Handle directory/file counts (e.g., "3 directories, 8 files")
		if !strings.Contains(line, " ") || strings.Contains(line, "directories") || strings.Contains(line, "files") {
			continue
		}

		// Determine the depth of the current line based on leading characters.
		depth := 0
		for i := 0; i < len(line); i++ {
			if slices.Contains([]uint8{' ', ' ', 226, 148, 156, 128, 130, 194}, line[i]) {
				depth++
			} else {
				break
			}
		}

		// Create a new node.
		newNode := &types.Node{
			Name: line[depth:],
		}

		depth = len(strings.Split(line, "   "))
		// Adjust the stack based on depth.
		if depth > lastDepth {
			// Child node, add to the last node in the stack.
			if len(nodeStack) > 0 {
				parent := nodeStack[len(nodeStack)-1]
				parent.IsDir = true
				parent.Children = append(parent.Children, newNode)
			}
		} else {
			// Go up the stack until we find the right parent depth.
			nodeStack = nodeStack[:depth]
			if len(nodeStack) > 0 {
				parent := nodeStack[len(nodeStack)-1]
				parent.Children = append(parent.Children, newNode)
			}
		}

		// Update stack and depth.
		nodeStack = append(nodeStack, newNode)
		lastDepth = depth
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return root, nil
}

// printTree recursively prints the tree structure for visualization.
func printTree(node *types.Node, prefix string) string {
	if node == nil {
		return ""
	}

	name := node.Name
	if node.IsDir {
		name += "/"
	}

	res := ""
	if len(node.Name) < 2 {
		res = prefix + name + "\n"
	} else {
		res = prefix + " " + name + "\n"
	}
	for _, child := range node.Children {
		childPrefix := prefix + "────"
		res += printTree(child, childPrefix)
	}
	return res
}

func TouchFile(node *types.Node, rootDir string, prefix string) error {
	file, err := os.OpenFile(rootDir+prefix+node.Name, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	if strings.HasSuffix(node.Name, ".cs") {
		_, _ = file.WriteString(templates.CsTemplate(node, prefix))
	}

	return file.Close()
}

func enforceArgs(path string, args *types.Args) bool {
	if _, err := os.Stat(path); err == nil {
		if args.Overwrite && !strings.HasSuffix(path, "/") {
			err = os.Remove(path)
			return err == nil
		} else {
			return false
		}
	}
	return true
}
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
func CreateTree(node *types.Node, root string, prefix string, args *types.Args) error {
	if node == nil {
		return nil
	}

	if node.Name != "." {
		// create the file/Folder
		if node.IsDir {
			if !pathExists(root + prefix + node.Name + "/") {
				err := os.Mkdir(root+prefix+node.Name, os.ModePerm)
				if err != nil {
					return err
				}
			}
		} else {
			exists := pathExists(root + prefix + node.Name)

			if exists && args.Overwrite {
				err := os.Remove(root + prefix + node.Name)
				if err != nil {
					return err
				}
			}

			if !exists || args.Overwrite {
				err := TouchFile(node, root, prefix)
				if err != nil {
					return err
				}
			}
		}

		prefix = prefix + node.Name + "/"
	}

	for _, child := range node.Children {
		err := CreateTree(child, root, prefix, args)
		if err != nil {
			return err
		}
	}

	return nil
}
