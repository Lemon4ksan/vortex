// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cfg constructs an optimized, lightweight Control Flow Graph (CFG) for Go AST.
//
// # Architecture Overview
//
// The CFG engine splits Go statement trees into sequential basic blocks, computes predecessor and successor
// edges across conditional branches, switches, and loops, and performs reachability and path analysis:
//
//	Go Function Body (ast.BlockStmt)
//	               |
//	               v
//	+-------------------------------------------------------------------+
//	|                             cfg.New                               |
//	|  - Splits statements into basic sequential Blocks                 |
//	|  - Resolves IfStmt, SwitchStmt, ForStmt, BranchStmt jumps        |
//	|  - Computes Predecessor (Preds) and Successor (Succs) edges       |
//	+-------------------------------------------------------------------+
//	               |
//	               v
//	+-------------------------------------------------------------------+
//	|                           CFG Graph                               |
//	|  Block 0 (Entry) ---> Block 1 (Condition) --True--> Block 2 (Body)|
//	|                             |                          |          |
//	|                           False                        v          |
//	|                             +--------------------> Block 3 (Exit) |
//	+-------------------------------------------------------------------+
//	               |
//	               v
//	+-------------------------------------------------------------------+
//	|                      Path & Loop Analysis                         |
//	|  - WalkPaths: acyclic execution path traversal                    |
//	|  - FindLoopBlocks: identification of iterative control blocks    |
//	|  - ReturnBlocks & NoReturn reachability validation                |
//	+-------------------------------------------------------------------+
//
// # Core Building Blocks
//
//   - [New]: Constructs a fully resolved [CFG] from an `ast.BlockStmt` or function body.
//   - [CFG]: Represents the complete graph structure of indexed basic blocks and return points.
//   - [Block]: Represents an individual sequential basic block of statements and outgoing edges.
//   - [BlockKind]: Categorizes block semantics (e.g. entry, branch condition, loop header, exit).
//   - [PathVisitor]: Callback invoked for each complete acyclic execution path during path traversal.
//
// # Usage Tiers
//
// ## Tier 1: Graph Construction
//
// Construct a CFG from a Go function declaration body:
//
//	graph := cfg.New(funcDecl.Body, nil)
//	entry := graph.Entry()
//	returns := graph.ReturnBlocks()
//
// ## Tier 2: Block Traversal & Reachability
//
// Inspect basic blocks, determine live blocks, and traverse control paths:
//
//	for _, block := range graph.Blocks {
//	    if !block.Live {
//	        // dead code detected
//	    }
//	    for _, succ := range block.Succs {
//	        _ = succ.Index
//	    }
//	}
//
// ## Tier 3: Execution Path Traversal & Loop Inspection
//
// Execute path walking and loop detection for compiler optimizations and borrow checking:
//
//	graph.WalkPaths(func(path []*cfg.Block) {
//	    // inspect sequence of blocks along execution path
//	})
//	loopBlocks := graph.FindLoopBlocks()
//
// # Concurrency & Thread Safety
//
// [CFG] instances are strictly immutable once constructed by [New]. They are 100% thread-safe
// for concurrent graph queries, path traversals, and inspections across multiple goroutines.
//
// # Performance & Zero-Allocation Profile
//
// Basic block branch pointers utilize inline small-array buffers (`succs2 [2]*Block`) for two-way
// conditional branches, eliminating slice heap allocations on typical if-else control structures.
package cfg
