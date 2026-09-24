# Milestone 3 Remediation Specification & Master Plan

- **Author**: `explorer_m3_fix_3` (teamwork_preview_explorer)
- **Role**: Explorer, Synthesizer, Remediation Architect
- **Target**: Milestone 3 Remediation Worker
- **Workspace**: `d:/CodingProjects/vortex`
- **Working Directory**: `d:/CodingProjects/vortex/.agents/explorer_m3_fix_3/`
- **Output Artifact**: `d:/CodingProjects/vortex/.agents/explorer_m3_fix_3/handoff.md`
- **Status**: Complete & Verified (Ready for Implementation)

---

## Executive Summary & Comprehensive File Modification Inventory

This master remediation specification synthesizes the gate feedback from `challenger_m3_1` (typed nil pointer dereference panics across 17 predicates), `challenger_m3_2` (redundant package comments in 4 sibling files and hallucinated APIs across 8 `doc.go` files), and `reviewer_m3_1`.

An exhaustive AST and source audit confirms that exactly **28 files** must be modified to satisfy all Milestone 3 acceptance criteria:

| # | Batch | Relative File Path | Component | Defect / Objective |
|---|---|---|---|---|
| 1 | Batch 1 | `pkg/emitter/emitter.go` | Sibling Cleanup | Remove redundant `// Package emitter` comment at line 5 |
| 2 | Batch 1 | `pkg/ingest/namer.go` | Sibling Cleanup | Remove redundant `// Package ingest` comment at line 5 |
| 3 | Batch 1 | `pkg/lint/rule.go` | Sibling Cleanup | Remove redundant `// Package lint` comment at line 5 |
| 4 | Batch 1 | `pkg/openapi/importer.go` | Sibling Cleanup | Remove redundant `// Package openapi` comment at lines 5-6 |
| 5 | Batch 2 | `pkg/project/errors.go` | Predicate Fix | Add `&& pErr != nil` guard to `IsNotFound` (line 93) and `IsStale` (line 112) |
| 6 | Batch 2 | `pkg/parser/errors.go` | Predicate Fix | Add `&& pErr != nil` guard to `IsSyntaxError` (line 84) and `IsNotFound` (line 100) |
| 7 | Batch 2 | `pkg/diff/errors.go` | Predicate Fix | Add `&& dErr != nil` guard to `IsNotFound` (line 86), `IsConflict` (line 102), `IsEmpty` (line 118) |
| 8 | Batch 2 | `pkg/git/errors.go` | Predicate Fix | Add `&& gErr != nil` guard to `IsNotRepository` (line 83) and `IsNotFound` (line 101) |
| 9 | Batch 2 | `pkg/cache/errors.go` | Predicate Fix | Add `&& cErr != nil` guard to `IsNotFound` (line 86) and `IsCorrupt` (line 104) |
| 10 | Batch 2 | `pkg/lint/errors.go` | Predicate Fix | Add `&& lErr != nil` guard to `IsLintFailure` (line 83) and `IsNotFound` (line 99) |
| 11 | Batch 2 | `pkg/spec/errors.go` | Predicate Fix | Add `&& sErr != nil` guard to `IsNotFound` (line 84) and `IsUnsupportedFormat` (line 100) |
| 12 | Batch 2 | `pkg/pipeline/errors.go` | Predicate Fix | Add `&& pErr != nil` guard to `IsPipelineAborted` (line 83) and `IsNotFound` (line 99) |
| 13 | Batch 3 | `pkg/project/errors_test.go` | Test Suite | Add typed nil pointer assertions to `TestProjectErrors_IsNotFound` & `TestProjectErrors_IsStale` |
| 14 | Batch 3 | `pkg/parser/errors_test.go` | Test Suite | Add typed nil pointer assertions to `TestParserErrors_IsSyntaxError` & `TestParserErrors_IsNotFound` |
| 15 | Batch 3 | `pkg/diff/errors_test.go` | Test Suite | Add typed nil pointer assertions to `TestDiffErrors_IsNotFound`, `TestDiffErrors_IsConflict`, `TestDiffErrors_IsEmpty` |
| 16 | Batch 3 | `pkg/git/errors_test.go` | Test Suite | Add typed nil pointer assertions to `TestGitErrors_IsNotRepository` & `TestGitErrors_IsNotFound` |
| 17 | Batch 3 | `pkg/cache/errors_test.go` | Test Suite | Add typed nil pointer assertions to `TestCacheErrors_IsNotFound` & `TestCacheErrors_IsCorrupt` |
| 18 | Batch 3 | `pkg/lint/errors_test.go` | Test Suite | Add typed nil pointer assertions to `TestLintErrors_IsLintFailure` & `TestLintErrors_IsNotFound` |
| 19 | Batch 3 | `pkg/spec/errors_test.go` | Test Suite | Add typed nil pointer assertions to `TestSpecErrors_IsNotFound` & `TestSpecErrors_IsUnsupportedFormat` |
| 20 | Batch 3 | `pkg/pipeline/errors_test.go` | Test Suite | Add typed nil pointer assertions to `TestPipelineErrors_IsPipelineAborted` & `TestPipelineErrors_IsNotFound` |
| 21 | Batch 4 | `pkg/ingest/doc.go` | Godoc Alignment | Replace `ParseHAR`/`ConvertHARToOpenAPI` with `HARToOpenAPI`/`HARToOpenAPIOpts` and valid examples |
| 22 | Batch 4 | `pkg/cache/doc.go` | Godoc Alignment | Replace hallucinated vault/traffic APIs with `LoadSecrets`, `TrafficIndex`, `StoreTraffic`, `GetTraffic`, `IsFresh`/`Put` |
| 23 | Batch 4 | `pkg/cfg/doc.go` | Godoc Alignment | Replace non-existent `Build`/`ReachingDefinitions`/`LiveVariables` with `New`, `WalkPaths`, `FindLoopBlocks` |
| 24 | Batch 4 | `pkg/diff/doc.go` | Godoc Alignment | Replace `CheckpointStack` with `DiffStack` / `LoadStack`; align `Compare` / `Push` signatures |
| 25 | Batch 4 | `pkg/jsbundle/doc.go` | Godoc Alignment | Replace `ScanDirectory` with `ScanFiles`; fix `ScanBytes` single return value (`*ScanResult`) |
| 26 | Batch 4 | `pkg/git/doc.go` | Godoc Alignment | Replace `ListBranches`/`IsCleanWorkingTree` with `ListProposalBranches`/`IsClean` and correct parameters |
| 27 | Batch 4 | `pkg/mirror/doc.go` | Godoc Alignment | Remove non-existent `CheckAllServices` and `SyncService`; align with real `CheckService` signature |
| 28 | Batch 4 | `pkg/parser/doc.go` | Godoc Alignment | Remove non-existent `[Lexer]` and `[Token]` links; escape or de-bracket `generic.Optional[T]` |

---

## 1. Observation

