// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser_test

import (
	"testing"

	"github.com/lemon4ksan/foundation/testing/require"

	"github.com/lemon4ksan/vortex/pkg/ir"
	"github.com/lemon4ksan/vortex/pkg/parser"
)

func TestParser_Adversarial_VariedGenericsResolution(t *testing.T) {
	src := `package testpkg

import (
	"context"
	"time"

	"github.com/lemon4ksan/foundation/generic"
)

// @aoni:dto
type ComplexGenericsDTO struct {
	OptString       generic.Optional[string]                         ` + "`query:\"opt_string\"`" + `
	OptIntUnqual    Optional[int]                                    ` + "`query:\"opt_int\"`" + `
	OptTime         Optional[time.Time]                              ` + "`query:\"opt_time\"`" + `
	OptPtr          generic.Optional[*string]                        ` + "`query:\"opt_ptr\"`" + `
	OptSlice        generic.Optional[[]int]                          ` + "`query:\"opt_slice\"`" + `
	OptMap          generic.Optional[map[string]any]                 ` + "`query:\"opt_map\"`" + `
	OptNested       generic.Optional[generic.Optional[string]]       ` + "`query:\"opt_nested\"`" + `
	Pair            generic.Pair[string, int]                        ` + "`query:\"pair\"`" + `
	Triple          generic.Triple[int, string, bool]                ` + "`query:\"triple\"`" + `
	Quad            generic.Quad[string, int, bool, float64]         ` + "`query:\"quad\"`" + `
	ParenOuter      (generic.Optional[string])                       ` + "`query:\"paren_outer\"`" + `
	ParenInner      generic.Optional[(time.Time)]                    ` + "`query:\"paren_inner\"`" + `
	OptChan         generic.Optional[<-chan int]                     ` + "`query:\"opt_chan\"`" + `
	OptFunc         generic.Optional[func(int) string]               ` + "`query:\"opt_func\"`" + `
}

// @aoni:service
type AdversarialService interface {
	// @get "/items"
	GetItems(
		ctx context.Context,
		filter *ComplexGenericsDTO,
		optStr generic.Optional[string],
		optInt Optional[int],
		optTime Optional[time.Time],
		optPtr *generic.Optional[string],
	) error

	// @delete "/items/{id}"
	DeleteItem(
		ctx context.Context,
		id string,
		purge generic.Optional[bool],
	) error

	// @post "/items"
	// @query dryRun
	CreateItem(
		ctx context.Context,
		req *ComplexGenericsDTO,
		dryRun generic.Optional[bool],
	) error

	// @post "/items/raw"
	PostRawOptional(
		ctx context.Context,
		payload generic.Optional[string],
	) error
}
`

	p := parser.NewParser()
	root, err := p.ParseSource("adversarial.go", []byte(src))
	require.NoError(t, err)
	require.NotNil(t, root)
	require.Len(t, root.Structs, 1)

	s := root.Structs[0]
	require.Equal(t, "ComplexGenericsDTO", s.Name)

	fieldMap := make(map[string]*ir.FieldIR)
	for _, f := range s.Fields {
		fieldMap[f.GoName] = f
	}

	// 1. OptString: generic.Optional[string]
	fOptStr := fieldMap["OptString"]
	require.NotNil(t, fOptStr)
	require.Equal(t, "generic.Optional[string]", fOptStr.Type.Name)
	require.Equal(t, "string", fOptStr.Type.ElemType)
	require.True(t, fOptStr.Type.IsCustomType)

	// 2. OptIntUnqual: Optional[int]
	fOptInt := fieldMap["OptIntUnqual"]
	require.NotNil(t, fOptInt)
	require.Equal(t, "Optional[int]", fOptInt.Type.Name)
	require.Equal(t, "int", fOptInt.Type.ElemType)
	require.True(t, fOptInt.Type.IsCustomType)

	// 3. OptTime: Optional[time.Time]
	fOptTime := fieldMap["OptTime"]
	require.NotNil(t, fOptTime)
	require.Equal(t, "Optional[time.Time]", fOptTime.Type.Name)
	require.Equal(t, "time.Time", fOptTime.Type.ElemType)
	require.True(t, fOptTime.Type.IsCustomType)

	// 4. OptPtr: generic.Optional[*string]
	fOptPtr := fieldMap["OptPtr"]
	require.NotNil(t, fOptPtr)
	require.Equal(t, "generic.Optional[*string]", fOptPtr.Type.Name)
	require.Equal(t, "*string", fOptPtr.Type.ElemType)
	require.True(t, fOptPtr.Type.IsCustomType)

	// 5. OptSlice: generic.Optional[[]int]
	fOptSlice := fieldMap["OptSlice"]
	require.NotNil(t, fOptSlice)
	require.Equal(t, "generic.Optional[[]int]", fOptSlice.Type.Name)
	require.Equal(t, "[]int", fOptSlice.Type.ElemType)
	require.True(t, fOptSlice.Type.IsCustomType)

	// 6. OptMap: generic.Optional[map[string]any]
	fOptMap := fieldMap["OptMap"]
	require.NotNil(t, fOptMap)
	require.Equal(t, "generic.Optional[map[string]any]", fOptMap.Type.Name)
	require.Equal(t, "map[string]any", fOptMap.Type.ElemType)
	require.True(t, fOptMap.Type.IsCustomType)

	// 7. OptNested: generic.Optional[generic.Optional[string]]
	fOptNested := fieldMap["OptNested"]
	require.NotNil(t, fOptNested)
	require.Equal(t, "generic.Optional[generic.Optional[string]]", fOptNested.Type.Name)
	require.Equal(t, "generic.Optional[string]", fOptNested.Type.ElemType)
	require.True(t, fOptNested.Type.IsCustomType)

	// 8. Pair: generic.Pair[string, int] (IndexListExpr with 2 type args)
	fPair := fieldMap["Pair"]
	require.NotNil(t, fPair)
	require.Equal(t, "generic.Pair[string, int]", fPair.Type.Name)
	require.Equal(t, "string, int", fPair.Type.ElemType)
	require.True(t, fPair.Type.IsCustomType)

	// 9. Triple: generic.Triple[int, string, bool] (IndexListExpr with 3 type args)
	fTriple := fieldMap["Triple"]
	require.NotNil(t, fTriple)
	require.Equal(t, "generic.Triple[int, string, bool]", fTriple.Type.Name)
	require.Equal(t, "int, string, bool", fTriple.Type.ElemType)
	require.True(t, fTriple.Type.IsCustomType)

	// 10. Quad: generic.Quad[string, int, bool, float64] (IndexListExpr with 4 type args)
	fQuad := fieldMap["Quad"]
	require.NotNil(t, fQuad)
	require.Equal(t, "generic.Quad[string, int, bool, float64]", fQuad.Type.Name)
	require.Equal(t, "string, int, bool, float64", fQuad.Type.ElemType)
	require.True(t, fQuad.Type.IsCustomType)

	// 11. ParenOuter: (generic.Optional[string])
	fParenOuter := fieldMap["ParenOuter"]
	require.NotNil(t, fParenOuter)
	require.Equal(t, "generic.Optional[string]", fParenOuter.Type.Name)
	require.Equal(t, "string", fParenOuter.Type.ElemType)
	require.True(t, fParenOuter.Type.IsCustomType)

	// 12. ParenInner: generic.Optional[(time.Time)]
	fParenInner := fieldMap["ParenInner"]
	require.NotNil(t, fParenInner)
	require.Equal(t, "generic.Optional[time.Time]", fParenInner.Type.Name)
	require.Equal(t, "time.Time", fParenInner.Type.ElemType)
	require.True(t, fParenInner.Type.IsCustomType)

	// 13. OptChan: generic.Optional[<-chan int]
	fOptChan := fieldMap["OptChan"]
	require.NotNil(t, fOptChan)
	require.Equal(t, "generic.Optional[<-chan int]", fOptChan.Type.Name)
	require.Equal(t, "<-chan int", fOptChan.Type.ElemType)
	require.True(t, fOptChan.Type.IsCustomType)

	// 14. OptFunc: generic.Optional[func(int) string]
	fOptFunc := fieldMap["OptFunc"]
	require.NotNil(t, fOptFunc)
	require.Equal(t, "generic.Optional[func(int) string]", fOptFunc.Type.Name)
	require.Equal(t, "func(int) string", fOptFunc.Type.ElemType)
	require.True(t, fOptFunc.Type.IsCustomType)

	// Verify Services and Methods
	require.Len(t, root.Services, 1)
	svc := root.Services[0]
	require.Equal(t, "AdversarialService", svc.Name)
	require.Len(t, svc.Methods, 4)

	// Method 0: GetItems
	mGet := svc.Methods[0]
	require.Equal(t, "GetItems", mGet.Name)
	require.Len(t, mGet.Params, 6)
	require.Equal(t, ir.LocContext, mGet.Params[0].Location)
	require.Equal(t, ir.LocQueryStruct, mGet.Params[1].Location) // filter is query struct
	require.Equal(t, ir.LocQuery, mGet.Params[2].Location)       // optStr is query param!
	require.Equal(t, ir.LocQuery, mGet.Params[3].Location)       // optInt is query param!
	require.Equal(t, ir.LocQuery, mGet.Params[4].Location)       // optTime is query param!
	require.Equal(t, ir.LocQuery, mGet.Params[5].Location)       // optPtr is query param!

	// Method 1: DeleteItem
	mDel := svc.Methods[1]
	require.Equal(t, "DeleteItem", mDel.Name)
	require.Len(t, mDel.Params, 3)
	require.Equal(t, ir.LocContext, mDel.Params[0].Location)
	require.Equal(t, ir.LocPath, mDel.Params[1].Location)  // id
	require.Equal(t, ir.LocQuery, mDel.Params[2].Location) // purge optional

	// Method 2: CreateItem (with explicit @query directive)
	mPost := svc.Methods[2]
	require.Equal(t, "CreateItem", mPost.Name)
	require.Len(t, mPost.Params, 3)
	require.Equal(t, ir.LocContext, mPost.Params[0].Location)
	require.Equal(t, ir.LocBody, mPost.Params[1].Location)  // req is body
	require.Equal(t, ir.LocQuery, mPost.Params[2].Location) // dryRun has @query directive -> LocQuery

	// Method 3: PostRawOptional (without directive, on POST custom type becomes LocBody)
	mRaw := svc.Methods[3]
	require.Equal(t, "PostRawOptional", mRaw.Name)
	require.Len(t, mRaw.Params, 2)
	require.Equal(t, ir.LocContext, mRaw.Params[0].Location)
	require.Equal(t, ir.LocBody, mRaw.Params[1].Location) // payload becomes LocBody
}
