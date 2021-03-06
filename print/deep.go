package print

import (
	"bytes"
	"fmt"
	"reflect"
)

// Deep is like printing with the %#v formatter of fmt, but it prints
// pointer fields recursively.
func Deep(x interface{}) string {
	b := &bytes.Buffer{}
	deep(b, reflect.ValueOf(x))
	return b.String()
}

func deep(b *bytes.Buffer, v reflect.Value) {
	i := v.Interface()
	t := v.Type()

	// GoStringer
	if g, ok := i.(fmt.GoStringer); ok {
		b.WriteString(g.GoString())
		return
	}

	// nil
	switch v.Kind() {
	case reflect.Interface, reflect.Map, reflect.Slice, reflect.Ptr:
		if v.IsNil() {
			b.WriteString("nil")
			return
		}
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		// Composite kinds
		b.WriteString(t.String())
		b.WriteRune('{')
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				if i > 0 {
					b.WriteString(", ")
				}
				deep(b, v.Index(i))
			}
		case reflect.Map:
			keys := v.MapKeys()
			for i, k := range keys {
				if i > 0 {
					b.WriteString(", ")
				}
				deep(b, k)
				b.WriteString(": ")
				deep(b, v.MapIndex(k))
			}
		case reflect.Struct:
			for i := 0; i < t.NumField(); i++ {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(t.Field(i).Name)
				b.WriteString(": ")
				deep(b, v.Field(i))
			}
		}
		b.WriteRune('}')
	case reflect.Ptr:
		b.WriteRune('&')
		deep(b, reflect.Indirect(v))
		return
	case reflect.Interface:
		deep(b, v.Elem())
		return
	default:
		fmt.Fprintf(b, "%#v", i)
		return
	}
}