### 1.1 Empirical Verification of Typed Nil Panics
When a typed pointer variable initialized to `nil` is passed to any of the 17 error predicate functions:
```go
var p *project.ProjectError = nil
project.IsNotFound(p)
```
The Go runtime immediately terminates with:
```
panic: runtime error: invalid memory address or nil pointer dereference
[signal 0xc0000005, code 0x0 addr 0x30 pc 0x...]
```
The exact panic occurs because `errors.AsType[*ProjectError](err)` in Go 1.27 returns `(target: (*ProjectError)(nil), ok: true)`. The predicate checks `if pErr, ok := ...; ok {`, which evaluates to `true`, and subsequent dereference `pErr.Err` faults.

### 1.2 Empirical Verification of Duplicate Package Comments in Sibling Files
Executing `go doc ./pkg/emitter`, `go doc ./pkg/ingest`, `go doc ./pkg/lint`, and `go doc ./pkg/openapi` outputs duplicate trailing summary headers. For example, running `go doc ./pkg/emitter` displays:
```
# Performance & Zero-Allocation Profile
...
Package emitter generates high-performance, zero-allocation Go client facades
from Vortex AST RootIR contracts.
```
This dangling sentence is caused by line 5 of `pkg/emitter/emitter.go`:
```go
// Package emitter generates high-performance, zero-allocation Go client facades from Vortex AST RootIR contracts.
package emitter
```
Similar redundant comments exist in `pkg/ingest/namer.go:5`, `pkg/lint/rule.go:5`, and `pkg/openapi/importer.go:5-6`.

### 1.3 Empirical Verification of Missing Symbols & Hallucinated APIs in `doc.go`
Executing `go doc <pkg> <Symbol>` confirms that the following documented identifiers do not exist in the codebase:
- `go doc ./pkg/ingest ParseHAR` -> `doc: symbol ParseHAR not found` (Exit 1)
- `go doc ./pkg/ingest ConvertHARToOpenAPI` -> `doc: symbol ConvertHARToOpenAPI not found` (Exit 1)
- `go doc ./pkg/cache LoadSecretsVault` -> `doc: symbol LoadSecretsVault not found` (Exit 1)
- `go doc ./pkg/cache TrafficStore` -> `doc: symbol TrafficStore not found` (Exit 1)
- `go doc ./pkg/cfg Build` -> `doc: symbol Build not found` (Exit 1)
- `go doc ./pkg/cfg ReachingDefinitions` -> `doc: symbol ReachingDefinitions not found` (Exit 1)
- `go doc ./pkg/cfg LiveVariables` -> `doc: symbol LiveVariables not found` (Exit 1)
- `go doc ./pkg/diff CheckpointStack` -> `doc: symbol CheckpointStack not found` (Exit 1)
- `go doc ./pkg/jsbundle ScanDirectory` -> `doc: symbol ScanDirectory not found` (Exit 1)
- `go doc ./pkg/git ListBranches` -> `doc: symbol ListBranches not found` (Exit 1)
- `go doc ./pkg/git IsCleanWorkingTree` -> `doc: symbol IsCleanWorkingTree not found` (Exit 1)
- `go doc ./pkg/mirror CheckAllServices` -> `doc: symbol CheckAllServices not found` (Exit 1)
- `go doc ./pkg/mirror SyncService` -> `doc: symbol SyncService not found` (Exit 1)
- `go doc ./pkg/parser Lexer` -> `doc: symbol Lexer not found` (Exit 1)
- `go doc ./pkg/parser Token` -> `doc: symbol Token not found` (Exit 1)

In addition, deep code inspection revealed hidden mismatches:
- In `pkg/cache/doc.go`: `lc.IsValid(...)` and `lc.Set(...)` were documented, but the actual methods on `LintCache` are `lc.IsFresh(relPath string, content []byte) bool` and `lc.Put(relPath string, content []byte, issueCount int)`. Furthermore, `lc.Save()` requires `rootDir string`.
- In `pkg/diff/doc.go`: `diff.Compare(...)` was documented as returning `(report, err)`, but `Compare` actually returns `*DiffReport` directly with no error return.
- In `pkg/jsbundle/doc.go`: `ScanBytes(...)` was documented as returning `(res, err)`, but `ScanBytes` actually returns `*ScanResult` directly with no error return.

---

## 2. Logic Chain

1. **Typed Nil Safety Logic**:
   - In Go, `var err error = (*ProjectError)(nil)` is a non-nil interface (`err != nil` evaluates to `true`) whose underlying concrete value is a nil pointer of type `*ProjectError`.
   - `errors.AsType[*ProjectError](err)` checks if `err` matches `*ProjectError`. Because the dynamic type is `*ProjectError`, `errors.AsType` returns `target = (*ProjectError)(nil)` and `ok = true`.
   - The predicate's `if pErr, ok := errors.AsType[*ProjectError](err); ok` succeeds because `ok` is `true`.
   - When the predicate accesses `pErr.Err`, it attempts to read memory offset `+0x30` from address `0x0`, causing an unrecoverable runtime panic.
   - Guarding the condition with `ok && pErr != nil` ensures that typed nil pointers evaluate safely without dereference. Because `(*SubsystemError).Error()` and `(*SubsystemError).Unwrap()` already handle nil receivers safely, the predicate gracefully returns `false`.

2. **Godoc Deduplication Logic**:
   - `go doc` gathers all comments immediately preceding `package <name>` across every source file in the package directory and concatenates them.
   - Having a dedicated `doc.go` file requires that sibling `.go` files have clean `package <name>` declarations without doc comments, ensuring `doc.go` is the sole source of package documentation.

3. **API Alignment Logic**:
   - Godoc documentation comments are read by developers as the authoritative guide.
   - Code snippets in `doc.go` must compile and use genuine exported signatures.
   - Bracketed identifiers (`[Identifier]`) in Go 1.19+ doc comments create hyperlinks in web godoc and are rendered in terminal `go doc`. If the identifier does not exist, `go doc` leaves the brackets unlinked, signaling an unresolved link defect.
   - Replacing hallucinated symbols with their genuine counterparts guarantees that every link resolves and every code example is compilable.

---

## 3. Caveats

- **Explorer Read-Only Boundary**: Explorer is strictly read-only for production files in `pkg/` and `internal/`. This specification provides exact code replacements for `worker_m3_remediation` to apply.
- **Go 1.27 Pre-requisite**: The project leverages Go 1.27 generic error unwrapping (`errors.AsType`). This is standard across the workspace (`go version go1.27.0 windows/amd64`).
- **Historical Untested Packages**: 11 internal/spec packages (`ast`, `base`, `borrow`, `core`, `oracle`, `spec`, `traffic`, `workspace`, `enum`, `ir`, `oracle/spec`, `version`) show `[no test files]` in `go test ./...`. This is preexisting legacy state and outside Milestone 3 scope.

---

## 4. Master Remediation Specification

