package bptree

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

// DumpStruct prints the exported fields of any struct (recursively).

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
				Key: 15,
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
				Key: 5,
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
				Key: 5,
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
				Key: 15,
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

// ========== GET METHOD TESTS ==========

// TestGet_HappyPath tests successful retrieval scenarios
func TestGet_HappyPath(t *testing.T) {
	tests := []struct {
		name          string
		setupTree     func() *BpTreeRootNode
		key           int
		expectedValue string
		expectError   bool
	}{
		{
			name: "Get single key from simple tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				return tree
			},
			key:           10,
			expectedValue: "value10",
			expectError:   false,
		},
		{
			name: "Get key from tree with multiple insertions",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:           20,
			expectedValue: "value20",
			expectError:   false,
		},
		{
			name: "Get first key in tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:           10,
			expectedValue: "value10",
			expectError:   false,
		},
		{
			name: "Get last key in tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:           30,
			expectedValue: "value30",
			expectError:   false,
		},
		{
			name: "Get key from tree with multiple internal nodes",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				// Create enough insertions to trigger multiple internal nodes
				for i := 0; i < 10; i++ {
					tree.Insert(i*10, fmt.Sprintf("value%d", i*10))
				}
				return tree
			},
			key:           50,
			expectedValue: "value50",
			expectError:   false,
		},
		{
			name: "Get key after node split",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				tree.Insert(25, "value25") // Triggers split
				return tree
			},
			key:           25,
			expectedValue: "value25",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()
			value, err := tree.Get(tt.key)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError && value != tt.expectedValue {
				t.Errorf("Expected value %s, got %s", tt.expectedValue, value)
			}
		})
	}
}

// TestGet_BasicIssues tests common error scenarios
func TestGet_BasicIssues(t *testing.T) {
	tests := []struct {
		name        string
		setupTree   func() *BpTreeRootNode
		key         int
		expectError bool
		description string
	}{
		{
			name: "Get from empty tree",
			setupTree: func() *BpTreeRootNode {
				return NewBpTree()
			},
			key:         10,
			expectError: true,
			description: "Should error when getting from empty tree",
		},
		{
			name: "Get non-existent key from populated tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:         25,
			expectError: true,
			description: "Should error when key doesn't exist",
		},
		{
			name: "Get key smaller than all keys in tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:         5,
			expectError: true,
			description: "Should error when key is smaller than all keys",
		},
		{
			name: "Get key larger than all keys in tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:         40,
			expectError: true,
			description: "Should error when key is larger than all keys",
		},
		{
			name: "Get from tree with single internal node",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				return tree
			},
			key:         20,
			expectError: true,
			description: "Should error when key doesn't exist in single internal node",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()
			value, err := tree.Get(tt.key)

			if tt.expectError && err == nil {
				t.Errorf("%s: Expected error but got none, value=%s", tt.description, value)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: Unexpected error: %v", tt.description, err)
			}
		})
	}
}

// TestGet_EdgeCases tests edge case scenarios
func TestGet_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		setupTree     func() *BpTreeRootNode
		key           int
		expectedValue string
		expectError   bool
		description   string
	}{
		{
			name: "Get with math.MinInt key",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(math.MinInt, "minValue")
				tree.Insert(0, "value0")
				tree.Insert(100, "value100")
				return tree
			},
			key:           math.MinInt,
			expectedValue: "minValue",
			expectError:   false,
			description:   "Should retrieve minimum integer key",
		},
		{
			name: "Get with math.MaxInt key",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(0, "value0")
				tree.Insert(100, "value100")
				tree.Insert(math.MaxInt, "maxValue")
				return tree
			},
			key:           math.MaxInt,
			expectedValue: "maxValue",
			expectError:   false,
			description:   "Should retrieve maximum integer key",
		},
		{
			name: "Get when internal nodes are at MAX_SIZE",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				return tree
			},
			key:           30,
			expectedValue: "value30",
			expectError:   false,
			description:   "Should retrieve key when internal node is at MAX_SIZE",
		},
		{
			name: "Get zero key",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(-10, "valueMinus10")
				tree.Insert(0, "value0")
				tree.Insert(10, "value10")
				return tree
			},
			key:           0,
			expectedValue: "value0",
			expectError:   false,
			description:   "Should retrieve zero key",
		},
		{
			name: "Get negative key",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(-50, "valueMinus50")
				tree.Insert(-10, "valueMinus10")
				tree.Insert(10, "value10")
				return tree
			},
			key:           -10,
			expectedValue: "valueMinus10",
			expectError:   false,
			description:   "Should retrieve negative key",
		},
		{
			name: "Get with duplicate keys - retrieves first occurrence", // TODO: Fix failure here
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10_first")
				tree.Insert(10, "value10_second")
				tree.Insert(10, "value10_third")
				tree.Insert(20, "value20")
				return tree
			},
			key:           10,
			expectedValue: "value10_first",
			expectError:   false,
			description:   "Should retrieve first occurrence of duplicate key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()
			value, err := tree.Get(tt.key)

			if tt.expectError && err == nil {
				t.Errorf("%s: Expected error but got none, value=%s", tt.description, value)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: Unexpected error: %v", tt.description, err)
			}
			if !tt.expectError && value != tt.expectedValue {
				t.Errorf("%s: Expected value %s, got %s", tt.description, tt.expectedValue, value)
			}
		})
	}
}

