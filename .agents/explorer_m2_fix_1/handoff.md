# Milestone 2 Concrete Remediation Blueprint

**Agent**: `explorer_m2_fix_1`  
**Working Directory**: `d:/CodingProjects/vortex/.agents/explorer_m2_fix_1/`  
**Milestone**: Milestone 2 — Restrained High-Craft CLI Presentation via `foundation/tuikit`  
**Target Date**: 2026-09-22T19:50:00Z  
**Role**: Teamwork Explorer (Read-Only Investigation & Remediation Strategy)  

---

## Executive Summary

Milestone 2 represents high craftsmanship: all 26 raw ANSI escape sequences have been eradicated from production source code, `foundation/tuikit` components (`Box`, `Table`, `Badge`, `ProbeTerminal`, `VisibleWidth`) are genuinely integrated, and NO_COLOR/pipe redirection is strictly enforced.

However, the forensic audit and adversarial reviews flagged 3 categories of defects:
1. **Informal Emoji & Arrow Leakage**: Remnants of `⚡`, `🤖`, `➔`, and `➜` in production code, causing 3 test failures in `cmd/vortex/adversarial_m2_test.go`.
2. **Linter Formatting Violations**: 5 issues detected by `golangci-lint run ./...` (3 `gci` section breaks in `cmd/vortex/app.go`, `internal/perf/prof.go`, `pkg/lint/format.go`, and 2 `golines` line wraps in `pkg/lint/format.go`, `internal/core/autopilot.go`).
3. **Concurrency & Quality Deficiencies**: Global state mutation of `tuikit.SetColorEnabled()` in `internal/text/render_terminal.go` (unsafe for concurrent renderers) and `len()` width calculation in `internal/workspace/doctor.go` (unsafe for wide Unicode glyphs).

Crucially, this investigation uncovered that replacing `➔` with `↳` in `pkg/diff/stack.go` and `internal/traffic/diff.go` directly affects assertions in `cmd/vortex/app_test.go`. Both the production source and the corresponding test assertions must be updated together to guarantee 100% test pass.

This document provides the exact file paths, line numbers, and verbatim Before/After code blocks for the implementation worker.

---

## 1. Observation

### 1.1 Forbidden Informal Emojis & Arrows
1. **`pkg/openapi/reconcile.go:59`**:
   ```go
   fmt.Fprintf(&sb, "⚡ [vortex merge] Merging %q into %q\n", s.SpecSource, targetPath)
   ```
   Invoked by `vortex spec import` CLI subcommand via `internal/spec/import.go:358, 528`. Leaks `⚡` to user terminal stdout.
2. **`pkg/tuple/analyzer.go:208`**:
   ```go
   fmt.Fprintf(&sb, "⚡ Vortex Tuple Saliency Analysis (%s across %d sample entries)\n\n",
   ```
   Outputs `⚡` during tuple saliency analysis table generation.
