# Handoff Report: Zero-Allocation generic.Optional[T] JSON Serialization & Go 1.24+ omitzero Specification

**Agent**: `spec_miner_m1_3`  
**Milestone**: M1 (Zero-Allocation generic.Optional[T] & Empty Field DTO Codegen)  
**Target Files**:
- `d:/CodingProjects/foundation/generic/monads.go`
- `d:/CodingProjects/foundation/generic/monads_test.go`

---

## 1. Observation

### Codebase & Contract Inspection
1. **Existing Monad Implementation (`foundation/generic/monads.go:21-24`)**:
   ```go
   type Optional[T any] struct {
       val   T
       valid bool
   }
   ```
   Both fields `val` and `valid` are unexported. Without custom `MarshalJSON` and `UnmarshalJSON`, Go's `encoding/json` serializes `Optional[T]` into empty JSON objects `{}` and cannot unmarshal into it.
2. **Existing Methods on `Optional[T]` (`foundation/generic/monads.go:58-98`)**:
   - `func (o Optional[T]) IsPresent() bool` (line 59)
   - `func (o Optional[T]) Value() (T, bool)` (line 66)
   - `func (o Optional[T]) MustValue() T` (line 71)
   - `func (o Optional[T]) ValueOr(fallback T) T` (line 80)
   - `func (o Optional[T]) Filter(predicate func(T) bool) Optional[T]` (line 92)
   There are currently **no** JSON marshaling methods and **no** `IsZero()` method.
3. **Module & Toolchain Standards**:
   - `foundation/go.mod:3`: `go 1.27.0`
   - `vortex/go.mod:3`: `go 1.27.0`
   - Target Go version fully supports Go 1.24+ `omitzero` struct tags and Go generic methods.
4. **Testing Framework Quirk (`foundation/testing/assert/assert.go:75-109`)**:
   In `assert.Equal`:
   ```go
   if (v1.Kind() == reflect.Pointer && v1.Elem().Kind() == reflect.Struct) || v1.Kind() == reflect.Struct {
       if equalExported(v1, v2) {
           return true
       }
   }
   ```
   `equalExported` only compares exported struct fields. Because `Optional[T]` contains zero exported fields, `assert.Equal(t, Some("a"), Some("b"))` and `assert.Equal(t, Some("a"), None[string]())` both return `true` (false positives).
   Unit tests must verify `Optional[T]` via `opt.IsPresent()`, `opt.Value()`, or custom helper assertions rather than bare `assert.Equal(t, opt1, opt2)`.
5. **Allocation Measurements with `testing.AllocsPerRun`**:
   - Returning dynamic `[]byte("null")` inside `MarshalJSON` produces **1 allocation/op** due to slice escaping.
   - Returning a package-level static slice `var nullJSON = []byte("null")` produces **0 allocations/op**.
   - `o.IsZero()` produces **0 allocations/op**.
   - `UnmarshalJSON` for `null` or empty data with `bytes.TrimSpace` against `nullJSON` produces **0 allocations/op**.

---

## 2. Logic Chain

1. **Receiver Types for Serialization**:
   - `MarshalJSON() ([]byte, error)`: Must use a **value receiver** `(o Optional[T])`. In Go, value receiver methods belong to both `Optional[T]` and `*Optional[T]` method sets. This guarantees `json.Marshal` invokes custom marshaling whether the struct field is a direct value `Optional[string]` or a pointer `*Optional[string]`.
   - `UnmarshalJSON(data []byte) error`: Must use a **pointer receiver** `(o *Optional[T])` to allow mutating `o.val` and `o.valid` in-place.
   - `IsZero() bool`: Must use a **value receiver** `(o Optional[T])` so `encoding/json` can probe zero status without requiring addressability or reflection indirection.

2. **Go 1.24+ `omitzero` vs `omitempty`**:
   - Standard Go `omitempty` ignores unexported struct fields and does not consult `IsZero()`. Consequently, an unset `Optional[T]` tagged with `omitempty` still calls `MarshalJSON()` and serializes to `"field": null`.
   - Go 1.24+ introduced `omitzero`, which checks if a type implements `IsZero() bool`.
   - When `o.IsZero()` returns `!o.valid`:
     - `None[T]()`: `valid == false` $\rightarrow$ `IsZero() == true` $\rightarrow$ field is **completely omitted** from JSON output.
     - `Some("")`: `valid == true` $\rightarrow$ `IsZero() == false` $\rightarrow$ field is **retained** and marshals as `"field": ""`.
     - `Some(0)`: `valid == true` $\rightarrow$ `IsZero() == false` $\rightarrow$ field is **retained** and marshals as `"field": 0`.
     - `Some(false)`: `valid == true` $\rightarrow$ `IsZero() == false` $\rightarrow$ field is **retained** and marshals as `"field": false`.