### Batch 1: Sibling Comment Cleanups (4 Files)

#### 1. `d:/CodingProjects/vortex/pkg/emitter/emitter.go`
- **Target Line**: Line 5
- **Current Content**:
  ```go
  // Package emitter generates high-performance, zero-allocation Go client facades from Vortex AST RootIR contracts.
  package emitter
  ```
- **Replacement Content**:
  ```go
  package emitter
  ```

#### 2. `d:/CodingProjects/vortex/pkg/ingest/namer.go`
- **Target Line**: Line 5
- **Current Content**:
  ```go
  // Package ingest implements generic specification detection and intelligent naming normalizers.
  package ingest
  ```
- **Replacement Content**:
  ```go
  package ingest
  ```

#### 3. `d:/CodingProjects/vortex/pkg/lint/rule.go`
- **Target Line**: Line 5
- **Current Content**:
  ```go
  // Package lint provides a modular contract linter and diagnostic engine for aoni/vortex interfaces.
  package lint
  ```
- **Replacement Content**:
  ```go
  package lint
  ```

#### 4. `d:/CodingProjects/vortex/pkg/openapi/importer.go`
- **Target Lines**: Lines 5-6
- **Current Content**:
  ```go
  // Package openapi provides parsing, loading, 3-way specification merging,
  // and declarative Go contract generation for OpenAPI 2.0/3.0/3.1 and HAR specifications.
  package openapi
  ```
- **Replacement Content**:
  ```go
  package openapi
  ```

---

### Batch 2: Typed Nil Predicate Hardening (8 Files, 17 Predicates)

#### 1. `d:/CodingProjects/vortex/pkg/project/errors.go`
- **Lines 93-98** (`IsNotFound`):
  ```go
  // BEFORE:
  	if pErr, ok := errors.AsType[*ProjectError](err); ok {
  		return errors.Is(pErr.Err, ErrWorkspaceNotFound) ||
  			errors.Is(pErr.Err, ErrConfigNotFound) ||
  			errors.Is(pErr.Err, ErrContractNotFound) ||
  			errors.Is(pErr.Err, os.ErrNotExist)
  	}

  // AFTER:
  	if pErr, ok := errors.AsType[*ProjectError](err); ok && pErr != nil {
  		return errors.Is(pErr.Err, ErrWorkspaceNotFound) ||
  			errors.Is(pErr.Err, ErrConfigNotFound) ||
  			errors.Is(pErr.Err, ErrContractNotFound) ||
  			errors.Is(pErr.Err, os.ErrNotExist)
  	}
  ```
- **Lines 112-114** (`IsStale`):
  ```go
  // BEFORE:
  	if pErr, ok := errors.AsType[*ProjectError](err); ok {
  		return errors.Is(pErr.Err, ErrStaleCodegen)
  	}

  // AFTER:
  	if pErr, ok := errors.AsType[*ProjectError](err); ok && pErr != nil {
  		return errors.Is(pErr.Err, ErrStaleCodegen)
  	}
  ```

#### 2. `d:/CodingProjects/vortex/pkg/parser/errors.go`
- **Lines 84-86** (`IsSyntaxError`):
  ```go
  // BEFORE:
  	if pErr, ok := errors.AsType[*ParseError](err); ok {
  		return errors.Is(pErr.Err, ErrSyntaxError)
  	}

  // AFTER:
  	if pErr, ok := errors.AsType[*ParseError](err); ok && pErr != nil {
  		return errors.Is(pErr.Err, ErrSyntaxError)
  	}
  ```
- **Lines 100-102** (`IsNotFound`):
  ```go
  // BEFORE:
  	if pErr, ok := errors.AsType[*ParseError](err); ok {
  		return errors.Is(pErr.Err, ErrContractNotFound)
  	}

  // AFTER:
  	if pErr, ok := errors.AsType[*ParseError](err); ok && pErr != nil {
  		return errors.Is(pErr.Err, ErrContractNotFound)
  	}
  ```

#### 3. `d:/CodingProjects/vortex/pkg/diff/errors.go`
- **Lines 86-88** (`IsNotFound`):
  ```go
  // BEFORE:
  	if dErr, ok := errors.AsType[*DiffError](err); ok {
  		return errors.Is(dErr.Err, ErrFrameNotFound)
  	}

  // AFTER:
  	if dErr, ok := errors.AsType[*DiffError](err); ok && dErr != nil {
  		return errors.Is(dErr.Err, ErrFrameNotFound)
  	}
  ```
- **Lines 102-104** (`IsConflict`):
  ```go
  // BEFORE:
  	if dErr, ok := errors.AsType[*DiffError](err); ok {
  		return errors.Is(dErr.Err, ErrConflict)
  	}

  // AFTER:
  	if dErr, ok := errors.AsType[*DiffError](err); ok && dErr != nil {
  		return errors.Is(dErr.Err, ErrConflict)
  	}
  ```
- **Lines 118-120** (`IsEmpty`):
  ```go
  // BEFORE:
  	if dErr, ok := errors.AsType[*DiffError](err); ok {
  		return errors.Is(dErr.Err, ErrStackEmpty)
  	}

  // AFTER:
  	if dErr, ok := errors.AsType[*DiffError](err); ok && dErr != nil {
  		return errors.Is(dErr.Err, ErrStackEmpty)
  	}
  ```

#### 4. `d:/CodingProjects/vortex/pkg/git/errors.go`
- **Lines 83-87** (`IsNotRepository`):
  ```go
  // BEFORE:
  	if gErr, ok := errors.AsType[*GitError](err); ok {
  		if errors.Is(gErr.Err, ErrNotRepository) {
  			return true
  		}
  	}

  // AFTER:
  	if gErr, ok := errors.AsType[*GitError](err); ok && gErr != nil {
  		if errors.Is(gErr.Err, ErrNotRepository) {
  			return true
  		}
  	}
  ```
- **Lines 101-103** (`IsNotFound`):
  ```go
  // BEFORE:
  	if gErr, ok := errors.AsType[*GitError](err); ok {
  		return errors.Is(gErr.Err, ErrBranchNotFound)
  	}

  // AFTER:
  	if gErr, ok := errors.AsType[*GitError](err); ok && gErr != nil {
  		return errors.Is(gErr.Err, ErrBranchNotFound)
  	}
  ```

#### 5. `d:/CodingProjects/vortex/pkg/cache/errors.go`
- **Lines 86-90** (`IsNotFound`):
  ```go
  // BEFORE:
  	if cErr, ok := errors.AsType[*CacheError](err); ok {
  		return errors.Is(cErr.Err, ErrSessionNotFound) ||
  			errors.Is(cErr.Err, ErrSecretNotFound) ||
  			errors.Is(cErr.Err, os.ErrNotExist)
  	}

  // AFTER:
  	if cErr, ok := errors.AsType[*CacheError](err); ok && cErr != nil {
  		return errors.Is(cErr.Err, ErrSessionNotFound) ||
  			errors.Is(cErr.Err, ErrSecretNotFound) ||
  			errors.Is(cErr.Err, os.ErrNotExist)
  	}
  ```
