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

	if pErr, ok := errors.AsType[*ProjectError](err); ok && pErr != nil {
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

	if pErr, ok := errors.AsType[*ProjectError](err); ok && pErr != nil {
		return errors.Is(pErr.Err, ErrStaleCodegen)
	}

	return false
}
