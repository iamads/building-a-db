package main

import "fmt"

func main() {
	// Create a sample B+ tree with multiple internal nodes and leaves
	root := &BpTreeRootNode{
		Children: []*BpTreeInternalNode{
			{
				Key: 10,
				Children: []*BpTreeLeafNode{
					{Key: 5, Value: "value5"},
					{Key: 10, Value: "value10"},
				},
			},
			{
				Key: 20,
				Children: []*BpTreeLeafNode{
					{Key: 15, Value: "value15"},
					{Key: 20, Value: "value20"},
				},
			},
			{
				Key: 30,
				Children: []*BpTreeLeafNode{
					{Key: 25, Value: "value25"},
					{Key: 30, Value: "value30"},
					{Key: 35, Value: "value35"},
				},
			},
		},
	}

	fmt.Println("B+ Tree Visualization:")
	fmt.Println("======================")
	root.PrettyPrint()

	fmt.Println("\nAnother example with more leaves:")
	fmt.Println("===================================")

	root2 := &BpTreeRootNode{
		Children: []*BpTreeInternalNode{
			{
				Key: 100,
				Children: []*BpTreeLeafNode{
					{Key: 50, Value: "val50"},
					{Key: 75, Value: "val75"},
					{Key: 100, Value: "val100"},
				},
			},
			{
				Key: 200,
				Children: []*BpTreeLeafNode{
					{Key: 150, Value: "val150"},
					{Key: 200, Value: "val200"},
				},
			},
		},
	}

	root2.PrettyPrint()
}
