package main

import (
	"reflect"
)

func toWasmType(s any) any {
	switch reflect.TypeOf(s).Kind()  {
	case reflect.Struct:
		v := reflect.ValueOf(s)
		transformed := make(map[string]any)

		for i := 0; i < v.NumField(); i++ {
			transformed[v.Type().Field(i).Name] = toWasmType(v.Field(i).Interface())
		}
		return transformed
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		v := reflect.ValueOf(s)
		transformed := make([]any, v.Len())

		for i := 0; i < v.Len(); i++ {
			transformed[i] = toWasmType(v.Index(i).Interface())
		}
		return transformed
	case reflect.Chan:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Pointer:
		panic("The given variable contained a type which is not supported for transformation")
	default:
		return s
	}
}
