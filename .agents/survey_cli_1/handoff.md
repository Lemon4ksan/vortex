# R2 Investigation Handoff: Restrained High-Craft CLI Presentation via foundation/tuikit

**Location**: `d:/CodingProjects/vortex/.agents/survey_cli_1/handoff.md`  
**Timestamp**: `2026-09-22T14:34:00Z`  
**Author**: `survey_cli_1` (teamwork_preview_explorer)  
**Parent**: `dc717d24-c5eb-4ae0-99fc-b085ebaedd2b` (parent)  
**Status**: Complete (Hard Handoff)

---

## 1. Observation

### 1.1 Foundation Tuikit Location & Primitives Inventory
`foundation/tuikit` is defined in the local sibling module `github.com/lemon4ksan/foundation` at `D:\CodingProjects\foundation\tuikit`, wired into the workspace via `D:\CodingProjects\go.work` and referenced in `d:\CodingProjects\vortex\go.mod` (`github.com/lemon4ksan/foundation v0.0.0-20260920191713-7709c688b2d7`).

The package provides a comprehensive, zero-external-dependency terminal UI presentation suite:

| Primitive | Source File | API / Contract | Description |
|---|---|---|---|
| **`tuikit.App`** | `app.go:26` | `NewApp(...)`, `RegisterCommand(...)`, `PrintUsage(w)` | Subcommand routing, grouped help screen formatting without ASCII table grids or color noise |
| **`tuikit.Command`** | `command.go:13` | `Run(ctx, args, stdout, stderr)` | Subcommand execution interface (100% matches `internal/base.Command`) |
| **`tuikit.Box`** | `box.go:78` | `NewBox(title, width)`, `AddLine()`, `AddRow()`, `AddDivider()`, `Render(w)` | Bordered panels (`BorderSingle` ┌─┐, `BorderRounded` ╭─╮, `BorderHeavy` ┏━┓) with auto-width and straight right borders |
| **`tuikit.Table`** | `table.go:23` | `NewTable(headers...)`, `AddRow(...)`, `SetAlignment(col, align)`, `Render(w)` | Columnar formatting with `VisibleWidth` calculation, alignment (`AlignLeft`, `AlignRight`, `AlignCenter`), and divider line |
| **`tuikit.Badge`** | `badge.go:8` | `Badge(text, colorFn)`, `BadgePass()`, `BadgePassed()`, `BadgeWarn()`, `BadgeFail()`, `BadgeInfo()`, `BadgeActive()`, `BadgeSync()` | Status badge badges with ANSI color wrappers |
| **`tuikit.VisibleWidth`** | `style.go:48` | `VisibleWidth(s string) int`, `StripANSI(s string) string` | Visual character width in monospaced terminal ignoring ANSI color escapes |
| **`tuikit.ColorEnabled`** | `style.go:37` | `ColorEnabled() bool`, `SetColorEnabled(bool)` | Global color toggle, initialized from `NO_COLOR != ""` or `TERM == "dumb"` |
| **`tuikit.ProbeTerminal`** | `probe.go:18` | `ProbeTerminal(w any) bool`, `ProbeTerminalFd(fd uintptr) bool` | Native OS terminal detection (Windows `GetConsoleMode` & mintty PTY check in `probe_windows.go`, Linux `ioctl`, BSD) |
| **`tuikit.IsInteractive`** | `probe.go:37` | `IsInteractive(w io.Writer) bool` | Returns `ProbeTerminal(w) && ColorEnabled()` |
| **`tuikit.RenderStep`** | `step.go:14` | `RenderStep(curr, total, label, status, width) string` | Pipeline stage formatting with dotted connector lines: `[1/3] Pre-flight Audit ..... ✔ PASSED` |
| **`tuikit.RenderDivider`**| `step.go:28` | `RenderDivider(width int) string` | Clean horizontal rule (`strings.Repeat("─", width)`) in gray |
| **`tuikit.RenderHeader`** | `step.go:38` | `RenderHeader(title string) string` | Emphasized bold title |
| **`tuikit.RenderBar`**    | `bar.go:13`  | `RenderBar(ratio float64, width int) string` | Proportional bar fill (`█` in cyan, `░` in gray) |
| **`tuikit.RenderSparkline`** | `bar.go:35` | `RenderSparkline(values []float64) string` | Inline trend graph from rune set `  ▂▃▄▅▆▇█` |
| **`tuikit.RenderTaxDecomposition`** | `bar.go:85` | `RenderTaxDecomposition(stages []TaxStage, barWidth int) string` | Latency / share stage breakdown table inside a Box |
| **`tuikit.Progress`**    | `progress.go:80`| `NewProgress(...)`, `Add()`, `Render()`, `Done()`, `Finish()` | Live progress bar with TTY spinner or periodic non-TTY log line (respects non-interactive redirects) |
| **`tuikit.FormatBytes`** | `progress.go:371`| `FormatBytes(bytes uint64) string` | Human-readable byte sizes (e.g. `4.2 MB`, `512 KB`) |
| **Styling Colors**      | `style.go:78`| `Bold`, `Dim`, `Italic`, `Underline`, `Red`, `Green`, `Yellow`, `Blue`, `Magenta`, `Cyan`, `Gray`, `White` | Zero-allocation styling functions that return plain text when `ColorEnabled() == false` |
| **Text Padding**        | `style.go:115`| `PadRight`, `PadLeft`, `PadCenter` | Visual width-aware text alignment |
| **Line Iterator**       | `style.go:152`| `LinesSeq(text string) iter.Seq[string]` | Zero-allocation line iterator supporting LF and CRLF |

