package main

import (
	"fmt"
	"reflect"
	"strings"
)

func DumpStruct(v interface{}) {
	dumpStruct(reflect.ValueOf(v), 0)
}

func dumpStruct(val reflect.Value, indent int) {
	if !val.IsValid() {
		fmt.Println(strings.Repeat("  ", indent) + "<invalid value>")
		return
	}

	if val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
		if val.IsNil() {
			fmt.Println(strings.Repeat("  ", indent) + "<nil>")
			return
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		t := val.Type()
		for i := 0; i < val.NumField(); i++ {
			fieldVal := val.Field(i)
			fieldType := t.Field(i)
			if !fieldType.IsExported() {
				continue // skip unexported fields
			}
			prefix := fmt.Sprintf("%s%s (%s): ", strings.Repeat("  ", indent), fieldType.Name, fieldVal.Type())
			if isSimpleKind(fieldVal.Kind()) {
				fmt.Printf("%s%v\n", prefix, fieldVal.Interface())
			} else {
				fmt.Println(prefix)
				dumpStruct(fieldVal, indent+1)
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			fmt.Printf("%s[%d]: ", strings.Repeat("  ", indent), i)
			elem := val.Index(i)
			if isSimpleKind(elem.Kind()) {
				fmt.Printf("%v\n", elem.Interface())
			} else {
				fmt.Println()
				dumpStruct(elem, indent+1)
			}
		}
	case reflect.Map:
		iter := val.MapRange()
		for iter.Next() {
			key := iter.Key()
			elem := iter.Value()
			fmt.Printf("%s%v: ", strings.Repeat("  ", indent), key.Interface())
			if isSimpleKind(elem.Kind()) {
				fmt.Printf("%v\n", elem.Interface())
			} else {
				fmt.Println()
				dumpStruct(elem, indent+1)
			}
		}
	default:
		fmt.Printf("%s%v\n", strings.Repeat("  ", indent), val.Interface())
	}
}

func isSimpleKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}
