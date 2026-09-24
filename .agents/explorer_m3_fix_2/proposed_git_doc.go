// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package git provides in-memory Git repository inspection querying branches, logs, and files with zero disk artifacts.
//
// # Architecture Overview
//
// The git package executes targeted Git CLI operations bounded by strict timeouts ([DefaultTimeout]), capturing
// file revisions directly into memory buffers without checking out worktrees or writing temporary disk files:
//
//	Git Working Tree / Repository
//	              |
//	              v
//	+-------------------------------------------------------------------+
//	|                   exec.CommandContext ("git", ...)                |
//	|  - Enforces DefaultTimeout (5s)                                   |
//	|  - Captures stdout directly into in-memory bytes.Buffer           |
//	+-------------------------------------------------------------------+
//	     |                   |                   |                  |
//	     v                   v                   v                  v
//	+----------+       +-----------+       +-----------+      +-----------+
//	| ShowFile |       |LogCommits |       |ListProposal|     |  IsClean  |
//	| (ref:path|       | (History) |       | (Branches)|      | (Status)  |
//	+----------+       +-----------+       +-----------+      +-----------+
//
// # Core Building Blocks
//
//   - [ShowFile]: Extracts the exact byte content of a file at a specific Git ref into an in-memory byte slice.
//   - [RootDir]: Resolves the top-level repository root directory.
//   - [CurrentBranch]: Retrieves the active Git branch name.
//   - [LogCommits]: Retrieves commit history metadata for a specific file or path.
//   - [ListProposalBranches]: Discovers local and remote proposal branches matching consumer prefixes.
//   - [IsClean]: Checks whether uncommitted changes exist in the working tree for a specific path or entire repository.
//   - [CommitInfo]: Commit metadata record containing hash, author, timestamp, and subject.
//   - [BranchProposal]: Structured descriptor for consumer feature branches proposing contract updates.
//
// # Usage Tiers
//
// ## Tier 1: In-Memory File Retrieval
//
// Read an earlier revision of an API contract without writing to disk:
//
//	data, err := git.ShowFile(ctx, rootDir, "HEAD~1", "pkg/api/service.go")
//	if err != nil {
//	    log.Fatalf("failed to retrieve file from git: %v", err)
//	}
//
// ## Tier 2: Commit History & Branch Inspection
//
// Inspect recent commits affecting a contract file or discover active branches:
//
//	commits, err := git.LogCommits(ctx, rootDir, "pkg/api/service.go", 5)
//	branches, err := git.ListProposalBranches(ctx, rootDir, nil)
//
// ## Tier 3: Working Tree Cleanliness
//
// Verify that the repository is clean before applying automated migrations:
//
//	isClean, err := git.IsClean(ctx, rootDir, "")
//
// # Concurrency & Thread Safety
//
// All functions in this package accept [context.Context] and invoke isolated CLI subprocesses with
// dedicated output buffers. They are 100% thread-safe across concurrent goroutines.
//
// # Performance & Zero-Allocation Profile
//
// File retrieval streams subprocess output directly into pre-allocated memory buffers, completely
// bypassing temporary disk files and intermediate working tree checkouts.
package git
