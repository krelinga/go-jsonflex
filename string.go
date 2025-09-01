package jsonflex

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// String returns a string representation of the given value.
//
// Only supports the following types:
// - All types rooted in Object.
// - JSON basic types (bool, float64, int32, string)
// - Slices of any other supported type.
func String(v any) string {
	return innerString(reflect.ValueOf(v))
}

func indent(in string) string {
	return strings.ReplaceAll(in, "\n", "\n  ")
}

func innerString(v reflect.Value) string {
	sb := strings.Builder{}
	switch v.Kind() {
	case reflect.Map:
		sb.WriteString("{\n")
		methods := make([]int, v.NumMethod())
		for methodNum := range v.Type().NumMethod() {
			methods[methodNum] = methodNum
		}
		slices.SortFunc(methods, func(a, b int) int {
			aName := v.Type().Method(a).Name
			bName := v.Type().Method(b).Name
			return strings.Compare(aName, bName)
		})
		for _, methodNum := range methods {
			method := v.Type().Method(methodNum)
			if method.Type.NumIn() != 1 {
				continue
			}
			if method.Type.NumOut() != 2 {
				continue
			}
			if method.Type.Out(1) != reflect.TypeFor[error]() {
				continue
			}
			outs := method.Func.Call([]reflect.Value{v})
			var outString string
			if outs[1].IsNil() {
				outString = indent(innerString(outs[0]))
			} else if errors.Is(outs[1].Interface().(error), ErrNullValue) {
				outString = "null"
			} else if errors.Is(outs[1].Interface().(error), ErrFieldNotFound) {
				continue
			} else {
				outString = fmt.Sprintf("error: %s", outs[1].Interface())
			}
			sb.WriteString(fmt.Sprintf("  %s: %s,\n", method.Name, outString))
		}
		sb.WriteString("}")
	case reflect.Slice:
		sb.WriteString("[\n")
		for i := 0; i < v.Len(); i++ {
			sb.WriteString(fmt.Sprintf("  %d: %s,\n", i, indent(innerString(v.Index(i)))))
		}
		sb.WriteString("]")
	case reflect.Bool, reflect.Int32, reflect.Float64:
		sb.WriteString(fmt.Sprintf("%v", v.Interface()))
	case reflect.String:
		sb.WriteString(fmt.Sprintf("%q", v.Interface()))
	default:
		sb.WriteString(fmt.Sprintf("unsupported type: %s", v.Type()))
	}
	return sb.String()
}
