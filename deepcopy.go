// Package deepcopy knows how to deep copy a value, using reflection.
package deepcopy

import (
	"errors"
	"fmt"
	"reflect"
)

// Copy returns a deepcopy of the specified object.
// Unexported fields of a struct are ignored and will not be copied.
// The types unsafe.Pointer and uintptr are not supported and they will cause a panic.
// A channel will point to original channel.
// Error will always be nil, and it is here just for compatibility with copystructure package.
func Copy(o interface{}) (interface{}, error) {
	return copyr(reflect.ValueOf(o)).Interface(), nil
}

// copyr deep copies a reflect value.
// We intentionally specify all supported types, so we panic for all unsupported.
func copyr(ov reflect.Value) reflect.Value {
	if !ov.IsValid() {
		panic(errors.New("invalid value"))
	}
	if t := ov.Type(); t.PkgPath() == "time" && t.Name() == "Time" {
		return copyTime(ov)
	}
	switch ov.Kind() {
	case reflect.Struct:
		return copyStruct(ov)
	case reflect.Ptr:
		return copyPointer(ov)
	case reflect.Slice:
		return copySlice(ov)
	case reflect.Map:
		return copyMap(ov)
	case reflect.Interface:
		return copyInterface(ov)
	case reflect.Array:
		return copyArray(ov)
	case reflect.Int, reflect.String, reflect.Int64, reflect.Float64, reflect.Bool, reflect.Uint, reflect.Uint64,
		reflect.Func, reflect.Chan, reflect.Float32,
		reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Complex64, reflect.Complex128,
		reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return ov
	}
	panic(fmt.Sprintf("unsupported type: %s", ov.Kind()))
}

func copyTime(ov reflect.Value) reflect.Value {
	return ov
}

func copyInterface(ov reflect.Value) reflect.Value {
	if ov.IsNil() {
		return ov
	}
	oc := reflect.New(ov.Type())
	oc.Elem().Set(copyr(ov.Elem()))
	return oc.Elem()
}

func copyPointer(ov reflect.Value) reflect.Value {
	if ov.IsNil() {
		return ov
	}
	oc := reflect.New(ov.Type().Elem())
	oc.Elem().Set(copyr(ov.Elem()))
	return oc
}

func copyStruct(ov reflect.Value) reflect.Value {
	oc := reflect.New(ov.Type()).Elem()
	for i := 0; i < ov.NumField(); i++ {
		fv := ov.Field(i)
		// we do not set unexported fields as runtime does not allow it
		// also, runtime does not allow assigning a zero value, in case of pointers
		if !fv.IsZero() && fv.CanInterface() {
			oc.Field(i).Set(copyr(ov.Field(i)))
		}
	}
	return oc
}

func copySlice(ov reflect.Value) reflect.Value {
	if ov.IsNil() {
		return ov
	}
	oc := reflect.MakeSlice(ov.Type(), 0, ov.Cap())
	for i := 0; i < ov.Len(); i++ {
		oc = reflect.Append(oc, copyr(ov.Index(i)))
	}
	return oc
}

func copyArray(ov reflect.Value) reflect.Value {
	array := reflect.New(ov.Type()).Elem()
	slice := array.Slice3(0, 0, array.Len())
	for i := 0; i < ov.Len(); i++ {
		slice = reflect.Append(slice, copyr(ov.Index(i)))
	}
	return array
}

func copyMap(ov reflect.Value) reflect.Value {
	if ov.IsNil() {
		return ov
	}
	oc := reflect.MakeMapWithSize(ov.Type(), ov.Len())
	iter := ov.MapRange()
	for iter.Next() {
		oc.SetMapIndex(copyr(iter.Key()), copyr(iter.Value()))
	}
	return oc
}

func isPrimitive(ot reflect.Type) bool {
	switch ot.Kind() {
	case reflect.Int, reflect.String, reflect.Int64, reflect.Float64, reflect.Bool, reflect.Uint, reflect.Uint64,
		reflect.Func, reflect.Chan, reflect.Float32,
		reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Complex64, reflect.Complex128,
		reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return true
	case reflect.Struct:
		for i := 0; i < ot.NumField(); i++ {
			if !isPrimitive(ot.Field(i).Type) {
				return false
			}
		}
		return true
	}
	return false
}
