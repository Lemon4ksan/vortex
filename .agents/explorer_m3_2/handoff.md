# Handoff Report: Standardized Sentinel Errors & Typed Error Predicates Architecture

- **Agent**: `explorer_m3_2` (teamwork_preview_explorer)
- **Objective**: Standardized Sentinel Errors and Typed Error Predicates architecture across 8 key subsystems of Vortex
- **Workspace**: `d:/CodingProjects/vortex`
- **Timestamp**: 2026-09-22T20:08:00Z
- **Milestone**: Milestone 3 — Benchmark-Grade Code Documentation & Architecture

---

## 1. Observation

### 1.1 Codebase Error Handling Inventory
A systematic inspection of the 8 target subsystems in `pkg/` revealed 0 declared sentinel errors (`var Err...`), 0 typed error structs implementing `error`, and widespread reliance on ad-hoc string formatting:

1. **`pkg/project`**:
   - `pkg/project/config.go:414`: `return errors.New("nil configuration")`
   - `pkg/project/config.go:553`: `fmt.Errorf("discovering workspace root: %w", err)`
   - `pkg/project/config.go:560`: `fmt.Errorf("reading %s: %w", configPath, err)`
   - `pkg/project/config.go:565`: `fmt.Errorf("parsing %s: %w", configPath, err)`
   - `pkg/project/config.go:577`: `fmt.Errorf("validating %s: %w", configPath, err)`
   - `pkg/project/config.go:782`: `errors.New(".vortex.yml already exists (use --force to overwrite)")`
   - `pkg/project/config.go:858`: `errors.New("cannot save config: empty RootDir and ConfigPath")`
   - `pkg/project/config.go:1121`: `errors.New("no .vortex.work workspace file found in hierarchy")`
   - `pkg/project/status.go:44-45`: `c.IsGenStale = true; c.GenStaleReason = "api.gen.go is STALE"`

2. **`pkg/parser`**:
   - `pkg/parser/parser.go:37`: `return nil, fmt.Errorf("read file %s: %w", filePath, err)`
   - `pkg/parser/parser.go:47`: `return nil, fmt.Errorf("read dir %s: %w", dirPath, err)`
   - `pkg/parser/parser.go:108`: `return nil, fmt.Errorf("aoni/codegen/parser: syntax error in %q: %w", filename, err)`

3. **`pkg/diff`**:
   - `pkg/diff/stack.go:173`: `return nil, errors.New("no files specified and no candidate files found in workspace")`
   - `pkg/diff/stack.go:284`: `return nil, errors.New("stack is empty, nothing to pop")`
   - `pkg/diff/stack.go:379`: `return nil, -1, errors.New("stack is empty")`
   - `pkg/diff/stack.go:414`: `return nil, -1, fmt.Errorf("frame matching label or tag %q not found in stack", query)`
   - `pkg/diff/stack.go:438`: `return nil, errors.New("at least 2 frames required in stack for adjacent diff")`
   - `pkg/diff/stack.go:453`: `return nil, errors.New("at least 2 frames required in stack for cumulative diff")`
   - `pkg/diff/stack.go:967`: `return nil, fmt.Errorf("target frame %q not found in stack", targetQuery)`

4. **`pkg/git`**:
   - `pkg/git/git.go:65`: `fmt.Errorf("git show %s failed: %s (err: %w)", target, strings.TrimSpace(stderr.String()), err)`
   - `pkg/git/git.go:89`: `fmt.Errorf("git merge-base %s %s failed: %s (err: %w)", refA, refB, ...)`
   - `pkg/git/git.go:125`: `fmt.Errorf("git branch failed: %s (err: %w)", strings.TrimSpace(stderr.String()), err)`
   - `pkg/git/git.go:223`: `fmt.Errorf("git log failed: %s (err: %w)", strings.TrimSpace(stderr.String()), err)`
   - `pkg/git/git.go:270`: `fmt.Errorf("git rev-parse HEAD failed: %s (err: %w)", strings.TrimSpace(stderr.String()), err)`
   - `pkg/git/git.go:300`: `fmt.Errorf("git status failed: %s (err: %w)", strings.TrimSpace(stderr.String()), err)`
   - `pkg/git/git.go:324`: `return "", errors.New("not a git repository (or any parent directory)")`

