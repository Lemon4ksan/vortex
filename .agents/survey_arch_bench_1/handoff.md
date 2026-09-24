# Handoff Report: R3 & R4 Architectural & Benchmark Survey

- **Agent**: `teamwork_preview_explorer` (`survey_arch_bench_1`)
- **Target**: R3 & R4 (Documentation, Error Architecture, Benchmarks & Test Baseline) of the Sovereign Upgrade of Vortex
- **Workspace**: `d:/CodingProjects/vortex`
- **Timestamp**: 2026-09-22T14:35:00Z

---

## 1. Observation

Direct, verbatim findings gathered via tool executions across the codebase:

### 1.1 Package & `doc.go` Inventory
Execution of `go list ./...` revealed 41 total Go packages across the workspace.
- **Root/CLI**:
  - `ast`: `doc.go` exists (108 lines, comprehensive with ASCII diagram & usage tiers).
  - `cmd/vortex`: entrypoint CLI application; no package `doc.go`.
- **Packages in `pkg/` (27 total)**:
  - **Existing `doc.go` (7 packages)**:
    - `pkg/analysis/doc.go` (31 lines; well-structured architecture overview & rule hierarchy).
    - `pkg/asyncapi/doc.go` (49 lines; covers dual-version normalization, trait merging, set algebra).
    - `pkg/emitter/doc.go` (**10 lines stub**; minimal 2-sentence description, lacks building blocks, diagrams, tiers).
    - `pkg/ir/doc.go` (**14 lines stub**; brief 4-scope list, lacks builders and type references).
    - `pkg/openapi/doc.go` (53 lines; covers 3-way merging, set algebra, supported RFC standards).
    - `pkg/optimizer/doc.go` (**20 lines**; lists 4 optimization passes, lacks diagrams and pipeline integration).
    - `pkg/parser/doc.go` (**10 lines stub**; minimal 2-sentence description).
  - **Missing `doc.go` (21 packages)**:
    - `pkg/builder` (programmatic IR builder & contract assembly)
    - `pkg/cache` (lint cache, encrypted secrets vault, compressed traffic store)
    - `pkg/cfg` (control flow graph & static dataflow analysis)
    - `pkg/diff` (semantic AST diffing & checkpoint stack snapshots)
    - `pkg/enum` (OpenAPI/TypeScript enum inference & validation)
    - `pkg/git` (git repository integration & commit inspection)
    - `pkg/history` (contract version evolution & ledger)
    - `pkg/ingest` (HAR 1.2 network capture ingestion)
    - `pkg/jsbundle` (JavaScript/TypeScript SDK bundling)
    - `pkg/lint` (static analysis contract lint engine)
    - `pkg/merge` (3-way spec reconciliation & conflict resolution)
    - `pkg/mirror` (upstream API contract synchronization & drift detection)
    - `pkg/oracle/gen` (JS AST & runtime oracle emitter)
    - `pkg/oracle/spec` (test harness validation specifications)
    - `pkg/patcher` (in-place Go AST patching & transformation)
    - `pkg/pipeline` (end-to-end compilation & deobfuscation pipeline)
    - `pkg/project` (workspace manifest `.vortex.yml`, drift & status engine)
    - `pkg/spec` (multi-spec manager & remote fetcher)
    - `pkg/sys` (system environment & platform capability inspection)
    - `pkg/tuple` (positional tuple inference & JSON array unwrapping)
    - `pkg/version` (version metadata & build flags)
- **Packages in `internal/` (11 total)**:
  - Existing `doc.go` (9 packages): `internal/ast`, `internal/base`, `internal/core` (stub, 8 lines), `internal/oracle` (stub, 8 lines), `internal/perf`, `internal/spec`, `internal/text`, `internal/traffic`, `internal/workspace`.
  - Missing `doc.go` (2 packages):
    - `internal/borrow` (linear resource & AST borrow checker)
    - `internal/inspector` (embedded web dashboard & live traffic inspector)