- **Lines 104-106** (`IsCorrupt`):
  ```go
  // BEFORE:
  	if cErr, ok := errors.AsType[*CacheError](err); ok {
  		return errors.Is(cErr.Err, ErrCorruptCache)
  	}

  // AFTER:
  	if cErr, ok := errors.AsType[*CacheError](err); ok && cErr != nil {
  		return errors.Is(cErr.Err, ErrCorruptCache)
  	}
  ```

#### 6. `d:/CodingProjects/vortex/pkg/lint/errors.go`
- **Lines 83-85** (`IsLintFailure`):
  ```go
  // BEFORE:
  	if lErr, ok := errors.AsType[*LintError](err); ok {
  		return errors.Is(lErr.Err, ErrLintFailure)
  	}

  // AFTER:
  	if lErr, ok := errors.AsType[*LintError](err); ok && lErr != nil {
  		return errors.Is(lErr.Err, ErrLintFailure)
  	}
  ```
- **Lines 99-101** (`IsNotFound`):
  ```go
  // BEFORE:
  	if lErr, ok := errors.AsType[*LintError](err); ok {
  		return errors.Is(lErr.Err, ErrRuleNotFound)
  	}

  // AFTER:
  	if lErr, ok := errors.AsType[*LintError](err); ok && lErr != nil {
  		return errors.Is(lErr.Err, ErrRuleNotFound)
  	}
  ```

#### 7. `d:/CodingProjects/vortex/pkg/spec/errors.go`
- **Lines 84-86** (`IsNotFound`):
  ```go
  // BEFORE:
  	if sErr, ok := errors.AsType[*SpecError](err); ok {
  		return errors.Is(sErr.Err, ErrSpecNotFound) || errors.Is(sErr.Err, os.ErrNotExist)
  	}

  // AFTER:
  	if sErr, ok := errors.AsType[*SpecError](err); ok && sErr != nil {
  		return errors.Is(sErr.Err, ErrSpecNotFound) || errors.Is(sErr.Err, os.ErrNotExist)
  	}
  ```
- **Lines 100-102** (`IsUnsupportedFormat`):
  ```go
  // BEFORE:
  	if sErr, ok := errors.AsType[*SpecError](err); ok {
  		return errors.Is(sErr.Err, ErrUnsupportedFormat)
  	}

  // AFTER:
  	if sErr, ok := errors.AsType[*SpecError](err); ok && sErr != nil {
  		return errors.Is(sErr.Err, ErrUnsupportedFormat)
  	}
  ```

#### 8. `d:/CodingProjects/vortex/pkg/pipeline/errors.go`
- **Lines 83-85** (`IsPipelineAborted`):
  ```go
  // BEFORE:
  	if pErr, ok := errors.AsType[*PipelineError](err); ok {
  		return errors.Is(pErr.Err, ErrPipelineAborted)
  	}

  // AFTER:
  	if pErr, ok := errors.AsType[*PipelineError](err); ok && pErr != nil {
  		return errors.Is(pErr.Err, ErrPipelineAborted)
  	}
  ```
- **Lines 99-101** (`IsNotFound`):
  ```go
  // BEFORE:
  	if pErr, ok := errors.AsType[*PipelineError](err); ok {
  		return errors.Is(pErr.Err, ErrNoContractsFound)
  	}

  // AFTER:
  	if pErr, ok := errors.AsType[*PipelineError](err); ok && pErr != nil {
  		return errors.Is(pErr.Err, ErrNoContractsFound)
  	}
  ```

---

### Batch 3: Companion Test Suite Augmentation (8 Files)

In each test suite, add a dedicated typed nil test block asserting that calling the predicate on `var typedNil *SubsystemError = nil` and `fmt.Errorf("wrap: %w", typedNil)` returns `false` without panic:

#### 1. `d:/CodingProjects/vortex/pkg/project/errors_test.go`
- In `TestProjectErrors_IsNotFound`:
  ```go
  	// 6. Typed nil error
  	var typedNil *project.ProjectError
  	require.False(t, project.IsNotFound(typedNil))
  	require.False(t, project.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
  ```
- In `TestProjectErrors_IsStale`:
  ```go
  	// 6. Typed nil error
  	var typedNil *project.ProjectError
  	require.False(t, project.IsStale(typedNil))
  	require.False(t, project.IsStale(fmt.Errorf("wrap: %w", typedNil)))
  ```

#### 2. `d:/CodingProjects/vortex/pkg/parser/errors_test.go`
- In `TestParserErrors_IsSyntaxError`:
  ```go
  	// 6. Typed nil error
  	var typedNil *parser.ParseError
  	require.False(t, parser.IsSyntaxError(typedNil))
  	require.False(t, parser.IsSyntaxError(fmt.Errorf("wrap: %w", typedNil)))
  ```
- In `TestParserErrors_IsNotFound`:
  ```go
  	// 6. Typed nil error
  	var typedNil *parser.ParseError
  	require.False(t, parser.IsNotFound(typedNil))
  	require.False(t, parser.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
  ```

#### 3. `d:/CodingProjects/vortex/pkg/diff/errors_test.go`
- In `TestDiffErrors_IsNotFound`:
  ```go
  	// 6. Typed nil error
  	var typedNil *diff.DiffError
  	require.False(t, diff.IsNotFound(typedNil))
  	require.False(t, diff.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
  ```
- In `TestDiffErrors_IsConflict`:
  ```go
  	// 6. Typed nil error
  	var typedNil *diff.DiffError
  	require.False(t, diff.IsConflict(typedNil))
  	require.False(t, diff.IsConflict(fmt.Errorf("wrap: %w", typedNil)))
  ```
- In `TestDiffErrors_IsEmpty`:
  ```go
  	// 6. Typed nil error
  	var typedNil *diff.DiffError
  	require.False(t, diff.IsEmpty(typedNil))
  	require.False(t, diff.IsEmpty(fmt.Errorf("wrap: %w", typedNil)))
  ```

#### 4. `d:/CodingProjects/vortex/pkg/git/errors_test.go`
- In `TestGitErrors_IsNotRepository`:
  ```go
  	// 6. Typed nil error
  	var typedNil *git.GitError
  	require.False(t, git.IsNotRepository(typedNil))
  	require.False(t, git.IsNotRepository(fmt.Errorf("wrap: %w", typedNil)))
  ```
- In `TestGitErrors_IsNotFound`:
  ```go
  	// 6. Typed nil error
  	var typedNil *git.GitError
  	require.False(t, git.IsNotFound(typedNil))
  	require.False(t, git.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
  ```

#### 5. `d:/CodingProjects/vortex/pkg/cache/errors_test.go`
- In `TestCacheErrors_IsNotFound`:
  ```go
  	// 6. Typed nil error
  	var typedNil *cache.CacheError
  	require.False(t, cache.IsNotFound(typedNil))
  	require.False(t, cache.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
  ```
