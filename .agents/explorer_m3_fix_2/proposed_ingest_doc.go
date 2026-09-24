// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ingest parses W3C HAR 1.2 traffic archives and synthesizes OpenAPI 3.x specifications.
//
// # Architecture Overview
//
// The ingest engine decodes network traffic captures from browsers or HTTP proxies, deduplicates repetitive
// endpoints, heuristics-infers dynamic path variables, and constructs valid OpenAPI 3.x specifications:
//
//	Browser / Proxy Network Export (*.har)
//	                  |
//	                  v
//	+-------------------------------------------------------------------+
//	|                           HARToOpenAPI                            |
//	|  - Decodes JSON log into HARLog and HAREntry stream               |
//	|  - Extracts Request, Response, Headers, and PostData              |
//	+-------------------------------------------------------------------+
//	                  |
//	                  v
//	+-------------------------------------------------------------------+
//	|                      Deduplicate & Parameterize                   |
//	|  - Clusters endpoints by URL structure                            |
//	|  - Heuristically infers path variables (/items/123 -> /items/{id})|
//	|  - Detects JSON / Form / Binary content types                     |
//	+-------------------------------------------------------------------+
//	                  |
//	                  v
//	+-------------------------------------------------------------------+
//	|                         OpenAPI Synthesis                         |
//	|  - Emits valid OpenAPI 3.x Document AST (*openapi.Document)       |
//	+-------------------------------------------------------------------+
//
// # Core Building Blocks
//
//   - [HARToOpenAPI]: Transforms raw W3C HAR 1.2 JSON bytes into an [*openapi.Document].
//   - [HARToOpenAPIOpts]: Ingests HAR logs with custom ignore patterns and route template rewrites via [IngestOptions].
//   - [DetectFormat]: Inspects raw payload bytes to identify OpenAPI, Swagger, Postman, or HAR format.
//   - [HARLog]: Root W3C HAR 1.2 container.
//   - [HAREntry]: Individual captured request-response transaction.
//   - [HARNV]: Key-value pair representing headers or query parameters.
//   - [HARPostData]: Request body payload model.
//   - [HARContent]: Response body payload model.
//   - [IngestOptions]: Configuration options for HAR ingestion rules and filters.
//
// # Usage Tiers
//
// ## Tier 1: Automated Spec Synthesis
//
// Ingest a browser HAR capture and generate an OpenAPI document:
//
//	doc, err := ingest.HARToOpenAPI(harBytes)
//	if err != nil {
//	    log.Fatalf("failed to synthesize OpenAPI from HAR: %v", err)
//	}
//	_ = doc.Paths
//
// ## Tier 2: Inspecting Captured Transactions
//
// Traverse parsed HAR entries to inspect headers and query parameters:
//
//	var harLog ingest.HARLog
//	if err := json.Unmarshal(harBytes, &harLog); err != nil {
//	    log.Fatalf("failed to decode HAR: %v", err)
//	}
//	for _, entry := range harLog.Log.Entries {
//	    log.Printf("[%s] %s -> %d", entry.Request.Method, entry.Request.URL, entry.Response.Status)
//	}
//
// ## Tier 3: Path Variable Heuristics & Filters
//
// The engine automatically detects UUIDs, integer IDs, and hashes in URL segments and parameterizes them
// into standard OpenAPI path template variables (e.g. `/v1/users/{userId}`). Custom rules are passed via [IngestOptions]:
//
//	opts := ingest.IngestOptions{
//	    IgnorePatterns: []string{"*.png", "*.css", "*.js"},
//	    RouteTemplates: []string{"/api/v1/users/{userId}"},
//	}
//	doc, err := ingest.HARToOpenAPIOpts(harBytes, opts)
//
// # Concurrency & Thread Safety
//
// [HARToOpenAPI], [HARToOpenAPIOpts], and [DetectFormat] are stateless pure functions. They are 100% thread-safe
// across concurrent goroutines.
//
// # Performance & Zero-Allocation Profile
//
// The parser processes transaction records sequentially, reusing internal string buffers to convert
// URL paths and header maps with minimal garbage collector pressure.
package ingest