5. **`pkg/cache`**:
   - `pkg/cache/traffic_store.go:59`: `fmt.Errorf("parsing traffic index: %w", err)`
   - `pkg/cache/traffic_store.go:169`: `fmt.Errorf("compressing traffic blob: %w", err)`
   - `pkg/cache/traffic_store.go:208`: `fmt.Errorf("vortex: traffic session %q not found in cache", idOrHash)`
   - `pkg/cache/traffic_store.go:250`: `fmt.Errorf("decompressing traffic blob: %w", err)`
   - `pkg/cache/secrets_vault.go:170`: `fmt.Errorf("parsing secrets vault: %w", err)`
   - `pkg/cache/lint_cache.go:48`: `fmt.Errorf("parsing lint cache: %w", err)`

6. **`pkg/lint`**:
   - `pkg/lint/engine.go:64`: `fmt.Errorf("apply fix for [%s] %s: %w", d.RuleID, d.Message, err)`
   - `pkg/lint/engine.go:91`: `errors.New("lint: pass context cannot be nil")`

7. **`pkg/spec`**:
   - Contains directive definitions and specifications; lacking dedicated exported error sentinels for missing specs, unknown formats, and empty registries.

8. **`pkg/pipeline`**:
   - `pkg/pipeline/ast.go:54, 63, 158, 240, 476`: `return errors.New("target file is required")` (repeated verbatim 5 times across AST methods)
   - `pkg/pipeline/traffic.go:88`: `return 0, errors.New("no contract files found to extract environment variables")`

### 1.2 Go 1.27 `errors.AsType` Verification
Command execution `go doc errors.AsType` verified the standard library signature and semantics:
```go
package errors // import "errors"

func AsType[E error](err error) (E, bool)
```
- `AsType` finds the first error in `err`'s unwrapping tree matching type `E`.
- If found, it returns `(E, true)`; otherwise `(zero[E], false)`.
- It eliminates the reflection, interface boxing, and pointer-to-pointer indirection required by legacy `errors.As(err, &target)`.

### 1.3 Test Suite Assertions Inspection
A search for `require.EqualError`, `require.Contains(..., err)`, and `err.Error()` across all `*_test.go` files showed:
- Only one test asserts on an error string in the entire workspace: `cmd/vortex/app_test.go:360` (`require.Contains(t, err.Error(), "unknown command")`), which is generated directly in `cmd/vortex/app.go` and does not depend on `pkg/` errors.
- No unit tests in `pkg/*` assert on brittle exact error strings.

### 1.4 Aoni & Foundation Precedent
Inspection of `d:/CodingProjects/aoni/errors.go` verified:
- Exported sentinel errors declared at package scope (`var Err... = errors.New("<pkg>: ...")`).
- Typed domain error structs with `Op`, `Path`/`URL`, `Err` implementing `Error() string` and `Unwrap() error`.
- Predicate functions `func Is<Predicate>(err error) bool` combining `errors.Is` with `errors.AsType[*<Domain>Error]`.

---

## 2. Logic Chain

1. **Need for Standardized Sentinel Errors**:
   Callers in CLI, linters, or build scripts currently cannot distinguish between a missing file (`os.ErrNotExist`), a malformed contract (`ErrInvalidConfig`), a syntax error (`ErrSyntaxError`), or a missing git repository (`ErrNotRepository`) without resorting to string parsing. Defining exported sentinel errors (`var Err... = errors.New("<pkg>: <concept>")`) provides stable identities for `errors.Is`.

2. **Uniform Structural Layout (`Op, Path, Key, Err`)**:
   Standardizing every subsystem error struct to:
   ```go
   type <Subsystem>Error struct {
       Op   string // Action being executed (e.g. "load", "parse", "show", "run_pass")
       Path string // Target path or file (e.g. ".vortex.yml", "pkg/api.go")
       Key  string // Subsystem identifier (contract name, rule ID, frame label, git ref)
       Err  error  // Underlying root cause or sentinel error
   }
   ```
   delivers total architectural symmetry across Vortex. Callers, loggers, and formatters can inspect any error uniformly without bespoke reflection.