3. **`pkg/oracle/gen/js_emitter.go:664`**:
   ```javascript
   console.log(`🤖 Vortex Universal Oracle [${PORT}] listening on http://127.0.0.1:${PORT}`);
   ```
   Emits `🤖` in the generated JavaScript oracle server startup banner.
4. **`pkg/diff/stack.go:871, 882`**:
   ```go
   fmt.Fprintf(&buf, "  • %s [Tag #%s]: %s ➔ %s (%s)\n", tr.StructName, tr.Tag, tr.OldField, tr.NewField, tr.GoType)
   ...
   fmt.Fprintf(&buf, "  • %s: %s ➔ %s\n", mr.Route, mr.OldMethod, mr.NewMethod)
   ```
   Uses informal dingbat arrow `➔` (U+2794).
5. **`internal/traffic/diff.go:683`**:
   ```go
   fmt.Fprintf(stdout, "  • %s%s: %s ➔ %s      ➜ vortex ast rename --type=%s --field=%s --to=<NAME>\n",
   ```
   Uses informal dingbat arrows `➔` (U+2794) and `➜` (U+279C).

### 1.2 Downstream Test Failures Caused by Emojis
Running `go test ./cmd/vortex` produces exit code 1 with 3 failing tests:
- `TestMilestone2_Adversarial_NoInformalEmojisInCodebase` (fails on `pkg/openapi/reconcile.go:59` and `pkg/tuple/analyzer.go:208`)
- `TestMilestone2_Adversarial_SpecImport_NoInformalEmojis` (fails because `vortex spec import` prints `⚡`)
- `TestMilestone2_Adversarial_TupleRenderTable_NoInformalEmojis` (fails because `r.RenderTable()` emits `⚡`)

### 1.3 Downstream Test Assertions on Arrows in `cmd/vortex/app_test.go`
Searching for `➔` across test files revealed that `cmd/vortex/app_test.go` specifically asserts on the exact arrow string:
- Lines 1728 & 1731:
  ```go
  require.Contains(t, out, "65536 ➔ 8192")
  ...
  require.Contains(t, out, "0.7 ➔ 1")
  ```
- Lines 1960–1973:
  ```go
  require.Contains(t, adjDiff, "Field4 ➔ MaxTokens")
  require.Contains(t, adjDiff, "RPCMethod1 ➔ GenerateContent")
  require.NotContains(t, adjDiff, "Field0 ➔ ModelName")
  ...
  require.Contains(t, cumDiff, "Field0 ➔ ModelName")
  require.Contains(t, cumDiff, "Field4 ➔ MaxTokens")
  require.Contains(t, cumDiff, "RPCMethod1 ➔ GenerateContent")
  ```
Updating `pkg/diff/stack.go` and `internal/traffic/diff.go` to use `↳` without updating `cmd/vortex/app_test.go` will break these tests. They must be updated in lockstep.

### 1.4 Linter Violations (`golangci-lint run ./...`)
Executing `golangci-lint run ./...` reports 5 distinct violations across the codebase:
```text
cmd\vortex\app.go:17:1: File is not properly formatted (gci)
	"github.com/lemon4ksan/foundation/tuikit"

internal\perf\prof.go:26:1: File is not properly formatted (gci)
	"github.com/lemon4ksan/foundation/tuikit"

pkg\lint\format.go:18:1: File is not properly formatted (gci)
	"github.com/lemon4ksan/foundation/tuikit"

pkg\lint\format.go:183:1: File is not properly formatted (golines)
		fmt.Fprintf(&buf, "\n%s\n",

internal\core\autopilot.go:533:1: File is not properly formatted (golines)
		fmt.Sprintf("Workspace is 100%% healthy, synchronized, and compiled [%v | 0 allocs]!", elapsed.Round(time.Millisecond)),
```
In `.golangci.yml`, `gci.sections` specifies:
```yaml
sections:
  - standard
  - default
  - prefix(github.com/lemon4ksan/vortex)
```
In `app.go`, `prof.go`, and `format.go`, `"github.com/lemon4ksan/foundation/tuikit"` belongs to `default`, but is grouped with `"github.com/lemon4ksan/vortex/..."` without a separating blank line.
In `format.go` and `autopilot.go`, statement lines exceed the configured `max-len: 120`.

### 1.5 Concurrency & Quality Deficiencies
1. **`internal/text/render_terminal.go:178–181, 216–219`**:
   ```go
   if !r.isColorActive() {
       prev := tuikit.ColorEnabled()
       tuikit.SetColorEnabled(false)
       _ = box.Render(w)
       tuikit.SetColorEnabled(prev)
   }
   ```
   Mutating package-level `tuikit.SetColorEnabled()` during render is a concurrency anti-pattern. If multiple goroutines render concurrently (e.g. one in color and one plain), a race window is created. Local stripping via `tuikit.StripANSI` eliminates this completely.
2. **`internal/workspace/doctor.go:324`**:
   ```go
   maxNameWidth := 0
   for _, c := range rep.Checks {
       if len(c.Name) > maxNameWidth {
           maxNameWidth = len(c.Name)
       }
   }
   ```
   Uses byte length `len()` for terminal column alignment instead of visual cell width `tui.VisibleWidth()`.

---

## 2. Logic Chain

1. **Step 1**: `ORIGINAL_REQUEST.md` (R2, Acceptance Criteria) mandates zero informal emojis (`⚡`, `✨`, `🔴`, `🟡`, `🔵`, `❌`, `⚠️`, `🚀`), zero raw ANSI escapes, adoption of clean Unicode glyphs (`✔`, `✖`, `◆`, `↳`, `—`), 100% clean test passes (`go test ./...`), and zero lint violations (`golangci-lint run`).
2. **Step 2**: Observations 1.1 and 1.2 demonstrate that `pkg/openapi/reconcile.go:59` and `pkg/tuple/analyzer.go:208` cause 3 test failures in `cmd/vortex/adversarial_m2_test.go` and leak `⚡` to user terminals. Replacing them with sovereign glyph `◆` restores aesthetic discipline and resolves the test failures.
3. **Step 3**: Observation 1.1 shows `pkg/oracle/gen/js_emitter.go:664` contains `🤖`, and `pkg/diff/stack.go:871, 882` & `internal/traffic/diff.go:683` contain `➔` / `➜`. Replacing them with `◆` and `↳` aligns the entire codebase with the sovereign design language.
4. **Step 4**: Observation 1.3 reveals that updating `➔` to `↳` changes CLI output asserted by `cmd/vortex/app_test.go:1728, 1731, 1960-1973`. Therefore, updating `app_test.go` alongside the production code is required to maintain a passing test suite.
5. **Step 5**: Observation 1.4 confirms that `cmd/vortex/app.go`, `internal/perf/prof.go`, and `pkg/lint/format.go` violate `gci` sectioning, and `pkg/lint/format.go:183` & `internal/core/autopilot.go:533` violate `golines`. Inserting blank lines and wrapping long statements resolves all 5 linter violations.
6. **Step 6**: Observation 1.5 proves that temporary global state mutation in `internal/text/render_terminal.go` can be cleanly replaced by rendering to a buffer and calling `tuikit.StripANSI()`, and `internal/workspace/doctor.go:324` can be made Unicode-safe by calling `tui.VisibleWidth(c.Name)`.
7. **Conclusion**: Executing the concrete remediation plan below will achieve 100% compliance with Milestone 2 criteria, with zero test failures and zero lint violations.

---

## 3. Caveats

1. **Implementation Delegation**: In strict accordance with the Teamwork Explorer contract ("Read-only investigation — do NOT implement code changes directly"), all edits must be applied by the implementation worker (`worker_m2_fix_1`).
2. **Workspace Isolation**: Running `go test` or `golangci-lint` in PowerShell must include `$env:GOWORK="off"` to prevent interference from sibling repository modules.
3. **Linter Cache Concurrency**: Running `golangci-lint` simultaneously with another agent can cause a file lock error on the cache. If encountered, use a private cache via `$env:GOLANGCI_LINT_CACHE="$env:TEMP\golangci-lint-fix"`.

---

## 4. Conclusion & Concrete Remediation Plan

The worker must apply the following edits across the 10 target files.

---

### Part 1: Informal Emoji & Arrow Decontamination

#### 1.1 `pkg/openapi/reconcile.go` (Line 59)
- **Target File**: `d:/CodingProjects/vortex/pkg/openapi/reconcile.go`
- **Lines**: 57–60
- **Action**: Replace `⚡ [vortex merge]` with `◆ [vortex merge]`.

**Before**:
```go
// Render formats a human-readable terminal report of the merge.
func (s MergeSummary) Render(targetPath string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "⚡ [vortex merge] Merging %q into %q\n", s.SpecSource, targetPath)
```

**After**:
```go
// Render formats a human-readable terminal report of the merge.
func (s MergeSummary) Render(targetPath string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "◆ [vortex merge] Merging %q into %q\n", s.SpecSource, targetPath)
```

---

#### 1.2 `pkg/tuple/analyzer.go` (Line 208)
- **Target File**: `d:/CodingProjects/vortex/pkg/tuple/analyzer.go`
- **Lines**: 204–210
- **Action**: Replace `⚡ Vortex Tuple Saliency Analysis` with `◆ Vortex Tuple Saliency Analysis`.

**Before**:
```go
// RenderTable renders a neutral, clean terminal table of the tuple index analysis.
func (r *TupleAnalysisReport) RenderTable() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "⚡ Vortex Tuple Saliency Analysis (%s across %d sample entries)\n\n",
		r.StructName, r.TotalSamples)
```

**After**:
```go
// RenderTable renders a neutral, clean terminal table of the tuple index analysis.
func (r *TupleAnalysisReport) RenderTable() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "◆ Vortex Tuple Saliency Analysis (%s across %d sample entries)\n\n",
		r.StructName, r.TotalSamples)
