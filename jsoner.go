package jsoner

import (
	"encoding/json"
	"reflect"
	"strings"
)

func Unmarshal(str string, target any) error {
	d := json.NewDecoder(strings.NewReader(string(str)))
	d.UseNumber()

	if err := d.Decode(target); err != nil {
		return err
	}

	replaceNumbers(target)

	return nil
}

func replaceNumbers(target any) {
	value := reflect.Indirect(reflect.ValueOf(target))

	switch value.Kind() {
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			currentValue := value.Field(i).Interface()
			switch v := currentValue.(type) {
			case json.Number:
				mayBeNum := toPrimitiveNumber(v)
				value.Field(i).Set(reflect.ValueOf(mayBeNum))
			default:
				replaceNumbers(v)
			}
		}
	case reflect.Map:
		for _, k := range value.MapKeys() {
			currentValue := value.MapIndex(k)
			switch v := currentValue.Interface().(type) {
			case json.Number:
				mayBeNum := toPrimitiveNumber(v)
				value.SetMapIndex(k, reflect.ValueOf(mayBeNum))
			default:
				replaceNumbers(v)
			}

		}

	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			currentValue := value.Index(i)
			switch v := currentValue.Interface().(type) {
			case json.Number:
				mayBeNum := toPrimitiveNumber(v)
				currentValue.Set(reflect.ValueOf(mayBeNum))
			default:
				replaceNumbers(v)
			}
		}
	case reflect.Interface:
		switch v := value.Interface().(type) {
		case json.Number:
			mayBeNum := toPrimitiveNumber(v)
			reflect.ValueOf(target).Elem().Set(reflect.ValueOf(mayBeNum))
		default:
			replaceNumbers(value.Interface())
		}
	}

}

func toPrimitiveNumber(n json.Number) any {
	if i, err := n.Int64(); err == nil {
		return int(i)
	}

	if f, err := n.Float64(); err == nil {
		return f
	}

	return n
}