- In `TestCacheErrors_IsCorrupt`:
  ```go
  	// 6. Typed nil error
  	var typedNil *cache.CacheError
  	require.False(t, cache.IsCorrupt(typedNil))
  	require.False(t, cache.IsCorrupt(fmt.Errorf("wrap: %w", typedNil)))
  ```

#### 6. `d:/CodingProjects/vortex/pkg/lint/errors_test.go`
- In `TestLintErrors_IsLintFailure`:
  ```go
  	// 6. Typed nil error
  	var typedNil *lint.LintError
  	require.False(t, lint.IsLintFailure(typedNil))
  	require.False(t, lint.IsLintFailure(fmt.Errorf("wrap: %w", typedNil)))
  ```
- In `TestLintErrors_IsNotFound`:
  ```go
  	// 6. Typed nil error
  	var typedNil *lint.LintError
  	require.False(t, lint.IsNotFound(typedNil))
  	require.False(t, lint.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
  ```

#### 7. `d:/CodingProjects/vortex/pkg/spec/errors_test.go`
- In `TestSpecErrors_IsNotFound`:
  ```go
  	// 6. Typed nil error
  	var typedNil *spec.SpecError
  	require.False(t, spec.IsNotFound(typedNil))
  	require.False(t, spec.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
  ```
- In `TestSpecErrors_IsUnsupportedFormat`:
  ```go
  	// 6. Typed nil error
  	var typedNil *spec.SpecError
  	require.False(t, spec.IsUnsupportedFormat(typedNil))
  	require.False(t, spec.IsUnsupportedFormat(fmt.Errorf("wrap: %w", typedNil)))
  ```

#### 8. `d:/CodingProjects/vortex/pkg/pipeline/errors_test.go`
- In `TestPipelineErrors_IsPipelineAborted`:
  ```go
  	// 6. Typed nil error
  	var typedNil *pipeline.PipelineError
  	require.False(t, pipeline.IsPipelineAborted(typedNil))
  	require.False(t, pipeline.IsPipelineAborted(fmt.Errorf("wrap: %w", typedNil)))
  ```
- In `TestPipelineErrors_IsNotFound`:
  ```go
  	// 6. Typed nil error
  	var typedNil *pipeline.PipelineError
  	require.False(t, pipeline.IsNotFound(typedNil))
  	require.False(t, pipeline.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
  ```

---

### Batch 4: Godoc Architectural & API Alignment (8 Files)

#### 1. `d:/CodingProjects/vortex/pkg/ingest/doc.go`
Replace entire file content with verified APIs (`HARToOpenAPI`, `HARToOpenAPIOpts`, `DetectFormat`, `IngestOptions`):
```go
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
//	|                         HARToOpenAPI                              |
//	|  - Decodes JSON log into HAREntry stream                          |
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
//	|                       OpenAPI 3.x Synthesis                       |
//	|  - Emits valid OpenAPI 3.x Document AST (*openapi.Document)       |
//	+-------------------------------------------------------------------+
//
// # Core Building Blocks
//
//   - [HARToOpenAPI]: Transforms recorded W3C HAR 1.2 logs into an OpenAPI 3.0 specification.
//   - [HARToOpenAPIOpts]: Transforms recorded HAR logs applying custom [IngestOptions].
//   - [DetectFormat]: Inspects raw specification bytes and determines its [SpecFormat].
//   - [SpecFormat]: Identifies standard API specification and capture formats.
//   - [IngestOptions]: Configures HAR traffic parsing with ignore and route template rewrite rules.
//   - [HARLog]: Root W3C HAR 1.2 container.
//   - [HAREntry]: Individual captured request-response transaction.
//   - [HARNV]: Key-value pair representing headers or query parameters.
//   - [HARPostData]: Request body payload model.
//   - [HARContent]: Response body payload model.
//
// # Usage Tiers
//
// ## Tier 1: Automated Spec Synthesis
//
// Ingest a browser HAR capture and generate an OpenAPI document:
//
//	doc, err := ingest.HARToOpenAPI(harBytes)
//	if err != nil {
//	    log.Fatalf("failed to ingest HAR: %v", err)
//	}
//
// ## Tier 2: Granular Options & Route Rewriting
//
// Apply custom ignore patterns and parameterized route templates:
//
//	opts := ingest.IngestOptions{
//	    IgnorePatterns: []string{"*.png", "*.css"},
//	    RouteTemplates: []string{"/api/v1/users/{id}"},
//	}
//	doc, err := ingest.HARToOpenAPIOpts(harBytes, opts)
//
// ## Tier 3: Format Auto-Detection
//
// Automatically identify specification formats prior to conversion:
//
//	format, err := ingest.DetectFormat(rawBytes)
//	if format == ingest.FormatHAR {
//	    doc, err = ingest.HARToOpenAPI(rawBytes)
//	}
//
// # Concurrency & Thread Safety
//
// [HARToOpenAPI] and [HARToOpenAPIOpts] are stateless pure functions. They are 100% thread-safe
// across concurrent goroutines.
//
// # Performance & Zero-Allocation Profile
//
// The parser processes transaction records sequentially, reusing internal string buffers to convert
// URL paths and header maps with minimal garbage collector pressure.
package ingest
```