3. **High-Performance Typed Predicates via Go 1.27 `errors.AsType`**:
   Traditional `errors.As` uses `reflect` and takes `any`. In Go 1.27, `errors.AsType[*<Subsystem>Error](err)` compiles to direct type assertions while traversing the unwrapped error tree. Each predicate `func Is<Predicate>(err error) bool` will:
   - Check `errors.Is(err, ErrSentinel)`
   - Extract `errors.AsType[*<Subsystem>Error](err)` to inspect `e.Err`
   - Provide fallback matching for legacy raw errors (e.g. `strings.Contains`) for backwards compatibility.

4. **Zero Heap Allocation & Unwrapping**:
   Implementing `Unwrap() error { return e.Err }` on every `<Subsystem>Error` ensures native interoperability with standard `errors.Is`, `errors.As`, and `errors.AsType`.

---

## 3. Caveats

1. **Pointer Receiver on `Unwrap()` and `Error()`**:
   `errors.AsType[*<Subsystem>Error](err)` requires `*<Subsystem>Error` as the generic type argument because the `Error()` and `Unwrap()` methods are declared on the pointer receiver. Returning non-pointer values `SubsystemError` must be avoided in constructors to maintain `error` interface satisfaction.
2. **`os.ErrNotExist` Interoperability**:
   `IsNotFound` predicates in `pkg/project`, `pkg/cache`, and `pkg/spec` should also check `errors.Is(err, os.ErrNotExist)` because low-level file operations directly bubble filesystem missing errors.
3. **No Production Modifications during Investigation**:
   Per Explorer constraints, this report provides complete production-ready source code designs in this document for implementer agents (`implementer_m3_*`) without directly editing production files.

---

## 4. Conclusion & Concrete Design

Below are the complete, concrete `errors.go` implementations for the 8 target subsystems.

---

### 4.1 Subsystem 1: `pkg/project/errors.go`

```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package project

import (
	"errors"
	"os"
	"strings"
)

var (
	// ErrWorkspaceNotFound is returned when workspace root or boundary (go.mod, .git, .vortex.yml) cannot be located.
	ErrWorkspaceNotFound = errors.New("project: workspace root not found")

	// ErrConfigNotFound is returned when .vortex.yml or .vortex.work configuration file does not exist.
	ErrConfigNotFound = errors.New("project: configuration file not found")

	// ErrInvalidConfig is returned when the configuration schema or contract definition fails validation.
	ErrInvalidConfig = errors.New("project: invalid configuration")

	// ErrContractNotFound is returned when a requested contract name or file cannot be resolved in the workspace.
	ErrContractNotFound = errors.New("project: contract not found")

	// ErrStaleCodegen is returned or reported when generated artifacts (.gen.go) are missing or out of sync with contract sources.
	ErrStaleCodegen = errors.New("project: generated code is stale")
)

// ProjectError represents an operational or structural error within the project workspace subsystem.
type ProjectError struct {
	Op   string // Operation (e.g. "load", "validate", "save", "find_root", "inspect")
	Path string // Filesystem path to config, workspace, or contract file
	Key  string // Contract name, configuration key, or service package
	Err  error  // Underlying cause or sentinel error
}

func (e *ProjectError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("project")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" [")
			sb.WriteString(e.Key)
			sb.WriteByte(']')
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString("[")
		sb.WriteString(e.Key)
		sb.WriteString("]: ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown project error")
	}

	return sb.String()
}

func (e *ProjectError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsNotFound reports whether err indicates a missing workspace, configuration, or contract resource.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrWorkspaceNotFound) ||
		errors.Is(err, ErrConfigNotFound) ||
		errors.Is(err, ErrContractNotFound) ||
		errors.Is(err, os.ErrNotExist) {
		return true
	}

	if pErr, ok := errors.AsType[*ProjectError](err); ok {
		return errors.Is(pErr.Err, ErrWorkspaceNotFound) ||
			errors.Is(pErr.Err, ErrConfigNotFound) ||
			errors.Is(pErr.Err, ErrContractNotFound) ||
			errors.Is(pErr.Err, os.ErrNotExist)
	}

	return false
}

// IsStale reports whether err indicates stale generated code.
func IsStale(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrStaleCodegen) {
		return true
	}

	if pErr, ok := errors.AsType[*ProjectError](err); ok {
		return errors.Is(pErr.Err, ErrStaleCodegen)
	}

	return false
}
```

