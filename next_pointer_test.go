package main

import (
	"fmt"
	"testing"
)

// ========== NEXT POINTER MANAGEMENT TESTS ==========

// Helper function to traverse the entire leaf chain via Next pointers
func traverseLeafChain(tree *BpTreeRootNode) []*BpTreeLeafNode {
	if len(tree.Children) == 0 {
		return []*BpTreeLeafNode{}
	}

	var result []*BpTreeLeafNode
	current := tree.Children[0].Children[0] // Start from first leaf of first internal node

	for current != nil {
		result = append(result, current)
		current = current.Next
	}

	return result
}

// Helper function to verify Next pointer chain is correct
func verifyNextChain(t *testing.T, tree *BpTreeRootNode, expectedKeys []int) {
	chain := traverseLeafChain(tree)

	if len(chain) != len(expectedKeys) {
		t.Errorf("Chain length mismatch: expected %d, got %d", len(expectedKeys), len(chain))
		return
	}

	for i, leaf := range chain {
		if leaf.Key != expectedKeys[i] {
			t.Errorf("At position %d: expected key %d, got %d", i, expectedKeys[i], leaf.Key)
		}
	}

	// Verify last node's Next is nil
	if len(chain) > 0 && chain[len(chain)-1].Next != nil {
		t.Errorf("Last leaf node's Next should be nil, but it's not")
	}
}

// TestNextPointer_HappyPath tests basic Next pointer functionality
func TestNextPointer_HappyPath(t *testing.T) {
	tests := []struct {
		name       string
		insertions []struct {
			key int
			val string
		}
		expectedKeys []int
		description  string
	}{
		{
			name: "Single insertion - Next should be nil",
			insertions: []struct {
				key int
				val string
			}{
				{10, "value10"},
			},
			expectedKeys: []int{10},
			description:  "First and only leaf node should have Next = nil",
		},
		{
			name: "Three sequential insertions",
			insertions: []struct {
				key int
				val string
			}{
				{10, "value10"},
				{20, "value20"},
				{30, "value30"},
			},
			expectedKeys: []int{10, 20, 30},
			description:  "Chain should be: 10→20→30→nil",
		},
		{
			name: "Insertions in reverse order",
			insertions: []struct {
				key int
				val string
			}{
				{30, "value30"},
				{20, "value20"},
				{10, "value10"},
			},
			expectedKeys: []int{10, 20, 30},
			description:  "Despite reverse insertion, chain should be: 10→20→30→nil",
		},
		{
			name: "Insertions in random order",
			insertions: []struct {
				key int
				val string
			}{
				{20, "value20"},
				{10, "value10"},
				{40, "value40"},
				{30, "value30"},
			},
			expectedKeys: []int{10, 20, 30, 40},
			description:  "Chain should be sorted: 10→20→30→40→nil",
		},
		{
			name: "Five insertions filling one internal node",
			insertions: []struct {
				key int
				val string
			}{
				{10, "value10"},
				{20, "value20"},
				{30, "value30"},
				{40, "value40"},
			},
			expectedKeys: []int{10, 20, 30, 40},
			description:  "Single internal node with MAX_SIZE elements",
		},
		{
			name: "Insertions with duplicates",
			insertions: []struct {
				key int
				val string
			}{
				{10, "value10_1"},
				{10, "value10_2"},
				{20, "value20"},
				{10, "value10_3"},
			},
			expectedKeys: []int{10, 10, 10, 20},
			description:  "Duplicate keys should maintain Next chain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := NewBpTree()

			for _, ins := range tt.insertions {
				err := tree.Insert(ins.key, ins.val)
				if err != nil {
					t.Fatalf("Insertion failed: %v", err)
				}
			}

			verifyNextChain(t, tree, tt.expectedKeys)
			t.Logf("%s: PASSED", tt.description)
		})
	}
}