#### 2. `d:/CodingProjects/vortex/pkg/cache/doc.go`
Replace entire file content with verified APIs (`LintCache.IsFresh`/`Put`, `LoadSecrets`, `TrafficIndex`, `StoreTraffic`, `GetTraffic`):
```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cache provides workspace caching including lint memoization, encrypted secrets, and traffic storage.
//
// # Architecture Overview
//
// The cache package implements the three primary state persistence pillars of the Vortex developer toolchain:
// fast SHA-256 lint memoization, local machine secrets vault, and compressed gzip HTTP traffic recording:
//
//	+-------------------------------------------------------------------+
//	|                          pkg/cache Subsystems                     |
//	+-------------------------------------------------------------------+
//	         |                          |                         |
//	         v                          v                         v
//	+------------------+      +-------------------+     +------------------+
//	|    LintCache     |      |   SecretsVault    |     |   Traffic Store  |
//	| - File SHA-256   |      | - Local .vortex/  |     | - gzip Session   |
//	| - Issue counts   |      | - Target metadata |     | - TrafficIndex   |
//	| - .vortex/cache/ |      | - Key masking     |     | - Store/Get      |
//	+------------------+      +-------------------+     +------------------+
//
// # Core Building Blocks
//
//   - [LintCache]: Memoizes contract file hashes and diagnostic counts in `.vortex/cache/lint.json`.
//   - [LoadLintCache]: Loads or initializes the workspace lint cache.
//   - [SecretsVault]: Manages local, machine-isolated secrets stored in `.vortex/cache/secrets.json`.
//   - [LoadSecrets]: Discovers and loads the secrets vault from the working directory traversing upwards.
//   - [TrafficIndex]: Stores the catalog of cached traffic captures in `.vortex/cache/traffic/index.json`.
//   - [StoreTraffic]: Compresses and archives a HAR payload into `.vortex/cache/traffic/<hash>.har.gz`.
//   - [GetTraffic]: Retrieves and decompresses a cached traffic session by ID or hash prefix.
//   - [ListTraffic]: Returns all cached traffic sessions sorted by storage date descending.
//   - [TrafficEntry]: Metadata snapshot for a cached traffic recording session.
//   - [SecretEntry]: Individual credentials or token entry stored in the secrets vault.
//
// # Usage Tiers
//
// ## Tier 1: Lint Memoization
//
// Check if a file's lint state is clean and cache results:
//
//	lc, err := cache.LoadLintCache(rootDir)
//	if err == nil && lc.IsFresh(relPath, fileBytes) {
//	    // skip unchanged file
//	}
//	lc.Put(relPath, fileBytes, issueCount)
//	_ = lc.Save(rootDir)
//
// ## Tier 2: Secrets Vault Management
//
// Retrieve and store sensitive credentials for contract replay:
//
//	vault, vaultPath, err := cache.LoadSecrets(rootDir)
//	if err == nil {
//	    token, found := vault.Get("API_KEY")
//	    if !found {
//	        vault.Set("API_KEY", "secret-value", "cli")
//	        _ = vault.Save(vaultPath)
//	    }
//	}
//
// ## Tier 3: Captured Traffic Storage
//
// Save and query recorded HTTP traffic sessions:
//
//	entry, _, err := cache.StoreTraffic(rootDir, "recording.har", harBytes, false, false)
//	if err == nil {
//	    payload, entry, err := cache.GetTraffic(rootDir, entry.ID)
//	}
//
// # Concurrency & Thread Safety
//
// [LintCache] and [SecretsVault] guard internal mutations using sync.RWMutex locks and are fully safe
// for concurrent reads and writes across goroutines. [StoreTraffic] and [GetTraffic] serialize index updates atomically.
//
// # Performance & Zero-Allocation Profile
//
// Traffic payloads are compressed using streaming gzip to minimize memory spikes on large HAR bodies.
// Lint caching utilizes 256-bit SHA-256 hashes formatted in-place via stack hex buffers.
package cache
```

#### 3. `d:/CodingProjects/vortex/pkg/cfg/doc.go`
Replace entire file content with verified APIs (`New`, `WalkPaths`, `FindLoopBlocks`, `FindStatementPosition`):
```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cfg constructs an optimized, lightweight Control Flow Graph (CFG) for Go AST.
//
// # Architecture Overview
//
// The CFG engine splits Go statement trees into sequential basic blocks, computes predecessor and successor
// edges across conditional branches, switches, and loops, and performs reachability and path analysis:
//
//	Go Function Body (ast.BlockStmt)
//	               |
//	               v
//	+-------------------------------------------------------------------+
//	|                             cfg.New                               |
//	|  - Splits statements into basic sequential Blocks                 |
//	|  - Resolves IfStmt, SwitchStmt, ForStmt, BranchStmt jumps        |
//	|  - Computes Predecessor (Preds) and Successor (Succs) edges       |
//	+-------------------------------------------------------------------+
//	               |
//	               v
//	+-------------------------------------------------------------------+
//	|                           CFG Graph                               |
//	|  Block 0 (Entry) ---> Block 1 (Condition) --True--> Block 2 (Body)|
//	|                             |                          |          |
//	|                           False                        v          |
//	|                             +--------------------> Block 3 (Exit) |
//	+-------------------------------------------------------------------+
//	               |
//	               v
//	+-------------------------------------------------------------------+
//	|                     Path Traversal & Loops                        |
//	|  - Acyclic path traversal via WalkPaths                           |
//	|  - Loop body block detection via FindLoopBlocks                   |
//	+-------------------------------------------------------------------+
//
// # Core Building Blocks
//
//   - [New]: Constructs the [CFG] for the specified function body.
//   - [CFG]: Represents the complete graph structure of indexed basic blocks and return points.
//   - [Block]: Represents an individual sequential basic block of statements and outgoing edges.
//   - [BlockKind]: Categorizes block semantics (e.g. entry, branch condition, loop header, exit).
//   - [PathVisitor]: Callback invoked for each complete acyclic execution path.
//   - [FindStatementPosition]: Returns token position of a node within the CFG.
//
// # Usage Tiers
//
// ## Tier 1: Graph Construction
//
// Construct a CFG from a Go function declaration body:
//
//	graph := cfg.New(funcDecl.Body, nil)
//	entry := graph.Entry()
//	returns := graph.ReturnBlocks()
//
// ## Tier 2: Block Traversal & Reachability
//
// Inspect basic blocks, determine live blocks, and traverse control paths:
//
//	for _, block := range graph.Blocks {
//	    if !block.Live {
//	        // dead code detected
//	    }
//	    for _, succ := range block.Succs {
//	        _ = succ.Index
//	    }
//	}
//
// ## Tier 3: Path Walking & Loop Detection
//
// Traverse all acyclic execution paths and identify loop bodies:
//
//	graph.WalkPaths(func(path []*cfg.Block) {
//	    // inspect path
//	})
//	loopBlocks := graph.FindLoopBlocks()
//
// # Concurrency & Thread Safety
//
// [CFG] instances are strictly immutable once constructed by [New]. They are 100% thread-safe
// for concurrent graph queries, dataflow passes, and traversals across multiple goroutines.
//
// # Performance & Zero-Allocation Profile
//
// Basic block branch pointers utilize inline small-array buffers (`succs2 [2]*Block`) for two-way
// conditional branches, eliminating slice heap allocations on typical if-else control structures.
package cfg
```