// ========== SEARCH METHOD TESTS ==========

// TestInternalNodeSearch_HappyPath tests the search method on internal nodes
func TestInternalNodeSearch_HappyPath(t *testing.T) {
	tests := []struct {
		name          string
		node          *BpTreeInternalNode
		searchKey     int
		expectedValue string
		expectError   bool
	}{
		{
			name: "Search existing key in middle",
			node: &BpTreeInternalNode{
				Key: 10,
				Children: []*BpTreeLeafNode{
					{Key: 10, Value: "value10"},
					{Key: 20, Value: "value20"},
					{Key: 30, Value: "value30"},
				},
			},
			searchKey:     20,
			expectedValue: "value20",
			expectError:   false,
		},
		{
			name: "Search first key",
			node: &BpTreeInternalNode{
				Key: 10,
				Children: []*BpTreeLeafNode{
					{Key: 10, Value: "value10"},
					{Key: 20, Value: "value20"},
				},
			},
			searchKey:     10,
			expectedValue: "value10",
			expectError:   false,
		},
		{
			name: "Search last key",
			node: &BpTreeInternalNode{
				Key: 10,
				Children: []*BpTreeLeafNode{
					{Key: 10, Value: "value10"},
					{Key: 20, Value: "value20"},
					{Key: 30, Value: "value30"},
				},
			},
			searchKey:     30,
			expectedValue: "value30",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leaf, err := tt.node.search(tt.searchKey)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError && leaf.Value != tt.expectedValue {
				t.Errorf("Expected value %s, got %s", tt.expectedValue, leaf.Value)
			}
		})
	}
}

// ========== GETRANGE METHOD TESTS ==========

// TestGetRange_HappyPath tests successful range query scenarios
func TestGetRange_HappyPath(t *testing.T) {
	tests := []struct {
		name           string
		setupTree      func() *BpTreeRootNode
		start          int
		end            int
		expectedValues []string
		expectError    bool
	}{
		{
			name: "Basic range query within single internal node",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				return tree
			},
			start:          20,
			end:            40,
			expectedValues: []string{"value20", "value30"},
			expectError:    false,
		},
		{
			name: "Range query spanning multiple internal nodes",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				// Create enough insertions to span multiple internal nodes
				for i := 0; i < 10; i++ {
					tree.Insert(i*10, fmt.Sprintf("value%d", i*10))
				}
				return tree
			},
			start:          20,
			end:            70,
			expectedValues: []string{"value20", "value30", "value40", "value50", "value60"},
			expectError:    false,
		},
		{
			name: "Range from tree start",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				tree.Insert(50, "value50")
				return tree
			},
			start:          10,
			end:            35,
			expectedValues: []string{"value10", "value20", "value30"},
			expectError:    false,
		},
		{
			name: "Range to tree end",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				tree.Insert(50, "value50")
				return tree
			},
			start:          30,
			end:            100,
			expectedValues: []string{"value30", "value40", "value50"},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()
			result, err := tree.GetRange(tt.start, tt.end)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError {
				if len(result) != len(tt.expectedValues) {
					t.Errorf("Expected %d values, got %d. Expected: %v, Got: %v",
						len(tt.expectedValues), len(result), tt.expectedValues, result)
				} else {
					for i, expected := range tt.expectedValues {
						if result[i] != expected {
							t.Errorf("At index %d: expected %s, got %s", i, expected, result[i])
						}
					}
				}
			}
		})
	}
}

