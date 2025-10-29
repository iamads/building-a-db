package main

import (
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
			tree:     &BpTreeRootNode{Children: []*BpTreeInternalNode{&BpTreeInternalNode{Key: 10}, &BpTreeInternalNode{Key: 20}, &BpTreeInternalNode{Key: 100}}},
			val:      5,
			expected: -1,
		},
		{
			name:     "Value after the first. node itself",
			tree:     &BpTreeRootNode{Children: []*BpTreeInternalNode{&BpTreeInternalNode{Key: 10}, &BpTreeInternalNode{Key: 20}, &BpTreeInternalNode{Key: 100}}},
			val:      11,
			expected: 0,
		},
		{
			name:     "Value after the middle node",
			tree:     &BpTreeRootNode{Children: []*BpTreeInternalNode{&BpTreeInternalNode{Key: 10}, &BpTreeInternalNode{Key: 20}, &BpTreeInternalNode{Key: 100}}},
			val:      70,
			expected: 1,
		},

		{
			name:     "Value is maxInt",
			tree:     &BpTreeRootNode{Children: []*BpTreeInternalNode{&BpTreeInternalNode{Key: 10}, &BpTreeInternalNode{Key: 20}, &BpTreeInternalNode{Key: 100}}},
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

// func TestFindLeafNodeParent(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		refs     []*BpTreeNode
// 		val      int
// 		expected int
// 	}{
// 		// Empty refs array
// 		{
// 			name:     "empty refs array",
// 			refs:     []*BpTreeNode{},
// 			val:      5,
// 			expected: -1,
// 		},

// 		// Single ref cases
// 		{
// 			name:     "single ref, value less than node value",
// 			refs:     []*BpTreeNode{{Value: 10}},
// 			val:      5,
// 			expected: -1,
// 		},
// 		{
// 			name:     "single ref, value equal to node value",
// 			refs:     []*BpTreeNode{{Value: 5}},
// 			val:      5,
// 			expected: 0,
// 		},
// 		{
// 			name:     "single ref, value greater than node value",
// 			refs:     []*BpTreeNode{{Value: 3}},
// 			val:      5,
// 			expected: 0,
// 		},

// 		// Two refs cases
// 		{
// 			name:     "two refs, value less than first",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}},
// 			val:      5,
// 			expected: -1,
// 		},
// 		{
// 			name:     "two refs, value between first and second",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}},
// 			val:      15,
// 			expected: 0,
// 		},
// 		{
// 			name:     "two refs, value equal to first",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}},
// 			val:      10,
// 			expected: 1,
// 		},
// 		{
// 			name:     "two refs, value equal to second",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}},
// 			val:      20,
// 			expected: 1,
// 		},
// 		{
// 			name:     "two refs, value greater than second",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}},
// 			val:      25,
// 			expected: 1,
// 		},

// 		// Three refs cases
// 		{
// 			name:     "three refs, value less than first",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}, {Value: 30}},
// 			val:      5,
// 			expected: -1,
// 		},
// 		{
// 			name:     "three refs, value between first and second",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}, {Value: 30}},
// 			val:      15,
// 			expected: 0,
// 		},
// 		{
// 			name:     "three refs, value between second and third",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}, {Value: 30}},
// 			val:      25,
// 			expected: 1,
// 		},
// 		{
// 			name:     "three refs, value greater than third",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}, {Value: 30}},
// 			val:      35,
// 			expected: 2,
// 		},

// 		// More complex cases with binary search behavior
// 		{
// 			name:     "four refs, binary search left side",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}, {Value: 30}, {Value: 40}},
// 			val:      15,
// 			expected: 0,
// 		},
// 		{
// 			name:     "four refs, binary search right side",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}, {Value: 30}, {Value: 40}},
// 			val:      35,
// 			expected: 2,
// 		},
// 		{
// 			name:     "five refs, middle value",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}, {Value: 30}, {Value: 40}, {Value: 50}},
// 			val:      30,
// 			expected: 2,
// 		},

// 		// Edge cases with duplicate values (though B+ tree typically doesn't allow duplicates)
// 		// ******** This is invalid as intermediate node will not have duplicates
// 		// {
// 		// 	name:     "duplicate values, find first occurrence",
// 		// 	refs:     []*BpTreeNode{{Value: 10}, {Value: 10}, {Value: 20}},
// 		// 	val:      10,
// 		// 	expected: 0,
// 		// },

// 		// Boundary values
// 		{
// 			name:     "minimum integer value",
// 			refs:     []*BpTreeNode{{Value: 0}, {Value: 10}},
// 			val:      -2147483648,
// 			expected: -1,
// 		},
// 		{
// 			name:     "maximum integer value",
// 			refs:     []*BpTreeNode{{Value: 10}, {Value: 20}},
// 			val:      2147483647,
// 			expected: 1,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result := findLeafNodeParent(tt.refs, tt.val)
// 			if result != tt.expected {
// 				t.Errorf("findLeafNodeParent(%v, %d) = %d; expected %d: Test case: %s",
// 					tt.refs, tt.val, result, tt.expected, tt.name)
// 			}
// 		})
// 	}
// }

// func TestFindLeafNodeParentIntegration(t *testing.T) {
// 	// Test that simulates insertion scenario
// 	// Create a tree structure and test finding parent for insertion

// 	// Create intermediate node with children
// 	intermediate := &BpTreeNode{
// 		Value: 25,
// 		Children: []*BpTreeNode{
// 			{Value: 10, Children: []*BpTreeNode{{Value: 5}, {Value: 15}}},
// 			{Value: 40, Children: []*BpTreeNode{{Value: 30}, {Value: 50}}},
// 		},
// 	}

// 	// Test finding parent for values that should go to left child
// 	leftParent := findLeafNodeParent(intermediate.Children, 12)
// 	if leftParent != 0 {
// 		t.Errorf("Expected 0 for value 12, got %d", leftParent)
// 	}

// 	// Test finding parent for values that should go to right child
// 	rightParent := findLeafNodeParent(intermediate.Children, 35)
// 	// Let's see what the function actually returns
// 	t.Logf("For value 35 with refs [10, 40], function returned: %d", rightParent)
// 	if rightParent != 0 {
// 		t.Errorf("Expected 1 for value 35, got %d", rightParent)
// 	}

// 	// Test value less than all children
// 	lowParent := findLeafNodeParent(intermediate.Children, 1)
// 	if lowParent != -1 {
// 		t.Errorf("Expected -1 for value 1, got %d", lowParent)
// 	}

// 	// Test value greater than all children
// 	highParent := findLeafNodeParent(intermediate.Children, 60)
// 	if highParent != 1 {
// 		t.Errorf("Expected 1 for value 60, got %d", highParent)
// 	}
// }

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
