<div align="center">

# vortex

### The Zero-Allocation AST Toolchain & Sovereign Contract Engine

_«A little copying is better than a little dependency — own your network contracts down to the AST»_

[![Go Version](https://img.shields.io/badge/go-1.27%2B-007d9c?logo=go&logoColor=white&style=flat-square)](https://go.dev/)
[![Go Reference](https://img.shields.io/badge/godoc-reference-007d9c?style=flat-square)](https://pkg.go.dev/github.com/lemon4ksan/vortex)
[![License](https://img.shields.io/badge/license-BSD--3--Clause-blue?style=flat-square)](LICENSE)
[![Zero-Alloc](https://img.shields.io/badge/memory-0%20B%2Fop%20%7C%200%20allocs-brightgreen?style=flat-square)](docs/VORTEX.md)
[![OpenAPI 3.1](https://img.shields.io/badge/spec-OpenAPI%203.1%20%26%20AsyncAPI-blueviolet?style=flat-square)](pkg/openapi)
[![AST Engine](https://img.shields.io/badge/compiler-Go%20AST%203--Way%20Merge-orange?style=flat-square)](pkg/parser)
[![Oracle Hub](https://img.shields.io/badge/attestation-Universal%20Oracle%20v2-yellow?style=flat-square)](pkg/oracle/gen)

**vortex** is the unified declarative contract toolchain, code generator, and traffic-driven reverse-engineering suite for [`aoni`](https://github.com/lemon4ksan/aoni). It operates directly on Go Abstract Syntax Trees (AST), treating idiomatic Go interface declarations as the single source of truth for REST, WebSocket, SSE, OpenAPI 3.1, AsyncAPI, and Protocol Buffer communications.

#### English • [Русский](README_RU.md) • [Architecture Specification](docs/VORTEX.md) • [License](LICENSE)

</div>

---

## The Vortex Manifesto: Sovereign Clients

For over two decades, network client development in software engineering has suffered from fragmentation:
* Standard `net/http` is wrapped in unmaintained third-party TLS forks to bypass basic bot filters.
* Brittle headless browser scripts are cobbled together to solve challenges, consuming gigabytes of RAM.
* Heavy, reflection-based third-party SDKs pollute dependency graphs, dragging in breaking upgrade cycles.
* When upstream platforms update internal endpoints, engineering teams are blocked waiting for external maintainers.

**The Sovereign Paradigm**: No engineering team should ever depend on third-party API wrapper libraries. Every production project must own its sovereign, zero-allocation API client generated directly into its own codebase (`pkg/api/`) from Go AST contracts, OpenAPI schemas, or live network traffic captures (`.har`).

---

## Key Capabilities

| Capability | Description | Command |
| :--- | :--- | :--- |
| **AST Compilation** | Generates reflection-free, zero-allocation Go clients from declarative interfaces. | `vortex gen` |
| **Traffic Reverse Engineering** | Ingests `.har` network captures or OpenAPI 3.1 specs with 3-way AST merge. | `vortex spec import` |
| **Universal Oracle v2** | Compiles browser attestation sidecars to bypass WAFs and Cloudflare invisibly. | `vortex oracle` |
| **Contract Quality & Linting** | Static linter verifying AST directives, route consistency, and types. | `vortex check --fix` |
| **In-Memory Mocking** | Generates zero-dependency virtual test servers for isolated unit & integration tests. | `vortex mock` |
| **Live Smoke Probing** | Rapidly tests live endpoints using contract secrets and renders latency tables. | `vortex smoke` |
| **Traffic Inspector UI** | Local UI for inspecting WebSocket frames, analyzing payloads, and debugging JA4 TLS. | `vortex traffic` |
| **AST Borrow Checker** | Verifies zero-allocation constraints and memory lifetimes across contract code. | `vortex borrow` |
| **Silicon Benchmarking** | Measures throughput, allocations, and CPU profiles against production baselines. | `vortex bench` |

---

## Installation

`vortex` requires Go version `1.27` or higher.

```bash
go install github.com/lemon4ksan/vortex/cmd/vortex@latest
```

Verify installation:
```bash
vortex --version
```

---

## Workflow & Examples

### 1. Declare Go Contract Interface

Define your API contract using clean Go interfaces and `@aoni` Godoc directives:

```go
// Package api defines the sovereign GitHub API contract.
package api

import (
	"context"
)

type User struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Bio   string `json:"bio"`
}

// GitHubService defines the declarative API contract.
//
// @aoni:service GitHubService
// @aoni:baseURL https://api.github.com
type GitHubService interface {
	// GetUser retrieves public user profile information.
	//
	// @aoni:method GET /users/{username}
	// @aoni:header Accept: application/vnd.github.v3+json
	// @aoni:cache TTL=5m
	GetUser(ctx context.Context, username string) (*User, error)
}
```

### 2. Compile Zero-Allocation Client (`vortex gen`)

Compile the interface into a high-performance `aoni`-backed client without runtime reflection:

```bash
# Compile contracts across the project
vortex gen ./...
```

Vortex generates `github_service.vortex.go` implementing the interface with exact wire layouts, zero heap allocations on the hot path, and built-in connection pooling.

### 3. Traffic Ingestion & Reverse Engineering (`vortex spec import`)

Convert real browser or mobile network captures directly into typed Go contracts:

```bash
# Record live traffic from browser into a HAR archive
vortex traffic record -out=session.har

# Ingest HAR into Go contracts with automatic schema inference and 3-way merge
vortex spec import -spec=session.har -pkg=api -add
```

Vortex automatically identifies endpoints, extracts URL variables, infers JSON request/response models, and merges changes into existing Go interfaces without overwriting manual edits.

### 4. Browser Attestation Sidecar (`vortex oracle`)

Deploy an automated browser attestation sidecar for APIs protected by Cloudflare Turnstile or anti-bot challenges:

```bash
# Compile an Oracle sidecar for a target web application
vortex oracle -name=portal https://portal.example.com
```

The Oracle sidecar orchestrates headless Chromium instances, solves challenges, extracts session state and `cf_clearance` cookies, and streams dynamic tokens directly into your sovereign client via IPC.

### 5. Instant Test Mock Server (`vortex mock`)

Generate a zero-dependency, in-memory mock HTTP server for integration testing:

```bash
vortex mock ./pkg/api/github_service.go -out=./pkg/api/mock_test.go
```

---

## CLI Command Reference

### Daily Core Commands

```bash
vortex                                 # Auto-pilot: audit, synchronize, and compile all contracts
vortex gen [path]                      # Compile zero-allocation Go clients from AST
vortex check [path]                    # Lint and validate contract directives (--fix to auto-repair)
vortex mock [path]                     # Generate virtual in-memory mock server for testing
vortex smoke [path]                    # Probe live endpoints using contract secrets
vortex env [path]                      # Scan contracts for ${VAR} and generate .env templates
```

### Domain Hubs

```bash
vortex spec import -spec=api.json      # Ingest OpenAPI 3.1 / HAR into Go contracts
vortex spec export -out=openapi.json   # Export Go contracts to OpenAPI 3.1 specification
vortex spec diff -spec=v2.json         # Compare local contract with remote spec for breaking changes
vortex oracle -name=bot https://...    # Compile browser attestation sidecar
vortex traffic                         # Launch local web Traffic Inspector UI
vortex ast split                       # Refactor and split monolithic interfaces
vortex perf bench                      # Run hardware benchmarks and allocation diagnostics
vortex borrow ./...                    # Run AST borrow checker for zero-alloc invariants
```

### Workspace Management

```bash
vortex init [name]                     # Initialize .vortex.yml workspace or contract template
vortex status                          # Show 360° synchronization health of workspace contracts
vortex doctor                          # Diagnose toolchain health, paths, and git synchronization
vortex explain @aoni:cache             # Display documentation and examples for a directive
vortex clean                           # Clean generated artifacts, profiles, and harnesses
```

---

## Architecture: The 6-Layer Engine

```
       +-------------------------------------------------------------+
       | 1. Ergonomic Surface: Declarative Go Interfaces & Godoc     |
       +------------------------------+------------------------------+
                                      |
       +------------------------------v------------------------------+
       | 2. Toolchain Layer: Vortex AST Engine & 3-Way Merge         |
       +------------------------------+------------------------------+
                                      |
       +------------------------------v------------------------------+
       | 3. Attestation Layer: Universal Oracle v2 (Browser Sidecars)|
       +------------------------------+------------------------------+
                                      |
       +------------------------------v------------------------------+
       | 4. Execution Core: aoni.Client & fast.Client (1.5M+ RPS)    |
       +------------------------------+------------------------------+
                                      |
       +------------------------------v------------------------------+
       | 5. Network Fidelity: L4 p0f Spoofing, uTLS, HTTP/2 Framing  |
       +------------------------------+------------------------------+
                                      |
       +------------------------------v------------------------------+
       | 6. Silicon Sympathy: 64-Byte Cache Lines, SIMD, Arenas      |
       +-------------------------------------------------------------+
```

---

## License

This project is licensed under the **BSD 3-Clause License**. See the [LICENSE](LICENSE) file for details.

Copyright (c) 2026 Lemon4ksan. All rights reserved.
