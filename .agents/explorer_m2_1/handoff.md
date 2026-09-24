# Milestone 2 Investigation Report: Restrained High-Craft CLI Presentation via foundation/tuikit

**Location**: `d:/CodingProjects/vortex/.agents/explorer_m2_1/handoff.md`  
**Timestamp**: `2026-09-22T15:08:00Z`  
**Author**: `explorer_m2_1` (explorer archetype)  
**Parent**: `dc717d24-c5eb-4ae0-99fc-b085ebaedd2b` (parent)  
**Milestone**: Milestone 2 (Restrained High-Craft CLI Presentation via foundation/tuikit)  
**Status**: Complete (Hard Handoff)

---

## 1. Observation

### 1.1 `internal/text/render_terminal.go` Raw ANSI Constants & Structure
Direct inspection of `internal/text/render_terminal.go` revealed:
- Lines 25–38 define **11 raw ANSI escape constants**:
  ```go
  // ANSI Escape sequences.
  const (
  	ansiReset     = "\033[0m"
  	ansiBold      = "\033[1m"
  	ansiDim       = "\033[2m"
  	ansiUnderline = "\033[4m"

  	ansiRed     = "\033[31m"
  	ansiGreen   = "\033[32m"
  	ansiYellow  = "\033[33m"
  	ansiBlue    = "\033[34m"
  	ansiMagenta = "\033[35m"
  	ansiCyan    = "\033[36m"
  	ansiWhite   = "\033[37m"
  )
  ```
- Lines 40–46 define `r.style(code, s string) string`:
  ```go
  func (r *TerminalRenderer) style(code, s string) string {
  	if !r.ColorEnabled || code == "" {
  		return s
  	}

  	return code + s + ansiReset
  }
  ```
- Lines 62–70 style headings by concatenating raw escape strings (`ansiBold+ansiWhite+ansiUnderline`, `ansiBold+ansiCyan`, `ansiBold`).
- Line 77 styles section headers (`ansiBold+ansiYellow`).
- Lines 83–93 style fields (`ansiBold`, `ansiCyan`).
- Lines 98–101 style list items (`ansiDim`, `ansiCyan`).
- Lines 110–117 style code blocks (`ansiDim`, `ansiCyan`).
- Lines 120–123 style quotes (`ansiDim`).
- Line 126 styles dividers: `fmt.Fprintf(w, "%s\n\n", r.style(ansiDim, "────────────────────────────────────────"))`.
- Lines 145–174 (`renderCallout`) format callouts with crude block characters `▍`:
  ```go
  fmt.Fprintf(w, "%s %s\n", r.style(ansiBold+color, "▍"), r.style(ansiBold+color, title))
  if n.Body != "" {
  	for line := range bytesconv.ScanTokens(n.Body, '\n') {
  		fmt.Fprintf(w, "%s   %s\n", r.style(color, "▍"), line)
  	}
  }
  ```
- Lines 176–216 (`renderTable`) calculate column widths using standard Go byte length `len(cell)` instead of visual character width, causing crooked columns when cells contain multi-byte UTF-8 runes (e.g. `✔`, `✖`, `—`). Furthermore, `renderTable` completely ignores the column alignment array `t.Aligns`, formatting every cell with left-aligned `%-*s`.

### 1.2 `internal/text/intent.go` Semantic Glyphs
Direct inspection of `internal/text/intent.go` lines 46–61 revealed that `Intent.Icon()` returns multi-byte emoji strings:
```go
func (i Intent) Icon() string {
	switch i {
	case IntentInfo:
		return "ℹ️"
	case IntentSuccess:
		return "✅"
	case IntentWarning:
		return "⚠️"
	case IntentDanger:
		return "❌"
	case IntentMuted:
		return "•"
	default:
		return ""
	}
}
```
Repo-wide search confirmed `✅` only exists in `intent.go:51` and `ℹ️` only exists in `intent.go:49`.

