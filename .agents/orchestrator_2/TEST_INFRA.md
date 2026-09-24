# E2E Test Infra: Vortex Sovereign Quality Suite

## Test Philosophy
- Requirement-driven testing validating zero-allocation guarantees, terminal presentation standards, documentation completeness, and end-to-end workspace health.
- Methodology: Category-Partition + Boundary Value Analysis + Pairwise Combinatorial Testing + Real-World Workload Testing.

## Feature Inventory & Test Matrix
| # | Feature | Requirement | Tier 1 (Feature) | Tier 2 (Boundary) | Tier 3 (Cross-Feature) | Tier 4 (Workload) |
|---|---------|-------------|:----------------:|:-----------------:|:----------------------:|:-----------------:|
| 1 | `Optional[string]` DTO Serialization | R1, R4 | ✓ (Some, None) | ✓ (`Some("")`, unicode, spaces) | ✓ (combined with int/bool) | ✓ (real DTO payload) |
| 2 | `Optional[int*]` DTO Serialization | R1, R4 | ✓ (int64, int32) | ✓ (0, MinInt, MaxInt) | ✓ (mixed query/form) | ✓ (real DTO payload) |
| 3 | `Optional[uint*]` DTO Serialization | R1, R4 | ✓ (uint64, uint8) | ✓ (0, MaxUint) | ✓ (mixed query/form) | ✓ (real DTO payload) |
| 4 | `Optional[float*]` DTO Serialization | R1, R4 | ✓ (float64, float32) | ✓ (0.0, -0.0, NaN, Inf) | ✓ (mixed query/form) | ✓ (real DTO payload) |
| 5 | `Optional[bool]` DTO Serialization | R1, R4 | ✓ (true, false) | ✓ (`FormatBoolInt` 1/0) | ✓ (mixed query/form) | ✓ (real DTO payload) |
| 6 | `Optional[time.Time]` Serialization | R1, R4 | ✓ (RFC3339) | ✓ (`time.Time{}` zero-time) | ✓ (mixed query/form) | ✓ (real DTO payload) |
| 7 | `Optional[T]` JSON Marshaling | R1 | ✓ (`MarshalJSON`) | ✓ (`Some("")`, `Some(0)`, `None`) | ✓ (struct round-trip) | ✓ (REST JSON DTO) |
| 8 | `Optional[T]` JSON Unmarshaling | R1 | ✓ (`UnmarshalJSON`) | ✓ (`null`, `""`, missing field) | ✓ (struct round-trip) | ✓ (REST JSON DTO) |
| 9 | Raw ANSI Escape Eradication | R2 | ✓ (zero escapes in pkg/internal) | ✓ (edge cases in status/format) | ✓ (colored rendering) | ✓ (CLI command output) |
| 10 | NO_COLOR & Non-TTY Redirection | R2 | ✓ (plain text output) | ✓ (`TERM=dumb`, file redirect) | ✓ (piped execution) | ✓ (full CLI runs) |
| 11 | Sovereign CLI Badges & Tables | R2 | ✓ (tuikit.Table & Badge) | ✓ (wide columns, non-ASCII) | ✓ (diagnostic rendering) | ✓ (`vortex status/check`) |
| 12 | Package `doc.go` Coverage | R3 | ✓ (all pkg/ have doc.go) | ✓ (stubs fully overhauled) | ✓ (cross-package links) | ✓ (godoc buildable) |
| 13 | Sentinel Errors & Predicates | R3 | ✓ (declared Err* variables) | ✓ (typed error structs) | ✓ (`errors.Is`, `errors.AsType`) | ✓ (toolchain diagnostics) |
| 14 | Zero-Alloc Regression Benchmarks | R4 | ✓ (`b.ReportAllocs()`) | ✓ (0 B/op, 0 allocs/op) | ✓ (`AppendQuery` & `AppendFormData`) | ✓ (high-throughput DTO) |

## Test Architecture
- Unit and Benchmark Test Runner: `go test -v -run Test...` and `go test -bench=BenchmarkDTO -benchmem ./pkg/emitter/...`
- Linter Runner: `golangci-lint run ./...`
- Workspace-wide Test Suite: `go test ./...` across all packages in `d:/CodingProjects/vortex` and `d:/CodingProjects/foundation/generic`
- Allocation Assertion: `testing.AllocsPerRun(1000, ...)` asserting 0.0 allocations for primitive optional operations.

## Coverage Thresholds
- Tier 1: ≥5 test cases per feature covering representative inputs
- Tier 2: Boundary value testing (empty string, zero, negative, MaxInt, MinInt, zero time, null)
- Tier 3: Cross-feature combinations (structs mixing multiple Optional types and non-optional types)
- Tier 4: Real-world workload benchmarks proving 0 allocs/op under sustained execution
