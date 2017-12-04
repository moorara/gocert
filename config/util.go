package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
)

const (
	askTemplate = "[%s] %s:"
)

// fillIn asks for input for empty fields
func fillIn(target interface{}, tagKey string, includeOmitted bool, ui cli.Ui) {
	// Get into top-level struct
	v := reflect.ValueOf(target).Elem() // reflect.Value --> v.Kind()
	t := v.Type()                       // reflect.Type --> t.Kind(), t.Name()

	// Iterate over struct fields
	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i) // reflect.Value --> vField.Kind(), vField.Interface(), vField.Type().Name(), vField.Type().Kind()
		tField := t.Field(i) // reflect.StructField --> tField.Name, tField.Type.Name(), tField.Type.Kind(), tField.Tag.Get(tag)

		kind := vField.Kind()
		value := vField.Interface()
		tag := tField.Tag.Get(tagKey)

		if includeOmitted || tag != "-" {
			if kind == reflect.Bool && vField.Bool() == false {
				str, err := ui.Ask(fmt.Sprintf(askTemplate, "boolean", tField.Name))
				if err == nil {
					b, err := strconv.ParseBool(str)
					if err == nil {
						vField.SetBool(b)
					}
				}
			} else if kind == reflect.Int && vField.Int() == 0 {
				str, err := ui.Ask(fmt.Sprintf(askTemplate, "number", tField.Name))
				if err == nil {
					n, err := strconv.ParseInt(str, 10, 32)
					if err == nil {
						vField.SetInt(n)
					}
				}
			} else if kind == reflect.Int64 && vField.Int() == 0 {
				str, err := ui.Ask(fmt.Sprintf(askTemplate, "number", tField.Name))
				if err == nil {
					n, err := strconv.ParseInt(str, 10, 64)
					if err == nil {
						vField.SetInt(n)
					}
				}
			} else if kind == reflect.String && vField.String() == "" {
				str, err := ui.Ask(fmt.Sprintf(askTemplate, "string", tField.Name))
				if err == nil && str != "" {
					vField.SetString(str)
				}
			} else if kind == reflect.Slice && vField.Len() == 0 {
				sliceKind := reflect.TypeOf(value).Elem().Kind()
				if sliceKind == reflect.String {
					list, err := ui.Ask(fmt.Sprintf(askTemplate, "string list", tField.Name))
					if err == nil && list != "" {
						slice := strings.Split(list, ",")
						vField.Set(reflect.ValueOf(slice))
					}
				}
			}
		}
	}
}
