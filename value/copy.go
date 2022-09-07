package value

import "reflect"

func Cast[T any](in any) (t T) {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()

	t = in.(T)
	return
}

func Copy(from interface{}, to interface{}) {
	toElem := reflect.ValueOf(to)
	fromElem := reflect.ValueOf(from)
	if elem := reflect.ValueOf(to); elem.Kind() == reflect.Ptr {
		toElem = elem.Elem()
	} else {
		toElem = elem
	}

	if elem := reflect.ValueOf(from); elem.Kind() == reflect.Ptr {
		fromElem = elem.Elem()
	} else {
		fromElem = elem
	}

	switch toElem.Kind() {
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String:
		if fromElem.Kind() == toElem.Kind() && toElem.CanSet() {
			toElem.Set(fromElem)
		}
	case reflect.Struct:
		if fromElem.Kind() == toElem.Kind() {
			for i := 0; i < toElem.NumField(); i++ {
				fromField := fromElem.FieldByName(toElem.Type().Field(i).Name)
				toField := toElem.Field(i)
				if fromField.Kind() == toField.Kind() && toField.CanSet() {
					toField.Set(fromField)
				}
			}
		}
	case reflect.Map:
		if fromElem.Kind() == toElem.Kind() && fromElem.String() == toElem.String() && toElem.CanSet() {
			toElem.Set(fromElem)
		}
	}
}