---

### 1.2 Comprehensive Inventory of Raw ANSI Escapes in `pkg/` and `internal/`
A repo-wide search confirmed that **exactly 26 raw ANSI escape sequences** exist across **3 files** in `pkg/` and `internal/`:

#### A. `internal/text/render_terminal.go` (11 occurrences)
```go
Line 26: 	ansiReset     = "\033[0m"
Line 27: 	ansiBold      = "\033[1m"
Line 28: 	ansiDim       = "\033[2m"
Line 29: 	ansiUnderline = "\033[4m"
Line 31: 	ansiRed     = "\033[31m"
Line 32: 	ansiGreen   = "\033[32m"
Line 33: 	ansiYellow  = "\033[33m"
Line 34: 	ansiBlue    = "\033[34m"
Line 35: 	ansiMagenta = "\033[35m"
Line 36: 	ansiCyan    = "\033[36m"
Line 37: 	ansiWhite   = "\033[37m"
```

#### B. `pkg/lint/format.go` (8 occurrences)
```go
Line 21: 	ansiReset  = "\033[0m"
Line 22: 	ansiBold   = "\033[1m"
Line 23: 	ansiDim    = "\033[2m"
Line 24: 	ansiRed    = "\033[31m"
Line 25: 	ansiGreen  = "\033[32m"
Line 26: 	ansiYellow = "\033[33m"
Line 27: 	ansiBlue   = "\033[34m"
Line 28: 	ansiCyan   = "\033[36m"
```

#### C. `pkg/project/status.go` (7 occurrences)
```go
Line 437: 	return code + text + "\033[0m"
Line 440: func ansiBold(color bool, text string) string   { return ansi(color, "\033[1m", text) }
Line 441: func ansiDim(color bool, text string) string    { return ansi(color, "\033[2m", text) }
Line 442: func ansiGreen(color bool, text string) string  { return ansi(color, "\033[32m", text) }
Line 443: func ansiYellow(color bool, text string) string { return ansi(color, "\033[33m", text) }
Line 444: func ansiRed(color bool, text string) string    { return ansi(color, "\033[31m", text) }
Line 445: func ansiCyan(color bool, text string) string   { return ansi(color, "\033[36m", text) }
```
*(Note: Makefile lines 9, 10, 72 also define ANSI variables for `make help`, but Makefile is build-system metadata, not part of `pkg/` or `internal/`).*