---

### 4.2 Subsystem 2: `pkg/parser/errors.go`

```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"errors"
	"strings"
)

var (
	// ErrSyntaxError indicates a Go syntax or DSL parsing failure in source code.
	ErrSyntaxError = errors.New("parser: syntax error")

	// ErrContractNotFound indicates that no declarative service contract was discovered in the parsed source.
	ErrContractNotFound = errors.New("parser: contract not found")

	// ErrInvalidDirective indicates an unrecognized, misplaced, or malformed DSL directive or argument.
	ErrInvalidDirective = errors.New("parser: invalid directive")

	// ErrUnresolvedType indicates a referenced type or generic argument cannot be bound or resolved.
	ErrUnresolvedType = errors.New("parser: unresolved type")
)

// ParseError represents a syntax or semantic error discovered during contract AST parsing.
type ParseError struct {
	Op   string // Operation (e.g. "parse_file", "parse_source", "bind_type", "parse_directive")
	Path string // Source file path
	Key  string // Line:column position, directive name, or symbol identifier
	Err  error  // Underlying syntax or sentinel error
}

func (e *ParseError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("parser")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteByte(':')
			sb.WriteString(e.Key)
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString(e.Key)
		sb.WriteString(": ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown parse error")
	}

	return sb.String()
}

func (e *ParseError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsSyntaxError reports whether err represents a syntax parsing error in Go or DSL code.
func IsSyntaxError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrSyntaxError) {
		return true
	}

	if pErr, ok := errors.AsType[*ParseError](err); ok {
		return errors.Is(pErr.Err, ErrSyntaxError)
	}

	return strings.Contains(err.Error(), "syntax error")
}

// IsNotFound reports whether err represents a missing contract in the parsed source.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrContractNotFound) {
		return true
	}

	if pErr, ok := errors.AsType[*ParseError](err); ok {
		return errors.Is(pErr.Err, ErrContractNotFound)
	}

	return false
}
```

---

### 4.3 Subsystem 3: `pkg/diff/errors.go`

```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"errors"
	"strings"
)

var (
	// ErrStackEmpty indicates an operation was invoked on an empty diff stack.
	ErrStackEmpty = errors.New("diff: stack is empty")

	// ErrFrameNotFound indicates that a requested stack frame label or index does not exist.
	ErrFrameNotFound = errors.New("diff: frame not found")

	// ErrInsufficientFrames indicates that adjacent or cumulative diff requires at least 2 frames.
	ErrInsufficientFrames = errors.New("diff: insufficient frames for diff")

	// ErrConflict indicates conflicting AST or schema migrations between diff frames.
	ErrConflict = errors.New("diff: conflicting schema changes")
)

// DiffError represents an error during diff calculation or stack manipulation.
type DiffError struct {
	Op   string // Operation (e.g. "push", "pop", "pop_to", "diff_adjacent", "diff_cumulative")
	Path string // Stack file path or affected target file
	Key  string // Frame index, frame label, or query
	Err  error  // Underlying or sentinel error
}

func (e *DiffError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("diff")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" [frame=")
			sb.WriteString(e.Key)
			sb.WriteByte(']')
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString("frame ")
		sb.WriteString(e.Key)
		sb.WriteString(": ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown diff error")
	}

	return sb.String()
}

func (e *DiffError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsNotFound reports whether err represents a missing stack frame.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrFrameNotFound) {
		return true
	}

	if dErr, ok := errors.AsType[*DiffError](err); ok {
		return errors.Is(dErr.Err, ErrFrameNotFound)
	}

	return false
}

// IsConflict reports whether err represents a schema or contract migration conflict.
func IsConflict(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrConflict) {
		return true
	}

	if dErr, ok := errors.AsType[*DiffError](err); ok {
		return errors.Is(dErr.Err, ErrConflict)
	}

	return false
}

// IsEmpty reports whether err represents an empty stack condition.
func IsEmpty(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrStackEmpty) {
		return true
	}

	if dErr, ok := errors.AsType[*DiffError](err); ok {
		return errors.Is(dErr.Err, ErrStackEmpty)
	}

	return strings.Contains(err.Error(), "stack is empty")
}
```

