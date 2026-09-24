// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lint

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/lemon4ksan/foundation/tuikit"

	"github.com/lemon4ksan/vortex/pkg/version"
)

// FormatReport writes a formatted terminal report of discovered diagnostics.
func FormatReport(w io.Writer, target string, report *Report) {
	if report == nil {
		return
	}

	useColor := tuikit.IsInteractive(w)
	var buf bytes.Buffer

	fmt.Fprintln(&buf, tuikit.RenderHeader("◆ Vortex Contract Inspector"))
	fmt.Fprintf(&buf, "%s\n\n", tuikit.Dim(fmt.Sprintf("Target: %s (%d services, %d methods across %d files)",
		target, report.ServicesChecked, report.MethodsChecked, report.FilesChecked)))

	if len(report.Diagnostics) == 0 {
		if report.SuppressedCount > 0 {
			fmt.Fprintf(
				&buf,
				"%s %s\n",
				tuikit.Bold(tuikit.Green("✔ All contracts are valid and synchronized!")),
				tuikit.Dim(fmt.Sprintf("(%d warnings suppressed via //vortex:ignore)", report.SuppressedCount)),
			)
		} else {
			fmt.Fprintln(&buf, tuikit.Bold(tuikit.Green("✔ All contracts are valid and synchronized!")))
		}

		out := buf.String()
		if !useColor {
			out = tuikit.StripANSI(out)
		}
		_, _ = io.WriteString(w, out)
		return
	}

	// Group by severity
	var errs, warns, infos []Diagnostic

	for _, d := range report.Diagnostics {
		switch d.Severity {
		case SeverityError:
			errs = append(errs, d)
		case SeverityWarning:
			warns = append(warns, d)
		default:
			infos = append(infos, d)
		}
	}

	sortDiagnostics(errs)
	sortDiagnostics(warns)
	sortDiagnostics(infos)

	if len(errs) > 0 {
		fmt.Fprintln(&buf, tuikit.Bold(tuikit.Red(fmt.Sprintf("◆ Errors (%d)", len(errs)))))

		for _, d := range errs {
			printDiagnostic(&buf, d, tuikit.Red)
		}

		fmt.Fprintln(&buf)
	}

	if len(warns) > 0 {
		fmt.Fprintln(&buf, tuikit.Bold(tuikit.Yellow(fmt.Sprintf("◆ Warnings & Suggestions (%d)", len(warns)))))

		for _, d := range warns {
			printDiagnostic(&buf, d, tuikit.Yellow)
		}

		fmt.Fprintln(&buf)
	}

	if len(infos) > 0 {
		fmt.Fprintln(&buf, tuikit.Bold(tuikit.Cyan(fmt.Sprintf("◆ Info (%d)", len(infos)))))

		for _, d := range infos {
			printDiagnostic(&buf, d, tuikit.Cyan)
		}

		fmt.Fprintln(&buf)
	}

	// Summary with colored severity breakdown and rule count table
	fmt.Fprintf(&buf, "%s ", tuikit.Bold("Summary:"))

	var parts []string
	if report.Errors() > 0 {
		parts = append(parts, tuikit.Red(fmt.Sprintf("%d error(s)", report.Errors())))
	}

	if report.Warnings() > 0 {
		parts = append(parts, tuikit.Yellow(fmt.Sprintf("%d warning(s)", report.Warnings())))
	}

	if report.FixableCount() > 0 {
		parts = append(parts, tuikit.Green(fmt.Sprintf("%d auto-fixable", report.FixableCount())))
	}

	if report.SuppressedCount > 0 {
		parts = append(parts, tuikit.Dim(fmt.Sprintf("%d suppressed", report.SuppressedCount)))
	}

	fmt.Fprintln(&buf, strings.Join(parts, ", "))

	type ruleStat struct {
		ruleID   string
		ruleName string
		severity Severity
		count    int
	}

	ruleMap := make(map[string]*ruleStat)
	for _, d := range report.Diagnostics {
		if s, exists := ruleMap[d.RuleID]; exists {
			s.count++
		} else {
			ruleMap[d.RuleID] = &ruleStat{
				ruleID:   d.RuleID,
				ruleName: d.RuleName,
				severity: d.Severity,
				count:    1,
			}
		}
	}

	var stats []*ruleStat
	for _, s := range ruleMap {
		stats = append(stats, s)
	}

	slices.SortFunc(stats, func(a, b *ruleStat) int {
		if a.count != b.count {
			return cmp.Compare(b.count, a.count)
		}

		return cmp.Compare(a.ruleID, b.ruleID)
	})

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

	if report.FixableCount() > 0 {
		fixMsg := fmt.Sprintf(
			"Run `vortex check --fix` to automatically resolve %d safe issue(s).",
			report.FixableCount(),
		)
		fmt.Fprintf(&buf, "\n%s\n", tuikit.Cyan(fixMsg))
	}

	out := buf.String()
	if !useColor {
		out = tuikit.StripANSI(out)
	}
	_, _ = io.WriteString(w, out)
}

