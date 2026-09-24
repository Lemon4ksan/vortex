// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package emitter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/lemon4ksan/vortex/pkg/ir"
)

func emitStructDTO(buf *bytes.Buffer, tracker *ImportTracker, s *ir.StructIR) {
	if !s.GenValueEncoder {
		return
	}

	tracker.Add("net/url")

	fmt.Fprintf(buf, "func (r *%s) AppendFormData(dst []byte) []byte {\n", s.Name)
	buf.WriteString("\tif r == nil {\n\t\treturn dst\n\t}\n\n")

	for _, f := range s.Fields {
		emitFieldFormData(buf, tracker, f)
	}

	buf.WriteString("\n\treturn dst\n}\n\n")

	fmt.Fprintf(buf, "func (r *%s) AppendQuery(dst []byte) []byte {\n", s.Name)
	buf.WriteString("\treturn r.AppendFormData(dst)\n}\n\n")

	// Also emit EncodeValues for url.Values interoperability
	fmt.Fprintf(buf, "func (r *%s) EncodeValues(vals url.Values) {\n", s.Name)
	buf.WriteString("\tif r == nil {\n\t\treturn\n\t}\n")

	for _, f := range s.Fields {
		emitFieldEncodeValues(buf, tracker, f)
	}

	buf.WriteString("}\n\n")

	if strings.Contains(buf.String(), "appendQueryEscape(") &&
		!strings.Contains(buf.String(), "func appendQueryEscape(") {
		emitDTOHelpers(buf)
	}
}

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