---

### 1.3 Inventory of CLI Presentation Components Requiring Modernization
Currently, CLI presentation across subcommands suffers from fragmented rendering, manual table alignment using byte counts, and ad-hoc string formatting:

| Component / Subcommand | File Path | Current Presentation Flaws | Required Modernization |
|---|---|---|---|
| **Linter Diagnostics** | `pkg/lint/format.go` | Hardcoded ANSI escapes; `⚡` banner; lacks `NO_COLOR` and non-TTY checks; ad-hoc summary list (`* Rule: count`) | Eliminate raw escapes; replace with `tuikit.Bold/Red/Yellow/Cyan`; render rule statistics using `tuikit.Table`; use `tuikit.RenderHeader` |
| **Workspace Guardian** | `pkg/project/status.go` | Hardcoded ANSI escapes; `⚡` banner; emoji status (`🔴`, `🟡`, `🔵`, `✨`); manual column width calculation with `len()` byte count | Use `tuikit.Table` for Contracts, Upstream Drift, Polyglot targets; replace emojis with restrained Unicode (`✔`, `▲`, `✖`, `↳`); use `tuikit.IsInteractive(w)` |
| **Semantic Document Renderer** | `internal/text/render_terminal.go` | 11 raw ANSI constants; manual table width using `len(cell)` causing rune misalignment; callouts rendered with crude block characters (`▍`) | Replace ANSI constants with `tuikit` styling functions; rewrite `renderTable` to use `tuikit.Table`; rewrite `renderCallout` using `tuikit.Box` with `BorderRounded` or `BorderSingle`; respect `tuikit.IsInteractive(w)` |
| **Intent Semantic Icons** | `internal/text/intent.go` | Returns emoji spam: `ℹ️`, `✅`, `⚠️`, `❌` | Replace with sovereign Unicode: `ℹ`, `✔`, `▲`, `✖`, `—` |
| **CLI App Engine** | `cmd/vortex/app.go` | Duplicates `tuikit.App` and `commandGroup` logic; does not probe stdout TTY or configure `tuikit.SetColorEnabled` at process entry | Forward global color configuration via `tuikit.ProbeTerminal(stdout)`; delegate usage screen or align styling with `tuikit.App` |
| **Autopilot Pipeline** | `internal/core/autopilot.go` | `⚡` title banner; `❌` error prefixes; `📊` section emoji; duplicates byte formatting (`formatByteSize`) | Use `tuikit.RenderStep` for all 3 stages; use `tuikit.FormatBytes`; replace `❌` with `✖`; format elapsed summary cleanly as `[12.4ms \| 0 allocs]` |
| **Doctor Diagnostics** | `internal/workspace/doctor.go` | `⚡` banner; `❌`, `⚠️`, `✨` conclusion emojis; manual width formatting for check rows | Use `tuikit.Table` for check list; replace conclusion emojis with clean Unicode status; use `tuikit.Badge("✖ FAIL", tuikit.Red)` |
| **Performance Profiler** | `internal/perf/prof.go` | `⚡`, `📊`, `🔬`, `⏱️` emojis in document nodes; manually constructed tax bar table | Replace section emojis with sovereign text; delegate latency tax decomposition directly to `tuikit.RenderTaxDecomposition`; format metrics as `[1.2ms \| 0 allocs]` |
| **Silicon Benchmarker** | `internal/perf/bench.go` | Manual ASCII dividers `===...===`; `⚡ Silicon Score:` emoji banner; manual `[1/7]` step strings | Use `tuikit.RenderStep` for 7 benchmark phases; use `tuikit.RenderDivider`; use `tuikit.Box` or `tuikit.Table` for results score card |
| **Tuple Saliency Analyzer** | `pkg/tuple/analyzer.go` | `⚡` title banner; ad-hoc string padding table with manual `─` underline | Modernize `RenderTable()` using `tuikit.Table` with column alignments |
| **Workspace Orchestrator** | `internal/workspace/work.go` | `⚡` title banner; `⚠️ Missing config`; `❌` / `✨` pipeline finish emoji | Modernize workspace multi-repo table using `tuikit.Table`; use `tuikit.RenderStep`; replace emojis with `✔`, `✖`, `▲` |
| **AST Subcommands (blame, history, log, tag)** | `internal/ast/*.go` | `⚡` banners in history, blame, log, tag; manual `strings.Repeat("─", ...)` table dividers; custom string truncation | Modernize tables via `tuikit.Table`; replace banners with `tuikit.RenderHeader` |
| **Spec & Upstream Sources** | `internal/spec/source.go` | `⚡` banner; `❌` errors across schema fetch failures; ad-hoc table formatting | Modernize list table with `tuikit.Table`; replace `❌` error indicators with `✖` |
| **Semantic Diff Stack** | `pkg/diff/stack.go` | `📦`, `⚡`, `➕`, `➖`, `📄` emoji headers | Replace with clean Unicode `◆ Tuple Field Renames`, `◆ RPC Renames`, `+ Added`, `- Removed` |

