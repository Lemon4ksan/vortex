# Handoff Report — Optional Monads & DTO Boundary Adversarial Testing

**Author**: `explorer_m4_2` (Optional Monads & Boundary Adversarial Test Explorer)  
**Date**: 2026-09-23T13:00:00Z  
**Target Milestone**: M4 (Performance Benchmarks & Adversarial Test Coverage)  
**Consumer**: `worker_m4` / `orchestrator_3`  

---

## 1. Observation

### 1.1 Foundation Monad Implementation (`foundation/generic/monads.go`)
- **`Optional[T]` Definition** (`lines 31-34`):
  ```go
  type Optional[T any] struct {
      val   T
      valid bool
  }
  ```
- **Go 1.24+ `omitzero` Support via `IsZero()`** (`lines 115-117`):
  ```go
  func (o Optional[T]) IsZero() bool {
      return !o.valid
  }
  ```
  Returns `true` when `o.valid == false` (`None[T]()`), enabling Go 1.24+ `encoding/json` with struct tag `json:",omitzero"` to omit the field. When `Some(...)` is present (even `Some("")`, `Some(0)`, `Some(false)`, `Some(0.0)`), `o.valid` is `true`, so `IsZero()` returns `false`, preventing premature omission.
- **`MarshalJSON()`** (`lines 121-127`):
  ```go
  func (o Optional[T]) MarshalJSON() ([]byte, error) {
      if !o.valid {
          return nullJSON, nil
      }
      return json.Marshal(o.val)
  }
  ```
  Empty/unset optional returns literal `[]byte("null")`.
