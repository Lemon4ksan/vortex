// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// NewTerminalRenderer constructs a fresh [TerminalRenderer] with ANSI colors enabled by default.
func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{ColorEnabled: true}
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

// Render writes the terminal-formatted representation of doc to w.
func (r *TerminalRenderer) Render(w io.Writer, doc *Document) error {
	if doc == nil || len(doc.Nodes) == 0 {
		return nil
	}

	for i, node := range doc.Nodes {
		switch n := node.(type) {
		case HeadingNode:
			title := n.Text
			if n.Icon != "" {
				title = n.Icon + " " + title
			}

			switch n.Level {
			case 1:
				fmt.Fprintf(w, "%s\n\n", r.style(title, tuikit.Underline, tuikit.White, tuikit.Bold))
			case 2:
				fmt.Fprintf(w, "%s\n\n", r.style(title, tuikit.Cyan, tuikit.Bold))
			default:
				fmt.Fprintf(w, "%s\n\n", r.style(title, tuikit.Bold))
			}

		case SectionNode:
			sec := n.Name
			if n.Icon != "" {
				sec = n.Icon + " " + sec
			}

			fmt.Fprintf(w, "%s\n", r.style(sec+":", tuikit.Yellow, tuikit.Bold))

		case ParagraphNode:
			fmt.Fprintf(w, "%s\n\n", n.Text)

		case FieldNode:
			key := r.style(n.Key, tuikit.Bold)
			val := n.Value

			switch n.Style {
			case FieldCode:
				val = r.style("`"+n.Value+"`", tuikit.Cyan)
			case FieldBold:
				val = r.style(n.Value, tuikit.Bold)
			}

			fmt.Fprintf(w, "  • %s: %s\n", key, val)

		case ListNode:
			for idx, item := range n.Items {
				if n.Kind == ListNumbered {
					fmt.Fprintf(w, "  %s %s\n", r.style(fmt.Sprintf("%d.", idx+1), tuikit.Dim), item)
				} else {
					fmt.Fprintf(w, "  %s %s\n", r.style("•", tuikit.Cyan), item)
				}
			}

			fmt.Fprintln(w)

		case CalloutNode:
			r.renderCallout(w, n)

		case CodeBlockNode:
			fmt.Fprintf(w, "  %s\n", r.style("┌── "+n.Language, tuikit.Dim))

			for line := range bytesconv.ScanTokens(n.Code, '\n') {
				fmt.Fprintf(w, "  %s %s\n", r.style("│", tuikit.Dim), r.style(line, tuikit.Cyan))
			}

			fmt.Fprintf(w, "  %s\n\n", r.style("└──", tuikit.Dim))

		case QuoteNode:
			for line := range bytesconv.ScanTokens(n.Text, '\n') {
				fmt.Fprintf(w, "  %s %s\n", r.style("▎", tuikit.Dim), r.style(line, tuikit.Dim))
			}

			fmt.Fprintln(w)

		case DividerNode:
			if r.isColorActive() {
				fmt.Fprintf(w, "%s\n\n", tuikit.RenderDivider(72))
			} else {
				fmt.Fprintf(w, "%s\n\n", strings.Repeat("─", 72))
			}

		case TableNode:
			r.renderTable(w, n)

		case RawNode:
			fmt.Fprint(w, n.Content)
		}

		if _, isField := node.(FieldNode); isField && i+1 < len(doc.Nodes) {
			if _, nextIsField := doc.Nodes[i+1].(FieldNode); !nextIsField {
				fmt.Fprintln(w)
			}
		}
	}

	return nil
}

func (r *TerminalRenderer) renderCallout(w io.Writer, n CalloutNode) {
	colorFn := tuikit.Cyan
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
	}

	title := n.Title
	if icon := n.Intent.Icon(); icon != "" {
		title = icon + " " + title
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

	if !r.isColorActive() {
		var buf strings.Builder
		_ = box.Render(&buf)
		_, _ = io.WriteString(w, tuikit.StripANSI(buf.String()))
	} else {
		_ = box.Render(w)
	}

	fmt.Fprintln(w)
}

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
		var buf strings.Builder
		_ = tbl.Render(&buf)
		_, _ = io.WriteString(w, tuikit.StripANSI(buf.String()))
	} else {
		_ = tbl.Render(w)
	}

	fmt.Fprintln(w)
}
