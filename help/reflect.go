package help

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
)

const (
	tagSecret     = "secret"
	defaultMinLen = 8

	// AskTemplate is used when asking for a value
	AskTemplate = "%s (%s):"
)

type askFunc func(query string) (string, error)

func secretOK(pass string, minLen int) bool {
	if len(pass) < minLen {
		return false
	}

	return true
}

func getAskFunc(tag string, ui cli.Ui) askFunc {
	tagOpts := strings.Split(tag, ",")
	obligation := tagOpts[0]

	if !util.IsStringIn(obligation, "required", "optional") {
		return ui.Ask
	}

	minLen := defaultMinLen
	if len(tagOpts) > 1 {
		n, err := strconv.ParseInt(tagOpts[1], 10, 32)
		if err == nil {
			minLen = int(n)
		}
	}

	return func(query string) (string, error) {
		secret, err := ui.AskSecret(query)
		if err != nil || !secretOK(secret, minLen) {
			ui.Error("Secret not valid.")
			return "", errors.New("secret not valid")
		}

		confirm, err := ui.AskSecret("CONFIRM " + query)
		if err != nil || secret != confirm {
			ui.Error("Secrets not matched.")
			return "", errors.New("secrets not matched")
		}

		return secret, nil
	}
}

func toIntSlice(list string) []int {
	vals := strings.Split(list, ",")
	intVals := make([]int, len(vals))

	for i, str := range vals {
		n, err := strconv.ParseInt(str, 10, 32)
		if err == nil {
			intVals[i] = int(n)
		}
	}

	return intVals
}

func toInt64Slice(list string) []int64 {
	vals := strings.Split(list, ",")
	int64Vals := make([]int64, len(vals))

	for i, str := range vals {
		n, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			int64Vals[i] = n
		}
	}

	return int64Vals
}

func toFloat32Slice(list string) []float32 {
	vals := strings.Split(list, ",")
	float32Vals := make([]float32, len(vals))

	for i, str := range vals {
		n, err := strconv.ParseFloat(str, 32)
		if err == nil {
			float32Vals[i] = float32(n)
		}
	}

	return float32Vals
}

func toFloat64Slice(list string) []float64 {
	vals := strings.Split(list, ",")
	float64Vals := make([]float64, len(vals))

	for i, str := range vals {
		n, err := strconv.ParseFloat(str, 64)
		if err == nil {
			float64Vals[i] = n
		}
	}

	return float64Vals
}

func toNetIPSlice(list string) []net.IP {
	vals := strings.Split(list, ",")
	netIPVals := make([]net.IP, len(vals))

	for i, str := range vals {
		netIPVals[i] = net.ParseIP(str)
	}

	return netIPVals
}

func askForStructV(v reflect.Value, tagKey string, ignoreOmitted bool, ui cli.Ui) error {
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

		// Check for omit empty tag
		tag := tField.Tag.Get(tagKey)
		if ignoreOmitted && strings.HasPrefix(tag, "-") {
			continue
		}

		secretTag := tField.Tag.Get(tagSecret)
		ask := getAskFunc(secretTag, ui)

		if kind == reflect.Struct {
			err := askForStructV(vField, tagKey, ignoreOmitted, ui)
			if err != nil {
				return err
			}
		} else if kind == reflect.Bool && vField.Bool() == false {
			str, err := ask(fmt.Sprintf(AskTemplate, name, "boolean"))
			if err != nil {
				return err
			}
			b, err := strconv.ParseBool(str)
			if err == nil {
				vField.SetBool(b)
			}
		} else if kind == reflect.Int && vField.Int() == 0 {
			str, err := ask(fmt.Sprintf(AskTemplate, name, "integer number"))
			if err != nil {
				return err
			}
			n, err := strconv.ParseInt(str, 10, 32)
			if err == nil {
				vField.SetInt(n)
			}
		} else if kind == reflect.Int64 && vField.Int() == 0 {
			str, err := ask(fmt.Sprintf(AskTemplate, name, "integer number"))
			if err != nil {
				return err
			}
			n, err := strconv.ParseInt(str, 10, 64)
			if err == nil {
				vField.SetInt(n)
			}
		} else if kind == reflect.Float32 && vField.Float() == 0 {
			str, err := ask(fmt.Sprintf(AskTemplate, name, "real number"))
			if err != nil {
				return err
			}
			n, err := strconv.ParseFloat(str, 32)
			if err == nil {
				vField.SetFloat(n)
			}
		} else if kind == reflect.Float64 && vField.Float() == 0 {
			str, err := ask(fmt.Sprintf(AskTemplate, name, "real number"))
			if err != nil {
				return err
			}
			n, err := strconv.ParseFloat(str, 64)
			if err == nil {
				vField.SetFloat(n)
			}
		} else if kind == reflect.String && vField.String() == "" {
			str, err := ask(fmt.Sprintf(AskTemplate, name, "string"))
			if err != nil {
				return err
			}
			if str != "" {
				vField.SetString(str)
			}
		} else if kind == reflect.Slice && vField.Len() == 0 {
			sliceType := reflect.TypeOf(value).Elem()
			sliceKind := sliceType.Kind()
			if sliceKind == reflect.Int {
				list, err := ask(fmt.Sprintf(AskTemplate, name, "integer numbers"))
				if err != nil {
					return err
				}
				if list != "" {
					intSlice := toIntSlice(list)
					vField.Set(reflect.ValueOf(intSlice))
				}
			} else if sliceKind == reflect.Int64 {
				list, err := ask(fmt.Sprintf(AskTemplate, name, "integer numbers"))
				if err != nil {
					return err
				}
				if list != "" {
					int64Slice := toInt64Slice(list)
					vField.Set(reflect.ValueOf(int64Slice))
				}
			} else if sliceKind == reflect.Float32 {
				list, err := ask(fmt.Sprintf(AskTemplate, name, "real numbers"))
				if err != nil {
					return err
				}
				if list != "" {
					float32Slice := toFloat32Slice(list)
					vField.Set(reflect.ValueOf(float32Slice))
				}
			} else if sliceKind == reflect.Float64 {
				list, err := ask(fmt.Sprintf(AskTemplate, name, "real numbers"))
				if err != nil {
					return err
				}
				if list != "" {
					float64Slice := toFloat64Slice(list)
					vField.Set(reflect.ValueOf(float64Slice))
				}
			} else if sliceKind == reflect.String {
				list, err := ask(fmt.Sprintf(AskTemplate, name, "string list"))
				if err != nil {
					return err
				}
				if list != "" {
					slice := strings.Split(list, ",")
					vField.Set(reflect.ValueOf(slice))
				}
			} else if sliceType.String() == "net.IP" {
				list, err := ask(fmt.Sprintf(AskTemplate, name, "string list"))
				if err != nil {
					return err
				}
				if list != "" {
					netIPSlice := toNetIPSlice(list)
					vField.Set(reflect.ValueOf(netIPSlice))
				}
			}
		}
	}

	return nil
}

// AskForStruct reads values from stdin for empty fields of a struct
func AskForStruct(target interface{}, tagKey string, ignoreOmitted bool, ui cli.Ui) error {
	// Get into top-level struct
	v := reflect.ValueOf(target).Elem()
	return askForStructV(v, tagKey, ignoreOmitted, ui)
}