---

### 1.4 Comprehensive Emoji Decontamination Inventory
A codebase search found **37 Go files** containing emoji decorations across `cmd/`, `internal/`, and `pkg/`:

| Emoji Symbol | Occurrences & Locations | Sovereign Replacement |
|---|---|---|
| **`⚡`** (Lightning) | 47 occurrences across `cmd/vortex`, `pkg/lint`, `pkg/project`, `pkg/diff`, `pkg/tuple`, `pkg/openapi`, `internal/ast`, `internal/core`, `internal/perf`, `internal/spec`, `internal/traffic`, `internal/workspace` | `◆` (U+25C6 Black Diamond) or clean `tuikit.RenderHeader` without icon |
| **`✨`** (Sparkles) | 8 occurrences in `pkg/project/status.go`, `internal/ast/accept.go`, `internal/core/autopilot.go`, `internal/oracle/cmd.go`, `internal/workspace/clean.go`, `internal/workspace/doctor.go`, `internal/workspace/work.go` | `✔` (U+2714 Heavy Check Mark) or understated success text: `✔ Workspace synchronized` |
| **`❌`** (Cross mark) | 16 occurrences in `internal/core/autopilot.go`, `internal/spec/source.go`, `internal/workspace/doctor.go`, `internal/workspace/work.go`, `internal/text/intent.go` | `✖` (U+2716 Heavy Multiplication X) or `tuikit.Badge("✖ FAIL", tuikit.Red)` |
| **`⚠️`** / **`⚠`** (Warning) | 12 occurrences in `pkg/project/status.go`, `internal/core/autopilot.go`, `internal/workspace/doctor.go`, `internal/workspace/work.go`, `internal/text/intent.go` | `▲` (U+25B2 Black Up-Pointing Triangle) or `tuikit.Badge("▲ WARN", tuikit.Yellow)` |
| **`🔴`** / **`🟡`** / **`🔵`** / **`🟢`** (Colored circles) | 5 occurrences in `pkg/project/status.go` (breaking drift, non-breaking update, git branch proposals) | Restrained Unicode: `✖ BREAKING` (Red), `▲ DRIFT` (Yellow), `↳` / `◆` (Cyan), `✔ IN SYNC` (Green) |
| **`🚀`** / **`📦`** / **`⚙️`** / **`💡`** (Action emojis) | `internal/core/autopilot.go:172-185`, `internal/text/bench_test.go`, `internal/text/builder_test.go`, `pkg/diff/stack.go:868` | Plain numbered steps: `[1]`, `[2]`, `[3]`, `↳ Tip: ...` |
| **`📊`** / **`🔬`** / **`⏱️`** / **`🔍`** (Section emojis) | `internal/core/autopilot.go:520`, `internal/perf/prof.go:357-388`, `internal/base/reporter.go:92` | Clean section headers: `◆ Summary`, `◆ Endpoint Ledger`, `◆ Latency Tax Decomposition` |

