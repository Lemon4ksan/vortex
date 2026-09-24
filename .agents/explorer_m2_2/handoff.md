# Milestone 2 Investigation Report: Restrained High-Craft CLI Presentation (`pkg/lint/format.go` & `pkg/project/status.go`)

- **Author**: `explorer_m2_2`
- **Working Directory**: `d:/CodingProjects/vortex/.agents/explorer_m2_2/`
- **Milestone**: Milestone 2 (Restrained High-Craft CLI Presentation via foundation/tuikit)
- **Status**: Complete (Hard Handoff)
- **Timestamp**: `2026-09-22T15:09:00Z`

---

## 1. Observation

### 1.1 Raw ANSI Escape Sequences Inventory

#### A. `pkg/lint/format.go` (8 raw ANSI escapes)
Lines 20–29 define 8 raw ANSI string constants:
```go
// ANSI color codes
const (
	ansiReset  = "\033[0m"   // Line 21
	ansiBold   = "\033[1m"   // Line 22
	ansiDim    = "\033[2m"   // Line 23
	ansiRed    = "\033[31m"  // Line 24
	ansiGreen  = "\033[32m"  // Line 25
	ansiYellow = "\033[33m"  // Line 26
	ansiBlue   = "\033[34m"  // Line 27
	ansiCyan   = "\033[36m"  // Line 28
)
```
These constants are referenced throughout `FormatReport` and `printDiagnostic`:
- Line 37: `fmt.Fprintf(w, "%s%s⚡ Vortex Contract Inspector%s\n", ansiBold, ansiCyan, ansiReset)`
- Line 38–39: `fmt.Fprintf(w, "%sTarget: %s ...%s\n\n", ansiDim, ..., ansiReset)`
- Line 45–52: `fmt.Fprintf(w, "%s%s✔ All contracts are valid and synchronized!%s %s(%d warnings suppressed ...)%s\n", ansiBold, ansiGreen, ansiReset, ansiDim, ..., ansiReset)`
- Line 54: `fmt.Fprintf(w, "%s%s✔ All contracts are valid and synchronized!%s\n", ansiBold, ansiGreen, ansiReset)`
- Line 79: `fmt.Fprintf(w, "%s%s◆ Errors (%d)%s\n", ansiBold, ansiRed, len(errs), ansiReset)`
- Line 82: `printDiagnostic(w, d, ansiRed)`
- Line 89: `fmt.Fprintf(w, "%s%s◆ Warnings & Suggestions (%d)%s\n", ansiBold, ansiYellow, len(warns), ansiReset)`
- Line 92: `printDiagnostic(w, d, ansiYellow)`
- Line 99: `fmt.Fprintf(w, "%s%s◆ Info (%d)%s\n", ansiBold, ansiBlue, len(infos), ansiReset)`
- Line 102: `printDiagnostic(w, d, ansiBlue)`
- Line 109: `fmt.Fprintf(w, "%sSummary:%s ", ansiBold, ansiReset)`
- Line 113: `parts = append(parts, fmt.Sprintf("%s%d error(s)%s", ansiRed, report.Errors(), ansiReset))`
- Line 117: `parts = append(parts, fmt.Sprintf("%s%d warning(s)%s", ansiYellow, report.Warnings(), ansiReset))`
- Line 121: `parts = append(parts, fmt.Sprintf("%s%d auto-fixable%s", ansiGreen, report.FixableCount(), ansiReset))`
- Line 125: `parts = append(parts, fmt.Sprintf("%s%d suppressed%s", ansiDim, report.SuppressedCount, ansiReset))`
- Line 159–161: `fmt.Fprintf(w, "\n%sRun `vortex check --fix` ...%s\n", ansiCyan, ..., ansiReset)`
- Line 193: `fmt.Fprintf(w, "  ↳ %s[%s:%s]%s %s%s%s\n", color, d.RuleID, d.RuleName, ansiReset, ansiBold, loc, ansiReset)`
- Line 197: `fmt.Fprintf(w, "    %s↳ Suggestion:%s %s\n", ansiCyan, ansiReset, d.Suggestion)`
- Line 201: `fmt.Fprintf(w, "    %s↳ To suppress:%s //vortex:ignore %s\n", ansiDim, ansiReset, d.RuleName)`