- **`UnmarshalJSON()`** (`lines 133-151`):
  ```go
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
  Correctly checks for `o == nil`, trims whitespace, treats `null` and empty payload as `None[T]()`, and populates `Some(v)` only on valid JSON unmarshaling into type `T`. If `json.Unmarshal` fails, existing `*o` is not corrupted.

### 1.2 DTO Code Generator (`pkg/emitter/dto.go`)
- **Generic Optional unwrapping** (`lines 50-62`):
  ```go
  func unwrapOptionalType(f *ir.FieldIR) (innerType string, isOptional bool) {
      name := strings.TrimPrefix(f.Type.Name, "*")
      if strings.HasPrefix(name, "generic.Optional[") && strings.HasSuffix(name, "]") {
          return strings.TrimSuffix(strings.TrimPrefix(name, "generic.Optional["), "]"), true
      }
      if strings.HasPrefix(name, "Optional[") && strings.HasSuffix(name, "]") {
          return strings.TrimSuffix(strings.TrimPrefix(name, "Optional["), "]"), true
      }
      if f.Type.ElemType != "" && strings.Contains(name, "Optional[") {
          return f.Type.ElemType, true
      }
      return "", false
  }
  ```
- **Nil Receiver Safety** (`lines 23-24, 36-37`):
  Both `AppendFormData` and `EncodeValues` emit `if r == nil { return dst }` / `if r == nil { return }`.
- **Empty String vs None Serialization** (`lines 191-198`):
  ```go
  case "string":
      fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
      buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
      fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
      buf.WriteString("\t\tif optVal != \"\" {\n")
      buf.WriteString("\t\t\tdst = appendQueryEscape(dst, optVal)\n")
      buf.WriteString("\t\t}\n")
      buf.WriteString("\t}\n")
  ```
  `Some("")` enters the branch (`ok == true`), appends `&` if needed, and appends `f.WireName+"="` (e.g. `q=`). Because `optVal != ""` is false, it does not append escaped value, yielding exactly `q=`. Conversely, `None[string]()` skips the entire block (`ok == false`), emitting 0 bytes.
- **Query Escaping Helper** (`lines 495-510`):
  Iterates byte-by-byte without string slicing or allocation:
  - Space `' '` -> `'+'`
  - Unreserved `a-z`, `A-Z`, `0-9`, `-`, `_`, `.`, `~` -> directly appended
  - All other bytes (including high bytes `0x80-0xFF` from multi-byte UTF-8 sequences and control bytes) -> `%XX` formatted using static hex lookup table `"0123456789ABCDEF"`
- **Slice Collections in Optional** (`lines 276-294` & `lines 463-475`):
  `case "[]int", "[]int64":` and `case "[]string":` iterate each item and append repeated keys `wire=val1&wire=val2` in `AppendFormData` and call `vals.Add(wire, val)` in `EncodeValues`.

### 1.3 Existing Test Coverage Status
- **`pkg/emitter/dto_test.go`**:
  - `TestEmitter_DTO_Emission`: verifies emitted syntax strings.
  - `TestEmitter_DTO_ExecutionAndZeroAlloc`: compiles a temp module with `SearchFilter`, testing `None` omission, `Some("")`, basic primitives, and explicit zeroes.
- **`pkg/emitter/dto_bench_test.go`**:
  - `TestDTO_AllPrimitives_ExactWireFormat`: tests `AllPrimitivesDTO` with 21 fields.
  - `TestDTO_AllPrimitives_ZeroAllocations`: verifies `testing.AllocsPerRun == 0`.
  - `TestDTO_SomeEmptyString_Only`, `TestDTO_AllNone_ZeroBytes`, `TestDTO_BoundaryValuesAndZeroes`, `TestDTO_QueryEscapingSpecialChars`.
  - Benchmarks: `Benchmark_AppendFormData_*` with `b.ReportAllocs()`.
- **`foundation/generic/monads_adversarial_test.go`**:
  - `TestOptional_Adversarial_NilPointerReceiver`: nil pointer receiver `UnmarshalJSON`.
  - `TestOptional_Adversarial_CorruptedData`: 12 malformed JSON payloads.
  - `TestOptional_Adversarial_WhitespaceAndEmpty`: 11 whitespace/null variations.
  - `TestOptional_Adversarial_NestedStructsRoundtrip`: `AdversarialParent` / `AdversarialChild`.
  - `TestOptional_Adversarial_NestedOptional`: `Optional[Optional[string]]`.
  - `TestOptional_Adversarial_OmitZeroExhaustive`: `BoundaryModel` with `json:",omitzero"`.
  - `TestOptional_Adversarial_ConcurrencyStress`: 100 readers & 100 unmarshal workers.

### 1.4 Observed Test Gaps
1. **DTO Nil Receiver Safety**: Neither `dto_test.go` nor `dto_bench_test.go` invokes `AppendFormData`, `AppendQuery`, or `EncodeValues` on a typed nil pointer receiver `(*DTO)(nil)`.
2. **Buffer Capacity Edge Cases & Realloc**: No tests verify behavior when pre-allocated slice capacity is zero (`cap == 0`), smaller than output (`cap == 5`), or oversized for reuse (`cap == 4096`).
3. **Pre-Populated Buffer Prefixing**: No tests verify that calling `AppendFormData` on a buffer that already contains data (e.g. `dst := []byte("existing=1")`) properly prepends `'&'` to the first emitted field (`existing=1&key=val`), or leaves `dst` untouched if DTO is empty (`None()`).
4. **Adversarial URL Escaping & Roundtrip**: `TestDTO_QueryEscapingSpecialChars` only tests ASCII query punctuation. It does not test Unicode (CJK, Cyrillic, German umlauts, Arabic), Emojis (`🚀🔥`), control chars (`\x00`, `\n`, `\r`, `\t`), unreserved chars (`-._~`), or verified lossless decoding via `net/url.QueryUnescape`.
5. **Multi-Field Empty String Permutations**: No tests verify all permutations of multiple `Some("")`, `None()`, and populated strings across consecutive struct fields (leading, middle, trailing, interleaved empty strings) to ensure delimiter `&` placement is never broken (e.g. no leading `&`, no double `&&`, no trailing `&`).
6. **Slice Collections in DTO**: Neither `SearchFilter` nor `AllPrimitivesDTO` includes `generic.Optional[[]int]` or `generic.Optional[[]string]`. These codegen branches (`lines 276-294` of `dto.go`) have 0% test coverage.
7. **Monad Direct Primitive Boundary Roundtrips**: `monads_test.go` and `monads_adversarial_test.go` test primitives inside structs with `omitzero`, but lack direct JSON marshal/unmarshal assertions for extreme values: `math.MinInt64`, `math.MaxUint64`, `-0.0`, `math.MaxFloat64`, `math.SmallestNonzeroFloat64`, unicode strings, and collection types (`Optional[[]int]`, `Optional[map[string]int]`).

---

## 2. Logic Chain

1. **Premise**: Sovereign benchmark-grade quality requires total test resilience for DTO emission and generic monads under adversarial conditions.
2. **Analysis of DTO Codegen (`pkg/emitter/dto.go`)**:
   - `dto.go` emits `if r == nil { return dst }` (`line 23`) and `if r == nil { return }` (`line 36`).
   - If an API consumer passes a nil DTO pointer, it should cleanly return the input slice or no-op on `url.Values` without panicking. An explicit adversarial test must assert this.
   - `AppendFormData` appends to `dst`. If `cap(dst) < requiredLen`, Go's built-in `append` reallocates the underlying array. If `cap(dst) >= requiredLen`, it performs zero heap allocations. Tests must assert that output is byte-identical across zero capacity, small capacity, exact capacity, and oversized capacity.
   - In `AppendFormData`, delimiter logic checks `if len(dst) > 0 { dst = append(dst, '&') }`. When `dst` has pre-existing content (`len(dst) > 0`), the first serialized field must have a leading `&`. When `dst` is empty, no leading `&` is emitted. When all fields are `None()`, `dst` must remain completely untouched. Tests must assert this invariant.
   - `appendQueryEscape` encodes non-ASCII bytes as `%XX`. Because UTF-8 multi-byte characters consist of bytes `>= 0x80`, every byte in a multi-byte sequence is percent-escaped. Standard `url.QueryUnescape` reverses this transformation. Therefore, decoding the encoded value must yield the identical original Unicode/Emoji/Control string. Tests must assert lossless roundtrip.
   - For consecutive optional string fields, empty strings (`Some("")`) emit `wire=`. Delimiters must join them strictly as `a=&b=&c=` or `a=&c=val` without double ampersands (`&&`) or stray `&`.
3. **Analysis of Monad Monads (`foundation/generic/monads.go`)**:
   - In Go 1.24+, `encoding/json` evaluates `IsZero() bool` when `omitzero` tag is present.
   - For `Some(zeroValue)` (e.g. `Some("")`, `Some(0)`, `Some(false)`, `Some(0.0)`), `IsZero()` is `false`. The field MUST be serialized in JSON.
   - For `None[T]()`, `IsZero()` is `true`. The field MUST be omitted.
   - When unmarshaling, `[]byte("null")`, `[]byte("")`, and whitespace payloads must reset the target optional to `None[T]()`. Conversely, JSON empty string `[]byte("\"\"")` must produce `Some("")`.
4. **Synthesis**:
   - Providing worker_m4 with complete, self-contained test implementations will close all 7 identified gaps immediately and ensure 100% adversarial verification.

---

## 3. Caveats

- **Explorer Read-Only Boundary**: As an Explorer, no production code or test files were modified. All test code provided in Section 4 is designed as drop-in implementations for `worker_m4`.
- **Operating System Environment**: Shell commands are executed under Windows PowerShell. Quoting in generated test files inside temp directories must use standard Go strings rather than PowerShell inline commands to prevent shell metacharacter expansion (such as `&` or `$`).
- **Benchmark Allocations**: When testing `testing.AllocsPerRun`, buffers must be declared outside closures (`var buf [1024]byte`) and passed as `buf[:0]` to avoid closure variable escape allocations.

---

## 4. Conclusion & Test Designs for worker_m4

`worker_m4` should implement two high-craft adversarial test suites:
1. **Suite A**: Add to `d:/CodingProjects/foundation/generic/monads_adversarial_test.go` to close Monad boundary coverage.
2. **Suite B**: Add to `d:/CodingProjects/vortex/pkg/emitter/dto_test.go` to close DTO codegen boundary, buffer, nil receiver, unicode escaping, and collection coverage.

### 4.1 Exact Code for `foundation/generic/monads_adversarial_test.go`

Add the following three test functions to `d:/CodingProjects/foundation/generic/monads_adversarial_test.go`:

```go
func TestOptional_Adversarial_PrimitiveBoundariesRoundtrip(t *testing.T) {
	// 1. Extreme Int64 Boundaries
	int64Cases := []int64{
		math.MinInt64,
		-9223372036854775807,
		-1000000000,
		-1,
		0,
		1,
		1000000000,
		math.MaxInt64,
	}
	for _, val := range int64Cases {
		opt := Some(val)
		b, err := json.Marshal(opt)
		assert.Nil(t, err)
		assert.Equal(t, fmt.Sprintf("%d", val), string(b))

		var decoded Optional[int64]
		err = json.Unmarshal(b, &decoded)
		assert.Nil(t, err)
		assert.True(t, decoded.IsPresent())
		assert.Equal(t, val, decoded.MustValue())
		assert.False(t, decoded.IsZero())
	}

	// 2. Extreme Uint64 Boundaries
	uint64Cases := []uint64{
		0,
		1,
		math.MaxUint32,
		math.MaxUint64,
	}
	for _, val := range uint64Cases {
		opt := Some(val)
		b, err := json.Marshal(opt)
		assert.Nil(t, err)
		assert.Equal(t, fmt.Sprintf("%d", val), string(b))

		var decoded Optional[uint64]
		err = json.Unmarshal(b, &decoded)
		assert.Nil(t, err)
		assert.True(t, decoded.IsPresent())
		assert.Equal(t, val, decoded.MustValue())
		assert.False(t, decoded.IsZero())
	}

	// 3. Float64 Extreme Boundaries
	float64Cases := []float64{
		-12345.6789,
		-1.0,
		-0.0,
		0.0,
		1.0,
		3.141592653589793,
		math.MaxFloat64,
		math.SmallestNonzeroFloat64,
	}
	for _, val := range float64Cases {
		opt := Some(val)
		b, err := json.Marshal(opt)
		assert.Nil(t, err)

		var decoded Optional[float64]
		err = json.Unmarshal(b, &decoded)
		assert.Nil(t, err)
		assert.True(t, decoded.IsPresent())
		assert.Equal(t, val, decoded.MustValue())
		assert.False(t, decoded.IsZero())
	}

	// 4. Unicode, Emojis, and Escape Strings
	strCases := []string{
		"",
		"simple",
		"spaces in between",
		"alpha & beta = gamma ? # 100% / test+case",
		"日本語のクエリ",
		"Привет мир",
		"Größe Überprüfung",
		"مرحبا بالعالم",
		"🚀🔥💻🎉",
		"line1\nline2\r\nline3\ttab",
		"\"quoted\" and \\backslash\\",
	}
	for _, val := range strCases {
		opt := Some(val)
		b, err := json.Marshal(opt)
		assert.Nil(t, err)

		var decoded Optional[string]
		err = json.Unmarshal(b, &decoded)
		assert.Nil(t, err)
		assert.True(t, decoded.IsPresent())
		assert.Equal(t, val, decoded.MustValue())
		assert.False(t, decoded.IsZero())
	}

	// 5. Booleans
	for _, val := range []bool{true, false} {
		opt := Some(val)
		b, err := json.Marshal(opt)
		assert.Nil(t, err)
		assert.Equal(t, fmt.Sprintf("%t", val), string(b))

		var decoded Optional[bool]
		err = json.Unmarshal(b, &decoded)
		assert.Nil(t, err)
		assert.True(t, decoded.IsPresent())
		assert.Equal(t, val, decoded.MustValue())
		assert.False(t, decoded.IsZero())
	}
}

