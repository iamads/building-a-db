package main

import (
	"fmt"
	"math"
	"testing"
)

func TestFindInternalPredecessor(t *testing.T) {
	tests := []struct {
		name     string
		tree     *BpTreeRootNode
		val      int
		expected int
	}{
		{
			name:     "no internal node",
			tree:     &BpTreeRootNode{Children: []*BpTreeInternalNode{}},
			val:      0,
			expected: -1,
		},
		{
			name:     "Value less than all internal nodes",
			tree:     &BpTreeRootNode{Children: []*BpTreeInternalNode{{Key: 10}, {Key: 20}, {Key: 100}}},
			val:      5,
			expected: -1,
		},
		{
			name:     "Value after the first. node itself",
			tree:     &BpTreeRootNode{Children: []*BpTreeInternalNode{{Key: 10}, {Key: 20}, {Key: 100}}},
			val:      11,
			expected: 0,
		},
		{
			name:     "Value after the middle node",
			tree:     &BpTreeRootNode{Children: []*BpTreeInternalNode{{Key: 10}, {Key: 20}, {Key: 100}}},
			val:      70,
			expected: 1,
		},

		{
			name:     "Value is maxInt",
			tree:     &BpTreeRootNode{Children: []*BpTreeInternalNode{{Key: 10}, {Key: 20}, {Key: 100}}},
			val:      math.MaxInt,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tree.findInternalPredecessor(tt.val)
			if result != tt.expected {
				t.Errorf("findInternalPredecessor(%d) = %d; expected %d: Test case: %s",
					tt.val, result, tt.expected, tt.name)
			}
		})
	}
}

// func BenchmarkFindLeafNodeParent(b *testing.B) {
// 	// Create a larger test case for benchmarking
// 	refs := make([]*BpTreeNode, 100)
// 	for i := 0; i < 100; i++ {
// 		refs[i] = &BpTreeNode{Value: i * 10}
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		findLeafNodeParent(refs, 450) // Test with a value in the middle
// 	}
// }

// func BenchmarkFindLeafNodeParentSmall(b *testing.B) {
// 	// Benchmark with small dataset
// 	refs := []*BpTreeNode{
// 		{Value: 10},
// 		{Value: 20},
// 		{Value: 30},
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		findLeafNodeParent(refs, 25)
// 	}
// }

func TestAddLeafNode(t *testing.T) {
	tests := []struct {
		name        string
		node        *BpTreeInternalNode
		key         int
		val         string
		expectedKey int
		expectedVal string
	}{
		{
			name: "Insert at beginning",
			node: &BpTreeInternalNode{
				Key: 10,
				Children: []*BpTreeLeafNode{
					{Key: 15, Value: "value15"},
					{Key: 20, Value: "value20"},
				},
			},
			key:         5,
			val:         "value5",
			expectedKey: 5,
			expectedVal: "value5",
		},
		{
			name: "Insert in middle",
			node: &BpTreeInternalNode{
				Key: 10,
				Children: []*BpTreeLeafNode{
					{Key: 5, Value: "value5"},
					{Key: 15, Value: "value15"},
				},
			},
			key:         10,
			val:         "value10",
			expectedKey: 10,
			expectedVal: "value10",
		},
		{
			name: "Insert at end",
			node: &BpTreeInternalNode{
				Key: 10,
				Children: []*BpTreeLeafNode{
					{Key: 5, Value: "value5"},
					{Key: 10, Value: "value10"},
				},
			},
			key:         20,
			val:         "value20",
			expectedKey: 20,
			expectedVal: "value20",
		},
		{
			name: "Key equals internal node key",
			node: &BpTreeInternalNode{
				Key: 10,
				Children: []*BpTreeLeafNode{
					{Key: 15, Value: "value15"},
				},
			},
			key:         10,
			val:         "value10",
			expectedKey: 10,
			expectedVal: "value10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.node.addLeafNode(tt.key, tt.val)

			// Check that the children array has the expected length
			expectedLength := len(tt.node.Children)
			if expectedLength != 1 && expectedLength != 2 && expectedLength != 3 {
				t.Errorf("Expected children length to be 1, 2, or 3, got %d", expectedLength)
			}

			// Find the inserted element
			found := false
			for _, child := range tt.node.Children {
				if child.Key == tt.expectedKey && child.Value == tt.expectedVal {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Expected key %d with value %s not found in children", tt.expectedKey, tt.expectedVal)
			}

			// Check that all children are sorted
			for i := 1; i < len(tt.node.Children); i++ {
				if tt.node.Children[i-1].Key > tt.node.Children[i].Key {
					t.Errorf("Children are not sorted: %d > %d at position %d",
						tt.node.Children[i-1].Key, tt.node.Children[i].Key, i)
				}
			}
		})
	}
}

func TestAddLeafNodeEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		node        *BpTreeInternalNode
		key         int
		val         string
		description string
	}{
		{
			name:        "Insert with minimum integer key",
			node:        &BpTreeInternalNode{Key: 0, Children: []*BpTreeLeafNode{}},
			key:         math.MinInt,
			val:         "min_value",
			description: "Testing with minimum possible integer key",
		},
		{
			name:        "Insert with maximum integer key",
			node:        &BpTreeInternalNode{Key: 0, Children: []*BpTreeLeafNode{}},
			key:         math.MaxInt,
			val:         "max_value",
			description: "Testing with maximum possible integer key",
		},
		{
			name: "Insert with existing key",
			node: &BpTreeInternalNode{
				Key: 10,
				Children: []*BpTreeLeafNode{
					{Key: 5, Value: "value5"},
					{Key: 15, Value: "value15"},
				},
			},
			key:         5,
			val:         "new_value5",
			description: "Testing insertion of duplicate key",
		},
		{
			name:        "Insert into node with one child",
			node:        &BpTreeInternalNode{Key: 10, Children: []*BpTreeLeafNode{{Key: 5, Value: "value5"}}},
			key:         10,
			val:         "value10",
			description: "Testing insertion into node that already has one child",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldLength := len(tt.node.Children)
			tt.node.addLeafNode(tt.key, tt.val)
			newLength := len(tt.node.Children)

			// Check that length increased by 1
			if newLength != oldLength+1 {
				t.Errorf("Expected length to increase from %d to %d, got %d", oldLength, oldLength+1, newLength)
			}

			// Check that the inserted element exists
			found := false
			for _, child := range tt.node.Children {
				if child.Key == tt.key && child.Value == tt.val {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Inserted key %d with value %s not found", tt.key, tt.val)
			}

			t.Logf("Test case: %s - %s", tt.name, tt.description)
		})
	}
}

func TestIsFullInternalNode(t *testing.T) {
	inode := BpTreeInternalNode{Key: 5, Children: []*BpTreeLeafNode{
		{Key: 5, Value: "abc"},
		{Key: 6, Value: "abcs"},
		{Key: 7, Value: "xyz"},
	}}

	inode1 := BpTreeInternalNode{Key: 5, Children: []*BpTreeLeafNode{
		{Key: 5, Value: "abc"},
		{Key: 6, Value: "abcs"},
	}}

	inode2 := BpTreeInternalNode{Key: 5, Children: []*BpTreeLeafNode{
		{Key: 7, Value: "xyz"},
	}}

	inode3 := BpTreeInternalNode{Key: 5, Children: []*BpTreeLeafNode{}}
	if inode.isFull() || inode1.isFull() || inode2.isFull() || inode3.isFull() {
		t.Errorf("Internal node is not full but marked as full")
	}

	inode4 := BpTreeInternalNode{Key: 5, Children: []*BpTreeLeafNode{
		{Key: 5, Value: "abc"},
		{Key: 6, Value: "abcs"},
		{Key: 7, Value: "xyz"},
		{Key: 8, Value: "xyzq"},
	}}

	if !inode4.isFull() {
		t.Errorf("Internal ndoe should be marked as full")
	}
}

