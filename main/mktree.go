package main

import (
	"bufio"
	"fmt"
	"slices"
	"strings"
)

func MakeTree(treeOutput string) {
	root, err := parseTree(treeOutput)
	if err != nil {
		fmt.Println("Error parsing tree output:", err)
		return
	}

	printTree(root, "")
}

type Node struct {
	Name     string
	Children []*Node
	IsDir    bool
}

// parseTree parses the output of a `tree` command and returns the root node of the tree structure.
func parseTree(treeOutput string) (*Node, error) {
	scanner := bufio.NewScanner(strings.NewReader(treeOutput))
	root := &Node{Name: ".", IsDir: true} // Root node
	nodeStack := []*Node{root}            // Stack to manage hierarchy
	lastDepth := -1                       // Track depth of the previous line

	for scanner.Scan() {
		line := scanner.Text()

		// Handle directory/file counts (e.g., "3 directories, 8 files")
		if strings.TrimSpace(line) == "" || strings.Contains(line, "directories") || strings.Contains(line, "files") {
			continue
		}

		// Determine the depth of the current line based on leading characters.
		depth := 0
		for i := 0; i < len(line); i++ {
			if line[i] == ' ' || line[i] == ' ' || slices.Contains([]uint8{226, 148, 156, 128, 130, 194}, line[i]) {
				depth++
			} else {
				break
			}
		}

		// Extract the name and check if it's a directory.
		trimmed := strings.TrimSpace(line[depth:])
		name := strings.TrimSuffix(trimmed, "/")

		// Create a new node.
		newNode := &Node{
			Name: name,
		}

		depth /= 8 // Each level is 4 characters of indentation but because we work with uint8s, it's somehow the double.
		// Adjust the stack based on depth.
		if depth > lastDepth {
			newNode.IsDir = true
			// Child node, add to the last node in the stack.
			if len(nodeStack) > 0 {
				parent := nodeStack[len(nodeStack)-1]
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
func printTree(node *Node, prefix string) {
	if node == nil {
		return
	}

	fmt.Println(prefix + node.Name)
	for _, child := range node.Children {
		childPrefix := prefix + "    "
		printTree(child, childPrefix)
	}
}