func TestOptional_Adversarial_CollectionsAndPointersRoundtrip(t *testing.T) {
	// 1. Slices
	optSlicePopulated := Some([]int{10, -20, 30})
	b, err := json.Marshal(optSlicePopulated)
	assert.Nil(t, err)
	assert.Equal(t, "[10,-20,30]", string(b))

	var decodedSlice Optional[[]int]
	err = json.Unmarshal(b, &decodedSlice)
	assert.Nil(t, err)
	assert.True(t, decodedSlice.IsPresent())
	assert.Equal(t, []int{10, -20, 30}, decodedSlice.MustValue())
	assert.False(t, decodedSlice.IsZero())

	// Empty slice
	optSliceEmpty := Some([]int{})
	b, err = json.Marshal(optSliceEmpty)
	assert.Nil(t, err)
	assert.Equal(t, "[]", string(b))

	var decodedEmptySlice Optional[[]int]
	err = json.Unmarshal(b, &decodedEmptySlice)
	assert.Nil(t, err)
	assert.True(t, decodedEmptySlice.IsPresent())
	assert.Equal(t, 0, len(decodedEmptySlice.MustValue()))
	assert.False(t, decodedEmptySlice.IsZero())

	// None slice
	optSliceNone := None[[]int]()
	b, err = json.Marshal(optSliceNone)
	assert.Nil(t, err)
	assert.Equal(t, "null", string(b))
	assert.True(t, optSliceNone.IsZero())

	// 2. Maps
	optMap := Some(map[string]int{"a": 1, "b": 2})
	b, err = json.Marshal(optMap)
	assert.Nil(t, err)

	var decodedMap Optional[map[string]int]
	err = json.Unmarshal(b, &decodedMap)
	assert.Nil(t, err)
	assert.True(t, decodedMap.IsPresent())
	assert.Equal(t, 2, len(decodedMap.MustValue()))
	assert.Equal(t, 1, decodedMap.MustValue()["a"])
	assert.Equal(t, 2, decodedMap.MustValue()["b"])
	assert.False(t, decodedMap.IsZero())

	// None map
	optMapNone := None[map[string]int]()
	b, err = json.Marshal(optMapNone)
	assert.Nil(t, err)
	assert.Equal(t, "null", string(b))
	assert.True(t, optMapNone.IsZero())

	// 3. Pointers
	x := 42
	optPtr := Some(&x)
	b, err = json.Marshal(optPtr)
	assert.Nil(t, err)
	assert.Equal(t, "42", string(b))
	assert.False(t, optPtr.IsZero())

	var decodedPtr Optional[*int]
	err = json.Unmarshal(b, &decodedPtr)
	assert.Nil(t, err)
	assert.True(t, decodedPtr.IsPresent())
	assert.Equal(t, 42, *decodedPtr.MustValue())

	// Nil pointer inside Some
	optNilPtr := Some[*int](nil)
	b, err = json.Marshal(optNilPtr)
	assert.Nil(t, err)
	assert.Equal(t, "null", string(b))
	assert.True(t, optNilPtr.IsPresent())
	assert.False(t, optNilPtr.IsZero())
}