// TestGetRange_BasicUseCases tests common use case scenarios
func TestGetRange_BasicUseCases(t *testing.T) {
	tests := []struct {
		name           string
		setupTree      func() *BpTreeRootNode
		start          int
		end            int
		expectedValues []string
		expectError    bool
		description    string
	}{
		{
			name: "Single element range",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			start:          20,
			end:            21,
			expectedValues: []string{"value20"},
			expectError:    false,
			description:    "Range that includes only one element",
		},
		{
			name: "Full tree range",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				return tree
			},
			start:          10,
			end:            50,
			expectedValues: []string{"value10", "value20", "value30", "value40"},
			expectError:    false,
			description:    "Range covering entire tree",
		},
		{
			name: "Range with duplicate keys",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10_first")
				tree.Insert(10, "value10_second")
				tree.Insert(20, "value20")
				tree.Insert(20, "value20_second")
				tree.Insert(30, "value30")
				return tree
			},
			start:          10,
			end:            25,
			expectedValues: []string{"value10_first", "value10_second", "value20", "value20_second"},
			expectError:    false,
			description:    "Verify all duplicates in range are returned",
		},
		{
			name: "Range with gaps in keys",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(30, "value30")
				tree.Insert(50, "value50")
				return tree
			},
			start:          10,
			end:            60,
			expectedValues: []string{"value10", "value30", "value50"},
			expectError:    false,
			description:    "Keys with gaps between them, verify all in range returned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()
			result, err := tree.GetRange(tt.start, tt.end)

			if tt.expectError && err == nil {
				t.Errorf("%s: Expected error but got none", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: Unexpected error: %v", tt.description, err)
			}
			if !tt.expectError {
				if len(result) != len(tt.expectedValues) {
					t.Errorf("%s: Expected %d values, got %d. Expected: %v, Got: %v",
						tt.description, len(tt.expectedValues), len(result), tt.expectedValues, result)
				} else {
					for i, expected := range tt.expectedValues {
						if result[i] != expected {
							t.Errorf("%s: At index %d: expected %s, got %s",
								tt.description, i, expected, result[i])
						}
					}
				}
			}
		})
	}
}

// TestGetRange_EdgeCases tests edge case scenarios
func TestGetRange_EdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		setupTree      func() *BpTreeRootNode
		start          int
		end            int
		expectedValues []string
		expectError    bool
		description    string
	}{
		{
			name: "Empty range (start == end)",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			start:          20,
			end:            20,
			expectedValues: []string{},
			expectError:    false,
			description:    "GetRange(20, 20) should return empty array",
		},
		{
			name: "Invalid range (start > end)",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			start:          30,
			end:            20,
			expectedValues: nil,
			expectError:    true,
			description:    "GetRange(30, 20) should return error",
		},
		{
			name: "Range completely outside tree bounds (both ends)",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			start:          100,
			end:            200,
			expectedValues: []string{},
			expectError:    false,
			description:    "Range [100, 200) where all keys < 100 should return empty",
		},
		{
			name: "Range outside tree bounds (start only)",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			start:          -100,
			end:            25,
			expectedValues: []string{"value10", "value20"},
			expectError:    false,
			description:    "GetRange(-100, 25) where min key is 10 should return values from 10",
		},
		{
			name: "Range outside tree bounds (end only)",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			start:          20,
			end:            1000,
			expectedValues: []string{"value20", "value30"},
			expectError:    false,
			description:    "GetRange(20, 1000) where max key is 30 should return to actual end",
		},
		{
			name: "Empty tree",
			setupTree: func() *BpTreeRootNode {
				return NewBpTree()
			},
			start:          10,
			end:            20,
			expectedValues: nil,
			expectError:    true,
			description:    "GetRange on empty tree should return error",
		},
		{
			name: "Single element tree - range includes element",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(15, "value15")
				return tree
			},
			start:          10,
			end:            20,
			expectedValues: []string{"value15"},
			expectError:    false,
			description:    "Single element tree with range that includes it",
		},
		{
			name: "Single element tree - range excludes element (before)",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(15, "value15")
				return tree
			},
			start:          20,
			end:            30,
			expectedValues: []string{},
			expectError:    false,
			description:    "Single element tree with range after element",
		},
		{
			name: "Single element tree - range excludes element (after)",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(15, "value15")
				return tree
			},
			start:          5,
			end:            10,
			expectedValues: []string{},
			expectError:    false,
			description:    "Single element tree with range before element",
		},
		{
			name: "Extreme values - MinInt to MaxInt",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(math.MinInt, "valueMin")
				tree.Insert(0, "value0")
				tree.Insert(100, "value100")
				tree.Insert(math.MaxInt, "valueMax")
				return tree
			},
			start:          math.MinInt,
			end:            math.MaxInt,
			expectedValues: []string{"valueMin", "value0", "value100"},
			expectError:    false,
			description:    "GetRange from MinInt to MaxInt (excludes MaxInt)",
		},
		{
			name: "Extreme values - range including MaxInt",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(0, "value0")
				tree.Insert(100, "value100")
				tree.Insert(math.MaxInt, "valueMax")
				return tree
			},
			start:          100,
			end:            math.MaxInt,
			expectedValues: []string{"value100"},
			expectError:    false,
			description:    "Range [100, MaxInt) excludes MaxInt key itself",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()
			result, err := tree.GetRange(tt.start, tt.end)

			if tt.expectError && err == nil {
				t.Errorf("%s: Expected error but got none, result=%v", tt.description, result)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: Unexpected error: %v", tt.description, err)
			}
			if !tt.expectError && tt.expectedValues != nil {
				if len(result) != len(tt.expectedValues) {
					t.Errorf("%s: Expected %d values, got %d. Expected: %v, Got: %v",
						tt.description, len(tt.expectedValues), len(result), tt.expectedValues, result)
				} else {
					for i, expected := range tt.expectedValues {
						if result[i] != expected {
							t.Errorf("%s: At index %d: expected %s, got %s",
								tt.description, i, expected, result[i])
						}
					}
				}
			}
		})
	}
}

