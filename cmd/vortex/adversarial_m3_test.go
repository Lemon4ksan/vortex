// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/lemon4ksan/foundation/testing/assert"
	"github.com/lemon4ksan/foundation/testing/require"

	"github.com/lemon4ksan/vortex/pkg/cache"
	"github.com/lemon4ksan/vortex/pkg/diff"
	"github.com/lemon4ksan/vortex/pkg/git"
	"github.com/lemon4ksan/vortex/pkg/lint"
	"github.com/lemon4ksan/vortex/pkg/parser"
	"github.com/lemon4ksan/vortex/pkg/pipeline"
	"github.com/lemon4ksan/vortex/pkg/project"
	"github.com/lemon4ksan/vortex/pkg/spec"
)

// helper to wrap an error N times
func wrapN(err error, n int) error {
	w := err
	for i := 0; i < n; i++ {
		w = fmt.Errorf("wrap_%d: %w", i, w)
	}
	return w
}

type predicateTest struct {
	name string
	fn   func(error) bool
}

func getAll17Predicates() []predicateTest {
	return []predicateTest{
		{"project.IsNotFound", project.IsNotFound},
		{"project.IsStale", project.IsStale},
		{"parser.IsSyntaxError", parser.IsSyntaxError},
		{"parser.IsNotFound", parser.IsNotFound},
		{"diff.IsNotFound", diff.IsNotFound},
		{"diff.IsConflict", diff.IsConflict},
		{"diff.IsEmpty", diff.IsEmpty},
		{"git.IsNotRepository", git.IsNotRepository},
		{"git.IsNotFound", git.IsNotFound},
		{"cache.IsNotFound", cache.IsNotFound},
		{"cache.IsCorrupt", cache.IsCorrupt},
		{"lint.IsLintFailure", lint.IsLintFailure},
		{"lint.IsNotFound", lint.IsNotFound},
		{"spec.IsNotFound", spec.IsNotFound},
		{"spec.IsUnsupportedFormat", spec.IsUnsupportedFormat},
		{"pipeline.IsPipelineAborted", pipeline.IsPipelineAborted},
		{"pipeline.IsNotFound", pipeline.IsNotFound},
	}
}

func getAll8TypedNils() []struct {
	pkg string
	err error
} {
	var pProject *project.ProjectError
	var pParser *parser.ParseError
	var pDiff *diff.DiffError
	var pGit *git.GitError
	var pCache *cache.CacheError
	var pLint *lint.LintError
	var pSpec *spec.SpecError
	var pPipeline *pipeline.PipelineError

	return []struct {
		pkg string
		err error
	}{
		{"project", pProject},
		{"parser", pParser},
		{"diff", pDiff},
		{"git", pGit},
		{"cache", pCache},
		{"lint", pLint},
		{"spec", pSpec},
		{"pipeline", pPipeline},
	}
}

// TestAdversarialM3_TypedNil_ExhaustiveMatrix verifies that every predicate
// returns false without panicking when given any typed nil or wrapped typed nil.
func TestAdversarialM3_TypedNil_ExhaustiveMatrix(t *testing.T) {
	t.Parallel()

	preds := getAll17Predicates()
	typedNils := getAll8TypedNils()

	// 1. Untyped nil across all 17
	for _, pred := range preds {
		t.Run("UntypedNil_"+pred.name, func(t *testing.T) {
			require.False(t, pred.fn(nil))
		})
	}

	// 2. Direct typed nil: 17 predicates * 8 typed nils = 136 test cases
	for _, pred := range preds {
		for _, tn := range typedNils {
			testName := fmt.Sprintf("Direct_%s_with_%sTypedNil", pred.name, tn.pkg)
			t.Run(testName, func(t *testing.T) {
				require.False(t, pred.fn(tn.err))
			})
		}
	}

	// 3. Wrapped typed nil: 17 predicates * 8 typed nils = 136 test cases
	for _, pred := range preds {
		for _, tn := range typedNils {
			testName := fmt.Sprintf("Wrapped_%s_with_%sTypedNil", pred.name, tn.pkg)
			t.Run(testName, func(t *testing.T) {
				wrapped := fmt.Errorf("outer_wrap: %w", tn.err)
				require.False(t, pred.fn(wrapped))
			})
		}
	}

	// 4. Double wrapped typed nil: 17 predicates * 8 typed nils = 136 test cases
	for _, pred := range preds {
		for _, tn := range typedNils {
			testName := fmt.Sprintf("DoubleWrapped_%s_with_%sTypedNil", pred.name, tn.pkg)
			t.Run(testName, func(t *testing.T) {
				wrapped := fmt.Errorf("outer2: %w", fmt.Errorf("outer1: %w", tn.err))
				require.False(t, pred.fn(wrapped))
			})
		}
	}

	// 5. 100-level wrapped typed nil across matching subsystem predicates
	for _, tn := range typedNils {
		t.Run("Deep100_WrappedTypedNil_"+tn.pkg, func(t *testing.T) {
			deepWrapped := wrapN(tn.err, 100)
			for _, pred := range preds {
				require.False(t, pred.fn(deepWrapped))
			}
		})
	}
}