### 1.3 `cmd/vortex/app.go` Boot Lifecycle & Terminal Probing
Direct inspection of `cmd/vortex/app.go` lines 99–118 revealed:
```go
func (a *App) Run(ctx context.Context, args []string) error {
	stdout := a.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	stderr := a.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}
...
```
Currently, `App.Run` neither probes `stdout` nor inspects `os.Getenv("NO_COLOR")`. Consequently:
- Subcommands and renderers continue emitting ANSI escapes even when `stdout` is redirected to a file or pipe (`vortex status > out.txt`).
- `NO_COLOR` environment variable is not enforced across all subcommands.

### 1.4 `foundation/tuikit` Capabilities & Contracts
Inspection of local package `github.com/lemon4ksan/foundation/tuikit` (in `d:/CodingProjects/foundation/tuikit`) revealed:
1. **Styling Functions** (`style.go:78–113`):
   - `tuikit.Bold(text string) string`
   - `tuikit.Dim(text string) string`
   - `tuikit.Underline(text string) string`
   - `tuikit.Red(text string) string`
   - `tuikit.Green(text string) string`
   - `tuikit.Yellow(text string) string`
   - `tuikit.Blue(text string) string`
   - `tuikit.Magenta(text string) string`
   - `tuikit.Cyan(text string) string`
   - `tuikit.Gray(text string) string`
   - `tuikit.White(text string) string`
   Each styling function checks `tuikit.ColorEnabled()`; if false, it returns `text` immediately without heap allocation or escape codes.
2. **Visual Width Calculation** (`style.go:48–51`):
   - `tuikit.VisibleWidth(s string) int`: strips ANSI escapes via `StripANSI` and calculates exact rune count via `utf8.RuneCountInString`.
3. **Table Primitive** (`table.go:23–75`):
   - `tuikit.NewTable(headers ...string) *Table`: initializes table with default indent of 2 spaces.
   - `SetAlignment(col int, align Alignment) *Table`: accepts `tuikit.AlignLeft`, `tuikit.AlignRight`, `tuikit.AlignCenter`.
   - `SetIndent(spaces int) *Table`: sets margin indentation.
   - `AddRow(cells ...string) *Table`: adds row data.
   - `Render(w io.Writer) error`: renders headers in bold, a horizontal divider in `tuikit.Gray`, and rows padded according to `VisibleWidth`.
4. **Box Primitive** (`box.go:78–127`):
   - `tuikit.NewBox(title string, innerWidth int) *Box`: constructs enclosed card.
   - `SetStyle(s BorderStyle) *Box`: supports `tuikit.BorderRounded` (`╭─╮│╰─╯`), `tuikit.BorderSingle` (`┌─┐│└─┘`), and `tuikit.BorderHeavy` (`┏━┓┃┗━┛`).
   - `AddLine(text string) *Box`: adds full-width line.
   - `AddDivider() *Box`: adds divider rule.
   - `Render(w io.Writer) error`: formats enclosed box with straight right borders based on `VisibleWidth`.
5. **Terminal Probe & NO_COLOR** (`probe.go:18–39`, `style.go:21–39`):
   - `tuikit.ProbeTerminal(w any) bool`: checks if `w` implements `Fd() uintptr` and tests console mode / mintty PTY status natively.
   - `tuikit.SetColorEnabled(enabled bool)`: globally enables/disables ANSI styling.
   - `tuikit.ColorEnabled() bool`: reports whether colors are globally active.
   - `tuikit.IsInteractive(w io.Writer) bool`: returns `ProbeTerminal(w) && ColorEnabled()`.
   - `tuikit.RenderDivider(width int) string`: returns `Gray(strings.Repeat("─", width))` (defaults to 72 width).

### 1.5 Critical Alignment Enum Inversion Observation
Comparing the alignment enum definitions in `internal/text/node.go` versus `foundation/tuikit/table.go` uncovered a vital divergence:
- In `internal/text/node.go:37-44`:
  ```go
  const (
  	AlignLeft Align = iota  // 0
  	AlignCenter             // 1
  	AlignRight              // 2
  )
  ```
- In `foundation/tuikit/table.go:16-20`:
  ```go
  const (
  	AlignLeft Alignment = iota // 0
  	AlignRight                 // 1
  	AlignCenter                // 2
  )
  ```
Notice that `AlignCenter` is index 1 in `text` but index 2 in `tuikit`, while `AlignRight` is index 2 in `text` but index 1 in `tuikit`. A naive type cast `tuikit.Alignment(align)` would silently invert center and right alignment!