### 1.2 Error Handling & Architecture Observations
- **Sentinel Errors**: Zero declared package sentinel errors (`var Err... = errors.New(...)`) exist anywhere in `pkg/` or `internal/`.
- **Custom Error Types**: Zero custom error types implementing `Error() string` exist across the entire repository.
- **Error Creation**: 140+ instances of ad-hoc `fmt.Errorf(...)` and 50+ instances of `errors.New(...)` are scattered across packages (e.g. `errors.New("target file is required")` repeated verbatim 5 times in `pkg/pipeline/ast.go:54, 63, 158, 240, 476`).
- **Call-Site Inspection**: Callers in `cmd/vortex` and tests rely on fragile substring checks (e.g. `require.Contains(t, err.Error(), "unknown command")` at `cmd/vortex/app_test.go:360`).
- **Foundation / Aoni Precedent**: Inspection of `D:/CodingProjects/aoni/errors.go` and `D:/CodingProjects/foundation/net/dns/error.go` revealed:
  - Exported sentinel errors (`var ErrNotFound = ...`, `var ErrNXDomain = ...`).
  - Typed domain error structs with `Unwrap() error` (`APIError`, `ResolutionError`, `BridgeError`).
  - Standard typed predicates using Go 1.27 `errors.AsType[T](err)` and `errors.Is(err, ...)` (`IsNotFound(err)`, `IsTimeout(err)`, `IsRateLimited(err)`).

### 1.3 Benchmark and Test Suite Observations
- **Total Test Files**: 42 `*_test.go` files exist.
- **Existing Benchmarks**: Only 6 benchmarks exist in the entire codebase:
  - `internal/text/bench_test.go`: `BenchmarkDocumentBuilder_Build`, `BenchmarkDocumentBuilder_RenderMarkdown`, `BenchmarkShield_ProtectAndRestore`.
  - `pkg/emitter/bitpack_test.go`: `BenchmarkPacketHeader_PackUint64`, `BenchmarkPacketHeader_PackTo`, `BenchmarkPacketHeader_BatchUnpack` (executed via subprocess `go test` runner).
  - `pkg/emitter/harness.go`: contains code generator for `Benchmark_<svcIdentifier>_<method>` using `for b.Loop()` and `b.ReportAllocs()`.
- **Untested Packages** (no test files):
  - `pkg/enum`, `pkg/ir`, `pkg/oracle/spec`, `pkg/pipeline`, `pkg/version`.
  - `internal/ast`, `internal/base`, `internal/borrow`, `internal/core`, `internal/oracle`, `internal/spec`, `internal/traffic`, `internal/workspace`.
- **Emitter DTO Testing Gap**:
  - `pkg/emitter/dto.go` generates `AppendFormData`, `AppendQuery`, and `EncodeValues`.
  - `pkg/emitter/dto_test.go` **does not exist**.
  - `pkg/emitter/emitter_test.go` only verifies substring occurrence in output Go strings (`require.Contains(t, codeStr, "func (r *User) EncodeValues(...)")`).
  - No compilation, execution, or allocation benchmarks exist for generated DTO serializers.

### 1.4 Baseline Build, Test, and Lint Status
- **`go test ./...`**:
  - Result: **PASS** (Exit code 0).
  - All 41 packages compile and pass existing test suites cleanly (cached/<1s for libraries, ~4.3s for `cmd/vortex`).
- **`golangci-lint run ./...`**:
  - Result: **PASS** (Exit code 0, 0 issues reported).
  - Validated with `.golangci.yml` v2 enabled linters (`govet`, `errcheck`, `staticcheck`, `revive`, `errorlint`, `gosec`, `gocritic`, `prealloc`, `perfsprint`, etc.).

### 1.5 DTO Emitter Zero-Allocation Deficiencies
Inspection of `pkg/emitter/dto.go` (lines 45-168 and 170-246) revealed:
1. **Reflection & Heap Boxing Anti-Pattern**:
   Lines 137-146 and lines 221-228:
   ```go
   if strings.HasPrefix(f.Type.Name, "generic.Optional[") || strings.HasPrefix(f.Type.Name, "Optional[") {
       tracker.Add("fmt")
       fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
       buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
       fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
       fmt.Fprintf(buf, "\t\tdst = append(dst, url.QueryEscape(fmt.Sprint(optVal))...)\n")
       buf.WriteString("\t}\n")
       return
   }
   ```
   - `fmt.Sprint(optVal)` allocates heap memory for every invocation by boxing `optVal` into `any` and performing runtime reflection.
   - `url.QueryEscape(...)` allocates a fresh string on the heap for every parameter.
2. **Missing Inner-Type Primitive Specialization**:
   The generator does not unwrap inner primitive types (`string`, `int`, `uint`, `float64`, `bool`, `time.Time`) of `generic.Optional[T]`.