#### B. `pkg/project/status.go` (7 raw ANSI escapes & helper functions)
Lines 432–446 define 7 raw ANSI escapes inside helper functions:
```go
func ansi(color bool, code, text string) string {
	if !color || text == "" {
		return text
	}

	return code + text + "\033[0m" // Line 437: Escape 1 (\033[0m)
}

func ansiBold(color bool, text string) string   { return ansi(color, "\033[1m", text) }  // Line 440: Escape 2 (\033[1m)
func ansiDim(color bool, text string) string    { return ansi(color, "\033[2m", text) }  // Line 441: Escape 3 (\033[2m)
func ansiGreen(color bool, text string) string  { return ansi(color, "\033[32m", text) } // Line 442: Escape 4 (\033[32m)
func ansiYellow(color bool, text string) string { return ansi(color, "\033[33m", text) } // Line 443: Escape 5 (\033[33m)
func ansiRed(color bool, text string) string    { return ansi(color, "\033[31m", text) }    // Line 444: Escape 6 (\033[31m)
func ansiCyan(color bool, text string) string   { return ansi(color, "\033[36m", text) }   // Line 445: Escape 7 (\033[36m)
```
These helpers are used 22 times across lines 97–423 in `StatusReport.Render(color bool)`.

---

### 1.2 Emoji and Presentation Clutter Inventory

| File | Line | Current Content | Clutter Type | Sovereign Replacement |
|---|---|---|---|---|
| `pkg/lint/format.go` | 37 | `⚡ Vortex Contract Inspector` | Emoji `⚡` | `tuikit.RenderHeader("◆ Vortex Contract Inspector")` |
| `pkg/project/status.go` | 97 | `⚡ Vortex API Guardian` | Emoji `⚡` | `tuikit.Bold(tuikit.Cyan("◆ Vortex API Guardian"))` or `tuikit.RenderHeader("◆ Vortex API Guardian")` |
| `pkg/project/status.go` | 157, 333, 342 | `⚠` | Warning symbol | `▲` or `tuikit.Badge("▲ STALE", tuikit.Yellow)` |
| `pkg/project/status.go` | 267 | `🔴 %s  %s\n` | Emoji `🔴` | `tuikit.Badge("✖ BREAKING", tuikit.Red)` |
| `pkg/project/status.go` | 279 | `🟡 %s  %s\n` | Emoji `🟡` | `tuikit.Badge("▲ DRIFT", tuikit.Yellow)` |
| `pkg/project/status.go` | 291 | `✔ %s  %s\n` | Plain check | `tuikit.Badge("✔ IN SYNC", tuikit.Green)` |
| `pkg/project/status.go` | 402 | `🔵 %s  by %s ...` | Emoji `🔵` | `↳` (styled via `tuikit.Cyan`) |
| `pkg/project/status.go` | 421, 423 | `✨ All systems nominal...` | Emoji `✨` | `✔ All systems nominal. Network layer is 100% synchronized.` |

---

### 1.3 Foundation Tuikit Primitives & Contracts

Inspected `foundation/tuikit` in sibling module `github.com/lemon4ksan/foundation/tuikit`:
1. `tuikit.RenderHeader(title string) string`: returns `Bold(title)`.
2. `tuikit.Table`:
   - Constructor: `NewTable(headers ...string) *Table`
   - Config: `SetAlignment(col, align)`, `SetIndent(spaces)`, `SetMinWidth(col, width)`
   - Data: `AddRow(cells ...string)`
   - Rendering: `Render(w io.Writer) error`, `String() string`
   - Visual width: Uses `VisibleWidth(cell)` and `StripANSI(cell)` for exact column alignment regardless of ANSI escapes and Unicode run lengths.
   - Headers behavior: When `len(t.Headers) == 0`, `Table.Render` omits the header line and divider rule, rendering rows with alignment and indent only.