// TestNextPointer_BasicUseCases tests common insertion patterns
func TestNextPointer_BasicUseCases(t *testing.T) {
	tests := []struct {
		name         string
		setupTree    func() *BpTreeRootNode
		expectedKeys []int
		description  string
	}{
		{
			name: "Insert at beginning of internal node",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(10, "value10") // Insert at beginning
				return tree
			},
			expectedKeys: []int{10, 20, 30},
			description:  "Inserting at beginning should update Next pointers correctly",
		},
		{
			name: "Insert in middle of internal node",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(30, "value30")
				tree.Insert(20, "value20") // Insert in middle
				return tree
			},
			expectedKeys: []int{10, 20, 30},
			description:  "Inserting in middle should maintain Next chain",
		},
		{
			name: "Insert at end of internal node",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30") // Insert at end
				return tree
			},
			expectedKeys: []int{10, 20, 30},
			description:  "Inserting at end should extend Next chain",
		},
		{
			name: "Multiple internal nodes - cross boundary linking",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				// Insert enough to create multiple internal nodes
				for i := 0; i < 10; i++ {
					tree.Insert(i*10, fmt.Sprintf("value%d", i*10))
				}
				return tree
			},
			expectedKeys: []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90},
			description:  "Next pointers should link across internal node boundaries",
		},
		{
			name: "Insert triggering new internal node at beginning",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(50, "value50")
				tree.Insert(60, "value60")
				tree.Insert(70, "value70")
				tree.Insert(10, "value10") // Should create new internal node at beginning
				return tree
			},
			expectedKeys: []int{10, 50, 60, 70},
			description:  "Creating new internal node at start should link to existing nodes",
		},
		{
			name: "Sequential insertions creating multiple internal nodes",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				for i := 1; i <= 15; i++ {
					tree.Insert(i, fmt.Sprintf("value%d", i))
				}
				return tree
			},
			expectedKeys: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			description:  "Large sequential insertions should maintain full Next chain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()
			verifyNextChain(t, tree, tt.expectedKeys)
			t.Logf("%s: PASSED", tt.description)
		})
	}
}

// TestNextPointer_EdgeCases tests boundary conditions
func TestNextPointer_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		setupTree    func() *BpTreeRootNode
		expectedKeys []int
		description  string
	}{
		{
			name: "Insert when inodeidx == -1 (before all keys)",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(50, "value50")
				tree.Insert(60, "value60")
				tree.Insert(5, "value5") // This triggers inodeidx == -1
				return tree
			},
			expectedKeys: []int{5, 50, 60},
			description:  "Insert before all keys should create new internal node and link properly",
		},
		{
			name: "Node split scenario",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				// Fill one internal node to MAX_SIZE
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				tree.Insert(50, "value50")
				// This should trigger a split
				tree.Insert(25, "value25")
				return tree
			},
			expectedKeys: []int{10, 20, 25, 30, 40, 50},
			description:  "Node split should preserve Next pointer chain",
		},
		{
			name: "Insert at exact boundary between internal nodes",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				// Create scenario with multiple internal nodes
				for i := 0; i < 8; i++ {
					tree.Insert(i*10, fmt.Sprintf("value%d", i*10))
				}
				// Insert at boundary
				tree.Insert(35, "value35")
				return tree
			},
			expectedKeys: []int{0, 10, 20, 30, 35, 40, 50, 60, 70},
			description:  "Insert at internal node boundary should maintain cross-boundary links",
		},
		{
			name: "Multiple duplicates at different positions",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10_1")
				tree.Insert(20, "value20_1")
				tree.Insert(10, "value10_2")
				tree.Insert(30, "value30_1")
				tree.Insert(20, "value20_2")
				tree.Insert(10, "value10_3")
				return tree
			},
			expectedKeys: []int{10, 10, 10, 20, 20, 30},
			description:  "Multiple duplicates should maintain sorted Next chain",
		},
		{
			name: "Insert causing multiple internal node creations",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				// Insert in pattern that creates multiple internal nodes
				keys := []int{100, 200, 10, 20, 150, 175, 50, 75}
				for _, key := range keys {
					tree.Insert(key, fmt.Sprintf("value%d", key))
				}
				return tree
			},
			expectedKeys: []int{10, 20, 50, 75, 100, 150, 175, 200},
			description:  "Complex insertion pattern should maintain complete Next chain",
		},
		{
			name: "Negative and positive keys",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(0, "value0")
				tree.Insert(-10, "valueMinus10")
				tree.Insert(10, "value10")
				tree.Insert(-20, "valueMinus20")
				tree.Insert(20, "value20")
				return tree
			},
			expectedKeys: []int{-20, -10, 0, 10, 20},
			description:  "Negative and positive keys should maintain sorted Next chain",
		},
		{
			name: "Large dataset - verify complete traversal",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				// TODO just for single level
				for i := 0; i < MAX_SIZE*MAX_SIZE; i++ {
					tree.Insert(i*5, fmt.Sprintf("value%d", i*5))
				}
				return tree
			},
			expectedKeys: []int{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75},
			description:  "Large dataset should have complete Next chain traversable from start to end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()
			verifyNextChain(t, tree, tt.expectedKeys)
			t.Logf("%s: PASSED", tt.description)
		})
	}
}

