package forms

import "reflect"

func String(stringOrPtr interface{}) string {
	v := reflect.ValueOf(stringOrPtr)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.String {
		return v.String()
	}

	return ""
}

func Bool(boolOrPtr interface{}) bool {
	v := reflect.ValueOf(boolOrPtr)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Bool {
		return v.Bool()
	}

	return false
}
