# Milestone 3 Error Architecture Adversarial Challenge Report

## 1. Observation

### 1.1 Empirical Verification Test Suite
The following commands were executed to establish baseline behavior:
1. `go test -v -count=1 ./pkg/project ./pkg/parser ./pkg/diff ./pkg/git ./pkg/cache ./pkg/lint ./pkg/spec ./pkg/pipeline`
   - **Result**: PASS (0 failures across all 8 packages).
2. `$env:GOWORK="off"; go test -count=1 ./...`
   - **Result**: PASS (all 41 packages passed, exit code 0).

### 1.2 Adversarial Stress Testing Results
An empirical test harness was compiled and run against all 8 `errors.go` implementations covering 64 adversarial test cases across four dimensions:
1. **Deep wrapping chains** (5, 10, 50, 100 levels of `fmt.Errorf("wrap_%d: %w", i, err)`):
   - **Result**: 32/32 tests PASSED. Sentinels and structured errors unwrap cleanly through 100 levels without truncation or stack overflow.
2. **Cross-subsystem typed wrapping** (nested typed errors spanning 3 to 8 subsystems, e.g., `PipelineError -> LintError -> ParseError -> ErrSyntaxError`):
   - **Result**: 4/4 tests PASSED. `errors.AsType` successfully extracts nested typed errors and predicates accurately inspect causal sentinels across subsystem boundaries.
3. **Unwrap termination & cycle safety**:
   - **Result**: 2/2 tests PASSED. 9-level heterogeneous error chain terminated in exactly 9 steps reaching `nil`; single-level unwrap stopped at sentinel.
4. **Untyped nil pointer safety & method calls**:
   - `Is<Predicate>(nil)` on untyped nil: 8/8 tests PASSED (returned `false`).
   - `(*SubsystemError)(nil).Error()`: 8/8 tests PASSED (returned `"<nil>"`).
   - `(*SubsystemError)(nil).Unwrap()`: 8/8 tests PASSED (returned `nil`).
   - Non-nil struct with `Err: nil`: 8/8 tests PASSED (returned `false`).
5. **Typed nil pointer safety (CRITICAL FAILURE)**:
   - **Result**: 25/25 tests FAILED with FATAL PANIC:
     `runtime error: invalid memory address or nil pointer dereference`
   - 17/17 direct predicate calls on typed nil pointer (`var p *SubsystemError = nil; Is<Predicate>(p)`) panicked.
   - 8/8 wrapped typed nil pointer calls (`w := fmt.Errorf("wrap: %w", p); Is<Predicate>(w)`) panicked.

### 1.3 Verbatim Panic Trace and Affected Lines
```
=== RUNNING EXHAUSTIVE 17-PREDICATE TYPED NIL STRESS HARNESS ===
FAIL (PANIC): project.IsNotFound(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): project.IsStale(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): parser.IsSyntaxError(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): parser.IsNotFound(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): diff.IsNotFound(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): diff.IsConflict(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): diff.IsEmpty(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): git.IsNotRepository(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): git.IsNotFound(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): cache.IsNotFound(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): cache.IsCorrupt(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): lint.IsLintFailure(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): lint.IsNotFound(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): spec.IsNotFound(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): spec.IsUnsupportedFormat(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): pipeline.IsPipelineAborted(typedNil) -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): pipeline.IsNotFound(typedNil) -> runtime error: invalid memory address or nil pointer dereference

=== WRAPPED TYPED NIL TESTS ===
FAIL (PANIC): WrappedTypedNil_Project -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): WrappedTypedNil_Parser -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): WrappedTypedNil_Diff -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): WrappedTypedNil_Git -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): WrappedTypedNil_Cache -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): WrappedTypedNil_Lint -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): WrappedTypedNil_Spec -> runtime error: invalid memory address or nil pointer dereference
FAIL (PANIC): WrappedTypedNil_Pipeline -> runtime error: invalid memory address or nil pointer dereference
```

Exact source code locations across all 8 packages:
1. `pkg/project/errors.go`:
   - Line 93: `if pErr, ok := errors.AsType[*ProjectError](err); ok { return errors.Is(pErr.Err, ...)`
   - Line 112: `if pErr, ok := errors.AsType[*ProjectError](err); ok { return errors.Is(pErr.Err, ErrStaleCodegen) }`
2. `pkg/parser/errors.go`:
   - Line 84: `if pErr, ok := errors.AsType[*ParseError](err); ok { return errors.Is(pErr.Err, ErrSyntaxError) }`
   - Line 100: `if pErr, ok := errors.AsType[*ParseError](err); ok { return errors.Is(pErr.Err, ErrContractNotFound) }`
3. `pkg/diff/errors.go`:
   - Line 86: `if dErr, ok := errors.AsType[*DiffError](err); ok { return errors.Is(dErr.Err, ErrFrameNotFound) }`
   - Line 102: `if dErr, ok := errors.AsType[*DiffError](err); ok { return errors.Is(dErr.Err, ErrConflict) }`
   - Line 118: `if dErr, ok := errors.AsType[*DiffError](err); ok { return errors.Is(dErr.Err, ErrStackEmpty) }`
4. `pkg/git/errors.go`:
   - Line 83: `if gErr, ok := errors.AsType[*GitError](err); ok { if errors.Is(gErr.Err, ErrNotRepository) ... }`
   - Line 101: `if gErr, ok := errors.AsType[*GitError](err); ok { return errors.Is(gErr.Err, ErrBranchNotFound) }`