---

### 4.4 Subsystem 4: `pkg/git/errors.go`

```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package git

import (
	"errors"
	"strings"
)

var (
	// ErrNotRepository indicates the directory is not inside a Git work tree or repository.
	ErrNotRepository = errors.New("git: not a git repository")

	// ErrBranchNotFound indicates that a requested branch, tag, or ref could not be resolved.
	ErrBranchNotFound = errors.New("git: branch not found")

	// ErrGitCommandFailed indicates that a git CLI subprocess exited with a non-zero exit status.
	ErrGitCommandFailed = errors.New("git: command execution failed")
)

// GitError represents an execution or resolution failure in the Git subsystem.
type GitError struct {
	Op   string // Operation (e.g. "show", "merge_base", "log", "branch", "rev_parse", "blame")
	Path string // Repository root or relative file path
	Key  string // Ref, branch name, or commit hash
	Err  error  // Underlying exec.ExitError or sentinel error
}

func (e *GitError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("git")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" (ref: ")
			sb.WriteString(e.Key)
			sb.WriteByte(')')
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString("ref ")
		sb.WriteString(e.Key)
		sb.WriteString(": ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown git error")
	}

	return sb.String()
}

func (e *GitError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsNotRepository reports whether err indicates that the target path is not a Git repository.
func IsNotRepository(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrNotRepository) {
		return true
	}

	if gErr, ok := errors.AsType[*GitError](err); ok {
		if errors.Is(gErr.Err, ErrNotRepository) {
			return true
		}
	}

	return strings.Contains(err.Error(), "not a git repository")
}
```

---

### 4.5 Subsystem 5: `pkg/cache/errors.go`

```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache

import (
	"errors"
	"os"
	"strings"
)

var (
	// ErrSessionNotFound indicates that a recorded traffic capture session ID was not found in the index.
	ErrSessionNotFound = errors.New("cache: session not found")

	// ErrSecretNotFound indicates that a requested secret rule or variable was not found in the vault.
	ErrSecretNotFound = errors.New("cache: secret not found")

	// ErrCorruptCache indicates that a cached artifact or index file is corrupt or has an invalid schema.
	ErrCorruptCache = errors.New("cache: corrupt cache data")
)

// CacheError represents an error during cache read, write, or index lookup.
type CacheError struct {
	Op   string // Operation (e.g. "load_session", "store_session", "load_vault", "save_lint")
	Path string // Cache file or storage directory path
	Key  string // Session ID, hash, or secret key
	Err  error  // Underlying IO or sentinel error
}

func (e *CacheError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("cache")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" [key=")
			sb.WriteString(e.Key)
			sb.WriteByte(']')
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString("key ")
		sb.WriteString(e.Key)
		sb.WriteString(": ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown cache error")
	}

	return sb.String()
}

func (e *CacheError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsNotFound reports whether err indicates a missing session, secret, or cache entry.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrSessionNotFound) ||
		errors.Is(err, ErrSecretNotFound) ||
		errors.Is(err, os.ErrNotExist) {
		return true
	}

	if cErr, ok := errors.AsType[*CacheError](err); ok {
		return errors.Is(cErr.Err, ErrSessionNotFound) ||
			errors.Is(cErr.Err, ErrSecretNotFound) ||
			errors.Is(cErr.Err, os.ErrNotExist)
	}

	return strings.Contains(err.Error(), "not found in cache")
}
```

---

### 4.6 Subsystem 6: `pkg/lint/errors.go`