---

### 1.5 NO_COLOR and Non-TTY Redirection Handling Analysis
1. **Existing status**:
   - `pkg/project/status.go`: Line 92 checks `os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb"`, but ignores whether output is redirected to a file or pipe. Furthermore, caller `internal/workspace/status.go:101` invokes `report.Render(true)` unconditionally.
   - `pkg/lint/format.go`: Completely ignores `NO_COLOR` and non-TTY. Emits raw ANSI escapes regardless of environment or redirection.
   - `internal/text/render_terminal.go`: `TerminalRenderer.ColorEnabled` defaults to `true`. Never inspects `NO_COLOR` or terminal descriptor.
   - `foundation/tuikit`: Implements `tuikit.ProbeTerminal(w)` and `tuikit.IsInteractive(w)`. Its `init()` automatically sets `ColorEnabled(false)` if `NO_COLOR != ""` or `TERM == "dumb"`.
2. **Required pattern**:
   - In `cmd/vortex/main.go` / `app.go`: At startup, probe `stdout` via `tuikit.ProbeTerminal(stdout)`. If `!tuikit.ProbeTerminal(stdout)`, invoke `tuikit.SetColorEnabled(false)`.
   - In all renderers: Check `tuikit.IsInteractive(w)` or rely on `tuikit.ColorEnabled()`.
   - Non-interactive / pipe output (`vortex status > file.txt`, `vortex check | grep E001`) automatically produces 100% clean, unadorned plaintext with zero ANSI escape codes and zero broken terminal control sequences.

---

### 1.6 Baseline Codebase Health Verification
- `go test ./...`: **PASS** across all packages in `cmd/`, `internal/`, and `pkg/` (0 failures).
- `golangci-lint run`: **PASS** (0 issues).

---

## 2. Logic Chain

1. **Premise**: Benchmark-grade CLI toolchains (matching aoni and foundation core) maintain zero raw ANSI escapes, strict NO_COLOR/non-TTY redirection safety, and a restrained sovereign presentation without cartoonish emoji noise.
2. **From Observation 1.2**: Exactly 26 raw ANSI escape sequences exist in the repository, isolated to `internal/text/render_terminal.go`, `pkg/lint/format.go`, and `pkg/project/status.go`.
3. **Inference**: Replacing these 26 occurrences with `foundation/tuikit` color functions (`tuikit.Bold`, `tuikit.Dim`, `tuikit.Red`, `tuikit.Green`, `tuikit.Yellow`, `tuikit.Cyan`, `tuikit.Gray`, `tuikit.White`) completely eliminates all raw escape strings from `pkg/` and `internal/`.
4. **From Observation 1.1 & 1.3**: Existing table rendering in `internal/text/render_terminal.go`, `pkg/project/status.go`, `pkg/tuple/analyzer.go`, and `internal/ast/*.go` relies on standard `len(s)` byte counting. Because UTF-8 glyphs and ANSI codes have byte lengths > 1, table columns regularly become crooked and misaligned.
5. **Inference**: Adopting `tuikit.Table` (which utilizes `tuikit.VisibleWidth` and `tuikit.StripANSI`) guarantees pixel-perfect columnar alignment under all conditions (plain, colored, UTF-8).
6. **From Observation 1.4**: Emojis like `⚡`, `✨`, `🔴`, `🟡`, `🔵`, `❌`, `⚠️`, `🚀` create visual clutter, inconsistent rendering across terminal emulators, and an informal "AI-generated" appearance.
7. **Inference**: Replacing them with restrained Unicode glyphs (`✔`, `✖`, `◆`, `▲`, `↳`, `—`) and bracketed execution statistics (`[1.2ms | 0 allocs]`) elevates Vortex to a sovereign, benchmark-grade aesthetic.
8. **From Observation 1.5**: Because `TerminalRenderer` and `pkg/lint/format.go` do not check `tuikit.IsInteractive(w)` or `tuikit.ColorEnabled()`, redirecting CLI output to CI loggers or files currently produces escape code corruption.
9. **Inference**: Wiring `tuikit.ProbeTerminal` at CLI boot and checking `tuikit.ColorEnabled()` in all renderers guarantees full compliance with the `NO_COLOR` standard and seamless pipe redirection.