5. `pkg/cache/errors.go`:
   - Line 86: `if cErr, ok := errors.AsType[*CacheError](err); ok { return errors.Is(cErr.Err, ...)`
   - Line 104: `if cErr, ok := errors.AsType[*CacheError](err); ok { return errors.Is(cErr.Err, ErrCorruptCache) }`
6. `pkg/lint/errors.go`:
   - Line 83: `if lErr, ok := errors.AsType[*LintError](err); ok { return errors.Is(lErr.Err, ErrLintFailure) }`
   - Line 99: `if lErr, ok := errors.AsType[*LintError](err); ok { return errors.Is(lErr.Err, ErrRuleNotFound) }`
7. `pkg/spec/errors.go`:
   - Line 84: `if sErr, ok := errors.AsType[*SpecError](err); ok { return errors.Is(sErr.Err, ...)`
   - Line 100: `if sErr, ok := errors.AsType[*SpecError](err); ok { return errors.Is(sErr.Err, ErrUnsupportedFormat) }`
8. `pkg/pipeline/errors.go`:
   - Line 83: `if pErr, ok := errors.AsType[*PipelineError](err); ok { return errors.Is(pErr.Err, ErrPipelineAborted) }`
   - Line 99: `if pErr, ok := errors.AsType[*PipelineError](err); ok { return errors.Is(pErr.Err, ErrNoContractsFound) }`

---

## 2. Logic Chain

1. **Premise**: In Go, an `error` interface value containing a typed nil pointer (e.g. `var p *ProjectError = nil; var err error = p`) evaluates `err != nil` to `true` because the dynamic type (`*ProjectError`) is non-nil.
2. **Standard Library Generic Unwrapping**: In Go 1.27, `errors.AsType[T error](err error) (target T, ok bool)` traverses the error chain looking for an error matching type `T`. When `err` holds `(*ProjectError)(nil)`, `errors.AsType[*ProjectError](err)` matches because `err` is indeed of type `*ProjectError`.
3. **Return Values**: `errors.AsType` returns `target = (*ProjectError)(nil)` and `ok = true`.
4. **Execution Flow**:
   - `if err == nil` does not trigger (because dynamic type is set).
   - `if errors.Is(err, ...)` evaluates `p.Unwrap()`, which returns `nil`, resulting in `false`.
   - The predicate enters `if pErr, ok := errors.AsType[*ProjectError](err); ok`.
   - Since `ok == true`, execution proceeds inside the block.
5. **Panic Trigger**: The code attempts to read `pErr.Err`. Because `pErr` is a nil pointer (`nil`), reading `pErr.Err` immediately triggers a hardware fault / runtime panic: `runtime error: invalid memory address or nil pointer dereference`.
6. **False Sense of Safety in Worker Tests**:
   - In `pkg/*/errors_test.go`, the worker asserted `Predicate(nil) == false`. This passed untyped nil `nil`, which was caught by `if err == nil { return false }` at the top of the function.
   - The worker never tested a typed nil pointer (`var p *SubsystemError = nil; Predicate(p)`).
7. **Severity & Blast Radius**: Any caller in Vortex or external consumers invoking a function that returns a typed error pointer, or passing a wrapped nil error to any predicate across any of the 8 subsystems, will crash the entire process.

---

## 3. Caveats

- **Untyped Nil Safety**: The functions are safe when passed an untyped `nil` literal.
- **Deep Wrapping Integrity**: For non-nil errors, deep wrapping up to 100 levels behaves impeccably and passes all checks.
- **Cross-Subsystem Resolution**: For non-nil typed errors, cross-subsystem unwrapping and type matching work as intended.
- **Scope Restriction**: Per Teamwork protocol, the challenger is strictly review-only and does not modify production code in `pkg/`. The remediation must be executed by the worker.

---

## 4. Conclusion & Verdict

**Verdict: REQUEST_CHANGES**

All 17 error predicate functions across all 8 packages (`pkg/project`, `pkg/parser`, `pkg/diff`, `pkg/git`, `pkg/cache`, `pkg/lint`, `pkg/spec`, `pkg/pipeline`) contain an identical critical defect: attempting to dereference `*.Err` after `errors.AsType` without verifying that the returned typed pointer is non-nil (`ok && pErr != nil`).

### Required Action Items for Worker:
1. In all 8 `pkg/*/errors.go` files, add a nil guard to every `errors.AsType` call:
   ```go
   // Change:
   if pErr, ok := errors.AsType[*ProjectError](err); ok {
       return errors.Is(pErr.Err, ...)
   }
   // To:
   if pErr, ok := errors.AsType[*ProjectError](err); ok && pErr != nil {
       return errors.Is(pErr.Err, ...)
   }
   ```
2. In all 8 `pkg/*/errors_test.go` files, add explicit typed nil pointer assertions to every test function:
   ```go
   var typedNil *ProjectError = nil
   if IsNotFound(typedNil) {
       t.Error("expected false for typed nil pointer")
   }
   if IsNotFound(fmt.Errorf("wrap: %w", typedNil)) {
       t.Error("expected false for wrapped typed nil pointer")
   }
   ```

---

## 5. Verification Method

1. Create a test case in any of the 8 packages:
   ```go
   func TestTypedNilPanic(t *testing.T) {
       var p *project.ProjectError = nil
       if project.IsNotFound(p) {
           t.Fatal("expected false")
       }
   }
   ```
2. Run `go test -v -run TestTypedNilPanic ./pkg/project`.
3. Invalidation condition: The test panics with `runtime error: invalid memory address or nil pointer dereference`.
4. Verification of fix: With `ok && pErr != nil`, the test passes with exit code 0.
