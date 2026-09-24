// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mirror synchronizes declarative contracts with untagged upstream Go source files via @mirror directives.
//
// # Architecture Overview
//
// The shadow mirror engine parses target Go source files referenced by `@mirror` directives, aligns method
// signatures, parameter types, and return values with declarative contracts, and identifies drift diagnostics:
//
//	Contract Interface (@mirror: "internal/legacy/service.go")
//	                 \                          /
//	                  v                        v
//	+-------------------------------------------------------------------+
//	|                           CheckService                            |
//	|  - Parses target Go AST at mirror path                            |
//	|  - Compares method signatures, parameter types, and return values |
//	|  - Compares DTO structs referenced by return types                |
//	+-------------------------------------------------------------------+
//	                                 |
//	                                 v
//	+-------------------------------------------------------------------+
//	|                       DriftDiagnostic Report                      |
//	|  - DriftMethodMissing: Missing in root source                     |
//	|  - DriftParamMismatch: Parameter type changed                     |
//	|  - DriftGhostMethod: Defined upstream, not exposed in wrapper     |
//	|  - DriftFieldMismatch: Struct field type discrepancy              |
//	+-------------------------------------------------------------------+
//
// # Core Building Blocks
//
//   - [CheckService]: Inspects an individual service contract against its target mirror source file.
//   - [DriftDiagnostic]: Structured diagnostic detailing a specific signature or field discrepancy.
//   - [DriftKind]: Categorization of mirror drift ([DriftMethodMissing], [DriftParamMismatch], [DriftGhostMethod], etc.).
//
// # Usage Tiers
//
// ## Tier 1: Service Drift Inspection
//
// Check an individual contract for synchronization drift:
//
//	diagnostics, err := mirror.CheckService(rootDir, contractPath, serviceIR, structs)
//	if err != nil {
//	    log.Fatalf("mirror check failed: %v", err)
//	}
//	for _, diag := range diagnostics {
//	    log.Printf("[%s] %s: %s", diag.Kind, diag.Service, diag.Message)
//	}
//
// ## Tier 2: Workspace-Wide Drift Auditing
//
// Audit all mirrored services across the parsed contract tree:
//
//	for _, svc := range rootIR.Services {
//	    diags, err := mirror.CheckService(rootDir, contractPath, svc, rootIR.Structs)
//	    if err != nil {
//	        log.Printf("error checking service %s: %v", svc.Name, err)
//	        continue
//	    }
//	    // process service diagnostics
//	}
//
// ## Tier 3: Diagnostic Filtering & Validation
//
// Filter diagnostics by severity to enforce strict synchronization in CI pipelines:
//
//	for _, diag := range diagnostics {
//	    if diag.Severity == "error" {
//	        log.Fatalf("blocking mirror drift in %s: %s", diag.Service, diag.Message)
//	    }
//	}
//
// # Concurrency & Thread Safety
//
// [CheckService] performs read-only AST comparisons and is 100% thread-safe across concurrent goroutines.
//
// # Performance & Zero-Allocation Profile
//
// The engine parses target mirror files into lightweight ASTs, comparing type signatures through
// fast identifier matching with minimal heap allocations.
package mirror
