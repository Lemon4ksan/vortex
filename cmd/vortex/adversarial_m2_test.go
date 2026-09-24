// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/lemon4ksan/foundation/testing/require"
	"github.com/lemon4ksan/foundation/tuikit"

	"github.com/lemon4ksan/vortex/internal/text"
	"github.com/lemon4ksan/vortex/pkg/lint"
	"github.com/lemon4ksan/vortex/pkg/project"
	"github.com/lemon4ksan/vortex/pkg/tuple"
)

// containsAnyANSI reports whether data contains any raw ANSI ESC byte (0x1b).
func containsAnyANSI(s string) bool {
	return strings.Contains(s, "\x1b") || strings.Contains(s, "\033")
}

func TestMilestone2_Adversarial_NoColor_TuikitStyles(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	prev := tuikit.ColorEnabled()
	tuikit.SetColorEnabled(false)
	defer tuikit.SetColorEnabled(prev)

	require.False(t, tuikit.ColorEnabled())

	styles := []string{
		tuikit.Bold("test-bold"),
		tuikit.Dim("test-dim"),
		tuikit.Italic("test-italic"),
		tuikit.Underline("test-underline"),
		tuikit.Red("test-red"),
		tuikit.Green("test-green"),
		tuikit.Yellow("test-yellow"),
		tuikit.Blue("test-blue"),
		tuikit.Magenta("test-magenta"),
		tuikit.Cyan("test-cyan"),
		tuikit.Gray("test-gray"),
		tuikit.White("test-white"),
		tuikit.Badge("INFO", tuikit.Cyan),
		tuikit.RenderHeader("Header Title"),
	}

	for _, s := range styles {
		if containsAnyANSI(s) {
			t.Fatalf("ANSI escape leaked when NO_COLOR active: %q", s)
		}
	}
}

func TestMilestone2_Adversarial_LintFormatReport_NonTTY_And_NoColor(t *testing.T) {
	report := &lint.Report{
		ServicesChecked: 3,
		MethodsChecked:  15,
		FilesChecked:    2,
		Diagnostics: []lint.Diagnostic{
			{
				RuleID:     "E001",
				RuleName:   "syntax-error",
				Severity:   lint.SeverityError,
				FilePath:   "api.go",
				Line:       10,
				Column:     2,
				Message:    "Invalid syntax",
				Suggestion: "Fix it",
			},
			{
				RuleID:     "W001",
				RuleName:   "param-lifting",
				Severity:   lint.SeverityWarning,
				FilePath:   "api.go",
				Line:       20,
				Column:     4,
				Message:    "Duplicated parameter",
				Suggestion: "Lift parameter",
			},
			{
				RuleID:   "I001",
				RuleName: "perf-note",
				Severity: lint.SeverityInfo,
				FilePath: "api.go",
				Line:     30,
				Column:   1,
				Message:  "Performance recommendation",
			},
		},
	}

	// 1. Piped / Non-TTY writer (bytes.Buffer) with normal settings
	var bufNonTTY bytes.Buffer
	lint.FormatReport(&bufNonTTY, "api.go", report)
	outNonTTY := bufNonTTY.String()
	require.NotEmpty(t, outNonTTY)
	if containsAnyANSI(outNonTTY) {
		t.Fatalf("ANSI escape sequence leaked to non-interactive writer in FormatReport: %q", outNonTTY)
	}

	// 2. With NO_COLOR=1 set
	t.Setenv("NO_COLOR", "1")
	prev := tuikit.ColorEnabled()
	tuikit.SetColorEnabled(false)
	defer tuikit.SetColorEnabled(prev)

	var bufNoColor bytes.Buffer
	lint.FormatReport(&bufNoColor, "api.go", report)
	outNoColor := bufNoColor.String()
	require.NotEmpty(t, outNoColor)
	if containsAnyANSI(outNoColor) {
		t.Fatalf("ANSI escape sequence leaked under NO_COLOR=1 in FormatReport: %q", outNoColor)
	}

	// Empty report test
	emptyReport := &lint.Report{ServicesChecked: 1, MethodsChecked: 1, FilesChecked: 1}
	var bufEmpty bytes.Buffer
	lint.FormatReport(&bufEmpty, "api.go", emptyReport)
	outEmpty := bufEmpty.String()
	require.NotEmpty(t, outEmpty)
	if containsAnyANSI(outEmpty) {
		t.Fatalf("ANSI escape sequence leaked in empty FormatReport: %q", outEmpty)
	}
}

