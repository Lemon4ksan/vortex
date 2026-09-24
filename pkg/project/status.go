// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package project

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lemon4ksan/foundation/pathkit"
	"github.com/lemon4ksan/foundation/tuikit"

	"github.com/lemon4ksan/vortex/pkg/diff"
	"github.com/lemon4ksan/vortex/pkg/git"
	"github.com/lemon4ksan/vortex/pkg/ingest"
	"github.com/lemon4ksan/vortex/pkg/openapi"
	"github.com/lemon4ksan/vortex/pkg/parser"
)

// PluginStatus reports the synchronization health of an external polyglot target.
type PluginStatus struct {
	Name    string `json:"name"`
	Out     string `json:"out"`
	IsStale bool   `json:"is_stale"`
	Message string `json:"message,omitempty"`
}

// ContractStatus captures the integrity and freshness state of an individual contract.
type ContractStatus struct {
	Name                  string         `json:"name"`
	Package               string         `json:"package"`
	File                  string         `json:"file"`
	GenFile               string         `json:"gen_file,omitempty"`
	ModelsFile            string         `json:"models_file,omitempty"`
	MethodsCount          int            `json:"methods_count"`
	DTOsCount             int            `json:"dtos_count"`
	Version               string         `json:"version,omitempty"`
	Source                string         `json:"source,omitempty"`
	IsGenStale            bool           `json:"is_gen_stale"`
	GenStaleReason        string         `json:"gen_stale_reason,omitempty"`
	UpstreamBreakingCount int            `json:"upstream_breaking_count"`
	UpstreamDriftCount    int            `json:"upstream_drift_count"`
	UpstreamGhostCount    int            `json:"upstream_ghost_count"`
	Plugins               []PluginStatus `json:"plugins,omitempty"`
}

// StatusReport summarizes the health of all monitored contracts across the workspace.
type StatusReport struct {
	WorkspaceRoot string               `json:"workspace_root"`
	ConfigPath    string               `json:"config_path,omitempty"`
	Contracts     []ContractStatus     `json:"contracts"`
	Proposals     []git.BranchProposal `json:"proposals,omitempty"`
	TotalMethods  int                  `json:"total_methods"`
	TotalDTOs     int                  `json:"total_dtos"`
	NextActions   []string             `json:"next_actions,omitempty"`
}

// StaleCount returns the number of contracts requiring code generation.
func (r *StatusReport) StaleCount() int {
	count := 0
	for _, c := range r.Contracts {
		if c.IsGenStale {
			count++
		}
	}

	return count
}

// BreakingDriftCount returns the number of breaking changes across all upstream specs.
func (r *StatusReport) BreakingDriftCount() int {
	count := 0
	for _, c := range r.Contracts {
		count += c.UpstreamBreakingCount
	}

	return count
}

// HasIssues returns true if any contract is stale or has breaking upstream changes.
func (r *StatusReport) HasIssues() bool {
	return r.StaleCount() > 0 || r.BreakingDriftCount() > 0
}