func TestSplitInode2(t *testing.T) {
	// this should be triggered only when internal node reaches max size

	inode := createNewInternalNode(5)
	for _, leaf := range []struct {
		key int
		val string
	}{
		{key: 5, val: "ABC"},
		{key: 6, val: "ABC1"},
		{key: 7, val: "ABC2"},
		{key: 8, val: "ABC3"},
	} {
		inode.addLeafNode(leaf.key, leaf.val)
	}

	left, right := splitInode2(inode)

	if left.Children[0].Key != 5 {
		t.Errorf("Left's fist child should be 5, it was=(%d)", left.Children[0].Key)
	}

	if left.Key != 5 {
		t.Errorf("Lefts key should be 5, it was %d", left.Key)
	}

	if right.Key != 7 {
		t.Errorf("Rights key should be 7, it was %d", right.Key)
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		name       string
		setupTree  func() *BpTreeRootNode
		insertions []struct {
			key int
			val string
		}
		expectedError bool
		validate      func(*testing.T, *BpTreeRootNode)
	}{
		// Test Case 1: First insertion into empty tree (no internal nodes)
		{
			name: "First insertion into empty tree",
			setupTree: func() *BpTreeRootNode {
				return NewBpTree()
			},
			insertions: []struct {
				key int
				val string
			}{
				{key: 10, val: "value10"},
			},
			expectedError: false,
			validate: func(t *testing.T, tree *BpTreeRootNode) {
				if len(tree.Children) != 1 {
					t.Errorf("Expected 1 internal node, got %d", len(tree.Children))
				}
				if tree.Children[0].Key != 10 {
					t.Errorf("Expected internal node key to be 10, got %d", tree.Children[0].Key)
				}
				if len(tree.Children[0].Children) != 1 {
					t.Errorf("Expected 1 leaf node, got %d", len(tree.Children[0].Children))
				}
				if tree.Children[0].Children[0].Key != 10 || tree.Children[0].Children[0].Value != "value10" {
					t.Errorf("Leaf node has incorrect key/value")
				}
			},
		},

		// // Test Case 2: Insert at the very start (key smaller than all existing)
		{
			name: "Insert at the very start",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			insertions: []struct {
				key int
				val string
			}{
				{key: 5, val: "value5"},
			},
			expectedError: false,
			validate: func(t *testing.T, tree *BpTreeRootNode) {
				// Should create a new internal node at the start
				if len(tree.Children) < 2 {
					t.Errorf("Expected at least 2 internal nodes, got %d", len(tree.Children))
				}
				if tree.Children[0].Key != 5 {
					t.Errorf("Expected first internal node key to be 5, got %d", tree.Children[0].Key)
				}
				// Verify the leaf exists
				found := false
				for _, inode := range tree.Children {
					for _, leaf := range inode.Children {
						if leaf.Key == 5 && leaf.Value == "value5" {
							found = true
							break
						}
					}
				}
				if !found {
					t.Errorf("Inserted key 5 not found in tree")
				}
			},
		},

		// // Test Case 3: Insert into a specific internal node
		{
			name: "Insert into specific internal node",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			insertions: []struct {
				key int
				val string
			}{
				{key: 15, val: "value15"},
			},
			expectedError: false,
			validate: func(t *testing.T, tree *BpTreeRootNode) {
				// Find the internal node that should contain key 15
				found := false
				for _, inode := range tree.Children {
					for _, leaf := range inode.Children {
						if leaf.Key == 15 && leaf.Value == "value15" {
							found = true
							break
						}
					}
				}
				if !found {
					t.Errorf("Key 15 not found in the tree")
				}
			},
		},

		// // Test Case 4: Split internal node when full
		{
			name: "Split internal node when full",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				// Fill a single internal node to MAX_SIZE (4)
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				return tree
			},
			insertions: []struct {
				key int
				val string
			}{
				{key: 25, val: "value25"}, // This should trigger a split
			},
			expectedError: false,
			validate: func(t *testing.T, tree *BpTreeRootNode) {
				// After split, we should have 2 internal nodes
				if len(tree.Children) != 2 {
					t.Errorf("Expected 2 internal nodes after split, got %d", len(tree.Children))
				}
				// Verify all keys are present
				allKeys := []int{10, 20, 25, 30, 40}
				for _, expectedKey := range allKeys {
					found := false
					for _, inode := range tree.Children {
						for _, leaf := range inode.Children {
							if leaf.Key == expectedKey {
								found = true
								break
							}
						}
					}
					if !found {
						t.Errorf("Key %d not found after split", expectedKey)
					}
				}
			},
		},

		// Test Case 5: Error when parent is full and can't split
		{
			name: "Error when parent is full",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				// Create MAX_SIZE internal nodes, each full
				// MAX_SIZE = 4, so we need 4 internal nodes with 4 children each
				for i := 0; i < MAX_SIZE; i++ {
					baseKey := i * 100
					for j := 0; j < MAX_SIZE; j++ {
						key := baseKey + j*10
						tree.Insert(key, fmt.Sprintf("value%d", key))
					}
				}
				return tree
			},
			insertions: []struct {
				key int
				val string
			}{
				{key: 5, val: "value5"}, // This should fail as root is full
			},
			expectedError: true,
			validate: func(t *testing.T, tree *BpTreeRootNode) {
				// Tree should still be at MAX_SIZE

				if len(tree.Children) != MAX_SIZE {
					t.Errorf("Expected %d internal nodes, got %d", MAX_SIZE, len(tree.Children))
				}
			},
		},

		// // Test Case 6: Multiple sequential insertions
		{
			name: "Multiple sequential insertions",
			setupTree: func() *BpTreeRootNode {
				return NewBpTree()
			},
			insertions: []struct {
				key int
				val string
			}{
				{key: 50, val: "value50"},
				{key: 30, val: "value30"},
				{key: 70, val: "value70"},
				{key: 20, val: "value20"},
				{key: 60, val: "value60"},
			},
			expectedError: false,
			validate: func(t *testing.T, tree *BpTreeRootNode) {
				expectedKeys := []int{20, 30, 50, 60, 70}
				for _, expectedKey := range expectedKeys {
					found := false
					for _, inode := range tree.Children {
						for _, leaf := range inode.Children {
							if leaf.Key == expectedKey {
								found = true
								break
							}
						}
					}
					if !found {
						t.Errorf("Key %d not found in tree", expectedKey)
					}
				}
			},
		},

		// // Test Case 7: Insert with minimum and maximum integer values
		{
			name: "Insert with extreme values",
			setupTree: func() *BpTreeRootNode {
				return NewBpTree()
			},
			insertions: []struct {
				key int
				val string
			}{
				{key: 0, val: "value0"},
				{key: math.MinInt, val: "valueMin"},
				{key: math.MaxInt, val: "valueMax"},
			},
			expectedError: false,
			validate: func(t *testing.T, tree *BpTreeRootNode) {
				keys := []int{0, math.MinInt, math.MaxInt}
				for _, expectedKey := range keys {
					found := false
					for _, inode := range tree.Children {
						for _, leaf := range inode.Children {
							if leaf.Key == expectedKey {
								found = true
								break
							}
						}
					}
					if !found {
						t.Errorf("Key %d not found in tree", expectedKey)
					}
				}
			},
		},

		// // Test Case 8: Insert duplicate keys
		{
			name: "Insert duplicate keys",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "original")
				return tree
			},
			insertions: []struct {
				key int
				val string
			}{
				{key: 10, val: "duplicate"},
			},
			expectedError: false,
			validate: func(t *testing.T, tree *BpTreeRootNode) {
				// Count how many times key 10 appears
				count := 0
				for _, inode := range tree.Children {
					for _, leaf := range inode.Children {
						if leaf.Key == 10 {
							count++
						}
					}
				}
				if count != 2 {
					t.Errorf("Expected 2 occurrences of key 10, got %d", count)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()

			var err error
			for _, ins := range tt.insertions {
				err = tree.Insert(ins.key, ins.val)
				if err != nil {
					break
				}
			}

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, tree)
			}
		})
	}
}
