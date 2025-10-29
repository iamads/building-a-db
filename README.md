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
  - [ ] We use separate intermediate nodes, and leaf nodes .... should I also have a seaparate root node?
  - [ ] We should define node size as constant, but later on move to allowing the user to decrare the size
  - [ ] Pretty print the tree
  - [ ] Insert
    - [ ] Find the correct intermediate node: only 1 level except root
    - [ ] If intermediate node does not exist insert intermediate node: only 1 level except root
      - [ ] At this point we need to think about should we redistribute intermediate nodes or not
    - [ ] If insertion of intermediate node  causes us to hit limit on parent, we have to introduce another level: Multilevel
      - [ ] Do we rebalance nodes b/w different nodes?
  - [ ] Get
    - [ ] Get element by particular id
    - [ ] get element by range queries
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