### 1.6 Current Test Assertions & Workspace Health
- Verification command `go test ./internal/text ./cmd/vortex`: **PASS** (0 failures).
- In `internal/text/builder_test.go`:
  - Lines 46–61: checks markdown callout tags (`> [!TIP]`, `> [!WARNING]`, etc.).
  - Lines 65–70: checks plain text callouts (`[SUCCESS]`, `[WARNING]`, `[DANGER]`).
  - Lines 74–78: checks terminal text (`"Release Notes"`, `"Artifacts"`, `"┌── go"`).
  - None of the tests assert on raw ANSI escape characters or specific emojis.
- In `cmd/vortex/app_test.go`:
  - Line 385: `require.Contains(t, stdout.String(), "⚡ Vortex API Git History")`
  - Line 396: `require.Contains(t, stdout.String(), "⚡ [vortex diff]")`
  - Line 452: `require.Contains(t, stdout.String(), "⚡ Vortex Auto-Pilot: Audit & Build Pipeline")`
- In `pkg/project/project_test.go`:
  - Line 352: `require.Contains(t, rendered, "⚡ Vortex API Guardian")`

---

## 2. Logic Chain

1. **Premise**: Benchmark-grade sovereign CLI architecture requires zero hardcoded ANSI escape sequences, 100% clean plaintext pipe redirection when non-interactive or under `NO_COLOR`, visual-width-aware columnar tables, and understated Unicode glyphs instead of emoji spam.
2. **From Observation 1.1 & 1.4**: All 11 raw ANSI escape constants in `render_terminal.go` can be eliminated by importing `foundation/tuikit` and utilizing its styling functions (`tuikit.Bold`, `tuikit.Dim`, `tuikit.Underline`, `tuikit.Red`, `tuikit.Green`, `tuikit.Yellow`, `tuikit.Cyan`, `tuikit.Gray`, `tuikit.White`).
3. **From Observation 1.1 & 1.4**: `TerminalRenderer.style` currently wraps strings in raw escape codes. By replacing it with a functional compositor `func (r *TerminalRenderer) style(s string, fns ...func(string) string) string` that checks `r.ColorEnabled && tuikit.ColorEnabled()`, the renderer seamlessly respects both per-instance configuration and global environment settings.
4. **From Observation 1.1, 1.4 & 1.5**: Existing table rendering in `render_terminal.go` computes widths via `len(cell)`, which misaligns columns whenever runes are multi-byte, and ignores `t.Aligns`. By refactoring `renderTable` to instantiate `tuikit.NewTable`, map `text.Align` explicitly to `tuikit.Alignment` (preventing the enum inversion noted in Observation 1.5), and delegate to `tbl.Render(w)`, tables gain pixel-perfect alignment and a clean gray divider rule.
5. **From Observation 1.1 & 1.4**: `renderCallout` currently uses crude block runes `▍`. Replacing it with `tuikit.NewBox(styledTitle, 0).SetStyle(tuikit.BorderRounded).SetIndent(2)` creates an enclosed diagnostic card with rounded borders (`╭─ title ─╮`), straight right edges, and auto-computed width based on `tuikit.VisibleWidth`.
6. **From Observation 1.2**: Replacing emoji return values in `internal/text/intent.go` with restrained Unicode glyphs (`IntentSuccess` -> `✔`, `IntentDanger` -> `✖`, `IntentWarning` -> `▲`, `IntentInfo` -> `ℹ`, `IntentMuted` -> `—`) eliminates emoji noise across all documents while maintaining full compatibility with existing tests.
7. **From Observation 1.3 & 1.4**: In `cmd/vortex/app.go`, checking `!tuikit.ProbeTerminal(stdout) || os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb"` at the entry of `App.Run` and calling `tuikit.SetColorEnabled(false)` guarantees that redirected streams (pipes, files) and non-colored terminals receive 100% clean plaintext with zero ANSI control sequences.
8. **From Observation 1.6**: Existing tests in `internal/text` pass unchanged. In `cmd/vortex/app_test.go`, tests use `*bytes.Buffer`, which is not a TTY; disabling color during non-TTY runs will keep text assertions clean. When title banner emojis are decontaminated across subcommands, the 3 test assertions in `app_test.go` and 1 assertion in `project_test.go` will be updated synchronously.