// TestGetRange_BugSpecific tests the specific bug where starting from middle of internal node fails
func TestGetRange_BugSpecific(t *testing.T) {
	tests := []struct {
		name           string
		setupTree      func() *BpTreeRootNode
		start          int
		end            int
		expectedValues []string
		expectError    bool
		description    string
	}{
		{
			name: "Start key in middle of internal node (missing break bug)",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(5, "value5")
				tree.Insert(10, "value10")
				tree.Insert(15, "value15")
				tree.Insert(20, "value20")
				return tree
			},
			start:          10,
			end:            25,
			expectedValues: []string{"value10", "value15", "value20"},
			expectError:    false,
			description:    "Bug in b+tree.go:334-338 - missing break causes loop to return LAST qualifying node instead of FIRST",
		},
		{
			name: "Start key matches second element in internal node",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			start:          20,
			end:            35,
			expectedValues: []string{"value20", "value30"},
			expectError:    false,
			description:    "Start from second element in internal node",
		},
		{
			name: "Start key matches third element in internal node",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				return tree
			},
			start:          30,
			end:            45,
			expectedValues: []string{"value30", "value40"},
			expectError:    false,
			description:    "Start from third element when internal node has 4 children",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()
			result, err := tree.GetRange(tt.start, tt.end)

			if tt.expectError && err == nil {
				t.Errorf("%s: Expected error but got none, result=%v", tt.description, result)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: Unexpected error: %v", tt.description, err)
			}
			if !tt.expectError {
				if len(result) != len(tt.expectedValues) {
					t.Errorf("%s: Expected %d values, got %d. Expected: %v, Got: %v",
						tt.description, len(tt.expectedValues), len(result), tt.expectedValues, result)
				} else {
					for i, expected := range tt.expectedValues {
						if result[i] != expected {
							t.Errorf("%s: At index %d: expected %s, got %s",
								tt.description, i, expected, result[i])
						}
					}
				}
			}
		})
	}
}