3. `tuikit.Badge(text string, colorFn func(string) string) string`: Returns `colorFn(text)` or `text` if `colorFn == nil`.
4. `tuikit.VisibleWidth(s string) int`: Computes `utf8.RuneCountInString(StripANSI(s))`.
5. `tuikit.StripANSI(s string) string`: Strips `\x1b\[[0-9;]*[a-zA-Z]`.
6. `tuikit.ColorEnabled() bool`: Reports whether color is enabled globally (initialized from `NO_COLOR` and `TERM=dumb`).
7. `tuikit.ProbeTerminal(w any) bool`: Checks if `w` is connected to a console/TTY (`probe_windows.go` / `probe_linux.go`).
8. `tuikit.IsInteractive(w io.Writer) bool`: Evaluates `ProbeTerminal(w) && ColorEnabled()`.

---

### 1.4 Test Assertions Audit in `pkg/lint` and `pkg/project`

#### A. `pkg/lint/lint_test.go`
Examined `TestEngine_FormatReport` at line 211:
```go
228: 	var buf bytes.Buffer
229: 	lint.FormatReport(&buf, "api.go", report)
230: 	output := buf.String()
231: 
232: 	require.Contains(t, output, "Vortex Contract Inspector")
233: 	require.Contains(t, output, "param-lifting")
234: 	require.Contains(t, output, "* W001 (param-lifting): 1")
```
- Line 232: `"Vortex Contract Inspector"` — **Passes** with `tuikit.RenderHeader("◆ Vortex Contract Inspector")`.
- Line 233: `"param-lifting"` — **Passes** (rule name present in diagnostic and table).
- Line 234: `"* W001 (param-lifting): 1"` — **FAILS** after refactoring. The rule summary format changes from bulleted text `* <rule>: <count>` to a structured `tuikit.Table` (`RULE | SEVERITY | COUNT`).

#### B. `pkg/project/project_test.go`
Examined `TestProject_StatusEngine` (line 91) and `TestProject_StatusRenderAlignment` (line 316):
```go
351: 	// Check headers and summary
352: 	require.Contains(t, rendered, "⚡ Vortex API Guardian")
353: 	require.Contains(t, rendered, "Workspace: D:/CodingProjects/g-man (4 services, 268 methods)")
354: 	require.Contains(t, rendered, "● Contracts & Generated Code:")
355: 
356: 	// Check rows
357: 	require.Contains(t, rendered, "WebAPI")
358: 	require.Contains(t, rendered, "(pkg/steam/webapi)")
359: 	require.Contains(t, rendered, "169 methods, 44 DTOs")
360: 	require.Contains(t, rendered, "100% in sync")
361: 
362: 	require.Contains(t, rendered, "Notifications")
363: 	require.Contains(t, rendered, "(pkg/steam/sys/notifications)")
364: 	require.Contains(t, rendered, "10 methods")
365: 	require.Contains(t, rendered, "api.gen.go is STALE")
```
- Line 352: `"⚡ Vortex API Guardian"` — **FAILS** when emoji `⚡` is replaced with `◆`.
- Lines 353–365: All row assertions **Pass** provided the columns in `tuikit.Table` contain these exact string segments.
- `TestProject_StatusEngine:150-152`:
  - `require.Contains(t, rendered, "API")` — **Passes**.
  - `require.Contains(t, rendered, "100% in sync")` — **Passes**.
  - `require.Contains(t, rendered, "All systems nominal")` — **Passes** (`✨` decontaminated to `✔`).

---

## 2. Logic Chain