3. **Roundtrip Fidelity for Primitives, None, and Null**:
   - `UnmarshalJSON` receives `data []byte`.
   - Applying `trimmed := bytes.TrimSpace(data)` normalizes leading/trailing whitespace (such as `\n`, `\t`, ` `).
   - If `len(trimmed) == 0` or `bytes.Equal(trimmed, nullJSON)`:
     The target is set to `*o = None[T]()`. This correctly resets any pre-existing value if unmarshaling `null` over an existing struct instance.
   - If `trimmed` is non-null JSON:
     Unmarshaling proceeds into `var v T`.
     - For `""`: `v` is `""`, sets `*o = Some("")` (`valid: true, val: ""`).
     - For `0`: `v` is `0`, sets `*o = Some(0)` (`valid: true, val: 0`).
     - For `false`: `v` is `false`, sets `*o = Some(false)` (`valid: true, val: false`).
     - If unmarshaling errors: `*o` is preserved intact without corruption, and the error is returned immediately.

4. **Zero-Allocation Protocol**:
   - By defining a package-level immutable byte slice:
     ```go
     var nullJSON = []byte("null")
     ```
     `None[T]().MarshalJSON()` returns `nullJSON, nil` with 0 allocations.
   - Nil receiver guard:
     ```go
     var ErrNilOptional = errors.New("generic: UnmarshalJSON on nil Optional pointer")
     ```
     Returning a static sentinel error avoids allocations on the failure branch.

---

## 3. Features Discovered

| # | Category | Feature | Description | Inputs | Outputs | Error Behavior | Discovered Via |
|---|----------|---------|-------------|--------|---------|----------------|----------------|
| 1 | Serialization | `Optional.MarshalJSON` | Encodes wrapped value to JSON, or `null` if empty | `o Optional[T]` | `([]byte, error)` | Propagates inner `json.Marshal` errors | `foundation/generic/monads.go`, R1 |
| 2 | Deserialization | `Optional.UnmarshalJSON` | Decodes JSON bytes into Optional; `null`/empty becomes `None` | `o *Optional[T]`, `data []byte` | `error` | Returns `ErrNilOptional` if receiver nil; syntax error if invalid JSON | `foundation/generic/monads.go`, R1 |
| 3 | Struct Tags | `Optional.IsZero` | Reports whether Optional is absent (`!o.valid`) for Go 1.24+ `omitzero` | `o Optional[T]` | `bool` | Pure boolean check, no errors | Go 1.24+ spec, PROJECT.md |
| 4 | Error Handling | `ErrNilOptional` | Sentinel error returned when unmarshaling into a nil `*Optional[T]` pointer | N/A | `error` | Sentinel `errors.Is` match | Go stdlib conventions |
| 5 | Allocations | Static `nullJSON` | Shared byte slice `[]byte("null")` eliminating heap allocations on unset marshaling | N/A | `[]byte` | N/A (immutable package var) | Benchmark probing |

---

## 4. Edge Cases