// TODO: fix buggy internal node search
func TestInternalNodeSearch_Issues(t *testing.T) {
	tests := []struct {
		name        string
		node        *BpTreeInternalNode
		searchKey   int
		expectError bool
		description string
	}{
		{
			name: "Search in empty internal node - panic",
			node: &BpTreeInternalNode{
				Key:      10,
				Children: []*BpTreeLeafNode{},
			},
			searchKey:   10,
			expectError: true,
			description: "TODO BUGGY: search() panics when Children slice is empty due to len(t.Children)-1",
		},
		{
			name: "Search key greater than last child",
			node: &BpTreeInternalNode{
				Key: 10,
				Children: []*BpTreeLeafNode{
					{Key: 10, Value: "value10"},
					{Key: 20, Value: "value20"},
				},
			},
			searchKey:   30,
			expectError: true,
			description: "Should error when key is greater than last child",
		},
		{
			name: "Search key less than first child",
			node: &BpTreeInternalNode{
				Key: 10,
				Children: []*BpTreeLeafNode{
					{Key: 20, Value: "value20"},
					{Key: 30, Value: "value30"},
				},
			},
			searchKey:   15,
			expectError: true,
			description: "TODO BUGGY: search() doesn't properly handle key less than first child",
		},
		{
			name: "Search non-existent key between existing keys",
			node: &BpTreeInternalNode{
				Key: 10,
				Children: []*BpTreeLeafNode{
					{Key: 10, Value: "value10"},
					{Key: 30, Value: "value30"},
				},
			},
			searchKey:   20,
			expectError: true,
			description: "Should error when key doesn't exist in range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectError {
						t.Errorf("%s: Unexpected panic: %v", tt.description, r)
					} else {
						t.Logf("%s: Expected panic occurred: %v", tt.description, r)
					}
				}
			}()

			_, err := tt.node.search(tt.searchKey)

			if tt.expectError && err == nil {
				t.Errorf("%s: Expected error but got none", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: Unexpected error: %v", tt.description, err)
			}
		})
	}
}

// ========== UPDATE METHOD TESTS ==========

// TestUpdate_HappyPath tests successful update scenarios
func TestUpdate_HappyPath(t *testing.T) {
	tests := []struct {
		name             string
		setupTree        func() *BpTreeRootNode
		key              int
		newValue         string
		expectError      bool
		verifyOldValue   string
		verifyFinalValue string
	}{
		{
			name: "Update single key in simple tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "original")
				return tree
			},
			key:              10,
			newValue:         "updated",
			expectError:      false,
			verifyOldValue:   "original",
			verifyFinalValue: "updated",
		},
		{
			name: "Update key from tree with multiple elements",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:              20,
			newValue:         "updated20",
			expectError:      false,
			verifyOldValue:   "value20",
			verifyFinalValue: "updated20",
		},
		{
			name: "Update first key in tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:              10,
			newValue:         "updatedFirst",
			expectError:      false,
			verifyOldValue:   "value10",
			verifyFinalValue: "updatedFirst",
		},
		{
			name: "Update last key in tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:              30,
			newValue:         "updatedLast",
			expectError:      false,
			verifyOldValue:   "value30",
			verifyFinalValue: "updatedLast",
		},
		{
			name: "Update middle key in tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				tree.Insert(50, "value50")
				return tree
			},
			key:              30,
			newValue:         "updatedMiddle",
			expectError:      false,
			verifyOldValue:   "value30",
			verifyFinalValue: "updatedMiddle",
		},
		{
			name: "Update key in tree with multiple internal nodes",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				// Create enough insertions to trigger multiple internal nodes
				for i := 0; i < 10; i++ {
					tree.Insert(i*10, fmt.Sprintf("value%d", i*10))
				}
				return tree
			},
			key:              50,
			newValue:         "updatedValue50",
			expectError:      false,
			verifyOldValue:   "value50",
			verifyFinalValue: "updatedValue50",
		},
		{
			name: "Update key after internal node split",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				tree.Insert(25, "value25") // Triggers split
				return tree
			},
			key:              25,
			newValue:         "updatedAfterSplit",
			expectError:      false,
			verifyOldValue:   "value25",
			verifyFinalValue: "updatedAfterSplit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()

			// Verify old value
			oldValue, err := tree.Get(tt.key)
			if err != nil {
				t.Fatalf("Failed to get old value: %v", err)
			}
			if oldValue != tt.verifyOldValue {
				t.Errorf("Expected old value %s, got %s", tt.verifyOldValue, oldValue)
			}

			// Perform update
			err = tree.Update(tt.key, tt.newValue)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify new value
			if !tt.expectError {
				newValue, err := tree.Get(tt.key)
				if err != nil {
					t.Errorf("Failed to get updated value: %v", err)
				}
				if newValue != tt.verifyFinalValue {
					t.Errorf("Expected updated value %s, got %s", tt.verifyFinalValue, newValue)
				}
			}
		})
	}
}

