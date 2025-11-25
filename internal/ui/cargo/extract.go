package cargo

import (
	"reflect"
	"strings"
)

// Extract - assign fields from stored struct using key-pointer pairs
// keys are case-insensitive, missing fields are ignored, type mismatches are ignored
func Extract[T any](args ...any) {
	key := typeKey[T]()
	extractKeysFrom[T](key, args...)
}

// ExtractByKey - assign fields from stored struct using key-pointer pairs
func ExtractByKey[T any](key string, args ...any) {
	extractKeysFrom[T](key, args...)
}

// ExtractInto - fill entire struct from stored struct, case-insensitive
// missing fields are zeroed, type mismatches are ignored
func ExtractInto[T any](outPtr any) {
	if outPtr == nil {
		return
	}

	ptrVal := reflect.ValueOf(outPtr)
	if ptrVal.Kind() != reflect.Ptr || ptrVal.IsNil() {
		return
	}

	elem := ptrVal.Elem()
	if elem.Kind() != reflect.Struct {
		return
	}

	key := typeKey[T]()
	value, ok := Get[T](key)
	if !ok {
		var zero T
		value = zero
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			var zero T
			rv = reflect.ValueOf(zero)
		} else {
			rv = rv.Elem()
		}
	}
	if rv.Kind() != reflect.Struct {
		return
	}

	// Собираем только поля верхнего уровня
	fields := make(map[string]reflect.Value)
	collectTopLevelFields(rv, rv.Type(), fields)

	ciLookup := make(map[string]string)
	for name := range fields {
		ciLookup[strings.ToLower(name)] = name
	}

	outType := elem.Type()
	for i := 0; i < outType.NumField(); i++ {
		f := outType.Field(i)
		if !f.IsExported() {
			continue
		}

		outField := elem.Field(i)
		if !outField.CanSet() {
			continue
		}

		lowerName := strings.ToLower(f.Name)
		structFieldName, exists := ciLookup[lowerName]
		if !exists {
			outField.Set(reflect.Zero(outField.Type()))
			continue
		}

		fieldVal := fields[structFieldName]
		if !fieldVal.Type().AssignableTo(outField.Type()) {
			outField.Set(reflect.Zero(outField.Type()))
			continue
		}

		outField.Set(fieldVal)
	}
}

func extractKeysFrom[T any](key string, args ...any) {
	if len(args)%2 != 0 {
		return
	}

	value, ok := Get[T](key)
	if !ok {
		var zero T
		value = zero
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			var zero T
			rv = reflect.ValueOf(zero)
		} else {
			rv = rv.Elem()
		}
	}
	if rv.Kind() != reflect.Struct {
		return
	}

	fields := make(map[string]reflect.Value)
	collectTopLevelFields(rv, rv.Type(), fields)

	ciLookup := make(map[string]string)
	for name := range fields {
		ciLookup[strings.ToLower(name)] = name
	}

	for i := 0; i < len(args); i += 2 {
		rawKey := args[i]
		rawPtr := args[i+1]

		keyStr, ok := rawKey.(string)
		if !ok || rawPtr == nil {
			continue
		}

		ptrVal := reflect.ValueOf(rawPtr)
		if ptrVal.Kind() != reflect.Ptr || ptrVal.IsNil() {
			continue
		}

		elem := ptrVal.Elem()
		fieldName, exists := ciLookup[strings.ToLower(keyStr)]
		if !exists {
			elem.Set(reflect.Zero(elem.Type()))
			continue
		}

		fieldVal := fields[fieldName]
		if !fieldVal.Type().AssignableTo(elem.Type()) {
			elem.Set(reflect.Zero(elem.Type()))
			continue
		}

		elem.Set(fieldVal)
	}
}

func collectTopLevelFields(rv reflect.Value, rt reflect.Type, out map[string]reflect.Value) {
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		fv := rv.Field(i)

		if !f.IsExported() {
			continue
		}

		if f.Anonymous {
			if fv.Kind() == reflect.Struct {
				collectTopLevelFields(fv, fv.Type(), out)
			} else if fv.Kind() == reflect.Ptr && fv.Type().Elem().Kind() == reflect.Struct && !fv.IsNil() {
				collectTopLevelFields(fv.Elem(), fv.Type().Elem(), out)
			}
			continue
		}

		out[f.Name] = fv
	}
}
