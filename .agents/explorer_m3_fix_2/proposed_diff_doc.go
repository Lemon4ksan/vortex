// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package diff implements semantic contract drift analysis between local Go interfaces and OpenAPI specifications.
//
// # Architecture Overview
//
// The diff engine computes semantic discrepancies between local declarative contract ASTs and remote
// OpenAPI/AsyncAPI specifications. Discrepancies are categorized by severity (Breaking, NonBreaking, Ghost),
// while a persistent [DiffStack] maintains snapshot history across contract iterations:
//
//	Local RootIR Contract           Remote OpenAPI Specification
//	          \                                /
//	           v                              v
//	    +--------------------------------------------+
//	    |                  Compare                   |
//	    |  - Normalizes paths (/users/{id} <-> :id)  |
//	    |  - Aligns operations by HTTP verb and path |
//	    +--------------------------------------------+
//	                          |
//	                          v
//	    +--------------------------------------------+
//	    |                DiffReport                  |
//	    |  - SeverityBreaking: Missing endpoint,     |
//	    |    incompatible type, missing param        |
//	    |  - SeverityNonBreaking: New optional field |
//	    |  - SeverityGhost: Removed endpoint         |
//	    +--------------------------------------------+
//	                          |
//	                          v
//	    +--------------------------------------------+
//	    |                 DiffStack                  |
//	    |  - Push, DiffFrames snapshot undo frames   |
//	    +--------------------------------------------+
//
// # Core Building Blocks
//
//   - [Compare]: Evaluates drift between a local [*ir.RootIR] contract and an [*openapi.Document].
//   - [CompareWithOptions]: Evaluates drift with customized comparison options.
//   - [DiffReport]: Aggregates all detected discrepancies and summary metrics.
//   - [DiffOptions]: Configures path normalization, case sensitivity, and tolerance rules.
//   - [DriftItem]: Individual contract difference record detailing paths, parameters, and severity.
//   - [DriftSeverity]: Classification of drift impact ([SeverityBreaking], [SeverityNonBreaking], [SeverityGhost]).
//   - [DiffStack]: In-memory and persistent snapshot stack managing iterative contract checkpoints in `.vortex/cache/diff_stack.json`.
//   - [LoadStack]: Loads or initializes the diff snapshot stack for a workspace root.
//   - [StackFrame]: Individual historical frame inside the snapshot stack.
//
// # Usage Tiers
//
// ## Tier 1: High-Level Contract Comparison
//
// Compare a local Go interface AST with a parsed remote OpenAPI document:
//
//	report := diff.Compare(localRootIR, remoteDoc, "local.go", "remote.json")
//	if report.HasBreaking() {
//	    log.Println("breaking changes detected")
//	}
//
// ## Tier 2: Granular Options & Filtering
//
// Apply customized comparison options such as enabling additive comparison mode:
//
//	opts := diff.DiffOptions{Additive: true}
//	report := diff.CompareWithOptions(localRootIR, remoteDoc, "local.go", "remote.json", opts)
//
// ## Tier 3: Checkpoint Stack Management
//
// Push and inspect contract snapshots during automated reconciliation workflows:
//
//	stack, err := diff.LoadStack(rootDir)
//	if err != nil {
//	    log.Fatalf("failed to load stack: %v", err)
//	}
//	frame, err := stack.Push("pre-merge", []string{"pkg/api/service.go"}, []string{"release"}, nil)
//	_ = stack.Save()
//
// # Concurrency & Thread Safety
//
// [Compare] and [CompareWithOptions] are pure, stateless functions safe for concurrent execution across
// multiple goroutines. [DiffStack] operations are synchronized via internal read-write mutexes.
//
// # Performance & Zero-Allocation Profile
//
// Path alignment and verb matching utilize stack-allocated string builders and pre-hashed endpoint maps,
// ensuring comparison passes complete in microseconds without excessive garbage collector pressure.
package diff