func TestMilestone2_Adversarial_ProjectStatus_NonTTY_And_NoColor(t *testing.T) {
	report := &project.StatusReport{
		WorkspaceRoot: "/workspace/demo",
		TotalMethods:  20,
		TotalDTOs:     10,
		Contracts: []project.ContractStatus{
			{
				Name:                  "ServiceA",
				File:                  "pkg/svc/api.go",
				MethodsCount:          10,
				DTOsCount:             5,
				IsGenStale:            false,
				Source:                "https://api.example.com/spec.json",
				UpstreamBreakingCount: 2,
			},
			{
				Name:           "ServiceB",
				File:           "pkg/svcb/api.go",
				MethodsCount:   10,
				DTOsCount:      5,
				IsGenStale:     true,
				GenStaleReason: "api.gen.go is STALE",
			},
		},
		NextActions: []string{"Run vortex gen"},
	}

	// 1. Render(false)
	plain := report.Render(false)
	if containsAnyANSI(plain) {
		t.Fatalf("ANSI escape leaked in StatusReport.Render(false): %q", plain)
	}

	// 2. Render(true) with NO_COLOR=1
	t.Setenv("NO_COLOR", "1")
	prev := tuikit.ColorEnabled()
	tuikit.SetColorEnabled(false)
	defer tuikit.SetColorEnabled(prev)

	coloredUnderNoColor := report.Render(true)
	if containsAnyANSI(coloredUnderNoColor) {
		t.Fatalf("ANSI escape leaked in StatusReport.Render(true) under NO_COLOR=1: %q", coloredUnderNoColor)
	}

	// 3. TERM=dumb
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "dumb")
	coloredUnderDumb := report.Render(true)
	if containsAnyANSI(coloredUnderDumb) {
		t.Fatalf("ANSI escape leaked in StatusReport.Render(true) under TERM=dumb: %q", coloredUnderDumb)
	}
}

func TestMilestone2_Adversarial_TerminalRenderer_NoColor(t *testing.T) {
	doc := text.NewDocument().
		Title("◆", "Test Document").
		Section("◆", "Details").
		Field("Key", "Value").
		FieldCode("CodeKey", "CodeVal").
		FieldBold("BoldKey", "BoldVal").
		Bullet("Item 1", "Item 2").
		Numbered("Step 1", "Step 2").
		Success("Success Title", "Success Body").
		Warning("Warning Title", "Warning Body").
		Danger("Danger Title", "Danger Body").
		Info("Info Title", "Info Body").
		Code("go", "func main() {}").
		Quote("Quote text").
		Divider().
		Table([]string{"H1", "H2"}, []string{"R1C1", "R1C2"})
	defer doc.Release()

	// 1. TerminalRenderer with ColorEnabled = false
	r := text.NewTerminalRenderer()
	r.ColorEnabled = false
	var buf bytes.Buffer
	err := r.Render(&buf, doc.Build())
	require.NoError(t, err)
	out := buf.String()
	if containsAnyANSI(out) {
		t.Fatalf("ANSI escape leaked in TerminalRenderer with ColorEnabled=false: %q", out)
	}

	// 2. TerminalRenderer under global NO_COLOR=1
	t.Setenv("NO_COLOR", "1")
	prev := tuikit.ColorEnabled()
	tuikit.SetColorEnabled(false)
	defer tuikit.SetColorEnabled(prev)

	rDefault := text.NewTerminalRenderer()
	buf.Reset()
	err = rDefault.Render(&buf, doc.Build())
	require.NoError(t, err)
	outNoColor := buf.String()
	if containsAnyANSI(outNoColor) {
		t.Fatalf("ANSI escape leaked in TerminalRenderer under NO_COLOR=1: %q", outNoColor)
	}
}

func TestMilestone2_Adversarial_CLIApp_PipedStdout_NoANSI(t *testing.T) {
	// Any CLI execution where Stdout is a pipe/buffer must never leak ANSI codes
	cmds := [][]string{
		{"--help"},
		{"--version"},
		{"list"},
		{"explain", "status"},
		{"example", "http"},
	}

	for _, cmdArgs := range cmds {
		var stdout, stderr bytes.Buffer
		app := newTestApp(&stdout, &stderr)
		err := app.Run(context.Background(), cmdArgs)
		require.NoError(t, err)

		out := stdout.String()
		if containsAnyANSI(out) {
			t.Fatalf("ANSI escape leaked in CLI app command %v output: %q", cmdArgs, out)
		}
		errOut := stderr.String()
		if containsAnyANSI(errOut) {
			t.Fatalf("ANSI escape leaked in CLI app command %v stderr: %q", cmdArgs, errOut)
		}
	}
}