1. **Premise**: Sovereign CLI presentation mandates zero raw ANSI escape literals (`\033[` / `\x1b[`), zero emoji spam, pixel-perfect columnar formatting via `tuikit.Table`, and strict non-TTY/`NO_COLOR` safety.
2. **From Observation 1.1**: Exactly 8 raw ANSI constants reside in `pkg/lint/format.go` (lines 20–29) and 7 raw escapes in `pkg/project/status.go` (lines 432–446).
3. **Inference**: Deleting these blocks and routing all color styling through `tuikit` (`tuikit.Bold`, `tuikit.Dim`, `tuikit.Red`, `tuikit.Green`, `tuikit.Yellow`, `tuikit.Cyan`) eliminates 100% of raw escapes from both packages.
4. **From Observation 1.2**: Emojis `⚡`, `✨`, `🔴`, `🟡`, `🔵`, `⚠` violate the restrained aesthetic.
5. **Inference**: Replacing them with restrained Unicode glyphs (`◆`, `✔`, `▲`, `✖`, `↳`) and structured badges (`tuikit.Badge("✖ BREAKING", tuikit.Red)`, `tuikit.Badge("▲ DRIFT", tuikit.Yellow)`, `tuikit.Badge("✔ IN SYNC", tuikit.Green)`) establishes benchmark-grade visual hierarchy.
6. **From Observation 1.1 & 1.3**: `pkg/project/status.go` spends ~150 lines performing manual `len()` byte counting to align columns across 4 sections (Contracts, Upstream Drift, Polyglot Targets, Proposals). Because Unicode runes (`✔`, `▲`, `↳`) have byte size > 1, manual byte-length padding produces misaligned columns.
7. **Inference**: Replacing the manual loops with `tuikit.Table` guarantees exact alignment via `tuikit.VisibleWidth` while simplifying codebase complexity by >100 lines of boilerplate.
8. **From Observation 1.3 & 1.4**: Unit tests run against `bytes.Buffer`. `tuikit.ProbeTerminal(&buf)` returns `false`, causing `tuikit.IsInteractive(&buf)` to return `false`.
9. **Inference**: Ensuring non-interactive destinations produce clean plaintext via `tuikit.StripANSI` guarantees that piped/file-redirected CLI output contains zero escape sequences, and unit test assertions on string content succeed predictably.
10. **From Observation 1.4**: Two specific test assertions (`pkg/lint/lint_test.go:234` and `pkg/project/project_test.go:352`) directly assert on the legacy bulleted rule format and the legacy emoji title banner.
11. **Inference**: These two test assertions must be updated synchronously with the refactoring to ensure `go test ./...` remains 100% green.

---

## 3. Caveats

1. **Synchronous Test Updates Required**: The implementer must update `pkg/lint/lint_test.go:234` and `pkg/project/project_test.go:352` at the exact same time as `pkg/lint/format.go` and `pkg/project/status.go`. If not updated synchronously, `go test ./pkg/...` will fail on string mismatch.
2. **`tuikit.Table` Headers vs Headerless Tables**:
   - For `pkg/lint/format.go` rule summary: Use a table WITH headers (`"RULE"`, `"SEVERITY"`, `"COUNT"`) and right-aligned count column.
   - For `pkg/project/status.go`: The section titles (`● Contracts & Generated Code:`, `● Upstream Drift:`, etc.) already serve as section headers. Constructing `tuikit.NewTable()` without column headers allows rows to align seamlessly without introducing redundant inner grid headers or divider lines that would disrupt existing row expectations.
3. **No Upstream Changes to `foundation/tuikit` Needed**: `tuikit.Badge(text, colorFn)` provides complete flexibility to render sovereign badges like `tuikit.Badge("✖ BREAKING", tuikit.Red)` without relying on the emoji-laden `BadgeFail()` or `BadgeWarn()` convenience helpers.

---

## 4. Conclusion & Refactoring Implementation Plan

### 4.1 Refactoring Plan for `pkg/lint/format.go`