3. **Empty String vs None Violation**:
   `generic.Some("")` must serialize as `field=`, while `generic.None()` must be omitted. Currently, `fmt.Sprint("")` returns `""`, which serializes `field=`, but untyped handling means `generic.None()` causes omission only through `Value()` check without formal contract test coverage.
4. **`time.Time` Heap Allocations**:
   Line 77: `url.QueryEscape(r.%s.Format(time.RFC3339))` performs two separate heap string allocations instead of zero-alloc buffer appending (`optVal.AppendFormat(dst, time.RFC3339)`).
5. **`foundation/generic/monads.go` Gap**:
   `Optional[T]` in `foundation/generic/monads.go` lacks `MarshalJSON` and `UnmarshalJSON` implementations.

---

## 2. Logic Chain

From the direct observations above, we establish the following deductive reasoning chain:

1. **Documentation Standard Gap**:
   - *Observation*: Only 7/27 packages in `pkg/` have `doc.go`, and 4 of those (`emitter`, `ir`, `parser`, `optimizer`) are 10-20 line stubs.
   - *Requirement*: aoni/foundation manifesto demands comprehensive Godoc package architecture with package overview, architecture/pipeline stages, ASCII diagrams, usage tiers (Tier 1 beginner, Tier 2 standard, Tier 3 advanced/performance), and Godoc type links `[Type]`.
   - *Deduction*: 21 packages in `pkg/` must receive new `doc.go` files, and 4 stub packages must be upgraded to full sovereign grade. In `internal/`, 2 packages (`borrow`, `inspector`) must receive new `doc.go` files.

2. **Error Architecture Gap**:
   - *Observation*: The toolchain currently has 0 sentinel errors, 0 typed error structs, and 190+ ad-hoc `fmt.Errorf` / `errors.New` strings.
   - *Requirement*: Standardize typed error predicates and sentinel errors across the toolchain.
   - *Deduction*: Subsystems must define typed error structs (e.g. `ParseError`, `ConfigError`, `DiffError`, `GitError`, `CacheError`, `LintError`) implementing `Error() string` and `Unwrap() error`, package-level sentinel errors (`ErrNotFound`, `ErrSyntaxError`, `ErrInvalidConfig`, etc.), and package-level predicates (`IsNotFound`, `IsSyntaxError`, `IsValidationError`, `IsConflict`) utilizing Go 1.27 `errors.AsType` and `errors.Is`.

3. **Performance & Zero-Alloc Emitter Gap**:
   - *Observation*: `pkg/emitter/dto.go` emits `fmt.Sprint(optVal)` and `url.QueryEscape(...)` for all `generic.Optional[T]` fields.
   - *Requirement*: DTO emission for `generic.Optional[T]` must produce exactly 0 allocs/op for primitives without calling `fmt.Sprint`.
   - *Deduction*: `pkg/emitter/dto.go` must be updated to parse the inner type `T` of `generic.Optional[T]` from `ir.FieldIR` and emit dedicated zero-allocation encoding paths using `strconv.AppendInt`, `strconv.AppendUint`, `strconv.AppendFloat`, `strconv.AppendBool`, direct query escaping to byte buffer, and explicit `field=` serialization for `generic.Some("")`.

4. **Benchmark & Test Harness Gap**:
   - *Observation*: Zero benchmarks exist for DTO emitters (`AppendQuery`, `AppendFormData`, `EncodeValues`), and no unit test compiles or executes generated DTO serialization code.
   - *Requirement*: Add zero-allocation regression benchmarks (`b.ReportAllocs()`) for all generated DTO emitters.
   - *Deduction*: A new test and benchmark suite (`pkg/emitter/dto_test.go` or generated test harness) must be constructed that emits DTO structs with varied primitive and optional fields, compiles them, executes unit assertions (including `testing.AllocsPerRun(1000, ...) == 0`), and runs regression benchmarks with `b.ReportAllocs()` and `for b.Loop()`.

---

## 3. Caveats

1. **External Dependency Modification**:
   - `foundation/generic/monads.go` belongs to `github.com/lemon4ksan/foundation`, located at `D:/CodingProjects/foundation`. Adding `MarshalJSON` and `UnmarshalJSON` to `generic.Optional[T]` requires editing `D:/CodingProjects/foundation/generic/monads.go` or submitting a PR/local replace in `go.mod`.
