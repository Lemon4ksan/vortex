// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spec_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/lemon4ksan/foundation/testing/assert"
	"github.com/lemon4ksan/foundation/testing/require"

	"github.com/lemon4ksan/vortex/pkg/spec"
)

func TestSpecErrors_IsNotFound(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, spec.IsNotFound(spec.ErrSpecNotFound))
	require.True(t, spec.IsNotFound(os.ErrNotExist))

	// 2. Wrapped sentinel
	require.True(t, spec.IsNotFound(fmt.Errorf("wrap: %w", spec.ErrSpecNotFound)))
	require.True(t, spec.IsNotFound(fmt.Errorf("wrap: %w", os.ErrNotExist)))

	// 3. Typed struct
	sErr := &spec.SpecError{
		Op:   "fetch",
		Path: "openapi.yaml",
		Err:  spec.ErrSpecNotFound,
	}
	require.True(t, spec.IsNotFound(sErr))

	sErrWrapped := fmt.Errorf("outer: %w", sErr)
	require.True(t, spec.IsNotFound(sErrWrapped))

	// 4. Negative case
	require.False(t, spec.IsNotFound(errors.New("unrelated error")))
	require.False(t, spec.IsNotFound(spec.ErrUnsupportedFormat))

	// 5. Nil error
	require.False(t, spec.IsNotFound(nil))
	var typedNil *spec.SpecError
	require.False(t, spec.IsNotFound(typedNil))
	require.False(t, spec.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
}

func TestSpecErrors_IsUnsupportedFormat(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, spec.IsUnsupportedFormat(spec.ErrUnsupportedFormat))

	// 2. Wrapped sentinel
	require.True(t, spec.IsUnsupportedFormat(fmt.Errorf("wrap: %w", spec.ErrUnsupportedFormat)))

	// 3. Typed struct
	sErr := &spec.SpecError{
		Op:  "parse",
		Key: "graphql",
		Err: spec.ErrUnsupportedFormat,
	}
	require.True(t, spec.IsUnsupportedFormat(sErr))

	// 4. Negative case
	require.False(t, spec.IsUnsupportedFormat(errors.New("unrelated error")))
	require.False(t, spec.IsUnsupportedFormat(spec.ErrSpecNotFound))

	// 5. Nil error
	require.False(t, spec.IsUnsupportedFormat(nil))
	var typedNil *spec.SpecError
	require.False(t, spec.IsUnsupportedFormat(typedNil))
	require.False(t, spec.IsUnsupportedFormat(fmt.Errorf("wrap: %w", typedNil)))
}

func TestSpecError_ErrorAndUnwrap(t *testing.T) {
	t.Parallel()

	// Nil receiver
	var nilErr *spec.SpecError
	assert.Equal(t, "<nil>", nilErr.Error())
	assert.Nil(t, nilErr.Unwrap())

	// Op, Path, Key, Err
	e1 := &spec.SpecError{
		Op:   "validate",
		Path: "spec.json",
		Key:  "openapi-3.1",
		Err:  spec.ErrEmptySpec,
	}
	assert.Equal(t, "spec/validate: spec.json (format: openapi-3.1): spec: specification is empty", e1.Error())
	assert.Equal(t, spec.ErrEmptySpec, e1.Unwrap())
	require.True(t, errors.Is(e1, spec.ErrEmptySpec))

	// Op, Path, Err (no Key)
	e2 := &spec.SpecError{
		Op:   "fetch",
		Path: "https://example.com/spec.json",
		Err:  spec.ErrSpecNotFound,
	}
	assert.Equal(t, "spec/fetch: https://example.com/spec.json: spec: specification not found", e2.Error())

	// Op, Key, Err (no Path)
	e3 := &spec.SpecError{
		Op:  "parse",
		Key: "proto2",
		Err: spec.ErrUnsupportedFormat,
	}
	assert.Equal(t, "spec/parse: format proto2: spec: unsupported specification format", e3.Error())

	// Op, nil Err
	e4 := &spec.SpecError{
		Op: "scan",
	}
	assert.Equal(t, "spec/scan: unknown spec error", e4.Error())
}
