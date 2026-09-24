// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"errors"
	"strings"
)

var (
	// ErrSyntaxError indicates a Go syntax or DSL parsing failure in source code.
	ErrSyntaxError = errors.New("parser: syntax error")

	// ErrContractNotFound indicates that no declarative service contract was discovered in the parsed source.
	ErrContractNotFound = errors.New("parser: contract not found")

	// ErrInvalidDirective indicates an unrecognized, misplaced, or malformed DSL directive or argument.
	ErrInvalidDirective = errors.New("parser: invalid directive")

	// ErrUnresolvedType indicates a referenced type or generic argument cannot be bound or resolved.
	ErrUnresolvedType = errors.New("parser: unresolved type")
)

// ParseError represents a syntax or semantic error discovered during contract AST parsing.
type ParseError struct {
	Op   string // Operation (e.g. "parse_file", "parse_source", "bind_type", "parse_directive")
	Path string // Source file path
	Key  string // Line:column position, directive name, or symbol identifier
	Err  error  // Underlying syntax or sentinel error
}

func (e *ParseError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("parser")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteByte(':')
			sb.WriteString(e.Key)
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString(e.Key)
		sb.WriteString(": ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown parse error")
	}

	return sb.String()
}

func (e *ParseError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsSyntaxError reports whether err represents a syntax parsing error in Go or DSL code.
func IsSyntaxError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrSyntaxError) {
		return true
	}

	if pErr, ok := errors.AsType[*ParseError](err); ok && pErr != nil {
		return errors.Is(pErr.Err, ErrSyntaxError)
	}

	return strings.Contains(err.Error(), "syntax error")
}

// IsNotFound reports whether err represents a missing contract in the parsed source.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrContractNotFound) {
		return true
	}

	if pErr, ok := errors.AsType[*ParseError](err); ok && pErr != nil {
		return errors.Is(pErr.Err, ErrContractNotFound)
	}

	return false
}