2. **`url.Values` Inherent Allocations**:
   - `url.Values` is defined by stdlib as `map[string][]string`. The `EncodeValues(vals url.Values)` method internally allocates map headers and slice storage when `vals.Set` or `vals.Add` is called. While `AppendQuery` and `AppendFormData` can achieve true 0 allocs/op with a stack-backed destination slice, `EncodeValues` cannot be 100% zero-alloc unless `vals` already holds allocated slice buckets. Zero-alloc guarantees strictly apply to `AppendQuery` and `AppendFormData`.
3. **Query Escaping Zero-Alloc Requirement**:
   - Standard library `url.QueryEscape(s)` allocates a new string. For `AppendQuery` to achieve 0 allocs on string fields containing characters that require escaping, the emitter must emit or link to an in-place/buffer-appending escaping helper (e.g. `appendEscapedQuery(dst, s)` or `urlkit.AppendEscape`).

---

## 4. Conclusion & Feature Specifications

### 4.1 Complete Package `doc.go` Action List

#### Group A: 21 Packages in `pkg/` Needing Brand New `doc.go`
1. `pkg/builder`: Programmatic builder and contract assembly pipeline for IR and client components.
2. `pkg/cache`: Workspace caching, AST lint rule result memoization, encrypted secrets vault, and compressed traffic session storage.
3. `pkg/cfg`: Control Flow Graph (CFG) analysis, static dataflow inspection, reachability analysis of API endpoints.
4. `pkg/diff`: Semantic contract diffing, 3-way structural diffing, checkpoint stack snapshots, and schema migration tracking.
5. `pkg/enum`: TypeScript/OpenAPI/Go enum inference, string-to-int enum mappings, type-safe enum constants and validators.
6. `pkg/git`: Git repository integration, branch/commit tracking, worktree state inspection, and change detection for contract versioning.
7. `pkg/history`: Contract evolution timeline, migration history ledger, semantic version bump analysis.
8. `pkg/ingest`: Network capture ingestion (HAR 1.2 archives), extracting HTTP request/response pairs into contract candidates.
9. `pkg/jsbundle`: JavaScript/TypeScript client SDK bundling and distribution packaging.
10. `pkg/lint`: Static analysis lint engine enforcing aoni contract conventions, REST guidelines, RFC 9110/RFC 6570 rules.
11. `pkg/merge`: 3-way specification merging, conflict resolution, and schema reconciliation.
12. `pkg/mirror`: Upstream API mirror synchronization, contract mirroring, and remote drift detection.
13. `pkg/oracle/gen`: JavaScript/Node.js AST and runtime oracle code emission for browser/protocol test harnesses.
14. `pkg/oracle/spec`: Specification oracle definitions, test harness validation specs.
15. `pkg/patcher`: In-place AST transformation and contract patching (applying non-destructive code updates).
16. `pkg/pipeline`: End-to-end compilation pipeline (Parse -> Normalize -> Analyze -> Optimize -> Emit).
17. `pkg/project`: `.vortex.yml` workspace manifest parsing, contract lifecycle management, drift detection, and project status.
18. `pkg/spec`: Unified specification management (OpenAPI, AsyncAPI, HAR), remote fetching, and contract registry.
19. `pkg/sys`: System environment detection, platform architecture inspection, file system isolation.
20. `pkg/tuple`: Positional tuple schema inference, JSON array to typed struct/tuple conversion.
21. `pkg/version`: Vortex toolchain version metadata, semantic build tags, and runtime capability flags.

#### Group B: 4 Packages in `pkg/` Needing Comprehensive Overhaul
1. `pkg/emitter`: Expand 10-line stub to cover DTO encoders, client structs, sub-requester routing, harness generation, and tuple decoders with ASCII flow diagrams.
2. `pkg/ir`: Expand 14-line stub to cover the 4 hierarchical scopes (Service, Method, Param, Return), type resolution rules, and AST-to-IR invariants.
3. `pkg/optimizer`: Expand 20 lines to detail sub-requester clustering, stack allocation sizing, query canonicalization, and CPU cache-line alignment with memory layout diagrams.
4. `pkg/parser`: Expand 10-line stub to detail directive parsing (@service, @get, @post, @unwrap), AST traversal, and syntax diagnostics.

#### Group C: 2 Internal Packages Needing `doc.go`
1. `internal/borrow`: Linear resource tracking, borrow-checker semantics for mutable request modifiers.
2. `internal/inspector`: Real-time HTTP capture dashboard, WebSocket telemetry streamer, and HAR replay engine.

---

### 4.2 Error Architecture Standardization Specification

Each subsystem will follow the uniform Aoni/Foundation pattern:

