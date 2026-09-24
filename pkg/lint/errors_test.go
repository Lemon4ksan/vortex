// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lint_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lemon4ksan/foundation/testing/assert"
	"github.com/lemon4ksan/foundation/testing/require"

	"github.com/lemon4ksan/vortex/pkg/lint"
)

func TestLintErrors_IsLintFailure(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, lint.IsLintFailure(lint.ErrLintFailure))

	// 2. Wrapped sentinel
	require.True(t, lint.IsLintFailure(fmt.Errorf("wrap: %w", lint.ErrLintFailure)))

	// 3. Typed struct
	lErr := &lint.LintError{
		Op:   "run_pass",
		Path: "service.go",
		Key:  "S001",
		Err:  lint.ErrLintFailure,
	}
	require.True(t, lint.IsLintFailure(lErr))

	lErrWrapped := fmt.Errorf("outer: %w", lErr)
	require.True(t, lint.IsLintFailure(lErrWrapped))

	// 4. Negative case
	require.False(t, lint.IsLintFailure(errors.New("unrelated error")))
	require.False(t, lint.IsLintFailure(lint.ErrRuleNotFound))

	// 5. Nil error
	require.False(t, lint.IsLintFailure(nil))
	var typedNil *lint.LintError
	require.False(t, lint.IsLintFailure(typedNil))
	require.False(t, lint.IsLintFailure(fmt.Errorf("wrap: %w", typedNil)))
}

func TestLintErrors_IsNotFound(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, lint.IsNotFound(lint.ErrRuleNotFound))

	// 2. Wrapped sentinel
	require.True(t, lint.IsNotFound(fmt.Errorf("wrap: %w", lint.ErrRuleNotFound)))

	// 3. Typed struct
	lErr := &lint.LintError{
		Op:  "lookup",
		Key: "UNKNOWN_RULE",
		Err: lint.ErrRuleNotFound,
	}
	require.True(t, lint.IsNotFound(lErr))

	// 4. Negative case
	require.False(t, lint.IsNotFound(errors.New("unrelated error")))
	require.False(t, lint.IsNotFound(lint.ErrFixFailed))

	// 5. Nil error
	require.False(t, lint.IsNotFound(nil))
	var typedNil *lint.LintError
	require.False(t, lint.IsNotFound(typedNil))
	require.False(t, lint.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
}

func TestLintError_ErrorAndUnwrap(t *testing.T) {
	t.Parallel()

	// Nil receiver
	var nilErr *lint.LintError
	assert.Equal(t, "<nil>", nilErr.Error())
	assert.Nil(t, nilErr.Unwrap())

	// Op, Path, Key, Err
	e1 := &lint.LintError{
		Op:   "apply_fix",
		Path: "contract.go",
		Key:  "P002",
		Err:  lint.ErrFixFailed,
	}
	assert.Equal(t, "lint/apply_fix: contract.go [rule=P002]: lint: automated fix failed", e1.Error())
	assert.Equal(t, lint.ErrFixFailed, e1.Unwrap())
	require.True(t, errors.Is(e1, lint.ErrFixFailed))

	// Op, Path, Err (no Key)
	e2 := &lint.LintError{
		Op:   "lint_file",
		Path: "types.go",
		Err:  lint.ErrLintFailure,
	}
	assert.Equal(t, "lint/lint_file: types.go: lint: static analysis checks failed", e2.Error())

	// Op, Key, Err (no Path)
	e3 := &lint.LintError{
		Op:  "enable_rule",
		Key: "CUSTOM01",
		Err: lint.ErrRuleNotFound,
	}
	assert.Equal(t, "lint/enable_rule: [CUSTOM01]: lint: rule not found", e3.Error())

	// Op, nil Err
	e4 := &lint.LintError{
		Op: "check",
	}
	assert.Equal(t, "lint/check: unknown lint error", e4.Error())
}
