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
		inode.addLeafNode(key, val)
		if !inode.isfull() {
			return
		}

		lastEle := inode.Children[len(inode.Children)-1]
		inode.Children = inode.Children[:MAX_SIZE+1]

		if inodeidx+1 <= MAX_SIZE {
			if len(t.Children) >= inodeidx+1 {
				// how should this traversal happen ?
				/**
				The issue is for max size of 4 we. have 4 internal nodesa

				1.        6.         10.          15 <-internal nodes

				1 3 4 5     6 7 8 9    10 11 12 13     15 16 <- leaf nodes

				Mmy tree still has some space but. if I want to insert key 2
				it will cause cascade effect in all the other internal nodes how should. I handle this?

				Answer:

				When internal node size > MAX_SIZE
				- We split internal node
				- We update parent node with the the splitted nodes
				- In case parent also reaches max size we split parent too
					- For single level internal nodes we will return error
					- For multi level internal nodes we will split parent too
						- This will also require parent reference in child node
				*/
			}
		}

	}
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
	inode := &BpTreeInternalNode{Key: key}
	return &inode
}

func (t *BpTreeInternalNode) isfull() bool {
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

// func (bpn *BpTreeNode) insertChild(n int) {
// 	if len(bpn.Children) == 0 {
// 		bpn.Children = append(bpn.Children, &BpTreeNode{Key: n})
// 		return
// 	}

// 	var i int
// 	for index, v := range bpn.Children {
// 		if v.Key > n {
// 			i = index - 1
// 			break
// 		}
// 	}

// 	bpn.Children = append(bpn.Children, &BpTreeNode{})
// 	copy(bpn.Children[i+1:], bpn.Children[i:])
// 	bpn.Children[i] = &BpTreeNode{Key: n}
// }

// // returns index of to be parent node
// func findLeafNodeParent(refs []*BpTreeNode, val int) int {
// 	// only current level for now
// 	// if 1 ref => then 0
// 	// if 2 ref => if less than second => 1st other wise second
// 	// if 3 refs

// 	cur := 0

// 	nodes := refs
// 	for len(nodes) > 2 {
// 		midpoint := len(nodes) / 2

// 		if val >= nodes[midpoint].Key {
// 			cur += midpoint
// 			nodes = nodes[midpoint:]
// 		} else {
// 			nodes = nodes[:midpoint]
// 		}
// 	}

// 	if len(nodes) == 0 {
// 		return -1 // Handle empty case, though probably not reached
// 	}

// 	if len(nodes) == 1 {
// 		if nodes[0].Key > val {
// 			return -1
// 		}
// 		return cur + 0
// 	}

// 	if len(nodes) == 2 {
// 		if nodes[0].Key > val {
// 			return -1
// 		}

// 		if nodes[0].Key < val && val < nodes[1].Key {
// 			return cur + 0
// 		}

// 		return cur + 1
// 	}

// 	return -1 // Fallback, should not reach
// }

// func (bp *BpTreeNode) Insert(n int) {
// 	// If it's the first element to be inserted then
// 	if len(bp.Children) == 0 {
// 		bp.Children = append(bp.Children, &BpTreeNode{Key: n})
// 		return
// 	}
// 	// find leaf node parent
// 	index := findLeafNodeParent(bp.Children, n)

// 	if index == -1 {
// 		// If index is -1 , I have to insert a new intermediate node
// 		// I also have to keep in mind that I can't go over the max size of intermediate nodes

// 		if len(bp.Children)+1 < MAX_SIZE {
// 			// add in the intermediate node
// 			// add the element here
// 			n := len(bp.Children)
// 			bp.Children = bp.Children[:n+1]
// 			copy(bp.Children[1:], bp.Children[:n])
// 			bp.Children[0] = &BpTreeNode{Key: n}
// 			bp.Children[0].Children = append(bp.Children[0].Children, &BpTreeNode{})
// 		} else {
// 			// create new intermediate node
// 			// reference it to parent
// 			// add element in leaf node

// 			// lets keep it like this for now
// 		}
// 	} else {
// 		// TODO: currently not checking leaf node
// 		bp.Children[index].insertChild(n) // insert child will take care of inserting the child in correct sorted order
// 	}
// }
