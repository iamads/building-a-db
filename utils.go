package main

import (
	"fmt"
	"strings"
)

// PrettyPrint prints the B+Tree structure in a visual tree format
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

	fmt.Println()

	// Collect all leaf nodes for positioning
	type leafInfo struct {
		key   int
		value string
	}

	var allLeaves []leafInfo
	leafPositions := make(map[int]int) // maps leaf index to its internal node parent index

	leafIdx := 0
	for inodeIdx, inode := range root.Children {
		for _, leaf := range inode.Children {
			allLeaves = append(allLeaves, leafInfo{key: leaf.Key, value: leaf.Value})
			leafPositions[leafIdx] = inodeIdx
			leafIdx++
		}
	}

	// Node width for each leaf
	nodeWidth := 8
	spacing := 2

	// Build leaf node strings
	var leafNodes []string
	for _, leaf := range allLeaves {
		leafNodes = append(leafNodes, fmt.Sprintf("L[%d]", leaf.key))
	}

	// Calculate positions
	totalWidth := len(leafNodes)*nodeWidth + (len(leafNodes)-1)*spacing

	// Print ROOT node
	rootKeys := make([]string, len(root.Children))
	for i, child := range root.Children {
		rootKeys[i] = fmt.Sprintf("%d", child.Key)
	}
	rootStr := fmt.Sprintf("ROOT[%s]", strings.Join(rootKeys, ","))
	rootPadding := (totalWidth - len(rootStr)) / 2
	if rootPadding < 0 {
		rootPadding = 0
	}
	fmt.Printf("%s%s\n", strings.Repeat(" ", rootPadding), rootStr)

	// Print connections from root to internal nodes
	var line1 strings.Builder

	// Calculate internal node positions
	inodePositions := make([]int, len(root.Children))
	for i := range root.Children {
		// Find the middle position of all leaves belonging to this internal node
		firstLeafIdx := -1
		lastLeafIdx := -1
		for lIdx, inodeIdx := range leafPositions {
			if inodeIdx == i {
				if firstLeafIdx == -1 {
					firstLeafIdx = lIdx
				}
				lastLeafIdx = lIdx
			}
		}

		if firstLeafIdx != -1 {
			firstPos := firstLeafIdx * (nodeWidth + spacing)
			lastPos := lastLeafIdx * (nodeWidth + spacing)
			inodePositions[i] = (firstPos + lastPos) / 2
		}
	}

	// Draw connection lines from root to internal nodes
	for pos := 0; pos < totalWidth; pos++ {
		isConnection := false
		connectionChar := " "

		for i := range root.Children {
			if inodePositions[i] == pos {
				if i == 0 {
					connectionChar = "/"
				} else if i == len(root.Children)-1 {
					connectionChar = "\\"
				} else {
					connectionChar = "|"
				}
				isConnection = true
				break
			}
		}

		if isConnection {
			line1.WriteString(connectionChar)
		} else {
			line1.WriteString(" ")
		}
	}
	fmt.Println(line1.String())

	// Print internal nodes
	var inodeLine strings.Builder
	for i, inode := range root.Children {
		nodeStr := fmt.Sprintf("IN[%d]", inode.Key)
		pos := inodePositions[i]

		// Pad to position
		for inodeLine.Len() < pos-(len(nodeStr)/2) {
			inodeLine.WriteString(" ")
		}
		inodeLine.WriteString(nodeStr)
	}
	fmt.Println(inodeLine.String())

	// Print connections from internal nodes to leaf nodes
	var connLine strings.Builder

	for lIdx := 0; lIdx < len(leafNodes); lIdx++ {
		pos := lIdx * (nodeWidth + spacing)
		inodeIdx := leafPositions[lIdx]

		// Pad to position
		for connLine.Len() < pos {
			connLine.WriteString(" ")
		}

		// Determine connection character
		// Count leaves for this internal node
		leavesForInode := 0
		firstLeafOfInode := -1
		for idx, iIdx := range leafPositions {
			if iIdx == inodeIdx {
				leavesForInode++
				if firstLeafOfInode == -1 {
					firstLeafOfInode = idx
				}
			}
		}

		localIdx := lIdx - firstLeafOfInode

		if leavesForInode == 1 {
			connLine.WriteString("|")
		} else if localIdx == 0 {
			connLine.WriteString("/")
		} else if localIdx == leavesForInode-1 {
			connLine.WriteString("\\")
		} else {
			connLine.WriteString("|")
		}
	}
	fmt.Println(connLine.String())

	// Print leaf nodes
	var leafLine strings.Builder
	for i, leafNode := range leafNodes {
		if i > 0 {
			leafLine.WriteString(strings.Repeat(" ", spacing))
		}
		// Pad node to fixed width
		leafLine.WriteString(leafNode)
		padding := nodeWidth - len(leafNode)
		if padding > 0 {
			leafLine.WriteString(strings.Repeat(" ", padding))
		}
	}
	fmt.Println(leafLine.String())

	fmt.Println()
}