// TestAdversarialM3_DeepWrapping_Sentinels verifies that all sentinels resolve
// correctly through 100 levels of wrapping.
func TestAdversarialM3_DeepWrapping_Sentinels(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		sentinel error
		pred     func(error) bool
	}{
		// project
		{"project.ErrWorkspaceNotFound", project.ErrWorkspaceNotFound, project.IsNotFound},
		{"project.ErrConfigNotFound", project.ErrConfigNotFound, project.IsNotFound},
		{"project.ErrContractNotFound", project.ErrContractNotFound, project.IsNotFound},
		{"project.os.ErrNotExist", os.ErrNotExist, project.IsNotFound},
		{"project.ErrStaleCodegen", project.ErrStaleCodegen, project.IsStale},

		// parser
		{"parser.ErrSyntaxError", parser.ErrSyntaxError, parser.IsSyntaxError},
		{"parser.ErrContractNotFound", parser.ErrContractNotFound, parser.IsNotFound},

		// diff
		{"diff.ErrFrameNotFound", diff.ErrFrameNotFound, diff.IsNotFound},
		{"diff.ErrConflict", diff.ErrConflict, diff.IsConflict},
		{"diff.ErrStackEmpty", diff.ErrStackEmpty, diff.IsEmpty},

		// git
		{"git.ErrNotRepository", git.ErrNotRepository, git.IsNotRepository},
		{"git.ErrBranchNotFound", git.ErrBranchNotFound, git.IsNotFound},

		// cache
		{"cache.ErrSessionNotFound", cache.ErrSessionNotFound, cache.IsNotFound},
		{"cache.ErrSecretNotFound", cache.ErrSecretNotFound, cache.IsNotFound},
		{"cache.os.ErrNotExist", os.ErrNotExist, cache.IsNotFound},
		{"cache.ErrCorruptCache", cache.ErrCorruptCache, cache.IsCorrupt},

		// lint
		{"lint.ErrLintFailure", lint.ErrLintFailure, lint.IsLintFailure},
		{"lint.ErrRuleNotFound", lint.ErrRuleNotFound, lint.IsNotFound},

		// spec
		{"spec.ErrSpecNotFound", spec.ErrSpecNotFound, spec.IsNotFound},
		{"spec.os.ErrNotExist", os.ErrNotExist, spec.IsNotFound},
		{"spec.ErrUnsupportedFormat", spec.ErrUnsupportedFormat, spec.IsUnsupportedFormat},

		// pipeline
		{"pipeline.ErrPipelineAborted", pipeline.ErrPipelineAborted, pipeline.IsPipelineAborted},
		{"pipeline.ErrNoContractsFound", pipeline.ErrNoContractsFound, pipeline.IsNotFound},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for _, depth := range []int{1, 5, 10, 50, 100} {
				w := wrapN(tc.sentinel, depth)
				require.True(t, tc.pred(w), fmt.Sprintf("expected true at depth %d for %s", depth, tc.name))
			}
		})
	}
}

