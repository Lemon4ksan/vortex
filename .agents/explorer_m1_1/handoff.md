# Handoff Report: IR Type Resolution in `pkg/parser/binder.go:extractGoType`

**Agent**: `explorer_m1_1` (Archetype: Explorer)  
**Task**: Milestone 1 (Zero-Allocation generic.Optional[T] & Empty Field DTO Codegen) — IR Generic Type Resolution Investigation  
**Working Directory**: `d:/CodingProjects/vortex/.agents/explorer_m1_1/`  
**Date**: 2026-09-22T14:43:00Z  

---

## 1. Observation

Direct code inspection of `d:/CodingProjects/vortex` reveals the following concrete facts:

### 1.1 `pkg/parser/binder.go:extractGoType` (Lines 943–1035)
In `pkg/parser/binder.go`:
```go
943: func (p *Parser) extractGoType(expr ast.Expr) ir.GoTypeIR {
944: 	var (
945: 		buf    bytes.Buffer
946: 		goType ir.GoTypeIR
947: 	)
948: 
949: 	switch t := expr.(type) {
950: 	case *ast.Ident:
951: 		goType.Name = t.Name
952: 		if isPrimitive(t.Name) {
953: 			goType.Underlying = t.Name
954: 		} else {
955: 			goType.IsCustomType = true
956: 		}
957: 
958: 	case *ast.StarExpr:
959: 		elem := p.extractGoType(t.X)
960: 		goType = elem
961: 		goType.IsPointer = true
962: 		goType.Name = "*" + elem.Name
963: 	case *ast.ArrayType:
964: 		elem := p.extractGoType(t.Elt)
965: 		goType = elem
966: 		goType.IsSlice = true
967: 		goType.ElemType = elem.Name
968: 		goType.Name = "[]" + elem.Name
969: 
970: 	case *ast.MapType:
971: 		k := p.extractGoType(t.Key)
972: 		v := p.extractGoType(t.Value)
973: 		goType.IsMap = true
974: 		goType.KeyType = k.Name
975: 		goType.ElemType = v.Name
976: 		goType.Name = fmt.Sprintf("map[%s]%s", k.Name, v.Name)
977: 
978: 	case *ast.SelectorExpr:
979: 		pkg := p.extractGoType(t.X)
980: 		goType.Package = pkg.Name
981: 		goType.Name = fmt.Sprintf("%s.%s", pkg.Name, t.Sel.Name)
982: 		goType.IsCustomType = true
983: 	case *ast.ChanType:
984: 		elem := p.extractGoType(t.Value)
985: 		goType = elem
986: 		goType.IsChannel = true
987: 		goType.Name = "<-chan " + elem.Name
988: 	case *ast.Ellipsis:
989: 		elem := p.extractGoType(t.Elt)
990: 		goType = elem
991: 		goType.IsVariadic = true
992: 		goType.Name = "..." + elem.Name
993: 	case *ast.FuncType:
...
1029: 	default:
1030: 		_ = buf
1031: 		goType.Name = "any"
1032: 	}
1033: 
1034: 	return goType
1035: }
```
**Observation Details**:
- `extractGoType` has NO case branches for `*ast.IndexExpr` or `*ast.IndexListExpr`.
- Any generic instantiation, including `generic.Optional[string]`, `Optional[int]`, `generic.Optional[time.Time]`, or `generic.Pair[string, int]`, evaluates as unrecognized and falls through to `default:` at line 1029.
- This results in `goType.Name = "any"`, `goType.ElemType = ""`, and `goType.IsCustomType = false`.

### 1.2 `pkg/ir/types.go:GoTypeIR` (Lines 624–636)
`pkg/ir/types.go`:
```go
624: type GoTypeIR struct {
625: 	Name         string
626: 	Package      string
627: 	IsPointer    bool
628: 	IsSlice      bool
629: 	IsMap        bool
630: 	IsChannel    bool
631: 	IsVariadic   bool
632: 	IsCustomType bool
633: 	Underlying   string
634: 	ElemType     string
635: 	KeyType      string
636: }
```
**Observation Details**:
- `GoTypeIR` already contains the `ElemType string` field, which is currently populated only for slices (`*ast.ArrayType`), maps (`*ast.MapType`), and callback return types (`*ast.FuncType`).
- No IR schema modification is required to store `ElemType` or `IsCustomType` for generic types.