func emitFieldFormData(buf *bytes.Buffer, tracker *ImportTracker, f *ir.FieldIR) {
	if innerType, isOpt := unwrapOptionalType(f); isOpt {
		emitOptionalFieldFormData(buf, tracker, f, innerType)
		return
	}

	switch f.Type.Name {
	case "string":
		fmt.Fprintf(buf, "\tif r.%s != \"\" {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		fmt.Fprintf(buf, "\t\tdst = appendQueryEscape(dst, r.%s)\n", f.GoName)
		buf.WriteString("\t}\n")

	case "int", "int64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif r.%s != 0 {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		fmt.Fprintf(buf, "\t\tdst = strconv.AppendInt(dst, int64(r.%s), 10)\n", f.GoName)
		buf.WriteString("\t}\n")

	case "uint", "uint32", "uint64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif r.%s != 0 {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		fmt.Fprintf(buf, "\t\tdst = strconv.AppendUint(dst, uint64(r.%s), 10)\n", f.GoName)
		buf.WriteString("\t}\n")

	case "time.Time":
		tracker.Add("time")
		fmt.Fprintf(buf, "\tif !r.%s.IsZero() {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		buf.WriteString("\t\tvar timeBuf [32]byte\n")
		fmt.Fprintf(buf, "\t\ttimeBytes := r.%s.AppendFormat(timeBuf[:0], time.RFC3339)\n", f.GoName)
		buf.WriteString("\t\tfor _, c := range timeBytes {\n")
		buf.WriteString("\t\t\tif c == ':' {\n")
		buf.WriteString("\t\t\t\tdst = append(dst, \"%3A\"...)\n")
		buf.WriteString("\t\t\t} else if c == '+' {\n")
		buf.WriteString("\t\t\t\tdst = append(dst, \"%2B\"...)\n")
		buf.WriteString("\t\t\t} else {\n")
		buf.WriteString("\t\t\t\tdst = append(dst, c)\n")
		buf.WriteString("\t\t\t}\n")
		buf.WriteString("\t\t}\n")
		buf.WriteString("\t}\n")

	case "bool":
		fmt.Fprintf(buf, "\tif r.%s {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")

		switch f.Formatter {
		case ir.FormatBoolInt:
			fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=1")
		case ir.FormatBoolFlag:
			fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName)
		default:
			fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=true")
		}

		buf.WriteString("\t}\n")

	case "[]int", "[]int64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tfor _, v := range r.%s {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		fmt.Fprintf(buf, "\t\tdst = strconv.AppendInt(dst, int64(v), 10)\n")
		buf.WriteString("\t}\n")

	case "[]string":
		fmt.Fprintf(buf, "\tfor _, v := range r.%s {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		fmt.Fprintf(buf, "\t\tdst = appendQueryEscape(dst, v)\n")
		buf.WriteString("\t}\n")

	case "any", "interface{}":
		tracker.Add("fmt")
		fmt.Fprintf(buf, "\tif r.%s != nil {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		fmt.Fprintf(buf, "\t\tdst = appendQueryEscape(dst, fmt.Sprint(r.%s))\n", f.GoName)
		buf.WriteString("\t}\n")

	case "float32", "float64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif r.%s != 0 {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		fmt.Fprintf(buf, "\t\tdst = strconv.AppendFloat(dst, float64(r.%s), 'f', -1, 64)\n", f.GoName)
		buf.WriteString("\t}\n")

	case "values.Int64String", "values.Uint64String", "values.Float64String", "values.BoolInt":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif r.%s != 0 {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		fmt.Fprintf(buf, "\t\tdst = strconv.AppendInt(dst, int64(r.%s), 10)\n", f.GoName)
		buf.WriteString("\t}\n")

	default:
		if strings.HasPrefix(f.Type.Name, "[]") || strings.HasPrefix(f.Type.Name, "map[") {
			return
		}

		if f.Type.IsPointer || strings.HasPrefix(f.Type.Name, "*") {
			tracker.Add("fmt")
			fmt.Fprintf(buf, "\tif r.%s != nil {\n", f.GoName)
			buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
			fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
			fmt.Fprintf(buf, "\t\tdst = appendQueryEscape(dst, fmt.Sprint(r.%s))\n", f.GoName)
			buf.WriteString("\t}\n")
		} else {
			tracker.Add("fmt")
			fmt.Fprintf(buf, "\tif strVal := fmt.Sprint(r.%s); strVal != \"\" && strVal != \"0\" {\n", f.GoName)
			buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
			fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
			fmt.Fprintf(buf, "\t\tdst = appendQueryEscape(dst, strVal)\n")
			buf.WriteString("\t}\n")
		}
	}
}

func emitOptionalFieldFormData(buf *bytes.Buffer, tracker *ImportTracker, f *ir.FieldIR, innerType string) {
	switch innerType {
	case "string":
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		buf.WriteString("\t\tif optVal != \"\" {\n")
		buf.WriteString("\t\t\tdst = appendQueryEscape(dst, optVal)\n")
		buf.WriteString("\t\t}\n")
		buf.WriteString("\t}\n")

	case "int", "int8", "int16", "int32", "int64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		buf.WriteString("\t\tdst = strconv.AppendInt(dst, int64(optVal), 10)\n")
		buf.WriteString("\t}\n")

	case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "byte":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		buf.WriteString("\t\tdst = strconv.AppendUint(dst, uint64(optVal), 10)\n")
		buf.WriteString("\t}\n")

	case "float32", "float64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		buf.WriteString("\t\tdst = strconv.AppendFloat(dst, float64(optVal), 'f', -1, 64)\n")
		buf.WriteString("\t}\n")

	case "bool":
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		switch f.Formatter {
		case ir.FormatBoolInt:
			buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
			buf.WriteString("\t\tif optVal {\n")
			fmt.Fprintf(buf, "\t\t\tdst = append(dst, %q...)\n", f.WireName+"=1")
			buf.WriteString("\t\t} else {\n")
			fmt.Fprintf(buf, "\t\t\tdst = append(dst, %q...)\n", f.WireName+"=0")
			buf.WriteString("\t\t}\n")
		case ir.FormatBoolFlag:
			buf.WriteString("\t\tif optVal {\n")
			buf.WriteString("\t\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
			fmt.Fprintf(buf, "\t\t\tdst = append(dst, %q...)\n", f.WireName)
			buf.WriteString("\t\t}\n")
		default:
			buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
			buf.WriteString("\t\tif optVal {\n")
			fmt.Fprintf(buf, "\t\t\tdst = append(dst, %q...)\n", f.WireName+"=true")
			buf.WriteString("\t\t} else {\n")
			fmt.Fprintf(buf, "\t\t\tdst = append(dst, %q...)\n", f.WireName+"=false")
			buf.WriteString("\t\t}\n")
		}

		buf.WriteString("\t}\n")

	case "time.Time":
		tracker.Add("time")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok && !optVal.IsZero() {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		buf.WriteString("\t\tvar timeBuf [32]byte\n")
		buf.WriteString("\t\ttimeBytes := optVal.AppendFormat(timeBuf[:0], time.RFC3339)\n")
		buf.WriteString("\t\tfor _, c := range timeBytes {\n")
		buf.WriteString("\t\t\tif c == ':' {\n")
		buf.WriteString("\t\t\t\tdst = append(dst, \"%3A\"...)\n")
		buf.WriteString("\t\t\t} else if c == '+' {\n")
		buf.WriteString("\t\t\t\tdst = append(dst, \"%2B\"...)\n")
		buf.WriteString("\t\t\t} else {\n")
		buf.WriteString("\t\t\t\tdst = append(dst, c)\n")
		buf.WriteString("\t\t\t}\n")
		buf.WriteString("\t\t}\n")
		buf.WriteString("\t}\n")

	case "values.Int64String", "values.Uint64String", "values.Float64String", "values.BoolInt":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		buf.WriteString("\t\tdst = strconv.AppendInt(dst, int64(optVal), 10)\n")
		buf.WriteString("\t}\n")

	case "[]int", "[]int64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		buf.WriteString("\t\tfor _, v := range optVal {\n")
		buf.WriteString("\t\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		buf.WriteString("\t\tdst = strconv.AppendInt(dst, int64(v), 10)\n")
		buf.WriteString("\t\t}\n")
		buf.WriteString("\t}\n")

	case "[]string":
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		buf.WriteString("\t\tfor _, v := range optVal {\n")
		buf.WriteString("\t\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		buf.WriteString("\t\tdst = appendQueryEscape(dst, v)\n")
		buf.WriteString("\t\t}\n")
		buf.WriteString("\t}\n")

	default:
		if strings.HasPrefix(innerType, "*") {
			tracker.Add("fmt")
			fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok && optVal != nil {\n", f.GoName)
			buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
			fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
			fmt.Fprintf(buf, "\t\tdst = appendQueryEscape(dst, fmt.Sprint(optVal))\n")
			buf.WriteString("\t}\n")
			return
		}

		tracker.Add("fmt")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		buf.WriteString("\t\tif len(dst) > 0 { dst = append(dst, '&') }\n")
		fmt.Fprintf(buf, "\t\tdst = append(dst, %q...)\n", f.WireName+"=")
		fmt.Fprintf(buf, "\t\tdst = appendQueryEscape(dst, fmt.Sprint(optVal))\n")
		buf.WriteString("\t}\n")
	}
}

func emitFieldEncodeValues(buf *bytes.Buffer, tracker *ImportTracker, f *ir.FieldIR) {
	if innerType, isOpt := unwrapOptionalType(f); isOpt {
		emitOptionalFieldEncodeValues(buf, tracker, f, innerType)
		return
	}

	switch f.Type.Name {
	case "string":
		fmt.Fprintf(buf, "\tif r.%s != \"\" {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, r.%s)\n", f.WireName, f.GoName)
		buf.WriteString("\t}\n")

	case "int", "int64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif r.%s != 0 {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, strconv.FormatInt(int64(r.%s), 10))\n", f.WireName, f.GoName)
		buf.WriteString("\t}\n")

	case "uint", "uint32", "uint64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif r.%s != 0 {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, strconv.FormatUint(uint64(r.%s), 10))\n", f.WireName, f.GoName)
		buf.WriteString("\t}\n")

	case "float32", "float64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif r.%s != 0 {\n", f.GoName)
		fmt.Fprintf(
			buf,
			"\t\tvals.Set(%q, strconv.FormatFloat(float64(r.%s), 'f', -1, 64))\n",
			f.WireName,
			f.GoName,
		)
		buf.WriteString("\t}\n")

	case "bool":
		fmt.Fprintf(buf, "\tif r.%s {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, \"true\")\n", f.WireName)
		buf.WriteString("\t}\n")

	case "time.Time":
		tracker.Add("time")
		fmt.Fprintf(buf, "\tif !r.%s.IsZero() {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, r.%s.Format(time.RFC3339))\n", f.WireName, f.GoName)
		buf.WriteString("\t}\n")

	case "[]int", "[]int64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tfor _, v := range r.%s {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Add(%q, strconv.FormatInt(int64(v), 10))\n", f.WireName)
		buf.WriteString("\t}\n")

	case "[]string":
		fmt.Fprintf(buf, "\tfor _, v := range r.%s {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Add(%q, v)\n", f.WireName)
		buf.WriteString("\t}\n")

	case "any", "interface{}":
		tracker.Add("fmt")
		fmt.Fprintf(buf, "\tif r.%s != nil {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, fmt.Sprint(r.%s))\n", f.WireName, f.GoName)
		buf.WriteString("\t}\n")

	case "values.Int64String", "values.Uint64String", "values.Float64String", "values.BoolInt":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif r.%s != 0 {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, strconv.FormatInt(int64(r.%s), 10))\n", f.WireName, f.GoName)
		buf.WriteString("\t}\n")

	default:
		if strings.HasPrefix(f.Type.Name, "[]") || strings.HasPrefix(f.Type.Name, "map[") {
			return
		}

		if f.Type.IsPointer || strings.HasPrefix(f.Type.Name, "*") {
			tracker.Add("fmt")
			fmt.Fprintf(buf, "\tif r.%s != nil {\n", f.GoName)
			fmt.Fprintf(buf, "\t\tvals.Set(%q, fmt.Sprint(r.%s))\n", f.WireName, f.GoName)
			buf.WriteString("\t}\n")
		} else {
			tracker.Add("fmt")
			fmt.Fprintf(buf, "\tif strVal := fmt.Sprint(r.%s); strVal != \"\" && strVal != \"0\" {\n", f.GoName)
			fmt.Fprintf(buf, "\t\tvals.Set(%q, strVal)\n", f.WireName)
			buf.WriteString("\t}\n")
		}
	}
}

func emitOptionalFieldEncodeValues(buf *bytes.Buffer, tracker *ImportTracker, f *ir.FieldIR, innerType string) {
	switch innerType {
	case "string":
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, optVal)\n", f.WireName)
		buf.WriteString("\t}\n")

	case "int", "int8", "int16", "int32", "int64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, strconv.FormatInt(int64(optVal), 10))\n", f.WireName)
		buf.WriteString("\t}\n")

	case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "byte":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, strconv.FormatUint(uint64(optVal), 10))\n", f.WireName)
		buf.WriteString("\t}\n")

	case "float32", "float64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		fmt.Fprintf(
			buf,
			"\t\tvals.Set(%q, strconv.FormatFloat(float64(optVal), 'f', -1, 64))\n",
			f.WireName,
		)
		buf.WriteString("\t}\n")

	case "bool":
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		switch f.Formatter {
		case ir.FormatBoolInt:
			buf.WriteString("\t\tif optVal {\n")
			fmt.Fprintf(buf, "\t\t\tvals.Set(%q, \"1\")\n", f.WireName)
			buf.WriteString("\t\t} else {\n")
			fmt.Fprintf(buf, "\t\t\tvals.Set(%q, \"0\")\n", f.WireName)
			buf.WriteString("\t\t}\n")
		default:
			buf.WriteString("\t\tif optVal {\n")
			fmt.Fprintf(buf, "\t\t\tvals.Set(%q, \"true\")\n", f.WireName)
			buf.WriteString("\t\t} else {\n")
			fmt.Fprintf(buf, "\t\t\tvals.Set(%q, \"false\")\n", f.WireName)
			buf.WriteString("\t\t}\n")
		}
		buf.WriteString("\t}\n")

	case "time.Time":
		tracker.Add("time")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok && !optVal.IsZero() {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, optVal.Format(time.RFC3339))\n", f.WireName)
		buf.WriteString("\t}\n")

	case "values.Int64String", "values.Uint64String", "values.Float64String", "values.BoolInt":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, strconv.FormatInt(int64(optVal), 10))\n", f.WireName)
		buf.WriteString("\t}\n")

	case "[]int", "[]int64":
		tracker.Add("strconv")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		buf.WriteString("\t\tfor _, v := range optVal {\n")
		fmt.Fprintf(buf, "\t\t\tvals.Add(%q, strconv.FormatInt(int64(v), 10))\n", f.WireName)
		buf.WriteString("\t\t}\n")
		buf.WriteString("\t}\n")

	case "[]string":
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		buf.WriteString("\t\tfor _, v := range optVal {\n")
		fmt.Fprintf(buf, "\t\t\tvals.Add(%q, v)\n", f.WireName)
		buf.WriteString("\t\t}\n")
		buf.WriteString("\t}\n")

	default:
		if strings.HasPrefix(innerType, "*") {
			tracker.Add("fmt")
			fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok && optVal != nil {\n", f.GoName)
			fmt.Fprintf(buf, "\t\tvals.Set(%q, fmt.Sprint(optVal))\n", f.WireName)
			buf.WriteString("\t}\n")
			return
		}

		tracker.Add("fmt")
		fmt.Fprintf(buf, "\tif optVal, ok := r.%s.Value(); ok {\n", f.GoName)
		fmt.Fprintf(buf, "\t\tvals.Set(%q, fmt.Sprint(optVal))\n", f.WireName)
		buf.WriteString("\t}\n")
	}
}

// emitDTOHelpers emits the zero-allocation query escaping helper into the emitted Go source code.
func emitDTOHelpers(buf *bytes.Buffer) {
	buf.WriteString(`func appendQueryEscape(dst []byte, s string) []byte {
	const hexUpper = "0123456789ABCDEF"
	for i := 0; i < len(s); i++ {
		c := s[i]
		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~' {
			dst = append(dst, c)
		} else if c == ' ' {
			dst = append(dst, '+')
		} else {
			dst = append(dst, '%', hexUpper[c>>4], hexUpper[c&15])
		}
	}
	return dst
}

`)
}