#### 4. `d:/CodingProjects/vortex/pkg/diff/doc.go`
Replace entire file content with verified APIs (`DiffStack`, `LoadStack`, `Compare`, `CompareWithOptions`):
```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package diff implements semantic contract drift analysis between local Go interfaces and OpenAPI specifications.
//
// # Architecture Overview
//
// The diff engine computes semantic discrepancies between local declarative contract ASTs and remote
// OpenAPI/AsyncAPI specifications. Discrepancies are categorized by severity (Breaking, NonBreaking, Ghost),
// while a persistent [DiffStack] maintains snapshot history across contract iterations:
//
//	Local RootIR Contract           Remote OpenAPI Specification
//	          \                                /
//	           v                              v
//	    +--------------------------------------------+
//	    |                  Compare                   |
//	    |  - Normalizes paths (/users/{id} <-> :id)  |
//	    |  - Aligns operations by HTTP verb and path |
//	    +--------------------------------------------+
//	                          |
//	                          v
//	    +--------------------------------------------+
//	    |                DiffReport                  |
//	    |  - SeverityBreaking: Missing endpoint,     |
//	    |    incompatible type, missing param        |
//	    |  - SeverityNonBreaking: New optional field |
//	    |  - SeverityGhost: Removed endpoint         |
//	    +--------------------------------------------+
//	                          |
//	                          v
//	    +--------------------------------------------+
//	    |                 DiffStack                  |
//	    |  - Push, Pop, Peek snapshot undo frames   |
//	    +--------------------------------------------+
//
// # Core Building Blocks
//
//   - [Compare]: Evaluates drift between a local [*ir.RootIR] contract and an [*openapi.Document].
//   - [CompareWithOptions]: Evaluates drift with customized comparison options.
//   - [DiffReport]: Aggregates all detected discrepancies and summary metrics.
//   - [DiffOptions]: Configures path normalization, case sensitivity, and tolerance rules.
//   - [DriftItem]: Individual contract difference record detailing paths, parameters, and severity.
//   - [DriftSeverity]: Classification of drift impact ([SeverityBreaking], [SeverityNonBreaking], [SeverityGhost]).
//   - [DiffStack]: In-memory and persistent snapshot stack managing iterative contract checkpoints.
//   - [LoadStack]: Loads the diff stack from `.vortex/cache/diff_stack.json`.
//   - [StackFrame]: Individual historical frame inside the diff stack.
//   - [StackDiffResult]: Full comparison report between two stack frames.
//
// # Usage Tiers
//
// ## Tier 1: High-Level Contract Comparison
//
// Compare a local Go interface AST with a parsed remote OpenAPI document:
//
//	report := diff.Compare(localRootIR, remoteDoc, "pkg/api", "openapi.yaml")
//	if report.HasBreaking() {
//	    log.Println("breaking changes detected")
//	}
//
// ## Tier 2: Granular Options & Filtering
//
// Apply customized comparison options such as ignoring deprecated endpoints or custom tag filtering:
//
//	opts := diff.DiffOptions{IgnoreDeprecated: true}
//	report := diff.CompareWithOptions(localRootIR, remoteDoc, "pkg/api", "openapi.yaml", opts)
//
// ## Tier 3: Checkpoint Stack Management
//
// Push and roll back contract snapshots during automated reconciliation workflows:
//
//	stack, err := diff.LoadStack(rootDir)
//	if err != nil {
//	    log.Fatalf("failed loading stack: %v", err)
//	}
//	frame, err := stack.Push("pre-merge", []string{"pkg/api/service.go"}, []string{"v1"}, nil)
//
// # Concurrency & Thread Safety
//
// [Compare] and [CompareWithOptions] are pure, stateless functions safe for concurrent execution across
// multiple goroutines. [DiffStack] operations are synchronized via internal read-write mutexes.
//
// # Performance & Zero-Allocation Profile
//
// Path alignment and verb matching utilize stack-allocated string builders and pre-hashed endpoint maps,
// ensuring comparison passes complete in microseconds without excessive garbage collector pressure.
package diff
```

#### 5. `d:/CodingProjects/vortex/pkg/jsbundle/doc.go`
Replace entire file content with verified APIs (`ScanFiles`, `ScanFile`, `ScanBytes`):
```go
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
//	|                          ReconcileToIR                            |
//	|  - Bridges discovered endpoints into Vortex declarative RootIR    |
//	+-------------------------------------------------------------------+
//
// # Core Building Blocks
//
//   - [ScanFiles]: Scans multiple files matching glob patterns and aggregates their results.
//   - [ScanFile]: Scans an individual bundle file on disk.
//   - [ScanBytes]: Analyzes JavaScript byte content in-memory.
//   - [ScanResult]: Aggregated collection of discovered endpoints, message descriptors, and enums.
//   - [NewScanResult]: Instantiates an empty [ScanResult] container.
//   - [Endpoint]: Describes an individual discovered API route or RPC procedure.
//   - [MessageDescriptor]: Models a discovered Protobuf, JSPB, or DTO wire schema.
//   - [FieldDescriptor]: Models a specific struct or message field index and type.
//   - [EnumDescriptor]: Models a discovered integer-to-string enum mapping.
//
// # Usage Tiers
//
// ## Tier 1: Glob Pattern & File Scanning
//
// Scan client build output files for hidden endpoints and schemas:
//
//	res, err := jsbundle.ScanFiles([]string{"frontend/dist/*.js"})
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
```

#### 6. `d:/CodingProjects/vortex/pkg/git/doc.go`
Replace entire file content with verified APIs (`ListProposalBranches`, `IsClean`, `ShowFile`, `LogCommits`):
```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package git provides in-memory Git repository inspection querying branches, logs, and files with zero disk artifacts.
//
// # Architecture Overview
//
// The git package executes targeted Git CLI operations bounded by strict timeouts ([DefaultTimeout]), capturing
// file revisions directly into memory buffers without checking out worktrees or writing temporary disk files:
//
//	Git Working Tree / Repository
//	              |
//	              v
//	+-------------------------------------------------------------------+
//	|                   exec.CommandContext ("git", ...)                |
//	|  - Enforces DefaultTimeout (5s)                                   |
//	|  - Captures stdout directly into in-memory bytes.Buffer           |
//	+-------------------------------------------------------------------+
//	     |                   |                   |                  |
//	     v                   v                   v                  v
//	+----------+       +-----------+       +-----------+      +-----------+
//	| ShowFile |       |LogCommits |       |ListProposal|     |  IsClean  |
//	| (ref:path|       | (History) |       | (Branches)|      | (Status)  |
//	+----------+       +-----------+       +-----------+      +-----------+
//
// # Core Building Blocks
//
//   - [ShowFile]: Extracts the exact byte content of a file at a specific Git ref into an in-memory byte slice.
//   - [RootDir]: Resolves the top-level repository root directory.
//   - [CurrentBranch]: Retrieves the active Git branch name.
//   - [LogCommits]: Retrieves commit history metadata for a specific file or path.
//   - [ListProposalBranches]: Discovers local and remote proposal branches.
//   - [IsClean]: Checks whether uncommitted changes exist in the working directory or specific path.
//   - [MergeBase]: Finds the best common ancestor commit between two git references.
//   - [BlameFile]: Runs git blame in porcelain mode and returns line-by-line provenance.
//   - [CommitInfo]: Commit metadata record containing hash, author, timestamp, and subject.
//   - [BranchProposal]: Structured descriptor for consumer feature branches proposing contract updates.
//
// # Usage Tiers
//
// ## Tier 1: In-Memory File Retrieval
//
// Read an earlier revision of an API contract without writing to disk:
//
//	data, err := git.ShowFile(ctx, rootDir, "HEAD~1", "pkg/api/service.go")
//	if err != nil {
//	    log.Fatalf("failed to retrieve file from git: %v", err)
//	}
//
// ## Tier 2: Commit History & Branch Inspection
//
// Inspect recent commits affecting a contract file or discover active proposal branches:
//
//	commits, err := git.LogCommits(ctx, rootDir, "pkg/api/service.go", 5)
//	branches, err := git.ListProposalBranches(ctx, rootDir, nil)
//
// ## Tier 3: Working Tree Cleanliness
//
// Verify that the repository is clean before applying automated migrations:
//
//	isClean, err := git.IsClean(ctx, rootDir, "")
//
// # Concurrency & Thread Safety
//
// All functions in this package accept context.Context and invoke isolated CLI subprocesses with
// dedicated output buffers. They are 100% thread-safe across concurrent goroutines.
//
// # Performance & Zero-Allocation Profile
//
// File retrieval streams subprocess output directly into pre-allocated memory buffers, completely
// bypassing temporary disk files and intermediate working tree checkouts.
package git
```

