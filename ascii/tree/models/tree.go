package models

import "github.com/pkg/errors"

type Tree struct {
	root *Node
}

type Node struct {
	data     string
	children []*Node
}

func NewNode(data string) *Node {
	return &Node{
		data: data,
	}
}

func (n *Node) AddChild(data string) *Node {
	node := &Node{
		data: data,
	}
	n.children = append(n.children, node)
	return node
}
func NewTree() *Tree {
	return &Tree{}
}

func (t *Tree) AddRoot(data string) (*Node, error) {
	if t.root == nil {
		node := NewNode(data)
		t.root = node
		return node, nil
	}
	return nil, errors.New("已有 root")
}

func (t Tree) Draw() string {
	if t.root == nil {
		return ""
	}
	strs := make([]string, 0, 10)
	current := make([]*Node, 0, 10)
	current = append(current, t.root)

	for _, node := range current {

	}
}