```

---

#### 1.3 `pkg/oracle/gen/js_emitter.go` (Line 664)
- **Target File**: `d:/CodingProjects/vortex/pkg/oracle/gen/js_emitter.go`
- **Lines**: 663–665
- **Action**: Replace `🤖` with `◆`.

**Before**:
```go
server.listen(PORT, '127.0.0.1', () => {
  console.log(`+"`🤖 Vortex Universal Oracle [%q] listening on http://127.0.0.1:${PORT}`"+`);
  initBrowser().catch(err => {
```

**After**:
```go
server.listen(PORT, '127.0.0.1', () => {
  console.log(`+"`◆ Vortex Universal Oracle [%q] listening on http://127.0.0.1:${PORT}`"+`);
  initBrowser().catch(err => {
```

---

#### 1.4 `pkg/diff/stack.go` (Lines 871, 882)
- **Target File**: `d:/CodingProjects/vortex/pkg/diff/stack.go`
- **Lines**: 867–886
- **Action**: Replace `➔` with `↳`.

**Before**:
```go
		if len(r.ASTEvolution.TupleRenames) > 0 {
			fmt.Fprintf(&buf, "◆ Tuple Field Renames (Deobfuscation Dictionary):\n")

			for _, tr := range r.ASTEvolution.TupleRenames {
				fmt.Fprintf(&buf, "  • %s [Tag #%s]: %s ➔ %s (%s)\n",
					tr.StructName, tr.Tag, tr.OldField, tr.NewField, tr.GoType)
			}

			fmt.Fprintf(&buf, "\n")
		}

		if len(r.ASTEvolution.MethodRenames) > 0 {
			fmt.Fprintf(&buf, "◆ RPC / Method Renames:\n")

			for _, mr := range r.ASTEvolution.MethodRenames {
				fmt.Fprintf(&buf, "  • %s: %s ➔ %s\n", mr.Route, mr.OldMethod, mr.NewMethod)
			}

			fmt.Fprintf(&buf, "\n")
		}
```

**After**:
```go
		if len(r.ASTEvolution.TupleRenames) > 0 {
			fmt.Fprintf(&buf, "◆ Tuple Field Renames (Deobfuscation Dictionary):\n")

			for _, tr := range r.ASTEvolution.TupleRenames {
				fmt.Fprintf(&buf, "  • %s [Tag #%s]: %s ↳ %s (%s)\n",
					tr.StructName, tr.Tag, tr.OldField, tr.NewField, tr.GoType)
			}

			fmt.Fprintf(&buf, "\n")
		}

		if len(r.ASTEvolution.MethodRenames) > 0 {
			fmt.Fprintf(&buf, "◆ RPC / Method Renames:\n")

			for _, mr := range r.ASTEvolution.MethodRenames {
				fmt.Fprintf(&buf, "  • %s: %s ↳ %s\n", mr.Route, mr.OldMethod, mr.NewMethod)
			}

			fmt.Fprintf(&buf, "\n")
		}
```

---

#### 1.5 `internal/traffic/diff.go` (Line 683)
- **Target File**: `d:/CodingProjects/vortex/internal/traffic/diff.go`
- **Lines**: 677–686
- **Action**: Replace `➔` and `➜` with `↳`.

**Before**:
```go
		for _, it := range items {
			tagInfo := ""
			if it.Tag != "" {
				tagInfo = fmt.Sprintf(" (tag %s)", it.Tag)
			}

			fmt.Fprintf(stdout, "  • %s%s: %s ➔ %s      ➜ vortex ast rename --type=%s --field=%s --to=<NAME>\n",
				it.Field, tagInfo, it.OldVal, it.NewVal, it.Struct, it.Field)
		}
```

**After**:
```go
		for _, it := range items {
			tagInfo := ""
			if it.Tag != "" {
				tagInfo = fmt.Sprintf(" (tag %s)", it.Tag)
			}

			fmt.Fprintf(stdout, "  • %s%s: %s ↳ %s      ↳ vortex ast rename --type=%s --field=%s --to=<NAME>\n",
				it.Field, tagInfo, it.OldVal, it.NewVal, it.Struct, it.Field)
		}
```

---

#### 1.6 `cmd/vortex/app_test.go` (Lines 1728, 1731, 1960–1973) — CRITICAL CORRELATED TEST FIX
- **Target File**: `d:/CodingProjects/vortex/cmd/vortex/app_test.go`
- **Action**: Align test assertions with the sovereign `↳` arrow glyph.

**Hunk 1 (Lines 1724–1732)**:
**Before**:
```go
	out := stdout.String()
	require.Contains(t, out, "Traffic Diff")
	require.Contains(t, out, "GenerateContentRequest")
	require.Contains(t, out, "Field4")
	require.Contains(t, out, "65536 ➔ 8192")
	require.Contains(t, out, "vortex ast rename --type=GenerateContentRequest --field=Field4 --to=<NAME>")
	require.Contains(t, out, "Field5")
	require.Contains(t, out, "0.7 ➔ 1")
```

**After**:
```go
	out := stdout.String()
	require.Contains(t, out, "Traffic Diff")
	require.Contains(t, out, "GenerateContentRequest")
	require.Contains(t, out, "Field4")
	require.Contains(t, out, "65536 ↳ 8192")
	require.Contains(t, out, "vortex ast rename --type=GenerateContentRequest --field=Field4 --to=<NAME>")
	require.Contains(t, out, "Field5")
	require.Contains(t, out, "0.7 ↳ 1")
```

**Hunk 2 (Lines 1959–1974)**:
**Before**:
```go
	adjDiff := stdout.String()
	require.Contains(t, adjDiff, "Field4 ➔ MaxTokens")
	require.Contains(t, adjDiff, "RPCMethod1 ➔ GenerateContent")
	require.NotContains(t, adjDiff, "Field0 ➔ ModelName") // Happened in Frame 1, not Frame 2

	// 6. Test stack diff --cumulative (Frame 2 vs Frame 0)
	stdout.Reset()

	err = app.Run(context.Background(), []string{"ast", "stack", "diff", "--cumulative"})
	require.NoError(t, err)

	cumDiff := stdout.String()
	require.Contains(t, cumDiff, "Field0 ➔ ModelName")
	require.Contains(t, cumDiff, "Field4 ➔ MaxTokens")
	require.Contains(t, cumDiff, "RPCMethod1 ➔ GenerateContent")
```

**After**:
```go
	adjDiff := stdout.String()
	require.Contains(t, adjDiff, "Field4 ↳ MaxTokens")
	require.Contains(t, adjDiff, "RPCMethod1 ↳ GenerateContent")
	require.NotContains(t, adjDiff, "Field0 ↳ ModelName") // Happened in Frame 1, not Frame 2

	// 6. Test stack diff --cumulative (Frame 2 vs Frame 0)
	stdout.Reset()

	err = app.Run(context.Background(), []string{"ast", "stack", "diff", "--cumulative"})
	require.NoError(t, err)

	cumDiff := stdout.String()
	require.Contains(t, cumDiff, "Field0 ↳ ModelName")
	require.Contains(t, cumDiff, "Field4 ↳ MaxTokens")
	require.Contains(t, cumDiff, "RPCMethod1 ↳ GenerateContent")
```

---

### Part 2: Linter Formatting Compliance

#### 2.1 `cmd/vortex/app.go` (Line 17)
- **Target File**: `d:/CodingProjects/vortex/cmd/vortex/app.go`
- **Lines**: 15–20
- **Action**: Add an empty line between `foundation/tuikit` and `vortex/internal/base` to satisfy `gci`.

**Before**:
```go
	"runtime"
	"strings"

	"github.com/lemon4ksan/foundation/tuikit"
	"github.com/lemon4ksan/vortex/internal/base"
)
```

**After**:
```go
	"runtime"
	"strings"

	"github.com/lemon4ksan/foundation/tuikit"

	"github.com/lemon4ksan/vortex/internal/base"
)
```

---

#### 2.2 `internal/perf/prof.go` (Line 26)
- **Target File**: `d:/CodingProjects/vortex/internal/perf/prof.go`
- **Lines**: 24–31
- **Action**: Add an empty line between `foundation/tuikit` and `vortex/...` imports to satisfy `gci`.

**Before**:
```go
	"time"

	"github.com/lemon4ksan/foundation/tuikit"
	"github.com/lemon4ksan/vortex/internal/base"
	"github.com/lemon4ksan/vortex/internal/text"
	"github.com/lemon4ksan/vortex/pkg/project"
)
```

**After**:
```go
	"time"

	"github.com/lemon4ksan/foundation/tuikit"

	"github.com/lemon4ksan/vortex/internal/base"
	"github.com/lemon4ksan/vortex/internal/text"
	"github.com/lemon4ksan/vortex/pkg/project"
)
```

---

#### 2.3 `pkg/lint/format.go` (Line 18)
- **Target File**: `d:/CodingProjects/vortex/pkg/lint/format.go`
- **Lines**: 16–21
- **Action**: Add an empty line between `foundation/tuikit` and `vortex/pkg/version` to satisfy `gci`.

**Before**:
```go
	"strings"

	"github.com/lemon4ksan/foundation/tuikit"
	"github.com/lemon4ksan/vortex/pkg/version"
)
```

**After**:
```go
	"strings"

	"github.com/lemon4ksan/foundation/tuikit"

	"github.com/lemon4ksan/vortex/pkg/version"
)
```

---

#### 2.4 `pkg/lint/format.go` (Line 183)
- **Target File**: `d:/CodingProjects/vortex/pkg/lint/format.go`
- **Lines**: 182–186
- **Action**: Break statement to satisfy `golines` line length rule.

**Before**:
```go
	if report.FixableCount() > 0 {
		fmt.Fprintf(&buf, "\n%s\n",
			tuikit.Cyan(fmt.Sprintf("Run `vortex check --fix` to automatically resolve %d safe issue(s).", report.FixableCount())))
	}
```

**After**:
```go
	if report.FixableCount() > 0 {
		fixMsg := fmt.Sprintf(
			"Run `vortex check --fix` to automatically resolve %d safe issue(s).",
			report.FixableCount(),
		)
		fmt.Fprintf(&buf, "\n%s\n", tuikit.Cyan(fixMsg))
	}
```

---

#### 2.5 `internal/core/autopilot.go` (Line 533)
- **Target File**: `d:/CodingProjects/vortex/internal/core/autopilot.go`
- **Lines**: 531–535
- **Action**: Break `fmt.Sprintf` across multiple lines to satisfy `golines` line length rule.

**Before**:
```go
	doc.Success(
		"Workspace Synchronized",
		fmt.Sprintf("Workspace is 100%% healthy, synchronized, and compiled [%v | 0 allocs]!", elapsed.Round(time.Millisecond)),
	)
```

**After**:
```go
	doc.Success(
		"Workspace Synchronized",
		fmt.Sprintf(
			"Workspace is 100%% healthy, synchronized, and compiled [%v | 0 allocs]!",
			elapsed.Round(time.Millisecond),
		),
	)
```

---

### Part 3: Code Quality & Concurrency Fixes

#### 3.1 `internal/text/render_terminal.go` (Lines 177–184 & 215–222)
- **Target File**: `d:/CodingProjects/vortex/internal/text/render_terminal.go`
- **Lines**: 176–185 and 214–223
- **Action**: Eradicate global state mutation `tuikit.SetColorEnabled(false)` / `tuikit.SetColorEnabled(prev)`. Use `strings.Builder` and `tuikit.StripANSI` locally.

**In `renderCallout` (Lines 176–185)**:
**Before**:
```go
	if !r.isColorActive() {
		prev := tuikit.ColorEnabled()
		tuikit.SetColorEnabled(false)
		_ = box.Render(w)
		tuikit.SetColorEnabled(prev)
	} else {
		_ = box.Render(w)
	}
```

**After**:
```go
	if !r.isColorActive() {
		var buf strings.Builder
		_ = box.Render(&buf)
		_, _ = io.WriteString(w, tuikit.StripANSI(buf.String()))
	} else {
		_ = box.Render(w)
	}
```

**In `renderTable` (Lines 214–223)**:
**Before**:
```go
	if !r.isColorActive() {
		prev := tuikit.ColorEnabled()
		tuikit.SetColorEnabled(false)
		_ = tbl.Render(w)
		tuikit.SetColorEnabled(prev)
	} else {
		_ = tbl.Render(w)
	}
```

**After**:
```go
	if !r.isColorActive() {
		var buf strings.Builder
		_ = tbl.Render(&buf)
		_, _ = io.WriteString(w, tuikit.StripANSI(buf.String()))
	} else {
		_ = tbl.Render(w)
	}
```

---

#### 3.2 `internal/workspace/doctor.go` (Line 324)
- **Target File**: `d:/CodingProjects/vortex/internal/workspace/doctor.go`
- **Lines**: 322–327
- **Action**: Replace `len(c.Name)` with `tui.VisibleWidth(c.Name)` to ensure proper visual column alignment.

**Before**:
```go
	maxNameWidth := 0
	for _, c := range rep.Checks {
		if len(c.Name) > maxNameWidth {
			maxNameWidth = len(c.Name)
		}
	}
```

**After**:
```go
	maxNameWidth := 0
	for _, c := range rep.Checks {
		if w := tui.VisibleWidth(c.Name); w > maxNameWidth {
			maxNameWidth = w
		}
	}
```

---

## 5. Verification Method

Once the worker has applied all edits, the implementation must be verified via the following sequence:

### 1. Adversarial Test Verification
```powershell
$env:GOWORK="off"
go test -v ./cmd/vortex -run TestMilestone2_Adversarial
```
*Expected*: All 10 adversarial tests PASS cleanly, specifically:
- `TestMilestone2_Adversarial_NoInformalEmojisInCodebase` PASS
- `TestMilestone2_Adversarial_SpecImport_NoInformalEmojis` PASS
- `TestMilestone2_Adversarial_TupleRenderTable_NoInformalEmojis` PASS

### 2. Full Test Suite Verification
```powershell
$env:GOWORK="off"
go test -count=1 ./...
```
*Expected*: All package test suites across the workspace pass with exit code 0.

### 3. Linter Compliance Verification
```powershell
$env:GOWORK="off"
golangci-lint run ./...
```
*Expected*: 0 issues found, exit code 0.

### 4. Emoji Cleanliness Verification
```powershell
git grep -n -E "⚡|✨|🔴|🟡|🔵|❌|⚠️|🚀|🤖" -- "*.go"
```
*Expected*: Matches only inside `cmd/vortex/adversarial_m2_test.go` assertion definitions. Zero matches in production code.

### 5. Invalidation Conditions
- If any test in `cmd/vortex/app_test.go` fails, verify that the arrow replacements `65536 ↳ 8192` and `Field4 ↳ MaxTokens` were correctly applied to `app_test.go`.
- If `golangci-lint` reports a `gci` error, check that an empty newline exists between `github.com/lemon4ksan/foundation/tuikit` and any `github.com/lemon4ksan/vortex/...` import.
- If `golangci-lint` reports a `golines` error, verify that line lengths in `pkg/lint/format.go` and `internal/core/autopilot.go` are under 120 columns.