// TestAdversarialM3_DeepWrapping_TypedStructs verifies that typed error structs
// containing underlying sentinels resolve correctly even after 100 levels of wrapping.
func TestAdversarialM3_DeepWrapping_TypedStructs(t *testing.T) {
	t.Parallel()

	structCases := []struct {
		name string
		err  error
		pred func(error) bool
	}{
		{
			name: "project.ProjectError",
			err:  &project.ProjectError{Op: "load", Err: project.ErrWorkspaceNotFound},
			pred: project.IsNotFound,
		},
		{
			name: "project.ProjectError_Stale",
			err:  &project.ProjectError{Op: "status", Err: project.ErrStaleCodegen},
			pred: project.IsStale,
		},
		{
			name: "parser.ParseError_Syntax",
			err:  &parser.ParseError{Op: "parse", Err: parser.ErrSyntaxError},
			pred: parser.IsSyntaxError,
		},
		{
			name: "parser.ParseError_NotFound",
			err:  &parser.ParseError{Op: "bind", Err: parser.ErrContractNotFound},
			pred: parser.IsNotFound,
		},
		{
			name: "diff.DiffError_NotFound",
			err:  &diff.DiffError{Op: "pop", Err: diff.ErrFrameNotFound},
			pred: diff.IsNotFound,
		},
		{
			name: "diff.DiffError_Conflict",
			err:  &diff.DiffError{Op: "diff", Err: diff.ErrConflict},
			pred: diff.IsConflict,
		},
		{
			name: "diff.DiffError_Empty",
			err:  &diff.DiffError{Op: "peek", Err: diff.ErrStackEmpty},
			pred: diff.IsEmpty,
		},
		{
			name: "git.GitError_NotRepo",
			err:  &git.GitError{Op: "status", Err: git.ErrNotRepository},
			pred: git.IsNotRepository,
		},
		{
			name: "git.GitError_NotFound",
			err:  &git.GitError{Op: "branch", Err: git.ErrBranchNotFound},
			pred: git.IsNotFound,
		},
		{
			name: "cache.CacheError_NotFound",
			err:  &cache.CacheError{Op: "get", Err: cache.ErrSessionNotFound},
			pred: cache.IsNotFound,
		},
		{
			name: "cache.CacheError_Corrupt",
			err:  &cache.CacheError{Op: "load", Err: cache.ErrCorruptCache},
			pred: cache.IsCorrupt,
		},
		{
			name: "lint.LintError_Failure",
			err:  &lint.LintError{Op: "run", Err: lint.ErrLintFailure},
			pred: lint.IsLintFailure,
		},
		{
			name: "lint.LintError_NotFound",
			err:  &lint.LintError{Op: "check", Err: lint.ErrRuleNotFound},
			pred: lint.IsNotFound,
		},
		{
			name: "spec.SpecError_NotFound",
			err:  &spec.SpecError{Op: "read", Err: spec.ErrSpecNotFound},
			pred: spec.IsNotFound,
		},
		{
			name: "spec.SpecError_Unsupported",
			err:  &spec.SpecError{Op: "parse", Err: spec.ErrUnsupportedFormat},
			pred: spec.IsUnsupportedFormat,
		},
		{
			name: "pipeline.PipelineError_Aborted",
			err:  &pipeline.PipelineError{Op: "run", Err: pipeline.ErrPipelineAborted},
			pred: pipeline.IsPipelineAborted,
		},
		{
			name: "pipeline.PipelineError_NotFound",
			err:  &pipeline.PipelineError{Op: "plan", Err: pipeline.ErrNoContractsFound},
			pred: pipeline.IsNotFound,
		},
	}

	for _, sc := range structCases {
		sc := sc
		t.Run(sc.name, func(t *testing.T) {
			for _, depth := range []int{1, 5, 10, 50, 100} {
				w := wrapN(sc.err, depth)
				require.True(t, sc.pred(w), fmt.Sprintf("expected true at depth %d for %s", depth, sc.name))
			}
		})
	}
}

