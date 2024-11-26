package main

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/huh"
	"slices"
	"strings"
)

func main() {
	tree := ""
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Insert your tree output:").
				Value(&tree),
		),
	).Run()
	if err != nil {
		return
	}

	makeTree(tree)
}

func makeTree(treeOutput string) {
	root, err := parseTree(treeOutput)
	if err != nil {
		fmt.Println("Error parsing tree output:", err)
		return
	}

	confirmTree(root)
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

func confirmTree(node *Node) bool {

	success := true
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Parsing Successful, is that alright?").
				Description(printTree(node, "")).
				Value(&success),
		),
	).Run()
	if err != nil {
		return false
	}

	return success
}

// printTree recursively prints the tree structure for visualization.
func printTree(node *Node, prefix string) string {
	if node == nil {
		return ""
	}

	res := ""
	if len(node.Name) < 2 {
		res = prefix + node.Name + "\n"
	} else {
		res = prefix + " " + node.Name + "\n"
	}
	for _, child := range node.Children {
		childPrefix := prefix + "────"
		res += printTree(child, childPrefix)
	}
	return res
}