#### Step 1: Import Foundation Tuikit and Purge Raw ANSI Constants
- Import `"github.com/lemon4ksan/foundation/tuikit"`.
- Delete `const ( ansiReset = ... ansiCyan = ... )` (lines 20–29).

#### Step 2: Refactor `FormatReport`
- Detect interactivity: `interactive := tuikit.IsInteractive(w)`.
- Buffer output using `var buf bytes.Buffer`.
- Header: Replace line 37 with `fmt.Fprintln(&buf, tuikit.RenderHeader("◆ Vortex Contract Inspector"))`.
- Target Subtitle: `tuikit.Dim(fmt.Sprintf("Target: %s (%d services, %d methods across %d files)", ...))`.
- Valid State: `tuikit.Bold(tuikit.Green("✔ All contracts are valid and synchronized!"))`.
- Error/Warning/Info Section Titles:
  - `tuikit.Bold(tuikit.Red(fmt.Sprintf("◆ Errors (%d)", len(errs))))`
  - `tuikit.Bold(tuikit.Yellow(fmt.Sprintf("◆ Warnings & Suggestions (%d)", len(warns))))`
  - `tuikit.Bold(tuikit.Cyan(fmt.Sprintf("◆ Info (%d)", len(infos))))`
- Summary Counts: Style parts using `tuikit.Red`, `tuikit.Yellow`, `tuikit.Green`, and `tuikit.Dim`.
- Rule Summary Table: Replace lines 130–157 with `tuikit.Table`:
  ```go
  if len(stats) > 0 {
      fmt.Fprintln(&buf)
      tbl := tuikit.NewTable("RULE", "SEVERITY", "COUNT")
      tbl.SetIndent(2)
      tbl.SetAlignment(2, tuikit.AlignRight)

      for _, s := range stats {
          ruleLabel := fmt.Sprintf("%s (%s)", s.ruleID, s.ruleName)
          sevLabel := string(s.severity)
          switch s.severity {
          case SeverityError:
              sevLabel = tuikit.Red(sevLabel)
          case SeverityWarning:
              sevLabel = tuikit.Yellow(sevLabel)
          default:
              sevLabel = tuikit.Cyan(sevLabel)
          }
          tbl.AddRow(ruleLabel, sevLabel, strconv.Itoa(s.count))
      }
      _ = tbl.Render(&buf)
  }
  ```
- Non-TTY / Redirection Safety:
  ```go
  out := buf.String()
  if !interactive {
      out = tuikit.StripANSI(out)
  }
  _, _ = io.WriteString(w, out)
  ```

#### Step 3: Refactor `printDiagnostic`
- Change signature: `func printDiagnostic(w io.Writer, d Diagnostic, colorFn func(string) string)`
- Replace ANSI formatting on line 193 with:
  ```go
  tag := fmt.Sprintf("[%s:%s]", d.RuleID, d.RuleName)
  if colorFn != nil {
      tag = colorFn(tag)
  }
  fmt.Fprintf(w, "  ↳ %s %s\n", tag, tuikit.Bold(loc))
  fmt.Fprintf(w, "    %s\n", d.Message)
  if d.Suggestion != "" && !strings.Contains(d.Suggestion, "vortex check --fix") {
      fmt.Fprintf(w, "    %s %s\n", tuikit.Cyan("↳ Suggestion:"), d.Suggestion)
  }
  if !d.Fixable() {
      fmt.Fprintf(w, "    %s //vortex:ignore %s\n", tuikit.Dim("↳ To suppress:"), d.RuleName)
  }
  ```

---

### 4.2 Refactoring Plan for `pkg/project/status.go`

#### Step 1: Import Foundation Tuikit and Purge Raw ANSI Helpers
- Import `"github.com/lemon4ksan/foundation/tuikit"`.
- Delete lines 432–446 (`ansi`, `ansiBold`, `ansiDim`, `ansiGreen`, `ansiYellow`, `ansiRed`, `ansiCyan`).