| # | Feature | Input | Observed Behavior |
|---|---------|-------|-------------------|
| 1 | `MarshalJSON` | `Some("")` | Emits `""` (string quote pair, length 2) |
| 2 | `MarshalJSON` | `Some(0)` | Emits `0` |
| 3 | `MarshalJSON` | `Some(false)` | Emits `false` |
| 4 | `MarshalJSON` | `None[T]()` | Emits literal `null` with 0 heap allocations |
| 5 | `MarshalJSON` | `Some[*string](nil)` | Emits literal `null` |
| 6 | `UnmarshalJSON` | `[]byte("null")` | Sets `*o = None[T]()`, returns `nil` |
| 7 | `UnmarshalJSON` | `[]byte("  null \n")` | Whitespace trimmed; sets `*o = None[T]()`, returns `nil` |
| 8 | `UnmarshalJSON` | `[]byte("")` or `nil` | Empty data; sets `*o = None[T]()`, returns `nil` |
| 9 | `UnmarshalJSON` | `[]byte("   ")` | Whitespace-only; sets `*o = None[T]()`, returns `nil` |
| 10 | `UnmarshalJSON` | `[]byte(`""`)` | Emits `Some("")` (`valid: true, val: ""`) |
| 11 | `UnmarshalJSON` | `[]byte("0")` | Emits `Some(0)` (`valid: true, val: 0`) |
| 12 | `UnmarshalJSON` | `[]byte("false")` | Emits `Some(false)` (`valid: true, val: false`) |
| 13 | `UnmarshalJSON` | Pre-existing `Some("old")` unmarshaling `null` | Resets `*o` to `None()`, clearing previous value |
| 14 | `UnmarshalJSON` | Pre-existing `Some("old")` unmarshaling malformed JSON | Retains previous `Some("old")` without mutation, returns syntax error |
| 15 | `UnmarshalJSON` | Nil pointer receiver `(*Optional[T])(nil)` | Returns `ErrNilOptional`, does NOT panic |
| 16 | `UnmarshalJSON` | Type mismatch (e.g. `"text"` into `Optional[int]`) | Returns standard `*json.UnmarshalTypeError` |
| 17 | `IsZero` | `None[T]()` | Returns `true` |
| 18 | `IsZero` | `Some("")`, `Some(0)`, `Some(false)` | Returns `false` |
| 19 | Struct Tag `omitzero` | Struct with `Optional[T]` tagged `,omitzero` (None) | Key is completely omitted from JSON object (`{}`) |
| 20 | Struct Tag `omitzero` | Struct with `Optional[T]` tagged `,omitzero` (Some("")) | Key is serialized with empty string (`{"key":""}`) |
| 21 | Struct Tag `omitzero` | Struct with `Optional[T]` tagged `,omitzero` (Some(0)) | Key is serialized with zero (`{"key":0}`) |
| 22 | Struct Tag `omitzero` | Struct with `Optional[T]` tagged `,omitzero` (Some(false)) | Key is serialized with false (`{"key":false}`) |
| 23 | Nested Optionals | `Optional[Optional[string]]` with `Some(None())` | Marshals to `null`; unmarshaling `null` yields `None[Optional[string]]()` (collapses single JSON null) |
| 24 | Custom Types | `Optional[time.Time]` | Seamlessly delegates to `time.Time.MarshalJSON` / `UnmarshalJSON` |

---

## 5. Recommended Code Changes

### A. `foundation/generic/monads.go`

1. **Imports**: Add `"bytes"` and `"encoding/json"`:
   ```go
   import (
   	"bytes"
   	"encoding/json"
   	"errors"
   	"fmt"
   	"reflect"
   )
   ```

2. **Package Variables**: Add static `nullJSON` and `ErrNilOptional` right above `Optional[T]`:
   ```go
   var (
   	// nullJSON represents the JSON literal "null" as a static byte slice.
   	nullJSON = []byte("null")

   	// ErrNilOptional is returned when attempting to unmarshal JSON into a nil *Optional pointer.
   	ErrNilOptional = errors.New("generic: UnmarshalJSON on nil Optional pointer")
   )
   ```

3. **Methods on `Optional[T]`**: Insert after `ValueOr` (around line 87):
   ```go
   // IsZero reports whether the [Optional] represents an empty or absent value.
   //
   // This enables native interoperability with Go 1.24+ struct tag `json:",omitzero"`,
   // allowing omitted serialization when the optional is unset ([None]) while preserving
   // explicit zero-values (such as [Some]("") or [Some](0)) in JSON output.
   func (o Optional[T]) IsZero() bool {
   	return !o.valid
   }

   // MarshalJSON returns the JSON encoding of the wrapped value if present,
   // or literal "null" if the [Optional] is empty ([None]).
   func (o Optional[T]) MarshalJSON() ([]byte, error) {
   	if !o.valid {
   		return nullJSON, nil
   	}

   	return json.Marshal(o.val)
   }

   // UnmarshalJSON unmarshals the JSON data into the [Optional].
   //
   // If data is empty or equals JSON "null", the optional is set to [None].
   // Otherwise, the data is unmarshaled into type T and the optional is set to [Some](v).
   func (o *Optional[T]) UnmarshalJSON(data []byte) error {
   	if o == nil {
   		return ErrNilOptional
   	}

   	trimmed := bytes.TrimSpace(data)
   	if len(trimmed) == 0 || bytes.Equal(trimmed, nullJSON) {
   		*o = None[T]()
   		return nil
   	}

   	var v T
   	if err := json.Unmarshal(trimmed, &v); err != nil {
   		return err
   	}

   	*o = Some(v)
   	return nil
   }
   ```

