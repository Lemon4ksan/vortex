// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/lemon4ksan/foundation/testing/assert"
	"github.com/lemon4ksan/foundation/testing/require"

	"github.com/lemon4ksan/vortex/pkg/cache"
)

func TestCacheErrors_IsNotFound(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, cache.IsNotFound(cache.ErrSessionNotFound))
	require.True(t, cache.IsNotFound(cache.ErrSecretNotFound))
	require.True(t, cache.IsNotFound(os.ErrNotExist))

	// 2. Wrapped sentinel
	require.True(t, cache.IsNotFound(fmt.Errorf("wrap: %w", cache.ErrSessionNotFound)))
	require.True(t, cache.IsNotFound(fmt.Errorf("wrap: %w", cache.ErrSecretNotFound)))
	require.True(t, cache.IsNotFound(fmt.Errorf("wrap: %w", os.ErrNotExist)))

	// 3. Typed struct
	cErr := &cache.CacheError{
		Op:  "load_session",
		Key: "sess-123",
		Err: cache.ErrSessionNotFound,
	}
	require.True(t, cache.IsNotFound(cErr))

	cErrWrapped := fmt.Errorf("outer: %w", cErr)
	require.True(t, cache.IsNotFound(cErrWrapped))

	// 3b. String fallback match
	require.True(t, cache.IsNotFound(errors.New("vortex: traffic session not found in cache")))

	// 4. Negative case
	require.False(t, cache.IsNotFound(errors.New("unrelated error")))
	require.False(t, cache.IsNotFound(cache.ErrCorruptCache))

	// 5. Nil error
	require.False(t, cache.IsNotFound(nil))
	var typedNil *cache.CacheError
	require.False(t, cache.IsNotFound(typedNil))
	require.False(t, cache.IsNotFound(fmt.Errorf("wrap: %w", typedNil)))
}

func TestCacheErrors_IsCorrupt(t *testing.T) {
	t.Parallel()

	// 1. Direct sentinel
	require.True(t, cache.IsCorrupt(cache.ErrCorruptCache))

	// 2. Wrapped sentinel
	require.True(t, cache.IsCorrupt(fmt.Errorf("wrap: %w", cache.ErrCorruptCache)))

	// 3. Typed struct
	cErr := &cache.CacheError{
		Op:   "load_vault",
		Path: ".vortex/secrets.vault",
		Err:  cache.ErrCorruptCache,
	}
	require.True(t, cache.IsCorrupt(cErr))

	// 4. Negative case
	require.False(t, cache.IsCorrupt(errors.New("unrelated error")))
	require.False(t, cache.IsCorrupt(cache.ErrSessionNotFound))

	// 5. Nil error
	require.False(t, cache.IsCorrupt(nil))
	var typedNil *cache.CacheError
	require.False(t, cache.IsCorrupt(typedNil))
	require.False(t, cache.IsCorrupt(fmt.Errorf("wrap: %w", typedNil)))
}

func TestCacheError_ErrorAndUnwrap(t *testing.T) {
	t.Parallel()

	// Nil receiver
	var nilErr *cache.CacheError
	assert.Equal(t, "<nil>", nilErr.Error())
	assert.Nil(t, nilErr.Unwrap())

	// Op, Path, Key, Err
	e1 := &cache.CacheError{
		Op:   "load_session",
		Path: ".vortex/traffic/idx",
		Key:  "hash-abc",
		Err:  cache.ErrSessionNotFound,
	}
	assert.Equal(t, "cache/load_session: .vortex/traffic/idx [key=hash-abc]: cache: session not found", e1.Error())
	assert.Equal(t, cache.ErrSessionNotFound, e1.Unwrap())
	require.True(t, errors.Is(e1, cache.ErrSessionNotFound))

	// Op, Path, Err (no Key)
	e2 := &cache.CacheError{
		Op:   "save_lint",
		Path: ".vortex/lint.json",
		Err:  cache.ErrCorruptCache,
	}
	assert.Equal(t, "cache/save_lint: .vortex/lint.json: cache: corrupt cache data", e2.Error())

	// Op, Key, Err (no Path)
	e3 := &cache.CacheError{
		Op:  "get_secret",
		Key: "API_KEY",
		Err: cache.ErrSecretNotFound,
	}
	assert.Equal(t, "cache/get_secret: key API_KEY: cache: secret not found", e3.Error())

	// Op, nil Err
	e4 := &cache.CacheError{
		Op: "flush",
	}
	assert.Equal(t, "cache/flush: unknown cache error", e4.Error())
}
