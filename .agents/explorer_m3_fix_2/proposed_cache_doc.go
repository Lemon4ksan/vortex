// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cache provides workspace caching including lint memoization, encrypted secrets, and traffic storage.
//
// # Architecture Overview
//
// The cache package implements the three primary state persistence pillars of the Vortex developer toolchain:
// fast SHA-256 lint memoization, masked workspace secrets storage, and compressed gzip HTTP traffic recording:
//
//	+-------------------------------------------------------------------+
//	|                          pkg/cache Subsystems                     |
//	+-------------------------------------------------------------------+
//	         |                          |                         |
//	         v                          v                         v
//	+------------------+      +-------------------+     +------------------+
//	|    LintCache     |      |   SecretsVault    |     |  Traffic Storage |
//	| - File SHA-256   |      | - JSON Storage    |     | - gzip Session   |
//	| - Issue counts   |      | - Secret masking  |     | - Indexed by     |
//	| - .vortex/cache/ |      | - Target mapping  |     |   hash & origins |
//	+------------------+      +-------------------+     +------------------+
//
// # Core Building Blocks
//
//   - [LintCache]: Memoizes contract file hashes and diagnostic counts in `.vortex/cache/lint.json`.
//   - [LoadLintCache]: Loads or initializes the workspace lint cache from disk.
//   - [SecretsVault]: Workspace credentials and authentication token store with secret masking.
//   - [LoadSecrets]: Discovers and loads the secrets vault from a directory tree.
//   - [TrafficIndex]: Catalog of cached traffic captures in `.vortex/cache/traffic/index.json`.
//   - [TrafficEntry]: Snapshot metadata for a cached, compressed traffic session.
//   - [StoreTraffic]: Compresses and archives a HAR payload into `.vortex/cache/traffic/<hash>.har.gz`.
//   - [GetTraffic]: Retrieves and decompresses a cached traffic session by ID or hash prefix.
//   - [ListTraffic]: Returns all cached traffic sessions sorted by storage date descending.
//
// # Usage Tiers
//
// ## Tier 1: Lint Memoization
//
// Check if a file's lint state is clean and cache results:
//
//	lc, err := cache.LoadLintCache(rootDir)
//	if err != nil {
//	    log.Fatalf("failed to load lint cache: %v", err)
//	}
//	if lc.IsFresh(filePath, contentBytes) {
//	    // skip unchanged file
//	}
//	lc.Put(filePath, contentBytes, issueCount)
//	_ = lc.Save(rootDir)
//
// ## Tier 2: Secrets Vault Management
//
// Retrieve and store sensitive credentials for contract replay:
//
//	vault, vaultPath, err := cache.LoadSecrets(rootDir)
//	if err != nil {
//	    log.Fatalf("failed to load secrets: %v", err)
//	}
//	token, found := vault.Get("API_KEY")
//	if !found {
//	    vault.Set("API_KEY", "secret-value", "user")
//	    _ = vault.Save(vaultPath)
//	}
//
// ## Tier 3: Captured Traffic Storage
//
// Save and query recorded HTTP traffic sessions:
//
//	entry, secrets, err := cache.StoreTraffic(rootDir, "capture.har", harBytes, false, true)
//	if err != nil {
//	    log.Fatalf("failed to store traffic: %v", err)
//	}
//	payload, entry, err := cache.GetTraffic(rootDir, entry.ID)
//
// # Concurrency & Thread Safety
//
// Both [LintCache] and [SecretsVault] guard internal mutations using `sync.RWMutex` locks and are fully safe
// for concurrent reads and writes across goroutines. Package-level traffic functions coordinate disk
// access atomically.
//
// # Performance & Zero-Allocation Profile
//
// Traffic payloads are compressed using streaming gzip to minimize memory spikes on large HAR bodies.
// Lint caching utilizes 256-bit SHA-256 hashes formatted in-place via stack hex buffers.
package cache