func TestOptional_Adversarial_DirectIsZeroMatrix(t *testing.T) {
	// Explicit Some(zeroValue) MUST report IsZero() == false
	assert.False(t, Some("").IsZero())
	assert.False(t, Some(0).IsZero())
	assert.False(t, Some(int64(0)).IsZero())
	assert.False(t, Some(uint(0)).IsZero())
	assert.False(t, Some(uint64(0)).IsZero())
	assert.False(t, Some(0.0).IsZero())
	assert.False(t, Some(false).IsZero())
	assert.False(t, Some(time.Time{}).IsZero())
	assert.False(t, Some([]int{}).IsZero())
	assert.False(t, Some(map[string]int{}).IsZero())
	assert.False(t, Some[*int](nil).IsZero())

	// Unset None[T]() MUST report IsZero() == true
	assert.True(t, None[string]().IsZero())
	assert.True(t, None[int]().IsZero())
	assert.True(t, None[int64]().IsZero())
	assert.True(t, None[uint]().IsZero())
	assert.True(t, None[uint64]().IsZero())
	assert.True(t, None[float64]().IsZero())
	assert.True(t, None[bool]().IsZero())
	assert.True(t, None[time.Time]().IsZero())
	assert.True(t, None[[]int]().IsZero())
	assert.True(t, None[map[string]int]().IsZero())
	assert.True(t, None[*int]().IsZero())
}
```

### 4.2 Exact Code for `pkg/emitter/dto_test.go`

Add the following comprehensive adversarial test to `d:/CodingProjects/vortex/pkg/emitter/dto_test.go`:

```go
func TestEmitter_DTO_Adversarial_FullSuite(t *testing.T) {
	src := `package advdto

import (
	"time"
	"github.com/lemon4ksan/foundation/generic"
)

// @aoni:dto
type FullAdversarialDTO struct {
	StrA        generic.Optional[string]    ` + "`query:\"a\"`" + `
	StrB        generic.Optional[string]    ` + "`query:\"b\"`" + `
	StrC        generic.Optional[string]    ` + "`query:\"c\"`" + `
	QueryEsc    generic.Optional[string]    ` + "`query:\"q\"`" + `
	IntVal      generic.Optional[int64]     ` + "`query:\"num\"`" + `
	UintVal     generic.Optional[uint64]    ` + "`query:\"unum\"`" + `
	FloatVal    generic.Optional[float64]   ` + "`query:\"flt\"`" + `
	BoolStd     generic.Optional[bool]      ` + "`query:\"b_std\"`" + `
	// @format bool_int
	BoolInt     generic.Optional[bool]      ` + "`query:\"b_int\"`" + `
	// @format flag
	BoolFlag    generic.Optional[bool]      ` + "`query:\"b_flag\"`" + `
	SliceInt    generic.Optional[[]int]     ` + "`query:\"items\"`" + `
	SliceStr    generic.Optional[[]string]  ` + "`query:\"tags\"`" + `
	TimeVal     generic.Optional[time.Time] ` + "`query:\"ts\"`" + `
}
`

	tmpDir := t.TempDir()
	foundationPath := filepath.ToSlash("d:/CodingProjects/foundation")
	goModContent := "module advdto\n\ngo 1.27.0\n\nrequire github.com/lemon4ksan/foundation v0.0.0\n\nreplace github.com/lemon4ksan/foundation => " + foundationPath + "\n"

	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0o600)
	require.NoError(t, err)

	dtoSrcPath := filepath.Join(tmpDir, "adv.go")
	err = os.WriteFile(dtoSrcPath, []byte(src), 0o600)
	require.NoError(t, err)

	p := parser.NewParser()
	root, err := p.ParseFile(dtoSrcPath)
	require.NoError(t, err)

	genBytes, err := emitter.Emit(root)
	require.NoError(t, err)

	genFile := filepath.Join(tmpDir, "adv.gen.go")
	err = os.WriteFile(genFile, genBytes, 0o600)
	require.NoError(t, err)

	testSrc := `package advdto

import (
	"math"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/lemon4ksan/foundation/generic"
)

// 1. Nil Receiver Safety
func TestAdversarialDTO_NilReceiverSafety(t *testing.T) {
	var nilDTO *FullAdversarialDTO
	var buf [256]byte

	outForm := nilDTO.AppendFormData(buf[:0])
	if len(outForm) != 0 {
		t.Fatalf("expected empty slice for nil receiver AppendFormData, got %q", string(outForm))
	}

	outQuery := nilDTO.AppendQuery(buf[:0])
	if len(outQuery) != 0 {
		t.Fatalf("expected empty slice for nil receiver AppendQuery, got %q", string(outQuery))
	}

	vals := make(url.Values)
	nilDTO.EncodeValues(vals)
	if len(vals) != 0 {
		t.Fatalf("expected untouched url.Values for nil receiver EncodeValues, got %v", vals)
	}
}

// 2. Empty String Combinations and Permutations
func TestAdversarialDTO_EmptyStringPermutations(t *testing.T) {
	var buf [256]byte

	// Case 1: All 3 fields are Some("") -> "a=&b=&c="
	dtoAllEmpty := &FullAdversarialDTO{
		StrA: generic.Some(""),
		StrB: generic.Some(""),
		StrC: generic.Some(""),
	}
	out := string(dtoAllEmpty.AppendFormData(buf[:0]))
	if out != "a=&b=&c=" {
		t.Fatalf("expected 'a=&b=&c=', got %q", out)
	}

	// Case 2: Middle field only Some("") -> "b=" (no leading or trailing &)
	dtoMiddleOnly := &FullAdversarialDTO{
		StrB: generic.Some(""),
	}
	out = string(dtoMiddleOnly.AppendFormData(buf[:0]))
	if out != "b=" {
		t.Fatalf("expected 'b=', got %q", out)
	}

	// Case 3: Leading empty, trailing populated -> "a=&c=end"
	dtoLeadEmpty := &FullAdversarialDTO{
		StrA: generic.Some(""),
		StrC: generic.Some("end"),
	}
	out = string(dtoLeadEmpty.AppendFormData(buf[:0]))
	if out != "a=&c=end" {
		t.Fatalf("expected 'a=&c=end', got %q", out)
	}

	// Case 4: Interleaved -> "a=&b=mid&c="
	dtoInterleaved := &FullAdversarialDTO{
		StrA: generic.Some(""),
		StrB: generic.Some("mid"),
		StrC: generic.Some(""),
	}
	out = string(dtoInterleaved.AppendFormData(buf[:0]))
	if out != "a=&b=mid&c=" {
		t.Fatalf("expected 'a=&b=mid&c=', got %q", out)
	}

	// Case 5: All None -> 0 bytes
	dtoAllNone := &FullAdversarialDTO{}
	outBytes := dtoAllNone.AppendFormData(buf[:0])
	if len(outBytes) != 0 {
		t.Fatalf("expected 0 bytes for all None, got %q", string(outBytes))
	}
}

// 3. Adversarial Unicode, Emojis, and Lossless Query Unescaping
func TestAdversarialDTO_UnicodeEmojisAndQueryUnescape(t *testing.T) {
	testCases := []string{
		"hello world",
		"alpha & beta = gamma ? # 100% / test+case",
		"日本語の検索クエリ",
		"Привет мир",
		"Größe Überprüfung",
		"مرحبا بالعالم",
		"🔥 VORTEX 🚀 benchmark ⚡",
		"\x00\n\r\t",
		"-._~",
		"@:;$,\"'<>\\^|{}",
	}

	var buf [1024]byte
	for _, tc := range testCases {
		dto := &FullAdversarialDTO{
			QueryEsc: generic.Some(tc),
		}
		raw := string(dto.AppendFormData(buf[:0]))
		if !strings.HasPrefix(raw, "q=") {
			t.Fatalf("expected prefix 'q=', got: %s", raw)
		}

		encodedVal := strings.TrimPrefix(raw, "q=")
		decodedVal, err := url.QueryUnescape(encodedVal)
		if err != nil {
			t.Fatalf("QueryUnescape failed for encoded value %q: %v", encodedVal, err)
		}
		if decodedVal != tc {
			t.Fatalf("roundtrip mismatch for %q:\nEncoded: %s\nDecoded: %s", tc, encodedVal, decodedVal)
		}

		// Verify zero allocations when buffer capacity is sufficient
		allocs := testing.AllocsPerRun(1000, func() {
			_ = dto.AppendFormData(buf[:0])
		})
		if allocs != 0 {
			t.Fatalf("expected 0 allocs for %q, got %v", tc, allocs)
		}
	}
}

// 4. Extreme Primitive Boundaries & Explicit Zeros
func TestAdversarialDTO_ExtremeBoundariesAndZeros(t *testing.T) {
	dto := &FullAdversarialDTO{
		IntVal:   generic.Some(int64(math.MinInt64)),
		UintVal:  generic.Some(uint64(math.MaxUint64)),
		FloatVal: generic.Some(0.0),
		BoolStd:  generic.Some(false),
		BoolInt:  generic.Some(false),
		BoolFlag: generic.Some(false),
	}

	var buf [512]byte
	out := string(dto.AppendFormData(buf[:0]))
	expected := "num=-9223372036854775808&unum=18446744073709551615&flt=0&b_std=false&b_int=0"
	if out != expected {
		t.Fatalf("mismatch for extreme boundaries:\nExpected: %s\nGot:      %s", expected, out)
	}

	// Verify Flag when true emits bare key
	dtoFlag := &FullAdversarialDTO{
		BoolFlag: generic.Some(true),
	}
	outFlag := string(dtoFlag.AppendFormData(buf[:0]))
	if outFlag != "b_flag" {
		t.Fatalf("expected 'b_flag', got %q", outFlag)
	}

	// Verify EncodeValues matches
	vals := make(url.Values)
	dto.EncodeValues(vals)
	if vals.Get("num") != "-9223372036854775808" {
		t.Fatalf("EncodeValues num mismatch: %q", vals.Get("num"))
	}
	if vals.Get("unum") != "18446744073709551615" {
		t.Fatalf("EncodeValues unum mismatch: %q", vals.Get("unum"))
	}
	if vals.Get("b_int") != "0" {
		t.Fatalf("EncodeValues b_int mismatch: %q", vals.Get("b_int"))
	}
}

// 5. Slice Collections in Optional
func TestAdversarialDTO_SliceCollections(t *testing.T) {
	dto := &FullAdversarialDTO{
		SliceInt: generic.Some([]int{10, -20, 30}),
		SliceStr: generic.Some([]string{"alpha", "beta & gamma"}),
	}

	var buf [512]byte
	out := string(dto.AppendFormData(buf[:0]))
	expected := "items=10&items=-20&items=30&tags=alpha&tags=beta+%26+gamma"
	if out != expected {
		t.Fatalf("mismatch for slice collections:\nExpected: %s\nGot:      %s", expected, out)
	}

	vals := make(url.Values)
	dto.EncodeValues(vals)
	if len(vals["items"]) != 3 || vals["items"][0] != "10" || vals["items"][1] != "-20" || vals["items"][2] != "30" {
		t.Fatalf("unexpected vals[items]: %v", vals["items"])
	}
	if len(vals["tags"]) != 2 || vals["tags"][0] != "alpha" || vals["tags"][1] != "beta & gamma" {
		t.Fatalf("unexpected vals[tags]: %v", vals["tags"])
	}

	// Empty slices must omit
	dtoEmptySlices := &FullAdversarialDTO{
		SliceInt: generic.Some([]int{}),
		SliceStr: generic.Some([]string{}),
	}
	outEmpty := dtoEmptySlices.AppendFormData(buf[:0])
	if len(outEmpty) != 0 {
		t.Fatalf("expected 0 bytes for empty slices, got: %q", string(outEmpty))
	}
}

// 6. Buffer Capacity Boundaries & Growth
func TestAdversarialDTO_BufferCapacitiesAndGrowth(t *testing.T) {
	dto := &FullAdversarialDTO{
		StrA:    generic.Some("hello"),
		IntVal:  generic.Some(int64(42)),
		BoolStd: generic.Some(true),
	}
	expected := "a=hello&num=42&b_std=true"

	// Zero capacity: make([]byte, 0, 0)
	outZero := string(dto.AppendFormData(make([]byte, 0, 0)))
	if outZero != expected {
		t.Fatalf("zero-cap failed: expected %q, got %q", expected, outZero)
	}

	// Tiny capacity: make([]byte, 0, 4)
	outTiny := string(dto.AppendFormData(make([]byte, 0, 4)))
	if outTiny != expected {
		t.Fatalf("tiny-cap failed: expected %q, got %q", expected, outTiny)
	}

	// Pre-populated buffer with existing content
	existing := []byte("prefix=ok")
	outPrefixed := string(dto.AppendFormData(existing))
	if outPrefixed != "prefix=ok&"+expected {
		t.Fatalf("prefixed buffer mismatch:\nExpected: prefix=ok&%s\nGot:      %s", expected, outPrefixed)
	}

	// Pre-populated buffer with all None DTO -> remains untouched
	dtoNone := &FullAdversarialDTO{}
	cleanExisting := []byte("prefix=ok")
	outNonePrefixed := string(dtoNone.AppendFormData(cleanExisting))
	if outNonePrefixed != "prefix=ok" {
		t.Fatalf("empty DTO should leave prefix untouched, got %q", outNonePrefixed)
	}
}
`

	testFile := filepath.Join(tmpDir, "adv_test.go")
	err = os.WriteFile(testFile, []byte(testSrc), 0o600)
	require.NoError(t, err)

	cmd := exec.Command("go", "test", "-v", ".")
	cmd.Dir = tmpDir
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go test failed in temporary directory %s: %v\nOutput:\n%s", tmpDir, err, string(outBytes))
	}

	t.Logf("Generated Adversarial DTO test output:\n%s", string(outBytes))
	require.Contains(t, string(outBytes), "PASS")
}
```

---

## 5. Verification Method

To independently verify the adversarial test suite:

1. **Verify Foundation Generic Tests**:
   ```pwsh
   cd d:\CodingProjects\foundation
   go test -v -race ./generic/...
   ```
   **Pass Condition**: All tests pass cleanly, including the new primitive boundaries, collections/pointers roundtrip, and direct `IsZero` matrix. Zero race conditions detected.

2. **Verify Vortex Emitter Tests**:
   ```pwsh
   cd d:\CodingProjects\vortex
   go test -v ./pkg/emitter/...
   ```
   **Pass Condition**: `TestEmitter_DTO_Adversarial_FullSuite` passes cleanly, demonstrating nil receiver safety, buffer growth/reuse, Unicode/Emoji URL unescaping, and multi-field empty string permutations.

3. **Verify DTO Zero-Allocation Regressions**:
   ```pwsh
   cd d:\CodingProjects\vortex
   go test -v -bench=. -benchmem ./pkg/emitter/...
   ```
   **Pass Condition**: All benchmark functions (`Benchmark_AppendFormData_*`, `Benchmark_AppendQuery_*`) report `0 B/op` and `0 allocs/op`.

4. **Full Workspace Gate**:
   ```pwsh
   cd d:\CodingProjects\vortex
   go test ./...
   ```
   **Pass Condition**: All packages pass with zero build errors and zero test failures.
