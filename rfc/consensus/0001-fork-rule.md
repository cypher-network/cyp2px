+ Feature name: `fork-rule`
+ Start date: 2021-06-05

# Summary

The fork rule defined the strategy for handling valid new blocks which are minted / accepted by the network, but may conflict with other minted / accepted blocks due to the distributed nature of the network.

# Motivation

An honest node follows a fork choice rule prescribed in the protocol. The purpose is to reach a consensus on which chain of blocks to maintain, in a distributed manner. Eventually, such chosen chain of blocks produces a final ledger of transactions. Under the Nakamoto protocol, a node appends the next generated block to the longest chain in its local view of the block tree. Unlike PoW systems, Nakamoto protocol can lead to serious security issues for PoS systems.

The fork rule described here shall mitigate these security issues and shall be followed by all honest cypher nodes.

# Detailed design

Definitions:

* `Block` - Smallest element in the collection of blocks forming the blockchain.

* `Block height` - 0-indexed position of a `Block` in the list of blocks forming a chain.

* `Chain` - List of `Block`s from `Block height` 0 up to the last accepted `Block`.

* `Chain height` - 1-indexed number of `Blocks` in a `Chain`.

* `Fork root` - Position in `Chain` where more than 1 `Block`s are valid.

* `Fork branch` - List of `Block`s originating from `Fork root` up to the last accepted `Block` in the branch.

* `Minimum Fork length` - `Fork branch` length threshold for consensus.

* `Block creation time` - Time needed to mint a `Block`.

* `Fork creation time` - Total time needed to mint all `Block`s in a given `Fork branch`.

* `Fork rule` - Set of rules to follow for the `Chain` selection given a collection of `Block`s.

## Ideal chain

The ideal chain is a chain with no forks. Let 0 be the chain's first block with height 0 and 5 be the chain's last accepted block with height 5. In example 01, the chain length is 6:


```
------------
  example 01

    1) 0 -- 1 -- 2 -- 3 -- 4 -- 5
------------
```


## Fork creation time

When distributed nodes synchronize their local chains with the rest of the network, this may lead to a fork. In example 02 there exist two branches where the block with height 4 is the fork root:


```
------------
  example 02

    1) 0 -- 1 -- 2 -- 3 -- 4 -- 5 -- 6 -- 7 -- 8 (17s)
    2)                      \-- 5 -- 6 -- 7 -- 8 (18s)
------------
```


The node uses `Rule 1` to decide which is the currently accepted chain.

**Rule 1** Fork branches with the shortest `Fork creation time` for the first `Minimum Fork length` blocks become the selected chain.

We assume a `Minimum fork length` of 4 in this document. According to `Rule 1, fork branch 1 is the currently selected chain in example 02, with a `Fork creation time` of 17 seconds.

This applies to forks of forks as well. In example 03 there is a second fork, originating from branch 2, block height 6. Branch 3 is the currently selected chain in that example, given the shortest `Fork creation time` of 16 seconds and a fork length of 4.


```
------------
  example 03

    1) 0 -- 1 -- 2 -- 3 -- 4 -- 5 -- 6 -- 7 -- 8 (17s)
    2)                      \-- 5 -- 6 -- 7 -- 8 (18s)
    3)                                \-- 7 -- 8 (16s)
------------
```


## Minimum Fork lengh

To prevent `Rule 1` from preventing any new blocks to become part of the selected chain, `Rule 2` defines the consensus for forks with a length smaller than `Minimum Fork length`:

**Rule 2** The longest fork branch with a length smaller than `Minimum Fork length` becomes the selected chain.

In example 04, both branches have a length smaller than 4. Therefore, branch 1 is the currently selected chain due to its length:


```
------------
  example 04

    1) 0 -- 1 -- 2 -- 3 -- 4 -- 5 -- 6
    2)                      \-- 5
------------
```


This applies to more complex fork structures as well. In example 05, branch 2 is the selected chain:


```
------------
  example 05

    1) 0 -- 1 -- 2 -- 3 -- 4 -- 5
    2)                 \-- 4 -- 5 -- 6
    3)                  \- 4 -- 5
    4)                      \-- 5
------------
```


## Pseudo algorithm

Both rules are applied starting from block 0, up to the end of the chain. Example 06 illustrates a situation where both rules need to be applied several times:


```
------------
  example 06

    1) 0 -- 1 -- 2 -- 3 -- 4 -- 5 -- 6 -- 7 -- 8 (17s) -- 9 -- 10
    2)                      \-- 5 -- 6 -- 7
    3)                                \-- 7 -- 8 (16s) -- 9 -- 10 -- 11 -- 12 (20s)
    4)                                                     \-- 10 -- 11       (18s)
------------
```


The evaluation is performed in the following order:

1. Start at block N (0)

2. Find the first fork root starting from block N

3. Evaluate all forks from fork root 4:

     1. Branch 1 and branch 3 have a minimum length of `Minimum fork length`, thus branch 2 is ignored and `Rule 1` is applied to branch 1 and 3.

     2. Branch 3 has the shortest `Fork creation time` of both branches, and becomes the selected chain.

     3. Continue search at block N (8)

4. Find the first fork root starting from block N

5. Evaluate all forks from fork root 9

     1. No branch has a minimum length of `Minimum fork length`, this `Rule 2` is applied to all branches.

     2. Branch 3 is the longest branch and thus becomes the selected chain.

     3. Continue search at block N (12)

6. Find the first fork root starting from block N

7. No more forks found, branch 3 is the selected chain with `Chain height` 13 (blocks 0 to 12).


A few moments later, the node has received some updates from the network and comes to the collection of blocks as depicted in example 7:


```
------------
  example 07

    1) 0 -- 1 -- 2 -- 3 -- 4 -- 5 -- 6 -- 7 -- 8 (17s) -- 9 -- 10
    2)                      \-- 5 -- 6 -- 7 -- 8 (15s)
    3)                                \-- 7 -- 8 (16s) -- 9 -- 10 -- 11 -- 12 (20s)
    4)                                                     \-- 10 -- 11       (18s)
------------
```


The evaluation is performed in the following order:

1. Start at block N (0)

2. Find the first fork root starting from block N

3. Evaluate all forks from fork root 4:

     1. All branches have a minimum length of `Minimum fork length`, thus `Rule 1` is applied to branches 1, 2 and 3.

     2. Branch 2 has the shortest `Fork creation time` of both branches, and becomes the selected chain.

     3. Continue search at block N (8)

4. No more forks found, branch 2 is the selected chain with `Chain height` 9 (blocks 0 to 8).