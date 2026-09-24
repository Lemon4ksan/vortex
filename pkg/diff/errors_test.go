// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lemon4ksan/foundation/testing/assert"
	"github.com/lemon4ksan/foundation/testing/require"

	"github.com/lemon4ksan/vortex/pkg/diff"
)

func TestDiffErrors_IsNotFound(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, diff.IsNotFound(diff.ErrFrameNotFound))

	// 2. Wrapped sentinel
	require.True(t, diff.IsNotFound(fmt.Errorf("wrap: %w", diff.ErrFrameNotFound)))

	// 3. Typed struct
	dErr := &diff.DiffError{
		Op:  "pop_to",
		Key: "tag-v1",
		Err: diff.ErrFrameNotFound,
	}
	require.True(t, diff.IsNotFound(dErr))

	dErrWrapped := fmt.Errorf("outer: %w", dErr)
	require.True(t, diff.IsNotFound(dErrWrapped))

	// 4. Negative case
	require.False(t, diff.IsNotFound(errors.New("unrelated error")))
	require.False(t, diff.IsNotFound(diff.ErrConflict))

	// 5. Nil error
	require.False(t, diff.IsNotFound(nil))
	var typedNil *diff.DiffError
	require.False(t, diff.IsNotFound(typedNil))
	require.False(t, diff.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
}

func TestDiffErrors_IsConflict(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, diff.IsConflict(diff.ErrConflict))

	// 2. Wrapped sentinel
	require.True(t, diff.IsConflict(fmt.Errorf("wrap: %w", diff.ErrConflict)))

	// 3. Typed struct
	dErr := &diff.DiffError{
		Op:  "diff_adjacent",
		Err: diff.ErrConflict,
	}
	require.True(t, diff.IsConflict(dErr))

	// 4. Negative case
	require.False(t, diff.IsConflict(errors.New("unrelated error")))
	require.False(t, diff.IsConflict(diff.ErrStackEmpty))

	// 5. Nil error
	require.False(t, diff.IsConflict(nil))
	var typedNil *diff.DiffError
	require.False(t, diff.IsConflict(typedNil))
	require.False(t, diff.IsConflict(fmt.Errorf("wrap: %w", typedNil)))
}

func TestDiffErrors_IsEmpty(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, diff.IsEmpty(diff.ErrStackEmpty))

	// 2. Wrapped sentinel
	require.True(t, diff.IsEmpty(fmt.Errorf("wrap: %w", diff.ErrStackEmpty)))

	// 3. Typed struct
	dErr := &diff.DiffError{
		Op:  "pop",
		Err: diff.ErrStackEmpty,
	}
	require.True(t, diff.IsEmpty(dErr))

	// 3b. String fallback match
	require.True(t, diff.IsEmpty(errors.New("stack is empty, nothing to pop")))

	// 4. Negative case
	require.False(t, diff.IsEmpty(errors.New("unrelated error")))
	require.False(t, diff.IsEmpty(diff.ErrFrameNotFound))

	// 5. Nil error
	require.False(t, diff.IsEmpty(nil))
	var typedNil *diff.DiffError
	require.False(t, diff.IsEmpty(typedNil))
	require.False(t, diff.IsEmpty(fmt.Errorf("wrap: %w", typedNil)))
}

func TestDiffError_ErrorAndUnwrap(t *testing.T) {
	t.Parallel()

	// Nil receiver
	var nilErr *diff.DiffError
	assert.Equal(t, "<nil>", nilErr.Error())
	assert.Nil(t, nilErr.Unwrap())

	// Op, Path, Key, Err
	e1 := &diff.DiffError{
		Op:   "pop_to",
		Path: "stack.json",
		Key:  "frame_3",
		Err:  diff.ErrFrameNotFound,
	}
	assert.Equal(t, "diff/pop_to: stack.json [frame=frame_3]: diff: frame not found", e1.Error())
	assert.Equal(t, diff.ErrFrameNotFound, e1.Unwrap())
	require.True(t, errors.Is(e1, diff.ErrFrameNotFound))

	// Op, Path, Err (no Key)
	e2 := &diff.DiffError{
		Op:   "save",
		Path: "stack.json",
		Err:  diff.ErrStackEmpty,
	}
	assert.Equal(t, "diff/save: stack.json: diff: stack is empty", e2.Error())

	// Op, Key, Err (no Path)
	e3 := &diff.DiffError{
		Op:  "inspect",
		Key: "idx_0",
		Err: diff.ErrInsufficientFrames,
	}
	assert.Equal(t, "diff/inspect: frame idx_0: diff: insufficient frames for diff", e3.Error())

	// Op, nil Err
	e4 := &diff.DiffError{
		Op: "peek",
	}
	assert.Equal(t, "diff/peek: unknown diff error", e4.Error())
}