// TestNextPointer_SpecificScenarios tests specific implementation requirements
func TestNextPointer_SpecificScenarios(t *testing.T) {
	t.Run("Empty tree - no Next pointers", func(t *testing.T) {
		tree := NewBpTree()
		chain := traverseLeafChain(tree)
		if len(chain) != 0 {
			t.Errorf("Empty tree should have no leaf nodes in chain")
		}
	})

	t.Run("First insertion - Next should be nil", func(t *testing.T) {
		tree := NewBpTree()
		tree.Insert(10, "value10")

		if len(tree.Children) == 0 || len(tree.Children[0].Children) == 0 {
			t.Fatal("Tree should have one internal node with one leaf")
		}

		firstLeaf := tree.Children[0].Children[0]
		if firstLeaf.Next != nil {
			t.Errorf("First and only leaf's Next should be nil, got %v", firstLeaf.Next)
		}
	})

	t.Run("Cross-boundary Next pointer verification", func(t *testing.T) {
		tree := NewBpTree()
		// Insert enough to create at least 2 internal nodes
		for i := 0; i < 10; i++ {
			tree.Insert(i*10, fmt.Sprintf("value%d", i*10))
		}

		// Find where internal nodes split and verify cross-boundary link
		if len(tree.Children) < 2 {
			t.Skip("Need at least 2 internal nodes for this test")
		}

		for i := 0; i < len(tree.Children)-1; i++ {
			currentInode := tree.Children[i]
			nextInode := tree.Children[i+1]

			lastLeafOfCurrent := currentInode.Children[len(currentInode.Children)-1]
			firstLeafOfNext := nextInode.Children[0]

			if lastLeafOfCurrent.Next != firstLeafOfNext {
				t.Errorf("Internal node %d's last leaf should point to internal node %d's first leaf", i, i+1)
				t.Errorf("Expected Next=%v, got Next=%v", firstLeafOfNext, lastLeafOfCurrent.Next)
			}
		}
	})

	t.Run("Insert at position 0 updates previous internal node", func(t *testing.T) {
		tree := NewBpTree()
		// Create multiple internal nodes
		for i := 0; i < 10; i++ {
			tree.Insert(i*10, fmt.Sprintf("value%d", i*10))
		}

		if len(tree.Children) < 2 {
			t.Skip("Need at least 2 internal nodes for this test")
		}

		// Get the first leaf of second internal node
		secondInodeFirstKey := tree.Children[1].Children[0].Key

		// Insert a key that goes to position 0 of second internal node
		// This should update the Next pointer of last leaf of first internal node
		newKey := secondInodeFirstKey - 1
		tree.Insert(newKey, fmt.Sprintf("value%d", newKey))

		// Verify the chain is still correct
		chain := traverseLeafChain(tree)
		for i := 0; i < len(chain)-1; i++ {
			if chain[i].Next != chain[i+1] {
				t.Errorf("Chain broken at position %d: key %d's Next should be key %d",
					i, chain[i].Key, chain[i+1].Key)
			}
		}
	})
}