// Render formats a colored, human-readable terminal dashboard.
func (r *StatusReport) Render(color bool) string {
	useColor := color && tuikit.ColorEnabled() && os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"

	var sb strings.Builder
	title := "◆ Vortex API Guardian"
	if useColor {
		title = tuikit.Bold(tuikit.Cyan(title))
	} else {
		title = tuikit.RenderHeader(title)
	}
	sb.WriteString(title + "\n")

	meta := fmt.Sprintf("(%d services, %d methods)", len(r.Contracts), r.TotalMethods)
	if useColor {
		meta = tuikit.Dim(meta)
	}
	fmt.Fprintf(&sb, "Workspace: %s %s\n\n", r.WorkspaceRoot, meta)

	if len(r.Contracts) == 0 {
		sb.WriteString("No service contracts detected. Run `vortex init` to configure workspace.\n")
		return sb.String()
	}

	contractsHeader := "● Contracts & Generated Code:"
	if useColor {
		contractsHeader = tuikit.Bold(contractsHeader)
	}
	sb.WriteString(contractsHeader + "\n")

	hasAnyDTOs := false
	hasAnyVersion := false
	for _, c := range r.Contracts {
		if c.DTOsCount > 0 {
			hasAnyDTOs = true
		}
		if c.Version != "" {
			hasAnyVersion = true
		}
	}

	tblContracts := tuikit.NewTable().SetIndent(2)

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

	// Upstream Drift section
	hasUpstream := false
	for _, c := range r.Contracts {
		if c.Source != "" {
			hasUpstream = true
			break
		}
	}

	if hasUpstream {
		hdr := "● Upstream Drift (OpenAPI / External):"
		if useColor {
			hdr = tuikit.Bold(hdr)
		}
		sb.WriteString(hdr + "\n")

		tblDrift := tuikit.NewTable().SetIndent(2)

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
	}

	// Polyglot plugins
	hasPlugins := false
	for _, c := range r.Contracts {
		if len(c.Plugins) > 0 {
			hasPlugins = true
			break
		}
	}

	if hasPlugins {
		hdr := "● Polyglot Targets:"
		if useColor {
			hdr = tuikit.Bold(hdr)
		}
		sb.WriteString(hdr + "\n")

		tblPoly := tuikit.NewTable().SetIndent(2)

		for _, c := range r.Contracts {
			for _, p := range c.Plugins {
				targetName := strings.ToUpper(p.Name) + " SDK"
				outStr := "(" + p.Out + ")"
				if useColor {
					outStr = tuikit.Dim(outStr)
				}

				var badge, statusStyled string
				if p.IsStale {
					if useColor {
						badge = tuikit.Badge("▲ STALE", tuikit.Yellow)
						statusStyled = tuikit.Yellow("Stale (rebuild required)")
					} else {
						badge = "▲ STALE"
						statusStyled = "Stale (rebuild required)"
					}
				} else {
					if useColor {
						badge = tuikit.Badge("✔ IN SYNC", tuikit.Green)
						statusStyled = tuikit.Green("Up to date")
					} else {
						badge = "✔ IN SYNC"
						statusStyled = "Up to date"
					}
				}

				tblPoly.AddRow(badge, targetName, outStr, statusStyled)
			}
		}

		sb.WriteString(tblPoly.String() + "\n")
	}

	if len(r.Proposals) > 0 {
		hdr := "● Incoming Consumer Proposals (Git Branches):"
		if useColor {
			hdr = tuikit.Bold(hdr)
		}
		sb.WriteString(hdr + "\n")

		tblProp := tuikit.NewTable().SetIndent(2)

		for _, prop := range r.Proposals {
			arrow := "↳"
			if useColor {
				arrow = tuikit.Cyan("↳")
			}

			remoteTag := ""
			if prop.IsRemote {
				remoteTag = " [remote]"
				if useColor {
					remoteTag = tuikit.Cyan(remoteTag)
				}
			}

			author := "@" + prop.Author
			authorWithBy := "by " + author

			dateStr := fmt.Sprintf("(%s)", prop.Date)
			if useColor {
				dateStr = tuikit.Dim(dateStr)
			}

			propDesc := prop.Name + remoteTag
			tblProp.AddRow(arrow, propDesc, authorWithBy, dateStr)
		}

		sb.WriteString(tblProp.String() + "\n")
	}

	if len(r.NextActions) > 0 {
		divider := tuikit.RenderDivider(67)
		if !useColor {
			divider = tuikit.StripANSI(divider)
		}
		sb.WriteString(divider + "\n")

		hdr := "Next Actions:"
		if useColor {
			hdr = tuikit.Bold(tuikit.Yellow(hdr))
		}
		sb.WriteString(hdr + "\n")

		for _, action := range r.NextActions {
			arrow := "↳"
			if useColor {
				arrow = tuikit.Cyan("↳")
			}

			fmt.Fprintf(&sb, "  %s %s\n", arrow, action)
		}
	} else {
		msg := "✔ All systems nominal. Network layer is 100% synchronized.\n"
		if useColor {
			msg = tuikit.Bold(tuikit.Green(msg))
		}

		sb.WriteString(msg)
	}

	result := sb.String()
	if !useColor {
		result = tuikit.StripANSI(result)
	}

	return result
}

// RenderJSON serializes the status report into JSON bytes.
func (r *StatusReport) RenderJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// StatusEngine inspects workspace contracts, freshness, and upstream drift.
type StatusEngine struct{}

// NewStatusEngine creates an initialized StatusEngine instance.
func NewStatusEngine() *StatusEngine {
	return &StatusEngine{}
}

