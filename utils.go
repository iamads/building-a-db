package main

import (
	"fmt"
	"strings"
)

// PrettyPrint prints the B+Tree structure in a vertical format
// Root node at top, internal nodes in middle, leaf nodes at bottom
func (root *BpTreeRootNode) PrettyPrint() {
	if root == nil {
		fmt.Println("Tree is empty (nil)")
		return
	}

	if len(root.Children) == 0 {
		fmt.Println("Tree is empty (no children)")
		return
	}

	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│            B+Tree Structure             │")
	fmt.Println("└─────────────────────────────────────────┘")
	fmt.Println()

	// Print Root Node
	fmt.Println("ROOT NODE:")
	fmt.Printf("  Number of internal nodes: %d\n", len(root.Children))

	// Collect internal node keys for display
	internalKeys := make([]int, len(root.Children))
	for i, child := range root.Children {
		internalKeys[i] = child.Key
	}
	fmt.Printf("  Internal node keys: %v\n", internalKeys)
	fmt.Println()

	// Print Internal Nodes
	fmt.Println("INTERNAL NODES:")
	fmt.Println(strings.Repeat("─", 60))

	for i, inode := range root.Children {
		fmt.Printf("\n[Internal Node %d] Key: %d\n", i, inode.Key)
		fmt.Printf("  └─ Leaf count: %d\n", len(inode.Children))

		if len(inode.Children) > 0 {
			fmt.Print("  └─ Leaf keys: [")
			for j, leaf := range inode.Children {
				if j > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%d", leaf.Key)
			}
			fmt.Println("]")
		}
	}

	fmt.Println()

	// Print Leaf Nodes
	fmt.Println("LEAF NODES:")
	fmt.Println(strings.Repeat("─", 60))

	leafCount := 0
	for inodeIdx, inode := range root.Children {
		if len(inode.Children) > 0 {
			fmt.Printf("\n[Internal Node %d - Key: %d] Leaves:\n", inodeIdx, inode.Key)
			for leafIdx, leaf := range inode.Children {
				fmt.Printf("  %d. Key: %-10d Value: %s\n", leafIdx+1, leaf.Key, leaf.Value)
				leafCount++
			}
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("═", 60))
	fmt.Printf("Total: %d internal nodes, %d leaf nodes\n", len(root.Children), leafCount)
	fmt.Println(strings.Repeat("═", 60))
}