// TestUpdate_BasicUseCases tests common use case scenarios
func TestUpdate_BasicUseCases(t *testing.T) {
	tests := []struct {
		name        string
		setupTree   func() *BpTreeRootNode
		key         int
		newValue    string
		expectError bool
		description string
	}{
		{
			name: "Update value to same value (idempotent)",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				return tree
			},
			key:         10,
			newValue:    "value10",
			expectError: false,
			description: "Updating to same value should work without issue",
		},
		{
			name: "Update value to empty string",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				return tree
			},
			key:         10,
			newValue:    "",
			expectError: false,
			description: "Should allow updating to empty string",
		},
		{
			name: "Update value to very long string",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "short")
				return tree
			},
			key:         10,
			newValue:    "This is a very long string with lots of characters to test if the update function can handle larger values without issues",
			expectError: false,
			description: "Should handle long string values",
		},
		{
			name: "Update same key multiple times sequentially",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "original")
				tree.Update(10, "update1")
				tree.Update(10, "update2")
				return tree
			},
			key:         10,
			newValue:    "update3",
			expectError: false,
			description: "Multiple sequential updates should work",
		},
		{
			name: "Update then retrieve to verify change persists",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:         20,
			newValue:    "persistentValue",
			expectError: false,
			description: "Updated value should persist after retrieval",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()
			err := tree.Update(tt.key, tt.newValue)

			if tt.expectError && err == nil {
				t.Errorf("%s: Expected error but got none", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: Unexpected error: %v", tt.description, err)
			}

			// Verify the value was actually updated
			if !tt.expectError {
				value, getErr := tree.Get(tt.key)
				if getErr != nil {
					t.Errorf("%s: Failed to retrieve updated value: %v", tt.description, getErr)
				}
				if value != tt.newValue {
					t.Errorf("%s: Expected value %s, got %s", tt.description, tt.newValue, value)
				}
			}
		})
	}
}

// TestUpdate_EdgeCases tests edge case scenarios
func TestUpdate_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		setupTree   func() *BpTreeRootNode
		key         int
		newValue    string
		expectError bool
		description string
	}{
		{
			name: "Update from empty tree",
			setupTree: func() *BpTreeRootNode {
				return NewBpTree()
			},
			key:         10,
			newValue:    "value10",
			expectError: true,
			description: "Should error when updating key in empty tree",
		},
		{
			name: "Update non-existent key from populated tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:         25,
			newValue:    "value25",
			expectError: true,
			description: "Should error when key doesn't exist",
		},
		{
			name: "Update key smaller than all keys in tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:         5,
			newValue:    "value5",
			expectError: true,
			description: "Should error when key is smaller than all keys",
		},
		{
			name: "Update key larger than all keys in tree",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				return tree
			},
			key:         40,
			newValue:    "value40",
			expectError: true,
			description: "Should error when key is larger than all keys",
		},
		{
			name: "Update with math.MinInt key",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(math.MinInt, "minValue")
				tree.Insert(0, "value0")
				tree.Insert(100, "value100")
				return tree
			},
			key:         math.MinInt,
			newValue:    "updatedMin",
			expectError: false,
			description: "Should update minimum integer key",
		},
		{
			name: "Update with math.MaxInt key",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(0, "value0")
				tree.Insert(100, "value100")
				tree.Insert(math.MaxInt, "maxValue")
				return tree
			},
			key:         math.MaxInt,
			newValue:    "updatedMax",
			expectError: false,
			description: "Should update maximum integer key",
		},
		{
			name: "Update with zero key",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(-10, "valueMinus10")
				tree.Insert(0, "value0")
				tree.Insert(10, "value10")
				return tree
			},
			key:         0,
			newValue:    "updatedZero",
			expectError: false,
			description: "Should update zero key",
		},
		{
			name: "Update with negative key",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(-50, "valueMinus50")
				tree.Insert(-10, "valueMinus10")
				tree.Insert(10, "value10")
				return tree
			},
			key:         -10,
			newValue:    "updatedMinus10",
			expectError: false,
			description: "Should update negative key",
		},
		{
			name: "Update after GetRange operation",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				// Perform a range query first
				tree.GetRange(10, 30)
				return tree
			},
			key:         20,
			newValue:    "updatedAfterRange",
			expectError: false,
			description: "Update should work after GetRange operation",
		},
		{
			name: "Update key that exists at internal node boundary",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20")
				tree.Insert(30, "value30")
				tree.Insert(40, "value40")
				tree.Insert(25, "value25") // May trigger split, 30 could be at boundary
				return tree
			},
			key:         30,
			newValue:    "updatedBoundary",
			expectError: false,
			description: "Should update key at internal node boundary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()
			err := tree.Update(tt.key, tt.newValue)

			if tt.expectError && err == nil {
				t.Errorf("%s: Expected error but got none", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: Unexpected error: %v", tt.description, err)
			}

			// Verify the value was actually updated (for non-error cases)
			if !tt.expectError {
				value, getErr := tree.Get(tt.key)
				if getErr != nil {
					t.Errorf("%s: Failed to retrieve updated value: %v", tt.description, getErr)
				}
				if value != tt.newValue {
					t.Errorf("%s: Expected value %s, got %s", tt.description, tt.newValue, value)
				}
			}
		})
	}
}