#### Step 2: Refactor `StatusReport.Render(color bool) string`
- Color calculation: `useColor := color && tuikit.ColorEnabled()`.
- Title:
  ```go
  title := "◆ Vortex API Guardian"
  if useColor {
      title = tuikit.Bold(tuikit.Cyan(title))
  } else {
      title = tuikit.RenderHeader(title)
  }
  sb.WriteString(title + "\n")
  ```
- Workspace metadata: `tuikit.Dim(fmt.Sprintf("(%d services, %d methods)", len(r.Contracts), r.TotalMethods))`.
- Section 1 (Contracts & Generated Code):
  Replace manual padding calculation (lines 110–229) with `tuikit.Table`:
  ```go
  tblContracts := tuikit.NewTable()
  tblContracts.SetIndent(2)

  for _, c := range r.Contracts {
      var iconStyled, descStyled string
      if c.IsGenStale {
          if useColor {
              iconStyled = tuikit.Badge("▲", tuikit.Yellow)
              descStyled = tuikit.Yellow(c.GenStaleReason)
          } else {
              iconStyled = "▲"
              descStyled = c.GenStaleReason
          }
      } else {
          if useColor {
              iconStyled = tuikit.Badge("✔", tuikit.Green)
              descStyled = tuikit.Green("100% in sync")
          } else {
              iconStyled = "✔"
              descStyled = "100% in sync"
          }
      }

      pathStr := "(" + filepath.ToSlash(filepath.Dir(c.File)) + ")"
      if useColor {
          pathStr = tuikit.Dim(pathStr)
      }

      methodsPart := fmt.Sprintf("%d methods", c.MethodsCount)
      if hasAnyDTOs && c.DTOsCount > 0 {
          methodsPart = fmt.Sprintf("%d methods, %d DTOs", c.MethodsCount, c.DTOsCount)
      }

      if hasAnyVersion {
          verStr := ""
          if c.Version != "" {
              verStr = fmt.Sprintf("[%s]", c.Version)
              if useColor {
                  verStr = tuikit.Cyan(verStr)
              }
          }
          tblContracts.AddRow(iconStyled, c.Name, pathStr, methodsPart, verStr, descStyled)
      } else {
          tblContracts.AddRow(iconStyled, c.Name, pathStr, methodsPart, descStyled)
      }
  }
  sb.WriteString(tblContracts.String() + "\n")
  ```
- Section 2 (Upstream Drift):
  Replace lines 243–296 with `tuikit.Table` + sovereign badges:
  ```go
  tblDrift := tuikit.NewTable()
  tblDrift.SetIndent(2)

  for _, c := range r.Contracts {
      if c.Source == "" {
          continue
      }
      var badge, desc string
      switch {
      case c.UpstreamBreakingCount > 0:
          text := fmt.Sprintf("%d BREAKING drift(s) detected with %s", c.UpstreamBreakingCount, c.Source)
          if useColor {
              badge = tuikit.Badge("✖ BREAKING", tuikit.Red)
              desc = tuikit.Red(text)
          } else {
              badge = "✖ BREAKING"
              desc = text
          }
      case c.UpstreamDriftCount > 0 || c.UpstreamGhostCount > 0:
          cnt := c.UpstreamDriftCount + c.UpstreamGhostCount
          text := fmt.Sprintf("%d non-breaking update(s) available in %s", cnt, c.Source)
          if useColor {
              badge = tuikit.Badge("▲ DRIFT", tuikit.Yellow)
              desc = tuikit.Yellow(text)
          } else {
              badge = "▲ DRIFT"
              desc = text
          }
      default:
          text := fmt.Sprintf("Up-to-date with %s (0 drift)", c.Source)
          if useColor {
              badge = tuikit.Badge("✔ IN SYNC", tuikit.Green)
              desc = tuikit.Green(text)
          } else {
              badge = "✔ IN SYNC"
              desc = text
          }
      }
      tblDrift.AddRow(badge, c.Name, desc)
  }
  sb.WriteString(tblDrift.String() + "\n")
  ```
