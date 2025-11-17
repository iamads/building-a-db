package db

import "encoding/binary"

// Backed by disk

/*
*
Node structure:
TODO: For now we will have one node for internal and leaf nodes but later it will be changed

NodeType 2 bytes
number of keys(keys) 2 bytes

poiters: nkeys * 8 bytes
offsets:  n keys * 2 bytes
key-values: (we will store key value pair)
unused space

Format of each KV pair
klen: 2 byte
vlen: 2byte
key: .....
value.......

Node size will be 4kb, which is typical os page size
*/
const HEADER = 4

const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	// assert(node1max <= BTREE_PAGE_SIZE)
}

// If we use use bnode as byte we will skip serialisation desrialisation cost
type BNode []byte

type Btree struct {
	root uint64

	get func(uint64) []byte
	new func([]byte) uint64
	del func(uint64)
}

// Little endian and Big Endian is way how we are storing data
// Little endian => Least significant byte
// Big Endian > Most significant byte
// TODO: Create note

// Decoding Header

const (
	BNODE_NODE = 1 // internal nodes
	BNODE_LEAF = 2 // leaf ndoes
)

func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

func (node BNode) getPtr(idx uint16) uint64 {
	assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) setPtr(idx uint16, val uint64) {
	// TODO
	// if idx < nkeys, I am probably doing an update
	// if idx == nkeys, I am probably doing an insert and should increase nkeys after this
	// What if idx > nkeys ??
}

// OFFSET related

// Each offset is the end of the KV pair relative to start of the 1st KV. The start offset of the 1st KV
// is just 0, so we use the end offset instead, which is the start offset of the next KV
// TODO:The above is from the book and I don't totally agree, we might change later

func offsetPos(node BNode, idx uint16) uint16 {
	assert(1 <= idx && idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}

func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16((node[offsetPos(node, idx):]))
}

func (node BNode) setOffset(idx uint16, offset uint16) {
	// find offset position
	// the at offset position: store the offset
	// TODO: unsure how that storage action will work
}

// KVPOS

func (node BNode) kvPos(idx uint16) uint16 {
	assert(idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:]) // TODO understand this
	return node[pos+4:][:klen]
}

func (node BNode) geVal(idx uint16) []byte {
	// TODO
	// find pos
	// find key len
	// find value length and return value
}

func leafInsert(new BNode, old BNode, idx uint16, key, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys()+1)

	// TODO: Understand flow
	nodeAppendRange(new,  old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, old.nkeys()-idx)
}


func nodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, val []byte) {
	new.setPtr(idx, ptr)

	pos := new.kvPos(idx)
	binary.LittleEndian.PutUint16(new[pos+0:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new[pos+2:], uint16(len(val)))

	copy(new[pos+4:], key)
	copy(new[pos+4+uint16(len(key)):], val)


	// set Offset is being used to set offset of the next key
	new.setOffset(idx +1, new.getOffset(idx)+4+uint16(len(key)+len(val)))
}

// copy multiple kvs
func nodeAppendRange(
	new BNode, old BNode, dstNew uint16, srcOld uint16, n uint16
) {
	// TODO

}


func nodeReplaceKidN(
	tree *Btree, new, old BNode, idx uint16
	kids ...BNode
) {
	inc := uint16(len(kids))
	new.setHeader(BNODE_NODE, old.nkeys() + inc -1)

	nodeAppendRange(new, old, 0, 0, idx)
	for i, node := range kids {
		nodeAppendKV(new, idx+uint16(i), tree.new(node), node.getKey(0), nil)
		// 								^position				^pointer				^key						^val
		nodeAppendRange(new, old, idx+inc, idx+1, old.nkeys() - (idx+1))
	}
}

// split the old ndoe into two nodes -> left right
func nodeSplit2(left, right, old BNode) {
	// code omitted
}

func nodeSplit3(old BNode) (uint16, [3]BNode) {
	if old.nbytes <= BTREE_PAGE_SIZE {
		old = old[:BTREE_PAGE_SIZE]
		return 1, [3]BNode{old} // no split
	}

	left := BNode(make([]byte, 2*BTREE_PAGE_SIZE))
	right := BNode(make([]byte, BTREE_PAGE_SIZE))
	nodeSplit2(left, right, old)

	if left.nbytes <= BTREE_PAGE_SIZE { // TODO: this indicates when we do split2, right is always less than PAGE SIZE and whatever remaining is put into left
		left := left[:BTREE_PAGE_SIZE]
		return 2. [3]BNode{left, right}
	}

	leftleft := BNode(make([]byte, BTREE_PAGE_SIZE))
	middle := BNODE(make([]byte, BTREE_PAGE_SIZE))

	nodeSplit2(leftleft, middle, left)

	assert(leftleft.nbytes() <= BTREE_PAGE_SIZE) // TODO: What happens if leftleft is not less than BTREE_PAGE_SIZE
	return 3, [3]BNode(leftleft, middle, right)
}


// insert a KV into a node, the result might be split
// the called is responsible, for deallocation the input node
// and splitting and allocating the result nodes
func leafInsert(tree *BTree, node BNode, key []byte, val []byte) BNode {
	// the result node
	// it's allowed to be bigger than 1 page and will be split if so

	new := BNode(make([]byte, 2 * BTREE_PAGE_SIZE))

	// where to insert key
	idx := nodeLookupLE(node, key)

	// act depending on node type

	switch node.btype() {
	case BNODE_LEAF:
		// leaf, node.getKey(idx) <= idx
		if bytes.Equal(key, node.getKey(idx)) {
			leafUpdate(new, node, idx, key, val)
		} else {
			// insert it after the position
			leafInsert(new, node, idx+1, key, val)
		}
	case BNODE_NODE:
		nodeInsert(tree, new, node, idx, key, val)
	default:
		panic("bad node!")
	}

	return new
}

// part of treeInsert()
func nodeInsert(
	tree *Btree, new, node BNode, idx uint16,
	key, val []byte
) {
	kpte := node.getPtr(idx)
	// recursive insertion to kid node
	knode := treeInsert(tree, tree.get(kptr), key, val)
	// split the result
	nsplit, split = nodeSplit3(knode)
	// deallocate the kid node
	tree.del(kptr)
	//update the kid links
	nodeReplaceKidN(tree, new, node, idx, split[:nplit])
}


func leafUpdate(new, node BNode, idx uint16, key, val []byte) {

}

func leafInsert(new, node BNode, idx uint16, key, val []byte) {
	
}

// insert a new key or update an existing key
func (tree *BTree) Insert(key []byte, val []byte)

// delete a key and returns whenther the key was there

func (tree *Btree) Delete(key []byte) bool



func (tree *BTree) Insert(key []byte, val []byte) {
		if tree.root == 0 {
			// create the first node
			root := BNode(make([]byte, BTREE_PAGE_SIZE))
			root.setHeader(BNODE_LEAF, 2)
			// TODO: I don't get this
			// a dummy key, this makes the tree cover the whole key space
			// thus a lookup can always find a containing node.

			nodeAppendKV(root, 0, 0, nil, nil)
			nodeAppendKV(root, 1, 0, key, val)
			tree.root = tree.new(root)
			return
		}

		node := treeIndert(tree, tree.get(tree.root), key, val) 
		nplsit, split := nodeSplit3(node)
		tree.del(tree.root)
		if nsplit > 1 {
			root := BNode(make([]byte, BTREE_PAGE_SIZE))
			root.setHeader(BNODE_NODE, nsplit)
			for i, knode := range split[:nsplit] {
				ptr, key := tree.new(knode), knode.getKey(0)
				nodeAppendKV(root, uint16(i), ptr, key, nil)
			}
			tree.root = tree.new(root)
		}else {
			tree.root = tree.new(split[0])
		}
}

// Merginfg  nodes

// remove a key from leaf Node
func leafDelete(new, old BNode, idx uint16)

// merge 2 nodes into 1
func nodeMerge(new, left, right BNode)

// replace 2 adjacent links with 1

func nodeReplace2Kid(new, old BNode, idx, ptr uint16, key []byte)

// should the updated kid me merged with siblig?

// Basically if a node has data < quarter of BTREE page size we want to merge sibling nodes
// return 0 for no merge, -1 for merge with leftand 1 to merge with right
func shouldMerge(tree *BTree, node BNode, idx uint16, updated BNode) (int, BNode) {
	if updated.nbytes() > BTREE_PAGE_SIZE/4 { // TODO: Why?
		return 0, BNode{}
	}

	if idx > 0{
		sibling := BNode(tree.get(node.getPrefix(idx-1)))
		merged := sibling.nbytes()+updated.nbytes() -HEADER
		if merged <= BTREE_PAGE_SIZE {
			return -1, sibling // left
		}
	}

	if idx+1 < node.nkeys() {
		sibling := BNode(tree.get(node.getPtr(idx+1)))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return 1, sibling
		}
	}

	return 0, BNode{}
}

// delete a key from the tree
func treeDelete(tree *Btree, node BNode, key []byte) BNodeleft


// delete a key from an internal node; part of tree Delete()
func nodeDelete(tree *Btree, node BNode, idx uint16, key []byte) BNode {
	// recurse into the kid
	kptr := node.getPtr(idx)
	updated := treeDelete(tree, tree.get(kptr), key)
	if len(updated) == 0 {
		return BNode{}
	}

	tree.Del(kptr)

	new := BNode(make([]byte, BTREE_PAGE_SIZE))
	
	//check for mergr
	mergeDir, sibling := shouldMerge(tree, node, idx, updated)
	switch {
		case mergeDir < 0: // left
			merged := BNode(make([]byte, BTREE_PAGE_SIZE))
			nodeMerge(merged, sibling, updated)
			tree.del(node.getPtr(idx - 1))
			nodeReplace2Kid(new, node, idx-1, tree.new(merged), merged.getKey(0))
		case mergeDir > 0:// right
			merged := BNode(make([]byte, BTREE_PAGE_SIZE))
			nodeMerge(merged, sibling, updated)
			tree.del(node.getPtr(idx - 1))
			nodeReplace2Kid(new, node, idx, tree.new(merged), merged.getKey(0))
		case mergeDir == 0 && updated.nkeys() == 0:
			assert(node.nkeys() == 1 && idx == 0) // 1 empty child but no sibling
			new.setHeader(BNODE_NODE, 0) // the parent becomes empty too
		case mergeDir == 0 && updated.nkeys() > 0: // no merge
			nodeReplaceKidN(tree, new, node, idx, updated)
		}
	return new	
}