// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package git_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lemon4ksan/foundation/testing/assert"
	"github.com/lemon4ksan/foundation/testing/require"

	"github.com/lemon4ksan/vortex/pkg/git"
)

func TestGitErrors_IsNotRepository(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, git.IsNotRepository(git.ErrNotRepository))

	// 2. Wrapped sentinel
	require.True(t, git.IsNotRepository(fmt.Errorf("wrap: %w", git.ErrNotRepository)))

	// 3. Typed struct
	gErr := &git.GitError{
		Op:   "rev_parse",
		Path: "/non-git-dir",
		Err:  git.ErrNotRepository,
	}
	require.True(t, git.IsNotRepository(gErr))

	gErrWrapped := fmt.Errorf("outer: %w", gErr)
	require.True(t, git.IsNotRepository(gErrWrapped))

	// 3b. String fallback match
	require.True(t, git.IsNotRepository(errors.New("fatal: not a git repository (or any of the parent directories)")))

	// 4. Negative case
	require.False(t, git.IsNotRepository(errors.New("unrelated error")))
	require.False(t, git.IsNotRepository(git.ErrBranchNotFound))

	// 5. Nil error
	require.False(t, git.IsNotRepository(nil))
	var typedNil *git.GitError
	require.False(t, git.IsNotRepository(typedNil))
	require.False(t, git.IsNotRepository(fmt.Errorf("wrap: %w", typedNil)))
}

func TestGitErrors_IsNotFound(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, git.IsNotFound(git.ErrBranchNotFound))

	// 2. Wrapped sentinel
	require.True(t, git.IsNotFound(fmt.Errorf("wrap: %w", git.ErrBranchNotFound)))

	// 3. Typed struct
	gErr := &git.GitError{
		Op:  "show",
		Key: "feature-branch",
		Err: git.ErrBranchNotFound,
	}
	require.True(t, git.IsNotFound(gErr))

	// 4. Negative case
	require.False(t, git.IsNotFound(errors.New("unrelated error")))
	require.False(t, git.IsNotFound(git.ErrNotRepository))

	// 5. Nil error
	require.False(t, git.IsNotFound(nil))
	var typedNil *git.GitError
	require.False(t, git.IsNotFound(typedNil))
	require.False(t, git.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
}

func TestGitError_ErrorAndUnwrap(t *testing.T) {
	t.Parallel()

	// Nil receiver
	var nilErr *git.GitError
	assert.Equal(t, "<nil>", nilErr.Error())
	assert.Nil(t, nilErr.Unwrap())

	// Op, Path, Key, Err
	e1 := &git.GitError{
		Op:   "show",
		Path: "repo",
		Key:  "HEAD~1",
		Err:  git.ErrGitCommandFailed,
	}
	assert.Equal(t, "git/show: repo (ref: HEAD~1): git: command execution failed", e1.Error())
	assert.Equal(t, git.ErrGitCommandFailed, e1.Unwrap())
	require.True(t, errors.Is(e1, git.ErrGitCommandFailed))

	// Op, Path, Err (no Key)
	e2 := &git.GitError{
		Op:   "root_dir",
		Path: "/workspace",
		Err:  git.ErrNotRepository,
	}
	assert.Equal(t, "git/root_dir: /workspace: git: not a git repository", e2.Error())

	// Op, Key, Err (no Path)
	e3 := &git.GitError{
		Op:  "resolve",
		Key: "origin/main",
		Err: git.ErrBranchNotFound,
	}
	assert.Equal(t, "git/resolve: ref origin/main: git: branch not found", e3.Error())

	// Op, nil Err
	e4 := &git.GitError{
		Op: "status",
	}
	assert.Equal(t, "git/status: unknown git error", e4.Error())
}
