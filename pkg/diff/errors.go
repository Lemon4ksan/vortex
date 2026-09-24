// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"errors"
	"strings"
)

var (
	// ErrStackEmpty indicates an operation was invoked on an empty diff stack.
	ErrStackEmpty = errors.New("diff: stack is empty")

	// ErrFrameNotFound indicates that a requested stack frame label or index does not exist.
	ErrFrameNotFound = errors.New("diff: frame not found")

	// ErrInsufficientFrames indicates that adjacent or cumulative diff requires at least 2 frames.
	ErrInsufficientFrames = errors.New("diff: insufficient frames for diff")

	// ErrConflict indicates conflicting AST or schema migrations between diff frames.
	ErrConflict = errors.New("diff: conflicting schema changes")
)

// DiffError represents an error during diff calculation or stack manipulation.
type DiffError struct {
	Op   string // Operation (e.g. "push", "pop", "pop_to", "diff_adjacent", "diff_cumulative")
	Path string // Stack file path or affected target file
	Key  string // Frame index, frame label, or query
	Err  error  // Underlying or sentinel error
}

func (e *DiffError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("diff")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" [frame=")
			sb.WriteString(e.Key)
			sb.WriteByte(']')
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString("frame ")
		sb.WriteString(e.Key)
		sb.WriteString(": ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown diff error")
	}

	return sb.String()
}

func (e *DiffError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsNotFound reports whether err represents a missing stack frame.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrFrameNotFound) {
		return true
	}

	if dErr, ok := errors.AsType[*DiffError](err); ok && dErr != nil {
		return errors.Is(dErr.Err, ErrFrameNotFound)
	}

	return false
}

// IsConflict reports whether err represents a schema or contract migration conflict.
func IsConflict(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrConflict) {
		return true
	}

	if dErr, ok := errors.AsType[*DiffError](err); ok && dErr != nil {
		return errors.Is(dErr.Err, ErrConflict)
	}

	return false
}

// IsEmpty reports whether err represents an empty stack condition.
func IsEmpty(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrStackEmpty) {
		return true
	}

	if dErr, ok := errors.AsType[*DiffError](err); ok && dErr != nil {
		return errors.Is(dErr.Err, ErrStackEmpty)
	}

	return strings.Contains(err.Error(), "stack is empty")
}
