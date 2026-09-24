# Milestone 2 Adversarial Verification Report: Glyph & Aesthetic Discipline

**Target**: Milestone 2 (Restrained High-Craft CLI Presentation via foundation/tuikit)  
**Agent**: `challenger_m2_2_rep` (EMPIRICAL CHALLENGER — critic, specialist)  
**Working Directory**: `d:/CodingProjects/vortex/.agents/challenger_m2_2_rep/`  
**Verdict**: **REQUEST_CHANGES**

---

## 1. Observation

### 1.1 Direct Leakage of Forbidden Informal Emoji `⚡` in User-Facing CLI Output
In `pkg/openapi/reconcile.go` line 59:
```go
// Render formats a human-readable terminal report of the merge.
func (s MergeSummary) Render(targetPath string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "⚡ [vortex merge] Merging %q into %q\n", s.SpecSource, targetPath)
```
This method is invoked directly in `internal/spec/import.go` lines 358 and 528:
```go
358:			fmt.Fprint(stdout, summary.Render(opts.targetOut))
...
528:	fmt.Fprint(stdout, summary.Render(opts.targetOut))
```
During CLI command execution `vortex spec import`, the terminal prints `⚡ [vortex merge] Merging ...` to `os.Stdout`.

### 1.2 Leakage of Forbidden Informal Emoji `⚡` in Terminal Table Report
In `pkg/tuple/analyzer.go` line 208:
```go
// RenderTable renders a neutral, clean terminal table of the tuple index analysis.
func (r *TupleAnalysisReport) RenderTable() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "⚡ Vortex Tuple Saliency Analysis (%s across %d sample entries)\n\n",
		r.StructName, r.TotalSamples)
```
The report header contains the forbidden informal emoji `⚡` and unstandardized manual divider lines instead of `foundation/tuikit`.

### 1.3 Informal Emoji `🤖` in Generated Sidecar Runtime
In `pkg/oracle/gen/js_emitter.go` line 664:
```javascript
server.listen(PORT, '127.0.0.1', () => {
  console.log(`🤖 Vortex Universal Oracle [${PORT}] listening on http://127.0.0.1:${PORT}`);
