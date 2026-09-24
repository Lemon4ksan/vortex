// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package parser extracts Go interface declarations, structs, and doc comment directives into an Unchecked IR.
//
// # Architecture Overview
//
// The parser inspects Go source files using standard go/parser and go/ast tools, tokenizing compiler directives
// prefixed with `@` in Godoc comments, validating them against the directive specification registry, and binding
// resolved types into an [*ir.RootIR] representation:
//
//	+-------------------------------------------------------------------+
//	|                         Go Source Files (.go)                     |
//	+-------------------------------------------------------------------+
//	                                   |
//	                                   v
//	                        +---------------------+
//	                        | go/parser & go/ast  |
//	                        +---------------------+
//	                                   |
//	                                   v
//	                        +---------------------+
//	                        |   Directive Lexer   |  <--- Tokenizes @service, @get, @unwrap
//	                        +---------------------+
//	                                   |
//	                                   v
//	                        +---------------------+
//	                        |   Directive Parser  |  <--- Validates against pkg/spec Registry
//	                        +---------------------+
//	                                   |
//	                                   v
//	                        +---------------------+
//	                        |    AST-to-IR Binder |  <--- Resolves generic.Optional, DTOs
//	                        +---------------------+
//	                                   |
//	                                   v
//	                        +---------------------+
//	                        |     *ir.RootIR      |  <--- Unchecked AST Ready for Analysis
//	                        +---------------------+
//
// # Core Building Blocks
//
//   - [Parser]: Main coordinator for file and package AST traversal and IR construction.
//   - [NewParser]: Instantiates a new parser equipped with an isolated token.FileSet.
//   - [ParseDirective]: High-performance directive tokenizer parsing `@directive(arg=val)` strings.
//   - [Directive]: Structured representation of a parsed directive and its argument key-value pairs.
//   - [ParsePathTemplate]: Decomposes an RFC 6570 URI template string into structured segments.
//   - [ParsePipeline]: Parses wire-transform pipeline expressions attached to directives.
//
// # Usage Tiers
//
// ## Tier 1: High-Level Ingestion
//
// Compilers and toolchains parse single files or full directory packages with standard entry points:
//
//	p := parser.NewParser()
//	rootIR, err := p.ParseFile("pkg/api/service.go")
//	if err != nil {
//	    log.Fatalf("parse failed: %v", err)
//	}
//
// ## Tier 2: In-Memory Source Parsing
//
// Language servers, tests, and memory buffers parse Go source bytes directly:
//
//	rootIR, err := p.ParseSource("virtual.go", sourceBytes)
//
// ## Tier 3: Low-Level Directive Extraction
//
// Linters and custom analyzers invoke [ParseDirective] directly to inspect Godoc comments:
//
//	d := parser.ParseDirective("// @service(name=Petstore, engine=aoni)")
//	if d != nil {
//	    _ = d.Name
//	    _ = d.Args["engine"]
//	}
//
// # Concurrency & Thread Safety
//
// An individual [Parser] instance maintains an internal `token.FileSet`. Separate [Parser] instances
// can safely parse distinct files or packages concurrently across goroutines. The resulting [*ir.RootIR]
// is mutable while being built, but should be treated as read-only once passed downstream to analysis.
//
// # Performance & Zero-Allocation Profile
//
// The directive parser scans raw byte slices of comment text directly without invoking regular expressions.
// Token boundaries and directive argument maps are constructed with minimal allocation overhead.
package parser
