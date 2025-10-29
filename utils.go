package main

// import "fmt"

// // PrettyPrint prints the B-tree structure in a visually appealing format
// func (bp *BpTreeNode) PrettyPrint() {
// 	if bp == nil {
// 		println("Tree is empty")
// 		return
// 	}

// 	printBpTree(bp, "", true)
// }

// // Helper function to recursively print the tree structure
// func printBpTree(node *BpTreeNode, prefix string, isLast bool) {
// 	if node == nil {
// 		return
// 	}

// 	// Print the current node
// 	var connector string
// 	if isLast {
// 		connector = "└── "
// 	} else {
// 		connector = "├── "
// 	}

// 	if node.Key != 0 { // Only print non-zero Keys (assuming 0 is placeholder)
// 		fmt.Printf("%s%s[%d]\n", prefix, connector, node.Key)
// 	}

// 	// Print children
// 	children := node.Children
// 	if len(children) > 0 {
// 		// Extend the prefix for children
// 		var childPrefix string
// 		if isLast {
// 			childPrefix = prefix + "    "
// 		} else {
// 			childPrefix = prefix + "│   "
// 		}

// 		for i, child := range children {
// 			isLastChild := i == len(children)-1
// 			printBpTree(child, childPrefix, isLastChild)
// 		}
// 	}
// }
