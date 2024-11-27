package lib

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Node struct {
	Name     string
	Children []*Node
	IsDir    bool
}

func MakeTree(treeOutput string) {
	root, err := parseTree(treeOutput)
	if err != nil {
		fmt.Println("Error parsing tree output:", err)
		return
	}

	if !confirmTree(root) {
		fmt.Println("Oops! Cancelling...")
		return
	}

	mode := ""
	err = CreateTree(root, "", &mode)
	if err != nil {
		return
	}
}

// parseTree parses the output of a `tree` command and returns the root node of the tree structure.
func parseTree(treeOutput string) (*Node, error) {
	scanner := bufio.NewScanner(strings.NewReader(treeOutput))
	root := &Node{Name: ".", IsDir: true} // Root node
	nodeStack := []*Node{root}            // Stack to manage hierarchy
	lastDepth := -1                       // Track depth of the previous line

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
		newNode := &Node{
			Name: line[depth:],
		}

		depth /= 8 // Each level is 4 characters of indentation but because we work with uint8s, it's somehow the double.
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
func printTree(node *Node, prefix string) string {
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

func TouchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	return file.Close()
}

func CreateTree(node *Node, prefix string, defaultMode *string) error {
	if node == nil {
		return nil
	}

	if node.Name != "." {
		// create the file/Folder
		if node.IsDir {
			fmt.Println(prefix + node.Name)
			err := os.Mkdir(prefix+node.Name, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			err := TouchFile(prefix + node.Name)
			if err != nil {
				return err
			}
		}
	}

	for _, child := range node.Children {
		childPrefix := prefix + node.Name + "/"

		err := CreateTree(child, childPrefix, defaultMode)
		if err != nil {
			return err
		}
	}

	return nil
}