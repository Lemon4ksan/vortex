// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package git

import (
	"errors"
	"strings"
)

var (
	// ErrNotRepository indicates the directory is not inside a Git work tree or repository.
	ErrNotRepository = errors.New("git: not a git repository")

	// ErrBranchNotFound indicates that a requested branch, tag, or ref could not be resolved.
	ErrBranchNotFound = errors.New("git: branch not found")

	// ErrGitCommandFailed indicates that a git CLI subprocess exited with a non-zero exit status.
	ErrGitCommandFailed = errors.New("git: command execution failed")
)

// GitError represents an execution or resolution failure in the Git subsystem.
type GitError struct {
	Op   string // Operation (e.g. "show", "merge_base", "log", "branch", "rev_parse", "blame")
	Path string // Repository root or relative file path
	Key  string // Ref, branch name, or commit hash
	Err  error  // Underlying exec.ExitError or sentinel error
}

func (e *GitError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("git")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" (ref: ")
			sb.WriteString(e.Key)
			sb.WriteByte(')')
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString("ref ")
		sb.WriteString(e.Key)
		sb.WriteString(": ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown git error")
	}

	return sb.String()
}

func (e *GitError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsNotRepository reports whether err indicates that the target path is not a Git repository.
func IsNotRepository(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrNotRepository) {
		return true
	}

	if gErr, ok := errors.AsType[*GitError](err); ok && gErr != nil {
		if errors.Is(gErr.Err, ErrNotRepository) {
			return true
		}
	}

	return strings.Contains(err.Error(), "not a git repository")
}

// IsNotFound reports whether err indicates a branch, tag, or ref not found.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrBranchNotFound) {
		return true
	}

	if gErr, ok := errors.AsType[*GitError](err); ok && gErr != nil {
		return errors.Is(gErr.Err, ErrBranchNotFound)
	}

	return false
}
