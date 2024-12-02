package types

type Node struct {
	Name     string
	Children []*Node
	IsDir    bool
}