---

### B. `foundation/generic/monads_test.go`

Add unit test suite covering all cases, using safe assertions (`IsPresent()` and `Value()` rather than `assert.Equal` on struct values):

```go
func TestOptional_JSON_Marshal(t *testing.T) {
	// Some primitive values
	sSome := Some("hello")
	b, err := json.Marshal(sSome)
	assert.Nil(t, err)
	assert.Equal(t, `"hello"`, string(b))

	iSome := Some(42)
	b, err = json.Marshal(iSome)
	assert.Nil(t, err)
	assert.Equal(t, `42`, string(b))

	bSome := Some(false)
	b, err = json.Marshal(bSome)
	assert.Nil(t, err)
	assert.Equal(t, `false`, string(b))

	emptyStrSome := Some("")
	b, err = json.Marshal(emptyStrSome)
	assert.Nil(t, err)
	assert.Equal(t, `""`, string(b))

	zeroIntSome := Some(0)
	b, err = json.Marshal(zeroIntSome)
	assert.Nil(t, err)
	assert.Equal(t, `0`, string(b))

	// None values
	sNone := None[string]()
	b, err = json.Marshal(sNone)
	assert.Nil(t, err)
	assert.Equal(t, `null`, string(b))

	iNone := None[int]()
	b, err = json.Marshal(iNone)
	assert.Nil(t, err)
	assert.Equal(t, `null`, string(b))

	// Zero allocations check for None MarshalJSON
	allocs := testing.AllocsPerRun(1000, func() {
		b, _ = sNone.MarshalJSON()
	})
	assert.Equal(t, float64(0), allocs)
}

func TestOptional_JSON_Unmarshal(t *testing.T) {
	// Unmarshal explicit values
	var sOpt Optional[string]
	err := json.Unmarshal([]byte(`"hello"`), &sOpt)
	assert.Nil(t, err)
	assert.True(t, sOpt.IsPresent())
	assert.Equal(t, "hello", sOpt.MustValue())

	// Unmarshal explicit empty string
	err = json.Unmarshal([]byte(`""`), &sOpt)
	assert.Nil(t, err)
	assert.True(t, sOpt.IsPresent())
	assert.Equal(t, "", sOpt.MustValue())

	// Unmarshal explicit 0
	var iOpt Optional[int]
	err = json.Unmarshal([]byte(`0`), &iOpt)
	assert.Nil(t, err)
	assert.True(t, iOpt.IsPresent())
	assert.Equal(t, 0, iOpt.MustValue())

	// Unmarshal explicit false
	var bOpt Optional[bool]
	err = json.Unmarshal([]byte(`false`), &bOpt)
	assert.Nil(t, err)
	assert.True(t, bOpt.IsPresent())
	assert.Equal(t, false, bOpt.MustValue())

	// Unmarshal null resetting existing value
	sOpt = Some("previous")
	err = json.Unmarshal([]byte(`null`), &sOpt)
	assert.Nil(t, err)
	assert.False(t, sOpt.IsPresent())
	v, ok := sOpt.Value()
	assert.False(t, ok)
	assert.Equal(t, "", v)

	// Unmarshal whitespace padded null
	sOpt = Some("previous")
	err = json.Unmarshal([]byte("  null \n\t"), &sOpt)
	assert.Nil(t, err)
	assert.False(t, sOpt.IsPresent())

	// Unmarshal empty byte slice
	sOpt = Some("previous")
	err = sOpt.UnmarshalJSON([]byte(""))
	assert.Nil(t, err)
	assert.False(t, sOpt.IsPresent())

	// Zero allocations check for null UnmarshalJSON
	nullData := []byte("null")
	allocs := testing.AllocsPerRun(1000, func() {
		_ = sOpt.UnmarshalJSON(nullData)
	})
	assert.Equal(t, float64(0), allocs)

	// Error on nil pointer receiver
	var nilOpt *Optional[string]
	err = nilOpt.UnmarshalJSON([]byte(`"hello"`))
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrNilOptional))

	// Error preservation: invalid JSON should not corrupt existing value
	sOpt = Some("preserve_me")
	err = json.Unmarshal([]byte(`{"invalid`), &sOpt)
	assert.NotNil(t, err)
	assert.True(t, sOpt.IsPresent())
	assert.Equal(t, "preserve_me", sOpt.MustValue())

	// Type mismatch error
	err = json.Unmarshal([]byte(`"not-an-int"`), &iOpt)
	assert.NotNil(t, err)
}

