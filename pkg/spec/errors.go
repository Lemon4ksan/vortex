// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spec

import (
	"errors"
	"os"
	"strings"
)

var (
	// ErrSpecNotFound indicates that an OpenAPI, AsyncAPI, or HAR spec definition was not found.
	ErrSpecNotFound = errors.New("spec: specification not found")

	// ErrUnsupportedFormat indicates that an unknown or invalid specification format was provided.
	ErrUnsupportedFormat = errors.New("spec: unsupported specification format")

	// ErrEmptySpec indicates that the provided specification document contains no paths, channels, or services.
	ErrEmptySpec = errors.New("spec: specification is empty")
)

// SpecError represents an error encountered while resolving or validating API specifications.
type SpecError struct {
	Op   string // Operation (e.g. "fetch", "parse", "validate", "lookup_directive")
	Path string // File path or URL of the specification
	Key  string // Specification format, directive name, or component key
	Err  error  // Underlying or sentinel error
}

func (e *SpecError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("spec")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" (format: ")
			sb.WriteString(e.Key)
			sb.WriteByte(')')
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString("format ")
		sb.WriteString(e.Key)
		sb.WriteString(": ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown spec error")
	}

	return sb.String()
}

func (e *SpecError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsNotFound reports whether err represents a missing specification file or resource.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrSpecNotFound) || errors.Is(err, os.ErrNotExist) {
		return true
	}

	if sErr, ok := errors.AsType[*SpecError](err); ok && sErr != nil {
		return errors.Is(sErr.Err, ErrSpecNotFound) || errors.Is(sErr.Err, os.ErrNotExist)
	}

	return false
}

// IsUnsupportedFormat reports whether err represents an unknown or invalid specification format.
func IsUnsupportedFormat(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrUnsupportedFormat) {
		return true
	}

	if sErr, ok := errors.AsType[*SpecError](err); ok && sErr != nil {
		return errors.Is(sErr.Err, ErrUnsupportedFormat)
	}

	return false
}
