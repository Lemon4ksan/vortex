// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package jsbundle inspects minified client JS/TS bundles to extract RPC endpoints, Protobuf, and JSPB wire schemas.
//
// # Architecture Overview
//
// The scanner parses client bundles (Webpack, Vite, Rollup output chunks), detects REST, gRPC-Web, Twirp,
// and tRPC endpoints, and extracts Protobuf/JSPB message field descriptors for IR synthesis:
//
//	Minified Client Bundle (*.js, *.ts, Webpack/Vite chunks)
//	                           |
//	                           v
//	+-------------------------------------------------------------------+
//	|                      ScanFiles / ScanFile                         |
//	|  - Lexes JS token streams for RPC call patterns                   |
//	|  - Detects gRPC-Web, Twirp, tRPC, and REST fetch calls            |
//	|  - Extracts Protobuf field numbers: `jspb.Message.getField(this,1)`|
//	+-------------------------------------------------------------------+
//	                           |
//	                           v
//	+-------------------------------------------------------------------+
//	|                         ScanResult AST                            |
//	|  - Endpoints: Route path, HTTP method, request/response models    |
//	|  - Messages: Field descriptors, indices, nested sub-messages      |
//	|  - Enums: Number-to-name reverse mappings                         |
//	+-------------------------------------------------------------------+
//	                           |
//	                           v
//	+-------------------------------------------------------------------+
//	|                        ReconcileServiceIR                         |
//	|  - Bridges discovered endpoints into Vortex declarative RootIR    |
//	+-------------------------------------------------------------------+
//
// # Core Building Blocks
//
//   - [ScanFiles]: Scans multiple files matching glob patterns and aggregates their results.
//   - [ScanFile]: Scans an individual bundle file on disk.
//   - [ScanBytes]: Analyzes JavaScript byte content in-memory and extracts schema metadata.
//   - [ScanResult]: Aggregated collection of discovered endpoints, message descriptors, and enums.
//   - [NewScanResult]: Instantiates an empty [ScanResult] container.
//   - [ReconcileServiceIR]: Enriches a ServiceIR using descriptors discovered in JavaScript bundles.
//   - [Endpoint]: Describes an individual discovered API route or RPC procedure.
//   - [MessageDescriptor]: Models a discovered Protobuf, JSPB, or DTO wire schema.
//   - [FieldDescriptor]: Models a specific struct or message field index and type.
//   - [EnumDescriptor]: Models a discovered integer-to-string enum mapping.
//
// # Usage Tiers
//
// ## Tier 1: Glob Pattern & File Scanning
//
// Scan client build bundles for hidden endpoints and schemas:
//
//	res, err := jsbundle.ScanFiles([]string{"frontend/dist/*.js", "frontend/dist/chunks/*.js"})
//	if err != nil {
//	    log.Fatalf("scanning failed: %v", err)
//	}
//	log.Printf("discovered %d endpoints", len(res.Endpoints))
//
// ## Tier 2: Combining Results from Chunks
//
// Concurrently scan multiple bundle chunks and merge findings:
//
//	combined := jsbundle.NewScanResult()
//	combined.Merge(chunkResult1)
//	combined.Merge(chunkResult2)
//
// ## Tier 3: In-Memory Byte Stream Scanning
//
// Scan streamed bundle chunks directly from memory:
//
//	res := jsbundle.ScanBytes(bundleBytes, "app.min.js")
//	log.Printf("discovered %d endpoints in chunk", len(res.Endpoints))
//
// # Concurrency & Thread Safety
//
// Scanning functions ([ScanFiles], [ScanFile], [ScanBytes]) are stateless and safe for parallel
// execution. [ScanResult.Merge] must be called with external synchronization if shared across goroutines.
//
// # Performance & Zero-Allocation Profile
//
// Streaming chunk processing minimizes memory consumption when scanning massive 100MB+ vendor JavaScript bundles.
package jsbundle
