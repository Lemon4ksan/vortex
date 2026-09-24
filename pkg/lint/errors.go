// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lint

import (
	"errors"
	"strings"
)

var (
	// ErrLintFailure is returned when contracts fail static analysis checks (e.g. SeverityError count > 0).
	ErrLintFailure = errors.New("lint: static analysis checks failed")

	// ErrRuleNotFound indicates that a requested rule ID is not registered in the lint registry.
	ErrRuleNotFound = errors.New("lint: rule not found")

	// ErrFixFailed indicates that an automated code fix could not be successfully applied.
	ErrFixFailed = errors.New("lint: automated fix failed")
)

// LintError represents an operational error in rule execution or automated fix application.
type LintError struct {
	Op   string // Operation (e.g. "run_pass", "apply_fix", "register_rule")
	Path string // Inspected source file path
	Key  string // Rule ID (e.g. "S001", "B002", "missing-context")
	Err  error  // Underlying diagnostic or sentinel error
}

func (e *LintError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("lint")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" [rule=")
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
		sb.WriteString("unknown lint error")
	}

	return sb.String()
}

func (e *LintError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsLintFailure reports whether err represents a static analysis check failure.
func IsLintFailure(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrLintFailure) {
		return true
	}

	if lErr, ok := errors.AsType[*LintError](err); ok && lErr != nil {
		return errors.Is(lErr.Err, ErrLintFailure)
	}

	return false
}

// IsNotFound reports whether err indicates a missing lint rule.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrRuleNotFound) {
		return true
	}

	if lErr, ok := errors.AsType[*LintError](err); ok && lErr != nil {
		return errors.Is(lErr.Err, ErrRuleNotFound)
	}

	return false
}