### 1.3 Downstream Consumptions in `pkg/emitter/dto.go`
`pkg/emitter/dto.go`:
- Lines 112–118:
  ```go
  case "any", "interface{}":
      tracker.Add("fmt")
      fmt.Fprintf(buf, "\tif r.%s != nil {\n", f.GoName)
      ...
      fmt.Fprintf(buf, "\t\tdst = append(dst, url.QueryEscape(fmt.Sprint(r.%s))...)\n", f.GoName)
      buf.WriteString("\t}\n")
  ```
- Lines 137–146:
  ```go
  default:
      if strings.HasPrefix(f.Type.Name, "generic.Optional[") || strings.HasPrefix(f.Type.Name, "Optional[") {
          tracker.Add("fmt")
          fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
          ...
      }
  ```
- Lines 210–214 and 221–228:
  ```go
  case "any", "interface{}":
      tracker.Add("fmt")
      fmt.Fprintf(buf, "\tif r.%s != nil {\n", f.GoName)
      fmt.Fprintf(buf, "\t\tvals.Set(%q, fmt.Sprint(r.%s))\n", f.WireName, f.GoName)
      buf.WriteString("\t}\n")
  default:
      if strings.HasPrefix(f.Type.Name, "generic.Optional[") || strings.HasPrefix(f.Type.Name, "Optional[") {
          ...
      }
  ```
**Observation Details**:
- Because `extractGoType` produces `Name: "any"`, `emitFieldFormData` and `emitFieldEncodeValues` match `case "any", "interface{}:"` at line 112/210 and **never** reach the `strings.HasPrefix(f.Type.Name, "generic.Optional[")` logic in `default:`!
- Even when reaching `default:`, `dto.go` currently invokes `fmt.Sprint(optVal)`, which causes runtime heap allocations, boxed interface conversions (`convT2E`), reflection overhead, and fails the zero-allocation requirement.

### 1.4 Downstream Consumptions Across Other Packages
- `pkg/emitter/helpers.go:16`: `formatType(t ir.GoTypeIR) string` returns `t.Name` if non-empty, otherwise `"any"`. Because `t.Name` was `"any"`, generated method signatures erased generic types to `any`.
- `pkg/emitter/buffer_writer.go:37`: `elemType := p.GoType.ElemType`. Only used when `p.Formatter` is `ir.FormatSliceComma`, `ir.FormatSliceSpace`, or `ir.FormatSlicePipe`.
- `pkg/emitter/imports.go:278`: `if strings.Contains(bodyCode, ref+".") { t.AddNamed(imp.Alias, imp.Path) }`. When `generic.Optional` is emitted in the body code, `imports.go` automatically adds `"github.com/lemon4ksan/foundation/generic"` from `root.Imports`.
- `pkg/parser/binder.go:isDTOQueryStruct` (Lines 1080–1104):
  ```go
  1080: func isDTOQueryStruct(name string) bool {
  1081: 	if isPrimitive(name) {
  1082: 		return false
  1083: 	}
  1084: 
  1085: 	switch name {
  1086: 	case "time.Time", "Time", "time.Duration", "Duration",
  1087: 		"values.Int64String", "values.Uint64String", "values.Float64String", "values.BoolInt",
  1088: 		"uuid.UUID", "UUID", "decimal.Decimal", "Decimal", "netip.Addr", "Addr":
  1089: 		return false
  1090: 	}
  1091: 
  1092: 	if strings.HasPrefix(name, "[]") || strings.HasPrefix(name, "map[") {
  1093: 		return false
  1094: 	}
  ...
  1103: 	return true
  1104: }
  ```
  **Critical Observation**:
  `isDTOQueryStruct` is invoked at lines 374 and 377:
  `case (m.HTTPMethod == "GET" || ...) && isDTOQueryStruct(goType.Name): param.Location = ir.LocQueryStruct; param.Formatter = ir.FormatCompiledEncode`
  If `goType.Name` is `"generic.Optional[string]"`, `isDTOQueryStruct` returns `true` because it is not in the excluded list. This causes a standalone optional method parameter on a `GET` method to be misclassified as a compiled DTO struct with `.AppendQuery(qBytes)`, causing code emission failure because `generic.Optional[T]` is a container monad, not a DTO query struct.

---

## 2. Logic Chain

