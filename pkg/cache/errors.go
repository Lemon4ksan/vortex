// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache

import (
	"errors"
	"os"
	"strings"
)

var (
	// ErrSessionNotFound indicates that a recorded traffic capture session ID was not found in the index.
	ErrSessionNotFound = errors.New("cache: session not found")

	// ErrSecretNotFound indicates that a requested secret rule or variable was not found in the vault.
	ErrSecretNotFound = errors.New("cache: secret not found")

	// ErrCorruptCache indicates that a cached artifact or index file is corrupt or has an invalid schema.
	ErrCorruptCache = errors.New("cache: corrupt cache data")
)

// CacheError represents an error during cache read, write, or index lookup.
type CacheError struct {
	Op   string // Operation (e.g. "load_session", "store_session", "load_vault", "save_lint")
	Path string // Cache file or storage directory path
	Key  string // Session ID, hash, or secret key
	Err  error  // Underlying IO or sentinel error
}

func (e *CacheError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString("cache")
	if e.Op != "" {
		sb.WriteByte('/')
		sb.WriteString(e.Op)
	}
	sb.WriteString(": ")

	if e.Path != "" {
		sb.WriteString(e.Path)
		if e.Key != "" {
			sb.WriteString(" [key=")
			sb.WriteString(e.Key)
			sb.WriteByte(']')
		}
		sb.WriteString(": ")
	} else if e.Key != "" {
		sb.WriteString("key ")
		sb.WriteString(e.Key)
		sb.WriteString(": ")
	}

	if e.Err != nil {
		sb.WriteString(e.Err.Error())
	} else {
		sb.WriteString("unknown cache error")
	}

	return sb.String()
}

func (e *CacheError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsNotFound reports whether err indicates a missing session, secret, or cache entry.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrSessionNotFound) ||
		errors.Is(err, ErrSecretNotFound) ||
		errors.Is(err, os.ErrNotExist) {
		return true
	}

	if cErr, ok := errors.AsType[*CacheError](err); ok && cErr != nil {
		return errors.Is(cErr.Err, ErrSessionNotFound) ||
			errors.Is(cErr.Err, ErrSecretNotFound) ||
			errors.Is(cErr.Err, os.ErrNotExist)
	}

	return strings.Contains(err.Error(), "not found in cache")
}

// IsCorrupt reports whether err indicates corrupt cache data.
func IsCorrupt(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrCorruptCache) {
		return true
	}

	if cErr, ok := errors.AsType[*CacheError](err); ok && cErr != nil {
		return errors.Is(cErr.Err, ErrCorruptCache)
	}

	return false
}
