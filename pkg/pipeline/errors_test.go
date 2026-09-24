// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipeline_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lemon4ksan/foundation/testing/assert"
	"github.com/lemon4ksan/foundation/testing/require"

	"github.com/lemon4ksan/vortex/pkg/pipeline"
)

func TestPipelineErrors_IsPipelineAborted(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, pipeline.IsPipelineAborted(pipeline.ErrPipelineAborted))

	// 2. Wrapped sentinel
	require.True(t, pipeline.IsPipelineAborted(fmt.Errorf("wrap: %w", pipeline.ErrPipelineAborted)))

	// 3. Typed struct
	pErr := &pipeline.PipelineError{
		Op:  "compile",
		Key: "stage_optimize",
		Err: pipeline.ErrPipelineAborted,
	}
	require.True(t, pipeline.IsPipelineAborted(pErr))

	pErrWrapped := fmt.Errorf("outer: %w", pErr)
	require.True(t, pipeline.IsPipelineAborted(pErrWrapped))

	// 4. Negative case
	require.False(t, pipeline.IsPipelineAborted(errors.New("unrelated error")))
	require.False(t, pipeline.IsPipelineAborted(pipeline.ErrTargetFileRequired))

	// 5. Nil error
	require.False(t, pipeline.IsPipelineAborted(nil))
	var typedNil *pipeline.PipelineError
	require.False(t, pipeline.IsPipelineAborted(typedNil))
	require.False(t, pipeline.IsPipelineAborted(fmt.Errorf("wrap: %w", typedNil)))
}

func TestPipelineErrors_IsNotFound(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, pipeline.IsNotFound(pipeline.ErrNoContractsFound))

	// 2. Wrapped sentinel
	require.True(t, pipeline.IsNotFound(fmt.Errorf("wrap: %w", pipeline.ErrNoContractsFound)))

	// 3. Typed struct
	pErr := &pipeline.PipelineError{
		Op:   "extract_env",
		Path: "pkg/api",
		Err:  pipeline.ErrNoContractsFound,
	}
	require.True(t, pipeline.IsNotFound(pErr))

	// 4. Negative case
	require.False(t, pipeline.IsNotFound(errors.New("unrelated error")))
	require.False(t, pipeline.IsNotFound(pipeline.ErrPipelineAborted))

	// 5. Nil error
	require.False(t, pipeline.IsNotFound(nil))
	var typedNil *pipeline.PipelineError
	require.False(t, pipeline.IsNotFound(typedNil))
	require.False(t, pipeline.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
}

func TestPipelineError_ErrorAndUnwrap(t *testing.T) {
	t.Parallel()

	// Nil receiver
	var nilErr *pipeline.PipelineError
	assert.Equal(t, "<nil>", nilErr.Error())
	assert.Nil(t, nilErr.Unwrap())

	// Op, Path, Key, Err
	e1 := &pipeline.PipelineError{
		Op:   "deobfuscate",
		Path: "contract.go",
		Key:  "tuple_unpack",
		Err:  pipeline.ErrPipelineAborted,
	}
	assert.Equal(t, "pipeline/deobfuscate: contract.go [stage=tuple_unpack]: pipeline: processing aborted", e1.Error())
	assert.Equal(t, pipeline.ErrPipelineAborted, e1.Unwrap())
	require.True(t, errors.Is(e1, pipeline.ErrPipelineAborted))

	// Op, Path, Err (no Key)
	e2 := &pipeline.PipelineError{
		Op:   "extract_env",
		Path: "pkg/api",
		Err:  pipeline.ErrNoContractsFound,
	}
	assert.Equal(t, "pipeline/extract_env: pkg/api: pipeline: no contracts found", e2.Error())

	// Op, Key, Err (no Path)
	e3 := &pipeline.PipelineError{
		Op:  "stage_run",
		Key: "init",
		Err: pipeline.ErrTargetFileRequired,
	}
	assert.Equal(t, "pipeline/stage_run: stage init: pipeline: target file is required", e3.Error())

	// Op, nil Err
	e4 := &pipeline.PipelineError{
		Op: "process",
	}
	assert.Equal(t, "pipeline/process: unknown pipeline error", e4.Error())
}
