# Original User Request

## 2026-09-22T14:19:08Z

Elevate Vortex to sovereign benchmark-grade quality matching aoni and foundation: zero-allocation DTO emitter logic with first-class generic.Optional[T] support, restrained high-craft CLI UX via tuikit (clean Unicode glyphs, zero AI-fluff/emoji spam), comprehensive Godoc package architecture, and rigorous benchmark test suites.

Working directory: d:/CodingProjects/vortex
Integrity mode: development

## Requirements

### R1. Zero-Allocation generic.Optional[T] & Empty Field DTO Codegen
Upgrade pkg/emitter/dto.go and the IR type resolution to eliminate all heap allocations (fmt.Sprint reflection anti-patterns):
- Directly unwrap inner primitive types (string, int, uint, float, bool, time.Time) using zero-alloc encoders (strconv.Append*, direct string escape).
- Natively support empty parameter serialization (field=) when wrapped in generic.Some(""), while omitting entirely when generic.None().
- Ensure interoperability between generic.Optional[T] and url.Values, form-data, and query parameters.
- Add MarshalJSON and UnmarshalJSON support to generic.Optional[T] in foundation/generic/monads.go for seamless JSON serialization.

### R2. Restrained High-Craft CLI Presentation via foundation/tuikit
Modernize cmd/vortex and CLI subcommands using foundation/tuikit with a disciplined, sovereign aesthetic:
- Remove hardcoded ANSI escape sequences scattered across pkg/lint/format.go, pkg/project/status.go, and internal/text/render_terminal.go.
- Replace ad-hoc terminal output with tuikit primitives: tuikit.Box for diagnostics, tuikit.Table for summaries and benchmarking, tuikit.Badge for status tags, and tuikit.VisibleWidth/NO_COLOR safety.
- Strictly avoid emoji spam or "AI-style" decorations; use clean, understated Unicode glyphs (✔, ✖, ◆, ↳, —), dim timestamps, and microsecond/byte statistics ([1.2ms | 0 allocs]).

### R3. Benchmark-Grade Code Documentation & Architecture
Align Vortex package documentation with aoni's forever-frozen core standards:
- Introduce comprehensive doc.go files across all undocumented packages in pkg/ (builder, cache, lint, pipeline, patcher, project, etc.) with usage tiers, ASCII diagrams, and godoc links.
- Standardize typed error predicates and sentinel errors across the toolchain.

### R4. Performance Benchmarks & Adversarial Test Coverage
- Add zero-allocation regression benchmarks (b.ReportAllocs()) for all generated DTO emitters (AppendQuery, AppendFormData, EncodeValues).
- Add comprehensive test suites for generic.Optional[T] fields across varied types and boundary conditions.

## Acceptance Criteria

### DTO Codegen & Zero-Alloc Guarantees
- [ ] DTO emission for generic.Optional[T] produces exactly 0 allocs/op for primitives without calling fmt.Sprint.
- [ ] Explicit empty string generic.Some("") serializes as key= in query and form payloads.
- [ ] Unset optional (generic.None()) appends nothing to the payload.
- [ ] generic.Optional[T] correctly marshals to and unmarshals from JSON matching the inner value or null.

### Terminal Aesthetics & Consistency
- [ ] Zero raw \033[ escapes in pkg/ or internal/; all styling uses foundation/tuikit.
- [ ] Terminal output respects NO_COLOR environment variable and non-TTY redirection.
- [ ] CLI displays clean, restrained status badges and structured tables without emoji spam.

### Codebase Health & Quality
- [ ] Every package in pkg/ contains an idiomatic doc.go following aoni/foundation documentation guidelines.
- [ ] go test ./... passes cleanly across the entire workspace.
- [ ] golangci-lint run reports zero lint violations.

## 2026-09-22T15:52:22Z

Resume execution from current state. Milestone 2 is ready for verification/gate, followed by Milestone 3 (Benchmark-Grade Code Documentation & Architecture) and Milestone 4 (Performance Benchmarks & Adversarial Test Coverage).

## 2026-09-22T19:22:38Z

Server restart completed. Please resume orchestration: finalize Milestone 2 review/gate, then proceed to Milestone 3 (Comprehensive Godoc Architecture across all pkg/ packages) and Milestone 4 (Zero-Alloc Performance Benchmarks).

## 2026-09-23T04:21:26Z

Server restart completed. Please resume orchestration: complete Milestone 3 (Godoc Architecture & Sentinel Errors) review/gate and proceed to Milestone 4 (Zero-Alloc Performance Benchmarks) and final acceptance gate.

## 2026-09-23T12:50:14Z

Server restart completed. Please resume orchestration: finalize Milestone 3 gate, execute Milestone 4 (Zero-Alloc Performance Benchmarks), run the final acceptance gate, and deliver the final report.