```go
// 1. Sentinel Errors
var (
    ErrNotFound      = errors.New("<pkg>: resource not found")
    ErrInvalidConfig = errors.New("<pkg>: invalid configuration")
    // ...
)

// 2. Typed Error Struct
type <Subsystem>Error struct {
    Op      string
    Path    string
    Message string
    Err     error
}

func (e *<Subsystem>Error) Error() string { ... }
func (e *<Subsystem>Error) Unwrap() error { return e.Err }

// 3. Typed Error Predicates
func IsNotFound(err error) bool {
    if errors.Is(err, ErrNotFound) {
        return true
    }
    if e, ok := errors.AsType[*<Subsystem>Error](err); ok {
        return errors.Is(e.Err, ErrNotFound)
    }
    return false
}
```

#### Subsystem Mapping Matrix:
| Subsystem Package | Key Sentinel Errors | Typed Error Struct | Typed Predicates |
|---|---|---|---|
| `pkg/project` | `ErrWorkspaceNotFound`, `ErrConfigNotFound`, `ErrInvalidConfig`, `ErrContractNotFound`, `ErrStaleCodegen` | `ProjectError` (`Op`, `Path`, `Key`, `Err`) | `IsNotFound(err)`, `IsStale(err)` |
| `pkg/parser` | `ErrSyntaxError`, `ErrContractNotFound`, `ErrInvalidDirective`, `ErrUnresolvedType` | `ParseError` (`File`, `Line`, `Col`, `Err`) | `IsSyntaxError(err)`, `IsNotFound(err)` |
| `pkg/diff` | `ErrStackEmpty`, `ErrFrameNotFound`, `ErrInsufficientFrames`, `ErrConflict` | `DiffError` (`Op`, `Frame`, `Err`) | `IsNotFound(err)`, `IsConflict(err)` |
| `pkg/git` | `ErrNotRepository`, `ErrBranchNotFound`, `ErrGitCommandFailed` | `GitError` (`Op`, `Target`, `Stderr`, `Err`) | `IsNotRepository(err)` |
| `pkg/cache` | `ErrSessionNotFound`, `ErrSecretNotFound`, `ErrCorruptCache` | `CacheError` (`Op`, `Key`, `Err`) | `IsNotFound(err)` |
| `pkg/lint` | `ErrLintFailure`, `ErrRuleNotFound`, `ErrFixFailed` | `LintError` (`RuleID`, `Target`, `Line`, `Err`) | `IsLintFailure(err)` |
| `pkg/spec` | `ErrSpecNotFound`, `ErrUnsupportedFormat`, `ErrEmptySpec` | `SpecError` (`Source`, `Format`, `Err`) | `IsNotFound(err)` |
| `pkg/pipeline` | `ErrTargetFileRequired`, `ErrNoContractsFound`, `ErrPipelineAborted` | `PipelineError` (`Stage`, `Target`, `Err`) | `IsPipelineAborted(err)` |

---

### 4.3 Benchmark & Zero-Alloc Regression Suite Inventory & Requirements

#### 1. DTO Codegen Engine Upgrades (`pkg/emitter/dto.go`)
- **Inner Type Extraction**: Extract inner type `T` from `generic.Optional[T]`.
- **Primitive Handling**:
  - `string`: If `generic.Some("")`, append `wireName=`. If `generic.Some("val")`, append `wireName=val` with zero-alloc byte escaping. If `generic.None()`, omit.
  - `int`, `int8`, `int16`, `int32`, `int64`: `strconv.AppendInt(dst, int64(optVal), 10)`.
  - `uint`, `uint8`, `uint16`, `uint32`, `uint64`: `strconv.AppendUint(dst, uint64(optVal), 10)`.
  - `float32`, `float64`: `strconv.AppendFloat(dst, float64(optVal), 'f', -1, 64)`.
  - `bool`: append `wireName=true` or `wireName=false` or `wireName=1`.
  - `time.Time`: `optVal.AppendFormat(dst, time.RFC3339)` (zero allocations).
- **Escape Helper**: Emit local helper `appendQueryEscape(dst []byte, s string) []byte` into emitted code, avoiding stdlib `url.QueryEscape` string heap allocation.

