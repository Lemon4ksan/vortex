// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipeline

import (
	"errors"
	"strings"
)

var (
	// ErrTargetFileRequired is returned when an AST or compilation pipeline is invoked without a target file.
	ErrTargetFileRequired = errors.New("pipeline: target file is required")

	// ErrNoContractsFound is returned when no contracts are discovered to execute the pipeline on.
	ErrNoContractsFound = errors.New("pipeline: no contracts found")

	// ErrPipelineAborted is returned when the compilation or deobfuscation pipeline halts due to an unrecoverable condition.
	ErrPipelineAborted = errors.New("pipeline: processing aborted")
)

// PipelineError represents a processing or compilation failure during pipeline execution.
type PipelineError struct {
	Op   string // Operation (e.g. "deobfuscate", "rename_field", "extract_env", "compile")
	Path string // Target contract file path
	Key  string // Pipeline stage name or transformation target
	Err  error  // Underlying or sentinel error
}

func (e *PipelineError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("pipeline")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" [stage=")
			sb.WriteString(e.Key)
			sb.WriteByte(']')
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString("stage ")
		sb.WriteString(e.Key)
		sb.WriteString(": ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown pipeline error")
	}

	return sb.String()
}

func (e *PipelineError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsPipelineAborted reports whether err represents an aborted pipeline run.
func IsPipelineAborted(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrPipelineAborted) {
		return true
	}

	if pErr, ok := errors.AsType[*PipelineError](err); ok && pErr != nil {
		return errors.Is(pErr.Err, ErrPipelineAborted)
	}

	return false
}

// IsNotFound reports whether err indicates that no contracts were found.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrNoContractsFound) {
		return true
	}

	if pErr, ok := errors.AsType[*PipelineError](err); ok && pErr != nil {
		return errors.Is(pErr.Err, ErrNoContractsFound)
	}

	return false
}
