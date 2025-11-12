package main

import (
	"building-a-db/bptree"
	"fmt"
)

func main() {
	// Create a sample B+ tree with multiple internal nodes and leaves
	root := &bptree.BpTreeRootNode{
		Children: []*bptree.BpTreeInternalNode{
			{
				Key: 10,
				Children: []*bptree.BpTreeLeafNode{
					{Key: 5, Value: "value5"},
					{Key: 10, Value: "value10"},
				},
			},
			{
				Key: 20,
				Children: []*bptree.BpTreeLeafNode{
					{Key: 15, Value: "value15"},
					{Key: 20, Value: "value20"},
				},
			},
			{
				Key: 30,
				Children: []*bptree.BpTreeLeafNode{
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

	root2 := &bptree.BpTreeRootNode{
		Children: []*bptree.BpTreeInternalNode{
			{
				Key: 100,
				Children: []*bptree.BpTreeLeafNode{
					{Key: 50, Value: "val50"},
					{Key: 75, Value: "val75"},
					{Key: 100, Value: "val100"},
				},
			},
			{
				Key: 200,
				Children: []*bptree.BpTreeLeafNode{
					{Key: 150, Value: "val150"},
					{Key: 200, Value: "val200"},
				},
			},
		},
	}

	root2.PrettyPrint()
}