---

## 3. Caveats

1. **Alignment Enum Divergence**: As discovered in Observation 1.5, `text.AlignCenter == 1` while `tuikit.AlignCenter == 2`, and `text.AlignRight == 2` while `tuikit.AlignRight == 1`. The implementer must NEVER use a direct cast `tuikit.Alignment(a)` and must instead use an explicit mapping function `toTuikitAlign(a text.Align) tuikit.Alignment`.
2. **TerminalRenderer Instance vs Global `tuikit` Toggle**: If an application instantiates `TerminalRenderer{ColorEnabled: false}`, but `tuikit.ColorEnabled()` is globally `true`, calling `tuikit.Table.Render(w)` could print ANSI bold headers unless the renderer guards the call by temporarily disabling global color or setting a plaintext mode.
3. **Banner Emoji Test Synchronicity**: The 3 test assertions in `cmd/vortex/app_test.go` (`"⚡ Vortex API Git History"`, `"⚡ [vortex diff]"`, `"⚡ Vortex Auto-Pilot: Audit & Build Pipeline"`) and 1 in `pkg/project/project_test.go` (`"⚡ Vortex API Guardian"`) must be updated synchronously when those subcommands are decontaminated.
4. **Left-Bar Rule Alternative**: While `tuikit.Box` with `BorderRounded` is recommended as the primary modern card format, an alternative sovereign left-bar rule (`│ Title\n│ Body`) is also documented below for situations where a compact, open format is desired.

---

## 4. Conclusion & Concrete Implementation Specifications

The investigation confirms that the transition to `foundation/tuikit` is fully supported by existing primitives, eliminates all 11 raw ANSI escapes in `render_terminal.go`, fixes table alignment bugs, introduces clean rounded callouts, eliminates emoji spam in `intent.go`, and guarantees NO_COLOR pipe redirection safety.

### 4.1 Implementation Specification: `internal/text/intent.go`

Replace lines 45–61 of `internal/text/intent.go` with:
```go
// Icon returns the canonical understated Unicode glyph corresponding to the intent.
func (i Intent) Icon() string {
	switch i {
	case IntentInfo:
		return "ℹ"
	case IntentSuccess:
		return "✔"
	case IntentWarning:
		return "▲"
	case IntentDanger:
		return "✖"
	case IntentMuted:
		return "—"
	default:
		return ""
	}
}
```

### 4.2 Implementation Specification: `internal/text/render_terminal.go`

#### A. Imports & ANSI Constants Eradication
Replace lines 7–47 with:
```go
package text

import (
	"fmt"
	"io"
	"strings"

	"github.com/lemon4ksan/foundation/silicon/bytesconv"
	"github.com/lemon4ksan/foundation/tuikit"
)

// TerminalRenderer converts a [Document] into ANSI colorized terminal output.
type TerminalRenderer struct {
	ColorEnabled bool
}

// NewTerminalRenderer constructs a fresh [TerminalRenderer] with ANSI colors enabled
// according to the global tuikit terminal state.
func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{ColorEnabled: tuikit.ColorEnabled()}
}

func (r *TerminalRenderer) isColorActive() bool {
	return r.ColorEnabled && tuikit.ColorEnabled()
}

func (r *TerminalRenderer) style(s string, fns ...func(string) string) string {
	if !r.isColorActive() || len(fns) == 0 || s == "" {
		return s
	}

	for _, fn := range fns {
		if fn != nil {
			s = fn(s)
		}
	}

	return s
}
```

#### B. Node Styling Updates in `Render()`
- **HeadingNode** (lines 62–69):
  ```go
  switch n.Level {
  case 1:
  	fmt.Fprintf(w, "%s\n\n", r.style(title, tuikit.Underline, tuikit.White, tuikit.Bold))
  case 2:
  	fmt.Fprintf(w, "%s\n\n", r.style(title, tuikit.Cyan, tuikit.Bold))
  default:
  	fmt.Fprintf(w, "%s\n\n", r.style(title, tuikit.Bold))
  }
  ```
- **SectionNode** (line 77):
  ```go
  fmt.Fprintf(w, "%s\n", r.style(sec+":", tuikit.Yellow, tuikit.Bold))
  ```
