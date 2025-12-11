package temp

/***
RESULT

Missing key bugs is happening in the code shared by author too
So I will not spend a lot of time right now to look into the issue
But the issue with node being invalid is not happening ...so we need to dig a little
deeper into it
*/

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"unsafe"

	tassert "github.com/stretchr/testify/assert"
)

// C is a test context that maintains both the BTree and a reference map
type C struct {
	tree  BTree
	ref   map[string]string
	pages map[uint64]BNode
}

// newC creates a new test context with an initialized BTree
func newC() *C {
	pages := map[uint64]BNode{}
	return &C{
		tree: BTree{
			get: func(ptr uint64) []byte {
				node, ok := pages[ptr]
				assertStatement(ok, "got page!")
				return node
			},
			new: func(node []byte) uint64 {
				assertStatement(BNode(node).nbytes() <= BTREE_PAGE_SIZE, "Node should be smaller than max node size")
				ptr := uint64(uintptr(unsafe.Pointer(&node[0])))
				assertStatement(pages[ptr] == nil, "page ptr should not point to nil")
				pages[ptr] = node
				return ptr
			},
			del: func(ptr uint64) {
				assertStatement(pages[ptr] != nil, "page ptr should not point to nil")
				delete(pages, ptr)
			},
		},
		ref:   map[string]string{},
		pages: pages,
	}
}

// assertStatement is a helper for test assertions
func assertStatement(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

// Helper: Add key-value to both tree and ref map
func (c *C) add(key, val string) error {
	err := c.tree.Insert([]byte(key), []byte(val))
	if err == nil {
		c.ref[key] = val
	}
	return err
}

// Helper: Delete key from both tree and ref map
func (c *C) del(key string) (bool, error) {
	deleted, err := c.tree.Delete([]byte(key))
	if err == nil && deleted {
		delete(c.ref, key)
	}
	return deleted, err
}

// Helper: Verify all nodes in tree satisfy size constraint
func (c *C) verifyNodeSizes(t *testing.T) {
	for ptr, node := range c.pages {
		tassert.True(t, node.nbytes() <= BTREE_PAGE_SIZE,
			"Node at ptr %d has size %d, exceeds max %d", ptr, node.nbytes(), BTREE_PAGE_SIZE)
	}
}

// Helper: Verify keys are sorted within a node
func verifyNodeKeysSorted(t *testing.T, node BNode) {
	nkeys := node.nkeys()
	for i := uint16(1); i < nkeys; i++ {
		prevKey := node.getKey(i - 1)
		currKey := node.getKey(i)
		cmp := bytes.Compare(prevKey, currKey)
		tassert.True(t, cmp <= 0,
			"Keys not sorted: key[%d]=%q should be <= key[%d]=%q", i-1, prevKey, i, currKey)
	}
}

// Helper: Recursively verify all keys are sorted in the tree
func (c *C) verifyKeysSorted(t *testing.T) {
	if c.tree.root == 0 {
		return
	}
	c.verifyNodeKeysSortedRecursive(t, c.tree.root)
}

func (c *C) verifyNodeKeysSortedRecursive(t *testing.T, ptr uint64) {
	node := BNode(c.tree.get(ptr))
	verifyNodeKeysSorted(t, node)

	if node.btype() == BNODE_NODE {
		// Recursively check children
		for i := uint16(0); i < node.nkeys(); i++ {
			childPtr := node.getPtr(i)
			c.verifyNodeKeysSortedRecursive(t, childPtr)
		}
	}
}

// Helper: Verify tree data matches ref map
func (c *C) verifyDataIntegrity(t *testing.T) {
	// Collect all keys from tree
	treeKeys := c.collectAllKeys()

	// Verify count matches
	tassert.Equal(t, len(c.ref), len(treeKeys),
		"Tree has %d keys but ref has %d keys", len(treeKeys), len(c.ref))

	// Verify each key in ref exists in tree
	for key := range c.ref {
		found := false
		for _, treeKey := range treeKeys {
			if bytes.Equal([]byte(key), treeKey) {
				found = true
				break
			}
		}
		tassert.True(t, found, "Key %q in ref but not in tree", key)
	}
}

// Helper: Collect all keys from tree via in-order traversal
func (c *C) collectAllKeys() [][]byte {
	if c.tree.root == 0 {
		return [][]byte{}
	}
	var keys [][]byte
	c.collectKeysRecursive(c.tree.root, &keys)
	return keys
}

func (c *C) collectKeysRecursive(ptr uint64, keys *[][]byte) {
	node := BNode(c.tree.get(ptr))

	if node.btype() == BNODE_LEAF {
		// Collect keys from leaf (skip sentinel at index 0)
		for i := uint16(1); i < node.nkeys(); i++ {
			key := node.getKey(i)
			if len(key) > 0 { // Skip empty keys
				*keys = append(*keys, key)
			}
		}
	} else {
		// Recursively collect from internal node children
		for i := uint16(0); i < node.nkeys(); i++ {
			c.collectKeysRecursive(node.getPtr(i), keys)
		}
	}
}

// Helper: Count total keys in tree
func (c *C) countKeys() int {
	return len(c.collectAllKeys())
}

// TestBTreeStressOperations high-volume test
func TestBTreeStressOperations(t *testing.T) {
	t.Run("1000 mixed operations", func(t *testing.T) {
		c := newC()

		// Perform 1000 random operations
		for i := 0; i < 1000; i++ {
			op := rand.Float32()
			key := fmt.Sprintf("key_%d", rand.Intn(500))

			if op < 0.5 { // 50% insert
				val := fmt.Sprintf("value_%d", i)
				c.add(key, val)
			} else if op < 0.8 { // 30% delete
				c.del(key)
			} else { // 20% update
				val := fmt.Sprintf("updated_%d", i)
				c.add(key, val)
			}

			// Periodic verification (every 100 ops)
			if i%100 == 99 {
				c.verifyKeysSorted(t)
				c.verifyNodeSizes(t)
			}
		}

		// Final comprehensive verification
		c.verifyKeysSorted(t)
		c.verifyNodeSizes(t)
		c.verifyDataIntegrity(t)

		t.Logf("Final state: %d keys in tree, %d pages allocated",
			c.countKeys(), len(c.pages))
	})
}

// TestBTreeNodeSizeInvariants continuously verifies node sizes
func TestBTreeNodeSizeInvariants(t *testing.T) {
	t.Run("Node sizes valid throughout operations", func(t *testing.T) {
		c := newC()

		// Perform 100 random operations
		for i := 0; i < 100; i++ {
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Logf("Panic at iteration %d: %v", i, r)
						t.Logf("Tree state: root=%d, pages=%d, ref keys=%d",
							c.tree.root, len(c.pages), len(c.ref))
						t.FailNow()
					}
				}()

				if rand.Float32() < 0.7 { // 70% inserts
					key := fmt.Sprintf("key_%d", rand.Intn(50))
					val := fmt.Sprintf("value_%d", i)
					err := c.add(key, val)
					if err != nil {
						t.Logf("Insert failed at iteration %d: %v", i, err)
					}

				} else { // 30% deletes
					key := fmt.Sprintf("key_%d", rand.Intn(50))
					_, err := c.del(key)
					if err != nil {
						t.Logf("Delete failed at iteration %d: %v", i, err)
					}
				}

				// Verify after EVERY operation
				c.verifyNodeSizes(t)
			}()
		}

		// Final verification
		c.verifyKeysSorted(t)
		// Only verify data integrity if tree has keys
		if c.tree.root != 0 && len(c.ref) > 0 {
			c.verifyDataIntegrity(t)
		}
	})
}