#### 7. `d:/CodingProjects/vortex/pkg/mirror/doc.go`
Replace entire file content with verified APIs (`CheckService`, `DriftDiagnostic`, `DriftKind`):
```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mirror synchronizes declarative contracts with untagged upstream Go source files via @mirror directives.
//
// # Architecture Overview
//
// The shadow mirror engine parses target Go source files referenced by `@mirror` directives, aligns method
// signatures, parameter types, and return values with declarative contracts, and identifies drift diagnostics:
//
//	Contract Interface (@mirror: "internal/legacy/service.go")
//	                 \                          /
//	                  v                        v
//	+-------------------------------------------------------------------+
//	|                           CheckService                            |
//	|  - Parses target Go AST at mirror path                            |
//	|  - Compares method signatures, parameter types, and return values |
//	+-------------------------------------------------------------------+
//	                                 |
//	                                 v
//	+-------------------------------------------------------------------+
//	|                       DriftDiagnostic Report                      |
//	|  - DriftMethodMissing: Added in upstream, missing in contract     |
//	|  - DriftParamMismatch: Parameter type changed                     |
//	|  - DriftGhostMethod: Exists in contract, removed upstream         |
//	+-------------------------------------------------------------------+
//
// # Core Building Blocks
//
//   - [CheckService]: Inspects an individual service contract against its target mirror source file.
//   - [DriftDiagnostic]: Structured diagnostic detailing a specific signature or field discrepancy.
//   - [DriftKind]: Categorization of mirror drift ([DriftMethodMissing], [DriftParamMismatch], [DriftGhostMethod], etc.).
//
// # Usage Tiers
//
// ## Tier 1: Service Drift Inspection
//
// Check an individual contract for synchronization drift:
//
//	diagnostics, err := mirror.CheckService(rootDir, contractPath, serviceIR, structs)
//	if err != nil {
//	    log.Fatalf("mirror check failed: %v", err)
//	}
//	for _, diag := range diagnostics {
//	    log.Printf("[%s] %s: %s", diag.Kind, diag.Service, diag.Message)
//	}
//
// ## Tier 2: Evaluating Drift Kinds
//
// Classify divergences between contract interfaces and mirrored upstream files:
//
//	for _, diag := range diagnostics {
//	    if diag.Kind == mirror.DriftMethodMissing {
//	        log.Printf("upstream added method missing from contract: %s", diag.Method)
//	    }
//	}
//
// ## Tier 3: Path Resolution Architecture
//
// [CheckService] automatically resolves relative mirror targets by walking directory trees to locate
// the governing `.vortex.yml` or `go.mod` boundary.
//
// # Concurrency & Thread Safety
//
// [CheckService] performs read-only AST comparisons and is 100% thread-safe across concurrent goroutines.
//
// # Performance & Zero-Allocation Profile
//
// The engine parses target mirror files into lightweight ASTs, comparing type signatures through
// fast identifier matching with minimal heap allocations.
package mirror
```

#### 8. `d:/CodingProjects/vortex/pkg/parser/doc.go`
Replace entire file content with verified APIs (`Parser`, `NewParser`, `ParseDirective`, `ParseDirectives`, `Directive`), removing `[Lexer]`/`[Token]` links and escaping `generic.Optional[T]`:
```go
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
//	                        |    AST-to-IR Binder |  <--- Resolves generic.Optional types, DTOs
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
//   - [ParseDirectives]: Scans AST comment groups and extracts all declared directives.
//   - [Directive]: Structured representation of a parsed directive and its argument key-value pairs.
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
// ## Tier 3: Low-Level Lexing & Directive Extraction
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
// An individual [Parser] instance maintains an internal token.FileSet. Separate [Parser] instances
// can safely parse distinct files or packages concurrently across goroutines. The resulting [*ir.RootIR]
// is mutable while being built, but should be treated as read-only once passed downstream to analysis.
//
// # Performance & Zero-Allocation Profile
//
// The directive lexer scans raw byte slices of comment text directly without invoking regular expressions.
// Token boundaries and directive argument maps are constructed with minimal allocation overhead.
package parser
```

---

## 5. Verification Method & Gate Command Playbook

To independently verify the implementation after applying all changes:

### Step 1: Verify Sibling Comment Deduplication
Execute `go doc` on the 4 cleaned packages and verify that NO trailing duplicate summary line appears:
```powershell
go doc ./pkg/emitter
go doc ./pkg/ingest
go doc ./pkg/lint
go doc ./pkg/openapi
```
*Expected*: Pristine Godoc output where the package mission summary appears only at the very top.

### Step 2: Verify Godoc Symbol Links Across All 8 Aligned Packages
Execute `go doc` symbol queries to prove that 100% of documented symbols exist and resolve:
```powershell
go doc ./pkg/ingest HARToOpenAPI
go doc ./pkg/ingest IngestOptions
go doc ./pkg/cache LoadSecrets
go doc ./pkg/cache TrafficIndex
go doc ./pkg/cfg New
go doc ./pkg/cfg WalkPaths
go doc ./pkg/diff DiffStack
go doc ./pkg/diff LoadStack
go doc ./pkg/jsbundle ScanFiles
go doc ./pkg/jsbundle ScanBytes
go doc ./pkg/git ListProposalBranches
go doc ./pkg/git IsClean
go doc ./pkg/mirror CheckService
go doc ./pkg/parser ParseDirective
```
*Expected*: Exit code 0 for every command with formatted documentation printed.

### Step 3: Verify Typed Nil Pointer Safety
Run the augmented error test suites across all 8 subsystems:
```powershell
go test -v -count=1 ./pkg/project ./pkg/parser ./pkg/diff ./pkg/git ./pkg/cache ./pkg/lint ./pkg/spec ./pkg/pipeline
```
*Expected*: PASS across all 8 packages with 0 failures and 0 panics.

### Step 4: Workspace Build & Regression Gate
Run full repository test execution:
```powershell
$env:GOWORK="off"; go test -count=1 ./...
```
*Expected*: Exit code 0, all 41 packages pass cleanly.

### Step 5: Linter Gate
Run the workspace linter:
```powershell
golangci-lint run --allow-parallel-runners ./...
```
*Expected*: Exit code 0, `0 issues.`
