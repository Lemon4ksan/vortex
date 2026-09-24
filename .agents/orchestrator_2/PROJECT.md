# Project: Sovereign Upgrade of Vortex

## Architecture
Vortex is a benchmark-grade API toolchain comprising contract parsing, AST manipulation, IR optimization, high-performance code emission, static analysis linting, and CLI tooling.
This project elevates Vortex to match `aoni` and `foundation` sovereign standards across:
1. **Zero-Allocation DTO Emission**: Unwrapping `generic.Optional[T]` primitives with zero heap allocations, full empty string parameter support (`field=`), and native JSON monad marshaling.
2. **Restrained High-Craft CLI Presentation**: Eradicating all 26 raw ANSI escapes, integrating `foundation/tuikit` (Box, Table, Badge, ProbeTerminal, VisibleWidth), eliminating emoji clutter, and guaranteeing NO_COLOR/pipe safety.
3. **Benchmark-Grade Documentation & Error Architecture**: Authoring comprehensive `doc.go` across 21 missing packages + 4 stubs + 2 internal packages, with ASCII diagrams, usage tiers, sentinel errors, and typed error predicates.
4. **Performance Benchmarks & Adversarial Test Suite**: Adding zero-allocation regression benchmarks (`b.ReportAllocs()`), `testing.AllocsPerRun == 0` validation, comprehensive boundary condition tests, and full workspace verification (`go test ./...` and `golangci-lint run`).

## Feature Inventory
| # | Feature | Description | Milestone | Source |
|---|---------|-------------|-----------|--------|
| 1 | F-R1-01 | IR Generic Type Resolution: support `*ast.IndexExpr` and `*ast.IndexListExpr` in `pkg/parser/binder.go` | M1 | survey_dto_spec_1 |
| 2 | F-R1-02 | Monad JSON Serialization: `MarshalJSON`, `UnmarshalJSON`, `IsZero` on `generic.Optional[T]` in `foundation/generic/monads.go` | M1 | survey_dto_spec_1 |
| 3 | F-R1-03 | Zero-Alloc DTO Primitive Emission: unwrap `generic.Optional[T]` with `strconv.Append*` in `pkg/emitter/dto.go` | M1 | survey_dto_spec_1 |
| 4 | F-R1-04 | Empty Field Serialization: `generic.Some("")` -> `wire=`, `generic.None()` -> omitted | M1 | survey_dto_spec_1 |
| 5 | F-R1-05 | Zero-Alloc Query Escaping: stack buffer helper in emitted DTO code eliminating `url.QueryEscape` heap allocations | M1 | survey_dto_spec_1 |
| 6 | F-R1-06 | Zero-Alloc `EncodeValues`: specialize primitive optionals in `EncodeValues` without `fmt.Sprint` | M1 | survey_dto_spec_1 |
| 7 | F-R2-01 | Raw ANSI Escape Eradication: eliminate 26 raw escapes across `render_terminal.go`, `format.go`, `status.go` | M2 | survey_cli_1 |
| 8 | F-R2-02 | Document Renderer Modernization: refactor `internal/text/render_terminal.go` & `intent.go` using `tuikit.Box`, `tuikit.Table`, `tuikit.VisibleWidth` | M2 | survey_cli_1 |
| 9 | F-R2-03 | Linter Diagnostics Modernization: refactor `pkg/lint/format.go` with `tuikit.RenderHeader`, `tuikit.Table`, and `tuikit.IsInteractive` | M2 | survey_cli_1 |
| 10 | F-R2-04 | Workspace Guardian Modernization: rewrite `pkg/project/status.go` using `tuikit.Table`, sovereign badges (`✖ BREAKING`, `▲ DRIFT`, `✔ IN SYNC`), non-TTY safety | M2 | survey_cli_1 |
| 11 | F-R2-05 | Codebase-Wide Emoji Decontamination: remove informal emojis (`⚡`, `✨`, `🔴`, `🟡`, `🔵`, `❌`, `⚠️`, `🚀`) across CLI and subcommands; use restrained Unicode glyphs | M2 | survey_cli_1 |
| 12 | F-R2-06 | Global NO_COLOR & Non-TTY Redirection: probe terminal at `cmd/vortex` boot; guarantee clean plaintext on pipe/file redirection | M2 | survey_cli_1 |
| 13 | F-R2-07 | Telemetry Formatting: use `tuikit.RenderTaxDecomposition`, `tuikit.FormatBytes`, and `[1.2ms \| 0 allocs]` metrics | M2 | survey_cli_1 |
| 14 | F-R3-01 | New Godoc Coverage: author 21 new `doc.go` files in `pkg/` following Aoni frozen-core documentation standards | M3 | survey_arch_bench_1 |
| 15 | F-R3-02 | Stub Godoc Overhaul: expand 4 stub `doc.go` files (`emitter`, `ir`, `optimizer`, `parser`) with ASCII diagrams & tiers | M3 | survey_arch_bench_1 |
| 16 | F-R3-03 | Internal Godoc Coverage: add `doc.go` for `internal/borrow` and `internal/inspector` | M3 | survey_arch_bench_1 |
| 17 | F-R3-04 | Sentinel Errors Architecture: declare package sentinel errors (`var Err... = errors.New(...)`) across key subsystems | M3 | survey_arch_bench_1 |
| 18 | F-R3-05 | Typed Error Structs & Predicates: declare typed error structs and predicates (`IsNotFound`, `IsSyntaxError`, etc.) using Go 1.27 `errors.AsType` | M3 | survey_arch_bench_1 |
| 19 | F-R4-01 | DTO Zero-Alloc Regression Benchmarks: create `b.ReportAllocs()` benchmarks for all generated DTO emitters | M4 | survey_arch_bench_1 |
| 20 | F-R4-02 | Zero-Alloc Assertions: implement unit tests verifying `testing.AllocsPerRun(1000, ...) == 0` for primitive optional serializers | M4 | survey_arch_bench_1 |
| 21 | F-R4-03 | Optional Adversarial Tests: comprehensive tests for `generic.Optional[T]` across all primitive types & boundary conditions | M4 | survey_arch_bench_1 |
| 22 | F-R4-04 | Full Workspace Acceptance Gate: verify 100% clean pass for `go test ./...` and `golangci-lint run ./...` | M4 | survey_arch_bench_1 |