- **FieldNode** (lines 83–92):
  ```go
  key := r.style(n.Key, tuikit.Bold)
  val := n.Value

  switch n.Style {
  case FieldCode:
  	val = r.style("`"+n.Value+"`", tuikit.Cyan)
  case FieldBold:
  	val = r.style(n.Value, tuikit.Bold)
  }

  fmt.Fprintf(w, "  • %s: %s\n", key, val)
  ```
- **ListNode** (lines 97–101):
  ```go
  if n.Kind == ListNumbered {
  	fmt.Fprintf(w, "  %s %s\n", r.style(fmt.Sprintf("%d.", idx+1), tuikit.Dim), item)
  } else {
  	fmt.Fprintf(w, "  %s %s\n", r.style("•", tuikit.Cyan), item)
  }
  ```
- **CodeBlockNode** (lines 110–116):
  ```go
  fmt.Fprintf(w, "  %s\n", r.style("┌── "+n.Language, tuikit.Dim))

  for line := range bytesconv.ScanTokens(n.Code, '\n') {
  	fmt.Fprintf(w, "  %s %s\n", r.style("│", tuikit.Dim), r.style(line, tuikit.Cyan))
  }

  fmt.Fprintf(w, "  %s\n\n", r.style("└──", tuikit.Dim))
  ```
- **QuoteNode** (lines 119–122):
  ```go
  for line := range bytesconv.ScanTokens(n.Text, '\n') {
  	fmt.Fprintf(w, "  %s %s\n", r.style("▎", tuikit.Dim), r.style(line, tuikit.Dim))
  }
  ```
- **DividerNode** (line 126):
  ```go
  if r.isColorActive() {
  	fmt.Fprintf(w, "%s\n\n", tuikit.RenderDivider(72))
  } else {
  	fmt.Fprintf(w, "%s\n\n", strings.Repeat("─", 72))
  }
  ```

#### C. `renderCallout` via `tuikit.Box` (Rounded Card)
Replace lines 145–174 with:
```go
func (r *TerminalRenderer) renderCallout(w io.Writer, n CalloutNode) {
	title := n.Title
	if icon := n.Intent.Icon(); icon != "" {
		title = icon + " " + title
	}

	var colorFn func(string) string
	switch n.Intent {
	case IntentSuccess:
		colorFn = tuikit.Green
	case IntentWarning:
		colorFn = tuikit.Yellow
	case IntentDanger:
		colorFn = tuikit.Red
	case IntentInfo:
		colorFn = tuikit.Cyan
	case IntentMuted:
		colorFn = tuikit.Gray
	default:
		colorFn = tuikit.Cyan
	}

	styledTitle := r.style(title, colorFn)

	box := tuikit.NewBox(styledTitle, 0).
		SetStyle(tuikit.BorderRounded).
		SetIndent(2)

	if n.Body != "" {
		for line := range bytesconv.ScanTokens(n.Body, '\n') {
			box.AddLine(line)
		}
	}

	// Guarantee clean uncolored borders if color is disabled
	if !r.isColorActive() {
		prev := tuikit.ColorEnabled()
		tuikit.SetColorEnabled(false)
		_ = box.Render(w)
		tuikit.SetColorEnabled(prev)
	} else {
		_ = box.Render(w)
	}

	fmt.Fprintln(w)
}
```
*(Alternative Sovereign Left-Bar Rule Note: If an open left-rule format without a full enclosure is desired, `r.style("│", colorFn, tuikit.Bold)` followed by `r.style(title, colorFn, tuikit.Bold)` can be used instead of `NewBox`).*

#### D. `renderTable` via `tuikit.Table`
Replace lines 176–216 with:
```go
func toTuikitAlign(a Align) tuikit.Alignment {
	switch a {
	case AlignRight:
		return tuikit.AlignRight
	case AlignCenter:
		return tuikit.AlignCenter
	default:
		return tuikit.AlignLeft
	}
}

func (r *TerminalRenderer) renderTable(w io.Writer, t TableNode) {
	if len(t.Headers) == 0 {
		return
	}

	tbl := tuikit.NewTable(t.Headers...).SetIndent(2)

	for i, align := range t.Aligns {
		tbl.SetAlignment(i, toTuikitAlign(align))
	}

	for _, row := range t.Rows {
		tbl.AddRow(row...)
	}

	if !r.isColorActive() {
		prev := tuikit.ColorEnabled()
		tuikit.SetColorEnabled(false)
		_ = tbl.Render(w)
		tuikit.SetColorEnabled(prev)
	} else {
		_ = tbl.Render(w)
	}

	fmt.Fprintln(w)
}
```

