package main

import "fmt"

const MAX_SIZE = 4

type BpTreeNode struct {
	Key      int
	Value    string
	Children []*BpTreeNode
}

func (bpn *BpTreeNode) insertChild(n int) {
	if len(bpn.Children) == 0 {
		bpn.Children = append(bpn.Children, &BpTreeNode{Key: n})
		return
	}

	var i int
	for index, v := range bpn.Children {
		if v.Key > n {
			i = index - 1
			break
		}
	}

	bpn.Children = append(bpn.Children, &BpTreeNode{})
	copy(bpn.Children[i+1:], bpn.Children[i:])
	bpn.Children[i] = &BpTreeNode{Key: n}
}

func NewBpTree() *BpTreeNode {
	tree := BpTreeNode{}
	return &tree
}

// We will use BpTreeNode for both intermediate and Leaf Nodes
// Lets scope ths just for positive integer``

/**
intermediate node range: [currentVal, next node's Val)
*/

// returns index of to be parent node
func findLeafNodeParent(refs []*BpTreeNode, val int) int {
	// only current level for now
	// if 1 ref => then 0
	// if 2 ref => if less than second => 1st other wise second
	// if 3 refs

	cur := 0

	nodes := refs
	for len(nodes) > 2 {
		midpoint := len(nodes) / 2

		if val >= nodes[midpoint].Key {
			cur += midpoint
			nodes = nodes[midpoint:]
		} else {
			nodes = nodes[:midpoint]
		}
	}

	if len(nodes) == 0 {
		return -1 // Handle empty case, though probably not reached
	}

	if len(nodes) == 1 {
		if nodes[0].Key > val {
			return -1
		}
		return cur + 0
	}

	if len(nodes) == 2 {
		if nodes[0].Key > val {
			return -1
		}

		if nodes[0].Key < val && val < nodes[1].Key {
			return cur + 0
		}

		return cur + 1
	}

	return -1 // Fallback, should not reach
}

func (bp *BpTreeNode) Insert(n int) {
	// If it's the first element to be inserted then
	if len(bp.Children) == 0 {
		bp.Children = append(bp.Children, &BpTreeNode{Key: n})
		return
	}
	// find leaf node parent
	index := findLeafNodeParent(bp.Children, n)

	if index == -1 {
		// If index is -1 , I have to insert a new intermediate node
		// I also have to keep in mind that I can't go over the max size of intermediate nodes

		if len(bp.Children)+1 < MAX_SIZE {
			// add in the intermediate node
			// add the element here
			n := len(bp.Children)
			bp.Children = bp.Children[:n+1]
			copy(bp.Children[1:], bp.Children[:n])
			bp.Children[0] = &BpTreeNode{Key: n}
			bp.Children[0].Children = append(bp.Children[0].Children, &BpTreeNode{})
		} else {
			// create new intermediate node
			// reference it to parent
			// add element in leaf node

			// lets keep it like this for now
		}
	} else {
		// TODO: currently not checking leaf node
		bp.Children[index].insertChild(n) // insert child will take care of inserting the child in correct sorted order
	}
}

// PrettyPrint prints the B-tree structure in a visually appealing format
func (bp *BpTreeNode) PrettyPrint() {
	if bp == nil {
		println("Tree is empty")
		return
	}

	printBpTree(bp, "", true)
}

// Helper function to recursively print the tree structure
func printBpTree(node *BpTreeNode, prefix string, isLast bool) {
	if node == nil {
		return
	}

	// Print the current node
	var connector string
	if isLast {
		connector = "└── "
	} else {
		connector = "├── "
	}

	if node.Key != 0 { // Only print non-zero Keys (assuming 0 is placeholder)
		fmt.Printf("%s%s[%d]\n", prefix, connector, node.Key)
	}

	// Print children
	children := node.Children
	if len(children) > 0 {
		// Extend the prefix for children
		var childPrefix string
		if isLast {
			childPrefix = prefix + "    "
		} else {
			childPrefix = prefix + "│   "
		}

		for i, child := range children {
			isLastChild := i == len(children)-1
			printBpTree(child, childPrefix, isLastChild)
		}
	}
}