### 2.1 AST Representation of Generics in Go 1.18+
1. In Go's `go/ast` syntax tree:
   - Single-type generic instantiation `T[U]` (such as `generic.Optional[string]`, `Optional[int]`, `generic.Optional[*Item]`, `generic.Optional[time.Time]`) is represented by `*ast.IndexExpr`:
     - `t.X` (`ast.Expr`): The generic type identifier or selector.
       - If qualified (e.g. `generic.Optional`), `t.X` is `*ast.SelectorExpr` where `t.X.X` is `*ast.Ident{Name: "generic"}` and `t.X.Sel` is `*ast.Ident{Name: "Optional"}`.
       - If unqualified (e.g. `Optional`), `t.X` is `*ast.Ident{Name: "Optional"}`.
     - `t.Index` (`ast.Expr`): The concrete type argument expression (e.g., `*ast.Ident{Name: "string"}`, `*ast.StarExpr`, `*ast.SelectorExpr`, etc.).
   - Multi-type generic instantiation `T[U, V, ...]` (such as `generic.Pair[string, int]`, `Result[Data, error]`, or edge-case single args with trailing commas like `Optional[string,]`) is represented by `*ast.IndexListExpr`:
     - `t.X` (`ast.Expr`): The generic base type expression.
     - `t.Indices` (`[]ast.Expr`): The slice of type argument expressions.
2. In `extractGoType`, `expr.(type)` lacks cases for `*ast.IndexExpr`, `*ast.IndexListExpr`, and `*ast.ParenExpr`.
3. Therefore, any field of type `generic.Optional[T]` matches `default:`, resulting in `goType.Name = "any"`.

### 2.2 AST Transformation & IR Field Resolution
1. To preserve the generic identity and its underlying element type:
   - For `*ast.IndexExpr`:
     - Recursively resolve `base := p.extractGoType(t.X)`.
     - Recursively resolve `elem := p.extractGoType(t.Index)`.
     - Inherit base attributes (`Package`, etc.).
     - Assign `goType.ElemType = elem.Name`.
     - Assign `goType.Name = fmt.Sprintf("%s[%s]", base.Name, elem.Name)`.
     - Set `goType.IsCustomType = true`.
   - For `*ast.IndexListExpr`:
     - Recursively resolve `base := p.extractGoType(t.X)`.
     - Recursively resolve each argument `it := p.extractGoType(idx)`.
     - If `len(indices) == 1`, set `goType.ElemType = indices[0]`; if multiple, set `goType.ElemType = strings.Join(indices, ", ")`.
     - Assign `goType.Name = fmt.Sprintf("%s[%s]", base.Name, strings.Join(indices, ", "))`.
     - Set `goType.IsCustomType = true`.
   - For `*ast.ParenExpr`:
     - Return `p.extractGoType(t.X)` to unwrap parentheses transparently.
2. By setting `goType.Name = "generic.Optional[string]"` and `goType.ElemType = "string"`:
   - `f.Type.Name` retains the exact type string matching `strings.HasPrefix(..., "generic.Optional[")`.
   - `f.Type.ElemType` gives downstream emitters the exact unwrapped primitive type (`"string"`, `"int"`, `"int64"`, `"uint64"`, `"float64"`, `"bool"`, `"time.Time"`).
   - Emitters can bypass `fmt.Sprint` entirely and branch on `f.Type.ElemType` to generate `strconv.Append*` and empty-string serialization without heap allocations.

### 2.3 Side-Effect Prevention
1. **Downstream Safety**:
   - `pkg/emitter/helpers.go`: Now emits the exact signature type (e.g. `generic.Optional[string]`) instead of `any`.
   - `pkg/emitter/imports.go`: Emits `import "github.com/lemon4ksan/foundation/generic"` whenever `generic.` appears in generated code.
   - `pkg/emitter/buffer_writer.go`: Slices continue to use `ElemType` as before; generic optionals are not slices (`IsSlice` is false).
2. **`isDTOQueryStruct` Misclassification Fix**:
   - At `binder.go:1080`, `isDTOQueryStruct` decides whether a method parameter should be encoded via `param.AppendQuery(qBytes)`.
   - Because `generic.Optional[T]` is not a primitive and not explicitly excluded in `isDTOQueryStruct`, it would return `true`.
   - Excluded prefix check `if strings.HasPrefix(name, "generic.Optional[") || strings.HasPrefix(name, "Optional[") { return false }` must be added to `isDTOQueryStruct` so standalone optional parameters are treated as individual query parameters rather than query DTO structs.