### 4.3 Implementation Specification: `cmd/vortex/app.go` Terminal Probing

In `cmd/vortex/app.go`, add `"github.com/lemon4ksan/foundation/tuikit"` to imports, and in `Run(ctx context.Context, args []string) error` right after setting default `stdout` and `stderr` (line 109):
```go
	stdout := a.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	stderr := a.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	// Probe terminal capabilities and enforce NO_COLOR:
	// When stdout is not an interactive terminal (e.g. redirected to a pipe or file),
	// or when NO_COLOR / TERM=dumb is set, disable colors globally so all subcommands,
	// tables, and renderers emit 100% clean, unadorned plaintext.
	if !tuikit.ProbeTerminal(stdout) || os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		tuikit.SetColorEnabled(false)
	}
```

### 4.4 Test Assertions & Compatibility Matrix

| File | Line | Current Assertion | Change Required |
|---|---|---|---|
| `internal/text/builder_test.go` | 46–78 | Checks `[!TIP]`, `[SUCCESS]`, `"Release Notes"`, `"┌── go"` | **None** (passes cleanly) |
| `cmd/vortex/app_test.go` | 385 | `require.Contains(t, stdout.String(), "⚡ Vortex API Git History")` | Update synchronously when `internal/ast/log.go` removes `⚡` to `"◆ Vortex API Git History"` or `"Vortex API Git History"` |
| `cmd/vortex/app_test.go` | 396 | `require.Contains(t, stdout.String(), "⚡ [vortex diff]")` | Update synchronously when `internal/traffic/diff.go` removes `⚡` to `"◆ [vortex diff]"` or `"[vortex diff]"` |
| `cmd/vortex/app_test.go` | 452 | `require.Contains(t, stdout.String(), "⚡ Vortex Auto-Pilot: Audit & Build Pipeline")` | Update synchronously when `internal/core/autopilot.go` removes `⚡` |
| `pkg/project/project_test.go` | 352 | `require.Contains(t, rendered, "⚡ Vortex API Guardian")` | Update synchronously when `pkg/project/status.go` removes `⚡` |
| `pkg/lint/lint_test.go` | 234 | `require.Contains(t, output, "* W001 (param-lifting): 1")` | Update if `pkg/lint/format.go` transitions rule stats to a `tuikit.Table` |

---

## 5. Verification Method

1. **Verify Raw ANSI Escape Absence in `internal/text`**:
   ```pwsh
   git grep -n -E '\\033\[|\\x1b\[' internal/text/
   # Must return 0 matches
   ```

2. **Verify Emoji Absence in `internal/text/intent.go`**:
   ```pwsh
   git grep -n -P '[\x{26A1}\x{2728}\x{1F534}\x{1F7E1}\x{1F535}\x{1F7E2}\x{274C}\x{2705}\x{26A0}]' internal/text/intent.go
   # Must return 0 matches
   ```

3. **Verify Table Alignment & Callout Rendering**:
   Run unit tests in `internal/text`:
   ```pwsh
   go test -v ./internal/text/...
   # All tests must pass with 0 failures
   ```

4. **Verify NO_COLOR & Pipe Redirection**:
   ```pwsh
   # 1. Pipe redirection test (must contain zero ESC bytes)
   go run ./cmd/vortex status > status_pipe_test.txt
   Get-Content status_pipe_test.txt -Raw | Select-String "\x1b"
   Remove-Item status_pipe_test.txt

   # 2. NO_COLOR test
   $env:NO_COLOR = "1"
   go run ./cmd/vortex version
   Remove-Item Env:\NO_COLOR
   ```

5. **Verify Full Workspace Integrity**:
   ```pwsh
   go test ./...
   golangci-lint run
   ```
