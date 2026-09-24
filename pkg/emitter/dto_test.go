// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package emitter_test

import (
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/lemon4ksan/foundation/generic"
	"github.com/lemon4ksan/foundation/testing/require"

	"github.com/lemon4ksan/vortex/pkg/emitter"
	"github.com/lemon4ksan/vortex/pkg/parser"
)

func TestZeroAlloc_AppendQuery_Primitives(t *testing.T) {
	dto := newFullBenchmarkDTO()
	var buf [1024]byte
	allocs := testing.AllocsPerRun(1000, func() {
		_ = dto.AppendQuery(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for AppendQuery, got %v", allocs)
	}
}

func TestZeroAlloc_AppendFormData_Primitives(t *testing.T) {
	dto := newFullBenchmarkDTO()
	var buf [1024]byte
	allocs := testing.AllocsPerRun(1000, func() {
		_ = dto.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for AppendFormData, got %v", allocs)
	}
}

func TestZeroAlloc_SomeEmptyString(t *testing.T) {
	dto := &BenchmarkTestDTO{EmptyStr: generic.Some("")}
	var buf [64]byte
	out := string(dto.AppendFormData(buf[:0]))
	if out != "empty_str=" {
		t.Fatalf("expected 'empty_str=', got %q", out)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = dto.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for Some(\"\"), got %v", allocs)
	}
}

func TestZeroAlloc_AllNone(t *testing.T) {
	dto := &BenchmarkTestDTO{}
	var buf [64]byte
	out := dto.AppendFormData(buf[:0])
	if len(out) != 0 {
		t.Fatalf("expected 0 bytes for AllNone, got %d", len(out))
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = dto.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for AllNone, got %v", allocs)
	}
}

func TestZeroAlloc_BoundaryValues(t *testing.T) {
	dto := &BenchmarkTestDTO{
		Int64Val:    generic.Some(int64(-9223372036854775808)),
		Uint64Val:   generic.Some(uint64(18446744073709551615)),
		Float64Val:  generic.Some(0.0),
		BoolVal:     generic.Some(false),
		BoolIntVal:  generic.Some(false),
		BoolFlagVal: generic.Some(false),
	}
	var buf [512]byte
	allocs := testing.AllocsPerRun(1000, func() {
		_ = dto.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for boundary values, got %v", allocs)
	}
}

func TestZeroAlloc_QueryEscaping(t *testing.T) {
	dto := &BenchmarkTestDTO{
		StrVal: generic.Some("alpha & beta = gamma ? # 100% / test+case"),
	}
	var buf [512]byte
	allocs := testing.AllocsPerRun(1000, func() {
		_ = dto.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for query escaping, got %v", allocs)
	}
}

func TestZeroAlloc_EncodeValues_None(t *testing.T) {
	dto := &BenchmarkTestDTO{}
	vals := make(url.Values)
	allocs := testing.AllocsPerRun(1000, func() {
		dto.EncodeValues(vals)
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for EncodeValues on None, got %v", allocs)
	}
}

func TestZeroAlloc_NilReceiver(t *testing.T) {
	var nilDTO *BenchmarkTestDTO
	var buf [64]byte
	allocs := testing.AllocsPerRun(1000, func() {
		_ = nilDTO.AppendFormData(buf[:0])
		_ = nilDTO.AppendQuery(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for nil receiver, got %v", allocs)
	}
}

func TestEmitter_DTO_Emission(t *testing.T) {
	src := `package filterapi

import (
	"time"
	"github.com/lemon4ksan/foundation/generic"
)

// @aoni:dto
type SearchFilter struct {
	Query    generic.Optional[string]    ` + "`query:\"q\"`" + `
	Page     generic.Optional[int]       ` + "`query:\"page\"`" + `
	Count    generic.Optional[uint64]    ` + "`query:\"count\"`" + `
	Score    generic.Optional[float64]   ` + "`query:\"score\"`" + `
	Active   generic.Optional[bool]      ` + "`query:\"active\"`" + `
	Created  generic.Optional[time.Time] ` + "`query:\"created\"`" + `
	PlainStr string                      ` + "`query:\"plain\"`" + `
}
`

	p := parser.NewParser()
	root, err := p.ParseSource("filter.go", []byte(src))
	require.NoError(t, err)
	require.NotNil(t, root)
	require.Len(t, root.Structs, 1)

	codeBytes, err := emitter.Emit(root)
	require.NoError(t, err)
	require.NotEmpty(t, codeBytes)

	codeStr := string(codeBytes)

	// Check methods emitted
	require.Contains(t, codeStr, "func (r *SearchFilter) AppendFormData(dst []byte) []byte")
	require.Contains(t, codeStr, "func (r *SearchFilter) AppendQuery(dst []byte) []byte")
	require.Contains(t, codeStr, "func (r *SearchFilter) EncodeValues(vals url.Values)")

	// Check zero-allocation helper emitted
	require.Contains(t, codeStr, "func appendQueryEscape(dst []byte, s string) []byte")

	// Check zero-allocation primitives in AppendFormData
	require.Contains(t, codeStr, `dst = append(dst, "q="...)`)
	require.Contains(t, codeStr, "dst = appendQueryEscape(dst, optVal)")
	require.Contains(t, codeStr, "strconv.AppendInt(dst, int64(optVal), 10)")
	require.Contains(t, codeStr, "strconv.AppendUint(dst, uint64(optVal), 10)")
	require.Contains(t, codeStr, "strconv.AppendFloat(dst, float64(optVal), 'f', -1, 64)")
	require.Contains(t, codeStr, `dst = append(dst, "active=true"...)`)
	require.Contains(t, codeStr, `dst = append(dst, "active=false"...)`)
	require.Contains(t, codeStr, "optVal.AppendFormat(timeBuf[:0], time.RFC3339)")

	// Check EncodeValues primitive formatting
	require.Contains(t, codeStr, `vals.Set("q", optVal)`)
	require.Contains(t, codeStr, `vals.Set("page", strconv.FormatInt(int64(optVal), 10))`)
	require.Contains(t, codeStr, `vals.Set("count", strconv.FormatUint(uint64(optVal), 10))`)
	require.Contains(t, codeStr, `vals.Set("score", strconv.FormatFloat(float64(optVal), 'f', -1, 64))`)
	require.Contains(t, codeStr, `vals.Set("active", "true")`)
	require.Contains(t, codeStr, `vals.Set("created", optVal.Format(time.RFC3339))`)

	// Ensure fmt.Sprint is NOT called for any of these primitives
	require.NotContains(t, codeStr, "fmt.Sprint(optVal)")
}

func TestEmitter_DTO_ExecutionAndZeroAlloc(t *testing.T) {
	src := `package filterapi

import (
	"time"
	"github.com/lemon4ksan/foundation/generic"
)

// @aoni:dto
type SearchFilter struct {
	Query   generic.Optional[string]    ` + "`query:\"q\"`" + `
	Page    generic.Optional[int]       ` + "`query:\"page\"`" + `
	Count   generic.Optional[uint64]    ` + "`query:\"count\"`" + `
	Score   generic.Optional[float64]   ` + "`query:\"score\"`" + `
	Active  generic.Optional[bool]      ` + "`query:\"active\"`" + `
	Created generic.Optional[time.Time] ` + "`query:\"created\"`" + `
}
`

	tmpDir := t.TempDir()
	foundationPath := filepath.ToSlash("d:/CodingProjects/foundation")
	goModContent := "module filterapi\n\ngo 1.27.0\n\nrequire github.com/lemon4ksan/foundation v0.0.0\n\nreplace github.com/lemon4ksan/foundation => " + foundationPath + "\n"

	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0o600)
	require.NoError(t, err)

	apiFile := filepath.Join(tmpDir, "filter.go")
	err = os.WriteFile(apiFile, []byte(src), 0o600)
	require.NoError(t, err)

	p := parser.NewParser()
	root, err := p.ParseFile(apiFile)
	require.NoError(t, err)

	genBytes, err := emitter.Emit(root)
	require.NoError(t, err)

	genFile := filepath.Join(tmpDir, "filter.gen.go")
	err = os.WriteFile(genFile, genBytes, 0o600)
	require.NoError(t, err)

	testSrc := `package filterapi

import (
	"net/url"
	"testing"
	"time"

	"github.com/lemon4ksan/foundation/generic"
)

func TestSearchFilter_NoneOmissionAndZeroAlloc(t *testing.T) {
	filter := &SearchFilter{}
	var buf [256]byte
	out := filter.AppendFormData(buf[:0])
	if len(out) != 0 {
		t.Fatalf("expected empty output for None, got %q", string(out))
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = filter.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for None, got %v", allocs)
	}
}

func TestSearchFilter_SomeEmptyString(t *testing.T) {
	filter := &SearchFilter{
		Query: generic.Some(""),
	}
	var buf [256]byte
	out := string(filter.AppendFormData(buf[:0]))
	if out != "q=" {
		t.Fatalf("expected 'q=', got %q", out)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = filter.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for Some(\"\"), got %v", allocs)
	}

	// Also verify EncodeValues
	vals := make(url.Values)
	filter.EncodeValues(vals)
	if vals.Encode() != "q=" {
		t.Fatalf("expected 'q=', got %q", vals.Encode())
	}
}

func TestSearchFilter_PrimitivesRoundtripAndZeroAlloc(t *testing.T) {
	tm := time.Date(2026, 9, 22, 12, 0, 0, 0, time.UTC)
	filter := &SearchFilter{
		Query:   generic.Some("hello world"),
		Page:    generic.Some(42),
		Count:   generic.Some(uint64(100)),
		Score:   generic.Some(3.14),
		Active:  generic.Some(true),
		Created: generic.Some(tm),
	}

	var buf [512]byte
	out := string(filter.AppendFormData(buf[:0]))
	expected := "q=hello+world&page=42&count=100&score=3.14&active=true&created=2026-09-22T12%3A00%3A00Z"
	if out != expected {
		t.Fatalf("mismatch: expected %q, got %q", expected, out)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = filter.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for primitives, got %v", allocs)
	}

	// Verify AppendQuery delegates to AppendFormData
	queryOut := string(filter.AppendQuery(buf[:0]))
	if queryOut != expected {
		t.Fatalf("AppendQuery mismatch: expected %q, got %q", expected, queryOut)
	}

	// Verify EncodeValues
	vals := make(url.Values)
	filter.EncodeValues(vals)
	if vals.Get("q") != "hello world" {
		t.Fatalf("expected q='hello world', got %q", vals.Get("q"))
	}
	if vals.Get("page") != "42" {
		t.Fatalf("expected page='42', got %q", vals.Get("page"))
	}
	if vals.Get("count") != "100" {
		t.Fatalf("expected count='100', got %q", vals.Get("count"))
	}
	if vals.Get("score") != "3.14" {
		t.Fatalf("expected score='3.14', got %q", vals.Get("score"))
	}
	if vals.Get("active") != "true" {
		t.Fatalf("expected active='true', got %q", vals.Get("active"))
	}
	if vals.Get("created") != "2026-09-22T12:00:00Z" {
		t.Fatalf("expected created='2026-09-22T12:00:00Z', got %q", vals.Get("created"))
	}
}

func TestSearchFilter_ExplicitFalseAndZero(t *testing.T) {
	filter := &SearchFilter{
		Page:   generic.Some(0),
		Count:  generic.Some(uint64(0)),
		Score:  generic.Some(0.0),
		Active: generic.Some(false),
	}

	var buf [256]byte
	out := string(filter.AppendFormData(buf[:0]))
	expected := "page=0&count=0&score=0&active=false"
	if out != expected {
		t.Fatalf("mismatch for explicit zeros: expected %q, got %q", expected, out)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = filter.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for explicit zeros, got %v", allocs)
	}
}
`

	testFile := filepath.Join(tmpDir, "filter_test.go")
	err = os.WriteFile(testFile, []byte(testSrc), 0o600)
	require.NoError(t, err)

	cmd := exec.Command("go", "test", "-v", ".")
	cmd.Dir = tmpDir
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go test failed in temporary directory %s: %v\nOutput:\n%s", tmpDir, err, string(outBytes))
	}

	t.Logf("Generated DTO test output:\n%s", string(outBytes))
	require.Contains(t, string(outBytes), "PASS")
}

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
