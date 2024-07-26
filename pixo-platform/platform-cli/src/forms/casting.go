package forms

import (
	"reflect"
	"strconv"
)

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

func Int(stringOrPtr interface{}) int {
	v := reflect.ValueOf(stringOrPtr)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return 0
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Int {
		return int(v.Int())
	}

	if v.Kind() == reflect.String {
		i, err := strconv.Atoi(v.String())
		if err != nil {
			return 0
		}
		return i
	}

	return 0
}

func IntSlice(sliceOrPtr interface{}) []int {
	v := reflect.ValueOf(sliceOrPtr)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Slice {
		var result []int
		for i := 0; i < v.Len(); i++ {
			result = append(result, int(v.Index(i).Int()))
		}
		return result
	}

	return nil
}

func StringSlice(sliceOrPtr interface{}) []string {
	v := reflect.ValueOf(sliceOrPtr)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Slice {
		var result []string
		for i := 0; i < v.Len(); i++ {
			result = append(result, v.Index(i).String())
		}
		return result
	}

	return nil
}