## Milestones
| # | Name | Scope | Dependencies | Status |
|---|------|-------|-------------|--------|
| M1 | Zero-Alloc generic.Optional[T] & DTO Codegen | Features 1–6: `pkg/parser/binder.go`, `pkg/emitter/dto.go`, `foundation/generic/monads.go` | none | DONE (0 allocs verified, tests & lint pass) |
| M2 | Restrained High-Craft CLI via tuikit | Features 7–13: `internal/text/render_terminal.go`, `internal/text/intent.go`, `pkg/lint/format.go`, `pkg/project/status.go`, `cmd/vortex`, emoji cleanup | none | PLANNED |
| M3 | Benchmark-Grade Documentation & Error Architecture | Features 14–18: `doc.go` across 21 pkg + 4 stubs + 2 internal, `errors.go` with sentinels and predicates | none | PLANNED |
| M4 | Performance Benchmarks, Adversarial Tests & E2E Gate | Features 19–22: `pkg/emitter/dto_bench_test.go`, `pkg/emitter/dto_test.go`, full workspace verification | M1, M2, M3 | PLANNED |

## Interface Contracts

### Parser ↔ Emitter Contract
- `pkg/parser/binder.go:extractGoType(expr ast.Expr) ir.GoTypeIR`:
  - For `generic.Optional[T]`: `GoTypeIR.Name` must be `"generic.Optional[T]"`, `GoTypeIR.ElemType` must be `"T"`, `GoTypeIR.IsCustomType` must be `true`.
  - For non-generic identifiers: unchanged.
- `pkg/emitter/dto.go:unwrapOptionalType(f *ir.FieldIR) (innerType string, isOptional bool)`:
  - Inspects `f.Type.Name` and `f.Type.ElemType`. Returns inner type (e.g. `"string"`, `"int64"`, `"bool"`) and `true` when wrapped in `Optional`.

### Foundation Monad Contract
- `foundation/generic/monads.go`:
  - `func (o Optional[T]) MarshalJSON() ([]byte, error)`: returns JSON bytes of `o.val` if valid, or `[]byte("null")` if unset.
  - `func (o *Optional[T]) UnmarshalJSON(data []byte) error`: sets `*o = None[T]()` if `data == "null"` or empty, else unmarshals into `T` and sets `*o = Some(v)`.
  - `func (o Optional[T]) IsZero() bool`: returns `!o.valid` (for Go 1.24+ `omitzero`).

### CLI Tuikit Contract
- `pkg/lint/format.go`, `pkg/project/status.go`, `internal/text/render_terminal.go`:
  - Zero raw `\033[` or `\x1b[` string literals.
  - Colorization only via `tuikit.Bold`, `tuikit.Dim`, `tuikit.Red`, `tuikit.Green`, `tuikit.Yellow`, `tuikit.Cyan`, `tuikit.Gray`, `tuikit.White`.
  - Tabular formatting only via `tuikit.NewTable(...)`.
  - Diagnostics/callouts via `tuikit.NewBox(...)`.
  - Respect `tuikit.IsInteractive(w)` and `tuikit.ColorEnabled()`.

### Error Architecture Contract
- Package sentinel errors: `var Err<Concept> = errors.New("<pkg>: <concept>")`
- Typed error struct: `type <Subsystem>Error struct { Op, Path string, Err error }` with `Error() string` and `Unwrap() error`.
- Predicate functions: `func Is<Predicate>(err error) bool` checking `errors.Is` and `errors.AsType[*<Subsystem>Error]`.

## Code Layout
- `d:/CodingProjects/vortex/pkg/emitter/dto.go` — DTO codegen logic
- `d:/CodingProjects/vortex/pkg/parser/binder.go` — AST to IR binder
- `d:/CodingProjects/foundation/generic/monads.go` — Generic monads (`Optional[T]`)
- `d:/CodingProjects/foundation/generic/monads_test.go` — Monad tests
- `d:/CodingProjects/vortex/internal/text/render_terminal.go` — Document renderer
- `d:/CodingProjects/vortex/internal/text/intent.go` — Intent glyphs
- `d:/CodingProjects/vortex/pkg/lint/format.go` — Lint output formatter
- `d:/CodingProjects/vortex/pkg/project/status.go` — Workspace status renderer
- `d:/CodingProjects/vortex/cmd/vortex/app.go` — CLI boot and probe
- `d:/CodingProjects/vortex/pkg/*/doc.go` — Package Godoc
- `d:/CodingProjects/vortex/pkg/*/errors.go` — Sentinel errors and typed predicates
- `d:/CodingProjects/vortex/pkg/emitter/dto_test.go` — DTO unit & boundary tests
- `d:/CodingProjects/vortex/pkg/emitter/dto_bench_test.go` — DTO zero-alloc regression benchmarks