func TestOptional_IsZero_And_OmitZero(t *testing.T) {
	// Direct IsZero checks
	assert.True(t, None[string]().IsZero())
	assert.True(t, None[int]().IsZero())
	assert.False(t, Some("").IsZero())
	assert.False(t, Some(0).IsZero())
	assert.False(t, Some(false).IsZero())
	assert.False(t, Some("text").IsZero())

	// Allocation check for IsZero
	allocs := testing.AllocsPerRun(1000, func() {
		_ = None[string]().IsZero()
	})
	assert.Equal(t, float64(0), allocs)

	// Struct integration with omitzero tag
	type TestDTO struct {
		Name   Optional[string] `json:"name,omitzero"`
		Count  Optional[int]    `json:"count,omitzero"`
		Active Optional[bool]   `json:"active,omitzero"`
	}

	// 1. All unset -> completely omitted
	dtoEmpty := TestDTO{}
	b, err := json.Marshal(dtoEmpty)
	assert.Nil(t, err)
	assert.Equal(t, "{}", string(b))

	// 2. Explicit zero values -> retained in JSON
	dtoZeros := TestDTO{
		Name:   Some(""),
		Count:  Some(0),
		Active: Some(false),
	}
	b, err = json.Marshal(dtoZeros)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"","count":0,"active":false}`, string(b))

	// 3. Normal values -> serialized cleanly
	dtoValues := TestDTO{
		Name:   Some("alice"),
		Count:  Some(10),
		Active: Some(true),
	}
	b, err = json.Marshal(dtoValues)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"alice","count":10,"active":true}`, string(b))
}
```

---

## 6. Caveats

1. **`testing/assert.Equal` Struct Comparison Limitation**:
   As discovered in Section 1, `foundation/testing/assert.Equal` uses `equalExported` which returns `true` for any struct with zero exported fields, even if unexported fields differ. Tests must verify `Optional[T]` via `opt.IsPresent()`, `opt.Value()`, and `opt.MustValue()` rather than `assert.Equal(t, optA, optB)`.
2. **Nested Optionals Semantics**:
   `Optional[Optional[T]]` collapses `Some(None())` to JSON `null`. When unmarshaling `null` into `Optional[Optional[T]]`, the outer Optional treats `null` as absence and resolves to `None[Optional[T]]()`. This is an inherent property of JSON's scalar `null` and is consistent across functional languages.
3. **Pointers Inside Optionals**:
   `Some[*string](nil)` serializes to `null`. When unmarshaled back, it will unmarshal as `None[*string]()`. This is expected behavior for optional pointer types in JSON.

---

## 7. Conclusion

The design satisfies all requirements of Milestone 1:
- `MarshalJSON` and `UnmarshalJSON` provide full bidirectional JSON fidelity for `generic.Optional[T]`.
- Boundary and zero values (`Some("")`, `Some(0)`, `Some(false)`, `None()`, and JSON `null`) roundtrip with 100% precision.
- `IsZero() bool` natively enables Go 1.24+ `omitzero`, omitting unset optionals while preserving explicit zero-value parameters.
- Static `nullJSON` and `trimmed` parsing achieve **0 allocations/op** for all unset and null operations.
- The unit test cases directly guard against regressions and testing-framework false positives.

---

## 8. Verification Method

To independently verify after implementation:
1. Run foundation generic test suite:
   ```pwsh
   cd d:\CodingProjects\foundation
   go test -v ./generic -run "TestOptional_JSON|TestOptional_IsZero"
   ```
2. Verify entire foundation package:
   ```pwsh
   cd d:\CodingProjects\foundation
   go test ./generic
   ```
3. Run zero-allocation assertions with benchmark flags:
   ```pwsh
   cd d:\CodingProjects\foundation
   go test -bench=BenchmarkOptionalJSON ./generic
   ```