```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lint

import (
	"errors"
	"strings"
)

var (
	// ErrLintFailure is returned when contracts fail static analysis checks (e.g. SeverityError count > 0).
	ErrLintFailure = errors.New("lint: static analysis checks failed")

	// ErrRuleNotFound indicates that a requested rule ID is not registered in the lint registry.
	ErrRuleNotFound = errors.New("lint: rule not found")

	// ErrFixFailed indicates that an automated code fix could not be successfully applied.
	ErrFixFailed = errors.New("lint: automated fix failed")
)

// LintError represents an operational error in rule execution or automated fix application.
type LintError struct {
	Op   string // Operation (e.g. "run_pass", "apply_fix", "register_rule")
	Path string // Inspected source file path
	Key  string // Rule ID (e.g. "S001", "B002", "missing-context")
	Err  error  // Underlying diagnostic or sentinel error
}

func (e *LintError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("lint")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" [rule=")
			sb.WriteString(e.Key)
			sb.WriteByte(']')
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString("[")
		sb.WriteString(e.Key)
		sb.WriteString("]: ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown lint error")
	}

	return sb.String()
}

func (e *LintError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsLintFailure reports whether err represents a static analysis check failure.
func IsLintFailure(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrLintFailure) {
		return true
	}

	if lErr, ok := errors.AsType[*LintError](err); ok {
		return errors.Is(lErr.Err, ErrLintFailure)
	}

	return false
}
```

---

### 4.7 Subsystem 7: `pkg/spec/errors.go`

```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spec

import (
	"errors"
	"os"
	"strings"
)

var (
	// ErrSpecNotFound indicates that an OpenAPI, AsyncAPI, or HAR spec definition was not found.
	ErrSpecNotFound = errors.New("spec: specification not found")

	// ErrUnsupportedFormat indicates that an unknown or invalid specification format was provided.
	ErrUnsupportedFormat = errors.New("spec: unsupported specification format")

	// ErrEmptySpec indicates that the provided specification document contains no paths, channels, or services.
	ErrEmptySpec = errors.New("spec: specification is empty")
)

// SpecError represents an error encountered while resolving or validating API specifications.
type SpecError struct {
	Op   string // Operation (e.g. "fetch", "parse", "validate", "lookup_directive")
	Path string // File path or URL of the specification
	Key  string // Specification format, directive name, or component key
	Err  error  // Underlying or sentinel error
}

func (e *SpecError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("spec")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" (format: ")
			sb.WriteString(e.Key)
			sb.WriteByte(')')
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString("format ")
		sb.WriteString(e.Key)
		sb.WriteString(": ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown spec error")
	}

	return sb.String()
}

func (e *SpecError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsNotFound reports whether err represents a missing specification file or resource.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrSpecNotFound) || errors.Is(err, os.ErrNotExist) {
		return true
	}

	if sErr, ok := errors.AsType[*SpecError](err); ok {
		return errors.Is(sErr.Err, ErrSpecNotFound) || errors.Is(sErr.Err, os.ErrNotExist)
	}

	return false
}
```

---

### 4.8 Subsystem 8: `pkg/pipeline/errors.go`

```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipeline

import (
	"errors"
	"strings"
)

var (
	// ErrTargetFileRequired is returned when an AST or compilation pipeline is invoked without a target file.
	ErrTargetFileRequired = errors.New("pipeline: target file is required")

	// ErrNoContractsFound is returned when no contracts are discovered to execute the pipeline on.
	ErrNoContractsFound = errors.New("pipeline: no contracts found")

	// ErrPipelineAborted is returned when the compilation or deobfuscation pipeline halts due to an unrecoverable condition.
	ErrPipelineAborted = errors.New("pipeline: processing aborted")
)

// PipelineError represents a processing or compilation failure during pipeline execution.
type PipelineError struct {
	Op   string // Operation (e.g. "deobfuscate", "rename_field", "extract_env", "compile")
	Path string // Target contract file path
	Key  string // Pipeline stage name or transformation target
	Err  error  // Underlying or sentinel error
}

func (e *PipelineError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("pipeline")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" [stage=")
			sb.WriteString(e.Key)
			sb.WriteByte(']')
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString("stage ")
		sb.WriteString(e.Key)
		sb.WriteString(": ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown pipeline error")
	}

	return sb.String()
}

func (e *PipelineError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsPipelineAborted reports whether err represents an aborted pipeline run.
func IsPipelineAborted(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrPipelineAborted) {
		return true
	}

	if pErr, ok := errors.AsType[*PipelineError](err); ok {
		return errors.Is(pErr.Err, ErrPipelineAborted)
	}

	return false
}
```