#### 2. Regression Benchmark Suite Specifications
Implement in `pkg/emitter/dto_bench_test.go` and `pkg/emitter/dto_test.go`:
```go
func BenchmarkDTO_AppendQuery_Primitives(b *testing.B)
func BenchmarkDTO_AppendQuery_Optional_Some(b *testing.B)
func BenchmarkDTO_AppendQuery_Optional_None(b *testing.B)
func BenchmarkDTO_AppendQuery_Optional_EmptyString(b *testing.B)
func BenchmarkDTO_AppendFormData_AllTypes(b *testing.B)
func BenchmarkDTO_EncodeValues(b *testing.B)
```
Each benchmark:
- Pre-allocates destination buffer: `buf := make([]byte, 0, 512)`
- Calls `b.ReportAllocs()`
- Uses `for b.Loop() { dst = dto.AppendQuery(buf[:0]) }`
- Mandatory unit test assertion using `testing.AllocsPerRun(1000, ...)`:
  ```go
  allocs := testing.AllocsPerRun(1000, func() {
      dst = dto.AppendQuery(buf[:0])
  })
  require.Equal(t, float64(0), allocs, "AppendQuery must produce exactly 0 allocs/op")
  ```

---

### 4.4 Enumerated Feature Inventory for Sovereign Upgrade

| Feature ID | Scope | Description | Priority |
|---|---|---|---|
| **FEAT-R3-01** | Godoc Coverage | Author 21 new `doc.go` files in `pkg/` following Aoni frozen core format (Summary, Architecture, Building Blocks, Quick Start, Tiers, Godoc links). | High |
| **FEAT-R3-02** | Godoc Expansion | Overhaul 4 stub `doc.go` files (`emitter`, `ir`, `optimizer`, `parser`) with ASCII diagrams and technical depth. | Medium |
| **FEAT-R3-03** | Internal Godoc | Add `doc.go` for `internal/borrow` and `internal/inspector`. | Low |
| **FEAT-R3-04** | Error Architecture | Create `errors.go` in key packages (`project`, `parser`, `diff`, `git`, `cache`, `lint`, `spec`, `pipeline`) declaring sentinel errors and typed error structs. | High |
| **FEAT-R3-05** | Error Predicates | Implement typed predicates (`IsNotFound`, `IsSyntaxError`, `IsValidationError`, `IsConflict`) using Go 1.27 `errors.AsType`. | High |
| **FEAT-R4-01** | Zero-Alloc DTO Codegen | Refactor `pkg/emitter/dto.go` to unwrap `generic.Optional[T]` primitives with `strconv.Append*` and eliminate `fmt.Sprint`. | Critical |
| **FEAT-R4-02** | Empty Parameter Support | Implement serialization of `generic.Some("")` as `wire_key=` and total omission of `generic.None()`. | Critical |
| **FEAT-R4-03** | Zero-Alloc Query Escape | Introduce stack/byte-buffer escaping helper to eliminate `url.QueryEscape` heap allocations in emitted DTOs. | High |
| **FEAT-R4-04** | Foundation JSON Monad | Add `MarshalJSON` and `UnmarshalJSON` to `generic.Optional[T]` in `foundation/generic/monads.go`. | High |
| **FEAT-R4-05** | DTO Benchmark Suite | Add zero-allocation regression benchmarks (`b.ReportAllocs()`) and `testing.AllocsPerRun` assertions for DTO serializers. | High |
| **FEAT-R4-06** | Optional Adversarial Tests | Add comprehensive unit tests covering `generic.Optional[T]` boundary conditions across all primitive types. | High |

---

## 5. Verification Method

To independently verify all findings and baseline assertions in this report:

1. **Verify Baseline Test Suite**:
   ```bash
   go test ./...
   ```
   *Expected outcome*: Exit code 0, all packages pass cleanly.

2. **Verify Baseline Linter**:
   ```bash
   golangci-lint run ./...
   ```
   *Expected outcome*: Exit code 0, 0 issues reported.

3. **Verify Missing `doc.go` Files**:
   Run in PowerShell:
   ```powershell
   go list -f '{{.Dir}}' ./pkg/... | ForEach-Object {
       $doc = Join-Path $_ "doc.go"
       if (-not (Test-Path $doc)) { Write-Output "Missing: $_" }
   }
   ```
   *Expected outcome*: 21 directories in `pkg/` listed as missing `doc.go`.

4. **Verify Absence of Sentinel Errors**:
   ```powershell
   grep -rn "var Err" pkg/ internal/
   ```
   *Expected outcome*: Zero matches.

5. **Verify Reflection Anti-Pattern in `dto.go`**:
   Inspect `d:/CodingProjects/vortex/pkg/emitter/dto.go` at lines 138-143 and 222-225 to observe `fmt.Sprint(optVal)`.
