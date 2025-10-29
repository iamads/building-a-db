package main

import (
	"fmt"
)

const MAX_SIZE = 4

// Root node will always point to intermediate nodes
type BpTreeRootNode struct {
	Children []*BpTreeInternalNode
}

func (t *BpTreeRootNode) isFull() bool {
	return len(t.Children) == MAX_SIZE
}

// BpTreeInternal Node will always point to Leaf node (For now we are talking about single level)

// Key will be +ve integer for now

/*
*

	Internal node references
	lets internal node has children: [ {key: a}, {key: b}, {key: c}]
	First internal node will reference => [a, b) children
	2nd internal node will reference => [b, c) children
	2rd internal node will reference => [c, math.MaxInt) children
*/
type BpTreeInternalNode struct {
	Key      int
	Children []*BpTreeLeafNode
}

// BptreeLeaf node will have the key value pair
type BpTreeLeafNode struct {
	Key   int
	Value string
}

// New Bptree root node will create a new BpTreeRootNode
func NewBpTree() *BpTreeRootNode {
	tree := BpTreeRootNode{}
	return &tree
}

/*
Insertion flow (single level)

- Inspect root nodes children
- scan for element which has value just less than on equal to node to be inserted (TODO: optimise to binary search )
- Once internal node is found
	- We check if insertion causes us to hit our MAX_SIZE
		- If it does not we plainly insert
		- If it does
			- We check if we can create a new internal node
				- If yes we create a new internal node and add leaf node
				- If not we return an error that no more elements can be inserted.

**/

/*
 * @topic: Handling Node Splitting and the "Cascade" Effect
 *
 * This note explains the logic for what happens when a node (leaf or internal)
 * overflows and how that split can propagate up the tree.
 *
 * --- The Splitting Process ---
 *
 * 1.  **Node is Full:** An insertion causes a node's key count to exceed MAX_SIZE.
 *
 * 2.  **Split:** The full node is split into two new nodes (let's call them L and R).
 *
 * 3.  **Promote Key:** A key is sent "up" to the parent node to act as a separator.
 * -   **If Leaf Split:** We *copy* the first key of the new R node to the parent. (Currently we do not have multiple values in leaf node, we can consider this at a later stage TODO)
 * -   **If Internal Split:** We *move* the middle key from the full node up to the parent.
 *
 * 4.  **Insert into Parent (The "Cascade"):**
 * -   The parent node now has to insert the promoted key and a pointer to the new R node.
 * -   **If the parent is NOT full:** The key/pointer are added, and the process stops.
 * (e.g., inserting key '2' in the example tree only splits a leaf and
 * adds key '4' to the root, which has space. No cascade occurs.)
 *
 * -   **If the parent IS full:** The parent is now *also* over MAX_SIZE.
 * We must repeat this entire process (from Step 2) on the parent.
 * This "split-and-promote" can continue all the way up.
 *
 * 5.  **Root Split (Tree Growth):**
 * -   If the cascade reaches the **Root** and the *Root itself splits*,
 * we create a **NEW, empty root** node.
 * -   This new root will contain the single key promoted from the old root.
 * -   The new root's two children will be the L and R nodes created
 * from splitting the old root.
 * -   This is the *only* way the B+ Tree increases in height.
 *
 * --- Implementation Notes ---
 *
 * 1.  **Parent Pointer:** For the cascade (Step 4) to work, every child node
 * *must* have a pointer/reference to its parent node.
 *
 */

func (t *BpTreeRootNode) Insert(key int, val string) error {
	inodeidx := t.findInternalPredecessor(key)

	if inodeidx == -1 {
		// TDOD: create new internal node
		if t.isFull() {
			return fmt.Errorf("Root node is full, we can't insert key=(%d) and val=(%s)", key, val)
		}
		t.Children = append(t.Children, &BpTreeInternalNode{})
		copy(t.Children[1:], t.Children[:len(t.Children)-1])
		newInode := createNewInternalNode(key)
		newInode.addLeafNode(key, val)
		t.Children[0] = newInode
	} else {
		inode := t.Children[inodeidx]
		if inode.isFull() {
			// If root has space we will split current node
			if len(t.Children) == MAX_SIZE {
				original, next := splitInode2(inode)
				t.Children[inodeidx] = original
				t.Children = append(t.Children, &BpTreeInternalNode{})
				copy(t.Children[inodeidx+2:], t.Children[inodeidx+1:])
				t.Children[inodeidx+1] = next
			} else {
				// else we return error
				return fmt.Errorf("Can't create new internal node for key=(%d) and val=(%s)", key, val)
			}
		} else {
			inode.addLeafNode(key, val)
		}
	}

	return nil
}

func splitInode2(inode *BpTreeInternalNode) (*BpTreeInternalNode, *BpTreeInternalNode) {
	q := MAX_SIZE / 2

	originalChildren := inode.Children[:q]
	nextChildren := inode.Children[q:]

	original := createNewInternalNode(originalChildren[0].Key)
	original.Children = originalChildren

	next := createNewInternalNode(nextChildren[0].Key)
	next.Children = nextChildren

	return original, next

}

// If all elemennts are greater than val -> -1
// If all elements are smaller than val -> len-1
// Otherwise find theh  predecessor i.e key which is just smaller than the current val

func (t *BpTreeRootNode) findInternalPredecessor(key int) int {
	if len(t.Children) == 0 {
		return -1
	}

	if t.Children[0].Key > key {
		return -1
	}

	if t.Children[len(t.Children)-1].Key <= key {
		return len(t.Children) - 1
	}

	index := 0
	for i, v := range t.Children {
		if v.Key > key {
			index = i - 1
			break
		}
	}
	return index
}

func createNewInternalNode(key int) *BpTreeInternalNode {
	inode := BpTreeInternalNode{Key: key}
	return &inode
}

func (t *BpTreeInternalNode) isFull() bool {
	return len(t.Children) == MAX_SIZE
}

func (t *BpTreeInternalNode) addLeafNode(key int, val string) {
	toInsertIdx := 0

	if key == t.Key {
		toInsertIdx = 0
	}

	if len(t.Children) > 0 && key > t.Children[len(t.Children)-1].Key {
		toInsertIdx = len(t.Children)
	}

	for i, v := range t.Children {
		if v.Key > key {
			toInsertIdx = i
			break
		}
	}

	t.Children = append(t.Children, &BpTreeLeafNode{})
	copy(t.Children[toInsertIdx+1:], t.Children[toInsertIdx:])
	t.Children[toInsertIdx] = &BpTreeLeafNode{Key: key, Value: val}
}