---

### 4.9 Subsystem Error Mapping Summary Matrix

| Package | Sentinels | Typed Error Struct | Typed Predicates (`errors.AsType`) |
|---|---|---|---|
| `pkg/project` | `ErrWorkspaceNotFound`, `ErrConfigNotFound`, `ErrInvalidConfig`, `ErrContractNotFound`, `ErrStaleCodegen` | `ProjectError{Op, Path, Key, Err}` | `IsNotFound(err)`, `IsStale(err)` |
| `pkg/parser` | `ErrSyntaxError`, `ErrContractNotFound`, `ErrInvalidDirective`, `ErrUnresolvedType` | `ParseError{Op, Path, Key, Err}` | `IsSyntaxError(err)`, `IsNotFound(err)` |
| `pkg/diff` | `ErrStackEmpty`, `ErrFrameNotFound`, `ErrInsufficientFrames`, `ErrConflict` | `DiffError{Op, Path, Key, Err}` | `IsNotFound(err)`, `IsConflict(err)`, `IsEmpty(err)` |
| `pkg/git` | `ErrNotRepository`, `ErrBranchNotFound`, `ErrGitCommandFailed` | `GitError{Op, Path, Key, Err}` | `IsNotRepository(err)` |
| `pkg/cache` | `ErrSessionNotFound`, `ErrSecretNotFound`, `ErrCorruptCache` | `CacheError{Op, Path, Key, Err}` | `IsNotFound(err)` |
| `pkg/lint` | `ErrLintFailure`, `ErrRuleNotFound`, `ErrFixFailed` | `LintError{Op, Path, Key, Err}` | `IsLintFailure(err)` |
| `pkg/spec` | `ErrSpecNotFound`, `ErrUnsupportedFormat`, `ErrEmptySpec` | `SpecError{Op, Path, Key, Err}` | `IsNotFound(err)` |
| `pkg/pipeline` | `ErrTargetFileRequired`, `ErrNoContractsFound`, `ErrPipelineAborted` | `PipelineError{Op, Path, Key, Err}` | `IsPipelineAborted(err)` |

---

## 5. Verification Method

Once implemented by `implementer_m3_*`, the error architecture must be verified as follows:

1. **Compilation Check**:
   ```bash
   go build ./pkg/...
   ```
   *Expected outcome*: Clean compilation with 0 syntax or generic type errors.

2. **Workspace Regression Check**:
   ```bash
   go test ./...
   ```
   *Expected outcome*: All existing unit tests pass without failure.

3. **Linter Gate**:
   ```bash
   golangci-lint run ./...
   ```
   *Expected outcome*: 0 warnings, zero unused variables, full compliance with `errorlint` and `revive`.

4. **Dedicated Predicate Unit Test Suite**:
   Create and execute test cases verifying that:
   - `errors.Is(err, ErrSentinel)` correctly matches wrapped instances.
   - `Is<Predicate>(err)` returns `true` for direct sentinels, wrapped typed errors, and `errors.AsType` matches.
   - `Is<Predicate>(err)` returns `false` for unrelated errors or `nil`.
   Example unit test template:
   ```go
   func TestProjectErrors(t *testing.T) {
       pErr := &project.ProjectError{
           Op:   "load",
           Path: ".vortex.yml",
           Err:  project.ErrConfigNotFound,
       }
       require.True(t, project.IsNotFound(pErr))
       require.True(t, errors.Is(pErr, project.ErrConfigNotFound))
       require.Contains(t, pErr.Error(), "project/load: .vortex.yml: project: configuration file not found")
   }
   ```
