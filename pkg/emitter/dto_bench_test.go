// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package emitter_test

import (
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	"github.com/lemon4ksan/foundation/generic"
	"github.com/lemon4ksan/foundation/testing/require"

	"github.com/lemon4ksan/vortex/pkg/emitter"
	"github.com/lemon4ksan/vortex/pkg/parser"
)

func BenchmarkAppendQuery_Primitives(b *testing.B) {
	dto := newFullBenchmarkDTO()
	var buf [1024]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendQuery(buf[:0])
	}
}

func BenchmarkAppendQuery_Optionals_Some(b *testing.B) {
	dto := &BenchmarkTestDTO{
		StrVal:   generic.Some("hello world"),
		EmptyStr: generic.Some(""),
		Int64Val: generic.Some(int64(42)),
		BoolVal:  generic.Some(true),
	}
	var buf [256]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendQuery(buf[:0])
	}
}

func BenchmarkAppendQuery_Optionals_None(b *testing.B) {
	dto := &BenchmarkTestDTO{}
	var buf [64]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendQuery(buf[:0])
	}
}

func BenchmarkAppendFormData_Primitives(b *testing.B) {
	dto := newFullBenchmarkDTO()
	var buf [1024]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendFormData(buf[:0])
	}
}

func BenchmarkAppendFormData_Optionals(b *testing.B) {
	dto := &BenchmarkTestDTO{
		StrVal:     generic.Some("vortex high-craft API serialization"),
		IntVal:     generic.Some(100),
		Float64Val: generic.Some(3.14159),
		BoolVal:    generic.Some(true),
	}
	var buf [256]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendFormData(buf[:0])
	}
}

func BenchmarkEncodeValues_ZeroAlloc(b *testing.B) {
	dto := &BenchmarkTestDTO{}
	vals := make(url.Values)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dto.EncodeValues(vals)
	}
}

