# DISPATCH — explorer_m4_3

**Role**: Milestone 4 Synthesis & Worker Execution Planner
**Working Directory**: `d:/CodingProjects/vortex/.agents/explorer_m4_3/`
**Workspace Root**: `d:/CodingProjects/vortex`

## Mandatory Reading
1. `d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md`
2. `d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md`
3. `pkg/emitter/dto.go`
4. `pkg/emitter/dto_bench_test.go`
5. `pkg/emitter/dto_test.go`

## Task Description
Synthesize the complete worker plan for Milestone 4 (Performance Benchmarks & Adversarial Test Coverage):
1. **Benchmark Suite Blueprint**:
   - Structure of `pkg/emitter/dto_bench_test.go`:
     - Benchmarks with `b.ReportAllocs()` covering:
       - `BenchmarkAppendQuery_Primitives` (all fields populated)
       - `BenchmarkAppendQuery_Optionals_Some`
       - `BenchmarkAppendQuery_Optionals_None`
       - `BenchmarkAppendFormData_Primitives`
       - `BenchmarkAppendFormData_Optionals`
       - `BenchmarkEncodeValues_ZeroAlloc`
   - Verification command: `go test -benchmem -run=^$ -bench=BenchmarkAppend.* ./pkg/emitter`.
2. **Zero-Alloc Unit Test Blueprint**:
   - In `pkg/emitter/dto_test.go`:
     - `testing.AllocsPerRun(1000, func() { ... }) == 0` assertions for primitive optional serialization paths.
     - Verification command: `go test -v ./pkg/emitter -run TestZeroAlloc`.
3. **Adversarial Edge-Case Suite Blueprint**:
   - `Some("")` -> `field=`
   - `None()` -> omitted
   - Full gamut of primitive types (`string`, `int`, `int64`, `uint`, `float64`, `bool`).
   - JSON monad serialization edge cases.
4. **Verification Gates & Standards**:
   - Full workspace test pass: `$env:GOWORK="off"; go test -count=1 ./...`
   - Zero lint violations: `golangci-lint run --allow-parallel-runners ./...`
5. **Worker Task Breakdown**:
   - Clear, step-by-step checklist for `worker_m4` with exact file targets.

Write your report to `d:/CodingProjects/vortex/.agents/explorer_m4_3/handoff.md`.
Do NOT modify production files (Explorer is read-only).
When done, notify parent via send_message.