---

## 3. Caveats

1. **Complex or Custom Inner Types**:
   When `T` in `generic.Optional[T]` is a custom struct or interface rather than a primitive (`string`, `int*`, `uint*`, `float*`, `bool`, `time.Time`), `ElemType` will contain the struct name (e.g. `"User"` or `"*User"`). Downstream DTO emission will fall back to JSON marshaling or `fmt.Sprint` for non-primitives. This is expected and standard Go behavior.
2. **Aliased Optional Types**:
   If an application imports `foundation/generic` with a local alias (e.g. `import opt "github.com/lemon4ksan/foundation/generic"` and uses `opt.Optional[string]`), `base.Name` will be `"opt.Optional"`, and `goType.Name` will be `"opt.Optional[string]"`. To ensure robust matching in `pkg/emitter/dto.go:unwrapOptionalType`, the check should verify `strings.HasSuffix(name, "]") && strings.Contains(name, "Optional[")`, or check `strings.HasPrefix(name, "generic.Optional[") || strings.HasPrefix(name, "Optional[")`.
3. **No Code Modifications Made in this Turn**:
   In strict compliance with the read-only exploration mandate, no source code in `pkg/` has been altered. All recommendations below are ready for immediate implementation by developer agents.

---

## 4. Conclusion

### Summary Assessment
1. **Root Cause**: The type erasure of `generic.Optional[T]` to `any` occurs exclusively in `pkg/parser/binder.go:extractGoType` (lines 943–1035) because `*ast.IndexExpr` and `*ast.IndexListExpr` fall into `default:`.
2. **Required Fields on `ir.GoTypeIR`**:
   - `Name`: `"generic.Optional[T]"` (or `"Optional[T]"`)
   - `ElemType`: `"T"` (inner type, e.g. `"string"`, `"int64"`, `"bool"`, `"time.Time"`)
   - `IsCustomType`: `true`
   - `Package`: `"generic"` (or `""` if unqualified)
3. **Downstream Compatibility**:
   - Downstream emitters in `pkg/emitter` are fully compatible and immediately benefit from having `ElemType` and concrete `Name`.
   - One critical prevention is required in `pkg/parser/binder.go:isDTOQueryStruct` (line 1080) to prevent standalone `generic.Optional[T]` method parameters from being classified as DTO query structs.

### Actionable Code Recommendations for `pkg/parser/binder.go`

#### Code Modification 1: `pkg/parser/binder.go:extractGoType` (around line 1028)

**Target Content**:
```go
		goType.Name = fmt.Sprintf("func(%s)%s", strings.Join(params, ", "), resStr)

	default:
		_ = buf
		goType.Name = "any"
	}

	return goType
```

**Replacement Content**:
```go
		goType.Name = fmt.Sprintf("func(%s)%s", strings.Join(params, ", "), resStr)

	case *ast.ParenExpr:
		return p.extractGoType(t.X)

	case *ast.IndexExpr:
		base := p.extractGoType(t.X)
		elem := p.extractGoType(t.Index)
		goType = base
		goType.ElemType = elem.Name
		goType.Name = fmt.Sprintf("%s[%s]", base.Name, elem.Name)
		goType.IsCustomType = true

	case *ast.IndexListExpr:
		base := p.extractGoType(t.X)
		indices := make([]string, 0, len(t.Indices))
		for _, idx := range t.Indices {
			it := p.extractGoType(idx)
			indices = append(indices, it.Name)
		}
		goType = base
		if len(indices) == 1 {
			goType.ElemType = indices[0]
		} else {
			goType.ElemType = strings.Join(indices, ", ")
		}
		goType.Name = fmt.Sprintf("%s[%s]", base.Name, strings.Join(indices, ", "))
		goType.IsCustomType = true

	default:
		_ = buf
		goType.Name = "any"
	}

	return goType
```

#### Code Modification 2: `pkg/parser/binder.go:isDTOQueryStruct` (around line 1080)

**Target Content**:
```go
func isDTOQueryStruct(name string) bool {
	if isPrimitive(name) {
		return false
	}

	switch name {
```

**Replacement Content**:
```go
func isDTOQueryStruct(name string) bool {
	if isPrimitive(name) {
		return false
	}

	if strings.HasPrefix(name, "generic.Optional[") || strings.HasPrefix(name, "Optional[") {
		return false
	}

	switch name {
```