// TestBTreeDataIntegrity verifies tree matches ref map
func TestBTreeDataIntegrity(t *testing.T) {
	t.Run("Tree data matches ref map", func(t *testing.T) {
		c := newC()

		// Insert 200 keys
		for i := 0; i < 200; i++ {
			key := fmt.Sprintf("key_%04d", i)
			val := fmt.Sprintf("value_%d", i)
			c.add(key, val)
		}

		// Verify integrity
		c.verifyDataIntegrity(t)

		// Update 50 keys
		for i := 0; i < 50; i++ {
			key := fmt.Sprintf("key_%04d", i)
			newVal := fmt.Sprintf("updated_%d", i)
			c.add(key, newVal)
		}

		// Verify integrity after updates
		c.verifyDataIntegrity(t)

		// Delete 100 keys
		for i := 50; i < 150; i++ {
			key := fmt.Sprintf("key_%04d", i)
			c.del(key)
		}

		// Verify integrity after deletes
		c.verifyDataIntegrity(t)
		tassert.Equal(t, 100, len(c.ref))
	})
}

// TestBTreeKeysSortedInvariant continuously verifies keys sorted
func TestBTreeKeysSortedInvariant(t *testing.T) {
	t.Run("Keys remain sorted throughout operations", func(t *testing.T) {
		c := newC()

		// Insert in random order
		keys := make([]string, 100)
		for i := 0; i < 100; i++ {
			keys[i] = fmt.Sprintf("%03d", rand.Intn(1000))
			c.add(keys[i], "value")
		}

		// Verify sorted
		c.verifyKeysSorted(t)

		// Delete random keys
		for i := 0; i < 50; i++ {
			if i < len(keys) {
				c.del(keys[i])
			}
		}

		// Verify still sorted
		c.verifyKeysSorted(t)
	})
}
