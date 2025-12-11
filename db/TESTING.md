# B+Tree Testing Documentation

## Overview

This document describes the comprehensive test suite for the B+tree implementation with mocked file I/O.

## Test Files

- **`db_test.go`** - Unit tests for individual node operations (low-level)
- **`btree_integration_test.go`** - Integration tests for full tree operations with mocked file I/O

## Integration Test Structure

### Test Helper: C Struct

The `C` struct provides mocked file I/O for testing:

```go
type C struct {
    tree  BTree                // B+tree instance
    ref   map[string]string    // Reference data for validation
    pages map[uint64]BNode     // Mocked page storage
}
```

**Key Methods:**
- `add(key, val)` - Insert into both tree and ref map
- `del(key)` - Delete from both tree and ref map
- `verifyKeysSorted()` - Assert all keys sorted in tree
- `verifyNodeSizes()` - Assert all nodes ≤ BTREE_PAGE_SIZE
- `verifyDataIntegrity()` - Assert tree matches ref map

## Test Categories

### 1. Insert Integration Tests

**File:** `TestBTreeInsertIntegration`

**Purpose:** Verify insertions maintain tree invariants

**Scenarios:**
- Sequential insertion (ascending/descending)
- Random insertion
- Duplicate key updates
- Max-sized keys/values

**Verifies:**
- ✅ Keys sorted after inserts
- ✅ All nodes ≤ 4096 bytes
- ✅ Data matches ref map

---

### 2. Split Trigger Tests

**File:** `TestBTreeSplitTriggers`

**Purpose:** Verify splits happen only when necessary

**Scenarios:**
- Fill node to just under PAGE_SIZE → NO split
- Add KV exceeding PAGE_SIZE → SPLIT occurs
- Cascade splits to parent
- Root split (height increase)

**Verifies:**
- ✅ Split triggered when `node.nbytes() > BTREE_PAGE_SIZE`
- ✅ All split nodes ≤ BTREE_PAGE_SIZE
- ✅ Keys sorted after split

**Implementation:**
```go
// Track page count to detect splits
before := len(c.pages)
c.add(key, largeValue)
after := len(c.pages)
splitOccurred := (after > before)
```

---

### 3. Delete Integration Tests

**File:** `TestBTreeDeleteIntegration`

**Purpose:** Verify deletions maintain tree invariants

**Scenarios:**
- Delete existing keys
- Delete non-existent keys
- Delete causing height reduction
- Delete all keys

**Verifies:**
- ✅ Deleted keys not retrievable
- ✅ Remaining keys sorted
- ✅ Data matches ref map

---

### 4. Merge Trigger Tests

**File:** `TestBTreeMergeTriggers`

**Purpose:** Verify merges happen only at underflow threshold

**Scenarios:**
- Delete until node at exactly 1024 bytes → NO merge
- Delete one more → MERGE occurs
- Merge left vs merge right
- Cannot merge (combined > PAGE_SIZE)

**Verifies:**
- ✅ Merge triggered when `node.nbytes() < BTREE_PAGE_SIZE/4`
- ✅ Merged node ≤ BTREE_PAGE_SIZE
- ✅ Keys sorted after merge (tests fix at db.go:509)

**Implementation:**
```go
// Track page count to detect merges
before := len(c.pages)
c.del(key)
after := len(c.pages)
mergeOccurred := (after < before)
```

---

### 5. Node Size Invariant Tests

**File:** `TestBTreeNodeSizeInvariants`

**Purpose:** Verify nodes never exceed size limit

**Scenarios:**
- 1000 random insert/delete operations
- Verify after EVERY operation
- Mix of large and small KV pairs

**Verifies:**
- ✅ All nodes ≤ BTREE_PAGE_SIZE at all times

**Implementation:**
```go
for i := 0; i < 1000; i++ {
    randomOperation()
    c.verifyNodeSizes(t)  // After EVERY op
}
```

---

### 6. Keys Sorted Invariant Tests

**File:** `TestBTreeKeysSortedInvariant`

**Purpose:** Verify keys remain sorted throughout operations

**Scenarios:**
- Insert random order, verify sorted
- Delete random keys, verify sorted
- Update keys, verify sorted

**Verifies:**
- ✅ Keys sorted within each node
- ✅ Keys sorted across entire tree

---

### 7. Data Integrity Tests

**File:** `TestBTreeDataIntegrity`

**Purpose:** Verify tree data matches reference map

**Scenarios:**
- Insert 1000 KV pairs
- Verify all keys retrievable with correct values
- Delete 500 keys
- Verify remaining 500 match ref

**Verifies:**
- ✅ Every key in ref exists in tree
- ✅ Every value matches ref[key]
- ✅ No phantom or missing keys

---

### 8. Stress Tests

**File:** `TestBTreeStressOperations`

**Purpose:** High-volume operations test

**Scenarios:**
- 10,000 random operations (50% insert, 30% delete, 20% update)
- Verify all invariants throughout

**Verifies:**
- ✅ No panics
- ✅ All invariants maintained
- ✅ Final state consistent

---

### 9. Tree Structure Tests

**File:** `TestBTreeTreeStructure`

**Purpose:** Verify structural properties

**Scenarios:**
- Monitor height changes
- Verify root promotion/demotion
- Verify pointer validity

**Verifies:**
- ✅ Height minimal for key count
- ✅ All pointers valid
- ✅ No dangling references

---

## Running Tests

```bash
# Run all tests
go test -v ./db

# Run only integration tests
go test -v ./db -run Integration

# Run specific test
go test -v ./db -run TestBTreeInsertIntegration
```

## Test Metrics

| Metric | Target | Measured By |
|--------|--------|-------------|
| Code Coverage | >80% | `go test -cover` |
| All Tests Pass | 100% | CI/CD |
| No Panics | 0 | Runtime |
| Memory Leaks | 0 | Page count validation |

## Key Invariants Tested

1. **Keys Sorted:** All keys in sorted order (in-order traversal)
2. **Node Size:** All nodes ≤ BTREE_PAGE_SIZE (4096 bytes)
3. **Split Trigger:** Split only when node > PAGE_SIZE
4. **Merge Trigger:** Merge only when node < PAGE_SIZE/4 (1024 bytes)
5. **Data Integrity:** Tree data matches reference map
6. **No Orphans:** All pages reachable or deallocated

## Debugging Failed Tests

If a test fails:

1. **Check which invariant failed:**
   - Keys not sorted → Check `nodeLookupLE`, insertion logic
   - Node too large → Check split logic
   - Data mismatch → Check insert/delete/update logic

2. **Use helper methods to inspect:**
   ```go
   c.verifyNodeSizes(t)      // Which node exceeded limit?
   c.verifyKeysSorted(t)     // Where is sort broken?
   c.verifyDataIntegrity(t)  // Which key has wrong value?
   ```

3. **Check tree structure:**
   ```go
   fmt.Printf("Root: %d\n", c.tree.root)
   fmt.Printf("Pages: %d\n", len(c.pages))
   fmt.Printf("Ref keys: %d\n", len(c.ref))
   ```

## Future Enhancements

- [ ] Concurrent access tests (if threading added)
- [ ] Disk I/O tests (when file backend implemented)
- [ ] Performance benchmarks
- [ ] Fuzz testing for edge cases
