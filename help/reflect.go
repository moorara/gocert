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
	placeholderSkip = "-"

	defaultMinLen = 8
	tagSecret     = "secret"
	tagDefault    = "default"

	promptTemplate        = "%s (type: %s):"
	promptDefaultTemplate = "%s (type: %s, default: %s):"
)

type (
	askFunc func(query string) (string, error)
)

func getPrompt(name, typeHint, defaultTag string) string {
	if defaultTag == "" {
		return fmt.Sprintf(promptTemplate, name, typeHint)
	} else {
		return fmt.Sprintf(promptDefaultTemplate, name, typeHint, defaultTag)
	}
}

func updateSkipList(skipList *[]string, field, value string) bool {
	if skipList != nil && *skipList != nil && value == placeholderSkip {
		*skipList = append(*skipList, field)
		return true
	}
	return false
}

func secretOK(pass string, minLen int) bool {
	if len(pass) < minLen {
		return false
	}

	return true
}

func getAskFunc(secretTag, defaultTag string, ui cli.Ui) askFunc {
	secretTagOpts := strings.Split(secretTag, ",")
	secretObligation := secretTagOpts[0]

	if !util.IsStringIn(secretObligation, "required", "optional") {
		// Ask function for non-secret values
		return func(query string) (string, error) {
			val, err := ui.Ask(query)
			if err != nil {
				return "", err
			}

			if val == "" {
				val = defaultTag
			}
			return val, nil
		}
	}

	minLen := defaultMinLen
	if len(secretTagOpts) > 1 {
		secretMinLen := secretTagOpts[1]
		n, err := strconv.ParseInt(secretMinLen, 10, 32)
		if err == nil {
			minLen = int(n)
		}
	}

	// Ask funtion for seccret values
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

func askForStructV(v reflect.Value, tagKey string, ignoreOmitted bool, skipList *[]string, ui cli.Ui) error {
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

		// Check if the field set to be skipped
		fullName := t.Name() + "." + name
		if skipList != nil && *skipList != nil && util.IsStringIn(fullName, *skipList...) {
			continue
		}

		secretTag := tField.Tag.Get(tagSecret)
		defaultTag := tField.Tag.Get(tagDefault)
		ask := getAskFunc(secretTag, defaultTag, ui)

		var err error
		var str string
		var b bool
		var i int64
		var f float64

		if kind == reflect.Struct {
			if err = askForStructV(vField, tagKey, ignoreOmitted, skipList, ui); err != nil {
				return err
			}
		} else if kind == reflect.Bool && vField.Bool() == false {
			if str, err = ask(getPrompt(name, "boolean", defaultTag)); err != nil {
				return err
			}
			if skipped := updateSkipList(skipList, fullName, str); !skipped {
				if b, err = strconv.ParseBool(str); err == nil {
					vField.SetBool(b)
				}
			}
		} else if kind == reflect.Int && vField.Int() == 0 {
			if str, err = ask(getPrompt(name, "integer number", defaultTag)); err != nil {
				return err
			}
			if skipped := updateSkipList(skipList, fullName, str); !skipped {
				if i, err = strconv.ParseInt(str, 10, 32); err == nil {
					vField.SetInt(i)
				}
			}
		} else if kind == reflect.Int64 && vField.Int() == 0 {
			if str, err = ask(getPrompt(name, "integer number", defaultTag)); err != nil {
				return err
			}
			if skipped := updateSkipList(skipList, fullName, str); !skipped {
				if i, err = strconv.ParseInt(str, 10, 64); err == nil {
					vField.SetInt(i)
				}
			}
		} else if kind == reflect.Float32 && vField.Float() == 0 {
			if str, err = ask(getPrompt(name, "real number", defaultTag)); err != nil {
				return err
			}
			if skipped := updateSkipList(skipList, fullName, str); !skipped {
				if f, err = strconv.ParseFloat(str, 32); err == nil {
					vField.SetFloat(f)
				}
			}
		} else if kind == reflect.Float64 && vField.Float() == 0 {
			if str, err = ask(getPrompt(name, "real number", defaultTag)); err != nil {
				return err
			}
			if skipped := updateSkipList(skipList, fullName, str); !skipped {
				if f, err = strconv.ParseFloat(str, 64); err == nil {
					vField.SetFloat(f)
				}
			}
		} else if kind == reflect.String && vField.String() == "" {
			if str, err = ask(getPrompt(name, "string", defaultTag)); err != nil {
				return err
			}
			if skipped := updateSkipList(skipList, fullName, str); !skipped {
				if str != "" {
					vField.SetString(str)
				}
			}
		} else if kind == reflect.Slice && vField.Len() == 0 {
			sliceType := reflect.TypeOf(value).Elem()
			sliceKind := sliceType.Kind()
			if sliceKind == reflect.Int {
				if str, err = ask(getPrompt(name, "integer numbers", defaultTag)); err != nil {
					return err
				}
				if skipped := updateSkipList(skipList, fullName, str); !skipped && str != "" {
					intSlice := toIntSlice(str)
					vField.Set(reflect.ValueOf(intSlice))
				}
			} else if sliceKind == reflect.Int64 {
				if str, err = ask(getPrompt(name, "integer numbers", defaultTag)); err != nil {
					return err
				}
				if skipped := updateSkipList(skipList, fullName, str); !skipped && str != "" {
					int64Slice := toInt64Slice(str)
					vField.Set(reflect.ValueOf(int64Slice))
				}
			} else if sliceKind == reflect.Float32 {
				if str, err = ask(getPrompt(name, "real numbers", defaultTag)); err != nil {
					return err
				}
				if skipped := updateSkipList(skipList, fullName, str); !skipped && str != "" {
					float32Slice := toFloat32Slice(str)
					vField.Set(reflect.ValueOf(float32Slice))
				}
			} else if sliceKind == reflect.Float64 {
				if str, err = ask(getPrompt(name, "real numbers", defaultTag)); err != nil {
					return err
				}
				if skipped := updateSkipList(skipList, fullName, str); !skipped && str != "" {
					float64Slice := toFloat64Slice(str)
					vField.Set(reflect.ValueOf(float64Slice))
				}
			} else if sliceKind == reflect.String {
				if str, err = ask(getPrompt(name, "string list", defaultTag)); err != nil {
					return err
				}
				if skipped := updateSkipList(skipList, fullName, str); !skipped && str != "" {
					slice := strings.Split(str, ",")
					vField.Set(reflect.ValueOf(slice))
				}
			} else if sliceType.String() == "net.IP" {
				if str, err = ask(getPrompt(name, "string list", defaultTag)); err != nil {
					return err
				}
				if skipped := updateSkipList(skipList, fullName, str); !skipped && str != "" {
					netIPSlice := toNetIPSlice(str)
					vField.Set(reflect.ValueOf(netIPSlice))
				}
			}
		}
	}

	return nil
}

// AskForStruct reads values from stdin for empty fields of a struct
func AskForStruct(target interface{}, tagKey string, ignoreOmitted bool, skipList *[]string, ui cli.Ui) error {
	// Get into top-level struct
	v := reflect.ValueOf(target).Elem()
	return askForStructV(v, tagKey, ignoreOmitted, skipList, ui)
}