// TestNextPointer_IdentifiedBugs tests for specific bugs found in implementation
func TestNextPointer_IdentifiedBugs(t *testing.T) {
	t.Run("safelyManageNextBoundaryLeafNodes bug - replaces child instead of updating Next", func(t *testing.T) {
		tree := NewBpTree()

		// Insert keys that will create at least 2 internal nodes
		tree.Insert(50, "value50")
		tree.Insert(60, "value60")
		tree.Insert(70, "value70")
		tree.Insert(80, "value80")

		// Insert at beginning to trigger inodeidx == -1
		tree.Insert(10, "value10")

		if len(tree.Children) < 2 {
			t.Logf("Created %d internal nodes", len(tree.Children))
		}

		// BUG: safelyManageNextBoundaryLeafNodes (line 108) does:
		// prev_internal_node.Children[len(prev_internal_node.Children)-1] = leafNode
		// This REPLACES the last child instead of updating its Next pointer

		// Verify the chain
		chain := traverseLeafChain(tree)
		expectedKeys := []int{10, 50, 60, 70, 80}

		if len(chain) != len(expectedKeys) {
			t.Errorf("BUG DETECTED: Expected %d nodes in chain, got %d", len(expectedKeys), len(chain))
			t.Errorf("Chain keys: %v", getKeysFromChain(chain))
		}

		for i := 0; i < len(chain)-1; i++ {
			if chain[i].Next != chain[i+1] {
				t.Errorf("BUG DETECTED: Chain broken at position %d", i)
				t.Errorf("Key %d's Next should point to key %d, but it doesn't",
					chain[i].Key, chain[i+1].Key)
			}
		}
	})

	t.Run("addLeafNode doesn't set Next when inserting at position 0", func(t *testing.T) {
		tree := NewBpTree()

		// First insert
		tree.Insert(20, "value20")
		tree.Insert(30, "value30")

		// Insert at beginning (position 0)
		tree.Insert(10, "value10")

		// BUG: addLeafNode only handles toInsertIdx > 0 case (lines 251-255)
		// When toInsertIdx == 0, Next pointer is not set

		firstLeaf := tree.Children[0].Children[0]
		if firstLeaf.Key == 10 && firstLeaf.Next == nil {
			t.Errorf("BUG DETECTED: First leaf (key=10) should have Next pointing to key=20, but Next is nil")
		}

		verifyNextChain(t, tree, []int{10, 20, 30})
	})

	t.Run("Node split doesn't preserve Next pointers", func(t *testing.T) {
		tree := NewBpTree()

		// Fill to trigger a split
		tree.Insert(10, "value10")
		tree.Insert(20, "value20")
		tree.Insert(30, "value30")
		tree.Insert(40, "value40")
		tree.Insert(25, "value25") // Should trigger split

		// BUG: splitInode2 doesn't handle Next pointers
		// The last leaf of 'original' should point to first leaf of 'next'

		chain := traverseLeafChain(tree)
		expectedKeys := []int{10, 20, 25, 30, 40}

		if len(chain) != len(expectedKeys) {
			t.Errorf("BUG DETECTED: After split, expected %d nodes, got %d", len(expectedKeys), len(chain))
		}

		// Check for broken chain
		for i := 0; i < len(chain)-1; i++ {
			if chain[i].Next != chain[i+1] {
				t.Errorf("BUG DETECTED: Chain broken after split at position %d", i)
				t.Errorf("Key %d (Next=%v) should point to key %d",
					chain[i].Key, chain[i].Next, chain[i+1].Key)
			}
		}
	})

	t.Run("Insert into second internal node at position 0 doesn't update previous inode", func(t *testing.T) {
		tree := NewBpTree()

		// Create scenario with 2 internal nodes
		for i := 0; i < 8; i++ {
			tree.Insert(i*10, fmt.Sprintf("value%d", i*10))
		}

		if len(tree.Children) < 2 {
			t.Skip("Need at least 2 internal nodes")
		}

		// Get boundary info
		firstInodeLastLeaf := tree.Children[0].Children[len(tree.Children[0].Children)-1]
		secondInodeFirstLeaf := tree.Children[1].Children[0]

		t.Logf("First inode last leaf: key=%d", firstInodeLastLeaf.Key)
		t.Logf("Second inode first leaf: key=%d", secondInodeFirstLeaf.Key)

		// Insert at position 0 of second internal node
		newKey := secondInodeFirstLeaf.Key - 1
		tree.Insert(newKey, fmt.Sprintf("value%d", newKey))

		// BUG: When inserting at position 0 of non-first internal node,
		// the previous internal node's last leaf's Next should be updated

		// Verify complete chain
		chain := traverseLeafChain(tree)
		for i := 0; i < len(chain)-1; i++ {
			if chain[i].Next != chain[i+1] {
				t.Errorf("BUG DETECTED: Chain broken at position %d (key %d → key %d)",
					i, chain[i].Key, chain[i+1].Key)
				if chain[i].Next != nil {
					t.Errorf("Currently points to key %d instead", chain[i].Next.Key)
				} else {
					t.Errorf("Currently points to nil")
				}
			}
		}
	})

	t.Run("Empty internal node children scenario", func(t *testing.T) {
		// This tests the case where addLeafNode might be called on a fresh internal node
		tree := NewBpTree()
		tree.Insert(10, "value10")

		firstLeaf := tree.Children[0].Children[0]

		// First leaf should have Next = nil
		if firstLeaf.Next != nil {
			t.Errorf("BUG DETECTED: First leaf in empty tree should have Next=nil, got Next=%v", firstLeaf.Next)
		}
	})

	t.Run("Cross-boundary link not established when creating new first internal node", func(t *testing.T) {
		tree := NewBpTree()

		// Create existing internal node
		tree.Insert(50, "value50")
		tree.Insert(60, "value60")

		// Insert with inodeidx == -1 (creates new internal node at beginning)
		tree.Insert(10, "value10")

		// BUG: When inodeidx == -1, new internal node is created at position 0
		// The new leaf node should have Next pointing to first leaf of what is now tree.Children[1]

		if len(tree.Children) >= 2 {
			newFirstLeaf := tree.Children[0].Children[0]
			oldFirstLeaf := tree.Children[1].Children[0]

			if newFirstLeaf.Next != oldFirstLeaf {
				t.Errorf("BUG DETECTED: New first leaf (key=%d) should have Next pointing to old first leaf (key=%d)",
					newFirstLeaf.Key, oldFirstLeaf.Key)
				t.Errorf("But Next is: %v", newFirstLeaf.Next)
			}
		}

		verifyNextChain(t, tree, []int{10, 50, 60})
	})
}

// Helper to get keys from chain for debugging
func getKeysFromChain(chain []*BpTreeLeafNode) []int {
	keys := make([]int, len(chain))
	for i, leaf := range chain {
		keys[i] = leaf.Key
	}
	return keys
}