func TestMilestone2_Adversarial_AllSubcommandsHelp_NoInformalEmojis(t *testing.T) {
	forbiddenEmojis := []string{"⚡", "✨", "🔴", "🟡", "🔵", "❌", "⚠️", "🚀", "🤖"}

	var stdout, stderr bytes.Buffer
	app := newTestApp(&stdout, &stderr)

	for _, cmd := range app.Commands {
		stdout.Reset()
		stderr.Reset()
		_ = app.Run(context.Background(), []string{"help", cmd.Name()})
		combined := stdout.String() + "\n" + stderr.String()
		for _, emoji := range forbiddenEmojis {
			if strings.Contains(combined, emoji) {
				t.Fatalf("Command 'help %s' leaked forbidden emoji %q:\n%s", cmd.Name(), emoji, combined)
			}
		}
	}
}

func TestMilestone2_Adversarial_NoInformalEmojisInCodebase(t *testing.T) {
	forbiddenEmojis := []string{"⚡", "✨", "🔴", "🟡", "🔵", "❌", "⚠️", "🚀"}

	rootDir, err := filepath.Abs("../..")
	require.NoError(t, err)

	type violation struct {
		file  string
		line  int
		text  string
		emoji string
	}
	var violations []violation

	dirs := []string{"cmd", "internal", "pkg"}
	for _, dir := range dirs {
		basePath := filepath.Join(rootDir, dir)
		walkErr := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			lines := strings.Split(string(data), "\n")
			for lineNum, line := range lines {
				for _, emoji := range forbiddenEmojis {
					if strings.Contains(line, emoji) {
						rel, _ := filepath.Rel(rootDir, path)
						violations = append(violations, violation{
							file:  rel,
							line:  lineNum + 1,
							text:  strings.TrimSpace(line),
							emoji: emoji,
						})
					}
				}
			}
			return nil
		})
		require.NoError(t, walkErr)
	}

	if len(violations) > 0 {
		var sb strings.Builder
		sb.WriteString("Forbidden informal emojis discovered in codebase production files:\n")
		for _, v := range violations {
			sb.WriteString(
				filepath.ToSlash(v.file) + ":" + strconv.Itoa(v.line) + " contains " + v.emoji + ": " + v.text + "\n",
			)
		}
		t.Fatalf("%s", sb.String())
	}
}

func TestMilestone2_Adversarial_SpecImport_NoInformalEmojis(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "api.go")
	specFile := filepath.Join(tmpDir, "spec.json")

	src := `package testapi

import (
	"context"
	"github.com/lemon4ksan/aoni"
)

// @aoni:service
type API interface {
	// @get "items/{id}"
	GetItem(ctx context.Context, id string) (map[string]any, error)
}
`
	require.NoError(t, os.WriteFile(srcFile, []byte(src), 0o600))

	spec := `{
  "openapi": "3.1.0",
  "info": {"title": "Test", "version": "v1.5.0"},
  "paths": {
    "/items/{id}": {
      "get": {
        "operationId": "api_v1_items_id_get",
        "parameters": [{"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}],
        "responses": {"200": {"description": "OK"}}
      }
    }
  }
}`
	require.NoError(t, os.WriteFile(specFile, []byte(spec), 0o600))

	var stdout, stderr bytes.Buffer
	app := newTestApp(&stdout, &stderr)

	err := app.Run(context.Background(), []string{"spec", "import", "-spec=" + specFile, "-out=" + srcFile})
	require.NoError(t, err)

	out := stdout.String()
	forbiddenEmojis := []string{"⚡", "✨", "🔴", "🟡", "🔵", "❌", "⚠️", "🚀"}
	for _, emoji := range forbiddenEmojis {
		if strings.Contains(out, emoji) {
			t.Fatalf("CLI 'spec import' stdout leaked forbidden informal emoji %q in output:\n%s", emoji, out)
		}
	}
}

func TestMilestone2_Adversarial_TupleRenderTable_NoInformalEmojis(t *testing.T) {
	r := &tuple.TupleAnalysisReport{
		StructName:   "TestTuple",
		TotalSamples: 10,
		Indices: []tuple.TupleIndexReport{
			{Index: 0, Occupancy: 1.0, NonNilCount: 10, InferredType: "string", DefaultName: "Field0"},
		},
	}
	out := r.RenderTable()
	forbiddenEmojis := []string{"⚡", "✨", "🔴", "🟡", "🔵", "❌", "⚠️", "🚀"}
	for _, emoji := range forbiddenEmojis {
		if strings.Contains(out, emoji) {
			t.Fatalf("TupleAnalysisReport.RenderTable leaked forbidden informal emoji %q:\n%s", emoji, out)
		}
	}
}

