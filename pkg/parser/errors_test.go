// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lemon4ksan/foundation/testing/assert"
	"github.com/lemon4ksan/foundation/testing/require"

	"github.com/lemon4ksan/vortex/pkg/parser"
)

func TestParserErrors_IsSyntaxError(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, parser.IsSyntaxError(parser.ErrSyntaxError))

	// 2. Wrapped sentinel
	require.True(t, parser.IsSyntaxError(fmt.Errorf("wrap: %w", parser.ErrSyntaxError)))

	// 3. Typed struct
	pErr := &parser.ParseError{
		Op:   "parse_file",
		Path: "service.go",
		Key:  "12:5",
		Err:  parser.ErrSyntaxError,
	}
	require.True(t, parser.IsSyntaxError(pErr))

	pErrWrapped := fmt.Errorf("outer: %w", pErr)
	require.True(t, parser.IsSyntaxError(pErrWrapped))

	// 3b. String fallback match
	require.True(t, parser.IsSyntaxError(errors.New("syntax error at EOF")))

	// 4. Negative case
	require.False(t, parser.IsSyntaxError(errors.New("unrelated io error")))
	require.False(t, parser.IsSyntaxError(parser.ErrContractNotFound))

	// 5. Nil error
	require.False(t, parser.IsSyntaxError(nil))
	var typedNil *parser.ParseError
	require.False(t, parser.IsSyntaxError(typedNil))
	require.False(t, parser.IsSyntaxError(fmt.Errorf("wrap: %w", typedNil)))
}

func TestParserErrors_IsNotFound(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, parser.IsNotFound(parser.ErrContractNotFound))

	// 2. Wrapped sentinel
	require.True(t, parser.IsNotFound(fmt.Errorf("wrap: %w", parser.ErrContractNotFound)))

	// 3. Typed struct
	pErr := &parser.ParseError{
		Op:   "bind",
		Path: "models.go",
		Err:  parser.ErrContractNotFound,
	}
	require.True(t, parser.IsNotFound(pErr))

	// 4. Negative case
	require.False(t, parser.IsNotFound(errors.New("unrelated error")))
	require.False(t, parser.IsNotFound(parser.ErrInvalidDirective))

	// 5. Nil error
	require.False(t, parser.IsNotFound(nil))
	var typedNil *parser.ParseError
	require.False(t, parser.IsNotFound(typedNil))
	require.False(t, parser.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
}

func TestParseError_ErrorAndUnwrap(t *testing.T) {
	t.Parallel()

	// Nil receiver
	var nilErr *parser.ParseError
	assert.Equal(t, "<nil>", nilErr.Error())
	assert.Nil(t, nilErr.Unwrap())

	// Op, Path, Key, Err
	e1 := &parser.ParseError{
		Op:   "parse_file",
		Path: "api.go",
		Key:  "42:10",
		Err:  parser.ErrSyntaxError,
	}
	assert.Equal(t, "parser/parse_file: api.go:42:10: parser: syntax error", e1.Error())
	assert.Equal(t, parser.ErrSyntaxError, e1.Unwrap())
	require.True(t, errors.Is(e1, parser.ErrSyntaxError))

	// Op, Path, Err (no Key)
	e2 := &parser.ParseError{
		Op:   "bind_type",
		Path: "schema.go",
		Err:  parser.ErrUnresolvedType,
	}
	assert.Equal(t, "parser/bind_type: schema.go: parser: unresolved type", e2.Error())

	// Op, Key, Err (no Path)
	e3 := &parser.ParseError{
		Op:  "parse_directive",
		Key: "@unknown",
		Err: parser.ErrInvalidDirective,
	}
	assert.Equal(t, "parser/parse_directive: @unknown: parser: invalid directive", e3.Error())

	// Op, nil Err
	e4 := &parser.ParseError{
		Op: "scan",
	}
	assert.Equal(t, "parser/scan: unknown parse error", e4.Error())
}