// TestAdversarialM3_CrossSubsystemWrapping tests deeply nested errors across multiple distinct subsystems.
func TestAdversarialM3_CrossSubsystemWrapping(t *testing.T) {
	t.Parallel()

	// PipelineError -> LintError -> ParseError -> ErrSyntaxError
	chain1 := &pipeline.PipelineError{
		Op: "compile",
		Err: &lint.LintError{
			Op: "validate",
			Err: &parser.ParseError{
				Op:  "parse_source",
				Err: parser.ErrSyntaxError,
			},
		},
	}

	// Must resolve innermost syntax error
	require.True(t, parser.IsSyntaxError(chain1))
	// PipelineError does NOT have ErrPipelineAborted
	require.False(t, pipeline.IsPipelineAborted(chain1))
	// LintError does NOT have ErrLintFailure
	require.False(t, lint.IsLintFailure(chain1))

	// ProjectError -> DiffError -> GitError -> ErrNotRepository
	chain2 := &project.ProjectError{
		Op: "sync",
		Err: &diff.DiffError{
			Op: "stack_diff",
			Err: &git.GitError{
				Op:  "rev_parse",
				Err: git.ErrNotRepository,
			},
		},
	}

	require.True(t, git.IsNotRepository(chain2))
	require.False(t, diff.IsConflict(chain2))
	require.False(t, project.IsNotFound(chain2))

	// All 8 subsystems chained sequentially
	chainAll := &project.ProjectError{
		Op: "all_0",
		Err: &cache.CacheError{
			Op: "all_1",
			Err: &diff.DiffError{
				Op: "all_2",
				Err: &git.GitError{
					Op: "all_3",
					Err: &lint.LintError{
						Op: "all_4",
						Err: &parser.ParseError{
							Op: "all_5",
							Err: &pipeline.PipelineError{
								Op: "all_6",
								Err: &spec.SpecError{
									Op:  "all_7",
									Err: spec.ErrUnsupportedFormat,
								},
							},
						},
					},
				},
			},
		},
	}

	require.True(t, spec.IsUnsupportedFormat(chainAll))
	require.False(t, spec.IsNotFound(chainAll))
	require.False(t, pipeline.IsPipelineAborted(chainAll))
	require.False(t, parser.IsSyntaxError(chainAll))
	require.False(t, lint.IsLintFailure(chainAll))
	require.False(t, git.IsNotRepository(chainAll))
	require.False(t, diff.IsConflict(chainAll))
	require.False(t, cache.IsCorrupt(chainAll))
	require.False(t, project.IsStale(chainAll))
}

// TestAdversarialM3_NilReceiverMethods verifies .Error() and .Unwrap() methods on typed nil receivers.
func TestAdversarialM3_NilReceiverMethods(t *testing.T) {
	t.Parallel()

	var pProject *project.ProjectError
	var pParser *parser.ParseError
	var pDiff *diff.DiffError
	var pGit *git.GitError
	var pCache *cache.CacheError
	var pLint *lint.LintError
	var pSpec *spec.SpecError
	var pPipeline *pipeline.PipelineError

	assert.Equal(t, "<nil>", pProject.Error())
	assert.Nil(t, pProject.Unwrap())

	assert.Equal(t, "<nil>", pParser.Error())
	assert.Nil(t, pParser.Unwrap())

	assert.Equal(t, "<nil>", pDiff.Error())
	assert.Nil(t, pDiff.Unwrap())

	assert.Equal(t, "<nil>", pGit.Error())
	assert.Nil(t, pGit.Unwrap())

	assert.Equal(t, "<nil>", pCache.Error())
	assert.Nil(t, pCache.Unwrap())

	assert.Equal(t, "<nil>", pLint.Error())
	assert.Nil(t, pLint.Unwrap())

	assert.Equal(t, "<nil>", pSpec.Error())
	assert.Nil(t, pSpec.Unwrap())

	assert.Equal(t, "<nil>", pPipeline.Error())
	assert.Nil(t, pPipeline.Unwrap())
}

// TestAdversarialM3_StructWithNilErr verifies behavior when struct fields are populated but Err is nil.
func TestAdversarialM3_StructWithNilErr(t *testing.T) {
	t.Parallel()

	preds := getAll17Predicates()

	emptyStructs := []error{
		&project.ProjectError{Op: "op", Path: "/path", Key: "k", Err: nil},
		&parser.ParseError{Op: "op", Path: "/path", Key: "k", Err: nil},
		&diff.DiffError{Op: "op", Path: "/path", Key: "k", Err: nil},
		&git.GitError{Op: "op", Path: "/path", Key: "k", Err: nil},
		&cache.CacheError{Op: "op", Path: "/path", Key: "k", Err: nil},
		&lint.LintError{Op: "op", Path: "/path", Key: "k", Err: nil},
		&spec.SpecError{Op: "op", Path: "/path", Key: "k", Err: nil},
		&pipeline.PipelineError{Op: "op", Path: "/path", Key: "k", Err: nil},
	}

	for _, s := range emptyStructs {
		for _, p := range preds {
			require.False(t, p.fn(s))
		}
	}
}
