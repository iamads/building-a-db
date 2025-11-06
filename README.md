building a db in golang

follow along: https://build-your-own.org/database/#table-of-contents

----------------------

B+tree principles:

- n-ary tree, node size is limited by a constant.
- Same height for all leaves.
- Split and merge for insertion and deletion.

----------------------

# Goal 1

- [ ] Building a in memory B+tree, which allows us to do insert, delete, get operation, range queries etc
  - [x] We use separate intermediate nodes, and leaf nodes .... should I also have a seaparate root node?
  - [x] We should define node size as constant,
  - [ ]  Later on move to allowing the user to decrare the size
  - [x] Pretty print the tree
  - [ ] Insert
    - [x] Find the correct intermediate node: only 1 level except root
    - [x] If intermediate node does not exist insert intermediate node: only 1 level except root
      - [x] At this point we need to think about should we redistribute intermediate nodes or not
    - [ ] If insertion of intermediate node  causes us to hit limit on parent, we have to introduce another level: Multilevel
      - [ ] ~~Do we rebalance nodes b/w different nodes?~~
  - [x] Get
    - [x] Get element by particular id
    - [ ] get element by range queries
      - [ ] To implemment range quueries: we should have leaf nodes like a linkekd. list (this is a point of B+tree)
      - [x] We need to channge how our inserts will work
        - [x] We will have to change leaf node struuct
        - [x] In case of insertion if we do not have any Leaf nodes we will have to mark Next as nil
        - [x] In case we are inserting after some leaf nodes we. will have to update the next of preceding leaf node and current leaf node
        - [x] In case we are inserting before are leaf nodes we will have to set next of current leaf node to first node
        - [ ] ~~Do we allslo handle this in case of inode~~
        - [x] Create insertion test cases where we check if next is correctly maintained
  - [ ] Delete
    - [ ] Delete key
  - [ ] Update
    - [ ] Update Element
  

# Goal 2
  - [ ] Model node struct as shown in the book
  - [ ] ..... Will dig deeper later


# Future goals
  - [ ] Multi process synchronisation
  - [ ] using disk to store data
  - [ ] parsing sql
  - [ ] JSON based storage
  - [ ] Analyse different DB storage engine
    - [ ] InnoDb
    - [ ] WiredTiger
    - [ ] BitCask
    - [ ] ...
  - [ ] 



# Notes
- reference this :https://github.com/cockroachdb/cockroach/blob/master/CLAUDE.md
- bazel ? -> https://www.youtube.com/watch?v=YiX0NpKL7ag
  