package forms

import (
	"context"
	"fmt"
	"reflect"
)

type QuestionType int

const (
	Input QuestionType = iota
	SensitiveInput
	Confirm
	Select
	SelectID
	MultiSelect
	MultiSelectIDs
)

type Question struct {
	Type         QuestionType
	Options      []Option
	GetItemsFunc func(context.Context) (interface{}, error)
	LabelFunc    func(interface{}) string
	Answer       interface{}
	Key          string
	Prompt       string
	Optional     bool
}

func (q *Question) GetOptions(ctx context.Context) error {
	if q.Options == nil && q.GetItemsFunc != nil {
		items, err := q.GetItemsFunc(ctx)
		if err != nil {
			return err
		}

		iterableItems := getIterableItems(items)
		options := make([]Option, len(iterableItems))
		for i, item := range iterableItems {
			var label string
			if q.LabelFunc != nil {
				label = q.LabelFunc(item)
			} else {
				label, _ = getLabel(item)
			}
			val, _ := getValue(item)
			options[i] = Option{
				Label: label,
				Value: fmt.Sprint(val),
			}
		}

		q.Options = options
		return nil
	}
	return nil
}

func getID(v interface{}) (id int, found bool) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Struct {
		fieldVal := val.FieldByName("ID")
		if fieldVal.IsValid() {
			id = int(fieldVal.Int())
			found = true
		}
	}

	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		fieldVal := val.Elem().FieldByName("ID")
		if fieldVal.IsValid() {
			id = int(fieldVal.Int())
			found = true
		}
	}

	return id, found && id > 0
}

func getValue(v interface{}) (string, bool) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Struct {
		fieldVal := val.FieldByName("Value")
		if fieldVal.IsValid() {
			value := fieldVal.String()
			return value, true
		}
	}

	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		fieldVal := val.Elem().FieldByName("Value")
		if fieldVal.IsValid() {
			value := fieldVal.String()
			return value, true
		}
	}

	if id, ok := getID(v); ok {
		return fmt.Sprint(id), true
	}

	return "", false
}

func getLabel(v interface{}) (string, bool) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Struct {
		fieldVal := val.FieldByName("Label")
		if fieldVal.IsValid() {
			label := fieldVal.String()
			return label, true
		}
	}

	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		fieldVal := val.Elem().FieldByName("Label")
		if fieldVal.IsValid() {
			label := fieldVal.String()
			return label, true
		}
	}

	if val.Kind() == reflect.Struct {
		fieldVal := val.FieldByName("Name")
		if fieldVal.IsValid() {
			label := fieldVal.String()
			return label, true
		}
	}

	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		fieldVal := val.Elem().FieldByName("Name")
		if fieldVal.IsValid() {
			label := fieldVal.String()
			return label, true
		}
	}

	return "", false
}

func getIterableItems(v interface{}) []interface{} {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Slice {
		items := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			items[i] = val.Index(i).Interface()
		}
		return items
	}

	return nil
}