func TestMilestone2_Adversarial_TuikitTable_AlignmentAndBorders(t *testing.T) {
	tbl := tuikit.NewTable("STATUS", "TARGET", "ACTIONS", "LATENCY")
	tbl.SetIndent(2)
	tbl.SetAlignment(0, tuikit.AlignCenter)
	tbl.SetAlignment(1, tuikit.AlignLeft)
	tbl.SetAlignment(2, tuikit.AlignLeft)
	tbl.SetAlignment(3, tuikit.AlignRight)

	tbl.AddRow(
		tuikit.Badge("✔ IN SYNC", tuikit.Green),
		"UserService",
		"↳ no action needed",
		"[0.8ms | 0 allocs]",
	)
	tbl.AddRow(
		tuikit.Badge("✖ BREAKING", tuikit.Red),
		"PaymentGateway",
		"↳ re-generate client",
		"[12.4ms | 2 allocs]",
	)
	tbl.AddRow(
		tuikit.Badge("▲ DRIFT", tuikit.Yellow),
		"InventorySpec",
		"↳ merge external changes",
		"[1.2ms | 0 allocs]",
	)
	tbl.AddRow("◆ IDLE", "—", "—", "[0.0ms | 0 allocs]")

	// Test 1: Interactive mode (color enabled)
	tuikit.SetColorEnabled(true)
	defer tuikit.SetColorEnabled(true)

	outColored := tbl.String()
	linesColored := strings.Split(strings.TrimRight(outColored, "\n"), "\n")
	require.True(t, len(linesColored) >= 6) // header + divider + 4 rows

	dividerLine := linesColored[1]
	dividerWidth := tuikit.VisibleWidth(dividerLine)

	for i, l := range linesColored {
		w := tuikit.VisibleWidth(l)
		if w != dividerWidth {
			t.Fatalf(
				"Table line %d visible width mismatch (got %d, expected %d):\nLine: %q\nFull output:\n%s",
				i,
				w,
				dividerWidth,
				l,
				outColored,
			)
		}
	}

	// Test 2: Pipe / Non-TTY mode (color disabled)
	tuikit.SetColorEnabled(false)
	outPlain := tbl.String()
	if containsAnyANSI(outPlain) {
		t.Fatalf("ANSI escape leaked in Table when color disabled:\n%s", outPlain)
	}

	linesPlain := strings.Split(strings.TrimRight(outPlain, "\n"), "\n")
	dividerWidthPlain := tuikit.VisibleWidth(linesPlain[1])
	for i, l := range linesPlain {
		w := tuikit.VisibleWidth(l)
		if w != dividerWidthPlain {
			t.Fatalf(
				"Plain table line %d visible width mismatch (got %d, expected %d):\nLine: %q\nFull output:\n%s",
				i,
				w,
				dividerWidthPlain,
				l,
				outPlain,
			)
		}
	}
}

func TestMilestone2_Adversarial_TuikitBox_FramingAndCorners(t *testing.T) {
	styles := []struct {
		name  string
		style tuikit.BorderStyle
	}{
		{"Single", tuikit.BorderSingle},
		{"Rounded", tuikit.BorderRounded},
		{"Heavy", tuikit.BorderHeavy},
	}

	for _, s := range styles {
		box := tuikit.NewBox("◆ Diagnostic Report", 60)
		box.SetStyle(s.style)
		box.SetIndent(4)
		box.AddLine("✔ Workspace integrity verified cleanly")
		box.AddLine("✖ Method contract drift in PaymentGateway")
		box.AddDivider()
		box.AddRow("Status Code", "HTTP 200 OK")
		box.AddRow("Execution Cost", "[1.2ms | 0 allocs]")
		box.AddLine("↳ Sub-task complete — zero allocations recorded")

		// 1. Color enabled
		tuikit.SetColorEnabled(true)
		outColored := box.String()
		linesColored := strings.Split(strings.TrimRight(outColored, "\n"), "\n")
		expectedW := 60 + 2 + 4 // inner 60 + 2 border columns + 4 indent = 66
		for i, l := range linesColored {
			w := tuikit.VisibleWidth(l)
			if w != expectedW {
				t.Fatalf(
					"Box (%s) line %d visible width %d != %d:\nLine: %q\nFull:\n%s",
					s.name,
					i,
					w,
					expectedW,
					l,
					outColored,
				)
			}
		}

		// 2. Color disabled / pipe
		tuikit.SetColorEnabled(false)
		outPlain := box.String()
		if containsAnyANSI(outPlain) {
			t.Fatalf("ANSI escape leaked in Box (%s) under pipe mode:\n%s", s.name, outPlain)
		}
		linesPlain := strings.Split(strings.TrimRight(outPlain, "\n"), "\n")
		for i, l := range linesPlain {
			w := tuikit.VisibleWidth(l)
			if w != expectedW {
				t.Fatalf(
					"Plain Box (%s) line %d visible width %d != %d:\nLine: %q\nFull:\n%s",
					s.name,
					i,
					w,
					expectedW,
					l,
					outPlain,
				)
			}
		}
	}
}