#### Recommended Unit Test for `pkg/parser/parser_test.go`
Add a dedicated test case verifying generic type resolution:
```go
func TestParser_GenericOptionalFields(t *testing.T) {
	src := `package testapi

import (
	"github.com/lemon4ksan/foundation/generic"
	"time"
)

// @aoni:dto
type SearchFilter struct {
	Query    generic.Optional[string]    ` + "`query:\"q\"`" + `
	Page     generic.Optional[int]       ` + "`query:\"page\"`" + `
	Active   generic.Optional[bool]      ` + "`query:\"active\"`" + `
	Created  generic.Optional[time.Time] ` + "`query:\"created\"`" + `
	Pair     generic.Pair[string, int]   ` + "`query:\"pair\"`" + `
}
`
	p := parser.NewParser()
	root, err := p.ParseSource("filter.go", []byte(src))
	require.NoError(t, err)
	require.NotNil(t, root)
	require.Len(t, root.Structs, 1)

	s := root.Structs[0]
	require.Equal(t, "SearchFilter", s.Name)
	require.Len(t, s.Fields, 5)

	// Field 0: generic.Optional[string]
	require.Equal(t, "Query", s.Fields[0].GoName)
	require.Equal(t, "generic.Optional[string]", s.Fields[0].Type.Name)
	require.Equal(t, "string", s.Fields[0].Type.ElemType)
	require.True(t, s.Fields[0].Type.IsCustomType)

	// Field 1: generic.Optional[int]
	require.Equal(t, "Page", s.Fields[1].GoName)
	require.Equal(t, "generic.Optional[int]", s.Fields[1].Type.Name)
	require.Equal(t, "int", s.Fields[1].Type.ElemType)
	require.True(t, s.Fields[1].Type.IsCustomType)

	// Field 2: generic.Optional[bool]
	require.Equal(t, "Active", s.Fields[2].GoName)
	require.Equal(t, "generic.Optional[bool]", s.Fields[2].Type.Name)
	require.Equal(t, "bool", s.Fields[2].Type.ElemType)
	require.True(t, s.Fields[2].Type.IsCustomType)

	// Field 3: generic.Optional[time.Time]
	require.Equal(t, "Created", s.Fields[3].GoName)
	require.Equal(t, "generic.Optional[time.Time]", s.Fields[3].Type.Name)
	require.Equal(t, "time.Time", s.Fields[3].Type.ElemType)
	require.True(t, s.Fields[3].Type.IsCustomType)

	// Field 4: generic.Pair[string, int] (IndexListExpr)
	require.Equal(t, "Pair", s.Fields[4].GoName)
	require.Equal(t, "generic.Pair[string, int]", s.Fields[4].Type.Name)
	require.Equal(t, "string, int", s.Fields[4].Type.ElemType)
	require.True(t, s.Fields[4].Type.IsCustomType)
}
```

---

## 5. Verification Method

### 5.1 Verification Commands
1. Run existing parser unit tests:
   ```pwsh
   go test -v ./pkg/parser/...
   ```
   (Must pass with code 0).
2. Run existing emitter unit tests:
   ```pwsh
   go test -v ./pkg/emitter/...
   ```
   (Must pass with code 0).
3. After applying the recommended code changes:
   Run the newly added `TestParser_GenericOptionalFields`:
   ```pwsh
   go test -v -run TestParser_GenericOptionalFields ./pkg/parser/...
   ```
   (Must pass with code 0).

### 5.2 Files to Inspect
- `d:/CodingProjects/vortex/pkg/parser/binder.go`: Verify `*ast.IndexExpr`, `*ast.IndexListExpr`, and `*ast.ParenExpr` are handled in `extractGoType`.
- `d:/CodingProjects/vortex/pkg/parser/binder.go`: Verify `isDTOQueryStruct` excludes `generic.Optional[` and `Optional[`.
- `d:/CodingProjects/vortex/pkg/emitter/dto.go`: Verify downstream usage of `f.Type.ElemType` and `unwrapOptionalType`.

### 5.3 Invalidation Conditions
This investigation is invalidated if:
- Go AST representation changes (not applicable in Go 1.18+ / 1.22+ / 1.24+).
- `ir.GoTypeIR` is refactored to replace `ElemType` with a new generic parameter slice structure.