---

## 3. Caveats

1. **Test String Assertions**: Several tests in `cmd/vortex/app_test.go` (e.g. lines 385, 396, 452) and `pkg/project/project_test.go` (line 352) assert on verbatim title strings like `"⚡ Vortex API Guardian"`, `"⚡ [vortex diff]"`, and `"⚡ Vortex Auto-Pilot: Audit & Build Pipeline"`. When the implementer updates these banners to the sovereign aesthetic (`◆ Vortex API Guardian` or `Vortex API Guardian`), these test assertions must be updated synchronously to prevent test regressions.
2. **`tuikit.BadgeWarn()` and `tuikit.BadgeFail()` in `foundation/tuikit`**: In `foundation/tuikit/badge.go`, `BadgeWarn()` returns `Yellow("⚠️ WARN")` and `BadgeFail()` returns `Red("❌ FAIL")`. Rather than using these two specific helpers with emoji prefixes, Vortex should either use `tuikit.Badge("▲ WARN", tuikit.Yellow)` and `tuikit.Badge("✖ FAIL", tuikit.Red)` via the general `tuikit.Badge(text, colorFn)` API, or upstream an aesthetic update to `foundation/tuikit`.
3. **ASCII Snapshot Stack Diagrams**: In `pkg/diff/stack.go:1001`, `RenderStackDiagram()` generates an ASCII tree structure (`├── [#0] base ... └── [#1] head`). This diagram is already clean ASCII and should remain unadorned.

---

## 4. Conclusion & Feature Inventory

The R2 modernization is highly tractable and architecturally clean: `foundation/tuikit` provides every required primitive, the raw ANSI escapes are confined to just 3 files, and the baseline test suite is 100% green.

### Enumerated Feature Inventory for R2 Implementation

#### [F-R2-01] Raw ANSI Escape Eradication
- Purge all 11 raw escape constants from `internal/text/render_terminal.go`.
- Purge all 8 raw escape constants from `pkg/lint/format.go`.
- Purge all 7 raw escape occurrences and custom `ansi*` helpers from `pkg/project/status.go`.
- Guarantee zero occurrences of `\033[` or `\x1b[` in `pkg/` and `internal/`.

#### [F-R2-02] Semantic Document Terminal Renderer Modernization (`internal/text`)
- Refactor `internal/text/render_terminal.go` to import and utilize `foundation/tuikit`.
- Replace `renderTable` with `tuikit.Table` leveraging column alignments and `VisibleWidth`.
- Replace `renderCallout` with `tuikit.Box` configured with `BorderRounded` or sovereign left-bar rules.
- Replace `DividerNode` rendering with `tuikit.RenderDivider(77)`.
- Update `internal/text/intent.go` to return sovereign icons: `IntentSuccess -> "✔"`, `IntentDanger -> "✖"`, `IntentWarning -> "▲"`, `IntentInfo -> "ℹ"`, `IntentMuted -> "—"`.

#### [F-R2-03] Sovereign Linter Diagnostic Presentation (`pkg/lint`)
- Modernize `pkg/lint/format.go` to use `tuikit.RenderHeader` and `tuikit.Bold/Red/Yellow/Cyan/Green`.
- Replace `printDiagnostic` ad-hoc strings with structured, indented diagnostics.
- Format the summary rule breakdown using `tuikit.Table` (`Rule`, `Count`, `Severity`).
- Ensure `FormatReport` respects `tuikit.IsInteractive(w)` and `tuikit.ColorEnabled()`.