func sortDiagnostics(diags []Diagnostic) {
	slices.SortStableFunc(diags, func(a, b Diagnostic) int {
		return cmp.Or(
			cmp.Compare(a.RuleID, b.RuleID),
			cmp.Compare(a.FilePath, b.FilePath),
			cmp.Compare(a.Line, b.Line),
			cmp.Compare(a.Column, b.Column),
		)
	})
}

func printDiagnostic(w io.Writer, d Diagnostic, colorFn func(string) string) {
	line := d.Line
	if line <= 0 {
		line = 1
	}

	col := d.Column
	if col <= 0 {
		col = 1
	}

	filePath := d.FilePath
	if abs, err := filepath.Abs(filePath); err == nil {
		filePath = abs
	}

	loc := fmt.Sprintf("%s:%d:%d", filePath, line, col)

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
}

func (d Diagnostic) Fixable() bool {
	return d.Fix != nil
}

// FormatGitHubActions outputs diagnostics as GitHub Actions Workflow Command annotations.
// Format: ::error file={name},line={line},title={title}::{message}
func FormatGitHubActions(w io.Writer, report *Report) {
	if report == nil {
		return
	}

	for _, d := range report.Diagnostics {
		level := "warning"
		switch d.Severity {
		case SeverityError:
			level = "error"
		case SeverityInfo:
			level = "notice"
		}

		filePath := filepath.ToSlash(d.FilePath)

		line := d.Line
		if line <= 0 {
			line = 1
		}

		fmt.Fprintf(w, "::%s file=%s,line=%d,title=%s::[%s] %s\n",
			level, filePath, line, d.RuleID, d.RuleName, d.Message)
	}
}

type sarifLog struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	InformationURI string      `json:"informationUri"`
	Rules          []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description sarifMsgWrapper `json:"shortDescription"`
}

type sarifResult struct {
	RuleID    string             `json:"ruleId"`
	Level     string             `json:"level"`
	Message   sarifMsgWrapper    `json:"message"`
	Locations []sarifLocationObj `json:"locations"`
}

type sarifMsgWrapper struct {
	Text string `json:"text"`
}

type sarifLocationObj struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifact `json:"artifactLocation"`
	Region           sarifRegion   `json:"region"`
}

type sarifArtifact struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine int `json:"startLine"`
}

// FormatSARIF outputs diagnostics adhering to OASIS SARIF v2.1.0 JSON format.
func FormatSARIF(w io.Writer, report *Report) error {
	if report == nil {
		return nil
	}

	ruleMap := make(map[string]bool)

	var (
		rules   []sarifRule
		results []sarifResult
	)

	for _, d := range report.Diagnostics {
		if !ruleMap[d.RuleID] {
			ruleMap[d.RuleID] = true
			rules = append(rules, sarifRule{
				ID:          d.RuleID,
				Name:        d.RuleName,
				Description: sarifMsgWrapper{Text: d.RuleName},
			})
		}

		level := "warning"
		switch d.Severity {
		case SeverityError:
			level = "error"
		case SeverityInfo:
			level = "note"
		}

		line := d.Line
		if line <= 0 {
			line = 1
		}

		results = append(results, sarifResult{
			RuleID:  d.RuleID,
			Level:   level,
			Message: sarifMsgWrapper{Text: d.Message},
			Locations: []sarifLocationObj{
				{
					PhysicalLocation: sarifPhysicalLocation{
						ArtifactLocation: sarifArtifact{
							URI: filepath.ToSlash(d.FilePath),
						},
						Region: sarifRegion{
							StartLine: line,
						},
					},
				},
			},
		})
	}

	log := sarifLog{
		Version: "2.1.0",
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:           "vortex",
						Version:        version.Number,
						InformationURI: "https://github.com/lemon4ksan/aoni",
						Rules:          rules,
					},
				},
				Results: results,
			},
		},
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	return encoder.Encode(log)
}
