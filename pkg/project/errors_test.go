// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package project_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/lemon4ksan/foundation/testing/assert"
	"github.com/lemon4ksan/foundation/testing/require"

	"github.com/lemon4ksan/vortex/pkg/project"
)

func TestProjectErrors_IsNotFound(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, project.IsNotFound(project.ErrWorkspaceNotFound))
	require.True(t, project.IsNotFound(project.ErrConfigNotFound))
	require.True(t, project.IsNotFound(project.ErrContractNotFound))
	require.True(t, project.IsNotFound(os.ErrNotExist))

	// 2. Wrapped sentinel
	require.True(t, project.IsNotFound(fmt.Errorf("wrap: %w", project.ErrWorkspaceNotFound)))
	require.True(t, project.IsNotFound(fmt.Errorf("wrap: %w", project.ErrConfigNotFound)))
	require.True(t, project.IsNotFound(fmt.Errorf("wrap: %w", project.ErrContractNotFound)))
	require.True(t, project.IsNotFound(fmt.Errorf("wrap: %w", os.ErrNotExist)))

	// 3. Typed struct
	pErr := &project.ProjectError{
		Op:   "find_root",
		Path: "/path/to/repo",
		Err:  project.ErrWorkspaceNotFound,
	}
	require.True(t, project.IsNotFound(pErr))

	pErrWrapped := fmt.Errorf("outer: %w", pErr)
	require.True(t, project.IsNotFound(pErrWrapped))

	// 4. Negative case
	require.False(t, project.IsNotFound(errors.New("unrelated error")))
	require.False(t, project.IsNotFound(project.ErrInvalidConfig))

	// 5. Nil error
	require.False(t, project.IsNotFound(nil))
	var typedNil *project.ProjectError
	require.False(t, project.IsNotFound(typedNil))
	require.False(t, project.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
}

func TestProjectErrors_IsStale(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, project.IsStale(project.ErrStaleCodegen))

	// 2. Wrapped sentinel
	require.True(t, project.IsStale(fmt.Errorf("wrap: %w", project.ErrStaleCodegen)))

	// 3. Typed struct
	pErr := &project.ProjectError{
		Op:  "check_status",
		Key: "petstore",
		Err: project.ErrStaleCodegen,
	}
	require.True(t, project.IsStale(pErr))

	// 4. Negative case
	require.False(t, project.IsStale(errors.New("unrelated error")))
	require.False(t, project.IsStale(project.ErrInvalidConfig))

	// 5. Nil error
	require.False(t, project.IsStale(nil))
	var typedNil *project.ProjectError
	require.False(t, project.IsStale(typedNil))
	require.False(t, project.IsStale(fmt.Errorf("wrap: %w", typedNil)))
}

func TestProjectError_ErrorAndUnwrap(t *testing.T) {
	t.Parallel()

	// Nil receiver
	var nilErr *project.ProjectError
	assert.Equal(t, "<nil>", nilErr.Error())
	assert.Nil(t, nilErr.Unwrap())

	// Op, Path, Key, Err
	e1 := &project.ProjectError{
		Op:   "load",
		Path: ".vortex.yml",
		Key:  "services",
		Err:  project.ErrInvalidConfig,
	}
	assert.Equal(t, "project/load: .vortex.yml [services]: project: invalid configuration", e1.Error())
	assert.Equal(t, project.ErrInvalidConfig, e1.Unwrap())
	require.True(t, errors.Is(e1, project.ErrInvalidConfig))

	// Op, Path, Err (no Key)
	e2 := &project.ProjectError{
		Op:   "discover",
		Path: "/root",
		Err:  project.ErrConfigNotFound,
	}
	assert.Equal(t, "project/discover: /root: project: configuration file not found", e2.Error())

	// Op, Key, Err (no Path)
	e3 := &project.ProjectError{
		Op:  "validate",
		Key: "contract-1",
		Err: project.ErrContractNotFound,
	}
	assert.Equal(t, "project/validate: [contract-1]: project: contract not found", e3.Error())

	// Op, nil Err
	e4 := &project.ProjectError{
		Op: "unknown",
	}
	assert.Equal(t, "project/unknown: unknown project error", e4.Error())
}