- Section 3 (Polyglot Targets):
  Replace lines 308–365 with `tuikit.Table`:
  - Stale badge: `tuikit.Badge("▲ STALE", tuikit.Yellow)` / `"▲ STALE"`.
  - In sync badge: `tuikit.Badge("✔ IN SYNC", tuikit.Green)` / `"✔ IN SYNC"`.
- Section 4 (Proposals):
  Replace lines 368–406 with `tuikit.Table`:
  - Arrow glyph: `tuikit.Cyan("↳")` / `"↳"` (replacing `🔵`).
- Summary / Next Actions:
  - Divider: `tuikit.RenderDivider(67)`.
  - Nominal Message: `✔ All systems nominal. Network layer is 100% synchronized.\n` (replacing `✨`).
- Plaintext Enforcement:
  ```go
  result := sb.String()
  if !useColor {
      result = tuikit.StripANSI(result)
  }
  return result
  ```

---

### 4.3 Exact Test Assertion Updates

#### A. `pkg/lint/lint_test.go`
```diff
--- a/pkg/lint/lint_test.go
+++ b/pkg/lint/lint_test.go
@@ -231,5 +231,5 @@ func TestEngine_FormatReport(t *testing.T) {
 	require.Contains(t, output, "Vortex Contract Inspector")
 	require.Contains(t, output, "param-lifting")
-	require.Contains(t, output, "* W001 (param-lifting): 1")
+	require.Contains(t, output, "W001 (param-lifting)")
 }
```
**Rationale**: `FormatReport` now presents rule statistics inside a structured `tuikit.Table` rather than bullet points `* Rule: Count`.

#### B. `pkg/project/project_test.go`
```diff
--- a/pkg/project/project_test.go
+++ b/pkg/project/project_test.go
@@ -349,5 +349,5 @@ func TestProject_StatusRenderAlignment(t *testing.T) {
 	rendered := report.Render(false)
 
 	// Check headers and summary
-	require.Contains(t, rendered, "⚡ Vortex API Guardian")
+	require.Contains(t, rendered, "◆ Vortex API Guardian")
 	require.Contains(t, rendered, "Workspace: D:/CodingProjects/g-man (4 services, 268 methods)")
```
**Rationale**: Replaced emoji `⚡` with sovereign diamond `◆`.

---

## 5. Verification Method

### 5.1 Verification Commands

1. **Verify Complete Elimination of Raw ANSI in `pkg/lint` and `pkg/project`**:
   ```pwsh
   # Must return 0 matches
   git grep -n -E '\\033\[|\\x1b\[' pkg/lint/ pkg/project/
   ```

2. **Verify Emoji Decontamination**:
   ```pwsh
   # Must return 0 matches
   git grep -n -P '[\x{26A1}\x{2728}\x{1F534}\x{1F7E1}\x{1F535}\x{1F7E2}]' pkg/lint/ pkg/project/
   ```

3. **Verify NO_COLOR and Redirected Plaintext Integrity**:
   ```pwsh
   # Test NO_COLOR
   $env:NO_COLOR = "1"
   go test -v ./pkg/lint/... ./pkg/project/...
   Remove-Item Env:\NO_COLOR
   ```

4. **Verify Full Unit Test Pass**:
   ```pwsh
   go test ./pkg/lint/... ./pkg/project/...
   go test ./...
   ```

5. **Verify Golangci-Lint Zero Violations**:
   ```pwsh
   golangci-lint run ./pkg/lint/... ./pkg/project/...
   ```

### 5.2 Invalidation Conditions
- Any occurrence of `\033[` or `\x1b[` remaining in `pkg/lint/format.go` or `pkg/project/status.go`.
- Failure to update `pkg/lint/lint_test.go:234` or `pkg/project/project_test.go:352`, causing unit test regressions.
- Misaligned columns in `status.go` under non-TTY or colored modes.
