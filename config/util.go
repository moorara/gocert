package config

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
)

const (
	tagSecret   = "secret"
	askTemplate = "[%s] %s:"
)

func secretOK(password string) bool {
	if len(password) < 6 {
		return false
	}

	return true
}

func getAskSecret(ui cli.Ui) func(query string) (string, error) {
	return func(query string) (string, error) {
		secret, err := ui.AskSecret(query)
		if err != nil || !secretOK(secret) {
			ui.Error("Secret is not valid (ignored).")
			return "", errors.New("secret is not valid")
		}

		confirm, err := ui.AskSecret("CONFIRM " + query)
		if err != nil || secret != confirm {
			ui.Error("Secrets not matching (ignored).")
			return "", errors.New("secrets not matching")
		}

		return secret, nil
	}
}

func toIntSlice(list string) []int {
	slice := strings.Split(list, ",")
	intSlice := make([]int, len(slice))

	for i, str := range slice {
		n, err := strconv.ParseInt(str, 10, 32)
		if err == nil {
			intSlice[i] = int(n)
		}
	}

	return intSlice
}

func toInt64Slice(list string) []int64 {
	slice := strings.Split(list, ",")
	int64Slice := make([]int64, len(slice))

	for i, str := range slice {
		n, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			int64Slice[i] = n
		}
	}

	return int64Slice
}

func toFloat32Slice(list string) []float32 {
	slice := strings.Split(list, ",")
	float32Slice := make([]float32, len(slice))

	for i, str := range slice {
		n, err := strconv.ParseFloat(str, 32)
		if err == nil {
			float32Slice[i] = float32(n)
		}
	}

	return float32Slice
}

func toFloat64Slice(list string) []float64 {
	slice := strings.Split(list, ",")
	float64Slice := make([]float64, len(slice))

	for i, str := range slice {
		n, err := strconv.ParseFloat(str, 64)
		if err == nil {
			float64Slice[i] = n
		}
	}

	return float64Slice
}

func fillInV(v reflect.Value, tagKey string, ignoreOmitted bool, ui cli.Ui) {
	// v: reflect.Value --> v.Kind()
	t := v.Type() // reflect.Type --> t.Kind(), t.Name()

	// Iterate over struct fields
	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i) // reflect.Value --> vField.Kind(), vField.Interface(), vField.Type().Name(), vField.Type().Kind()
		tField := t.Field(i) // reflect.StructField --> tField.Name, tField.Type.Name(), tField.Type.Kind(), tField.Tag.Get(tag)

		// Skip unexported fields
		if !vField.CanSet() {
			continue
		}

		name := tField.Name
		kind := vField.Kind()
		value := vField.Interface()
		tag := tField.Tag.Get(tagKey)
		tagOpts := strings.Split(tag, ",")

		if ignoreOmitted && tagOpts[0] == "-" {
			continue
		}

		ask := ui.Ask
		if tField.Tag.Get(tagSecret) == "true" {
			ask = getAskSecret(ui)
		}

		// fmt.Printf("--> dealing with %+v\n", name)

		if kind == reflect.Struct {
			fillInV(vField, tagKey, ignoreOmitted, ui)
		} else if kind == reflect.Bool && vField.Bool() == false {
			str, err := ask(fmt.Sprintf(askTemplate, "boolean", name))
			if err == nil {
				b, err := strconv.ParseBool(str)
				if err == nil {
					vField.SetBool(b)
				}
			}
		} else if kind == reflect.Int && vField.Int() == 0 {
			str, err := ask(fmt.Sprintf(askTemplate, "number", name))
			if err == nil {
				n, err := strconv.ParseInt(str, 10, 32)
				if err == nil {
					vField.SetInt(n)
				}
			}
		} else if kind == reflect.Int64 && vField.Int() == 0 {
			str, err := ask(fmt.Sprintf(askTemplate, "number", name))
			if err == nil {
				n, err := strconv.ParseInt(str, 10, 64)
				if err == nil {
					vField.SetInt(n)
				}
			}
		} else if kind == reflect.Float32 && vField.Float() == 0 {
			str, err := ask(fmt.Sprintf(askTemplate, "number", name))
			if err == nil {
				n, err := strconv.ParseFloat(str, 32)
				if err == nil {
					vField.SetFloat(n)
				}
			}
		} else if kind == reflect.Float64 && vField.Float() == 0 {
			str, err := ask(fmt.Sprintf(askTemplate, "number", name))
			if err == nil {
				n, err := strconv.ParseFloat(str, 64)
				if err == nil {
					vField.SetFloat(n)
				}
			}
		} else if kind == reflect.String && vField.String() == "" {
			str, err := ask(fmt.Sprintf(askTemplate, "string", name))
			if err == nil && str != "" {
				vField.SetString(str)
			}
		} else if kind == reflect.Slice && vField.Len() == 0 {
			sliceKind := reflect.TypeOf(value).Elem().Kind()
			if sliceKind == reflect.Int {
				list, err := ask(fmt.Sprintf(askTemplate, "number list", name))
				if err == nil && list != "" {
					intSlice := toIntSlice(list)
					vField.Set(reflect.ValueOf(intSlice))
				}
			} else if sliceKind == reflect.Int64 {
				list, err := ask(fmt.Sprintf(askTemplate, "number list", name))
				if err == nil && list != "" {
					int64Slice := toInt64Slice(list)
					vField.Set(reflect.ValueOf(int64Slice))
				}
			} else if sliceKind == reflect.Float32 {
				list, err := ask(fmt.Sprintf(askTemplate, "number list", name))
				if err == nil && list != "" {
					float32Slice := toFloat32Slice(list)
					vField.Set(reflect.ValueOf(float32Slice))
				}
			} else if sliceKind == reflect.Float64 {
				list, err := ask(fmt.Sprintf(askTemplate, "number list", name))
				if err == nil && list != "" {
					float64Slice := toFloat64Slice(list)
					vField.Set(reflect.ValueOf(float64Slice))
				}
			} else if sliceKind == reflect.String {
				list, err := ask(fmt.Sprintf(askTemplate, "string list", name))
				if err == nil && list != "" {
					slice := strings.Split(list, ",")
					vField.Set(reflect.ValueOf(slice))
				}
			}
		}
	}
}

func fillIn(target interface{}, tagKey string, ignoreOmitted bool, ui cli.Ui) {
	// Get into top-level struct
	v := reflect.ValueOf(target).Elem()
	fillInV(v, tagKey, ignoreOmitted, ui)
}