// Inspect evaluates contracts in the workspace against disk state and upstream sources.
func (e *StatusEngine) Inspect(cfg *Config, contracts []ContractConfig) *StatusReport {
	report := &StatusReport{
		WorkspaceRoot: cfg.RootDir,
		ConfigPath:    cfg.ConfigPath,
		Contracts:     make([]ContractStatus, 0, len(contracts)),
		NextActions:   make([]string, 0),
	}

	p := parser.NewParser()

	var (
		staleServiceNames []string
		driftServiceNames []string
	)

	for _, ct := range contracts {
		absFile := filepath.Join(cfg.RootDir, ct.File)

		srcInfo, err := os.Stat(absFile)
		if err != nil {
			continue
		}

		root, err := p.ParseFile(absFile)
		if err != nil || len(root.Services) == 0 {
			continue
		}

		totalMethods := 0
		for _, s := range root.Services {
			totalMethods += len(s.Methods)
		}

		svcName := ct.Name
		if svcName == "" {
			svcName = root.Services[0].Name
			if len(root.Services) > 1 {
				svcName = fmt.Sprintf("%s (+%d)", root.Services[0].Name, len(root.Services)-1)
			}
		}

		status := ContractStatus{
			Name:         svcName,
			Package:      root.PackageName,
			File:         ct.File,
			GenFile:      ct.Gen,
			ModelsFile:   ct.Models,
			MethodsCount: totalMethods,
			DTOsCount:    len(root.Structs),
			Version:      root.Services[0].Version,
			Source:       root.Services[0].Source,
		}

		if ct.Upstream != nil && ct.Upstream.Source != "" && status.Source == "" {
			status.Source = ct.Upstream.Source
		}

		report.TotalMethods += status.MethodsCount
		report.TotalDTOs += status.DTOsCount

		// 1. Check generated code freshness
		if ct.Gen != "" {
			absGen := filepath.Join(cfg.RootDir, ct.Gen)

			genInfo, genErr := os.Stat(absGen)
			if genErr != nil {
				status.IsGenStale = true
				status.GenStaleReason = "api.gen.go is missing"
			} else if srcInfo.ModTime().After(genInfo.ModTime()) {
				status.IsGenStale = true
				status.GenStaleReason = "api.gen.go is STALE"
			}
		}

		if !status.IsGenStale && ct.Models != "" {
			absModels := filepath.Join(cfg.RootDir, ct.Models)

			modelsInfo, modelsErr := os.Stat(absModels)
			if modelsErr == nil && srcInfo.ModTime().After(modelsInfo.ModTime()) {
				status.IsGenStale = true
				status.GenStaleReason = "models.gen.go is STALE"
			} else if modelsErr != nil && cfg.ConfigPath != "" {
				status.IsGenStale = true
				status.GenStaleReason = "models.gen.go is missing"
			}
		}

		if status.IsGenStale {
			staleServiceNames = append(staleServiceNames, status.Name)
		}

		// 2. Check upstream drift if source is available
		if status.Source != "" {
			pSrc := pathkit.New(status.Source)

			upstreamPath := pSrc.FilePath()
			if pSrc.IsFile() && !pSrc.IsAbs() {
				upstreamPath = filepath.Join(cfg.RootDir, upstreamPath)
			}

			rawSpecBytes, readErr := os.ReadFile(upstreamPath)
			specFormat, _ := ingest.DetectFormat(rawSpecBytes)

			if readErr == nil && (specFormat == ingest.FormatOpenAPI3 || specFormat == ingest.FormatSwagger2) {
				if doc, docErr := openapi.LoadSpec(upstreamPath, nil); docErr == nil {
					diffReport := diff.Compare(root, doc, ct.File, status.Source)
					status.UpstreamBreakingCount = diffReport.BreakingCount()
					status.UpstreamDriftCount = diffReport.NonBreakingCount()
					status.UpstreamGhostCount = diffReport.GhostCount()

					if status.UpstreamBreakingCount > 0 || status.UpstreamDriftCount > 0 {
						driftServiceNames = append(
							driftServiceNames,
							fmt.Sprintf("%s (%s)", status.Name, status.Source),
						)
					}
				}
			} else if upInfo, upErr := os.Stat(upstreamPath); upErr == nil {
				// If not standard OpenAPI (e.g. proprietary JSON/YAML dump), check file modification timestamp
				if upInfo.ModTime().After(srcInfo.ModTime()) {
					status.UpstreamDriftCount = 1
					if ct.Upstream != nil && ct.Upstream.Generate != "" {
						report.NextActions = append(
							report.NextActions,
							fmt.Sprintf(
								"Run `%s` to rebuild contract from updated upstream dump (%s)",
								ct.Upstream.Generate,
								status.Source,
							),
						)
					} else {
						driftServiceNames = append(
							driftServiceNames,
							fmt.Sprintf("%s (%s)", status.Name, status.Source),
						)
					}
				}
			}
		}

		// 3. Check polyglot plugins
		for _, pl := range ct.Plugins {
			absPlOut := filepath.Join(cfg.RootDir, pl.Out)
			plInfo, plErr := os.Stat(absPlOut)

			pStatus := PluginStatus{
				Name:    pl.Name,
				Out:     pl.Out,
				IsStale: false,
			}
			if plErr != nil || srcInfo.ModTime().After(plInfo.ModTime()) {
				pStatus.IsStale = true
			}

			status.Plugins = append(status.Plugins, pStatus)
		}

		report.Contracts = append(report.Contracts, status)
	}

	// Build actionable NextActions
	if len(staleServiceNames) > 0 {
		report.NextActions = append(report.NextActions,
			fmt.Sprintf("Run `vortex gen` to rebuild stale Go code (%s)", strings.Join(staleServiceNames, ", ")),
		)
	}

	if len(driftServiceNames) > 0 {
		report.NextActions = append(
			report.NextActions,
			fmt.Sprintf(
				"Run `vortex oapi import` to reconcile upstream changes (%s)",
				strings.Join(driftServiceNames, ", "),
			),
		)
	}

	// 4. Discover incoming proposals from Git branches
	if proposals, pErr := git.ListProposalBranches(
		context.Background(),
		cfg.RootDir,
		nil,
	); pErr == nil &&
		len(proposals) > 0 {
		report.Proposals = proposals
		for _, prop := range proposals {
			report.NextActions = append(
				report.NextActions,
				fmt.Sprintf(
					"Run `vortex review %s` or `vortex accept %s` to inspect/merge proposal",
					prop.Name,
					prop.Name,
				),
			)
		}
	}

	return report
}