// TestUpdate_DuplicateKeys tests the specific bug where duplicate keys only update first occurrence
func TestUpdate_DuplicateKeys(t *testing.T) {
	tests := []struct {
		name                  string
		setupTree             func() *BpTreeRootNode
		key                   int
		newValue              string
		expectError           bool
		verifyFirstOccurrence string
		verifyOtherValues     []string
		description           string
	}{
		{
			name: "Update when duplicate keys exist - only first is updated",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10_first")
				tree.Insert(10, "value10_second")
				tree.Insert(20, "value20")
				return tree
			},
			key:                   10,
			newValue:              "updatedValue10",
			expectError:           false,
			verifyFirstOccurrence: "updatedValue10",
			verifyOtherValues:     []string{"value10_second"},
			description:           "BUG: Update only modifies first occurrence of duplicate keys",
		},
		{
			name: "Update with multiple duplicate keys",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10_first")
				tree.Insert(10, "value10_second")
				tree.Insert(20, "value20")
				return tree
			},
			key:                   10,
			newValue:              "updatedValue10",
			expectError:           false,
			verifyFirstOccurrence: "updatedValue10",
			verifyOtherValues:     []string{"value10_second"},
			description:           "BUG: With multiple duplicates, only first is updated",
		},
		{
			name: "Update duplicate key at tree boundary",
			setupTree: func() *BpTreeRootNode {
				tree := NewBpTree()
				tree.Insert(10, "value10")
				tree.Insert(20, "value20_first")
				tree.Insert(20, "value20_second")
				tree.Insert(30, "value30")
				return tree
			},
			key:                   20,
			newValue:              "updatedValue20",
			expectError:           false,
			verifyFirstOccurrence: "updatedValue20",
			verifyOtherValues:     []string{"value20_second"},
			description:           "BUG: Duplicate keys at non-first position only update first occurrence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := tt.setupTree()

			tree.PrettyPrint()

			// Perform update
			err := tree.Update(tt.key, tt.newValue)
			tree.PrettyPrint()

			if tt.expectError && err == nil {
				t.Errorf("%s: Expected error but got none", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: Unexpected error: %v", tt.description, err)
			}

			if !tt.expectError {
				// Get returns first occurrence
				firstValue, err := tree.Get(tt.key)
				if err != nil {
					t.Errorf("%s: Failed to get first occurrence: %v", tt.description, err)
				}
				if firstValue != tt.verifyFirstOccurrence {
					t.Errorf("%s: Expected first occurrence to be %s, got %s",
						tt.description, tt.verifyFirstOccurrence, firstValue)
				}

				// Verify other occurrences remain unchanged by checking all values with this key
				// We'll iterate through the tree to find all occurrences
				var allValuesForKey []string
				for _, inode := range tree.Children {
					for _, leaf := range inode.Children {
						if leaf.Key == tt.key {
							allValuesForKey = append(allValuesForKey, leaf.Value)
						}
					}
				}

				// First should be updated
				if len(allValuesForKey) > 0 && allValuesForKey[0] != tt.verifyFirstOccurrence {
					t.Errorf("%s: First occurrence not updated correctly", tt.description)
				}

				// Check that other values remain unchanged
				if len(allValuesForKey) > 1 {
					for i, expectedOther := range tt.verifyOtherValues {
						if i+1 < len(allValuesForKey) && allValuesForKey[i+1] != expectedOther {
							t.Errorf("%s: Expected other occurrence at index %d to be %s, got %s",
								tt.description, i, expectedOther, allValuesForKey[i+1])
						}
					}
				}

				t.Logf("%s: CONFIRMED - Only first occurrence updated. All values for key %d: %v",
					tt.description, tt.key, allValuesForKey)
			}
		})
	}
}