#### [F-R2-04] Sovereign Workspace Guardian Presentation (`pkg/project`)
- Rewrite `StatusReport.Render` in `pkg/project/status.go` using `tuikit.Table` for Contracts, Upstream Drift, and Polyglot SDK targets.
- Replace colored circle emojis (`🔴`, `🟡`, `🔵`, `✨`) with sovereign indicators (`✖ BREAKING`, `▲ DRIFT`, `↳`, `✔ IN SYNC`).
- Eliminate manual column width padding algorithms.
- Enforce `tuikit.IsInteractive(w)` so piped output is automatically plain text.

#### [F-R2-05] Codebase-Wide Emoji Decontamination
- Remove all title banner emojis (`⚡`, `✨`) across `cmd/vortex`, `internal/core`, `internal/ast`, `internal/perf`, `internal/spec`, `internal/traffic`, `internal/workspace`.
- Replace error emojis (`❌`) with `✖`.
- Replace warning emojis (`⚠️`) with `▲`.
- Replace action emojis (`🚀`, `📦`, `⚙️`, `💡`) in `autopilot.go` with clean numbered brackets.
- Standardize step rendering across multi-phase workflows (`autopilot`, `work`, `bench`) via `tuikit.RenderStep`.

#### [F-R2-06] Global NO_COLOR & Non-TTY Redirection Architecture
- In `cmd/vortex/app.go`: At execution start, check `tuikit.ProbeTerminal(stdout)`. If non-interactive or `NO_COLOR` is active, call `tuikit.SetColorEnabled(false)`.
- Ensure all subcommands writing to `stdout` or `stderr` inherit this state.
- Guarantee that `vortex <command> > file.txt` produces zero ANSI sequences.

#### [F-R2-07] Benchmark & Latency Telemetry Formatting
- In `internal/perf/prof.go`: Delegate latency tax breakdown directly to `tuikit.RenderTaxDecomposition`.
- In `internal/core/autopilot.go` and `internal/perf/prof.go`: Replace ad-hoc byte formatting with `tuikit.FormatBytes`.
- Standardize execution timing summaries to bracketed sovereign statistics: `[1.2ms | 0 allocs]` or `[4.8ms | 100% zero-alloc]`.

---

## 5. Verification Method

To independently verify the survey observations and future implementation:

1. **Verify Raw ANSI Escape Absence**:
   ```pwsh
   # Must return 0 lines in pkg/ and internal/
   git grep -n -E '\\033\[|\\x1b\[' pkg/ internal/
   ```

2. **Verify Emoji Absence**:
   ```pwsh
   # Must return 0 lines in cmd/, pkg/, and internal/ (excluding markdown/docs)
   git grep -n -P '[\x{26A1}\x{2728}\x{1F534}\x{1F7E1}\x{1F535}\x{1F7E2}\x{274C}\x{2705}\x{26A0}\x{1F680}\x{1F4E6}\x{1F4A1}\x{1F389}\x{1F50D}\x{1F6E0}\x{1F525}\x{1F4CC}\x{1F4CA}\x{1F52C}\x{23F1}]' cmd/ pkg/ internal/
   ```

3. **Verify NO_COLOR & Non-TTY Redirection**:
   ```pwsh
   # Test NO_COLOR environment variable
   $env:NO_COLOR = "1"
   go run ./cmd/vortex status
   Remove-Item Env:\NO_COLOR

   # Test pipe / redirection (should contain zero ESC bytes \x1b)
   go run ./cmd/vortex status > status_test.txt
   Get-Content status_test.txt -Raw | Select-String "\x1b"
   Remove-Item status_test.txt
   ```

4. **Verify Full Test Suite & Lint Integrity**:
   ```pwsh
   go test ./...
   golangci-lint run
   ```