func TestEmitter_DTO_AllPrimitives_Comprehensive(t *testing.T) {
	src := `package benchdto

import (
	"time"
	"github.com/lemon4ksan/foundation/generic"
)

// @aoni:dto
type AllPrimitivesDTO struct {
	StrVal      generic.Optional[string]    ` + "`query:\"str\"`" + `
	EmptyStr    generic.Optional[string]    ` + "`query:\"empty_str\"`" + `
	NoneStr     generic.Optional[string]    ` + "`query:\"none_str\"`" + `
	IntVal      generic.Optional[int]       ` + "`query:\"i\"`" + `
	Int8Val     generic.Optional[int8]      ` + "`query:\"i8\"`" + `
	Int16Val    generic.Optional[int16]     ` + "`query:\"i16\"`" + `
	Int32Val    generic.Optional[int32]     ` + "`query:\"i32\"`" + `
	Int64Val    generic.Optional[int64]     ` + "`query:\"i64\"`" + `
	UintVal     generic.Optional[uint]      ` + "`query:\"u\"`" + `
	Uint8Val    generic.Optional[uint8]     ` + "`query:\"u8\"`" + `
	Uint16Val   generic.Optional[uint16]    ` + "`query:\"u16\"`" + `
	Uint32Val   generic.Optional[uint32]    ` + "`query:\"u32\"`" + `
	Uint64Val   generic.Optional[uint64]    ` + "`query:\"u64\"`" + `
	UintptrVal  generic.Optional[uintptr]   ` + "`query:\"uptr\"`" + `
	ByteVal     generic.Optional[byte]      ` + "`query:\"bval\"`" + `
	Float32Val  generic.Optional[float32]   ` + "`query:\"f32\"`" + `
	Float64Val  generic.Optional[float64]   ` + "`query:\"f64\"`" + `
	BoolVal     generic.Optional[bool]      ` + "`query:\"b\"`" + `
	// @format bool_int
	BoolIntVal  generic.Optional[bool]      ` + "`query:\"b_int\"`" + `
	// @format flag
	BoolFlagVal generic.Optional[bool]      ` + "`query:\"b_flag\"`" + `
	TimeVal     generic.Optional[time.Time] ` + "`query:\"ts\"`" + `
}
`

	tmpDir := t.TempDir()
	foundationPath := filepath.ToSlash("d:/CodingProjects/foundation")
	goModContent := "module benchdto\n\ngo 1.27.0\n\nrequire github.com/lemon4ksan/foundation v0.0.0\n\nreplace github.com/lemon4ksan/foundation => " + foundationPath + "\n"

	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0o600)
	require.NoError(t, err)

	dtoSrcPath := filepath.Join(tmpDir, "dto.go")
	err = os.WriteFile(dtoSrcPath, []byte(src), 0o600)
	require.NoError(t, err)

	p := parser.NewParser()
	root, err := p.ParseFile(dtoSrcPath)
	require.NoError(t, err)

	genBytes, err := emitter.Emit(root)
	require.NoError(t, err)

	genFile := filepath.Join(tmpDir, "dto.gen.go")
	err = os.WriteFile(genFile, genBytes, 0o600)
	require.NoError(t, err)

	testAndBenchSrc := `package benchdto

import (
	"math"
	"net/url"
	"testing"
	"time"

	"github.com/lemon4ksan/foundation/generic"
)

var fixedTime = time.Date(2026, 9, 22, 14, 30, 0, 0, time.UTC)

func makeFullDTO() *AllPrimitivesDTO {
	return &AllPrimitivesDTO{
		StrVal:      generic.Some("hello world"),
		EmptyStr:    generic.Some(""),
		NoneStr:     generic.None[string](),
		IntVal:      generic.Some(42),
		Int8Val:     generic.Some(int8(-8)),
		Int16Val:    generic.Some(int16(-1600)),
		Int32Val:    generic.Some(int32(-320000)),
		Int64Val:    generic.Some(int64(-64000000000)),
		UintVal:     generic.Some(uint(100)),
		Uint8Val:    generic.Some(uint8(255)),
		Uint16Val:   generic.Some(uint16(65535)),
		Uint32Val:   generic.Some(uint32(4294967295)),
		Uint64Val:   generic.Some(uint64(18446744073709551615)),
		UintptrVal:  generic.Some(uintptr(123456)),
		ByteVal:     generic.Some(byte(65)),
		Float32Val:  generic.Some(float32(2.718)),
		Float64Val:  generic.Some(3.1415926535),
		BoolVal:     generic.Some(true),
		BoolIntVal:  generic.Some(true),
		BoolFlagVal: generic.Some(true),
		TimeVal:     generic.Some(fixedTime),
	}
}

func TestDTO_AllPrimitives_ExactWireFormat(t *testing.T) {
	dto := makeFullDTO()
	var buf [1024]byte
	out := string(dto.AppendFormData(buf[:0]))

	// Check that empty string serialized as empty_str=
	if !contains(out, "empty_str=") {
		t.Fatalf("expected 'empty_str=' in output, got: %s", out)
	}

	// Check that none_str is completely absent
	if contains(out, "none_str") {
		t.Fatalf("expected 'none_str' to be omitted, got: %s", out)
	}

	// Check individual primitive representations
	checks := []string{
		"str=hello+world",
		"empty_str=",
		"i=42",
		"i8=-8",
		"i16=-1600",
		"i32=-320000",
		"i64=-64000000000",
		"u=100",
		"u8=255",
		"u16=65535",
		"u32=4294967295",
		"u64=18446744073709551615",
		"uptr=123456",
		"bval=65",
		"b=true",
		"b_int=1",
		"b_flag",
		"ts=2026-09-22T14%3A30%3A00Z",
	}

	for _, c := range checks {
		if !contains(out, c) {
			t.Errorf("missing expected serialized chunk %q in output: %s", c, out)
		}
	}

	// AppendQuery should equal AppendFormData
	queryOut := string(dto.AppendQuery(buf[:0]))
	if queryOut != out {
		t.Fatalf("AppendQuery output does not match AppendFormData:\nQuery: %s\nForm:  %s", queryOut, out)
	}

	// Verify EncodeValues
	vals := make(url.Values)
	dto.EncodeValues(vals)
	if vals.Get("empty_str") != "" || !vals.Has("empty_str") {
		t.Fatalf("expected empty_str to be present as empty string in url.Values, got: %q", vals.Get("empty_str"))
	}
	if vals.Has("none_str") {
		t.Fatalf("expected none_str to NOT be present in url.Values")
	}
	if vals.Get("str") != "hello world" {
		t.Fatalf("expected str='hello world', got %q", vals.Get("str"))
	}
	if vals.Get("i64") != "-64000000000" {
		t.Fatalf("expected i64='-64000000000', got %q", vals.Get("i64"))
	}
	if vals.Get("u64") != "18446744073709551615" {
		t.Fatalf("expected u64='18446744073709551615', got %q", vals.Get("u64"))
	}
	if vals.Get("b_int") != "1" {
		t.Fatalf("expected b_int='1', got %q", vals.Get("b_int"))
	}
}

func TestDTO_AllPrimitives_ZeroAllocations(t *testing.T) {
	dto := makeFullDTO()
	var buf [1024]byte

	// Verify AppendFormData 0 allocs
	allocs := testing.AllocsPerRun(1000, func() {
		_ = dto.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for AppendFormData, got %v", allocs)
	}

	// Verify AppendQuery 0 allocs
	allocs = testing.AllocsPerRun(1000, func() {
		_ = dto.AppendQuery(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for AppendQuery, got %v", allocs)
	}
}

func TestDTO_SomeEmptyString_Only(t *testing.T) {
	dto := &AllPrimitivesDTO{
		EmptyStr: generic.Some(""),
	}
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

func TestDTO_AllNone_ZeroBytes(t *testing.T) {
	dto := &AllPrimitivesDTO{}
	var buf [64]byte
	out := dto.AppendFormData(buf[:0])
	if len(out) != 0 {
		t.Fatalf("expected 0 bytes for AllNone, got %d bytes: %q", len(out), string(out))
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = dto.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for AllNone, got %v", allocs)
	}
}

func TestDTO_BoundaryValuesAndZeroes(t *testing.T) {
	dto := &AllPrimitivesDTO{
		Int64Val:    generic.Some(int64(math.MinInt64)),
		Uint64Val:   generic.Some(uint64(math.MaxUint64)),
		IntVal:      generic.Some(0),
		Float64Val:  generic.Some(0.0),
		BoolVal:     generic.Some(false),
		BoolIntVal:  generic.Some(false),
		BoolFlagVal: generic.Some(false),
	}

	var buf [512]byte
	out := string(dto.AppendFormData(buf[:0]))

	checks := []string{
		"i=0",
		"i64=-9223372036854775808",
		"u64=18446744073709551615",
		"f64=0",
		"b=false",
		"b_int=0",
	}
	for _, c := range checks {
		if !contains(out, c) {
			t.Errorf("missing expected boundary chunk %q in output: %s", c, out)
		}
	}
	// Flag format should be omitted when false
	if contains(out, "b_flag") {
		t.Fatalf("b_flag should be omitted when false, got: %s", out)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = dto.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for boundary values, got %v", allocs)
	}
}

func TestDTO_QueryEscapingSpecialChars(t *testing.T) {
	dto := &AllPrimitivesDTO{
		StrVal: generic.Some("alpha & beta = gamma ? # 100% / test+case"),
	}

	var buf [512]byte
	out := string(dto.AppendFormData(buf[:0]))
	expected := "str=alpha+%26+beta+%3D+gamma+%3F+%23+100%25+%2F+test%2Bcase"
	if out != expected {
		t.Fatalf("query escaping mismatch:\nGot:      %s\nExpected: %s", out, expected)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = dto.AppendFormData(buf[:0])
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs for escaped string, got %v", allocs)
	}
}

func Benchmark_AppendFormData_AllPrimitives(b *testing.B) {
	dto := makeFullDTO()
	var buf [1024]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendFormData(buf[:0])
	}
}

func Benchmark_AppendQuery_AllPrimitives(b *testing.B) {
	dto := makeFullDTO()
	var buf [1024]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendQuery(buf[:0])
	}
}

func Benchmark_AppendFormData_SomeEmptyString(b *testing.B) {
	dto := &AllPrimitivesDTO{
		EmptyStr: generic.Some(""),
	}
	var buf [128]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendFormData(buf[:0])
	}
}

func Benchmark_AppendFormData_AllNone(b *testing.B) {
	dto := &AllPrimitivesDTO{}
	var buf [128]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendFormData(buf[:0])
	}
}

func Benchmark_AppendFormData_EscapedString(b *testing.B) {
	dto := &AllPrimitivesDTO{
		StrVal: generic.Some("complex & string with symbols + spaces = 100%"),
	}
	var buf [512]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendFormData(buf[:0])
	}
}

func Benchmark_AppendQuery_SomeEmptyString(b *testing.B) {
	dto := &AllPrimitivesDTO{
		EmptyStr: generic.Some(""),
	}
	var buf [128]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendQuery(buf[:0])
	}
}

func Benchmark_AppendQuery_AllNone(b *testing.B) {
	dto := &AllPrimitivesDTO{}
	var buf [128]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendQuery(buf[:0])
	}
}

func Benchmark_AppendQuery_EscapedString(b *testing.B) {
	dto := &AllPrimitivesDTO{
		StrVal: generic.Some("complex & string with symbols + spaces = 100%"),
	}
	var buf [512]byte
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dto.AppendQuery(buf[:0])
	}
}

func Benchmark_EncodeValues_AllNone(b *testing.B) {
	dto := &AllPrimitivesDTO{}
	vals := make(url.Values)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dto.EncodeValues(vals)
	}
}

func Benchmark_EncodeValues_AllPrimitives(b *testing.B) {
	dto := makeFullDTO()
	vals := make(url.Values)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dto.EncodeValues(vals)
	}
}

func contains(s, substr string) bool {
	return urlQueryContains(s, substr)
}

func urlQueryContains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
`

	testFile := filepath.Join(tmpDir, "dto_test.go")
	err = os.WriteFile(testFile, []byte(testAndBenchSrc), 0o600)
	require.NoError(t, err)

	// Run go test -v -bench=. -benchmem . in the generated module
	cmd := exec.Command("go", "test", "-v", "-bench=.", "-benchmem", ".")
	cmd.Dir = tmpDir
	outBytes, err := cmd.CombinedOutput()
	outputStr := string(outBytes)

	if err != nil {
		t.Fatalf("go test and benchmark failed in %s: %v\nOutput:\n%s", tmpDir, err, outputStr)
	}

	t.Logf("Empirical Test & Benchmark Output:\n%s", outputStr)

	// Verify all unit tests passed
	require.Contains(t, outputStr, "PASS: TestDTO_AllPrimitives_ExactWireFormat")
	require.Contains(t, outputStr, "PASS: TestDTO_AllPrimitives_ZeroAllocations")
	require.Contains(t, outputStr, "PASS: TestDTO_SomeEmptyString_Only")
	require.Contains(t, outputStr, "PASS: TestDTO_AllNone_ZeroBytes")
	require.Contains(t, outputStr, "PASS: TestDTO_BoundaryValuesAndZeroes")
	require.Contains(t, outputStr, "PASS: TestDTO_QueryEscapingSpecialChars")

	// Parse benchmark results and verify EXACTLY 0 B/op and 0 allocs/op
	benchRegex := regexp.MustCompile(`(Benchmark\w+)-\d+\s+\d+\s+([0-9.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op`)
	matches := benchRegex.FindAllStringSubmatch(outputStr, -1)
	require.NotEmpty(t, matches, "expected benchmark results in output")

	for _, m := range matches {
		benchName := m[1]
		nsOp := m[2]
		bOp, _ := strconv.Atoi(m[3])
		allocsOp, _ := strconv.Atoi(m[4])

		t.Logf("Benchmark %-45s: %8s ns/op | %3d B/op | %3d allocs/op", benchName, nsOp, bOp, allocsOp)

		if benchName == "Benchmark_EncodeValues_AllPrimitives" {
			// EncodeValues on populated fields allocates slices via stdlib url.Values.Set
			// and strings via strconv/time formatting, but must NOT call fmt.Sprint (bounded to <= 35 allocs/op).
			require.LessOrEqualf(
				t,
				allocsOp,
				35,
				"Benchmark %s allocs must be bounded by url.Values map sets and strconv formats",
				benchName,
			)
		} else {
			// All AppendFormData, AppendQuery, and EncodeValues_AllNone benchmarks must produce EXACTLY 0 B/op and 0 allocs/op.
			require.Equalf(t, 0, bOp, "Benchmark %s must produce 0 B/op", benchName)
			require.Equalf(t, 0, allocsOp, "Benchmark %s must produce 0 allocs/op", benchName)
		}
	}
}

func BenchmarkEmitter_DTO_EmitCodegen(b *testing.B) {
	src := `package benchdto

import (
	"time"
	"github.com/lemon4ksan/foundation/generic"
)

// @aoni:dto
type BenchmarkDTO struct {
	StrVal   generic.Optional[string]    ` + "`query:\"str\"`" + `
	IntVal   generic.Optional[int64]     ` + "`query:\"val\"`" + `
	TimeVal  generic.Optional[time.Time] ` + "`query:\"ts\"`" + `
}
`

	p := parser.NewParser()
	root, err := p.ParseSource("bench.go", []byte(src))
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		out, err := emitter.Emit(root)
		if err != nil || len(out) == 0 {
			b.Fatalf("Emit failed: %v", err)
		}
	}
}
