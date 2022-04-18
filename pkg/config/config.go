package config

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

const tagName = "config"

// LoadFromMap loads the tagged struct fields from the given map.
func LoadFromMap(m map[string]string, out interface{}) error {
	return LoadFromMapP(m, "", out)
}

// LoadFromMapP loads the tagged struct fields from the given map with keys
// prefixed with the given prefix.
func LoadFromMapP(m map[string]string, prefix string, out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}

		args := strings.Split(tag, ",")
		if len(args) == 0 {
			continue
		}

		key := args[0]
		required := false
		if len(args) > 1 {
			required = args[1] == "required"
		}

		fullKey := prefix + key
		if val, ok := m[fullKey]; ok {
			err := convertStringAndSetField(val, v.Field(i))
			if err != nil {
				return err
			}
		} else if required {
			return errors.New("config: missing required key " + fullKey)
		}
	}

	return nil
}

func convertStringAndSetField(s string, f reflect.Value) error {
	switch f.Kind() {
	case reflect.String:
		f.SetString(s)
	case reflect.Int:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		f.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		f.SetBool(b)
	default:
		return errors.New("config: unsupported key type " + f.Kind().String())
	}
	return nil
}