```
Contains the informal `🤖` robot emoji.

### 1.4 Non-Sovereign Dingbat Arrows in CLI Diff Output
In `pkg/diff/stack.go` lines 871 and 882:
```go
871:	fmt.Fprintf(&buf, "  • %s [Tag #%s]: %s ➔ %s (%s)\n", tr.StructName, tr.Tag, tr.OldField, tr.NewField, tr.GoType)
...
882:	fmt.Fprintf(&buf, "  • %s: %s ➔ %s\n", mr.Route, mr.OldMethod, mr.NewMethod)
```
In `internal/traffic/diff.go` line 683:
```go
683:	fmt.Fprintf(stdout, "  • %s%s: %s ➔ %s      ➜ vortex ast rename --type=%s --field=%s --to=<NAME>\n",
```
These use non-sovereign dingbat arrows `➔` (U+2794) and `➜` (U+279C) instead of the prescribed sovereign arrow `↳` (U+21B3) or `—`.

### 1.5 Empirical Reproduction via Adversarial Test Harness
We added empirical reproduction tests to `cmd/vortex/adversarial_m2_test.go` and executed `go test -v ./cmd/vortex -run TestMilestone2`:
```text
=== RUN   TestMilestone2_Adversarial_NoColor_TuikitStyles
--- PASS: TestMilestone2_Adversarial_NoColor_TuikitStyles (0.00s)
=== RUN   TestMilestone2_Adversarial_LintFormatReport_NonTTY_And_NoColor
--- PASS: TestMilestone2_Adversarial_LintFormatReport_NonTTY_And_NoColor (0.00s)
=== RUN   TestMilestone2_Adversarial_ProjectStatus_NonTTY_And_NoColor
--- PASS: TestMilestone2_Adversarial_ProjectStatus_NonTTY_And_NoColor (0.00s)
=== RUN   TestMilestone2_Adversarial_TerminalRenderer_NoColor
--- PASS: TestMilestone2_Adversarial_TerminalRenderer_NoColor (0.00s)
=== RUN   TestMilestone2_Adversarial_CLIApp_PipedStdout_NoANSI
--- PASS: TestMilestone2_Adversarial_CLIApp_PipedStdout_NoANSI (0.00s)
=== RUN   TestMilestone2_Adversarial_AllSubcommandsHelp_NoInformalEmojis
--- PASS: TestMilestone2_Adversarial_AllSubcommandsHelp_NoInformalEmojis (0.03s)
=== RUN   TestMilestone2_Adversarial_NoInformalEmojisInCodebase
    adversarial_m2_test.go:334: Forbidden informal emojis discovered in codebase production files:
        pkg/openapi/reconcile.go:59 contains ⚡: fmt.Fprintf(&sb, "⚡ [vortex merge] Merging %q into %q\n", s.SpecSource, targetPath)
        pkg/tuple/analyzer.go:208 contains ⚡: fmt.Fprintf(&sb, "⚡ Vortex Tuple Saliency Analysis (%s across %d sample entries)\n\n",
--- FAIL: TestMilestone2_Adversarial_NoInformalEmojisInCodebase (0.07s)
=== RUN   TestMilestone2_Adversarial_SpecImport_NoInformalEmojis
    adversarial_m2_test.go:383: CLI 'spec import' stdout leaked forbidden informal emoji "⚡" in output:
        ⚡ [vortex merge] Merging "...\spec.json" into "...\api.go"
          Upstream Spec Version: v1.5.0
        
          [~] 1 endpoint(s) updated (custom types & directives preserved):
              • GetItem
        
        ✔ Successfully reconciled Go contract AST.
--- FAIL: TestMilestone2_Adversarial_SpecImport_NoInformalEmojis (0.08s)
=== RUN   TestMilestone2_Adversarial_TupleRenderTable_NoInformalEmojis
    adversarial_m2_test.go:400: TupleAnalysisReport.RenderTable leaked forbidden informal emoji "⚡":
        ⚡ Vortex Tuple Saliency Analysis (TestTuple across 10 sample entries)
        
          INDEX  OCCUPANCY    TYPE         STATUS / NAME            SAMPLE VALUES FROM TRAFFIC
          ─────────────────────────────────────────────────────────────────────────────────────────────
          [ 0]   100% (10)    string       Field0                   <always nil / unused>
        
        Tip: Use `vortex ast rename --type=TestTuple --field=<Index> --to=<Name>` to assign semantic names.
--- FAIL: TestMilestone2_Adversarial_TupleRenderTable_NoInformalEmojis (0.00s)
=== RUN   TestMilestone2_Adversarial_TuikitTable_AlignmentAndBorders
--- PASS: TestMilestone2_Adversarial_TuikitTable_AlignmentAndBorders (0.00s)
=== RUN   TestMilestone2_Adversarial_TuikitBox_FramingAndCorners
--- PASS: TestMilestone2_Adversarial_TuikitBox_FramingAndCorners (0.00s)
FAIL
FAIL	github.com/lemon4ksan/vortex/cmd/vortex	0.544s
FAIL
```

### 1.6 Verification of `tuikit.Table` and `tuikit.Box`
Both `TestMilestone2_Adversarial_TuikitTable_AlignmentAndBorders` and `TestMilestone2_Adversarial_TuikitBox_FramingAndCorners` passed:
- `tuikit.Table` correctly computes column widths using `VisibleWidth` (stripping ANSI codes and counting runes), properly aligning headers, horizontal dividers, and multi-byte sovereign runes (`✔`, `✖`, `▲`, `—`, `↳`).
- `tuikit.Box` guarantees straight rectangular right borders across `BorderSingle`, `BorderRounded`, and `BorderHeavy` styles under both interactive terminal styling and pipe redirection.
- Non-TTY / pipe execution produces clean plaintext with zero ANSI escape codes across all CLI subcommands (`TestMilestone2_Adversarial_CLIApp_PipedStdout_NoANSI`).

### 1.7 Direct Linter Failures from Milestone 2 Modifications
Execution of `golangci-lint run ./...` directly failed with 3 issues in files touched during Milestone 2:
```text
internal\perf\prof.go:26:1: File is not properly formatted (gci)
	"github.com/lemon4ksan/foundation/tuikit"
^
pkg\lint\format.go:18:1: File is not properly formatted (gci)
	"github.com/lemon4ksan/foundation/tuikit"
^
internal\core\autopilot.go:533:1: File is not properly formatted (golines)
		fmt.Sprintf("Workspace is 100%% healthy, synchronized, and compiled [%v | 0 allocs]!", elapsed.Round(time.Millisecond)),
^
3 issues:
* gci: 2
* golines: 1
```
This violates Acceptance Criteria: `golangci-lint run reports zero lint violations.`

---

## 2. Logic Chain

1. **Premise 1 (Requirement R2 / F-R2-05)**: ORIGINAL_REQUEST.md explicitly mandates:
   > "Strictly avoid emoji spam or 'AI-style' decorations; use clean, understated Unicode glyphs (✔, ✖, ◆, ↳, —)..."
   And the task dispatch explicitly mandates:
   > "1. Empirically scan the codebase and test CLI output for informal emojis (`⚡`, `✨`, `🔴`, `🟡`, `🔵`, `❌`, `⚠️`, `🚀`). Verify none remain in user-facing CLI output."
2. **Observation 1.1**: Direct invocation of `vortex spec import` prints `⚡ [vortex merge] Merging %q into %q\n` to `stdout` from `pkg/openapi/reconcile.go:59`.
3. **Inference 1**: `vortex spec import` violates the core acceptance criteria of Milestone 2 because informal emoji `⚡` is emitted directly into user-facing CLI output.
4. **Observation 1.2 & 1.3**: `pkg/tuple/analyzer.go:208` emits `⚡` and `pkg/oracle/gen/js_emitter.go:664` emits `🤖`.
5. **Observation 1.4**: `pkg/diff/stack.go` and `internal/traffic/diff.go` retain informal dingbat arrows `➔` and `➜` rather than the prescribed sovereign glyph `↳`.
6. **Observation 1.5**: The empirical test harness in `cmd/vortex/adversarial_m2_test.go` directly fails with code 1, reproducing the emoji leakages in CLI stdout and in report strings.
7. **Observation 1.7**: `golangci-lint run ./...` directly fails with 3 issues (`gci` import formatting in `internal/perf/prof.go` and `pkg/lint/format.go`, `golines` line length in `internal/core/autopilot.go`).
8. **Conclusion**: Milestone 2 cannot be approved in its current state. Changes are required to eliminate the remaining informal emojis (`⚡`, `🤖`), replace non-sovereign arrows with `↳`, and resolve the 3 lint issues.

---

## 3. Caveats

- Implementation code was strictly not modified by this agent in accordance with the review-only constraint.
- `tuikit` component primitives (`tuikit.Table`, `tuikit.Box`, `tuikit.Badge`, `tuikit.VisibleWidth`) were verified to function with mathematical correctness and zero ANSI leakage under pipe mode.
- Existing tests across all packages in `go test ./...` pass, confirming that the regression is aesthetic/glyph compliance and linting formatting rather than functional logic breakage.

---

## 4. Conclusion

**Verdict: REQUEST_CHANGES**

The following concrete fixes must be applied by the worker agent:
1. **`pkg/openapi/reconcile.go:59`**:
   Replace:
   ```go
   fmt.Fprintf(&sb, "⚡ [vortex merge] Merging %q into %q\n", s.SpecSource, targetPath)
   ```
   With:
   ```go
   fmt.Fprintf(&sb, "◆ [vortex merge] Merging %q into %q\n", s.SpecSource, targetPath)
   ```
2. **`pkg/tuple/analyzer.go:208`**:
   Replace:
   ```go
   fmt.Fprintf(&sb, "⚡ Vortex Tuple Saliency Analysis (%s across %d sample entries)\n\n",
   ```
   With:
   ```go
   fmt.Fprintf(&sb, "◆ Vortex Tuple Saliency Analysis (%s across %d sample entries)\n\n",
   ```
3. **`pkg/oracle/gen/js_emitter.go:664`**:
   Replace:
   ```javascript
   console.log(`🤖 Vortex Universal Oracle [%q] listening on http://127.0.0.1:${PORT}`);
   ```
   With:
   ```javascript
   console.log(`◆ Vortex Universal Oracle [%q] listening on http://127.0.0.1:${PORT}`);
   ```
4. **`pkg/diff/stack.go:871, 882` & `internal/traffic/diff.go:683`**:
   Replace `➔` and `➜` with `↳` or `->`.
5. **Linting Formatting (`golangci-lint run ./...`)**:
   - `internal/perf/prof.go`: Fix `gci` import grouping for `"github.com/lemon4ksan/foundation/tuikit"`.
   - `pkg/lint/format.go`: Fix `gci` import grouping for `"github.com/lemon4ksan/foundation/tuikit"`.
   - `internal/core/autopilot.go:533`: Split long line to satisfy `golines`.

---

## 5. Verification Method

To independently verify these findings:
1. Run the empirical adversarial test suite:
   ```powershell
   go test -v ./cmd/vortex -run TestMilestone2
   ```
   Expected prior to fix: `TestMilestone2_Adversarial_NoInformalEmojisInCodebase`, `TestMilestone2_Adversarial_SpecImport_NoInformalEmojis`, and `TestMilestone2_Adversarial_TupleRenderTable_NoInformalEmojis` FAIL with verbatim `⚡` error output.
   Expected after fix: All tests in `TestMilestone2*` PASS cleanly.
2. Run workspace lint suite:
   ```powershell
   golangci-lint run ./...
   ```
   Expected prior to fix: 3 issues (`gci` in `prof.go`, `format.go`, `golines` in `autopilot.go`).
   Expected after fix: 0 issues.
3. Run full workspace tests:
   ```powershell
   go test ./...
   ```
4. Run codebase-wide emoji scan:
   ```powershell
   git grep -nE "[⚡✨🔴🟡🔵❌⚠️🚀🤖]"
   ```
   Must return zero lines outside of historical markdown logs.

